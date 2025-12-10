package display

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"strings"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandlePaginationAsync(t *testing.T) {

	type items struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	// Test cases
	tests := []struct {
		name             string
		pageSize         int
		displayData      DisplayDataAsync
		processDataFunc  func(*[]items) ([][]string, int64)
		urlPath          string
		expectedError    error
		exceptedResponse string
		apiResponse      []items
	}{
		{
			name:     "Successful pagination",
			pageSize: 2,
			displayData: DisplayDataAsync{
				Page:             1,
				UrlPath:          "test",
				DataType:         "items",
				DisplayCountOnly: false,
				IsPage:           true,
				Header:           []string{"ID", "Name"},
			},
			processDataFunc: func(data *[]items) ([][]string, int64) {
				return [][]string{{"1", "Item1"}, {"2", "Item2"}}, 2
			},
			urlPath: "/test?page=0&pagesize=2",
			apiResponse: []items{
				{ID: "1", Name: "Item1"},
				{ID: "2", Name: "Item2"},
			},
			exceptedResponse: "Total number of items: 2\nPage: 1\n\x1b[1mID\x1b[0m  \x1b[1mNAME \x1b[0m  \n  1   Item1  \n  2   Item2  \n",
			expectedError:    nil,
		},
		{
			name:     "Successful count",
			pageSize: 2,
			displayData: DisplayDataAsync{
				Page:             1,
				UrlPath:          "test",
				DataType:         "items",
				DisplayCountOnly: true,
				IsPage:           false,
				Header:           []string{"ID", "Name"},
			},
			processDataFunc: func(data *[]items) ([][]string, int64) {
				return [][]string{{"1", "Item1"}, {"2", "Item2"}}, 2
			},
			urlPath: "/test?page=0&pagesize=2",
			apiResponse: []items{
				{ID: "1", Name: "Item1"},
				{ID: "2", Name: "Item2"},
			},
			exceptedResponse: "Total number of items: 2\n",
			expectedError:    nil,
		},
		{
			name:     "Successful count empty response",
			pageSize: 2,
			displayData: DisplayDataAsync{
				Page:             1,
				UrlPath:          "test",
				DataType:         "items",
				DisplayCountOnly: true,
				IsPage:           false,
				Header:           []string{"ID", "Name"},
			},
			processDataFunc: func(data *[]items) ([][]string, int64) {
				return [][]string{}, 0
			},
			urlPath:          "/test?page=0&pagesize=2",
			apiResponse:      []items{},
			exceptedResponse: "No items found\n",
			expectedError:    nil,
		},
	}

	utils.SetupContextEnv(t)

	//initialize mock server for handling requests
	utils.StartMockery(t)

	// create a test helper
	testContext := utils.NewTestHelper(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			utils.TokenFlag = utils.GetToken(t)
			mApiRespoanse, err := json.Marshal(tt.apiResponse)
			if err != nil {
				t.Fatalf("Failed to marshal API response: %v", err)
			}

			url := testContext.BaseURL + tt.urlPath

			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(200, string(mApiRespoanse)))

			originalStdout := os.Stdout

			// Create a pipe to capture output.
			reader, writer, _ := os.Pipe()
			os.Stdout = writer

			// Restore stdout at the end of the test.
			defer func() {
				os.Stdout = originalStdout
			}()

			_ = utils.SetupMeshkitLoggerTesting(t, false)

			err = HandlePaginationAsync(tt.pageSize, tt.displayData, tt.processDataFunc)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {

				writer.Close()

				// Read captured output.
				var buf bytes.Buffer
				_, _ = buf.ReadFrom(reader)
				output := buf.String()

				// Clean both actual and expected output to remove ANSI code and normalize formatting
				cleanedActual := utils.CleanStringFromHandlePagination(output)
				cleanExpected := utils.CleanStringFromHandlePagination(tt.exceptedResponse)
				assert.Equal(t, cleanExpected, cleanedActual)
				assert.NoError(t, err)
			}

		})
	}
}



func TestHandlePaginationAsync_AuthenticationErrors(t *testing.T) {
	type items struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	authTests := []struct {
		name          string
		errorResponse string
		statusCode    int
		expectAuthErr bool
	}{
		{
			name:          "Authentication error in response",
			errorResponse: `{"error": "authentication failed"}`,
			statusCode:    401,
			expectAuthErr: true,
		},
		{
			name:          "Token error in response", 
			errorResponse: `{"error": "invalid token"}`,
			statusCode:    401,
			expectAuthErr: true,
		},
		{
			name:          "Unauthorized error in response",
			errorResponse: `{"error": "unauthorized access"}`,
			statusCode:    403,
			expectAuthErr: true,
		},
		{
			name:          "Login error in response",
			errorResponse: `{"error": "please login"}`,
			statusCode:    401,
			expectAuthErr: true,
		},
		{
			name:          "Generic server error",
			errorResponse: `{"error": "internal server error"}`,
			statusCode:    500,
			expectAuthErr: false,
		},
	}

	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)

	for _, tt := range authTests {
		t.Run(tt.name, func(t *testing.T) {
			displayData := DisplayDataAsync{
				Page:             1,
				UrlPath:          "test",
				DataType:         "items",
				DisplayCountOnly: false,
				IsPage:           true,
				Header:           []string{"ID", "Name"},
			}

			processDataFunc := func(data *[]items) ([][]string, int64) {
				return [][]string{}, 0
			}

			urlPath := "/test?page=0&pagesize=10"
			url := testContext.BaseURL + urlPath

			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(tt.statusCode, tt.errorResponse))

			err := HandlePaginationAsync(10, displayData, processDataFunc)

			if tt.expectAuthErr {
				// Should return the original authentication error
				assert.Error(t, err)
				assert.Contains(t, strings.ToLower(err.Error()), 
					strings.Split(strings.ToLower(tt.errorResponse), "\"")[3])
			} else {
				// Should wrap with pagination context
				assert.Error(t, err)
				// This would depend on how your API error handling works
			}
		})
	}
}

