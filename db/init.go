package db

import (
	"database/sql"
	"fmt"

	// Need to import the postgres driver.
	_ "github.com/lib/pq"
)

var db *sql.DB

// DatabaseFlagOptions defines the flag options for the database.
type DatabaseFlagOptions struct {
	DatabaseUser     string `long:"db_user" description:"The database user." required:"true"`
	DatabasePassword string `long:"db_password" description:"The database user's password." required:"true"`
	DatabaseName     string `long:"db_name" description:"The database name." required:"true"`
}

// Init initializes the database instance.
func Init(options *DatabaseFlagOptions) (*sql.DB, error) {
	var err error
	connStr := fmt.Sprintf(
		// TODO: Manage sslmode in prod.
		"postgres://%s:%s@localhost/%s?sslmode=disable",
		options.DatabaseUser,
		options.DatabasePassword,
		options.DatabaseName,
	)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
