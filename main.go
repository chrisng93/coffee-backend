package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/chrisng93/coffee-backend/api"
	"github.com/chrisng93/coffee-backend/yelp"
	flags "github.com/jessevdk/go-flags"
)

type flagOptions struct {
	Port string `long:"port" description:"The port for the server to run on." default:"8080" required:"false"`
}

type combinedOptions struct {
	flagOptions
	yelp.YelpFlagOptions
}

var options combinedOptions
var yelpClient *yelp.YelpClient

func main() {
	options = combinedOptions{}
	_, err := flags.Parse(&options)
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	yelpClient, err = yelp.InitClient(&options.YelpFlagOptions)
	if err != nil {
		log.Fatalf("Error initializing Yelp client: %v", err)
	}

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
