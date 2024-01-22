package handlers

import (
	"awesome_web_app/models"
	"awesome_web_app/settings"
	"database/sql"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Models struct {
	UserModel *models.UserModel
}

type App struct {
	settings          *settings.Settings
	models            *Models
	googleOAuthConfig *oauth2.Config
	sessions          *sessions.CookieStore
	formDecoder       *schema.Decoder
}

func NewApp(settings *settings.Settings, db *sql.DB) *App {
	return &App{
		settings: settings,
		models: &Models{
			UserModel: &models.UserModel{DB: db},
		},
		googleOAuthConfig: &oauth2.Config{
			ClientID:     settings.Google.ClientID,
			ClientSecret: settings.Google.ClientSecret,
			RedirectURL:  settings.BaseUrl + "/auth/callback/",
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		sessions:    sessions.NewCookieStore([]byte(settings.SessionSecret)),
		formDecoder: schema.NewDecoder(),
	}
}
