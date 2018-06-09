package yelp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type YelpFlagOptions struct {
	YelpBaseURL  string `long:"yelp_base_url" description:"Yelp Base URL." default:"https://api.yelp.com" required:"false"`
	YelpClientID string `long:"yelp_client_id" description:"Yelp Client ID." default:"" required:"true"`
	YelpAPIKey   string `long:"yelp_api_key" description:"Yelp API Key." default:"" required:"true"`
}

type YelpClient struct {
	baseURL    *url.URL
	clientID   string
	apiKey     string
	httpClient *http.Client
}

type SearchBusinessResponse struct {
	Total      int64      `json:"total"`
	Businesses []Business `json:"businesses"`
}

type Business struct {
	Rating float64 `json:"rating"`
}

func InitClient(flags *YelpFlagOptions) (*YelpClient, error) {
	baseURL, err := url.Parse(flags.YelpBaseURL)
	if err != nil {
		return nil, err
	}
	return &YelpClient{
		baseURL:    baseURL,
		clientID:   flags.YelpClientID,
		apiKey:     flags.YelpAPIKey,
		httpClient: &http.Client{},
	}, nil
}

func (c *YelpClient) SearchBusinesses() ([]Business, error) {
	// Create request.
	rel := &url.URL{Path: "/v3/businesses/search"}
	u := c.baseURL.ResolveReference(rel)
	fmt.Println(u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// Generate query string.
	q := req.URL.Query()
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
