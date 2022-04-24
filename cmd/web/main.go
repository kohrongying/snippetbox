package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	
	// handlers 
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// slash is a catch all. eg. /foo, /bash --> home
	mux.HandleFunc("/", home) 
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}