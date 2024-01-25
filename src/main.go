package main

import (
	"awesome_web_app/db"
	"awesome_web_app/handlers"
	"awesome_web_app/settings"
	gorillaHandlers "github.com/gorilla/handlers"
	"log"
	"net/http"
)

func main() {
	// load settings
	config, err := settings.LoadSettings()
	if err != nil {
		log.Fatal(err)
	}

	// connect to database
	dbConn, err := db.Connect(config.DB)
	if err != nil {
		log.Fatalf("Could not connect to database: %v\n", err)
	}

	// create app
	app := handlers.NewApp(config, dbConn)

	// configure handlers
	app.AddHandler("/", "index", &gorillaHandlers.MethodHandler{
		"GET": app.AuthMiddleware(app.IndexPage),
	})

	app.AddHandler("/auth/login/", "login", &gorillaHandlers.MethodHandler{
		"GET": http.HandlerFunc(app.GoogleLogin),
	})
	app.AddHandler("/auth/logout/", "logout", &gorillaHandlers.MethodHandler{
		"POST": http.HandlerFunc(app.Logout),
	})
	app.AddHandler("/auth/callback/", "authCallback", &gorillaHandlers.MethodHandler{
		"GET": http.HandlerFunc(app.GoogleCallback),
	})

	app.AddHandler("/profile/", "profile", &gorillaHandlers.MethodHandler{
		"GET":  app.AuthRequiredMiddleware(app.ProfilePage),
		"POST": app.AuthRequiredMiddleware(app.HandleProfileForm),
	})

	app.AddHandler("/users/", "userList", &gorillaHandlers.MethodHandler{
		"GET": app.AuthRequiredMiddleware(app.UserList),
	})
	app.AddHandler("/users/{username}/", "userPage", &gorillaHandlers.MethodHandler{
		"GET": app.AuthRequiredMiddleware(app.UserPage),
	})

	// start server
	app.Serve()
}
