package yelp

import (
	"fmt"
	"strconv"

	"github.com/chrisng93/coffee-backend/util"
)

// MaxResults defines the maximum number of results Yelp allows you to see for one query.
const MaxResults = 1000

const searchBusinessesURL = "/v3/businesses/search"

// SearchBusinessesParams defines the parameters used for calling Yelp's /v3/businesses/search
// endpoint.
type SearchBusinessesParams struct {
	Location   string
	SearchTerm string
	Categories string
}

// searchBusinessesResponse defines the response from calling Yelp's /v3/businesses/search endpoint.
type searchBusinessesResponse struct {
	Total      int64       `json:"total"`
	Businesses []*Business `json:"businesses"`
}

// SearchBusinesses calls Yelp's /v3/businesses/search endpoint to get a list of businesses.
func (c *Client) SearchBusinesses(params *SearchBusinessesParams) ([]*Business, error) {
	limit := int64(50)
	var numTries int64
	// Default number of total businesses to MaxResults - this will change once we get a response
	// back from the Yelp API with the actual number of total businesses in the query.
	numTotalBusinesses := int64(MaxResults)
	var businesses []*Business

	for (int64(len(businesses)) < numTotalBusinesses || numTries != 0) && numTries*limit < MaxResults {
		req, err := util.CreateGetRequest(&util.GetRequestParams{
			BaseURL: c.baseURL,
			Path:    searchBusinessesURL,
			APIKey:  c.apiKey,
			QueryParams: map[string]string{
				"limit":    strconv.FormatInt(limit, 10),
				"offset":   strconv.FormatInt(numTries*limit, 10),
				"location": params.Location,
				"term":     params.SearchTerm,
			},
		})
		if err != nil {
			return nil, err
		}

		// Send request and decode body.
		resp, err := c.httpClient.Do(req)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("Error calling Yelp API: %v", resp.Status)
		}
		if err != nil {
			return nil, err
		}
		var searchBusinessesResponse searchBusinessesResponse
		err = util.UnmarshalResponseBody(resp, &searchBusinessesResponse)
		if err != nil {
			return nil, err
		}

		numTotalBusinesses = searchBusinessesResponse.Total
		businesses = append(businesses, searchBusinessesResponse.Businesses...)
		numTries++
	}

	return businesses, nil
}
