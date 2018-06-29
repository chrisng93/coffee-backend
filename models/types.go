package models

import (
	"time"
)

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
	// Keys below are added through the Yelp business details endpoint.
	Photos   []string `json:"photos"`
	Location struct {
		State          string   `json:"state"`
		DisplayAddress []string `json:"display_address"`
	} `json:"location"`
	Price string   `json:"price"`
	Phone string   `json:"phone"`
	Hours []*Hours `json:"hours"`
	// TODO: Add these fields when Instagram added.
	// IsInstagrammable bool `json:"is_instagrammable"`
	// InstagramID *string `json:"instagram_id"`
}

// Hours defines a set of operating hours for a Yelp business.
type Hours struct {
	HoursType string `json:"hours_type"`
	Open      []struct {
		IsOvernight bool   `json:"is_overnight"`
		Start       string `json:"start"`
		End         string `json:"end"`
		Day         int64  `json:"day"`
	} `json:"open"`
	IsOpenNow bool `json:"is_open_now"`
}
