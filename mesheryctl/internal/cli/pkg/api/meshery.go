package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Generic function to fetch data from Mesehry server needs to be type of meshery data ApiResponse
func Fetch[T any](url string) (*T, error) {
	resp, err := makeRequest(url, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	return generateDataFromBodyResponse[T](resp)
}

func Delete(url string) (*http.Response, error) {
	return makeRequest(url, http.MethodDelete, nil, nil)
}

// Add sends a POST request to the given URL path with the provided body and optional headers.
// headers may be nil. Header keys/values will be added to the http.Request before dispatch.
func Add(url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return makeRequest(url, http.MethodPost, body, headers)
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
func makeRequest(urlPath string, httpMethod string, body io.Reader, headers map[string]string) (*http.Response, error) {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return nil, utils.ErrLoadConfig(err)
	}

	baseUrl := mctlCfg.GetBaseMesheryURL()

	req, err := utils.NewRequest(httpMethod, fmt.Sprintf("%s/%s", baseUrl, urlPath), body)
	if err != nil {
		return nil, err
	}

	// Add any provided headers to the request. This is important for callers
	// that construct bodies externally (e.g., multipart form) and need to
	// preserve Content-Type or other headers.
for k, v := range headers {
		if http.CanonicalHeaderKey(k) == "Content-Type" {
			// Ensure we set Content-Type directly (replace any existing)
			req.Header.Set(k, v)
		} else {
			req.Header.Add(k, v)
		}
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		if meshkitErr, ok := err.(*mErrors.Error); ok {
			if meshkitErr.Code == utils.ErrFailRequestCode {
				return nil, utils.ErrFailRequest(errors.New("Request failed.\nEnsure meshery server is available."))
			}
		}
		return nil, err
	}

	return resp, nil
}
