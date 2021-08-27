package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	searchString string
	pageSize     = 25
)

var profileCmd = &cobra.Command{
	Use:   "profile [profile-name]",
	Short: "List Performance profiles",
	Long:  `List all the available performance profiles`,
	Args:  cobra.MinimumNArgs(0),
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

		profileURL := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles"

		// set default tokenpath for command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		if len(args) > 0 {
			// Merge args to get profile-name
			searchString = strings.Join(args, "%20")
		}

		data, body, err := fetchPerformanceProfiles(profileURL, searchString)
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
			utils.PrintToTable([]string{"Name", "ID", "RESULTS", "Load-Generator", "Endpoint", "Duration", "Last-Run"}, data)
		} else {
			log.Info("No Performance Profiles to display")
		}

		return nil
	},
}

// Fetch all the profiles
func fetchPerformanceProfiles(url, searchString string) ([][]string, []byte, error) {
	client := &http.Client{}
	var response *models.PerformanceProfilesAPIResponse

	// update the url
	tempURL := fmt.Sprintf("%s?page_size=%d", url, pageSize)
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
		return nil, nil, errors.New("failed to make a request")
	}

	// failsafe for not being authenticated
	if utils.ContentTypeIsHTML(resp) {
		return nil, nil, errors.New("invalid authentication token")
	}
	// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
	if resp.StatusCode != 200 {
		return nil, nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
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

	for _, profile := range response.Profiles {
		if profile.LastRun != nil {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], profile.Endpoints[0], profile.Duration, profile.LastRun.Time.Format("2006-01-02 15:04:05")})
		} else {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], profile.Endpoints[0], profile.Duration, ""})
		}
	}

	return data, body, nil
}
