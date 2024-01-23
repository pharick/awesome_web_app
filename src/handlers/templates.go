package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	// TODO: Cache templates
	templates, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		return err
	}

	templates = append(templates, "templates/layout.html")
	templates = append(templates, "templates/"+tmpl+".html")

	t, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "layout", data)
	return err
}
