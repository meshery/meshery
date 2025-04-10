package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

// Generic function to fetch data from Mesehry server needs to be type of meshery data ApiResponse
func Fetch[T any](url string) (*T, error) {
	return makeRequest[T](url, http.MethodGet, nil)
}

// Send a Http request to meshery server from mesheryctl cli
func makeRequest[T any](urlPath string, httpMethod string, body io.Reader) (*T, error) {
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
