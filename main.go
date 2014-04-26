package main

import (
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var build string
var conf = flag.String("conf", "blog.conf", "JSON configuration")
var gen = flag.Bool("gen", false, "generate a new config as conf is set")
var config Config

func main() {
	flag.Parse()

	if *gen {
		config.GenerateConfig(*conf)
		log.Println("generated new config at", *conf)
		return
	}

	config.LoadConfig(*conf)

	log.Printf("version %s", build)
	pool = newPool(config.RedisHost, config.RedisPass)

	for i := 0; i < 1; i++ {
		p := create()
		p.Title = "Test"
		p.Content = "test content :: " + p.UUID
		p.Author = "somebody"
		p.save()
		p.load()
		//p.delete()
	}

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler)            // index
	r.HandleFunc("/p/{uuid}", contentHandler)  // display a post with title
	r.HandleFunc("/login", loginHandler)       // fire up Outh2
	r.HandleFunc("/logout", logoutHander)      // ''
	r.HandleFunc("/callback", callbackHandler) // Outh2 callback addy
	r.HandleFunc("/secret", secretPageHandler) // simple login testing handler

	http.Handle("/", r)
	if err := http.ListenAndServe(config.WWWHost, nil); err != nil {
		log.Panicf("%v\n", err)
	}

}
