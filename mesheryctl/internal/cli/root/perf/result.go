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

package perf

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gofrs/uuid"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type resultStruct struct {
	Name          string
	StartTime     *time.Time
	LatenciesMs   *models.LatenciesMs
	QPS           int
	URL           string
	UserID        *uuid.UUID
	Duration      string
	MesheryID     *uuid.UUID
	LoadGenerator string
}

var (
	pageNumber       int
	viewSingleResult bool
)

var linkDocPerfResult = map[string]string{
	"link":    "![perf-result-usage](/assets/img/mesheryctl/perf-result.png)",
	"caption": "Usage of mesheryctl perf result",
}

var resultCmd = &cobra.Command{
	Use:   "result [profile-name]",
	Short: "List performance test results",
	Long:  `List all the available test results of a performance profile`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// List Test results (maximum 25 results)
mesheryctl perf result saturday-profile

// View other set of performance results with --page (maximum 25 results)
mesheryctl perf result saturday-profile --page 2

// View single performance result with detailed information
mesheryctl perf result saturday-profile --view
`,
	Annotations: linkDocPerfResult,
	RunE: func(cmd *cobra.Command, args []string) error {
		// used for searching performance profile
		var searchString, profileID string
		// setting up for error formatting
		cmdUsed = "result"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrMesheryConfig(err)
		}

		// Throw error if a profile name is not provided
		if len(args) == 0 {
			return ErrNoProfileName()
		}

		// handles spaces in args if quoted args passed
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, " ", "%20")
		}
		// Merge args to get profile-name
		searchString = strings.Join(args, "%20")

		profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), searchString, pageSize, pageNumber-1)
		if err != nil {
			return err
		}

		if len(profiles) == 0 {
			utils.Log.Info("No Performance Profiles found with the given name")
			return nil
		}

		data := profilesToStringArrays(profiles)

		if len(profiles) == 1 {
			// found only one profile with matching name
			profileID = data[0][1]
		} else {
			// user prompt to select profile
			selectedProfileIndex, err := userPrompt("profile", "Found multiple profiles with given name, select a profile", data)
			if err != nil {
				return err
			}
			// ids got shifted with 1 in userPrompt()
			profileID = data[selectedProfileIndex][2]
		}

		results, _, err := fetchPerformanceProfileResults(mctlCfg.GetBaseMesheryURL(), profileID, pageSize, pageNumber-1)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			utils.Log.Info("No Test Results to display")
			return nil
		}

		// get performance results in format of string arrays and resultStruct
		data, expandedData := performanceResultsToStringArrays(results)
		if len(expandedData) == 0 {
			utils.Log.Info("No test results to display")
			return nil
		}

		if outputFormatFlag != "" {
			body, _ := json.Marshal(results)
			if outputFormatFlag == "yaml" {
				body, _ = yaml.JSONToYAML(body)
			} else if outputFormatFlag != "json" {
				return ErrInvalidOutputChoice()
			}
			utils.Log.Info(string(body))
		} else if !viewSingleResult { // print all results
			utils.PrintToTable([]string{"NAME", "MESH", "QPS", "DURATION", "P50", "P99.9", "START-TIME"}, data)
		} else {
			index := 0
			// if more than one result exist ask for index
			if len(data) > 1 {
				index, err = userPrompt("result", "Select Performance-test result to expand", data)
				if err != nil {
					return err
				}
			}
			a := expandedData[index]
			fmt.Printf("Name: %v\n", a.Name)
			fmt.Printf("UserID: %s\n", a.UserID.String())
			fmt.Printf("Endpoint: %v\n", a.URL)
			fmt.Printf("QPS: %v\n", a.QPS)
			fmt.Printf("Test run duration: %v\n", a.Duration)
			fmt.Printf("Latencies _ms: Avg: %v, Max: %v, Min: %v, P50: %v, P90: %v, P99: %v\n", a.LatenciesMs.Average, a.LatenciesMs.Max, a.LatenciesMs.Min, a.LatenciesMs.P50, a.LatenciesMs.P90, a.LatenciesMs.P99)
			fmt.Printf("Start Time: %v\n", fmt.Sprintf("%d-%d-%d %d:%d:%d", int(a.StartTime.Month()), a.StartTime.Day(), a.StartTime.Year(), a.StartTime.Hour(), a.StartTime.Minute(), a.StartTime.Second()))
			fmt.Printf("Meshery ID: %v\n", a.MesheryID.String())
			fmt.Printf("Load Generator: %v\n", a.LoadGenerator)
		}
		return nil
	},
}

// Fetch results for a specific profile
func fetchPerformanceProfileResults(baseURL, profileID string, pageSize, pageNumber int) ([]models.PerformanceResult, []byte, error) {
	var response *models.PerformanceResultsAPIResponse

	url := baseURL + "/api/user/performance/profiles/" + profileID + "/results"

	tempURL := fmt.Sprintf("%s?pagesize=%d&page=%d", url, pageSize, pageNumber)

	req, err := utils.NewRequest("GET", tempURL, nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, ErrFailUnmarshal(err)
	}
	return response.Results, body, nil
}

// change performance results into string arrays(for tabular format printing) and profileStruct (to print single performance result)
func performanceResultsToStringArrays(results []models.PerformanceResult) ([][]string, []resultStruct) {
	var data [][]string
	var expendedData []resultStruct

	// append data for single profile
	for _, result := range results {
		serviceMesh := "No Mesh"
		p50 := ""
		p99_9 := ""
		var P50, P90, P99 float64

		if result.Mesh != "" {
			serviceMesh = result.Mesh
		}
		qps := fmt.Sprintf("%d", int(result.RunnerResults.QPS))
		duration := result.RunnerResults.RequestedDuration
		if len(result.RunnerResults.DurationHistogram.Percentiles) > 0 {
			p50 = fmt.Sprintf("%.8f", result.RunnerResults.DurationHistogram.Percentiles[0].Value)
			p99_9 = fmt.Sprintf("%.8f", result.RunnerResults.DurationHistogram.Percentiles[len(result.RunnerResults.DurationHistogram.Percentiles)-1].Value)
		}
		startTime := result.TestStartTime.Format("2006-01-02 15:04:05")
		data = append(data, []string{result.Name, serviceMesh, qps, duration, p50, p99_9, startTime})

		if len(result.RunnerResults.DurationHistogram.Percentiles) > 3 {
			P50 = result.RunnerResults.DurationHistogram.Percentiles[0].Value
			P90 = result.RunnerResults.DurationHistogram.Percentiles[2].Value
			P99 = result.RunnerResults.DurationHistogram.Percentiles[3].Value
		}

		// append data for extended output
		name := "None"
		userid := uuid.Nil
		mesheryid := uuid.Nil
		url := "None"
		loadGenerator := "None"

		if result.UserID != nil {
			userid = *result.UserID
		}

		if result.MesheryID != nil {
			mesheryid = *result.MesheryID
		}

		if result.Name != "" {
			name = result.Name
		}

		if result.RunnerResults.URL != "" {
			url = result.RunnerResults.URL
		}

		if result.RunnerResults.LoadGenerator != "" {
			loadGenerator = result.RunnerResults.LoadGenerator
		}

		a := resultStruct{
			Name:     name,
			UserID:   (*uuid.UUID)(userid.Bytes()),
			URL:      url,
			QPS:      int(result.RunnerResults.QPS),
			Duration: result.RunnerResults.RequestedDuration,
			LatenciesMs: &models.LatenciesMs{
				Average: result.RunnerResults.DurationHistogram.Average,
				Max:     result.RunnerResults.DurationHistogram.Max,
				Min:     result.RunnerResults.DurationHistogram.Min,
				P50:     P50,
				P90:     P90,
				P99:     P99,
			},
			StartTime:     result.TestStartTime,
			MesheryID:     (*uuid.UUID)(mesheryid.Bytes()),
			LoadGenerator: loadGenerator,
		}

		expendedData = append(expendedData, a)
	}

	return data, expendedData
}

func init() {
	resultCmd.Flags().BoolVarP(&viewSingleResult, "view", "", false, "(optional) View single performance results with more info")
	resultCmd.Flags().IntVarP(&pageNumber, "page", "p", 1, "(optional) List next set of performance results with --page (default = 1)")
}
