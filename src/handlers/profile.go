package handlers

import (
	"awesome_web_app/models"
	"net/http"
)

type ProfileData struct {
	Title       string
	CurrentUser *models.User
}

type ProfileForm struct {
	Username string
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var form ProfileForm
		err = a.formDecoder.Decode(&form, r.PostForm)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		currentUser, err = a.models.UserModel.Update(currentUser.Id, form.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	_ = renderTemplate(w, "profile", ProfileData{Title: "Your Profile", CurrentUser: currentUser})
}
