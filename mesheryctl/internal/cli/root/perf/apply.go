// Copyright 2023 Layer5, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package perf

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
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
	loadTestBody       string
	additionalOptions  string
	certPath           string
	disableCert        bool
)

var linkDocPerfApply = map[string]string{
	"link":    "![perf-apply-usage](/assets/img/mesheryctl/perf-apply.png)",
	"caption": "Usage of mesheryctl perf apply",
}

var applyCmd = &cobra.Command{
	Use:   "apply [profile-name]",
	Short: "Run a Performance test",
	Long:  `Run Performance test using existing profiles or using flags`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Execute a Performance test with the specified performance profile
mesheryctl perf apply meshery-profile [flags]

// Execute a Performance test with creating a new performance profile
mesheryctl perf apply meshery-profile-new --url "https://google.com"

// Execute a Performance test creating a new performance profile and pass certificate to be used 
mesheryctl perf apply meshery-profile-new --url "https://google.com" --cert-path path/to/cert.pem

// Execute a performance profile without using the certificate present in the profile
mesheryctl perf apply meshery-profile --url "https://google.com" --disable-cert

// Run Performance test using SMP compatible test configuration
// If the profile already exists, the test will be run overriding the values with the ones provided in the configuration file
mesheryctl perf apply meshery-profile -f path/to/perf-config.yaml

// Run performance test using SMP compatible test configuration and override values with flags
mesheryctl perf apply meshery-profile -f path/to/perf-config.yaml [flags]

// Choice of load generator - fortio, wrk2 or nighthawk (default: fortio)
mesheryctl perf apply meshery-profile --load-generator wrk2

// Execute a Performance test with specified queries per second
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --qps 30

// Execute a Performance test with specified service mesh
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --mesh istio

// Execute a Performance test creating a new performance profile and pass options to the load generator used
// If any options are already present in the profile or passed through flags, the --options flag will take precedence over the profile and flag options 
// Options for nighthawk - https://github.com/layer5io/getnighthawk/blob/v1.0.5/pkg/proto/options.pb.go#L882-L1018
// Options for fortio - https://github.com/fortio/fortio/blob/v1.57.0/fhttp/httprunner.go#L77-L84
// Options for wrk2 - https://github.com/layer5io/gowrk2/blob/v0.6.1/api/gowrk2.go#L47-L53
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options [filepath|json-string]
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options path/to/options.json
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator nighthawk --options '{"requests_per_second": 10, "max_pending_requests": 5}'
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator fortio --options '{"MethodOverride": "POST"}'
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator wrk2 --options '{"DurationInSeconds": 15, "Thread": 3}'
	`,
	Annotations: linkDocPerfApply,
	RunE: func(cmd *cobra.Command, args []string) error {
		userResponse := false

		// setting up for error formatting
		cmdUsed = "apply"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return ErrMesheryConfig(err)
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

			if loadTestBody == "" {
				loadTestBody = testClient.Body
			}
		}

		// Run test based on flags
		if testName == "" {
			utils.Log.Debug("Test Name not provided")
			testName = utils.StringWithCharset(8)
			utils.Log.Debug("Using random test name: ", testName)
		}

		// Throw error if a profile name is not provided
		if len(args) == 0 {
			return ErrNoProfileName()
		}

		// Invalid number of arguments
		if len(args) > 1 {
			return ErrorArgumentOverflow()
		}

		// handles spaces in args if quoted args passed
		for i, arg := range args {
			args[i] = strings.ReplaceAll(arg, " ", "%20")
		}
		// join all args to form profile name
		profileName = strings.Join(args, "%20")

		// Check if the profile name is valid, if not prompt the user to create a new one
		log.Debug("Fetching performance profile")
		profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), profileName, pageSize, pageNumber-1)
		if err != nil {
			return err
		}

		index := 0
		if len(profiles) == 0 {
			// if the provided performance profile does not exist, prompt the user to create a new one

			// skip asking confirmation if -y flag used
			if utils.SilentFlag {
				userResponse = true
			} else {
				// ask user for confirmation
				userResponse = utils.AskForConfirmation("Profile with name '" + profileName + "' does not exist. Do you want to create a new one")
			}

			if userResponse {
				profileID, profileName, err = createPerformanceProfile(mctlCfg)
				if err != nil {
					return err
				}
			} else {
				return ErrNoProfileFound()
			}
		} else {
			if len(profiles) == 1 { // if single performance profile found set profileID
				profileID = profiles[0].ID.String()
			} else { // multiple profiles found with matching name ask for profile index
				data := profilesToStringArrays(profiles)
				index, err = userPrompt("profile", "Enter index of the profile", data)
				if err != nil {
					return err
				}
				profileID = profiles[index].ID.String()
			}
			// what if user passed profile-name but didn't passed the url
			// we use url from performance profile
			if testURL == "" {
				testURL = profiles[index].Endpoints[0]
			}

			// reset profile name without %20
			// pull test configuration from the profile only if a test configuration is not provided and flag is not set
			if filePath == "" {
				profileName = profiles[index].Name

				if loadGenerator == "" {
					loadGenerator = profiles[index].LoadGenerators[0]
				}

				if concurrentRequests == "" {
					concurrentRequests = strconv.Itoa(profiles[index].ConcurrentRequest)
				}

				if qps == "" {
					qps = strconv.Itoa(profiles[index].QPS)
				}

				if testDuration == "" {
					testDuration = profiles[index].Duration
				}

				if testMesh == "" {
					testMesh = profiles[index].ServiceMesh
				}

				if loadTestBody == "" {
					loadTestBody = profiles[index].RequestBody
				}
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

		req, err = utils.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/run", nil)
		if err != nil {
			return err
		}

		q := req.URL.Query()

		if !disableCert {
			q.Add("cert", "true")
		}
		q.Add("name", testName)
		q.Add("loadGenerator", loadGenerator)
		q.Add("c", concurrentRequests)
		q.Add("url", testURL)
		q.Add("qps", qps)
		q.Add("reqBody", loadTestBody)

		durLen := len(testDuration)

		q.Add("dur", string(testDuration[durLen-1]))
		q.Add("t", string(testDuration[:durLen-1]))

		if testMesh != "" {
			q.Add("mesh", testMesh)
		}
		req.URL.RawQuery = q.Encode()

		utils.Log.Info("Initiating Performance test ...")

		resp, err := utils.MakeRequest(req)

		if err != nil {
			return err
		}

		defer utils.SafeClose(resp.Body)
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		utils.Log.Debug(string(data))

		utils.Log.Info("Test Completed Successfully!")
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVar(&testURL, "url", "", "(optional) Endpoint URL to test (required with --profile)")
	applyCmd.Flags().StringVar(&testName, "name", "", "(optional) Name of the Test")
	applyCmd.Flags().StringVar(&testMesh, "mesh", "", "(optional) Name of the Service Mesh")
	applyCmd.Flags().StringVar(&qps, "qps", "", "(optional) Queries per second")
	applyCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "", "(optional) Number of Parallel Requests")
	applyCmd.Flags().StringVar(&testDuration, "duration", "", "(optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration")
	applyCmd.Flags().StringVar(&loadGenerator, "load-generator", "", "(optional) Load-Generator to be used (fortio/wrk2/nighthawk)")
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "(optional) File containing SMP-compatible test configuration. For more, see https://github.com/layer5io/service-mesh-performance-specification")
	applyCmd.Flags().StringVarP(&loadTestBody, "body", "b", "", "(optional) Load test body. Can be a filepath/string")
	applyCmd.Flags().StringVar(&additionalOptions, "options", "", "(optional) Additional options to be passed to the load generator. Can be a json string or a filepath containing json")
	applyCmd.Flags().StringVar(&certPath, "cert-path", "", "(optional) Path to the certificate to be used for the load test")
	applyCmd.Flags().BoolVar(&disableCert, "disable-cert", false, "(optional) Do not use certificate present in the profile")
}

func createPerformanceProfile(mctlCfg *config.MesheryCtlConfig) (string, string, error) {
	utils.Log.Debug("Creating new performance profile inside function")

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

	if loadTestBody != "" {
		// Check if the loadTestBody is a filepath or a string
		if _, err := os.Stat(loadTestBody); err == nil {
			utils.Log.Info("Reading test body from file")
			bodyFile, err := os.ReadFile(loadTestBody)
			if err != nil {
				return "", "", ErrReadFilepath(err)
			}
			loadTestBody = string(bodyFile)
		}
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
		"request_body":       loadTestBody,
		"request_cookies":    "",
		"request_headers":    "",
		"content_type":       "",
	}

	if additionalOptions != "" {
		// Check if the additionalOptions is a filepath or a string
		if _, err := os.Stat(additionalOptions); err == nil {
			optFile, err := os.ReadFile(additionalOptions)
			if err != nil {
				return "", "", errors.New("unable to read options file. " + err.Error())
			}
			additionalOptions = string(optFile)
		}

		// Check if the additionalOptions is a valid json
		if !govalidator.IsJSON(additionalOptions) {
			return "", "", errors.New("invalid json passed as options")
		}

		values["metadata"] = map[string]interface{}{
			"additional_options": additionalOptions,
		}
	}

	if fileInfo, err := os.Stat(certPath); err == nil {
		certFile, err := os.ReadFile(certPath)
		if err != nil {
			return "", "", errors.New("unable to read certificate file. " + err.Error())
		}
		certData := string(certFile)
		certName := fileInfo.Name()

		// check if values["metadata"] is nil
		if values["metadata"] == nil {
			values["metadata"] = map[string]interface{}{}
		}

		values["metadata"].(map[string]interface{})["ca_certificate"] = map[string]interface{}{
			"name":     certName,
			"certName": certData,
		}
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", "", ErrFailMarshal(err)
	}
	req, err := utils.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles", bytes.NewBuffer(jsonValue))

	if err != nil {
		return "", "", err
	}

	resp, err := utils.MakeRequest(req)

	if err != nil {
		return "", "", err
	}

	var response *models.PerformanceProfile

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

	utils.Log.Debug("New profile created")
	return profileID, profileName, nil
}
