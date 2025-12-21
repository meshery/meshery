package api

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/stretchr/testify/assert"
)

type makeRequestTestCase struct {
	name             string
	httpMethod       string
	urlPath          string
	hasToken         bool
	expectedResponse *http.Response
	expectedError    error
}

func TestMakeRequest_Failures(t *testing.T) {
	tests := []makeRequestTestCase{
		{
			name:             "No Token Provided",
			httpMethod:       http.MethodGet,
			urlPath:          "/",
			hasToken:         false,
			expectedResponse: nil,
			expectedError:    utils.ErrAttachAuthToken(fmt.Errorf("Not Set does not exist")),
		},
		{
			name:             "Invalid URL",
			httpMethod:       http.MethodGet,
			urlPath:          "/",
			hasToken:         true,
			expectedResponse: nil,
			expectedError:    utils.ErrFailRequest(fmt.Errorf("parse \"http://%%41:8080/\": invalid URL escape \"%%41\"")),
		},
	}
	utils.SetupContextEnv(t)
	//initialize mock server for handling requests
	utils.StartMockery(t)
	// create a test helper
	testContext := utils.NewTestHelper(t)

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			if tt.hasToken {
				utils.TokenFlag = utils.GetToken(t)
			}

			url := testContext.BaseURL + tt.urlPath

			httpmock.RegisterResponder(tt.httpMethod, url, nil)

			defer func() {
				utils.TokenFlag = "Not Set"
			}()

			_ = utils.SetupMeshkitLoggerTesting(t, false)

			resp, err := makeRequest(url, tt.httpMethod, nil)
			if tt.expectedError != nil {
				assert.Equal(t, reflect.TypeOf(err), reflect.TypeOf(tt.expectedError))
				assert.Equal(t, errors.GetCode(err), errors.GetCode(tt.expectedError))
				assert.Equal(t, errors.GetCause(err), errors.GetCause(tt.expectedError))
				assert.Equal(t, errors.GetCause(err), errors.GetCause(tt.expectedError))
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.Equal(t, tt.expectedResponse, resp)
			}

		})
	}
}
