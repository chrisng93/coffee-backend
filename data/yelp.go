package data

import (
	"fmt"

	"github.com/chrisng93/coffee-backend/yelp"
)

func getYelpData(yelpClient *yelp.Client) {
	// TODO: Find keywords for searches, call Yelp's API, and insert/update in database.
	businesses, yelpErr := yelpClient.SearchBusinesses()
	fmt.Println(businesses, yelpErr)
}
