package main

import (
	"html/template"
	"log"
	"net/http"

	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
)

func tocHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	t.ExecuteTemplate(w, "head", nil)
	t.ExecuteTemplate(w, "bar", nil)
	t.ExecuteTemplate(w, "toc-head", nil)

	keys := store.GetPosts()

	// for each key we add a list element
	for _, key := range keys.Val() {

		p, err := store.Get(key, false)
		if err != nil {
			log.Println(err)
			continue
		}

		if key != "index" {
			t.ExecuteTemplate(w, "toc-item", p)
		}

	}

	t.ExecuteTemplate(w, "toc-foot", nil)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := store.Get(id, true)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", p)
	t.ExecuteTemplate(w, "foot", p)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	//p, err := store.Get("index", true)
	p, err := store.GetLatest()
	if err != nil {
		log.Println(err)
		return
	}

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", p)
	t.ExecuteTemplate(w, "foot", p)
}
