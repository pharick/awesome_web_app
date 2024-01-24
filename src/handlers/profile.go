package handlers

import (
	"awesome_web_app/models"
	"github.com/gorilla/csrf"
	"net/http"
)

type ProfileForm struct {
	Username string
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := r.Context().Value("user").(*models.User)

	if r.Method == http.MethodPost {
		var form ProfileForm
		err := a.formDecoder.Decode(&form, r.PostForm)
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

	a.renderTemplate(w, "profile", map[string]any{
		"Title":          "Your Profile",
		"CurrentUser":    currentUser,
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}
