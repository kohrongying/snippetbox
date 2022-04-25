package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// create custom loggers
	// local date and local time joined using bitwise OR |
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)


	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}
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

	// Initialise http.Server struct
	srv := &http.Server {
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}