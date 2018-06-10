package models

import "time"

// Coordinates defines lat/lng coordinates.
type Coordinates struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// CoffeeShop maps to the coffeeshop.shop data model.
type CoffeeShop struct {
	ID                int64        `json:"id"`
	LastUpdated       time.Time    `json:"last_updated"`
	Name              string       `json:"name"`
	Coordinates       *Coordinates `json:"coordinates"`
	YelpID            string       `json:"yelp_id"`
	YelpURL           string       `json:"yelp_url"`
	HasGoodCoffee     bool         `json:"has_good_coffee"`
	IsGoodForStudying bool         `json:"is_good_for_studying"`
	// TODO: Add these fields when Instagram added.
	// IsInstagrammable bool `json:"is_instagrammable"`
	// InstagramID *string `json:"instagram_id"`
}
