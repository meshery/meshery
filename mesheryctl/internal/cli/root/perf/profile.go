// Copyright Meshery Authors
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

package perf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	termbox "github.com/nsf/termbox-go"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	pageSize          = 25
	viewSingleProfile bool
)

var profileCmd = &cobra.Command{
	Use:   "profile [profile-name]",
	Short: "List performance profiles",
	Long:  `List all the available performance profiles`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// List performance profiles (maximum 25 profiles)
mesheryctl perf profile

// List performance profiles with search (maximum 25 profiles)
mesheryctl perf profile test 2

// View single performance profile with detailed information
mesheryctl perf profile test --view
`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Check for valid output Format
		if outputFormatFlag != "" {
			outputFormatFlag = strings.ToLower(outputFormatFlag)
			if !slices.Contains(validOutputFormats, outputFormatFlag) {
				return utils.ErrInvalidArgument(fmt.Errorf(invalidOutputFormatMsg, outputFormatFlag))
			}
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// used for searching performance profile
		var searchString string
		// setting up for error formatting
		cmdUsed = "profile"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		// handles spaces in args if quoted args passed
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, " ", "%20")
		}
		// Merge args to get profile-name
		searchString = strings.Join(args, "%20")

		profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), searchString, pageSize, pageNumber-1)
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		if len(profiles) == 0 {
			utils.Log.Info("No Performance Profiles to display")
			return nil
		}

		// get profiles as string arrays for printing tabular format profiles
		data := profilesToStringArrays(profiles)

		// print in json/yaml format
		if outputFormatFlag != "" {
			outputFormatterFactory := display.OutputFormatterFactory[[]models.PerformanceProfile]{}
			outputFormatter, err := outputFormatterFactory.New(outputFormatFlag, profiles)
			if err != nil {
				return err
			}

			outputFormatter.Display()
		} else if !viewSingleProfile { // print all profiles
			utils.PrintToTable([]string{"Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data, nil)
		} else { // print single profile
			index := 0
			// if profiles more than one profile, ask for profile index
			if len(profiles) > 1 {
				index, err = userPrompt("profile", "Enter index of the profile", data)
				if err != nil {
					return ErrUserPrompt(err)
				}
			}

			a := profiles[index]

			fmt.Printf("Name: %v\n", a.Name)
			fmt.Printf("ID: %s\n", a.ID.String())
			fmt.Printf("Total Results: %d\n", a.TotalResults)
			fmt.Printf("Endpoint: %v\n", a.Endpoints[0])
			fmt.Printf("Load Generators: %v\n", a.LoadGenerators[0])
			fmt.Printf("Test run duration: %v\n", a.Duration)
			fmt.Printf("QPS: %d\n", a.QPS)
			fmt.Printf("Infrastructure: %v\n", a.ServiceMesh)
			if a.LastRun != nil {
				fmt.Printf("Last Run: %v\n", a.LastRun.Time.Format("2006-01-02 15:04:05"))
			} else {
				fmt.Printf("Last Run: %v\n", "nil")
			}

			if _, ok := a.Metadata["additional_options"]; ok {
				var out bytes.Buffer
				err := json.Indent(&out, []byte(a.Metadata["additional_options"].(string)), "", "  ")
				if err != nil {
					return ErrFailMarshal(err)
				}
				fmt.Printf("Load generator options:\n%s\n", out.String())
			}
		}

		return nil
	},
}

// Fetch performance profiles
func fetchPerformanceProfiles(baseURL, searchString string, pageSize, pageNumber int) ([]models.PerformanceProfile, []byte, error) {
	var response *models.PerformanceProfilesAPIResponse

	url := baseURL + "/api/user/performance/profiles"

	// update the url
	url = fmt.Sprintf("%s?pagesize=%d&page=%d", url, pageSize, pageNumber)
	if searchString != "" {
		url = url + "&search=" + searchString
	}

	// utils.Log.Debug(url)

	req, err := utils.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, utils.ErrReadResponseBody(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, ErrFailUnmarshal(err)
	}

	return response.Profiles, body, nil
}

// add profiles as string arrays to print in a tabular format
func profilesToStringArrays(profiles []models.PerformanceProfile) [][]string {
	var data [][]string

	for _, profile := range profiles {
		// adding profile to data for list output
		if profile.LastRun != nil {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], profile.LastRun.Time.Format("2006-01-02 15:04:05")})
		} else {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], ""})
		}
	}

	return data
}

func userPrompt(key string, label string, data [][]string) (int, error) {
	err := termbox.Init()
	if err != nil {
		return -1, err
	}
	for i, a := range data {
		data[i] = append([]string{strconv.Itoa(i)}, a...)
	}

	if key == "result" {
		utils.PrintToTable([]string{"Index", "Name", "Mesh", "QPS", "Duration", "P50", "P99.9", "Start-Time"}, data, nil)
	} else {
		utils.PrintToTable([]string{"Index", "Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data, nil)
	}

	fmt.Printf("\n")
	validate := func(input string) error {
		index, err := strconv.Atoi(input)
		if err != nil {
			return err
		}
		if index < 0 || index >= len(data) {
			return errors.New("Invalid index")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	if err != nil {
		termbox.Close()
		return -1, fmt.Errorf("prompt failed %v", err)
	}

	termbox.Close()

	index, err := strconv.Atoi(result)
	if err != nil {
		return -1, errors.Wrap(err, "Failed to convert result from prompt")
	}
	return index, nil
}

func init() {
	profileCmd.Flags().BoolVarP(&viewSingleProfile, "view", "", false, "(optional) View single performance profile with more info")
	profileCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "(optional) List next set of performance results with --page (default = 1)")
}
