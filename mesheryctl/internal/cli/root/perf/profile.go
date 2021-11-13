package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/manifoldco/promptui"
	termbox "github.com/nsf/termbox-go"
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
	Short: "List performance profiles",
	Long:  `List all the available performance profiles`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// List performance profiles (maximum 25 profiles)	
mesheryctl perf profile

// List performance profiles with search (maximum 25 profiles)
mesheryctl perf profile test profile 2 -t "/Downloads/auth.json"

// View single performance profile with detailed information
mesheryctl perf profile --view -t "/Downloads/auth.json"
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// used for searching performance profile
		var searchString string
		// setting up for error formatting
		cmdUsed = "profile"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrMesheryConfig(err)
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
			return nil
		}

		if outputFormatFlag != "" {
			var tempStruct *models.PerformanceProfilesAPIResponse
			err := json.Unmarshal(body, &tempStruct)
			if err != nil {
				return ErrFailUnmarshal(err)
			}
			body, _ = json.Marshal(tempStruct.Profiles)
			if outputFormatFlag == "yaml" {
				body, _ = yaml.JSONToYAML(body)
			} else if outputFormatFlag != "json" {
				return ErrInvalidOutputChoice()
			}
			log.Info(string(body))
		} else if !expand {
			utils.PrintToTable([]string{"Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data)
		} else {
			// if data consists only one profile, directly print profile
			index := 0
			if len(data) > 1 {
				index, err = userPrompt("profile", "Enter index of the profile", data)
				if err != nil {
					return err
				}
			}

			a := expandedData[index]

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
		}

		return nil
	},
}

// Fetch all the profiles
func fetchPerformanceProfiles(url, searchString string) ([][]string, []profileStruct, []byte, error) {
	client := &http.Client{}
	var response *models.PerformanceProfilesAPIResponse

	// update the url
	tempURL := fmt.Sprintf("%s?page_size=%d&page=%d", url, pageSize, resultPage-1)
	if searchString != "" {
		tempURL = tempURL + "&search=" + searchString
	}

	req, _ := http.NewRequest("GET", tempURL, nil)

	err := utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return nil, nil, nil, ErrAttachAuthToken(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, nil, ErrFailRequest(err)
	}

	// failsafe for not being authenticated
	if utils.ContentTypeIsHTML(resp) {
		return nil, nil, nil, ErrUnauthenticated()
	}
	// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
	if resp.StatusCode != 200 {
		return nil, nil, nil, ErrFailReqStatus(resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, nil, nil, ErrFailUnmarshal(err)
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
	profileCmd.Flags().BoolVarP(&expand, "view", "", false, "(optional) View single performance profile with more info")
	profileCmd.Flags().IntVarP(&resultPage, "page", "p", 1, "(optional) List next set of performance results with --page (default = 1)")
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
		utils.PrintToTable([]string{"Index", "Name", "Mesh", "QPS", "Duration", "P50", "P99.9", "Start-Time"}, data)
	} else {
		utils.PrintToTable([]string{"Index", "Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data)
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
		return -1, err
	}
	return index, nil
}
