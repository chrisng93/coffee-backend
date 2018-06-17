package api

import (
	"github.com/gorilla/mux"
	"googlemaps.github.io/maps"

	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/db"
)

var databaseOps *db.DatabaseOps
var yelpClient *yelp.Client
var googleMapsClient *maps.Client

// Init initializes the router instance and sets handlers for the API.
func Init(dbops *db.DatabaseOps, yc *yelp.Client, gm *maps.Client) *mux.Router {
	r := mux.NewRouter()
	databaseOps = dbops
	yelpClient = yc
	googleMapsClient = gm

	r.HandleFunc("/coffee_shop", getAllCoffeeShopsHandler).Methods("GET")
	r.HandleFunc("/coffee_shop/{id}", getSingleCoffeeShopHandler).Methods("GET")
	r.HandleFunc("/isochrone", getIsochronesHandler).Methods("GET")
	return r
}
