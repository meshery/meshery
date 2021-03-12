package mesh

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Operation is the common body type to be passed for Mesh Ops
type Operation struct {
	Adapter    string `json:"adapter"`
	CustomBody string `json:"customBody"`
	DeleteOp   string `json:"deleteOp"`
	Namespace  string `json:"namespace"`
	Query      string `json:"query"`
}

var spec string
var adapterURL string
var namespace string
var tokenPath string
var err error

// validateCmd represents the service mesh validation command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate conformance to service mesh standards",
	Args:  cobra.NoArgs,
	Long:  `Validate service mesh conformance to different standard specifications`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		log.Infof("Starting service mesh validation...")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		path := mctlCfg.GetBaseMesheryURL() + "/api/mesh/ops"
		method := "POST"

		data := url.Values{}
		data.Set("adapter", adapterURL)
		data.Set("customBody", "")
		data.Set("deleteOp", "")
		data.Set("namespace", namespace)

		// Choose which specification to use for conformance test
		switch spec {
		case "smi":
			{
				data.Set("query", "smiConformanceTest")
				break
			}
		case "smp":
			{
				return errors.New("support for SMP coming in a future release")
			}
		default:
			{
				return errors.New("specified specification not found or not yet supported")
			}
		}

		payload := strings.NewReader(data.Encode())

		client := &http.Client{}
		req, err := http.NewRequest(method, path, payload)

		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		log.Infof(string(body))

		return nil
	},
}

func init() {
	validateCmd.Flags().StringVarP(&spec, "spec", "s", "smi", "specification to be used for conformance test")
	_ = validateCmd.MarkFlagRequired("spec")
	validateCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-osm:10010", "Adapter to use for validation")
	_ = validateCmd.MarkFlagRequired("adapter")
	validateCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the validation tests and sample workload")
	validateCmd.Flags().StringVarP(&tokenPath, "tokenPath", "t", "", "Path to token for authenticating to Meshery API")
	_ = validateCmd.MarkFlagRequired("tokenPath")
}
