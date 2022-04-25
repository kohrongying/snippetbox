package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	// same as mux.Handle("/", &home{}) or
	// same as mux.Handle("/", http.HandlerFunc(home))
	// home -> func (h *home) ServeHTTP(w, r)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// slash is a catch all. eg. /foo, /bash --> home
	mux.HandleFunc("/", app.home)

	fileServer := http.FileServer(http.Dir("./ui/static/")) // static file server
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	return mux
}