package handlers

import (
	"awesome_web_app/models"
	"log"
	"net/http"
)

type IndexData struct {
	Title       string
	CurrentUser *models.User
}

func (a *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)
	err := renderTemplate(w, "index", IndexData{Title: "Home", CurrentUser: currentUser})
	if err != nil {
		log.Println(err)
	}
}
