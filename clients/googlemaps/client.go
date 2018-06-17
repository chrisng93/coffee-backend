package googlemaps

import "googlemaps.github.io/maps"

// GoogleMapsFlagOptions defines the flag options for the GoogleMaps Client.
type GoogleMapsFlagOptions struct {
	GoogleMapsAPIKey string `long:"google_maps_api_key" description:"Google Maps API Key." default:"" required:"true"`
}

// InitClient initializes the Google Maps client.
func InitClient(flags *GoogleMapsFlagOptions) (*maps.Client, error) {
	return maps.NewClient(maps.WithAPIKey(flags.GoogleMapsAPIKey))
}
