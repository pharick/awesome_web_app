package handlers

import (
	"html/template"
	"net/http"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) error {
	t, err := template.ParseFiles("templates/layout.html", "templates/"+tmpl+".html")
	if err != nil {
		return err
	}
	err = t.ExecuteTemplate(w, "layout", data)
	return err
}
