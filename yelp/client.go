package yelp

import (
	"encoding/json"
	"fmt"
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

// SearchBusinessResponse defines the response from calling Yelp's /v3/businesses/search endpoint.
type SearchBusinessResponse struct {
	Total      int64      `json:"total"`
	Businesses []Business `json:"businesses"`
}

// Business defines the information related to a Yelp business.
type Business struct {
	Rating float64 `json:"rating"`
}

// SearchBusinesses calls Yelp's /v3/businesses/search endpoint to get a list of businesses.
func (c *Client) SearchBusinesses() ([]Business, error) {
	// Create request.
	rel := &url.URL{Path: "/v3/businesses/search"}
	u := c.baseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Generate query string.
	q := req.URL.Query()
	// TODO: Deal with pagination for large responses - offset + limit query params.
	// TODO: Add search filter and other params.
	q.Add("location", "Manhattan")
	req.URL.RawQuery = q.Encode()

	// Send request and decode body.
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchBusinessResponse SearchBusinessResponse
	err = json.NewDecoder(resp.Body).Decode(&searchBusinessResponse)
	if err != nil {
		return nil, err
	}
	return searchBusinessResponse.Businesses, nil
}
