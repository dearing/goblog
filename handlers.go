package main

import (
	"html/template"
	"log"
	"net/http"

	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
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
		Admin: false,
		Post:  p,
	}

	t.ExecuteTemplate(w, "post.html", page)
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
		Admin: false,
		Post:  p,
	}

	t.ExecuteTemplate(w, "post.html", page)
}
