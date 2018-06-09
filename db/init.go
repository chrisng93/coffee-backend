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
	DatabaseHost     string `long:"db_host" description:"The database host." default:"localhost" required:"false"`
	DatabasePort     int64  `long:"db_port" description:"The database port." default:"5432" required:"false"`
	DatabaseUser     string `long:"db_user" description:"The database user." required:"true"`
	DatabasePassword string `long:"db_password" description:"The database user's password." required:"true"`
	DatabaseName     string `long:"db_name" description:"The database name." default:"postgres" required:"false"`
}

// Init initializes the database instance.
func Init(options *DatabaseFlagOptions) (*sql.DB, error) {
	var err error
	connStr := fmt.Sprintf(
		// TODO: Manage sslmode in prod.
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		options.DatabaseHost,
		options.DatabasePort,
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
