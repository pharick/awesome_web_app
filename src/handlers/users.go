package handlers

import (
	"awesome_web_app/models"
	"github.com/gorilla/mux"
	"net/http"
)

func (a *App) UserList(w http.ResponseWriter, r *http.Request) {
	users, err := a.models.UserModel.GetListWithUsername()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser, _ := r.Context().Value("user").(*models.User)
	a.renderTemplate(w, r, "userList", "User List", map[string]any{
		"CurrentUser": currentUser,
		"Users":       users,
	})
}

func (a *App) UserPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	user, err := a.models.UserModel.GetByUsername(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil || !user.Username.Valid {
		http.NotFound(w, r)
		return
	}

	currentUser, _ := r.Context().Value("user").(*models.User)
	a.renderTemplate(w, r, "userPage", user.Username.String, map[string]any{
		"CurrentUser": currentUser,
		"User":        user,
	})
}
