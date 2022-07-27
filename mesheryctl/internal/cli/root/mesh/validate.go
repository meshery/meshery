package mesh

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
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
var watch bool

// validateCmd represents the service mesh validation command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate conformance to service mesh standards",
	Example: `
// Validate conformance to service mesh standards
mesheryctl mesh validate --adapter [name of the adapter] --tokenPath [path to token for authentication] --spec [specification to be used for conformance test] --namespace [namespace to be used]

! Refer below image link for usage
* Usage of mesheryctl mesh validate
# ![mesh-validate-usage](/assets/img/mesheryctl/mesh-validate.png)
	`,
	Long: `Validate service mesh conformance to different standard specifications`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		log.Infof("Verifying prerequisites...")

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err)
		}

		prefs, err := utils.GetSessionData(mctlCfg)
		if err != nil {
			log.Fatalln(err)
		}
		//resolve adapterUrl to adapter Location
		for _, adapter := range prefs.MeshAdapters {
			adapterName := strings.Split(adapter.Location, ":")
			if adapterName[0] == adapterURL {
				adapterURL = adapter.Location
				meshName = adapter.Location
			}
		}
		//sync with available adapters
		if err = validateAdapter(mctlCfg, meshName); err != nil {
			log.Fatalln(err)
		}
		log.Info("verified prerequisites")
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return errors.New(utils.MeshError(fmt.Sprintf("'%s' is an invalid argument for 'mesheryctl mesh validate'. Use 'mesheryctl mesh validate --help' to display usage guide.\n", args[0])))
		}

		log.Infof("Starting service mesh validation...")

			// Choose which specification to use for conformance test
			var query string
			switch spec {
			case "smi":
				query = "smi_conformance"
			case "istio-vet":
				if adapterURL == "meshery-istio:10000" {
					query = "istio-vet"
					break
				}
				return errors.New("only Istio supports istio-vet operation")
			default:
				return errors.New("specification not found or not yet supported")
			}

			_, err = sendOperationRequest(mctlCfg, query, false)
			if err != nil {
				log.Fatalln(err)
			}

			if watch {
				log.Infof("Verifying Operation")
				_, err = waitForValidateResponse(mctlCfg, "SMI conformance test")
				if err != nil {
					log.Fatalln(err)
				}
			}

			return nil
		},
	}
)

func init() {
	validateCmd.Flags().StringVarP(
		&spec,
		"spec",
		"s",
		"smi",
		"Specification to be used for conformance test [smi (default), istio-vet]",
	)
	validateCmd.Flags().StringVarP(
		&namespace, "namespace", "n", "default",
		"Kubernetes namespace where the mesh is deployed",
	)
	validateCmd.Flags().BoolVarP(
		&watch, "watch", "w", false,
		"Watch for events and verify operation (in beta testing)",
	)
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
