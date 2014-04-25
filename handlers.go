package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type ErrorPage struct {
	Code int
	Text string
}

func generateError(code int) *ErrorPage {
	return &ErrorPage{
		Code: code,
		Text: http.StatusText(code),
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseGlob("template/*")
	if err != nil {
		errorHandler(w, r, 404)
		log.Println(err)
		return
	}
	t.ExecuteTemplate(w, "index.html", nil)
}

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
		errorHandler(w, r, 404)
		return
	}

	p := &Page{}
	p.UUID = uuid
	p.load()

	t.ExecuteTemplate(w, "page.html", p)
}

func errorHandler(w http.ResponseWriter, r *http.Request, code int) {

	e := generateError(code)

	t, err := template.ParseGlob("template/*")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Println(err)
		return
	}
	log.Println(t)

	//http.Error(w, e.Text, e.Code)
	t.ExecuteTemplate(w, "error.html", e)
}
