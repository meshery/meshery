package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

// Generic function to fetch data from Mesehry server needs to be type of meshery data ApiResponse
func Fetch[T any](url string) (*T, error) {
	return makeRequest[T](url, http.MethodGet, nil)
}

// Send a Http request to meshery server from mesheryctl cli
func makeRequest[T any](url string, httpMethod string, body io.Reader) (*T, error) {
	req, err := utils.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	// defers the closing of the response body after its use, ensuring that the resources are properly released.
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResponse T
	err = json.Unmarshal(data, &apiResponse)
	if err != nil {
		return nil, err
	}

	return &apiResponse, nil
}
