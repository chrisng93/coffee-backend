package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrisng93/coffee-backend/api"
	"github.com/chrisng93/coffee-backend/yelp"
	flags "github.com/jessevdk/go-flags"
)

// App-level flag options.
type flagOptions struct {
	Port string `long:"port" description:"The port for the server to run on." default:"8080" required:"false"`
}

// Flag options for app and API clients.
type combinedOptions struct {
	flagOptions
	yelp.FlagOptions
}

var options combinedOptions
var yelpClient *yelp.Client

func main() {
	// Parse flags.
	options = combinedOptions{}
	_, err := flags.Parse(&options)
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	// Initialize Yelp API.
	yelpClient, err = yelp.InitClient(&options.FlagOptions)
	if err != nil {
		log.Fatalf("Error initializing Yelp client: %v", err)
	}

	// TODO: Get rid of this - it's just for testing.
	businesses, yelpErr := yelpClient.SearchBusinesses()
	fmt.Println(businesses, yelpErr)

	// Initialize router.
	router := api.Init()
	err = http.ListenAndServe(fmt.Sprintf(":%s", options.Port), corsMiddleware(router))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers:", "Origin, Content-Type, X-Auth-Token, Authorization")
			w.Header().Set("Content-Type", "application/json")
			return
		}

		next.ServeHTTP(w, r)
	})
}
