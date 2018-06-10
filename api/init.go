package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/chrisng93/coffee-backend/db"
)

// Init initializes the router instance and sets handlers for the API.
func Init(databaseOps *db.DatabaseOps) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/coffee_shop", func(w http.ResponseWriter, r *http.Request) { getAllCoffeeShopsHandler(w, r, databaseOps) }).Methods("GET")
	r.HandleFunc("/coffee_shop/{id}", getSingleCoffeeShopHandler).Methods("GET")
	return r
}

func getAllCoffeeShopsHandler(w http.ResponseWriter, r *http.Request, databaseOps *db.DatabaseOps) {
	// TODO: Implement pagination and query params once we get more coffee shops.
	coffeeShops, err := databaseOps.GetCoffeeShops()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffeeShops)
}

func getSingleCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint not implemented."))
}
