package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello from snippetbox"))
}

func main() {
	mux := http.NewServeMux()
	
	// slash is a catch all. eg. /foo, /bash --> home
	mux.HandleFunc("/", home) 

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}