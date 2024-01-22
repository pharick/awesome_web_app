package db

import (
	"awesome_web_app/settings"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func Connect(settings settings.DatabaseSettings) (db *sql.DB, err error) {
	connStr := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		settings.Host,
		settings.Port,
		settings.User,
		settings.Password,
		settings.Database,
	)
	db, err = sql.Open("postgres", connStr)
	return
}
