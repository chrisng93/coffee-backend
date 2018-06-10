package models

import "time"

// Coordinates defines lat/lng coordinates.
type Coordinates struct {
	Latitude  float64
	Longitude float64
}

// CoffeeShop maps to the coffeeshop.shop data model.
type CoffeeShop struct {
	ID                int64
	LastUpdated       time.Time
	Name              string
	Coordinates       *Coordinates
	YelpID            string
	YelpURL           string
	HasGoodCoffee     bool
	IsGoodForStudying bool
	// TODO: Add these fields when Instagram added.
	// IsInstagrammable bool
	// InstagramID *string
}
