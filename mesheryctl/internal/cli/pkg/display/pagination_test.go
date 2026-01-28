package display

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestHandlePaginationAsync(t *testing.T) {

	type items struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type apiResponse struct {
		Code  int     `json:"code"`
		Items []items `json:"items"`
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
		apiResponse      apiResponse
		hasToken         bool
	}{
		{
			name:     "Given_All_Requirements_Met_When_HandlePaginationAsync_Then_Successful_Response",
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
			apiResponse: apiResponse{
				Code: 200,
				Items: []items{
					{ID: "1", Name: "Item1"},
					{ID: "2", Name: "Item2"},
				},
			},
			exceptedResponse: "Total number of items: 2\nPage: 1\n\x1b[1mID\x1b[0m  \x1b[1mNAME \x1b[0m  \n  1   Item1  \n  2   Item2  \n",
			expectedError:    nil,
			hasToken:         true,
		},
		{
			name:     "Given_All_Requirements_Met_When_HandlePaginationAsync_Then_Successful_Count_Response",
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
			apiResponse: apiResponse{
				Code: 200,
				Items: []items{
					{ID: "1", Name: "Item1"},
					{ID: "2", Name: "Item2"},
				},
			},
			exceptedResponse: "Total number of items: 2\n",
			expectedError:    nil,
			hasToken:         true,
		},
		{
			name:     "Given_All_Requirements_Met_When_HandlePaginationAsync_Then_Successful_Count_Empty_Response",
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
			urlPath: "/test?page=0&pagesize=2",
			apiResponse: apiResponse{
				Code:  200,
				Items: []items{},
			},
			exceptedResponse: "No items found\n",
			expectedError:    nil,
			hasToken:         true,
		},
		{
			name:     "Given_Missing_Token_When_HandlePaginationAsync_Then_Error_Returned",
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
			urlPath: "/test?page=0&pagesize=2",
			apiResponse: apiResponse{
				Code:  400,
				Items: []items{},
			},
			exceptedResponse: "",
			expectedError:    utils.ErrAttachAuthToken(fmt.Errorf("Not Set does not exist")),
			hasToken:         false,
		},
		{
			name:     "Given_Meshery_Server_Is_Not_Reachable_When_HandlePaginationAsync_Then_Error_Returned",
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
			urlPath: "/test?page=0&pagesize=2",
			apiResponse: apiResponse{
				Code:  302,
				Items: []items{},
			},
			exceptedResponse: "",
			expectedError:    utils.ErrInvalidToken(),
			hasToken:         true,
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

			mApiResponse, err := json.Marshal(tt.apiResponse.Items)
			if err != nil {
				t.Fatalf("Failed to marshal API response: %v", err)
			}

			url := testContext.BaseURL + tt.urlPath

			httpmock.RegisterResponder("GET", url,
				httpmock.NewStringResponder(tt.apiResponse.Code, string(mApiResponse)))

			originalStdout := os.Stdout

			// Create a pipe to capture output.
			reader, writer, _ := os.Pipe()
			os.Stdout = writer

			// Restore stdout at the end of the test.
			defer func() {
				os.Stdout = originalStdout
				utils.TokenFlag = "Not Set"
			}()

			_ = utils.SetupMeshkitLoggerTesting(t, false)

			err = HandlePaginationAsync(tt.pageSize, tt.displayData, tt.processDataFunc)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {

				_ = writer.Close()

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
