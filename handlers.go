package main

import (
	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
	"github.com/russross/blackfriday"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func tocHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", nil)
	t.ExecuteTemplate(w, "bar", nil)
	t.ExecuteTemplate(w, "toc-head", nil)

	keys := store.Keys("post:*")

	// for each key we add a list element
	for _, element := range keys.Val() {

		key := strings.Replace(element, "post:", "", 1)

		p, err := store.Get(key)
		if err != nil {
			log.Println(err)
		}

		log.Println(p)

		t.ExecuteTemplate(w, "toc-item", p)

	}

	t.ExecuteTemplate(w, "toc-foot", nil)
}

// display content from storage
func contentHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id := vars["id"]

	p, err := store.Get(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", getHTML(p.Content))
	t.ExecuteTemplate(w, "foot", p)
}

func getHTML(content string) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(content)))
}

// Load and display an article from our redis db.
func indexHandler(w http.ResponseWriter, r *http.Request) {

	p, err := store.Get("1")
	if err != nil {
		log.Printf("error : %v\n", err)
		return
	}

	//t, err := template.ParseFiles("templates/common.html", "templates/article.html")
	t, err := template.ParseGlob(config.TemplateFolder + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", getHTML(p.Content))
	t.ExecuteTemplate(w, "foot", p)
}
