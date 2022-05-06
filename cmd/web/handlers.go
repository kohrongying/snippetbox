package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	
	"github.com/kohrongying/snippetbox/internal/models"
	"github.com/kohrongying/snippetbox/internal/validator"

	"github.com/julienschmidt/httprouter"
)

type snippetCreateForm struct {
	Title 		string `form:"title"` //note cannot have space
	Content		string `form:"content"`
	Expires 	int	   `form:"expires"`
	validator.Validator `form:"-"` //this struct inherits all fields/methods of Validator
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// EDIT: httprouter handles this for us
	// check if current request path exactly matches "/"
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	snippets, err := app.snippets.Latest()
	if err != nil { 
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets
	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// EDIT: handled by httprouter
	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
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

	data := app.newTemplateData(r)
	data.Snippet = snippet
	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request)  {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{
		Expires: 365,
	}
	app.render(w, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// EDIT: done by httprouter 
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", "POST")
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	var form snippetCreateForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MaxCharacters(form.Title, 100), "title", "This field cannot be more than 100 characters long")
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedInt(form.Expires, 1, 7, 365), "expires", "This field must be equal to 1, 7 or 364")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}
	
	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.infoLog.Printf("id is %d", id)
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}