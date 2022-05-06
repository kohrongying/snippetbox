package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"
)

// for middleware, change return type from servemux to http.handler
func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})


	fileServer := http.FileServer(http.Dir("./ui/static/")) // static file server
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	// middleware chain for dynamic routes
	dynamic := alice.New(app.sessionManager.LoadAndSave)

	// same as router.Handle("/", &home{}) or
	// same as router.Handle("/", http.HandlerFunc(home))
	// home -> func (h *home) ServeHTTP(w, r)
	// servemux doesnt support clean url (view/:id), using httprouter
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView)) 
	router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// slash is a catch all. eg. /foo, /bash --> home
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	
	// create middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}