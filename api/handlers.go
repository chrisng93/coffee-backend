package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/models"

	"github.com/gorilla/mux"
)

func getAllCoffeeShopsHandler(w http.ResponseWriter, r *http.Request) {
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

	additionalYelpDetails, err := yelpClient.GetBusinessDetails(coffeeShop.YelpID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	addAdditonalYelpDetails(coffeeShop, additionalYelpDetails)

	// TODO: Fetch more info from Instagram API.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffeeShop)
}

// TODO: Having to convert the Yelp business data to coffee shop model data in the handler
// isn't the cleanest. Figure out a better way to do this.
func addAdditonalYelpDetails(coffeeShop *models.CoffeeShop, yelpBusiness *yelp.Business) {
	coffeeShop.Photos = yelpBusiness.Photos
	coffeeShop.Location = yelpBusiness.Location
	coffeeShop.Price = yelpBusiness.Price
	coffeeShop.Phone = yelpBusiness.Phone
	coffeeShop.Hours = yelpBusiness.Hours
}
