package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/chrisng93/coffee-backend/clients/yelp"
	"github.com/chrisng93/coffee-backend/models"

	"github.com/gorilla/mux"
)

// DefaultWalkingTimeMin is the default walking time in minutes.
const DefaultWalkingTimeMin = 5

func getAllCoffeeShopsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement pagination and query params once we get more coffee shops.
	coffeeShops, err := databaseOps.GetCoffeeShops()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffeeShops)
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

func getSingleCoffeeShopHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	coffeeShop, err := databaseOps.GetCoffeeShop(id)
	if coffeeShop.ID == 0 {
		http.Error(w, "Coffee shop not found", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	additionalYelpDetails, err := yelpClient.GetBusinessDetails(coffeeShop.YelpID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	addAdditonalYelpDetails(coffeeShop, additionalYelpDetails)

	// TODO: Fetch more info from Instagram API.

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coffeeShop)
}

func getOriginFromRaw(originRaw string) (*models.Coordinates, error) {
	split := strings.Split(originRaw, ",")
	lat, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return nil, err
	}
	lng, err := strconv.ParseFloat(split[1], 64)
	if err != nil {
		return nil, err
	}
	return &models.Coordinates{
		Latitude:  lat,
		Longitude: lng,
	}, nil
}

func getIsochronesHandler(w http.ResponseWriter, r *http.Request) {
	// Format origin.
	originRaw := r.URL.Query().Get("origin")
	if originRaw == "" {
		http.Error(w, "Please set an origin", http.StatusBadRequest)
		return
	}
	origin, err := getOriginFromRaw(originRaw)
	if err != nil {
		http.Error(
			w,
			"Could not parse origin. String should be formatted \"{lat},{lng}\"",
			http.StatusBadRequest,
		)
		return
	}

	// Format travel time.
	walkingTimeMinString := r.URL.Query().Get("walking_time_min")
	if walkingTimeMinString == "" {
		walkingTimeMinString = fmt.Sprintf("%v", DefaultWalkingTimeMin)
	}
	walkingTimeMin, err := strconv.ParseInt(walkingTimeMinString, 10, 64)
	if err != nil {
		http.Error(w, "Please set a valid travel time", http.StatusBadRequest)
		return
	}

	iso, err := calculateIsochrones(googleMapsClient, origin, walkingTimeMin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(iso)
}
