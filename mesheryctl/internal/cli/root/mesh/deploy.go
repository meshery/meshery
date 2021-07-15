package mesh

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	deployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "deploy a service mesh in the kubernetes cluster",
		Args:  cobra.MinimumNArgs(1),
		Long:  `deploy service mesh in the connected kubernetes cluster`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Verifying prerequisites...")
			if len(args) < 1 {
				return errors.New("no service mesh specified")
			}

			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			if err = validateAdapter(mctlCfg, tokenPath, args[0]); err != nil {
				return errors.Wrap(err, "adapter not valid")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := utils.CreateDefaultSpinner("Installation started", "\nInstallation complete")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			s.Start()
			_, err = sendDeployRequest(mctlCfg, args[0], false)
			if err != nil {
				return errors.Wrap(err, "error installing service mesh")
			}
			s.Stop()

			//log.Infof("Verifying Installation")
			//s.Start()
			//_, err = waitForDeployResponse(mctlCfg, args[0])
			//if err != nil {
			//	return errors.Wrap(err, "error verifying installation")
			//}
			//s.Stop()

			return nil
		},
	}
)

func init() {
	deployCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-istio:10000", "Adapter to use for installation")
	deployCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the validation tests and sample workload")
	deployCmd.Flags().StringVarP(&tokenPath, "tokenPath", "t", "", "Path to token for authenticating to Meshery API")
	_ = deployCmd.MarkFlagRequired("tokenPath")
}

func sendDeployRequest(mctlCfg *config.MesheryCtlConfig, query string, delete bool) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/mesh/ops"
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
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func waitForDeployResponse(mctlCfg *config.MesheryCtlConfig, query string) (string, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/events?client=cli"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, path, nil)
	if err != nil {
		return "", err
	}

	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func validateAdapter(mctlCfg *config.MesheryCtlConfig, tokenPath string, name string) error {
	prefs, err := utils.GetSessionData(mctlCfg, tokenPath)
	if err != nil {
		return err
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
		return errors.New("No adapters found")
	}

	prompt := promptui.Select{
		Label: "Select an Adapter from the list",
		Items: adapterNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return err
	}

	adapterURL = adapterNames[i]
	return nil
}
