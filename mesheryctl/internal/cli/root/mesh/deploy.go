package mesh

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	smp "github.com/layer5io/service-mesh-performance/spec"
	"github.com/manifoldco/promptui"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	meshName  string
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service mesh to the Kubernetes cluster",
		Args:  cobra.MinimumNArgs(0),
		Long:  `Deploy a service mesh to the connected Kubernetes cluster`,
		Example: `
// Deploy a service mesh from an interactive on the default namespace
mesheryctl mesh deploy

// Deploy Linkerd mesh on a specific namespace
mesheryctl mesh deploy --adapter meshery-linkerd --namespace linkerd-ns

// Deploy Linkerd mesh and wait for it to be deployed
mesheryctl mesh deploy --adapter meshery-linkerd --watch`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				log.Fatalln(err)
			}

			// set default tokenpath for command.
			if tokenPath == "" {
				tokenPath = constants.GetCurrentAuthToken()
			}

			if len(args) < 1 {
				meshName, err = validateMesh(mctlCfg, tokenPath, "")
			} else {
				meshName, err = validateMesh(mctlCfg, tokenPath, args[0])
			}
			if err != nil {
				log.Fatalln(err)
			}

			if err = validateAdapter(mctlCfg, tokenPath, meshName); err != nil {
				// ErrValidatingAdapter
				log.Fatalln(err)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				log.Fatalln(err)
			}

			_, err = sendDeployRequest(mctlCfg, meshName, false)
			if err != nil {
				log.Fatalln(err)
			}

			if watch {
				log.Infof("Verifying Operation")
				_, err = waitForDeployResponse(mctlCfg, "mesh is now installed")
				if err != nil {
					log.Fatalln(err)
				}
			}

			return nil
		},
	}
)

func init() {
	deployCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-istio:10000", "Adapter to use for installation")
	deployCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the validation tests and sample workload")
	deployCmd.Flags().StringVarP(&tokenPath, "tokenPath", "t", "", "Path to token for authenticating to Meshery API")
	deployCmd.Flags().BoolVarP(&watch, "watch", "w", false, "Watch for events and verify operation (in beta testing)")
}

func sendDeployRequest(mctlCfg *config.MesheryCtlConfig, query string, delete bool) (string, error) {
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

	payload := strings.NewReader(data.Encode())

	client := &http.Client{}
	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		return "", ErrCreatingDeployRequest(err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", ErrAddingAuthDetails(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return "", ErrCreatingDeployRequest(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func waitForDeployResponse(mctlCfg *config.MesheryCtlConfig, query string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/events?client=cli_deploy"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return "", ErrCreatingDeployResponseRequest(err)
	}

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", ErrAddingAuthDetails(err)
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

func validateAdapter(mctlCfg *config.MesheryCtlConfig, tokenPath string, name string) error {
	prefs, err := utils.GetSessionData(mctlCfg, tokenPath)
	if err != nil {
		return ErrGettingSessionData(err)
	}

	adapterNames := []string{}
	for _, adapter := range prefs.MeshAdapters {
		if adapter.Name == name || adapter.Name == strings.ToUpper(name) {
			adapterNames = append(adapterNames, adapter.Location)
		}
		if adapter.Location == adapterURL {
			return nil
		}
	}

	if len(adapterNames) == 0 {
		return ErrNoAdapters
	}

	prompt := promptui.Select{
		Label: "Select an Adapter from the list",
		Items: adapterNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return ErrPrompt(err)
	}

	adapterURL = adapterNames[i]
	return nil
}

func validateMesh(mctlCfg *config.MesheryCtlConfig, tokenPath string, name string) (string, error) {
	if name != "" {
		if _, ok := smp.ServiceMesh_Type_value[name]; ok {
			return strings.ToLower(name), nil
		}
		if _, ok := smp.ServiceMesh_Type_value[strings.ToUpper(name)]; ok {
			return strings.ToLower(name), nil
		}
	}

	prefs, err := utils.GetSessionData(mctlCfg, tokenPath)
	if err != nil {
		// ErrGettingSessionData
		return "", ErrGettingSessionData(err)
	}

	meshNameMap := make(map[string]struct{})
	meshNames := []string{}
	for _, adapter := range prefs.MeshAdapters {
		if _, ok := meshNameMap[adapter.Name]; !ok {
			meshNames = append(meshNames, adapter.Name)
		}
	}

	if len(meshNames) == 0 {
		return "", ErrNoAdapters
	}

	prompt := promptui.Select{
		Label: "Select a Service Mesh from the list",
		Items: meshNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", ErrPrompt(err)
	}

	return strings.ToLower(meshNames[i]), nil
}
