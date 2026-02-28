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

func TestListAsyncPagination(t *testing.T) {
	type items struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type apiResponse struct {
		Code  int     `json:"code"`
		Items []items `json:"items"`
	}

	tests := []struct {
		name             string
		displayData      DisplayDataAsync
		processDataFunc  listRowBuilder[[]items]
		urlPath          string
		expectedError    error
		exceptedResponse string
		apiResponse      apiResponse
		hasToken         bool
	}{
		{
			name: "Given_All_Requirements_Met_When_ListAsyncPagination_Then_Successful_Response",
			displayData: DisplayDataAsync{
				PageSize:         2,
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
			name: "Given_All_Requirements_Met_When_ListAsyncPagination_Then_Successful_Count_Response",
			displayData: DisplayDataAsync{
				PageSize:         2,
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
			name: "Given_All_Requirements_Met_When_ListAsyncPagination_Then_Successful_Count_Empty_Response",
			displayData: DisplayDataAsync{
				PageSize:         2,
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
			name: "Given_Missing_Token_When_ListAsyncPagination_Then_Error_Returned",
			displayData: DisplayDataAsync{
				PageSize:         2,
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
			name: "Given_Meshery_Server_Is_Not_Reachable_When_ListAsyncPagination_Then_Error_Returned",
			displayData: DisplayDataAsync{
				PageSize:         2,
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
	utils.StartMockery(t)
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

			err = ListAsyncPagination(tt.displayData, tt.processDataFunc)

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

func TestPromptAsyncPagination(t *testing.T) {
	type testItem struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	type promptAPIResponse struct {
		Items []testItem `json:"items"`
	}

	extractItems := func(data *promptAPIResponse) ([]testItem, int64) {
		return data.Items, int64(len(data.Items))
	}

	formatLabel := func(rows []testItem) []string {
		labels := []string{}
		for _, r := range rows {
			labels = append(labels, fmt.Sprintf("%s (%s)", r.Name, r.ID))
		}
		return labels
	}

	tests := []struct {
		name          string
		searchTerm    string
		apiItems      []testItem
		apiStatusCode int
		hasToken      bool
		expectError   bool
		expectItem    *testItem
		errContains   string
	}{
		{
			name:          "Given_Zero_Results_When_PromptAsyncPagination_Then_ErrNotFound",
			searchTerm:    "nonexistent",
			apiItems:      []testItem{},
			apiStatusCode: 200,
			hasToken:      true,
			expectError:   true,
			errContains:   "No model with name 'nonexistent' found",
		},
		{
			name:       "Given_Single_Result_When_PromptAsyncPagination_Then_AutoSelect",
			searchTerm: "istio",
			apiItems: []testItem{
				{ID: "abc-123", Name: "Istio Base"},
			},
			apiStatusCode: 200,
			hasToken:      true,
			expectError:   false,
			expectItem:    &testItem{ID: "abc-123", Name: "Istio Base"},
		},
		{
			name:          "Given_Missing_Token_When_PromptAsyncPagination_Then_AuthError",
			searchTerm:    "istio",
			apiItems:      []testItem{},
			apiStatusCode: 400,
			hasToken:      false,
			expectError:   true,
			errContains:   "Not Set does not exist",
		},
		{
			name:          "Given_Invalid_Token_When_PromptAsyncPagination_Then_InvalidTokenError",
			searchTerm:    "istio",
			apiItems:      []testItem{},
			apiStatusCode: 302,
			hasToken:      true,
			expectError:   true,
		},
	}

	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testContext := utils.NewTestHelper(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasToken {
				utils.TokenFlag = utils.GetToken(t)
			}

			defer func() {
				utils.TokenFlag = "Not Set"
			}()

			mApiResponse, err := json.Marshal(promptAPIResponse{Items: tt.apiItems})
			if err != nil {
				t.Fatalf("Failed to marshal API response: %v", err)
			}

			mockURL := testContext.BaseURL + "/test?page=0&pagesize=10&search=" + tt.searchTerm
			httpmock.RegisterResponder("GET", mockURL,
				httpmock.NewStringResponder(tt.apiStatusCode, string(mApiResponse)))

			_ = utils.SetupMeshkitLoggerTesting(t, false)

			var result testItem
			err = PromptAsyncPagination(
				DisplayDataAsync{
					UrlPath:    "test",
					SearchTerm: tt.searchTerm,
				},
				formatLabel,
				extractItems,
				&result,
			)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.expectItem != nil {
					assert.Equal(t, tt.expectItem.ID, result.ID)
					assert.Equal(t, tt.expectItem.Name, result.Name)
				}
			}
		})
	}
}
