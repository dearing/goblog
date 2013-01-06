package main

import (
	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strings"
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

	keys := store.Keys("post:*")

	// for each key we add a list element
	for _, element := range keys.Val() {

		key := strings.TrimLeft(element, "post:")

		p, err := store.Get(key)
		if err != nil {
			log.Println(err)
		}

		t.ExecuteTemplate(w, "toc-item", p)

	}

	t.ExecuteTemplate(w, "toc-foot", nil)
}

func contentHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := store.Get(id)
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

	p, err := store.Get("index.md")
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
