package display

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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
				IsPage:           false,
				Header:           []string{"ID", "Name"},
			},
			processDataFunc: func(data *[]items) ([][]string, int64) {
				return [][]string{{"1", "Item1"}, {"2", "Item2"}}, 2
			},
			urlPath: "/test?page=1&page_size=2",
			apiResponse: []items{
				{ID: "1", Name: "Item1"},
				{ID: "2", Name: "Item2"},
			},
			exceptedResponse: "Total number of items:2\nPage: 1\n  \x1b[1mID\x1b[0m  \x1b[1mNAME \x1b[0m  \n  1   Item1  \n  2   Item2  \n",
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
			urlPath: "/test?page=1&page_size=2",
			apiResponse: []items{
				{ID: "1", Name: "Item1"},
				{ID: "2", Name: "Item2"},
			},
			exceptedResponse: "Total number of items:2\n",
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

			httpmock.RegisterResponder("GET", testContext.BaseURL+tt.urlPath,
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

				assert.Equal(t, tt.exceptedResponse, output)
				assert.NoError(t, err)
			}

		})
	}
}
