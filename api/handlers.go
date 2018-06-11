package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/chrisng93/coffee-backend/db"
)

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

func getSingleCoffeeShopHandler(w http.ResponseWriter, r *http.Request, databaseOps *db.DatabaseOps) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	coffeeShop, err := databaseOps.GetCoffeeShop(id)
	if coffeeShop.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Coffee shop not found"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffeeShop)
}
