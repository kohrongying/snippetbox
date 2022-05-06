package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kohrongying/snippetbox/internal/models"

	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/lib/pq"
)

type application struct {
	errorLog 		*log.Logger
	infoLog  		*log.Logger
	snippets 		*models.SnippetModel
	templateCache 	map[string]*template.Template
	formDecoder		*form.Decoder
	sessionManager	*scs.SessionManager
}

func main() {
	// command line flag
	addr := flag.String("addr", ":4000", "HTTP network address")
	connStr := flag.String("connStr", "user=snippetboxweb dbname=snippetbox sslmode=disable", "PSQL data source name")
	flag.Parse()

	// create custom loggers
	// local date and local time joined using bitwise OR |
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// connect to postgresDB
	db, err := openDB(*connStr)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// initialise new template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)	
	}

	// initialize form decoder instance
	formDecoder := form.NewDecoder()

	// initialize session manager
	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// Initialise http.Server struct
	srv := &http.Server {
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}