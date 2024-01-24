package handlers

import (
	"awesome_web_app/models"
	"net/http"
)

func (a *App) IndexPage(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)
	a.renderTemplate(w, r, "index", "Home", map[string]any{
		"CurrentUser": currentUser,
	})
}
