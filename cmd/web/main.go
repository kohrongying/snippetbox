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

	// Initialise http.Server struct
	srv := &http.Server {
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}