package mesh

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

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
	meshName   = ""
	namespace  string
	watch      bool
	MeshCmd    = &cobra.Command{
		Use:   "mesh",
		Short: "Service Mesh Lifecycle Management",
		Long:  "Provisioning, configuration, and on-going operational management of service meshes",
		Example: `  // Deploy Linkerd:
  mesheryctl mesh deploy linkerd -n linkerd-ns -t /path/to/token -w
  
  // Run SMI conformance tests against Linkerd:
  mesheryctl mesh validate linkerd -s smi -n linkerd-ns
  
  // Remove Linkerd:
  mesheryctl mesh remove linkerd -n linkerd-ns
  
Run "mesheryctl mesh [command] -h" for more details on each command.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			// if `mesh` command is ran without any subcommands, show Help and exit
			if cmd.HasSubCommands() {
				return cmd.Help()
			}

			// get the meshery config
			mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				log.Fatalln(err)
			}

			if len(args) > 0 {
				// if a mesh name was provided, convert it the same format as adapter.Name
				// all args are joined, converted to upper case, and non-word characters replaced with "_"
				// examples:
				//     args = Linkerd -> ["LINKERD"] -> "LINKERD"
				//     args = nginx service mesh -> ["nginx", "service", "mesh"] -> "NGINX_SERVICE_MESH"
				r, _ := regexp.Compile(`\W`)
				meshName = r.ReplaceAllString(strings.ToUpper(strings.Join(args, "_")), "_")
			}

			// verify the specified mesh is valid
			// if no mesh was specified, the user will be prompted to select one
			meshName, err = validateMesh(mctlCfg, meshName)
			if err != nil {
				log.Fatalln(err)
			}

			// ensure the mesh's adapter is available and update adapterURL if so
			if err = validateAdapter(mctlCfg, meshName); err != nil {
				// ErrValidatingAdapter
				log.Fatalln(err)
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
		return "", fmt.Errorf("%s is not a valid mesh name or is unsupported", meshName)
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

func sendOperationRequest(mctlCfg *config.MesheryCtlConfig, query string, delete bool) (string, error) {
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
		return "", err
	}
	return string(body), nil
}
