package main

import (
	"log"
	"net/http"
	"strings"
)

// Generate a simple list of article titles and links from redis.
// TODO: better naming scheme as at PUSHALL
func tocHandler(w http.ResponseWriter, r *http.Request) {
	title := "table of contents"

	keys := client.Keys("articles/*")

	templates.ExecuteTemplate(w, "head", title)
	templates.ExecuteTemplate(w, "bar", nil)
	templates.ExecuteTemplate(w, "toc-head", nil)

	// for each key we add a list element
	for _, element := range keys.Val() {
		url := strings.Replace(element, ".text", "", 1)
		url = strings.Replace(url, "articles/", "", 1)
		templates.ExecuteTemplate(w, "toc-item", url)
	}

	templates.ExecuteTemplate(w, "toc-foot", nil)
}

// Load and display an article from our redis db.
func articleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract a meaningful title from the path.
	title := r.URL.Path[len("/article/"):]

	p, err := pull("articles/" + title + ".text")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	//t, err := template.ParseFiles("templates/common.html", "templates/article.html")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("error : %v\n", err)
		return
	}

	templates.ExecuteTemplate(w, "head", p)
	templates.ExecuteTemplate(w, "bar", p)
	templates.ExecuteTemplate(w, "article", p)
	templates.ExecuteTemplate(w, "foot", p)
}
