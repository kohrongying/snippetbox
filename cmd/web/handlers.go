package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check if current request path exactly matches "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// parse template and catch error
	files := []string{
		"./ui/html/base.tmpl", // base must be first
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}
	ts, err := template.ParseFiles(files...) //pass as variadic parameter
	if err != nil {
		app.serverError(w, err) // home handler is method against application, can access its fields
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil) // use content of base template to response which invokes title and main templates
	if err != nil {
		app.serverError(w,err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with id %d", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// check request method
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 turt"
	content := "Turt is climbing mount fuji.\nSlowly and slowly!"
	expires := 2
	
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Printf("id is %d", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}