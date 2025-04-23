package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

// Generic function to fetch data from Mesehry server needs to be type of meshery data ApiResponse
func Fetch[T any](url string) (*T, error) {

	resp, err := makeRequest(url, http.MethodGet, nil)

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

func Add(url string, payload interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	_, err = makeRequest(url, http.MethodPost, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return nil
}

// Send a Http request to meshery server from mesheryctl cli
func makeRequest(url string, httpMethod string, body io.Reader) (*http.Response, error) {
	req, err := utils.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
