package perf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
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
	resultPage int
)

var resultCmd = &cobra.Command{
	Use:   "result profile-name",
	Short: "List performance test results",
	Long:  `List all the available test results of a performance profile`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// List Test results (maximum 25 results)	
mesheryctl perf result saturday-profile 

// View other set of performance results with --page (maximum 25 results)
mesheryctl perf result saturday-profile --page 2 

// View single performance result with detailed information
mesheryctl perf result saturday-profile --view 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// used for searching performance profile
		var searchString, profileID string
		// setting up for error formatting
		cmdUsed = "result"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrMesheryConfig(err)
		}

		profileURL := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles"

		// Merge args to get result-name
		searchString = strings.Join(args, "%20")

		data, _, _, err := fetchPerformanceProfiles(profileURL, searchString)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			log.Info("No Performance Profiles found with the given name")
			return nil
		}

		if len(data) == 1 {
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

		resultURL := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles/" + profileID + "/results"

		data, expandedData, body, err := fetchPerformanceProfileResults(resultURL)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			log.Info("No Test Results to display")
			return nil
		}

		if outputFormatFlag != "" {
			var tempStruct *models.PerformanceResultsAPIResponse
			err = json.Unmarshal(body, &tempStruct)
			if err != nil {
				return ErrFailUnmarshal(err)
			}
			body, _ = json.Marshal(tempStruct.Results)
			if outputFormatFlag == "yaml" {
				body, _ = yaml.JSONToYAML(body)
			} else if outputFormatFlag != "json" {
				return ErrInvalidOutputChoice()
			}
			log.Info(string(body))
		} else if !expand {
			utils.PrintToTable([]string{"NAME", "MESH", "QPS", "DURATION", "P50", "P99.9", "START-TIME"}, data)
		} else {
			// if data consists only one profile, directly print profile
			index := 0
			if len(data) > 1 {
				index, err = userPrompt("result", "Select Performance-test result to exapand", data)
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
func fetchPerformanceProfileResults(url string) ([][]string, []resultStruct, []byte, error) {
	client := &http.Client{}
	var response *models.PerformanceResultsAPIResponse
	tempURL := fmt.Sprintf("%s?pageSize=%d&page=%d", url, pageSize, resultPage-1)

	req, err := utils.NewRequest("GET", tempURL, nil)
	if err != nil {
		return nil, nil, nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, ErrFailRequest(err)
	}
	// failsafe for no authentication
	if utils.ContentTypeIsHTML(resp) {
		return nil, nil, nil, ErrUnauthenticated()
	}
	// failsafe for bad api call
	if resp.StatusCode != 200 {
		return nil, nil, nil, ErrFailReqStatus(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, nil, ErrFailUnmarshal(err)
	}

	var data [][]string
	var expendedData []resultStruct

	// append data for single profile
	for _, result := range response.Results {
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
		a := resultStruct{
			Name:     result.Name,
			UserID:   result.UserID,
			URL:      result.RunnerResults.URL,
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
			MesheryID:     result.MesheryID,
			LoadGenerator: result.RunnerResults.LoadGenerator,
		}

		expendedData = append(expendedData, a)
	}

	return data, expendedData, body, nil
}

func init() {
	resultCmd.Flags().BoolVarP(&expand, "view", "", false, "(optional) View single performance results with more info")
	resultCmd.Flags().IntVarP(&resultPage, "page", "p", 1, "(optional) List next set of performance results with --page (default = 1)")
}
