package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
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
	Use:   "list",
	Short: "List Performance profiles",
	Long:  `List all the available performance profiles`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		client := &http.Client{}
		if len(args) == 0 {
			req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/results", nil)
			if err != nil {
				return errors.Wrapf(err, utils.PerfError("failed to create a request"))
			}
			err = utils.AddAuthDetails(req, tokenPath)
			if err != nil {
				return err
			}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, utils.PerfError("failed to read response body"))
			}
			var response models.PerformanceAPIResponse
			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal response body")
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

		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/results", nil)
		if err != nil {
			return errors.Wrapf(err, utils.PerfError("Failed to fetch performance results"))
		}
		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		var response models.PerformanceAPIResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}
		var data [][]string
		for _, result := range response.Results {
			serviceMesh := "No Mesh"
			if result.Mesh != "" {
				serviceMesh = result.Mesh
			}
			startTime := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(result.TestStartTime.Month()), result.TestStartTime.Day(), result.TestStartTime.Year(), result.TestStartTime.Hour(), result.TestStartTime.Minute(), result.TestStartTime.Second())
			P50 := result.RunnerResults.DurationHistogram.Percentiles[0].Value
			P99_9 := result.RunnerResults.DurationHistogram.Percentiles[len(result.RunnerResults.DurationHistogram.Percentiles)-1].Value
			data = append(data, []string{result.Name, serviceMesh, startTime, fmt.Sprintf("%f", result.RunnerResults.Qps), result.RunnerResults.Duration, fmt.Sprintf("%f", P50), fmt.Sprintf("%f", P99_9)})
		}
		utils.PrintToTable([]string{"NAME", "MESH", "START-TIME", "QPS", "DURATION", "P50", "P99.9"}, data)
		return nil
	},
}

func init() {
	listCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", utils.AuthConfigFile, "(optional) Path to meshery auth config")
}
