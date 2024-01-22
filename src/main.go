package main

import (
	"awesome_web_app/db"
	"awesome_web_app/handlers"
	"awesome_web_app/settings"
	"fmt"
	"github.com/gorilla/mux"
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
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/", app.AuthMiddleware(app.IndexPage))

	r.HandleFunc("/auth/login/", app.GoogleLogin)
	r.HandleFunc("/auth/callback/", app.GoogleCallback)
	r.HandleFunc("/auth/logout/", app.Logout)

	r.HandleFunc("/profile/", app.AuthRequiredMiddleware(app.Profile))

	r.HandleFunc("/users/", app.AuthRequiredMiddleware(app.UserList))
	r.HandleFunc("/users/{username}/", app.AuthRequiredMiddleware(app.UserPage))

	// start server
	log.Printf("Starting server on :%v\n", config.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%v", config.Port), r)
	if err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
