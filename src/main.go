package main

import (
	"awesome_web_app/db"
	"awesome_web_app/handlers"
	"awesome_web_app/settings"
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
	app.AddHandler("/", "index", map[string]http.HandlerFunc{
		"GET": app.AuthMiddleware(app.IndexPage),
	})

	app.AddHandler("/auth/login/", "login", map[string]http.HandlerFunc{
		"GET": app.GoogleLogin,
	})
	app.AddHandler("/auth/logout/", "logout", map[string]http.HandlerFunc{
		"POST": app.Logout,
	})
	app.AddHandler("/auth/callback/", "authCallback", map[string]http.HandlerFunc{
		"GET": app.GoogleCallback,
	})

	app.AddHandler("/profile/", "profile", map[string]http.HandlerFunc{
		"GET":  app.AuthRequiredMiddleware(app.ProfilePage),
		"POST": app.AuthRequiredMiddleware(app.HandleProfileForm),
	})

	app.AddHandler("/users/", "userList", map[string]http.HandlerFunc{
		"GET": app.AuthRequiredMiddleware(app.UserList),
	})
	app.AddHandler("/users/{username}/", "userPage", map[string]http.HandlerFunc{
		"GET": app.AuthRequiredMiddleware(app.UserPage),
	})

	// start server
	app.Serve()
}
