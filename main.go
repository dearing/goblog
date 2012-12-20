package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/russross/blackfriday"
	"github.com/vmihailenco/redis"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Article struct {
	Title string
	Body  template.HTML
}

func push(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	client.Set(filename, string(blackfriday.MarkdownCommon(body)))
	log.Printf("pushed %s", filename)
	return err
}

func pushall(folder string) error {

	names, _ := ioutil.ReadDir("articles")

	log.Printf("%v", names)

	for _, element := range names {
		body, _ := ioutil.ReadFile("articles/" + element.Name())
		client.Set("articles/"+element.Name(), string(blackfriday.MarkdownCommon(body)))
		log.Printf("pushed %s", "articles/"+element.Name())
	}

	return nil
}

func load(title string) (*Article, error) {

	key := fmt.Sprintf("articles/" + title + ".text")
	log.Printf("%v", key)
	if !client.Exists(key).Val() {
		log.Printf("no found: %v", key)
	}

	output := client.Get(key).Val()
	log.Printf("%v", output)
	return &Article{Title: title, Body: template.HTML(output)}, nil
}

func tocHandler(w http.ResponseWriter, r *http.Request) {
	title := "table of contents"

	keys := client.Keys("*")
	log.Printf("%v", keys.Val())
	for _, element := range keys.Val() {
		log.Printf("%v", element)
	}

	names, err := ioutil.ReadDir("articles")
	if err != nil {
		log.Printf("tocHandler: %v", err)
		return
	}

	//t, err := template.ParseFiles("templates/common.html", "templates/toc.html")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("toc-t : %v", err)
		return
	}

	templates.ExecuteTemplate(w, "head", title)
	templates.ExecuteTemplate(w, "bar", nil)
	templates.ExecuteTemplate(w, "toc-head", nil)
	for _, element := range names {
		if !element.IsDir() {
			url := strings.Replace(element.Name(), ".text", "", 1)
			templates.ExecuteTemplate(w, "toc-item", url)
		}
	}

	templates.ExecuteTemplate(w, "toc-foot", nil)
}

func articleHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/article/"):]

	p, err := load(title)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
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

var host = flag.String("host", ":8080", "host to bind to")
var root = flag.String("root", "wwwroot", "webserver document root folder")
var templates = template.Must(template.ParseGlob("templates/*.html"))
var client = redis.NewTCPClient("192.168.1.150:6379", "", -1)

func main() {

	defer client.Close()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				log.Println("event:", ev)
				push(ev.Name)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch("articles")
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	pushall("articles")

	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/toc/", tocHandler)

	fmt.Printf("listening on %s // root=%s\r\n", *host, *root)

	if err = http.ListenAndServe(*host, nil); err != nil {
		log.Printf("error : %v", err)
	}
}
