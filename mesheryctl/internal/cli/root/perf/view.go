package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type profileStruct struct {
	Name           string
	ID             *uuid.UUID
	Endpoints      []string
	QPS            int
	Duration       string
	Loadgenerators []string
}

type resultStruct struct {
	Name          string
	StartTime     *time.Time
	LatenciesMs   *models.LatenciesMs
	QPS           int
	URL           string
	UserID        *uuid.UUID
	Duration      int
	MesheryID     *uuid.UUID
	LoadGenerator string
}

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view perf profile",
	Long:  `See the configuration of your performance profile`,
	Example: `
	View performance profile with
	mesheryctl perf view <profile-name>

	View performance results with
	mesheryctl perf view <profile-id> <test-name>
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		if len(args) == 2 {
			return viewResults(args, mctlCfg)
		}
		proName := args[0]
		var req *http.Request
		url := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles?page_size=25&search=" + proName
		var response *models.PerformanceProfilePage
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
		}
		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", proName)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}

		if err = json.Unmarshal(data, &response); err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		var a profileStruct
		if response.TotalCount == 0 {
			return errors.New("profile does not exit. Please get a profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles")
		}
		for _, profile := range response.Profiles {
			if response.Profiles == nil {
				return errors.New("profile name not provide. Please get a profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles")
			}
			a = profileStruct{
				Name:           profile.Name,
				ID:             profile.ID,
				Endpoints:      profile.Endpoints,
				QPS:            profile.QPS,
				Duration:       profile.Duration,
				Loadgenerators: profile.LoadGenerators,
			}
			data, err = json.MarshalIndent(&a, "", "  ")
			if err != nil {
				return err
			}
			if outputFormatFlag == "yaml" {
				data, err = yaml.JSONToYAML(data)
				if err != nil {
					return errors.Wrap(err, "failed to convert json to yaml")
				}
			} else if outputFormatFlag == "" {
				fmt.Printf("Name: %v\n", a.Name)
				fmt.Printf("ID: %s\n", a.ID.String())
				fmt.Printf("Endpoint: %v\n", a.Endpoints[0])
				fmt.Printf("Load Generators: %v\n", a.Loadgenerators[0])
				fmt.Printf("Test run duration: %v\n", a.Duration)
				fmt.Println("#####################")
				continue
			} else if outputFormatFlag != "json" {
				return errors.New("output-format not supported, use [json|yaml]")
			}
			log.Info(string(data))

		}
		return nil
	},
}

func viewResults(args []string, mctlCfg *config.MesheryCtlConfig) error {
	profileID := args[0]
	testName := args[1]
	var req *http.Request
	var response *models.PerformanceResultsAPIResponse

	url := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles/" + profileID + "/results?pageSize=25&search=" + testName
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
	}
	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.Errorf("Profile id `%s` is invalid. Please verify profile-id and try again. Use `mesheryctl perf list` to see a list of performance profiles.", profileID)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, utils.PerfError("failed to read response body"))
	}

	if err = json.Unmarshal(data, &response); err != nil {
		return errors.Wrap(err, "failed to unmarshal response body")
	}

	var a resultStruct
	if response.TotalCount == 0 {
		return errors.New("results does not exit. Please run a profile test and try again. Use `mesheryctl perf list` to see a list of performance profiles")
	}
	for _, result := range response.Results {
		if response.Results == nil {
			return errors.New("profile name not provide. Please get a profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles")
		}
		a = resultStruct{
			Name:     result.Name,
			UserID:   result.UserID,
			URL:      result.RunnerResults.URL,
			QPS:      int(result.RunnerResults.QPS),
			Duration: int(result.RunnerResults.Duration),
			LatenciesMs: &models.LatenciesMs{
				Average: result.RunnerResults.DurationHistogram.Average,
				Max:     result.RunnerResults.DurationHistogram.Max,
				Min:     result.RunnerResults.DurationHistogram.Min,
				P50:     result.RunnerResults.DurationHistogram.Percentiles[0].Value,
				P90:     result.RunnerResults.DurationHistogram.Percentiles[2].Value,
				P99:     result.RunnerResults.DurationHistogram.Percentiles[3].Value,
			},
			StartTime:     result.TestStartTime,
			MesheryID:     result.MesheryID,
			LoadGenerator: result.RunnerResults.LoadGenerator,
		}
		data, err = json.MarshalIndent(&a, "", "  ")
		if err != nil {
			return err
		}
		if outputFormatFlag == "yaml" {
			data, err = yaml.JSONToYAML(data)
			if err != nil {
				return errors.Wrap(err, "failed to convert json to yaml")
			}
		} else if outputFormatFlag == "" {
			fmt.Printf("Name: %v\n", a.Name)
			fmt.Printf("UserID: %s\n", a.UserID.String())
			fmt.Printf("Endpoint: %v\n", a.URL)
			fmt.Printf("QPS: %v\n", a.QPS)
			fmt.Printf("Test run duration: %v\n", a.Duration)
			fmt.Printf("Latencies _ms: Avg: %v, Max: %v, Min: %v, P50: %v, P90: %v, P99: %v\n", a.LatenciesMs.Average, a.LatenciesMs.Max, a.LatenciesMs.Min, a.LatenciesMs.P50, a.LatenciesMs.P90, a.LatenciesMs.P99)
			fmt.Printf("Start Time: %v\n", fmt.Sprintf("%d-%d-%d %d:%d:%d", int(a.StartTime.Month()), a.StartTime.Day(), a.StartTime.Year(), a.StartTime.Hour(), a.StartTime.Minute(), a.StartTime.Second()))
			fmt.Printf("Meshery ID: %v\n", a.MesheryID.String())
			fmt.Printf("Load Generator: %v\n", a.LoadGenerator)
			fmt.Println("#####################")
			continue
		} else if outputFormatFlag != "json" {
			return errors.New("output-format not supported, use [json|yaml]")
		}
		log.Info(string(data))
	}
	return nil
}

func init() {
	_ = viewCmd.MarkFlagRequired("token")
}
