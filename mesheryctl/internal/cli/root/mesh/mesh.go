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

package mesh

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands []*cobra.Command
)

// MeshCmd represents the Performance Management CLI command
var (
	adapterURL string
	err        error
	mctlCfg    *config.MesheryCtlConfig
	namespace  string
	watch      bool
	MeshCmd    = &cobra.Command{
		Use:   "mesh",
		Short: "Cloud Native Lifecycle Management",
		Long:  "Provisioning, configuration, and on-going operational management of service meshes",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// if `mesh` command is ran without any subcommands, show Help and exit
			if cmd.HasSubCommands() {
				return cmd.Help()
			}

			// get the meshery config
			mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				utils.Log.Error(err)
				return nil
			}

			if len(args) > 0 {
				// if a mesh name was provided, convert it the same format as adapter.Name
				// all args are joined, converted to upper case, and non-word characters replaced with "_"
				// examples:
				//     args = Linkerd -> ["LINKERD"] -> "LINKERD"
				//     args = nginx service mesh -> ["nginx", "service", "mesh"] -> "NGINX_SERVICE_MESH"
				r, _ := regexp.Compile(`\W`)

				// If the mesh name is single word and without any spaces eg: istio, linkerd
				if len(args) == 1 {
					meshName = strings.ToUpper(args[0])
				}
				meshName = r.ReplaceAllString(strings.ToUpper(strings.Join(args, "_")), "_")
			}

			// verify the specified mesh is valid
			// if no mesh was specified, the user will be prompted to select one
			meshName, err = validateMesh(mctlCfg, meshName)
			if err != nil {
				utils.Log.Error(err)
			}

			// ensure the mesh's adapter is available and update adapterURL if so
			if err = validateAdapter(mctlCfg, meshName); err != nil {
				utils.Log.Error(err)
			}
			return nil
		},
	}
)

func init() {
	availableSubcommands = []*cobra.Command{validateCmd, deployCmd, removeCmd}
	MeshCmd.AddCommand(availableSubcommands...)
	MeshCmd.PersistentFlags().StringVarP(
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
	// if a mesh name is provided, verify it is valid
	if meshName != "" {
		if _, ok := smp.ServiceMesh_Type_value[meshName]; ok {
			return meshName, nil
		}
		// return an error if the provided mesh name is invalid
		// this prevents it from dropping into interactive mode
		// in case the command is being ran by automation
		return "", ErrValidMeshName(meshName)
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
		Label: "Select a Service Mesh from the list",
		Items: meshNames,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return "", ErrPrompt(err)
	}

	// return the selected mesh name
	return meshNames[i], nil
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
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", utils.ErrReadResponseBody(err)
	}
	return string(body), nil
}

func waitForDeployResponse(mctlCfg *config.MesheryCtlConfig, query string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/events?client=cli_deploy"
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
	defer res.Body.Close()

	event, err := utils.ConvertRespToSSE(res)
	if err != nil {
		return "", ErrCreatingDeployResponseStream(err)
	}

	timer := time.NewTimer(time.Duration(1200) * time.Second)
	eventChan := make(chan string)

	//Run a goroutine to wait for the response
	go func() {
		for i := range event {
			if strings.Contains(i.Data.Details, query) {
				eventChan <- "successful"
				log.Infof("%s\n%s\n", i.Data.Summary, i.Data.Details)
			} else if strings.Contains(i.Data.Details, "Error") {
				eventChan <- "error"
				log.Infof("%s\n", i.Data.Summary)
			}
		}
	}()

	select {
	case <-timer.C:
		return "", ErrTimeoutWaitingForDeployResponse
	case event := <-eventChan:
		if event != "successful" {
			return "", ErrFailedDeployingMesh
		}
	}

	return "", nil
}

func waitForValidateResponse(mctlCfg *config.MesheryCtlConfig, query string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/events?client=cli_validate"
	method := "GET"
	client := &http.Client{}
	req, err := utils.NewRequest(method, path, nil)
	req.Header.Add("Accept", "text/event-stream")
	if err != nil {
		return "", ErrCreatingDeployResponseRequest(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", ErrCreatingValidateRequest(err)
	}

	event, err := utils.ConvertRespToSSE(res)
	if err != nil {
		return "", ErrCreatingValidateResponseStream(err)
	}

	timer := time.NewTimer(time.Duration(1200) * time.Second)
	eventChan := make(chan string)

	//Run a goroutine to wait for the response
	go func() {
		for i := range event {
			if strings.Contains(i.Data.Summary, query) {
				eventChan <- "successful"
				log.Infof("%s\n%s", i.Data.Summary, i.Data.Details)
			} else if strings.Contains(i.Data.Details, "error") {
				eventChan <- "error"
				log.Infof("%s", i.Data.Summary)
			}
		}
	}()

	select {
	case <-timer.C:
		return "", ErrTimeoutWaitingForValidateResponse
	case event := <-eventChan:
		if event != "successful" {
			return "", ErrSMIConformanceTestsFailed
		}
	}

	return "", nil
}
