/*
	Copyright (c) 2012 Jacob Dearing
*/
package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vmihailenco/redis"
	"html/template"
	"log"
	"net/http"
)

// Represents and blog post.
// In the future timestamps, author information etc will be implemented.
type Article struct {
	Title string        // just the title
	Body  template.HTML // we consider the storage to be safe enough to generate HTML from (after markdown processing)
}

var conf = flag.String("conf", "blog.conf", "JSON configuration")
var generate = flag.Bool("generate", false, "generate a new config as conf is set")
var config Config
var client *redis.Client

//  MAIN
func main() {

	flag.Parse()

	if *generate {
		config.GenerateConfig(*conf)
		return
	}

	config.LoadConfig(*conf)

	if config.Verbose {
		log.Println("configuration loaded from " + *conf)
	}

	// Initialize contact with the server using our arguments or defaults.
	client = redis.NewTCPClient(config.RedisHost, config.RedisPass, config.RedisDB)
	defer client.Close()

	// If we can ping wihtout an error then we can move on.
	if ping := client.Ping(); ping.Err() != nil {
		log.Panicf("%v", ping.Err())
	}

	// Watch our content for changes
	go watch()

	// Push our working content to our redis db
	log.Println("pushing everything...")
	pushall(config.ContentFolder)

	//	Setup our handlers and get cracking...
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/toc/", tocHandler)
	r.HandleFunc("/p/{title}", contentHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.WWWRoot)))
	http.Handle("/", r)

	if err := http.ListenAndServe(config.WWWHost, nil); err != nil {
		log.Fatalf("%v", err)
	}

	fmt.Printf("listening on %s // root=%s\n", config.WWWHost, config.WWWRoot)
}
