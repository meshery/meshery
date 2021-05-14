package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ProfileStruct to store profile related data
type ProfileStruct struct {
	ID      string
	LastRun *time.Time
	Results uint
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Performance profiles",
	Long:    `List all the available performance profiles`,
	Args:    cobra.MaximumNArgs(1),
	Example: "mesheryctl perf list \nmesheryctl perf list [profile-id]",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		if len(args) == 0 {
			response, err := fetchPerformanceAPIResponse(mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/results", 0, "")
			if err != nil {
				return err
			}
			if response == nil {
				return errors.New("failed to fetch results")
			}
			if response.TotalCount > 10 {
				response, err = fetchPerformanceAPIResponse(mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/results", response.TotalCount, "")
				if err != nil {
					return err
				}
			}
			// map to store profile related data wrt profile-id
			dataMap := make(map[string]*ProfileStruct)
			for _, result := range response.Results {
				// id assigned to the tests performed before introducing profiles
				id := ""
				if result.PerformanceProfile != nil {
					id = result.PerformanceProfile.String()
				}
				if _, present := dataMap[id]; !present {
					//Add new profile struct for new profile-id
					dataMap[id] = &ProfileStruct{
						ID:      id,
						LastRun: result.TestStartTime,
						Results: 1,
					}
				} else {
					// Increase the result count for the performance profile
					dataMap[id].Results = dataMap[id].Results + 1
					// Update time if the current result's LastRun is after the stored-one
					if dataMap[id].LastRun.After(*result.TestStartTime) {
						dataMap[id].LastRun = result.TestStartTime
					}
				}

			}

			var data [][]string
			for _, profile := range dataMap {
				id := profile.ID
				results := profile.Results
				lastRun := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(profile.LastRun.Month()), profile.LastRun.Day(), profile.LastRun.Year(), profile.LastRun.Hour(), profile.LastRun.Minute(), profile.LastRun.Second())
				data = append(data, []string{id, strconv.FormatUint(uint64(results), 10), lastRun})
			}
			utils.PrintToTable([]string{"ID", "RESULTS", "LAST-RUN"}, data)

			return nil
		}
		// Output results of a performance profile
		profileID := args[0]
		response, err := fetchPerformanceAPIResponse(mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/results", 0, profileID)
		if err != nil {
			return err
		}
		if response == nil {
			return errors.New("failed to fetch results")
		}
		if response.TotalCount > 10 {
			response, err = fetchPerformanceAPIResponse(mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/results", response.TotalCount, profileID)
			if err != nil {
				return err
			}
		}

		var data [][]string
		for _, result := range response.Results {
			serviceMesh := "No Mesh"
			if result.Mesh != "" {
				serviceMesh = result.Mesh
			}
			startTime := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(result.TestStartTime.Month()), result.TestStartTime.Day(), result.TestStartTime.Year(), result.TestStartTime.Hour(), result.TestStartTime.Minute(), result.TestStartTime.Second())
			p50 := result.RunnerResults.DurationHistogram.Percentiles[0].Value
			p99_9 := result.RunnerResults.DurationHistogram.Percentiles[len(result.RunnerResults.DurationHistogram.Percentiles)-1].Value
			timeDuration := strings.SplitAfterN(strconv.FormatUint(uint64(result.RunnerResults.Duration), 10), "", 5)
			duration := fmt.Sprintf("%s%s.%s%ss", timeDuration[0], timeDuration[1], timeDuration[2], timeDuration[3])
			data = append(data, []string{result.Name, serviceMesh, startTime, fmt.Sprintf("%f", result.RunnerResults.QPS), duration, fmt.Sprintf("%f", p50), fmt.Sprintf("%f", p99_9)})
		}
		utils.PrintToTable([]string{"NAME", "MESH", "START-TIME", "QPS", "DURATION", "P50", "P99.9"}, data)
		return nil
	},
}

func fetchPerformanceAPIResponse(url string, resultCount uint, profileID string) (*models.PerformanceAPIResponse, error) {
	client := &http.Client{}
	var response *models.PerformanceAPIResponse

	for i := 0; i < int(resultCount/25)+1; i++ {
		tempURL := fmt.Sprintf("%s?pageSize=25&page=%d", url, i)
		req, err := http.NewRequest("GET", tempURL, nil)
		if err != nil {
			return nil, errors.Wrapf(err, utils.PerfError("Failed to fetch performance results"))
		}
		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return nil, err
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != 200 {
			// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
			if profileID != "" {
				return nil, errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", profileID)
			}
			return nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		// for the first request unmarshal to "response"
		if i == 0 {
			err = json.Unmarshal(body, &response)
			if err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal response body")
			}
			continue
		}
		var temp *models.PerformanceAPIResponse
		err = json.Unmarshal(body, &temp)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal response body")
		}
		// append the results if no results then break
		if temp != nil && len(temp.Results) > 0 {
			response.Results = append(response.Results, temp.Results...)
		} else {
			break
		}
	}

	return response, nil
}

func init() {
	listCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", utils.AuthConfigFile, "(optional) Path to meshery auth config")
}
