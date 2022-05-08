package main

import (
	"crypto/tls"
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
	sessionManager.Cookie.Secure = true // Set to true to ensure cookie sent over HTTPS

	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
		templateCache: templateCache,
		formDecoder: formDecoder,
		sessionManager: sessionManager,
	}

	// Initialize tls.Config to hold settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion: tls.VersionTLS12,
    	MaxVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	// Initialise http.Server struct
	srv := &http.Server {
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		TLSConfig: tlsConfig,
		IdleTimeout: time.Minute, // keep alive connections to close after 1 min
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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