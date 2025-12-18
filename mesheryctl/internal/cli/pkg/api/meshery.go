package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
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

func Add(url string, body io.Reader) (*http.Response, error) {
	return makeRequest(url, http.MethodPost, body)
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
		if meshkitErr, ok := err.(*mErrors.Error); ok {
			if meshkitErr.Code == utils.ErrFailRequestCode {
				endpoint := mctlCfg.Contexts[mctlCfg.CurrentContext].Endpoint
				errCtx := fmt.Sprintf("Unable to connect to Meshery server at %s (current context).", endpoint)
				failedReqErr := utils.ErrFailRequest(fmt.Errorf("%s", errCtx))
				errRemediation := mErrors.GetRemedy(failedReqErr)
				return nil, utils.ErrFailRequest(fmt.Errorf("%s\n%s\n%s", errCtx, errRemediation, generateErrorReferenceDetails("ErrFailRequestCode", utils.ErrFailRequestCode)))
			}
		}
		return nil, err
	}

	return resp, nil
}

func generateErrorReferenceDetails(referenceCodeName, code string) string {
	codeNumber := strings.Split(code, "-")[1]
	return fmt.Sprintf("\nFor additional details see https://docs.meshery.io/reference/error-codes#%s-%s", referenceCodeName, codeNumber)
}
