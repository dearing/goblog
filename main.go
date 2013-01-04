/*
	Blogging with go, markdown and redis.
	Copyright (c) 2012 Jacob Dearing
*/
package main

import (
	"flag"
	"fmt"
	store "github.com/dearing/blog/storage/redis"
	"github.com/gorilla/mux"
	//"html/template"
	"log"
	"net/http"
)

var conf = flag.String("conf", "blog.conf", "JSON configuration")
var generate = flag.Bool("generate", false, "generate a new config as conf is set")
var config Config

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
	store.Connect(config.RedisHost, config.RedisPass, config.RedisDB)
	store.LoadDirectory(config.ContentFolder)

	//	Setup our handlers and get cracking...
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/toc/", tocHandler)
	r.HandleFunc("/p/{id}", contentHandler)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.WWWRoot)))
	http.Handle("/", r)

	if err := http.ListenAndServe(config.WWWHost, nil); err != nil {
		log.Printf("%v\n", err)
	}

	fmt.Printf("listening on %s // root=%s\n", config.WWWHost, config.WWWRoot)
}
