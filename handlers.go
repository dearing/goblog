/*
	Copyright (c) 2012 Jacob Dearing
*/
package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Generate a simple list of article titles and links from redis.
// TODO: better naming scheme as at PUSHALL
func tocHandler(w http.ResponseWriter, r *http.Request) {
	title := "table of contents"

	keys := client.Keys(*content + "/*")

	t, err := template.ParseGlob(*templates + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", title)
	t.ExecuteTemplate(w, "bar", nil)
	t.ExecuteTemplate(w, "toc-head", nil)

	// for each key we add a list element
	for _, element := range keys.Val() {
		if element != *content+"/index.md" {

			url := strings.Replace(element, ".md", "", 1)
			url = strings.Replace(url, *content+"/", "", 1)
			t.ExecuteTemplate(w, "toc-item", url)
		}
	}

	t.ExecuteTemplate(w, "toc-foot", nil)
}

// Load and display an article from our redis db.
func contentHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	title := vars["title"]

	key := fmt.Sprintf("%s/%s.md",*content,title)

	p, err := pull(key)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t, err := template.ParseGlob(*templates + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", p)
	t.ExecuteTemplate(w, "foot", p)
}

// Load and display an article from our redis db.
func indexHandler(w http.ResponseWriter, r *http.Request) {

	src := fmt.Sprintf("%s/index.md", *content)

	p, err := pull(src)
	if err != nil {
		log.Printf("error : %v\n", err)
		return
	}

	//t, err := template.ParseFiles("templates/common.html", "templates/article.html")
	t, err := template.ParseGlob(*templates + "/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", p)
	t.ExecuteTemplate(w, "foot", p)
}
