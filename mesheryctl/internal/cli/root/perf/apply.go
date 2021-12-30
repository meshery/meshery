package perf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	SMP "github.com/layer5io/service-mesh-performance/spec"
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
	profileID          string
	req                *http.Request
)

var applyCmd = &cobra.Command{
	Use:   "apply [profile-name | --file] --flags",
	Short: "Run a Performance test",
	Long:  `Run Performance test using existing profiles or using flags`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Execute a Performance test with the specified performance profile
mesheryctl perf apply meshery-profile --flags

// Execute a Performance test with creating a new performance profile
mesheryctl perf apply meshery-profile-new --url "https://google.com"

// Run Performance test using SMP compatible test configuration
mesheryctl perf apply -f perf-config.yaml

// Run performance test using SMP compatible test configuration and override values with flags
mesheryctl perf apply -f <filepath> --flags

// Choice of load generator - fortio or wrk2 (default: fortio)
mesheryctl perf apply meshery-test --load-generator wrk2

// Execute a Performance test with specified queries per second
mesheryctl perf apply local-perf --url https://192.168.1.15/productpage --qps 30

// Execute a Performance test with specified service mesh
mesheryctl perf apply local-perf --url https://192.168.1.15/productpage --mesh istio
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{}
		userResponse := false

		// setting up for error formatting
		cmdUsed = "apply"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrMesheryConfig(err)
		}

		// set default tokenpath for command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		// Importing SMP Configuration from the file
		// TODO: Refactor: Move checks to a single location and consolidate for file, flags and performance profile
		if filePath != "" {
			// Read the test configuration file
			smpConfig, err := os.ReadFile(filePath)
			if err != nil {
				return ErrReadFilepath(err)
			}

			testConfig := models.PerformanceTestConfigFile{}

			err = yaml.Unmarshal(smpConfig, &testConfig)
			if err != nil {
				return ErrFailUnmarshalFile(err)
			}

			if testConfig.Config == nil || testConfig.ServiceMesh == nil {
				return ErrInvalidTestConfigFile()
			}

			testClient := testConfig.Config.Clients[0]

			// Override values from configuration if passed on as flags
			if testName == "" {
				testName = testConfig.Config.Name
			}

			if testURL == "" {
				testURL = testClient.EndpointUrls[0]
			}

			if testMesh == "" {
				testMesh = SMP.ServiceMesh_Type_name[int32(testConfig.ServiceMesh.Type)]
			}

			if qps == "" {
				qps = strconv.FormatInt(testClient.Rps, 10)
			}

			if concurrentRequests == "" {
				concurrentRequests = strconv.Itoa(int(testClient.Connections))
			}

			if testDuration == "" {
				testDuration = testConfig.Config.Duration
			}

			if loadGenerator == "" {
				loadGenerator = testClient.LoadGenerator
			}
		}

		// Run test based on flags
		if testName == "" {
			log.Debug("Test Name not provided")
			testName = utils.StringWithCharset(8)
			log.Debug("Using random test name: ", testName)
		}

		// Throw error if a profile name is not provided
		if len(args) == 0 {
			return ErrNoProfileName()
		}
		profileName = strings.Join(args, "%20")

		// Check if the profile name is valid, if not prompt the user to create a new one
		log.Debug("Fetching performance profile")

		req, _ = http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles?search="+profileName, nil)

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return ErrAttachAuthToken(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return ErrFailRequest(err)
		}

		var response *models.PerformanceProfilesAPIResponse
		// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
		if resp.StatusCode != 200 {
			return ErrFailReqStatus(resp.StatusCode)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return ErrFailUnmarshal(err)
		}

		index := 0
		if len(response.Profiles) == 0 {
			// if the provided performance profile does not exist, prompt the user to create a new one

			// skip asking confirmation if -y flag used
			if utils.SilentFlag {
				userResponse = true
			} else {
				// ask user for confirmation
				userResponse = utils.AskForConfirmation("Profile with name '" + profileName + "' does not exist. Do you want to create a new one")
			}

			if userResponse {
				profileID, profileName, err = createPerformanceProfile(client, mctlCfg)
				if err != nil {
					return err
				}
			} else {
				return ErrNoProfileFound()
			}
		} else {
			if len(response.Profiles) == 1 {
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
			// pull test configuration from the profile only if a test configuration is not provided
			if filePath == "" {
				profileName = response.Profiles[index].Name
				loadGenerator = response.Profiles[index].LoadGenerators[0]
				concurrentRequests = strconv.Itoa(response.Profiles[index].ConcurrentRequest)
				qps = strconv.Itoa(response.Profiles[index].QPS)
				testDuration = response.Profiles[index].Duration
				testMesh = response.Profiles[index].ServiceMesh
			}
		}

		if testURL == "" {
			return ErrNoTestURL()
		}

		log.Debugf("performance profile is: %s", profileName)
		log.Debugf("test-url set to %s", testURL)

		// Method to check if the entered Test URL is valid or not
		if validURL := govalidator.IsURL(testURL); !validURL {
			return ErrNotValidURL()
		}

		req, _ = http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/run", nil)

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

		log.Info("Initiating Performance test ...")

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return ErrAttachAuthToken(err)
		}

		resp, err = client.Do(req)
		if err != nil {
			return ErrFailRequest(err)
		}
		if utils.ContentTypeIsHTML(resp) {
			return ErrFailTestRun()
		}
		if resp.StatusCode != 200 {
			return ErrFailTestRun()
		}

		defer utils.SafeClose(resp.Body)
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		log.Debug(string(data))

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
		fmt.Printf("Enter the index of profile: ")
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
	applyCmd.Flags().StringVar(&testURL, "url", "", "(optional) Endpoint URL to test (required with --profile)")
	applyCmd.Flags().StringVar(&testName, "name", "", "(optional) Name of the Test")
	applyCmd.Flags().StringVar(&testMesh, "mesh", "", "(optional) Name of the Service Mesh")
	applyCmd.Flags().StringVar(&qps, "qps", "", "(optional) Queries per second")
	applyCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "", "(optional) Number of Parallel Requests")
	applyCmd.Flags().StringVar(&testDuration, "duration", "", "(optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration")
	applyCmd.Flags().StringVar(&loadGenerator, "load-generator", "", "(optional) Load-Generator to be used (fortio/wrk2)")
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "(optional) file containing SMP-compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification")
}

func createPerformanceProfile(client *http.Client, mctlCfg *config.MesheryCtlConfig) (string, string, error) {
	log.Debug("Creating new performance profile inside function")

	if profileName == "" {
		return "", "", ErrNoProfileName()
	}

	// ask for test url first
	if testURL == "" {
		return "", "", ErrNoTestURL()
	}

	// Method to check if the entered Test URL is valid or not
	if validURL := govalidator.IsURL(testURL); !validURL {
		return "", "", ErrNotValidURL()
	}

	if testMesh == "" {
		testMesh = "None"
	}

	if qps == "" {
		qps = "0"
	}

	if concurrentRequests == "" {
		concurrentRequests = "1"
	}

	if testDuration == "" {
		testDuration = "30s"
	}

	if loadGenerator == "" {
		loadGenerator = "fortio"
	}

	convReq, err := strconv.Atoi(concurrentRequests)
	if err != nil {
		return "", "", errors.New("failed to convert concurrent-request")
	}
	convQPS, err := strconv.Atoi(qps)
	if err != nil {
		return "", "", errors.New("failed to convert qps")
	}
	values := map[string]interface{}{
		"concurrent_request": convReq,
		"duration":           testDuration,
		"endpoints":          []string{testURL},
		"load_generators":    []string{loadGenerator},
		"name":               profileName,
		"qps":                convQPS,
		"service_mesh":       testMesh,
		"request_body":       "",
		"request_cookies":    "",
		"request_headers":    "",
		"content_type":       "",
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", "", ErrFailMarshal(err)
	}
	req, _ = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles", bytes.NewBuffer(jsonValue))

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", "", ErrAttachAuthToken(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", "", ErrFailRequest(err)
	}

	var response *models.PerformanceProfile
	// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
	if resp.StatusCode != 200 {
		return "", "", ErrFailReqStatus(resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", "", ErrFailUnmarshal(err)
	}
	profileID = response.ID.String()
	profileName = response.Name

	log.Debug("New profile created")
	return profileID, profileName, nil
}
