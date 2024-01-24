package handlers

import (
	"awesome_web_app/models"
	"net/http"
)

type ProfileForm struct {
	Username string `validate:"required,min=3,max=40"`
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)

	if r.Method == http.MethodPost {
		var form ProfileForm

		err := a.ParseForm(r, &form)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = a.ValidateForm(&form, w, r)
		if err != nil {
			http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
			return
		}

		currentUser, err = a.models.UserModel.Update(currentUser.Id, form.Username)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	a.renderTemplate(w, r, "profile", "Your Profile", map[string]any{
		"CurrentUser": currentUser,
	})
}
