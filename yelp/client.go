package yelp

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type YelpFlagOptions struct {
	YelpBaseURL  string `long:"yelp_base_url" description:"Yelp Base URL." default:"https://api.yelp.com/v3" required:"false"`
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
	Total      string     `json:"total"`
	Businesses []Business `json:"businesses"`
}

type Business struct {
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
	rel := &url.URL{Path: "/businesses/search"}
	u := c.baseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

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
