package main

import (
	"awesome_web_app/db"
	"awesome_web_app/handlers"
	"awesome_web_app/settings"
	"log"
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
	app.AddHandler("/", app.AuthMiddleware(app.IndexPage), "index")

	app.AddHandler("/auth/login/", app.GoogleLogin, "login")
	app.AddHandler("/auth/logout/", app.Logout, "logout")
	app.AddHandler("/auth/callback/", app.GoogleCallback, "auth_callback")

	app.AddHandler("/profile/", app.AuthRequiredMiddleware(app.Profile), "profile")

	app.AddHandler("/users/", app.AuthRequiredMiddleware(app.UserList), "user_list")
	app.AddHandler("/users/{username}/", app.AuthRequiredMiddleware(app.UserPage), "user_page")

	// start server
	app.Serve()
}
