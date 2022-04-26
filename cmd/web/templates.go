package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/kohrongying/snippetbox/internal/models"
)

// Define this as a holding structure for dynamic data
type templateData struct {
	Snippet 	*models.Snippet
	Snippets	[]*models.Snippet
	CurrentYear int
}


func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// custom template functions can take multiple parameters but only return one value (+err)
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error)  {
	// initialise map to act as cache
	cache := map[string]*template.Template{}
	
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err	
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// create empty template set, use Funcs method to register the FuncMap then parse base.tmpl
		emptyTs := template.New(name)
		ts, err := emptyTs.Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		// parse partials on template set
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to map
		cache[name] = ts
	}
	return cache, nil
}