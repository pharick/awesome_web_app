package handlers

import (
	"awesome_web_app/models"
	"context"
	"net/http"

	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func (a *App) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := a.sessions.Get(r, "session")
		userId := session.Values["userId"]
		if userId == nil {
			next(w, r)
			return
		}

		user, err := a.models.UserModel.GetById(userId.(int))
		if err != nil {
			delete(session.Values, "userId")
			_ = session.Save(r, w)
			next(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next(w, r.WithContext(ctx))
	}
}

func (a *App) AuthRequiredMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return a.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		_, ok := r.Context().Value("user").(*models.User)
		if !ok {
			http.Redirect(w, r, a.generateUrl("login"), http.StatusFound)
			return
		}
		next(w, r)
	})
}

func (a *App) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := a.googleOAuthConfig.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func (a *App) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")

	token, err := a.googleOAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := a.googleOAuthConfig.Client(r.Context(), token)
	google, err := oauth2.NewService(r.Context(), option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userinfo, err := google.Userinfo.Get().Do()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := a.models.UserModel.GetByEmail(userinfo.Email)
	if user == nil {
		user, err = a.models.UserModel.Create(userinfo.Email)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, _ := a.sessions.Get(r, "session")
	session.Values["userId"] = user.Id
	_ = session.Save(r, w)

	http.Redirect(w, r, a.generateUrl("profile"), http.StatusFound)
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := a.sessions.Get(r, "session")
	session.Options.MaxAge = -1
	_ = session.Save(r, w)
	http.Redirect(w, r, a.generateUrl("index"), http.StatusSeeOther)
}
