package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	// check if current request path exactly matches "/"
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// parse template and catch error
	files := []string{
		"./ui/html/base.tmpl", // base must be first
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...) //pass as variadic parameter
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil) // use content of base template to response which invokes title and main templates
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with id %d", id)
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