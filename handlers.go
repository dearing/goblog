/*
	Jacob Dearing
*/
package main

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

// Generate a simple list of article titles and links from redis.
// TODO: better naming scheme as at PUSHALL
func tocHandler(w http.ResponseWriter, r *http.Request) {
	title := "table of contents"

	keys := client.Keys("articles/*")

	t, err := template.ParseFiles("templates/common.html", "templates/toc.html")
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
		url := strings.Replace(element, ".md", "", 1)
		url = strings.Replace(url, "articles/", "", 1)
		t.ExecuteTemplate(w, "toc-item", url)
	}

	t.ExecuteTemplate(w, "toc-foot", nil)
}

// Load and display an article from our redis db.
func articleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract a meaningful title from the path.
	title := r.URL.Path[len("/blog/"):]

	p, err := pull("articles/" + title + ".md")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t, err := template.ParseFiles("templates/common.html", "templates/article.html")
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

	p, err := pull("articles/index.md")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	t, err := template.ParseFiles("templates/common.html", "templates/article.html")
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
