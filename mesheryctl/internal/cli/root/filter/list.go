// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pageSize   = 25
	pageNumber int
	verbose    bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List filters",
	Long:  `Display list of all available filter files.`,
	Example: `
// List all WASM filter files present
mesheryctl filter list	(maximum 25 filters)

// Search for filter
mesheryctl filter list Test (maximum 25 filters)

// Search for filter with space
mesheryctl filter list 'Test Filter' (maximum 25 filters)
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}
		var searchString string
		if len(args) > 0 {
			searchString = strings.ReplaceAll(args[0], " ", "%20")
		}
		response, err := fetchFilters(mctlCfg.GetBaseMesheryURL(), searchString, pageSize, pageNumber-1)
		if err != nil {
			utils.Log.Error(ErrFetchFilter(err))
			return nil
		}

		if len(args) > 0 && len(response.Filters) == 0 {
			utils.Log.Info("No WASM Filter to display with name :", strings.Join(args, " "))
			return nil
		} else if len(response.Filters) == 0 {
			utils.Log.Info("No WASM Filter to display")
			return nil
		}

		tokenObj, err := utils.ReadToken(utils.TokenFlag)
		if err != nil {
			return errors.New(utils.FilterListError("error reading token\nUse 'mesheryctl filter list --help' to display usage guide\n" + err.Error()))
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		// Initialize pagination variables
		startIndex := 0
		endIndex := min(pageSize, len(data))
		totalFilters := int(response.TotalCount)

		for {
			// Clear the terminal screen and print pagination details
			utils.ClearLine()
			fmt.Printf("Total number of filters: %d\n", totalFilters)
			fmt.Printf("Page: %d of %d\n", (startIndex/pageSize)+1, (totalFilters+pageSize-1)/pageSize)

			// Display the current page of filters
			displayCurrentPage(data[startIndex:endIndex], verbose)

			// Check if we've displayed all filters
			if endIndex == totalFilters {
				break
			}

			// Wait for user input to navigate pages
			fmt.Println("Press Enter or ↓ to continue, Esc or Ctrl+C to exit")

			keysEvents, err := keyboard.GetKeys(10)
			if err != nil {
				return err
			}
			defer func() {
				_ = keyboard.Close()
			}()

			event := <-keysEvents
			if event.Err != nil {
				utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
				break
			}

			// Handle user input
			if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
				break
			} else if event.Key == keyboard.KeyEnter || event.Key == keyboard.KeyArrowDown {
				startIndex += pageSize
				endIndex = min(startIndex+pageSize, totalFilters)
				pageNumber = startIndex/pageSize + 1

				// Fetch the next set of filters
				response, err = fetchFilters(mctlCfg.GetBaseMesheryURL(), searchString, pageSize, pageNumber-1)
				if err != nil {
					utils.Log.Error(ErrFetchFilter(err))
					return err
				}
				data = processData(response.Filters, verbose, provider)
			}
		}

		return nil
	},
}

func displayCurrentPage(filters [][]string, verbose bool) {
	if verbose {
		utils.PrintToTableWithFooter([]string{"FILTER ID", "USER ID", "NAME", "CREATED", "UPDATED"}, filters, []string{"Total", fmt.Sprintf("%d", len(filters)), "", "", ""})
	} else {
		utils.PrintToTableWithFooter([]string{"FILTER ID", "NAME", "CREATED", "UPDATED"}, filters, []string{"Total", fmt.Sprintf("%d", len(filters)), "", ""})
	}
}

func processData(filters []models.MesheryFilter, verbose bool, provider string) [][]string {
	var data [][]string

	// Iterate through the filters and format the data based on the verbose flag and provider
	for _, v := range filters {
		// Define variables for the fields you want to display
		filterID := v.ID.String()
		filterName := v.Name
		createdAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", v.CreatedAt.Year(), v.CreatedAt.Month(), v.CreatedAt.Day(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
		updatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", v.UpdatedAt.Year(), v.UpdatedAt.Month(), v.UpdatedAt.Day(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())

		// Depending on the verbose flag, display additional information
		if verbose {
			if provider == "None" {
				data = append(data, []string{filterID, filterName, createdAt, updatedAt})
			} else {
				userID := "null"
				if v.UserID != nil {
					userID = *v.UserID
				}
				data = append(data, []string{filterID, userID, filterName, createdAt, updatedAt})
			}
		} else {
			filterNameTrimmed := strings.TrimSuffix(filterName, filepath.Ext(filterName))
			filterIDShort := utils.TruncateID(filterID) // Assuming TruncateID is a utility function that shortens the ID
			data = append(data, []string{filterIDShort, filterNameTrimmed, createdAt, updatedAt})
		}
	}

	return data
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// Pagination(making multiple requests) to retrieve filter Data in batches
func fetchFilters(baseURL, searchString string, pageSize, pageNumber int) (*models.FiltersAPIResponse, error) {
	var response *models.FiltersAPIResponse

	url := baseURL + "/api/filter"

	url = fmt.Sprintf("%s?pagesize=%d&page=%d", url, pageSize, pageNumber)
	if searchString != "" {
		url = url + "&search=" + searchString
	}

	utils.Log.Debug(url)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, utils.ErrCreatingRequest(err)
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, utils.ErrRequestResponse(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.FilterListError("failed to read response body"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Server returned with status code: " + fmt.Sprint(resp.StatusCode) + "\n" + "Response: " + string(body))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, utils.ErrUnmarshal(err)
	}
	return response, nil
}
func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and filter file identifiers")
	listCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "(optional) List next set of filters with --page (default = 1)")
}
