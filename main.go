package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// check if current request path exactly matches "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("hello from snippetbox"))
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a specific snippet"))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// check request method
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		http.Error(w, "Method not supported" , http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("create a new snippet"))
}

func main() {
	mux := http.NewServeMux()
	
	// handlers 
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)
	// slash is a catch all. eg. /foo, /bash --> home
	mux.HandleFunc("/", home) 

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}