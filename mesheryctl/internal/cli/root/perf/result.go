package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var resultCmd = &cobra.Command{
	Use:   "profile profile-id [result-name]",
	Short: "List Performance profiles",
	Long:  `List all the available performance profiles`,
	Args:  cobra.MinimumNArgs(1),
	Example: `
// List performance profiles (maximum 25 profiles)	
mesheryctl perf profile 
// List performance profiles with search (maximum 25 profiles)
mesheryctl perf profile test profile 2 
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		resultURL := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles/" + args[0] + "/results"

		// set default tokenpath for command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		if len(args) > 0 {
			// Merge args to get result-name
			searchString = strings.Join(args, "%20")
		}

		data, body, err := fetchPerformanceProfileResults(resultURL, args[0], searchString)
		if err != nil {
			return err
		}

		if outputFormatFlag != "" {
			if outputFormatFlag == "yaml" {
				if body, err = yaml.JSONToYAML(body); err != nil {
					return errors.Wrap(err, "failed to convert json to yaml")
				}
			} else if outputFormatFlag != "json" {
				return errors.New("output-format choice invalid, use [json|yaml]")
			}
			log.Info(string(body))
			return nil
		}

		if len(data) > 0 {
			utils.PrintToTable([]string{"NAME", "MESH", "START-TIME", "QPS", "DURATION", "P50", "P99.9"}, data)
		} else {
			log.Info("No Performance Profiles to display")
		}
		return nil
	},
}

// Fetch results for a specific profile
func fetchPerformanceProfileResults(url, profileID, searchString string) ([][]string, []byte, error) {
	client := &http.Client{}
	var response *models.PerformanceResultsAPIResponse
	tempURL := fmt.Sprintf("%s?pageSize=%d", url, pageSize)
	if searchString != "" {
		tempURL = tempURL + "&search=" + searchString
	}
	req, _ := http.NewRequest("GET", tempURL, nil)

	err := utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return nil, nil, errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, errors.Wrapf(err, utils.PerfError("Failed to fetch performance results"))
	}
	// failsafe for no authentication
	if utils.ContentTypeIsHTML(resp) {
		return nil, nil, errors.New("invalid authentication token")
	}
	// failsafe for bad api call
	if resp.StatusCode != 200 {
		return nil, nil, errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", profileID)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	var data [][]string

	// append data for single profile
	for _, result := range response.Results {
		serviceMesh := "No Mesh"
		if result.Mesh != "" {
			serviceMesh = result.Mesh
		}
		url := result.RunnerResults.URL
		qps := fmt.Sprintf("%.8f", result.RunnerResults.QPS)
		duration := fmt.Sprintf("%d", uint64(result.RunnerResults.Duration))
		p50 := fmt.Sprintf("%.8f", result.RunnerResults.DurationHistogram.Percentiles[0].Value)
		p99_9 := fmt.Sprintf("%.8f", result.RunnerResults.DurationHistogram.Percentiles[len(result.RunnerResults.DurationHistogram.Percentiles)-1].Value)
		startTime := result.TestStartTime.Format("2006-01-02 15:04:05")
		data = append(data, []string{result.Name, serviceMesh, url, qps, duration, p50, p99_9, startTime})
	}

	return data, body, nil
}
