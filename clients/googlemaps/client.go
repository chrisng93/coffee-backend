package googlemaps

import (
	"strings"

	"googlemaps.github.io/maps"
)

// Current key index.
var currKey int

// Maximum key index.
var maxKeys int

// Slice of the Google Maps API keys. We pass in multiple API keys because the Distance Matrix API
// has a low quota.
var keys []string

// GoogleMapsFlagOptions defines the flag options for the GoogleMaps Client.
type GoogleMapsFlagOptions struct {
	GoogleMapsAPIKey string `long:"google_maps_api_key" description:"Google Maps API Key." default:"" required:"true"`
}

// InitClient initializes the Google Maps client.
func InitClient(flags *GoogleMapsFlagOptions) (*maps.Client, error) {
	keys = strings.Split(flags.GoogleMapsAPIKey, ",")
	maxKeys = len(keys)
	return maps.NewClient(maps.WithAPIKey(keys[currKey]))
}

// RotateClient creates a Google Maps client with a new API key to deal with the quota issue.
func RotateClient() (*maps.Client, error) {
	if currKey == maxKeys-1 {
		currKey = 0
	} else {
		currKey++
	}
	return maps.NewClient(maps.WithAPIKey(keys[currKey]))
}
