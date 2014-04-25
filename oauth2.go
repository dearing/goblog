package main

import (
	"encoding/json"
	//"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/goauth2/oauth"
	"github.com/gorilla/securecookie"
)

var oauth_config = oauth.Config{
	ClientId:     config.ClientID,
	ClientSecret: config.ClientSecret,
	Scope:        "",
	AuthURL:      "https://github.com/login/oauth/authorize",
	TokenURL:     "https://github.com/login/oauth/access_token",
	RedirectURL:  config.RedirectURL,
}

// The blockKey is optional, used to encrypt the cookie value -- set it to nil to not use encryption.
// If set, the length must correspond to the block size of the encryption algorithm.
// For AES, used by default, valid lengths are 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256.
// SEE: http://www.gorillatoolkit.org/pkg/securecookie

var hashKey = []byte(securecookie.GenerateRandomKey(32))
var blockKey = []byte(securecookie.GenerateRandomKey(32))
var s = securecookie.New(hashKey, blockKey)

// Start the login process if the user doesn't have a AUTH valid cookie.
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if !validateCookie(w, r) {

		// Get a new random string for our STATE, store it in a secure cookie
		// on the client and start our authentication process.
		state := uuid.New()
		setStateCookie(w, r, state)
		url := oauth_config.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)

	} else {

		// No need to reauthenticate if our present cookie is still good.
		http.Redirect(w, r, "/", http.StatusAccepted)
	}
}

// This callback contains a temp CODE and our passed STATE from Github
func callbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")   // temp CODE for authentication
	state := r.FormValue("state") // secret between us and Github

	if state != getState(w, r) {
		if config.Verbose {
			log.Println("Client's stored state did not match post data from Github. Failing login attempt.")
		}
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}

	// We attempt to exchange the temp CODE we got for a real access token...
	t := &oauth.Transport{Config: &oauth_config}
	_, e := t.Exchange(code)
	if e != nil {
		log.Println(e)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	c := t.Client()

	// Now we have a client with a valid access token.
	// We use this access token to request the current authenticated user information.
	// This will be the same as a public read to the Github API v3 however, the /user URL
	// alone only returns the authenticated user according to the token.  So this should
	// suffice as proof that the user IS a valid Github USER.

	resp, _ := c.Get("https://api.github.com/user")
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	// We have a valid Github USER but is it our ADMIN?
	var info map[string]interface{}
	if err := json.Unmarshal(contents, &info); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	// Now here we look at the login field in the JSON the server sent us for our ADMIN check.
	// TODO:  I'm not comfortable bypassing strict TYPE checking by using the wildcard interface{}
	//        which matches anything.  Create a strict type with values we care about and move on.
	login := info["login"]
	if login != config.AdminLogin {

		// USER ain't our ADMIN, move along.
		log.Println("Github user did not match admin configuration.")
		removeCookies(w, r)
		http.Redirect(w, r, "/", http.StatusUnauthorized)

	} else {

		// Set a cookie now to retain AUTH status after cleaning up our state cookie and perhaps any previous invalid cookies.
		log.Println("Admin logged in.")
		removeCookies(w, r)
		setAuthCookie(w, r, login)
		http.Redirect(w, r, "/", http.StatusAccepted)
	}

}

// Delete cookies; send client home
func logoutHander(w http.ResponseWriter, r *http.Request) {
	removeCookies(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

// SECRET page for testing only; no need to redirect to login
func secretPageHandler(w http.ResponseWriter, r *http.Request) {
	if validateCookie(w, r) {
		http.Redirect(w, r, "/", http.StatusAccepted)
	} else {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
}

// 5 minutes to keep state cookie for cross-site forgery attempts
// I could be storing this at the server but this *is* a secure cookie right?
func setStateCookie(w http.ResponseWriter, r *http.Request, state string) {
	if encoded, err := s.Encode("state", state); err == nil {
		cookie := &http.Cookie{
			Name:    "state",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(5 * time.Minute), // t+5m
		}
		http.SetCookie(w, cookie)
	}
}

// Set our AUTH secure cookie for 1 hour
// TODO: clean it up - at present we are simply dumping all user data into the cookie.
func setAuthCookie(w http.ResponseWriter, r *http.Request, login interface{}) {

	value := map[string]interface{}{
		"login": login,
		"host":  r.Host,
	}

	if encoded, err := s.Encode("auth", value); err == nil {
		cookie := &http.Cookie{
			Name:    "auth",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(time.Hour), // t+60m
		}
		http.SetCookie(w, cookie)
	}
}

// Remove all cookies we have set on this client.
func removeCookies(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookies() {

		// MaxAge < 0 means delete cookie now
		// http://golang.org/pkg/net/http/#Cookie
		v.MaxAge = -1
		http.SetCookie(w, v)
	}
}

// Just returns the state stored in cookie, or "" if nadda.
// BUG (jacob): What if github returns state:"" and we have nothing stored?
//              Then one equals the other but code:?? is still needed to move on.
//			    Still, seems a bit sloppy for my taste...
func getState(w http.ResponseWriter, r *http.Request) (state string) {

	cookie, err := r.Cookie("state")
	if err != nil {
		log.Println(err)
	}

	err = s.Decode("state", cookie.Value, &state)
	if err != nil {
		log.Println(err)
	}
	return state
}

// Returns true if the cookie the client provided checks out with configured ADMIN
// Cleanup invalid cookies if anything doesn't check out.
// TODO: clean this shit up
func validateCookie(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("auth")
	if err != nil {
		log.Println(err)
		removeCookies(w, r)
		return false
	}

	value := make(map[string]interface{})
	err = s.Decode("auth", cookie.Value, &value)
	if err != nil {
		log.Println(err)
		removeCookies(w, r)
		return false
	}

	if value["login"] == config.AdminLogin {
		return true
	}
	removeCookies(w, r)
	return false
}
