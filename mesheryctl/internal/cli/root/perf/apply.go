package perf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	profileName        string
	testURL            string
	testName           string
	testMesh           string
	qps                string
	concurrentRequests string
	testDuration       string
	loadGenerator      string
	filePath           string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Run a Performance test",
	Long:  `Run Performance test using existing profiles or using flags`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
	Execute a Performance test with the specified performance profile
	mesheryctl perf apply <profile id> --flags

	Execute a Performance test without a specified performance profile
	mesheryctl perf apply --profile <profile-name> --url <url>

	Run Performance test using SMP compatible test configuration
	mesheryctl perf apply -f <filepath>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		client := &http.Client{}
		var profileID string

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		// Importing SMP Configuration from the file
		if filePath != "" {
			smpConfig, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}

			req, err = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/perf/load-test-smps", bytes.NewBuffer(smpConfig))
			if err != nil {
				return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
			}
		} else {
			// Run test based on flags
			if testName == "" {
				log.Debug("Test Name not provided")
				testName = utils.StringWithCharset(8)
				log.Debug("Using random test name: ", testName)
			}

			// If --profile passed than create a profile
			if len(args) == 0 { // First need to create a profile id
				log.Debug("Creating new performance profile")

				if profileName == "" {
					return errors.New(utils.PerfError("please enter a profile-name"))
				}

				// ask for test url first
				if testURL == "" {
					return errors.New(utils.PerfError("please enter a test URL"))
				}

				// Method to check if the entered Test URL is valid or not
				if validURL := govalidator.IsURL(testURL); !validURL {
					return errors.New(utils.PerfError("please enter a valid test URL"))
				}

				convReq, err := strconv.Atoi(concurrentRequests)
				if err != nil {
					return errors.New("failed to convert concurrent-request")
				}
				convQPS, err := strconv.Atoi(qps)
				if err != nil {
					return errors.New("failed to convert qps")
				}
				values := map[string]interface{}{
					"concurrent_request": convReq,
					"duration":           testDuration,
					"endpoints":          []string{testURL},
					"load_generators":    []string{loadGenerator},
					"name":               profileName,
					"qps":                convQPS,
					"service_mesh":       testMesh,
				}

				jsonValue, err := json.Marshal(values)
				if err != nil {
					return err
				}
				req, err = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles", bytes.NewBuffer(jsonValue))
				if err != nil {
					return err
				}

				err = utils.AddAuthDetails(req, tokenPath)
				if err != nil {
					return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
				}

				resp, err := client.Do(req)
				if err != nil {
					return err
				}

				var response *models.PerformanceProfile
				// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
				if resp.StatusCode != 200 {
					return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return errors.Wrap(err, utils.PerfError("failed to read response body"))
				}
				err = json.Unmarshal(body, &response)
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal response body")
				}
				profileID = response.ID.String()
				profileName = response.Name

				log.Debug("New profile created")
			} else { // set profile-name from args
				// Merge args to get profile-name
				profileName = strings.Join(args, "%20")

				// search and fetch performance profile with profile-name
				log.Debug("Fetching performance profile")

				req, err = http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles?search="+profileName, nil)
				if err != nil {
					return err
				}

				err = utils.AddAuthDetails(req, tokenPath)
				if err != nil {
					return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
				}

				resp, err := client.Do(req)
				if err != nil {
					return err
				}

				var response *models.PerformanceProfilesAPIResponse
				// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
				if resp.StatusCode != 200 {
					return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
				}
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return errors.Wrap(err, utils.PerfError("failed to read response body"))
				}
				err = json.Unmarshal(body, &response)
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal response body")
				}

				index := 0
				if len(response.Profiles) == 0 {
					return errors.New("no profiles found with the given name")
				} else if len(response.Profiles) == 1 {
					profileID = response.Profiles[0].ID.String()
				} else {
					// Multiple profiles with same name
					index = multipleProfileConfirmation(response.Profiles)
					profileID = response.Profiles[index].ID.String()
				}
				// what if user passed profile-name but didn't passed the url
				// we use url from performance profile
				if testURL == "" {
					testURL = response.Profiles[index].Endpoints[0]
				}

				// reset profile name without %20
				profileName = response.Profiles[index].Name
			}

			if testURL == "" {
				return errors.New(utils.PerfError("please enter a test URL"))
			}

			log.Debugf("performance profile is: %s", profileName)
			log.Debugf("test-url set to %s", testURL)

			// Method to check if the entered Test URL is valid or not
			if validURL := govalidator.IsURL(testURL); !validURL {
				return errors.New(utils.PerfError("please enter a valid test URL"))
			}

			req, err = http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/run", nil)
			if err != nil {
				return err
			}

			q := req.URL.Query()

			q.Add("name", testName)
			q.Add("loadGenerator", loadGenerator)
			q.Add("c", concurrentRequests)
			q.Add("url", testURL)
			q.Add("qps", qps)

			durLen := len(testDuration)

			q.Add("dur", string(testDuration[durLen-1]))
			q.Add("t", string(testDuration[:durLen-1]))

			if testMesh != "" {
				q.Add("mesh", testMesh)
			}
			req.URL.RawQuery = q.Encode()
		}

		log.Info("Initiating Performance test ...")

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
		}

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, utils.PerfError(fmt.Sprintf("failed to make request to %s", testURL)))
		}
		if utils.ContentTypeIsHTML(resp) {
			return errors.New("failed to run test")
		}
		if resp.StatusCode != 200 {
			return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
		}

		defer utils.SafeClose(resp.Body)
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		log.Debug(string(data))

		if err := utils.UpdateAuthDetails(tokenPath); err != nil {
			return errors.Wrap(err, utils.PerfError("failed to update auth details"))
		}

		log.Info("Test Completed Successfully!")
		return nil
	},
}

func multipleProfileConfirmation(profiles []models.PerformanceProfile) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("Endpoint: %v\n", a.Endpoints[0])
		fmt.Printf("Load Generators: %v\n", a.LoadGenerators[0])
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of profile?? ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		index, err := strconv.Atoi(response)
		if err != nil {
			log.Info(err)
		}
		if index < 0 || index >= len(profiles) {
			log.Info("Invalid index")
		} else {
			return index
		}
	}
}

func init() {
	applyCmd.Flags().StringVar(&testURL, "url", "", "(required/optional) Endpoint URL to test")
	applyCmd.Flags().StringVar(&testName, "name", "", "(optional) Name of the Test")
	applyCmd.Flags().StringVar(&profileName, "profile", "", "(required/optional) Name for the new Performance Profile")
	applyCmd.Flags().StringVar(&testMesh, "mesh", "None", "(optional) Name of the Service Mesh")
	applyCmd.Flags().StringVar(&qps, "qps", "0", "(optional) Queries per second")
	applyCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "1", "(optional) Number of Parallel Requests")
	applyCmd.Flags().StringVar(&testDuration, "duration", "30s", "(optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/-ParseDuration")
	applyCmd.Flags().StringVar(&loadGenerator, "load-generator", "fortio", "(optional) Load-Generator to be used (fortio/wrk2)")
	applyCmd.Flags().StringVar(&filePath, "file", "", "(optional) file containing SMP-compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification")

	_ = listCmd.MarkFlagRequired("token")
}
