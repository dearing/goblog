package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	//"time"
)

func contentHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseGlob("template/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}

	vars := mux.Vars(r)
	uuid := vars["uuid"]

	if !exists(uuid) {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	p := &Page{}
	p.UUID = uuid
	p.load()

	t.ExecuteTemplate(w, "page.html", p)
}
