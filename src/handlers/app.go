package handlers

import (
	"awesome_web_app/models"
	"awesome_web_app/settings"
	"database/sql"
	"fmt"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Models struct {
	UserModel *models.UserModel
}

type App struct {
	settings          *settings.Settings
	models            *Models
	router            *mux.Router
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
		router: mux.NewRouter().StrictSlash(true),
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

func (a *App) AddHandler(url string, handler http.HandlerFunc, name string) {
	a.router.HandleFunc(url, handler).Name(name)
}

func (a *App) Serve() {
	log.Printf("Starting server on :%v\n", a.settings.Port)
	csrfMiddleware := csrf.Protect([]byte(a.settings.CSRFSecret), csrf.Secure(false)) // TODO: Set Secure to true
	err := http.ListenAndServe(
		fmt.Sprintf(":%v", a.settings.Port),
		csrfMiddleware(a.ParseFormMiddleware(a.router)),
	)
	if err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func (a *App) URLGenerator() func(string, ...string) string {
	return func(name string, pairs ...string) string {
		route := a.router.Get(name)
		if route == nil {
			log.Printf("Route not found: %s", name)
			return ""
		}

		url, err := route.URL(pairs...)
		if err != nil {
			log.Printf("Error generating URL for route %s: %v", name, err)
			return ""
		}

		return url.String()
	}
}

func (a *App) renderTemplate(w http.ResponseWriter, tmpl string, data map[string]any) {
	// TODO: Cache templates
	templates, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		log.Printf("Error loading templates: %v", err)
	}
	templates = append(templates, "templates/layout.html")
	templates = append(templates, "templates/"+tmpl+".html")

	data["URL"] = a.URLGenerator()

	t, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
	}
	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

func (a *App) ParseFormMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			r.PostForm.Del("gorilla.csrf.Token") // TODO: get name from settings
		}
		next.ServeHTTP(w, r)
	})
}
