package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	log.Println("started.")
	pool = newPool("virtual-arch:6379", "")

	for i := 0; i < 3; i++ {
		p := create()
		p.Title = "Test"
		p.Content = "test content :: " + p.UUID
		p.Author = "jacob.dearing@gmail.com"
		p.save()
		p.load()
		//p.delete()
	}

	r := mux.NewRouter()
	//r.HandleFunc("/", indexHandler) // index
	//r.HandleFunc("/toc", tocHandler)            // table of contents
	r.HandleFunc("/{uuid}", contentHandler) // display a post with title
	//r.HandleFunc("/e/{id}", editContentHandler) // edit a post
	//r.HandleFunc("/s/{id}", saveContentHandler) // save a post
	//r.HandleFunc("/login", loginHandler)        // fire up Outh2
	//r.HandleFunc("/logout", logoutHander)       // ''
	//r.HandleFunc("/callback", callbackHandler)  // Outh2 callback addy
	//r.HandleFunc("/secret", secretPageHandler)  // simple login testing handler

	http.Handle("/", r)
	if err := http.ListenAndServe("localhost:9000", nil); err != nil {
		log.Panicln("%v\n", err)
	}

}
