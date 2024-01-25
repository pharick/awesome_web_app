package handlers

import (
	"awesome_web_app/models"
	"net/http"
)

func (a *App) ProfilePage(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)

	a.renderTemplate(w, r, "profile", "Your Profile", map[string]any{
		"CurrentUser": currentUser,
	})
}

type ProfileForm struct {
	Username string `validate:"required,min=3,max=40"`
}

func (a *App) HandleProfileForm(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)

	var form ProfileForm
	validationErrors, err := a.parseForm(r, &form)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := a.sessions.Get(r, "flash")
	if validationErrors != nil {
		for _, fieldError := range validationErrors {
			session.AddFlash(fieldError.Error()) // TODO: translate validation errors
		}
		_ = session.Save(r, w)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
	user, _ := a.models.UserModel.GetByUsername(form.Username)
	if user.Id == currentUser.Id {
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}
	if user != nil {
		session.AddFlash("Username already taken")
		_ = session.Save(r, w)
		http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
		return
	}

	currentUser, err = a.models.UserModel.Update(currentUser.Id, form.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.URL.Path, http.StatusSeeOther)
}
