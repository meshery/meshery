package mesh

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
var token string
var mesheryProvider string
var err error

// validateCmd represents the service mesh validation command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate Service Mesh",
	Args:  cobra.NoArgs,
	Long:  `Validate Service Meshes based on specified specification`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {

		log.Infof("Mesh Validation check started")

		path := "http://localhost:9081/api/mesh/ops"
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
		default:
			{
				return errors.New("specification not found")
			}
		}

		payload := strings.NewReader(data.Encode())

		client := &http.Client{}
		req, err := http.NewRequest(method, path, payload)

		if err != nil {
			return err
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded;charset=UTF-8")
		req.AddCookie(&http.Cookie{
			Name:  "meshery-provider",
			Value: mesheryProvider,
		})
		req.AddCookie(&http.Cookie{
			Name:  "token",
			Value: token,
		})

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
	validateCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-osm:10010", "Adapter url used for Conformance")
	_ = validateCmd.MarkFlagRequired("adapter")
	validateCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the workload")
	validateCmd.Flags().StringVarP(&token, "token", "t", "", "Token used for authenticating Meshery API")
	_ = validateCmd.MarkFlagRequired("token")
	validateCmd.Flags().StringVarP(&mesheryProvider, "provider", "p", "", "Provider used for Meshery")
	_ = validateCmd.MarkFlagRequired("provider")
}
