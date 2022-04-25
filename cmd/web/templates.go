package main

import (
	"github.com/kohrongying/snippetbox/internal/models"
)

// Define this as a holding structure for dynamic data
type templateData struct {
	Snippet 	*models.Snippet
	Snippets	[]*models.Snippet
}