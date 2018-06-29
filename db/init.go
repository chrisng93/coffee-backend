package db

import (
	"database/sql"
	"fmt"
	"time"

	// Need to import the postgres driver.
	_ "github.com/lib/pq"
)

// MaxPings is the maximum amount of pings to the database before we should error out.
const MaxPings = 30

// PingIntervalSec is the interval with which we should ping the database to check for connection.
const PingIntervalSec = 1

var db *sql.DB

// DatabaseFlagOptions defines the flag options for the database.
type DatabaseFlagOptions struct {
	DatabaseHost     string `long:"db_host" description:"The database host." default:"127.0.0.1" required:"false"`
	DatabasePort     int64  `long:"db_port" description:"The database port." default:"5432" required:"false"`
	DatabaseUser     string `long:"db_user" description:"The database user." required:"true"`
	DatabasePassword string `long:"db_password" description:"The database user's password." required:"true"`
	DatabaseName     string `long:"db_name" description:"The database name." default:"postgres" required:"false"`
}

// DatabaseOps is a struct that wraps the database object. Database operations can *only* be used
// elsewhere in the app via this struct.
type DatabaseOps struct {
	db *sql.DB
}

// Init initializes the database instance.
func Init(options *DatabaseFlagOptions) (*DatabaseOps, error) {
	connStr := fmt.Sprintf(
		// TODO: Manage sslmode in prod.
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		options.DatabaseHost,
		options.DatabasePort,
		options.DatabaseUser,
		options.DatabasePassword,
		options.DatabaseName,
	)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = checkConnection(db)
	if err != nil {
		return nil, err
	}

	return &DatabaseOps{db: db}, nil
}

// Ping the database to see if we have an actual connection. sql.Open doesn't actually check to
// see if there's a connection - it only really checks to see if the arguments are valid.
func checkConnection(db *sql.DB) error {
	var err error
	tries := 0
	for tries < MaxPings {
		err = db.Ping()
		if err == nil {
			break
		}
		tries++
		time.Sleep(PingIntervalSec * time.Second)
	}
	return err
}
