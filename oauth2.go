/*	
	-------------------------------------------------------------
		OAUTH2 Bullshit
		Experimenting with github OAUTH
	-------------------------------------------------------------
*/

package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"code.google.com/p/goauth2/oauth"
	//store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/securecookie"
)

var oauth_config oauth.Config
var NewState = make(chan string)

func initOauth2() {
	oauth_config = oauth.Config{
		ClientId:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Scope:        "",
		AuthURL:      "https://github.com/login/oauth/authorize",
		TokenURL:     "https://github.com/login/oauth/access_token",
		RedirectURL:  config.RedirectURL,
	}

	go func() {
		h := sha1.New()
		for {
			h.Write([]byte(time.Now().String()))
			NewState <- fmt.Sprintf("%X", h.Sum(nil))
		}
	}()

}

var hashKey = []byte(securecookie.GenerateRandomKey(32))
var blockKey = []byte(securecookie.GenerateRandomKey(32))
var s = securecookie.New(hashKey, blockKey)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if !validateCookie(w, r) {
		state := <-NewState
		setStateCookie(w, r, state)
		url := oauth_config.AuthCodeURL(state)

		log.Println(state)
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		http.Redirect(w, r, "/", http.StatusAccepted)
	}
}

// This callback contains a temp code and our passed secret from Github
func callbackHandler(w http.ResponseWriter, r *http.Request) {

	code := r.FormValue("code")   // temp code for authentication
	state := r.FormValue("state") // secret between us and authenticator

	log.Println(code, state)

	if state != getState(w, r) {
		log.Println("state did not match")
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}

	// We attempt to exchange the temp code we got for a real access token...
	t := &oauth.Transport{Config: &oauth_config}
	_, e := t.Exchange(code)
	if e != nil {
		log.Println(e)
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
	c := t.Client()

	// Now we have a client with a valid access token.
	// We use this access token to request the current authenticated user information.
	// This will be the same an public read to the Github API v3 however, the /user URL
	// alone only returns the authenticated user according to the token.  So this should
	// suffice as proof that the user IS a valid Github User.

	resp, _ := c.Get("https://api.github.com/user")
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	// We have a valid Github User but is it our admin?
	var info map[string]interface{}
	if err := json.Unmarshal(contents, &info); err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusInternalServerError)
	}

	// Now here we look at the login string in the JSON the server sent us for our ADMIN check.
	login := info["login"]
	if login != config.AdminLogin {

		// USER ain't our admin, move along.
		log.Println("Github user did not match admin.")
		removeCookies(w, r)
		http.Redirect(w, r, "/", http.StatusUnauthorized)

	} else {

		// We need drop a cookie now to retain login status
		log.Println("Admin logged in.")
		removeCookies(w, r)
		setAuthCookie(w, r, login)
		http.Redirect(w, r, "/", http.StatusAccepted)
	}

}

// Delete login cookie
func logoutHander(w http.ResponseWriter, r *http.Request) {
	removeCookies(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Testing secret page
func secretPageHandler(w http.ResponseWriter, r *http.Request) {
	if validateCookie(w, r) {
		http.Redirect(w, r, "/", http.StatusAccepted)
	} else {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
	}
}

// 5 minutes to keep state cookie for cross-site forgery attempts
// I could be storing this internally but this *is* a secure cookie right?
func setStateCookie(w http.ResponseWriter, r *http.Request, state string) {
	if encoded, err := s.Encode("state", state); err == nil {
		cookie := &http.Cookie{
			Name:    "state",
			Value:   encoded,
			Path:    "/",
			Expires: time.Now().Add(5 * time.Minute),
		}
		http.SetCookie(w, cookie)
	}
}

// Set a secure cookie
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
			Expires: time.Now().Add(time.Hour),
		}
		cookie.Expires = time.Now().Add(time.Hour)
		http.SetCookie(w, cookie)
	}
}

// Remove a secure cookie
func removeCookies(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookies() {

		// MaxAge < 0 means delete cookie now
		// http://golang.org/pkg/net/http/#Cookie
		v.MaxAge = -1
		http.SetCookie(w, v)
	}
}

// Returns true if the cookie the client provided checks out.
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

// Returns true if the cookie the client provided checks out.
func validateCookie(w http.ResponseWriter, r *http.Request) bool {
	if cookie, err := r.Cookie("auth"); err == nil {
		value := make(map[string]interface{})
		if err = s.Decode("auth", cookie.Value, &value); err == nil {
			if value["login"] == config.AdminLogin {
				return true
			}
		} else {
			log.Println(err)
		}
	}
	return false
}
