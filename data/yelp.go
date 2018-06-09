package data

import (
	"database/sql"
	"log"

	coffee_db "github.com/chrisng93/coffee-backend/db"
	"github.com/chrisng93/coffee-backend/yelp"
)

// ExcludedCoffeeShops is a list of coffee shops not to include.
var ExcludedCoffeeShops = []string{"Starbucks", "Dunkin' Donuts"}

func getAndUpsertYelpData(db *sql.DB, yelpClient *yelp.Client) error {
	// TODO: Find keywords for searches, call Yelp's API, and insert/update in database.
	bestCoffeeShops, err := yelpClient.SearchBusinesses(&yelp.SearchBusinessesParams{
		Location:   "Lower Manhattan",
		SearchTerm: "best coffee shops",
		Categories: "coffee,coffeeroasteries",
	})
	if err != nil {
		log.Printf("Error getting data from Yelp: %v", err)
	}
	return coffee_db.InsertOrUpdateYelpData(db, filterCoffeeShops(bestCoffeeShops))
}

func filterCoffeeShops(coffeeShops []*yelp.Business) []*yelp.Business {
	var filteredCoffeeShops []*yelp.Business
	for _, coffeeShop := range coffeeShops {
		if includeCoffeeShop(coffeeShop) {
			filteredCoffeeShops = append(filteredCoffeeShops, coffeeShop)
		}
	}
	return filteredCoffeeShops
}

func includeCoffeeShop(coffeeShop *yelp.Business) bool {
	for _, excludedCoffeeShop := range ExcludedCoffeeShops {
		// Don't inlcude coffee shop if it appears in the excluded coffee shop list.
		if coffeeShop.Name == excludedCoffeeShop {
			return false
		}
	}
	// TODO: Do a better job of this. This is very crude.
	return (coffeeShop.Rating > 4 && coffeeShop.ReviewCount > 50) || coffeeShop.ReviewCount > 200
}
