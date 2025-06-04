package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

// Generic function to fetch data from Mesehry server needs to be type of meshery data ApiResponse
func Fetch[T any](url string) (*T, error) {
	resp, err := makeRequest(url, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}
	return generateDataFromBodyResponse[T](resp)
}

func Delete(url string) (*http.Response, error) {
	return makeRequest(url, http.MethodDelete, nil)
}

func generateDataFromBodyResponse[T any](response *http.Response) (*T, error) {
	// defers the closing of the response body after its use, ensuring that the resources are properly released.
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var apiResult T
	err = json.Unmarshal(data, &apiResult)
	if err != nil {
		return nil, err
	}

	return &apiResult, nil
}

// Send a Http request to meshery server from mesheryctl cli
func makeRequest(urlPath string, httpMethod string, body io.Reader) (*http.Response, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, utils.ErrLoadConfig(err)
	}

	baseUrl := mctlCfg.GetBaseMesheryURL()

	req, err := utils.NewRequest(httpMethod, fmt.Sprintf("%s/%s", baseUrl, urlPath), body)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
