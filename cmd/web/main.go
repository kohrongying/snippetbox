package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	
	// handlers 
	// same as mux.Handle("/", &home{}) or
	// same as mux.Handle("/", http.HandlerFunc(home))
	// home -> func (h *home) ServeHTTP(w, r)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// slash is a catch all. eg. /foo, /bash --> home
	mux.HandleFunc("/", home) 
	
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Printf("Starting server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}