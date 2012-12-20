package main

import (
	"flag"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Title string
	Body  template.HTML
}

func load(title string) (*Article, error) {
	filename := title + ".text"

	log.Printf("read %s\n", filename)

	body, err := ioutil.ReadFile("articles/" + filename)
	if err != nil {
		return nil, err
	}

	output := blackfriday.MarkdownCommon(body)

	return &Article{Title: title, Body: template.HTML(output)}, nil
}

func tocHandler(w http.ResponseWriter, r *http.Request) {
	//title := "table of contents"

	names, err := ioutil.ReadDir("articles")
	if err != nil {
		log.Printf("tocHandler: %v", err)
		return
	}

	t, err := template.ParseFiles("templates/common.html", "templates/toc.html")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("toc-t : %v", err)
		return
	}

	for _, element := range names {
		if !element.IsDir() {
			log.Printf("toc: %s", element.Name())
			t.ExecuteTemplate(w, "toc", element)
		}
	}
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/article/"):]

	p, err := load(title)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("error : %v\n", err)
		return
	}

	t, err := template.ParseFiles("templates/common.html", "templates/article.html")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("error : %v\n", err)
		return
	}

	t.ExecuteTemplate(w, "head", p)
	t.ExecuteTemplate(w, "bar", p)
	t.ExecuteTemplate(w, "article", p)
	t.ExecuteTemplate(w, "foot", p)
}

var host = flag.String("host", ":8080", "host to bind to")
var root = flag.String("root", "wwwroot", "webserver document root folder")

func main() {

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/toc/", tocHandler)

	fmt.Printf("listening on %s // root=%s\r\n", *host, *root)

	if err := http.ListenAndServe(*host, nil); err != nil {
		log.Printf("error : %v", err)
	}
}
