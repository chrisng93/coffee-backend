package api

import (
	"net/http"

	"github.com/chrisng93/coffee-backend/db"

	"github.com/gorilla/mux"
)

// Init initializes the router instance and sets handlers for the API.
func Init(databaseOps *db.DatabaseOps) *mux.Router {
	r := mux.NewRouter()
	// TODO: Write higher order function to include databaseOps in handler.
	r.HandleFunc("/coffee_shop", func(w http.ResponseWriter, r *http.Request) { getAllCoffeeShopsHandler(w, r, databaseOps) }).Methods("GET")
	r.HandleFunc("/coffee_shop/{id}", func(w http.ResponseWriter, r *http.Request) { getSingleCoffeeShopHandler(w, r, databaseOps) }).Methods("GET")
	return r
}
