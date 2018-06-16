package yelp

import "github.com/chrisng93/coffee-backend/models"

// Business defines the information related to a Yelp business.
type Business struct {
	Name        string `json:"name"`
	Coordinates struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"coordinates"`
	YelpID string `json:"id"`
	// Rating and ReviewCount are used to filter for "good" coffee shops.
	Rating      float64 `json:"rating"`
	ReviewCount int64   `json:"review_count"`
	URL         string  `json:"url"`
	// Keys below are used for the business details endpoint.
	Photos   []string `json:"photos"`
	Location struct {
		DisplayAddress []string `json:"display_address"`
	} `json:"location"`
	Price string          `json:"price"`
	Phone string          `json:"phone"`
	Hours []*models.Hours `json:"hours"`
}
