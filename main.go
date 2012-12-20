package main

import (
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	//"time"
)

type Article struct {
	Title string
	Body  template.HTML
}

func loadPage(title string) (*Article, error) {
	filename := title + ".txt"
	fmt.Printf("loading article %s\n", filename)
	body, err := ioutil.ReadFile("articles/" + filename)
	if err != nil {
		return nil, err
	}
	return &Article{Title: title, Body: template.HTML(string(body))}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/article/"):]
	p, _ := loadPage(title)
	t, _ := template.ParseFiles("templates/article.html")
	t.Execute(w, p)
}

var host = flag.String("host", ":8080", "host to bind to")
var root = flag.String("root", "wwwroot", "webserver document root folder")

func main() {

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/article/", viewHandler)

	fmt.Printf("listening on %s // root=%s\r\n", *host, *root)

	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Printf("error : %v", err)
	}
}
