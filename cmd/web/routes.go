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

	// middleware chain for dynamic routes with session manager
	// add noSurf middleware on all dynamic routes to protect from CSRF
	// add authenticate middleware
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// same as router.Handle("/", &home{}) or
	// same as router.Handle("/", http.HandlerFunc(home))
	// home -> func (h *home) ServeHTTP(w, r)
	// servemux doesnt support clean url (view/:id), using httprouter
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	// slash is a catch all. eg. /foo, /bash --> home
	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))

	// protected (authenticated-only) application routes using new requireAuthentication middleware
	protected := dynamic.Append(app.requireAuthentication)
	router.Handler(http.MethodGet, "/snippet/view/:id", protected.ThenFunc(app.snippetView)) 
	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreate))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePost))

	// create middleware chain
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}