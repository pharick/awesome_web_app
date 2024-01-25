package handlers

import (
	"awesome_web_app/models"
	"awesome_web_app/settings"
	"database/sql"
	"errors"
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/csrf"
	gorillaHandlers "github.com/gorilla/handlers"
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
	settings                   *settings.Settings
	models                     *Models
	router                     *mux.Router
	googleOAuthConfig          *oauth2.Config
	sessions                   *sessions.CookieStore
	formDecoder                *schema.Decoder
	validator                  *validator.Validate
	validationErrorsTranslator *ut.Translator
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
		validator:   validator.New(),
	}
}

func (a *App) AddHandler(url string, name string, handlers map[string]http.HandlerFunc) {
	convertedHandlers := make(map[string]http.Handler)
	for method, handlerFunc := range handlers {
		convertedHandlers[method] = handlerFunc
	}
	a.router.Handle(url, gorillaHandlers.MethodHandler(convertedHandlers)).Name(name)
}

func (a *App) Serve() {
	log.Printf("Starting server on :%v\n", a.settings.Port)
	csrfMiddleware := csrf.Protect([]byte(a.settings.CSRFSecret), csrf.Secure(false), csrf.Path("/")) // TODO: Set Secure to true
	err := http.ListenAndServe(
		fmt.Sprintf(":%v", a.settings.Port),
		csrfMiddleware(a.router),
	)
	if err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}

func (a *App) generateUrl(name string, pairs ...string) string {
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

func (a *App) renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	tmpl string,
	title string,
	data map[string]any,
) {
	// TODO: Cache templates
	templates, err := filepath.Glob("templates/partials/*.html")
	if err != nil {
		log.Printf("Error loading templates: %v", err)
	}
	templates = append(templates, "templates/layout.html")
	templates = append(templates, "templates/"+tmpl+".html")

	data["Title"] = title
	data["URL"] = a.generateUrl
	data[csrf.TemplateTag] = csrf.TemplateField(r)

	flashSession, _ := a.sessions.Get(r, "flash")
	flashes := flashSession.Flashes()
	_ = flashSession.Save(r, w)
	data["Flashes"] = flashes

	t, err := template.ParseFiles(templates...)
	if err != nil {
		log.Printf("Error parsing templates: %v", err)
	}
	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
	}
}

var validationErrorMessages = map[string]string{
	"required": "Field %s is required.",
	"min":      "Field %s must have at least %s characters.",
	"max":      "Field %s must have at most %s characters.",
}

func (a *App) parseForm(r *http.Request, dst any) ([]string, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	r.PostForm.Del("gorilla.csrf.Token") // TODO: get name from settings

	err = a.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		return nil, err
	}

	err = a.validator.Struct(dst)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errorMessages []string
			for _, fieldError := range validationErrors {
				var errorMessage string
				if fieldError.Param() != "" {
					errorMessage = fmt.Sprintf(
						validationErrorMessages[fieldError.Tag()],
						fieldError.Field(),
						fieldError.Param(),
					)
				} else {
					errorMessage = fmt.Sprintf(
						validationErrorMessages[fieldError.Tag()],
						fieldError.Field(),
					)
				}
				errorMessages = append(errorMessages, errorMessage)
			}
			return errorMessages, nil
		}
		return nil, err
	}
	return nil, nil
}
