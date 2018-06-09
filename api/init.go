package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Init initializes the router instance and sets handlers for the API.
func Init() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/coffee_shop", getAllCoffeeShopsHandler).Methods("GET")
	r.HandleFunc("/coffee_shop/:id", getSingleCoffeeShopHandler).Methods("GET")
	return r
}

func getAllCoffeeShopsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint not implemented."))
}

func getSingleCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Endpoint not implemented."))
}
