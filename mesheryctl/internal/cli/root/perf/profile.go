package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/layer5io/meshery/internal/sql"

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
	pageSize = 25
	expand   bool
)

type profileStruct struct {
	Name           string
	ID             *uuid.UUID
	TotalResults   int
	Endpoints      []string
	QPS            int
	Duration       string
	Loadgenerators []string
	ServiceMesh    string
	LastRun        *sql.Time
}

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
		// used for searching performance profile
		var searchString string

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

		data, expandedData, body, err := fetchPerformanceProfiles(profileURL, searchString)
		if err != nil {
			return err
		}

		if len(data) == 0 {
			log.Info("No Performance Profiles to display")
		}

		if outputFormatFlag != "" {
			if outputFormatFlag == "yaml" {
				body, _ = yaml.JSONToYAML(body)
			} else if outputFormatFlag != "json" {
				return errors.New("output-format choice invalid, use [json|yaml]")
			}
			log.Info(string(body))
			return nil
		} else if !expand {
			utils.PrintToTable([]string{"Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data)
		} else {
			for _, a := range expandedData {
				fmt.Printf("Name: %v\n", a.Name)
				fmt.Printf("ID: %s\n", a.ID.String())
				fmt.Printf("Total Results: %d\n", a.TotalResults)
				fmt.Printf("Endpoint: %v\n", a.Endpoints[0])
				fmt.Printf("Load Generators: %v\n", a.Loadgenerators[0])
				fmt.Printf("Test run duration: %v\n", a.Duration)
				fmt.Printf("QPS: %d\n", a.QPS)
				fmt.Printf("Service Mesh: %v\n", a.ServiceMesh)
				if a.LastRun != nil {
					fmt.Printf("Last Run: %v\n", a.LastRun.Time.Format("2006-01-02 15:04:05"))
				} else {
					fmt.Printf("Last Run: %v\n", "nil")
				}
				fmt.Println("#####################")
			}
		}

		return nil
	},
}

// Fetch all the profiles
func fetchPerformanceProfiles(url, searchString string) ([][]string, []profileStruct, []byte, error) {
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
		return nil, nil, nil, errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, errors.New("failed to make a request")
	}

	// failsafe for not being authenticated
	if utils.ContentTypeIsHTML(resp) {
		return nil, nil, nil, errors.New("invalid authentication token")
	}
	// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
	if resp.StatusCode != 200 {
		return nil, nil, nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "failed to unmarshal response body")
	}
	var data [][]string
	var expendedData []profileStruct

	for _, profile := range response.Profiles {
		// adding stuff to data for list output
		if profile.LastRun != nil {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], profile.LastRun.Time.Format("2006-01-02 15:04:05")})
		} else {
			data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), profile.LoadGenerators[0], ""})
		}
		// adding stuff to expendedData for expended output
		a := profileStruct{
			Name:           profile.Name,
			ID:             profile.ID,
			TotalResults:   profile.TotalResults,
			Endpoints:      profile.Endpoints,
			QPS:            profile.QPS,
			Duration:       profile.Duration,
			Loadgenerators: profile.LoadGenerators,
			LastRun:        profile.LastRun,
			ServiceMesh:    profile.ServiceMesh,
		}
		expendedData = append(expendedData, a)
	}

	return data, expendedData, body, nil
}

func init() {
	profileCmd.Flags().BoolVarP(&expand, "expand", "e", false, "(optional) Expand the output")
}
