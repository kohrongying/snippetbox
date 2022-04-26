package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"github.com/kohrongying/snippetbox/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check if current request path exactly matches "/"
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	snippets, err := app.snippets.Latest()
	if err != nil { 
		app.serverError(w, err)
		return
	}

	app.render(w, http.StatusOK, "home.tmpl", &templateData{
		Snippets: snippets,
	})
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord)  {
			app.notFound(w)
			return
		} else {
			app.serverError(w, err)
			return
		}
	}

	app.render(w, http.StatusOK, "view.tmpl", &templateData{
		Snippet: snippet,
	})
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