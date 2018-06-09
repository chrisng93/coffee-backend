package yelp

import (
	"net/http"
	"net/url"
)

// FlagOptions defines the flag options for the Yelp Client.
type FlagOptions struct {
	YelpBaseURL  string `long:"yelp_base_url" description:"Yelp Base URL." default:"https://api.yelp.com" required:"false"`
	YelpClientID string `long:"yelp_client_id" description:"Yelp Client ID." default:"" required:"true"`
	YelpAPIKey   string `long:"yelp_api_key" description:"Yelp API Key." default:"" required:"true"`
}

// Client contains information needed to make calls to Yelp's API and an HTTP client to actually
// call Yelp's API.
type Client struct {
	baseURL    *url.URL
	clientID   string
	apiKey     string
	httpClient *http.Client
}

// InitClient initializes a Yelp Client given a set of flags.
func InitClient(flags *FlagOptions) (*Client, error) {
	baseURL, err := url.Parse(flags.YelpBaseURL)
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:    baseURL,
		clientID:   flags.YelpClientID,
		apiKey:     flags.YelpAPIKey,
		httpClient: &http.Client{},
	}, nil
}

// TODO: Create cron job to call Yelp's API every day to businesses.
