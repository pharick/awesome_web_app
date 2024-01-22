package handlers

import (
	"awesome_web_app/models"
	"github.com/gorilla/mux"
	"net/http"
)

type UserListData struct {
	Title       string
	CurrentUser *models.User
	Users       []models.User
}

func (a *App) UserList(w http.ResponseWriter, r *http.Request) {
	users, err := a.models.UserModel.GetListWithUsername()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	currentUser, _ := r.Context().Value("user").(*models.User)
	_ = renderTemplate(w, "userList", UserListData{Title: "User List", CurrentUser: currentUser, Users: users})
}

type UserPageData struct {
	Title       string
	CurrentUser *models.User
	User        *models.User
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
	_ = renderTemplate(w, "userPage", UserPageData{Title: user.Username.String, CurrentUser: currentUser, User: user})
}
