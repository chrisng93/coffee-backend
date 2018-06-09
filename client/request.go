package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GetRequestParams defines the parameters needed to create a GET request.
type GetRequestParams struct {
	BaseURL     *url.URL
	Path        string
	APIKey      string
	QueryParams map[string]string
}

// CreateGetRequest is a helper function for creating a GET request given a set of parameters.
func CreateGetRequest(params *GetRequestParams) (
	*http.Request, error) {
	// Create request.
	rel := &url.URL{Path: params.Path}
	u := params.BaseURL.ResolveReference(rel)
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", params.APIKey))

	// Generate query string.
	q := req.URL.Query()
	for param, value := range params.QueryParams {
		q.Add(param, value)
	}
	req.URL.RawQuery = q.Encode()

	return req, nil
}

// UnmarshalResponseBody gets the body from a response.
func UnmarshalResponseBody(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	err := json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
