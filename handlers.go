package main

import (
	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Page struct {
	Admin bool
	Post  store.Post
}

func tocHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	keys := store.GetPosts()
	posts := make(map[string]store.Post)
	for _, key := range keys.Val() {

		p, err := store.Get(key, false)
		if err != nil {
			log.Println(err)
			continue
		}

		posts[key] = p

	}

	t.ExecuteTemplate(w, "toc.html", posts)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := store.Get(id, true)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	page := &Page{
		Admin: validateCookie(w, r),
		Post:  p,
	}

	t.ExecuteTemplate(w, "post.html", page)
}

func editContentHandler(w http.ResponseWriter, r *http.Request) {

	if !validateCookie(w, r) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := store.GetRaw(id, false)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	page := &Page{
		Admin: validateCookie(w, r),
		Post:  p,
	}

	t.ExecuteTemplate(w, "edit.html", page)
}

func saveContentHandler(w http.ResponseWriter, r *http.Request) {

	if !validateCookie(w, r) {
		http.Redirect(w, r, "/", http.StatusUnauthorized)
		return
	}

	p := store.Post{
		ID:       r.PostFormValue("id"),
		Title:    r.PostFormValue("title"),
		Content:  template.HTML(r.PostFormValue("content")), //posible bug, dunno yet
		Modified: time.Now(),
	}

	log.Println(p)

	store.New(p)

	http.Redirect(w, r, "/p/"+p.ID, http.StatusFound)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	p, err := store.GetLatest()
	if err != nil {
		log.Println(err)
		return
	}

	page := &Page{
		Admin: validateCookie(w, r),
		Post:  p,
	}

	log.Println(page.Admin)
	t.ExecuteTemplate(w, "post.html", page)
}
