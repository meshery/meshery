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

package adapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
)

var availableSubcommands []*cobra.Command

// errEventStreamClosed is the cause when the stream ends without a result.
var errEventStreamClosed = errors.New("event stream closed before the operation finished")

// AdapterCmd represents the Performance Management CLI command
var (
	adapterURL string
	namespace  string
	watch      bool
	AdapterCmd = &cobra.Command{
		Use:   "adapter",
		Short: "Connect and use Meshery adapters",
		Long: `Provisioning, configuration, and on-going operational management of cloud and cloud native infrastructure.
	Find more information at: https://docs.meshery.io/reference/mesheryctl#command-reference`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// if `adapter` command is run without any subcommands, show Help and exit
			if cmd.HasSubCommands() {
				return cmd.Help()
			}
			return nil
		},
	}
)

func init() {
	availableSubcommands = []*cobra.Command{validateCmd, deployCmd, removeCmd}
	AdapterCmd.AddCommand(availableSubcommands...)
	AdapterCmd.PersistentFlags().StringVarP(
		&utils.TokenFlag, "token", "t", "",
		"Path to token for authenticating to Meshery API",
	)
}

func validateAdapter(mctlCfg *config.MesheryCtlConfig, meshName string) error {
	// get details about the current meshery session
	prefs, err := utils.GetSessionData(mctlCfg)
	if err != nil {
		return ErrGettingSessionData(err)
	}

	// search for the mesh's adapter and update adapterURL accordingly
	for _, adapter := range prefs.MeshAdapters {
		if adapter.Name == meshName {
			adapterURL = adapter.Location
			return nil
		}
	}

	// return an error if the mesh's adapter was not found
	return ErrNoAdapters
}

func validateMesh(mctlCfg *config.MesheryCtlConfig, meshName string) (string, error) {
	// A mesh/technology is now identified by a Meshery Registry model name
	// (free-form) rather than a fixed service-mesh enum, so any explicitly
	// provided name is accepted as-is. This also avoids dropping into
	// interactive mode when the command is run by automation.
	if meshName != "" {
		return meshName, nil
	}

	// get details about the current meshery session
	prefs, err := utils.GetSessionData(mctlCfg)
	if err != nil {
		return "", ErrGettingSessionData(err)
	}

	// get a list of meshes from available adapters
	meshNames := []string{}
	for _, adapter := range prefs.MeshAdapters {
		meshNames = append(meshNames, adapter.Name)
	}

	// return an error if no adapters were found
	if len(meshNames) == 0 {
		return "", ErrNoAdapters
	}

	// allow the user to select a mesh from the list of available ones
	prompt := promptui.Select{
		Label: "Select a component from the list",
		Items: meshNames,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return "", ErrPrompt(err)
	}

	// return the selected mesh name
	return meshNames[i], nil
}

// parseAndValidateMesh parses positional args into a normalized mesh name,
// prompts the user if no mesh was specified, and validates that an adapter is available.
func parseAndValidateMesh(mctlCfg *config.MesheryCtlConfig, args []string) error {
	meshName = ""
	if len(args) > 0 {
		// if a mesh name was provided, convert it to the same format as adapter.Name
		// all args are joined, converted to upper case, and non-word characters replaced with "_"
		// examples:
		//     args = Linkerd -> ["LINKERD"] -> "LINKERD"
		//     args = nginx service mesh -> ["nginx", "service", "mesh"] -> "NGINX_SERVICE_MESH"
		r, _ := regexp.Compile(`\W`)
		meshName = r.ReplaceAllString(strings.ToUpper(strings.Join(args, "_")), "_")
	}

	var err error
	// verify the specified mesh is valid; if no mesh was specified, prompt the user
	meshName, err = validateMesh(mctlCfg, meshName)
	if err != nil {
		return err
	}

	// ensure the mesh's adapter is available and update adapterURL if so
	return validateAdapter(mctlCfg, meshName)
}

func sendOperationRequest(mctlCfg *config.MesheryCtlConfig, query string, delete bool, spec string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/system/adapter/operation"
	method := "POST"
	data := url.Values{}
	data.Set("adapter", adapterURL)
	data.Set("query", query)
	data.Set("customBody", "")
	data.Set("namespace", namespace)
	if delete {
		data.Set("deleteOp", "on")
	} else {
		data.Set("deleteOp", "")
	}

	// Choose which specification to use for conformance test
	// Used in validate command
	switch spec {
	case "smi":
		{
			data.Set("query", "smi_conformance")
			break
		}
	case "istio-vet":
		{
			if adapterURL == "meshery-istio:10000" {
				data.Set("query", "istio-vet")
				break
			}
			return "", errors.New("only Istio supports istio-vet operation")
		}
	case "null":
		{
			// If validate command isn't used
			// This must be used only for deploy or remove operations
			break
		}
	default:
		{
			return "", errors.New("specified specification not found or not yet supported")
		}
	}

	payload := strings.NewReader(data.Encode())

	client := &http.Client{}
	req, err := utils.NewRequest(method, path, payload)
	if err != nil {
		return "", ErrCreatingValidateRequest(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	res, err := client.Do(req)
	if err != nil {
		return "", ErrCreatingDeployRequest(err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", utils.ErrReadResponseBody(err)
	}
	// A rejected request carries no operation ID, so without this check
	// --watch would follow another operation's events.
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", utils.ErrResponseStatus(res.StatusCode)
	}
	return string(body), nil
}

// operationIDFromResponse returns the operation ID in the response to an
// adapter operation request, or "" when the server doesn't send one.
func operationIDFromResponse(body string) string {
	var response struct {
		OperationID string `json:"operationId"`
	}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return ""
	}
	return response.OperationID
}

// isOperationEvent reports whether an event may be the result of the
// operation. Without an operation ID on both sides, every event may be.
func isOperationEvent(data utils.EventData, operationID string) bool {
	return operationID == "" || data.OperationID == "" || data.OperationID == operationID
}

// checkEventStream returns an error unless res is an open event stream, so a
// response that can never carry a result fails right away.
func checkEventStream(res *http.Response) error {
	if res.StatusCode != http.StatusOK {
		return utils.ErrResponseStatus(res.StatusCode)
	}
	if contentType := res.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/event-stream") {
		return fmt.Errorf("expected an event stream, got content type %q", contentType)
	}
	return nil
}

// waitForDeployResponse waits on the event stream for the result of an adapter
// operation whose event details contain query. When operationID is set, events
// of other operations are ignored.
func waitForDeployResponse(mctlCfg *config.MesheryCtlConfig, query, operationID string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/system/adapter/events?client=cli_deploy"
	method := "GET"
	client := &http.Client{}
	req, err := utils.NewRequest(method, path, nil)
	if err != nil {
		return "", ErrCreatingDeployRequest(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", ErrCreatingDeployResponseRequest(err)
	}
	defer func() { _ = res.Body.Close() }()

	if err := checkEventStream(res); err != nil {
		return "", ErrCreatingDeployResponseStream(err)
	}

	event, err := utils.ConvertRespToSSE(res)
	if err != nil {
		return "", ErrCreatingDeployResponseStream(err)
	}

	timer := time.NewTimer(1200 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return "", ErrTimeoutWaitingForDeployResponse
		case item, ok := <-event:
			if !ok {
				return "", ErrCreatingDeployResponseStream(errEventStreamClosed)
			}
			if !isOperationEvent(item.Data, operationID) {
				continue
			}
			if strings.Contains(item.Data.Details, query) {
				utils.Log.Infof("%s\n%s\n", item.Data.Summary, item.Data.Details)
				return "", nil
			}
			if strings.Contains(item.Data.Details, "Error") {
				utils.Log.Infof("%s\n", item.Data.Summary)
				return "", ErrFailedDeployingMesh
			}
		}
	}
}

// waitForValidateResponse waits on the event stream for the result of a
// validation whose event summary contains query. When operationID is set,
// events of other operations are ignored.
func waitForValidateResponse(mctlCfg *config.MesheryCtlConfig, query, operationID string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/system/adapter/events?client=cli_validate"
	method := "GET"
	client := &http.Client{}
	req, err := utils.NewRequest(method, path, nil)
	if err != nil {
		return "", ErrCreatingValidateResponseRequest(err)
	}
	req.Header.Add("Accept", "text/event-stream")

	res, err := client.Do(req)
	if err != nil {
		return "", ErrCreatingValidateRequest(err)
	}
	defer func() { _ = res.Body.Close() }()

	if err := checkEventStream(res); err != nil {
		return "", ErrCreatingValidateResponseStream(err)
	}

	event, err := utils.ConvertRespToSSE(res)
	if err != nil {
		return "", ErrCreatingValidateResponseStream(err)
	}

	timer := time.NewTimer(1200 * time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			return "", ErrTimeoutWaitingForValidateResponse
		case item, ok := <-event:
			if !ok {
				return "", ErrCreatingValidateResponseStream(errEventStreamClosed)
			}
			if !isOperationEvent(item.Data, operationID) {
				continue
			}
			if strings.Contains(item.Data.Summary, query) {
				utils.Log.Infof("%s\n%s", item.Data.Summary, item.Data.Details)
				return "", nil
			}
			if strings.Contains(item.Data.Details, "error") {
				utils.Log.Infof("%s", item.Data.Summary)
				return "", ErrSMIConformanceTestsFailed
			}
		}
	}
}
