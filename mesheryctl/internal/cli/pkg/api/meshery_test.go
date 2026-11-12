package api

import (
	"fmt"
	"net/http"
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

	testContext := utils.InitTestEnvironment(t)

	tests := []makeRequestTestCase{
		{
			name:             "Given no token provided when received a GET request then an error is returned",
			httpMethod:       http.MethodGet,
			urlPath:          "/",
			hasToken:         false,
			expectedResponse: nil,
			expectedError:    utils.ErrAttachAuthToken(fmt.Errorf("Not Set does not exist")),
		},
		{
			name:             "Given server not reachable when received a GET request then an error is returned",
			httpMethod:       http.MethodGet,
			urlPath:          "/",
			hasToken:         true,
			expectedResponse: nil,
			expectedError: func() error {
				errCtx := fmt.Sprintf("Unable to connect to Meshery server at %s (current context).", testContext.BaseURL)
				failedReqErr := utils.ErrFailRequest(fmt.Errorf("%s", errCtx))
				errRemediation := errors.GetRemedy(failedReqErr)
				return utils.ErrFailRequest(fmt.Errorf("%s\n%s\n%s", errCtx, errRemediation, generateErrorReferenceDetails("ErrFailRequestCode", utils.ErrFailRequestCode)))
			}(),
		},
	}

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

			resp, err := makeRequest(url, tt.httpMethod, nil, nil)
			if tt.expectedError != nil {
				utils.AssertMeshkitErrorsEqual(t, tt.expectedError, err)
			} else {
				assert.Equal(t, tt.expectedResponse, resp)
			}

		})
	}
}
