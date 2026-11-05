// Copyright Meshery Authors
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
	"fmt"
	"io"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/ghodss/yaml"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	servermodels "github.com/meshery/meshery/server/models"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"

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
	"link":    "![perf-apply-usage](../../../images/perf-apply.png)",
	"caption": "Usage of mesheryctl perf apply",
}

// splitTestDuration converts a Go duration into the value and unit query
// params of a test run. The server only reads a whole number of hours, minutes
// or seconds, so the duration is always sent in seconds, and one that can't be
// expressed that way is rejected rather than run for a different length.
func splitTestDuration(duration string) (string, string, error) {
	normalized := strings.ToLower(strings.TrimSpace(duration))
	if normalized == "" {
		return "30", "s", nil
	}

	// A number without a unit means seconds.
	if _, err := strconv.Atoi(normalized); err == nil {
		normalized += "s"
	}

	parsed, err := time.ParseDuration(normalized)
	if err != nil || parsed < time.Second || parsed%time.Second != 0 {
		return "", "", utils.ErrFlagsInvalid(fmt.Errorf("invalid test duration %q: use a whole number of seconds of at least 1s, such as 30s, 5m or 1h30m", duration))
	}

	return strconv.FormatInt(int64(parsed/time.Second), 10), "s", nil
}

var applyCmd = &cobra.Command{
	Use:   "apply [profile-name]",
	Short: "Run a Performance test",
	Long: `Run Performance test using existing profiles or using flags.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/perf/apply`,
	Args: cobra.MinimumNArgs(0),
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

// Execute a Performance test with specified queries per second
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --qps 30

// Execute a Performance test with specified infrastructure
mesheryctl perf apply meshery-profile --url https://192.168.1.15/productpage --mesh istio

// Execute a Performance test creating a new performance profile and pass options to the load generator used
// If any options are already present in the profile or passed through flags, the --options flag will take precedence over the profile and flag options
// Options for fortio - https://github.com/fortio/fortio/blob/v1.57.0/fhttp/httprunner.go#L77-L84
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options [filepath|json-string]
mesheryctl perf apply meshery-profile-new --url "https://google.com" --options path/to/options.json
mesheryctl perf apply meshery-profile-new --url "https://google.com" --load-generator fortio --options '{"MethodOverride": "POST"}'
	`,
	Annotations: linkDocPerfApply,
	RunE: func(cmd *cobra.Command, args []string) error {
		userResponse := false
		isExistingProfile := false

		// setting up for error formatting
		cmdUsed = cmd.Name()

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		// Validate SMP configuration file early if provided
		if filePath != "" {
			if _, err := parseSMPConfigFile(filePath); err != nil {
				return err
			}
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
		utils.Log.Debug("Fetching performance profile")
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
				userResponse, err = utils.AskForConfirmation("Profile with name '" + profileName + "' does not exist. Do you want to create a new one")
				if err != nil {
					return err
				}
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

			params, err := resolvePerformanceParams(currentFlagParams(), filePath, &profiles[index], false)
			if err != nil {
				return err
			}
			params.ApplyToGlobals()
			profileName = profiles[index].Name
		}

		utils.Log.Debugf("performance profile is: %s", profileName)
		utils.Log.Debugf("test-url set to %s", testURL)

		req, err = utils.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/run", nil)
		if err != nil {
			return utils.ErrCreatingRequest(err)
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

		testDurationValue, testDurationUnit, err := splitTestDuration(testDuration)
		if err != nil {
			return err
		}
		q.Add("dur", testDurationUnit)
		q.Add("t", testDurationValue)

		if testMesh != "" {
			q.Add("mesh", testMesh)
		}
		req.URL.RawQuery = q.Encode()

		// Update the saved profile last, so a command that fails validation
		// leaves it unchanged.
		if isExistingProfile {
			if err := applyMetadataFlagsToProfile(mctlCfg, profileID); err != nil {
				return err
			}
		}

		utils.Log.Info("Initiating Performance test ...")

		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		defer utils.SafeClose(resp.Body)
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return utils.ErrReadResponseBody(err)
		}
		utils.Log.Debug(string(data))

		utils.Log.Info("Test Completed!")
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVar(&testURL, "url", "", "(optional) Endpoint URL to test (required with --profile)")
	applyCmd.Flags().StringVar(&testName, "name", "", "(optional) Name of the Test")
	applyCmd.Flags().StringVar(&testMesh, "mesh", "", "(optional) Name of the infrastructure")
	applyCmd.Flags().StringVar(&qps, "qps", "", "(optional) Queries per second")
	applyCmd.Flags().StringVar(&concurrentRequests, "concurrent-requests", "", "(optional) Number of Parallel Requests")
	applyCmd.Flags().StringVar(&testDuration, "duration", "", "(optional) Length of test (e.g. 10s, 5m, 2h). For more, see https://golang.org/pkg/time/#ParseDuration")
	applyCmd.Flags().StringVar(&loadGenerator, "load-generator", "", "(optional) Load-Generator to be used (fortio)")
	applyCmd.Flags().StringVarP(&filePath, "file", "f", "", "(optional) File containing SMP-compatible test configuration. For more, see https://smp-spec.io")
	applyCmd.Flags().StringVarP(&loadTestBody, "body", "b", "", "(optional) Load test body. Can be a filepath/string")
	applyCmd.Flags().StringVar(&additionalOptions, "options", "", "(optional) Additional options to be passed to the load generator. Can be a json string or a filepath containing json")
	applyCmd.Flags().StringVar(&certPath, "cert-path", "", "(optional) Path to the certificate to be used for the load test")
	applyCmd.Flags().BoolVar(&disableCert, "disable-cert", false, "(optional) Do not use certificate present in the profile")
}

// createPerformanceProfile creates a new performance profile on the Meshery server
// using resolved and validated parameters.
func createPerformanceProfile(mctlCfg *config.MesheryCtlConfig) (string, string, error) {
	utils.Log.Debug("Creating new performance profile inside function")

	if profileName == "" {
		return "", "", ErrNoProfileName()
	}

	params, err := resolvePerformanceParams(currentFlagParams(), filePath, nil, true)
	if err != nil {
		return "", "", err
	}
	params.ApplyToGlobals()

	convReq, _ := strconv.Atoi(concurrentRequests)
	convQPS, _ := strconv.Atoi(qps)

	values := map[string]interface{}{
		"concurrentRequest": convReq,
		"duration":          testDuration,
		"endpoints":         []string{testURL},
		"loadGenerators":    []string{loadGenerator},
		"name":              profileName,
		"qps":               convQPS,
		"serviceMesh":       testMesh,
		"requestBody":       loadTestBody,
		"requestCookies":    "",
		"requestHeaders":    "",
		"contentType":       "",
	}

	metadata, err := profileMetadataFromFlags()
	if err != nil {
		return "", "", err
	}
	if len(metadata) > 0 {
		values["metadata"] = metadata
	}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", "", ErrFailMarshal(err)
	}

	response, err := savePerformanceProfile(mctlCfg, jsonValue)
	if err != nil {
		return "", "", err
	}

	profileID = response.ID.String()
	profileName = response.Name

	utils.Log.Debug("New profile created")
	return profileID, profileName, nil
}

// profileMetadataFromFlags returns the profile metadata set by --options and
// --cert-path. It is empty when neither flag is set.
func profileMetadataFromFlags() (map[string]interface{}, error) {
	metadata := map[string]interface{}{}

	if additionalOptions != "" {
		options := additionalOptions
		// Check if the additional options are a filepath or a string
		if _, err := os.Stat(options); err == nil {
			optFile, err := os.ReadFile(options)
			if err != nil {
				return nil, ErrReadFilepath(errors.Wrap(err, "unable to read options file. "))
			}
			options = string(optFile)
		}

		// Check if the additional options are valid json
		if !govalidator.IsJSON(options) {
			return nil, ErrInvalidJSONFile()
		}

		metadata["additional_options"] = options
	}

	if certPath != "" {
		certFile, err := os.ReadFile(certPath)
		if err != nil {
			return nil, ErrReadFilepath(errors.New("Unable to read certificate file. " + err.Error()))
		}

		metadata["ca_certificate"] = map[string]interface{}{
			"name": filepath.Base(certPath),
			"file": string(certFile),
		}
	}

	return metadata, nil
}

// applyMetadataFlagsToProfile saves --options and --cert-path to an existing
// profile. The server reads both from the saved profile when a test runs, so
// they have no effect unless the profile is updated first.
func applyMetadataFlagsToProfile(mctlCfg *config.MesheryCtlConfig, profileID string) error {
	metadata, err := profileMetadataFromFlags()
	if err != nil {
		return err
	}
	if len(metadata) == 0 {
		return nil
	}

	// Saving replaces the whole profile and the profile list leaves out fields
	// such as the schedule, so start from the full profile.
	profile, err := api.Fetch[perfprofile.PerformanceProfile]("api/performance/profiles/" + profileID)
	if err != nil {
		return err
	}

	updatedMetadata := make(map[string]any, len(profile.Metadata)+len(metadata))
	maps.Copy(updatedMetadata, profile.Metadata)
	maps.Copy(updatedMetadata, metadata)
	profile.Metadata = updatedMetadata

	jsonValue, err := json.Marshal(profile)
	if err != nil {
		return ErrFailMarshal(err)
	}

	if _, err := savePerformanceProfile(mctlCfg, jsonValue); err != nil {
		return err
	}

	utils.Log.Infof("Updated performance profile %s with the provided options", profile.Name)
	return nil
}

// savePerformanceProfile creates a profile, or updates the existing one when
// the payload carries its id, and returns the profile the server saved.
func savePerformanceProfile(mctlCfg *config.MesheryCtlConfig, payload []byte) (*perfprofile.PerformanceProfile, error) {
	req, err := utils.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/performance/profiles", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, ErrPerfProfileServer(err)
	}

	var response *perfprofile.PerformanceProfile

	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrReadResponseBody(err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		serverErrMsg := strings.TrimSpace(string(body))
		if serverErrMsg == "" {
			serverErrMsg = fmt.Sprintf("server returned status: %s", resp.Status)
		}
		return nil, ErrPerfProfileServer(errors.New(serverErrMsg))
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, ErrFailUnmarshal(err)
	}

	if response == nil {
		return nil, ErrPerfProfileServer(errors.New("server returned empty performance profile"))
	}

	return response, nil
}

// PerformanceTestParams holds the consolidated performance test parameters
// resolved across CLI flags, configuration files, and performance profiles.
type PerformanceTestParams struct {
	Name               string
	URL                string
	Mesh               string
	QPS                string
	ConcurrentRequests string
	Duration           string
	LoadGenerator      string
	Body               string
}

// FallbackWith populates unset fields in p with non-empty values from fallback.
func (p *PerformanceTestParams) FallbackWith(fallback PerformanceTestParams) {
	if p.Name == "" {
		p.Name = fallback.Name
	}
	if p.URL == "" {
		p.URL = fallback.URL
	}
	if p.Mesh == "" {
		p.Mesh = fallback.Mesh
	}
	if p.QPS == "" {
		p.QPS = fallback.QPS
	}
	if p.ConcurrentRequests == "" {
		p.ConcurrentRequests = fallback.ConcurrentRequests
	}
	if p.Duration == "" {
		p.Duration = fallback.Duration
	}
	if p.LoadGenerator == "" {
		p.LoadGenerator = fallback.LoadGenerator
	}
	if p.Body == "" {
		p.Body = fallback.Body
	}
}

// Validate validates the parameters: URL must be non-empty and valid,
// QPS and ConcurrentRequests must be valid integers, and Body if pointing
// to an existing file on disk is loaded into memory.
func (p *PerformanceTestParams) Validate() error {
	if p.URL == "" {
		return ErrNoTestURL()
	}
	if !govalidator.IsURL(p.URL) {
		return ErrNotValidURL()
	}
	if p.QPS != "" {
		if _, err := strconv.Atoi(p.QPS); err != nil {
			return ErrConvertQPS()
		}
	}
	if p.ConcurrentRequests != "" {
		if _, err := strconv.Atoi(p.ConcurrentRequests); err != nil {
			return ErrConvertConcurrentRequest()
		}
	}
	if p.Body != "" {
		if _, err := os.Stat(p.Body); err == nil {
			utils.Log.Info("Reading test body from file")
			bodyFile, err := os.ReadFile(p.Body)
			if err != nil {
				return ErrReadFilepath(err)
			}
			p.Body = string(bodyFile)
		}
	}
	return nil
}

// ApplyToGlobals updates the package-level variables with the values in p.
func (p *PerformanceTestParams) ApplyToGlobals() {
	testName = p.Name
	testURL = p.URL
	testMesh = p.Mesh
	qps = p.QPS
	concurrentRequests = p.ConcurrentRequests
	testDuration = p.Duration
	loadGenerator = p.LoadGenerator
	loadTestBody = p.Body
}

// currentFlagParams returns the parameters currently set by command-line flags.
func currentFlagParams() PerformanceTestParams {
	return PerformanceTestParams{
		Name:               testName,
		URL:                testURL,
		Mesh:               testMesh,
		QPS:                qps,
		ConcurrentRequests: concurrentRequests,
		Duration:           testDuration,
		LoadGenerator:      loadGenerator,
		Body:               loadTestBody,
	}
}

// defaultProfileParams returns standard default parameters used for performance tests.
func defaultProfileParams() PerformanceTestParams {
	return PerformanceTestParams{
		Mesh:               "None",
		QPS:                "0",
		ConcurrentRequests: "1",
		Duration:           "30s",
		LoadGenerator:      "fortio",
	}
}

// parseSMPConfigFile reads and parses an SMP test configuration file into PerformanceTestParams.
func parseSMPConfigFile(path string) (PerformanceTestParams, error) {
	smpConfig, err := os.ReadFile(path)
	if err != nil {
		return PerformanceTestParams{}, ErrReadFilepath(errors.Wrap(err, "Unable to read test configuration file. \n"))
	}

	testConfig := servermodels.PerformanceTestConfigFile{}
	if err = yaml.Unmarshal(smpConfig, &testConfig); err != nil {
		return PerformanceTestParams{}, ErrFailUnmarshalFile(err)
	}

	if testConfig.Config == nil || len(testConfig.Config.Clients) == 0 {
		return PerformanceTestParams{}, ErrInvalidTestConfigFile()
	}

	client := testConfig.Config.Clients[0]
	return PerformanceTestParams{
		Name:               testConfig.Config.Name,
		URL:                firstString(client.EndpointUrls),
		Mesh:               testConfig.ServiceMesh,
		QPS:                strconv.Itoa(client.Rps),
		ConcurrentRequests: strconv.Itoa(client.Connections),
		Duration:           testConfig.Config.Duration,
		LoadGenerator:      client.LoadGenerator,
		Body:               client.Body,
	}, nil
}

// extractProfileParams converts a PerformanceProfile into PerformanceTestParams.
func extractProfileParams(p perfprofile.PerformanceProfile) PerformanceTestParams {
	return PerformanceTestParams{
		URL:                firstString(p.Endpoints),
		Mesh:               p.ServiceMesh,
		QPS:                strconv.Itoa(p.QPS),
		ConcurrentRequests: strconv.Itoa(p.ConcurrentRequest),
		Duration:           p.Duration,
		LoadGenerator:      firstString(p.LoadGenerators),
		Body:               p.RequestBody,
	}
}

// resolvePerformanceParams consolidates performance test parameters across CLI flags,
// an optional SMP configuration file, an optional performance profile, and defaults.
// Precedence: Flags > File > Profile > Defaults.
func resolvePerformanceParams(flags PerformanceTestParams, smpFilePath string, profile *perfprofile.PerformanceProfile, applyDefaults bool) (PerformanceTestParams, error) {
	resolved := flags

	if smpFilePath != "" {
		fileParams, err := parseSMPConfigFile(smpFilePath)
		if err != nil {
			return PerformanceTestParams{}, err
		}
		resolved.FallbackWith(fileParams)
	}

	if profile != nil {
		profileParams := extractProfileParams(*profile)
		// If an SMP file was provided, profile only fills URL if still unset;
		// otherwise, profile fills all unset parameters.
		if smpFilePath != "" {
			if resolved.URL == "" {
				resolved.URL = profileParams.URL
			}
		} else {
			resolved.FallbackWith(profileParams)
		}
	}

	if applyDefaults {
		resolved.FallbackWith(defaultProfileParams())
	}

	if resolved.Name == "" {
		utils.Log.Debug("Test Name not provided")
		resolved.Name = utils.StringWithCharset(8)
		utils.Log.Debug("Using random test name: ", resolved.Name)
	}

	if err := resolved.Validate(); err != nil {
		return PerformanceTestParams{}, err
	}

	return resolved, nil
}
