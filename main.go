package main

import (
	"fmt"
	"log"
	"net/http"

	flags "github.com/jessevdk/go-flags"
	"googlemaps.github.io/maps"

	"github.com/chrisng93/coffee-backend/api"
	"github.com/chrisng93/coffee-backend/clients/googlemaps"
	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/data"
	"github.com/chrisng93/coffee-backend/db"
)

// App-level flag options.
type flagOptions struct {
	Port string `long:"port" description:"The port for the server to run on." default:"8080" required:"false"`
}

// Flag options for app and API clients.
type combinedOptions struct {
	flagOptions
	db.DatabaseFlagOptions
	yelp.YelpFlagOptions
	googlemaps.GoogleMapsFlagOptions
}

var options combinedOptions
var yelpClient *yelp.Client
var googleMapsClient *maps.Client

func main() {
	// Parse flags.
	options = combinedOptions{}
	_, err := flags.Parse(&options)
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Initialize database.
	databaseOps, err := db.Init(&options.DatabaseFlagOptions)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Initialize Yelp API.
	yelpClient, err = yelp.InitClient(&options.YelpFlagOptions)
	if err != nil {
		log.Fatalf("Error initializing Yelp client: %v", err)
	}

	// Initialize Google Maps API.
	googleMapsClient, err = googlemaps.InitClient(&options.GoogleMapsFlagOptions)
	if err != nil {
		log.Fatalf("Error initializing Google Maps client: %v", err)
	}

	// Start cron jobs for getting data from Yelp and Instagram.
	go data.InitializeCronJobs(databaseOps, yelpClient)

	// Initialize router.
	router := api.Init(databaseOps, yelpClient, googleMapsClient)
	err = http.ListenAndServe(fmt.Sprintf(":%s", options.Port), corsMiddleware(router))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	} else {
		log.Printf("Started server on port %v", options.Port)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
		w.Header().Set("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
