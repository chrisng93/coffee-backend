package yelp

import (
	"fmt"

	"github.com/chrisng93/coffee-backend/util"
)

func createGetBusinessDetailsURL(id string) string {
	return fmt.Sprintf("/v3/businesses/%s", id)
}

// GetBusinessDetails calls Yelp's /v3/businesses/{id} endpoint to get more details on a
// specific business.
func (c *Client) GetBusinessDetails(id string) (*Business, error) {
	req, err := util.CreateGetRequest(&util.GetRequestParams{
		BaseURL: c.baseURL,
		Path:    createGetBusinessDetailsURL(id),
		APIKey:  c.apiKey,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Error calling Yelp API: %v", resp.Status)
	}
	if err != nil {
		return nil, err
	}
	var business Business
	err = util.UnmarshalResponseBody(resp, &business)
	if err != nil {
		return nil, err
	}

	return &business, nil
}
