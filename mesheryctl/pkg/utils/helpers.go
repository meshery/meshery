package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/browser"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

const (
	// Meshery Docker Deployment URLs
	dockerComposeWebURL         = "https://api.github.com/repos/docker/compose/releases/latest"
	defaultDockerComposeVersion = "1.24.1/docker-compose"
	dockerComposeBinaryURL      = "https://github.com/docker/compose/releases/download/"
	dockerComposeBinary         = "/usr/local/bin/docker-compose"

	// Meshery Kubernetes Deployment URLs
	baseConfigURL = "https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/"
	OperatorURL   = baseConfigURL + "manifests/default.yaml"
	BrokerURL     = baseConfigURL + "samples/meshery_v1alpha1_broker.yaml"
	MeshsyncURL   = baseConfigURL + "samples/meshery_v1alpha1_meshsync.yaml"

	// Documentation URLs
	docsBaseURL       = "https://docs.meshery.io/"
	rootUsageURL      = docsBaseURL + "reference/mesheryctl"
	perfUsageURL      = docsBaseURL + "reference/mesheryctl/perf"
	systemUsageURL    = docsBaseURL + "reference/mesheryctl/system"
	systemStopURL     = docsBaseURL + "reference/mesheryctl/system/stop"
	systemUpdateURL   = docsBaseURL + "reference/mesheryctl/system/update"
	systemResetURL    = docsBaseURL + "reference/mesheryctl/system/reset"
	systemStatusURL   = docsBaseURL + "reference/mesheryctl/system/status"
	systemRestartURL  = docsBaseURL + "reference/mesheryctl/system/restart"
	meshUsageURL      = docsBaseURL + "reference/mesheryctl/mesh"
	expUsageURL       = docsBaseURL + "reference/mesheryctl/exp"
	filterUsageURL    = docsBaseURL + "reference/mesheryctl/filter"
	filterImportURL   = docsBaseURL + "reference/mesheryctl/filter/import"
	filterDeleteURL   = docsBaseURL + "reference/mesheryctl/filter/delete"
	filterListURL     = docsBaseURL + "reference/mesheryctl/filter/list"
	filterViewURL     = docsBaseURL + "reference/mesheryctl/filter/view"
	patternUsageURL   = docsBaseURL + "reference/mesheryctl/pattern"
	patternViewURL    = docsBaseURL + "reference/mesheryctl/pattern/view"
	appUsageURL       = docsBaseURL + "reference/mesheryctl/app"
	appViewURL        = docsBaseURL + "reference/mesheryctl/app/view"
	contextDeleteURL  = docsBaseURL + "reference/mesheryctl/system/context/delete"
	contextViewURL    = docsBaseURL + "reference/mesheryctl/system/context/view"
	contextCreateURL  = docsBaseURL + "reference/mesheryctl/system/context/create"
	contextUsageURL   = docsBaseURL + "reference/mesheryctl/system/context"
	channelUsageURL   = docsBaseURL + "reference/mesheryctl/system/channel"
	channelSetURL     = docsBaseURL + "reference/mesheryctl/system/channel/set"
	channelSwitchURL  = docsBaseURL + "reference/mesheryctl/system/channel/switch"
	channelViewURL    = docsBaseURL + "reference/mesheryctl/system/channel/view"
	providerUsageURL  = docsBaseURL + "reference/mesheryctl/system/provider"
	providerViewURL   = docsBaseURL + "reference/mesheryctl/system/provider/view"
	providerListURL   = docsBaseURL + "reference/mesheryctl/system/provider/list"
	providerSetURL    = docsBaseURL + "reference/mesheryctl/system/provider/set"
	providerResetURL  = docsBaseURL + "reference/mesheryctl/system/provider/reset"
	providerSwitchURL = docsBaseURL + "reference/mesheryctl/system/provider/switch"
	tokenUsageURL     = docsBaseURL + "reference/mesheryctl/system/token"
	modelUsageURL     = docsBaseURL + "reference/mesheryctl/system/model"
	modelListURL      = docsBaseURL + "reference/mesheryctl/system/model/list"
	modelViewURL      = docsBaseURL + "reference/mesheryctl/system/model/view"

	// Meshery Server Location
	EndpointProtocol = "http"
)

type cmdType string

const (
	cmdRoot           cmdType = "root"
	cmdPerf           cmdType = "perf"
	cmdMesh           cmdType = "mesh"
	cmdSystem         cmdType = "system"
	cmdSystemStop     cmdType = "system stop"
	cmdSystemUpdate   cmdType = "system update"
	cmdSystemReset    cmdType = "system reset"
	cmdSystemStatus   cmdType = "system status"
	cmdSystemRestart  cmdType = "system restart"
	cmdExp            cmdType = "exp"
	cmdFilter         cmdType = "filter"
	cmdFilterImport   cmdType = "filter import"
	cmdFilterDelete   cmdType = "filter delete"
	cmdFilterList     cmdType = "filter list"
	cmdFilterView     cmdType = "filter view"
	cmdPattern        cmdType = "pattern"
	cmdPatternView    cmdType = "pattern view"
	cmdApp            cmdType = "app"
	cmdAppView        cmdType = "app view"
	cmdContext        cmdType = "context"
	cmdContextDelete  cmdType = "delete"
	cmdContextCreate  cmdType = "create"
	cmdContextView    cmdType = "context view"
	cmdChannel        cmdType = "channel"
	cmdChannelSet     cmdType = "channel set"
	cmdChannelSwitch  cmdType = "channel switch"
	cmdChannelView    cmdType = "channel view"
	cmdProvider       cmdType = "provider"
	cmdProviderSet    cmdType = "provider set"
	cmdProviderSwitch cmdType = "provider switch"
	cmdProviderView   cmdType = "provider view"
	cmdProviderList   cmdType = "provider list"
	cmdProviderReset  cmdType = "provider reset"
	cmdToken          cmdType = "token"
	cmdModel          cmdType = "model"
	cmdModelList      cmdType = "model list"
	cmdModelView      cmdType = "model view"
)

const (
	HelmChartURL          = "https://meshery.io/charts/"
	HelmChartName         = "meshery"
	HelmChartOperatorName = "meshery-operator"
)

var (
	// ResetFlag indicates if a reset is required
	ResetFlag bool
	// SkipResetFlag indicates if fetching the updated manifest files is required
	SkipResetFlag bool
	// MesheryDefaultHost is the default host on which Meshery is exposed
	MesheryDefaultHost = "localhost"
	// MesheryDefaultPort is the default port on which Meshery is exposed
	MesheryDefaultPort = 9081
	// MesheryEndpoint is the default URL in which Meshery is exposed
	MesheryEndpoint = fmt.Sprintf("http://%s:%v", MesheryDefaultHost, MesheryDefaultPort)
	// MesheryFolder is the default relative location of the meshery config
	// related configuration files.
	MesheryFolder = ".meshery"
	// DockerComposeFile is the default location within the MesheryFolder
	// where the docker compose file is located.
	DockerComposeFile = "meshery.yaml"
	// AuthConfigFile is the location of the auth file for performing perf testing
	AuthConfigFile = "auth.json"
	// DefaultConfigPath is the detail path to mesheryctl config
	DefaultConfigPath = "config.yaml"
	// MesheryNamespace is the namespace to which Meshery is deployed in the Kubernetes cluster
	MesheryNamespace = "meshery"
	// MesheryDeployment is the name of a Kubernetes manifest file required to setup Meshery
	// check https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
	MesheryDeployment = "meshery-deployment.yaml"
	// MesheryService is the name of a Kubernetes manifest file required to setup Meshery
	// check https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
	MesheryService = "meshery-service.yaml"
	//MesheryOperator is the file for default Meshery operator
	//check https://github.com/layer5io/meshery-operator/blob/master/config/manifests/default.yaml
	MesheryOperator = "default.yaml"
	//MesheryOperatorBroker is the file for the Meshery broker
	//check https://github.com/layer5io/meshery-operator/blob/master/config/samples/meshery_v1alpha1_broker.yaml
	MesheryOperatorBroker = "meshery_v1alpha1_broker.yaml"
	//MesheryOperatorMeshsync is the file for the Meshery Meshsync Operator
	//check https://github.com/layer5io/meshery-operator/blob/master/config/samples/meshery_v1alpha1_meshsync.yaml
	MesheryOperatorMeshsync = "meshery_v1alpha1_meshsync.yaml"
	// ServiceAccount is the name of a Kubernetes manifest file required to setup Meshery
	// check https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
	ServiceAccount = "service-account.yaml"
	// To upload with param name
	ParamName = "k8sfile"
	// kubeconfig file name
	KubeConfigYaml = "kubeconfig.yaml"
	// ViperCompose is an instance of viper for docker-compose
	ViperCompose = viper.New()
	// ViperMeshconfig is an instance of viper for the meshconfig file
	ViperMeshconfig = viper.New()
	// SilentFlag skips waiting for user input and proceeds with default options
	SilentFlag bool
	// PlatformFlag sets the platform for the initial config file
	PlatformFlag string
	// Paths to kubeconfig files
	ConfigPath string
	KubeConfig string
	// KeepNamespace indicates if the namespace should be kept when Meshery is uninstalled
	KeepNamespace bool
	// TokenFlag sets token location passed by user with --token
	TokenFlag = "Not Set"
	// global logger variable
	Log logger.Handler
)

var CfgFile string

// TODO: add "meshery-perf" as a component

// ListOfComponents returns the list of components available
var ListOfComponents = []string{}

// TemplateContext is the template context provided when creating a config file
var TemplateContext = config.Context{
	Endpoint:   EndpointProtocol + "://localhost:9081",
	Token:      "default",
	Platform:   "kubernetes",
	Components: ListOfComponents,
	Channel:    "stable",
	Version:    "latest",
	Provider:   "Meshery",
}

var Services = map[string]Service{
	"meshery": {
		Image:  "layer5/meshery:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Environment: []string{
			"PROVIDER_BASE_URLS=https://meshery.layer5.io",
			"ADAPTER_URLS=meshery-istio:10000 meshery-linkerd:10001 meshery-consul:10002 meshery-nsm:10004 meshery-app-mesh:10005 meshery-kuma:10007 meshery-osm:10009 meshery-traefik-mesh:10006 meshery-nginx-sm:10010 meshery-cilium:10012",
			"EVENT=mesheryLocal",
			"PORT=9081",
		},
		Volumes: []string{"$HOME/.kube:/home/appuser/.kube:ro", "$HOME/.minikube:$HOME/.minikube:ro"},
		Ports:   []string{"9081:9081"},
	},
	"meshery-istio": {
		Image:  "layer5/meshery-istio:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10000:10000"},
	},
	"meshery-linkerd": {
		Image:  "layer5/meshery-linkerd:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10001:10001"},
	},
	"meshery-consul": {
		Image:  "layer5/meshery-consul:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10002:10002"},
	},
	"meshery-nsm": {
		Image:  "layer5/meshery-nsm:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10004:10004"},
	},
	"meshery-app-mesh": {
		Image:  "layer5/meshery-app-mesh:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10005:10005"},
	},
	"meshery-traefik-mesh": {
		Image:  "layer5/meshery-traefik-mesh:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10006:10006"},
	},
	"meshery-kuma": {
		Image:  "layer5/meshery-kuma:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10007:10007"},
	},
	"meshery-osm": {
		Image:  "layer5/meshery-osm:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10009:10009"},
	},
	"meshery-nginx-sm": {
		Image:  "layer5/meshery-nginx-sm:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10010:10010"},
	},
	"meshery-cilium": {
		Image:  "layer5/meshery-cilium:stable-latest",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
		Ports:  []string{"10012:10012"},
	},
	"watchtower": {
		Image:  "containrrr/watchtower",
		Labels: []string{"com.centurylinklabs.watchtower.enable=true"},
	},
}

// TemplateToken is the template token provided when creating a config file
var TemplateToken = config.Token{
	Name:     "default",
	Location: AuthConfigFile,
}

func BackupConfigFile(cfgFile string) {
	dir, file := filepath.Split(cfgFile)
	extension := filepath.Ext(file)
	bakLocation := filepath.Join(dir, file[:len(file)-len(extension)]+".bak.yaml")
	err := os.Rename(cfgFile, bakLocation)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Create(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
}

const tokenName = "token"
const providerName = "meshery-provider"

var seededRand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// StringWithCharset generates a random string with a given length
func StringWithCharset(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	// + "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// SafeClose is a helper function help to close the io
func SafeClose(co io.Closer) {
	if cerr := co.Close(); cerr != nil {
		log.Error(cerr)
	}
}

func prereq() ([]byte, []byte, error) {
	ostype, err := exec.Command("uname", "-s").Output()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not find os type")
	}

	osarch, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return nil, nil, errors.Wrap(err, "could not find os arch type")
	}

	return ostype, osarch, nil
}

// SetFileLocation to set absolute path
func SetFileLocation() error {
	// Find home directory.
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrap(err, "failed to get users home directory")
	}
	MesheryFolder = filepath.Join(home, MesheryFolder)
	DockerComposeFile = filepath.Join(MesheryFolder, DockerComposeFile)
	AuthConfigFile = filepath.Join(MesheryFolder, AuthConfigFile)
	DefaultConfigPath = filepath.Join(MesheryFolder, DefaultConfigPath)
	return nil
}

// NavigateToBroswer naviagtes to the endpoint displaying Meshery UI in the broswer.
func NavigateToBrowser(endpoint string) error {
	err := browser.OpenURL(endpoint)
	return err
}

// UploadFileWithParams returns a request configured to upload files with other values
func UploadFileWithParams(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if err = file.Close(); err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	_, err = part.Write(fileContents)
	if err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

// IsValidSubcommand checks if the passed subcommand is supported by the parent command
func IsValidSubcommand(available []*cobra.Command, sub string) bool {
	for _, s := range available {
		if sub == s.CalledAs() {
			return true
		}
	}
	return false
}

// ContentTypeIsHTML Checks if the response is an HTML resposnse
func ContentTypeIsHTML(resp *http.Response) bool {
	ctString := strings.Split(resp.Header.Get("Content-Type"), ";")
	if len(ctString) < 1 {
		return false
	}
	if ctString[0] == "text/html" {
		return true
	}
	return false
}

// UpdateMesheryContainers runs the update command for meshery client
func UpdateMesheryContainers() error {
	log.Info("Updating Meshery now...")

	start := exec.Command("docker-compose", "-f", DockerComposeFile, "pull")
	start.Stdout = os.Stdout
	start.Stderr = os.Stderr
	if err := start.Run(); err != nil {
		return errors.Wrap(err, SystemError("failed to start meshery"))
	}
	return nil
}

// AskForConfirmation asks the user for confirmation. A user must type in "yes" or "no" and then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as confirmations. If the input is not recognized, it will ask again. The function does not return until it gets a valid response from the user.
func AskForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]? ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}

// CreateConfigFile creates config file in Meshery Folder
func CreateConfigFile() error {
	if _, err := os.Stat(DefaultConfigPath); os.IsNotExist(err) {
		_, err := os.Create(DefaultConfigPath)
		if err != nil {
			return err
		}
	}
	return nil
}

// ValidateURL validates url provided for meshery backend to mesheryctl context
func ValidateURL(URL string) error {
	ParsedURL, err := url.ParseRequestURI(URL)
	if err != nil {
		return err
	}
	if ParsedURL.Scheme != "http" && ParsedURL.Scheme != "https" {
		return fmt.Errorf("%s is not a supported protocol", ParsedURL.Scheme)
	}
	return nil
}

// TruncateID shortens an id to 8 characters
func TruncateID(id string) string {
	ShortenedID := id[0:8]
	return ShortenedID
}

// PrintToTable prints the given data into a table format
func PrintToTable(header []string, data [][]string) {
	// The tables are formatted to look similar to how it looks in say `kubectl get deployments`
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header) // The header of the table
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // The data in the table
	table.Render()         // Render the table
}

// PrintToTableWithFooter prints the given data into a table format but with a footer
func PrintToTableWithFooter(header []string, data [][]string, footer []string) {
	// The tables are formatted to look similar to how it looks in say `kubectl get deployments`
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header) // The header of the table
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // The data in the table
	table.SetFooter(footer)
	table.Render() // Render the table
}

// StringContainedInSlice returns the index in which a string is a substring in a list of strings
func StringContainedInSlice(str string, slice []string) int {
	for index, ele := range slice {
		// Return index even if only a part of the string is present
		if strings.Contains(ele, str) {
			return index
		}
	}
	return -1
}

// StringInSlice checks if a string is present in a slice
func StringInSlice(str string, slice []string) bool {
	for _, ele := range slice {
		if ele == str {
			return true
		}
	}
	return false
}

// GetID returns a array of IDs from meshery server endpoint /api/{configurations}
func GetID(mesheryServerUrl, configuration string) ([]string, error) {
	url := mesheryServerUrl + "/api/" + configuration + "?page_size=10000"
	configType := configuration + "s"
	var idList []string
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return idList, err
	}

	res, err := MakeRequest(req)
	if err != nil {
		return idList, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return idList, ErrReadResponseBody(err)
	}
	var dat map[string]interface{}
	if err = json.Unmarshal(body, &dat); err != nil {
		return idList, ErrUnmarshal(errors.Wrap(err, "failed to unmarshal response body"))
	}
	if dat == nil {
		return idList, ErrNotFound(errors.New("no data found"))
	}
	if dat[configType] == nil {
		return idList, ErrNotFound(errors.New("no results found"))
	}
	for _, config := range dat[configType].([]interface{}) {
		idList = append(idList, config.(map[string]interface{})["id"].(string))
	}
	return idList, nil
}

// GetName returns a of name:id from meshery server endpoint /api/{configurations}
func GetName(mesheryServerUrl, configuration string) (map[string]string, error) {
	url := mesheryServerUrl + "/api/" + configuration + "?page_size=10000"
	configType := configuration + "s"
	nameIdMap := make(map[string]string)
	req, err := NewRequest("GET", url, nil)
	if err != nil {
		return nameIdMap, err
	}

	res, err := MakeRequest(req)
	if err != nil {
		return nameIdMap, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nameIdMap, ErrReadResponseBody(err)
	}
	var dat map[string]interface{}
	if err = json.Unmarshal(body, &dat); err != nil {
		return nameIdMap, ErrUnmarshal(errors.Wrap(err, "failed to unmarshal response body"))
	}
	if dat == nil {
		return nameIdMap, ErrNotFound(errors.New("no data found"))
	}
	if dat[configType] == nil {
		return nameIdMap, ErrNotFound(errors.New("no results found"))
	}
	for _, config := range dat[configType].([]interface{}) {
		nameIdMap[config.(map[string]interface{})["name"].(string)] = config.(map[string]interface{})["id"].(string)
	}
	return nameIdMap, nil
}

// Delete configuration from meshery server endpoint /api/{configurations}/{id}
func DeleteConfiguration(mesheryServerUrl, id, configuration string) error {
	url := mesheryServerUrl + "/api/" + configuration + "/" + id
	req, err := NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = MakeRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// ValidId - Check if args is a valid ID or a valid ID prefix and returns the full ID
func ValidId(mesheryServerUrl, args string, configuration string) (string, bool, error) {
	isID := false
	configID, err := GetID(mesheryServerUrl, configuration)
	if err == nil {
		for _, id := range configID {
			if strings.HasPrefix(id, args) {
				args = id
			}
		}
	} else {
		return "", false, err
	}
	isID, err = regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$", args)
	if err != nil {
		return "", false, ErrInvalidNameOrID(err)
	}
	return args, isID, nil
}

// ValidId - Check if args is a valid name or a valid name prefix and returns the full name and ID
func ValidName(mesheryServerUrl, args string, configuration string) (string, string, bool, error) {
	isName := false
	nameIdMap, err := GetName(mesheryServerUrl, configuration)

	if err != nil {
		return "", "", false, err
	}

	fullName := ""
	ID := ""

	for name := range nameIdMap {
		if strings.HasPrefix(name, args) {
			fullName = name
			ID = nameIdMap[name]
			isName = true
		}
	}

	return fullName, ID, isName, nil
}

// AskForInput asks the user for an input and checks if it is in the available values
func AskForInput(prompt string, allowed []string) string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s %s: ", prompt, allowed)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if StringInSlice(response, allowed) {
			return response
		}
		log.Fatalf("Invalid respose %s. Allowed responses %s", response, allowed)
	}
}

// ParseURLGithub checks URL and returns raw repo, path, error
func ParseURLGithub(URL string) (string, string, error) {
	// GitHub URL:
	// - https://github.com/layer5io/meshery/blob/master/.goreleaser.yml
	// - https://raw.githubusercontent.com/layer5io/meshery/master/.goreleaser.yml
	parsedURL, err := url.Parse(URL)
	if err != nil {
		return "", "", ErrParsingUrl(fmt.Errorf("failed to retrieve file from URL: %s", URL))
	}
	host := parsedURL.Host
	path := parsedURL.Path
	path = strings.Replace(path, "/blob/", "/", 1)
	paths := strings.Split(path, "/")
	if host == "github.com" {
		if len(paths) < 5 {
			return "", "", ErrParsingUrl(fmt.Errorf("failed to retrieve file from URL: %s", URL))
		}
		resURL := "https://" + host + strings.Join(paths[:4], "/")
		return resURL, strings.Join(paths[4:], "/"), nil
	} else if host == "raw.githubusercontent.com" {
		if len(paths) < 5 {
			return "", "", ErrParsingUrl(fmt.Errorf("failed to retrieve file from URL: %s", URL))
		}
		resURL := "https://" + "raw.githubusercontent.com" + path
		return resURL, "", nil
	}
	return URL, "", ErrParsingUrl(errors.New("only github urls are supported"))
}

// PrintToTableInStringFormat prints the given data into a table format but return as a string
func PrintToTableInStringFormat(header []string, data [][]string) string {
	// The tables are formatted to look similar to how it looks in say `kubectl get deployments`
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader(header) // The header of the table
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data) // The data in the table
	table.Render()         // Render the table

	return tableString.String()
}

// Indicate an ongoing Process at a given time on CLI
func CreateDefaultSpinner(suffix string, finalMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Suffix = " " + suffix
	s.FinalMSG = finalMsg + "\n"
	return s
}

// Get Meshery Session Data/Details (Adapters)
func GetSessionData(mctlCfg *config.MesheryCtlConfig) (*models.Preference, error) {
	path := mctlCfg.GetBaseMesheryURL() + "/api/system/sync"
	method := "GET"
	client := &http.Client{}
	req, err := NewRequest(method, path, nil)
	if err != nil {
		return nil, ErrCreatingRequest(err)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, ErrRequestResponse(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, ErrReadResponseBody(err)
	}

	prefs := &models.Preference{}
	err = utils.Unmarshal(string(body), prefs)
	if err != nil {
		return nil, errors.New("Failed to process JSON data. Please sign into Meshery")
	}

	return prefs, nil
}

// ContainsStringPrefix takes a string slice and a string and returns true if it is present
func ContainsStringPrefix(arr []string, str string) bool {
	for _, el := range arr {
		if strings.HasPrefix(el, str) {
			return true
		}
	}

	return false
}

// TransformYAML takes in:
//
//	yamlByt - YAML Byte slice that needs to be modified
//	transform - function that will be executed on that value, the returned value will replace the current value
//	keys - takes in a series of keys which are supposed to be nested, numbers can also be passed to access an array
func TransformYAML(yamlByt []byte, transform func(interface{}) (interface{}, error), keys ...string) ([]byte, error) {
	var data map[string]interface{}

	err := yaml.Unmarshal(yamlByt, &data)
	if err != nil {
		return nil, err
	}

	data = RecursiveCastMapStringInterfaceToMapStringInterface(data)

	val, ok := MapGet(data, keys...)
	if !ok {
		return nil, fmt.Errorf("invalid path")
	}

	transformed, err := transform(val)
	if err != nil {
		return nil, err
	}

	MapSet(data, transformed, keys...)

	return yaml.Marshal(data)
}

// MapGet takes in the map keys - each key goes one level deeper in the map
func MapGet(mp map[string]interface{}, key ...string) (interface{}, bool) {
	if mp == nil {
		return nil, false
	}

	if len(key) == 0 {
		return mp, true
	}

	if len(key) == 1 {
		val, ok := mp[key[0]]
		return val, ok
	}

	val, ok := mp[key[0]]
	if !ok {
		return mp, false
	}

	switch v := val.(type) {
	case map[string]interface{}:
		return MapGet(v, key[1:]...)
	case []interface{}:
		// Check if we can find key in the nested structure
		if len(key) < 2 {
			return mp, false
		}

		// Check if the key[1] is of type uint, if it is then
		keyNum, err := strconv.Atoi(key[1])
		if err != nil {
			return mp, false
		}

		if keyNum >= len(v) {
			return mp, false
		}

		valMapM, ok := v[keyNum].(map[string]interface{})
		if !ok {
			return mp, false
		}

		return MapGet(valMapM, key[2:]...)
	case []map[string]interface{}:
		// Check if we can find key in the nested structure
		if len(key) < 2 {
			return mp, false
		}

		// Check if the key[1] is of type uint, if it is then
		keyNum, err := strconv.Atoi(key[1])
		if err != nil {
			return mp, false
		}

		if keyNum >= len(v) {
			return mp, false
		}

		return MapGet(v[keyNum], key[2:]...)
	}

	return mp, true
}

// MapSet takes in the map that needs to be manipulated, the value that needs to
// be assgined to be assigned and the key - each key goes one level deeper in the map
func MapSet(mp map[string]interface{}, value interface{}, key ...string) {
	var _mapSet func(map[string]interface{}, interface{}, ...string) map[string]interface{}

	_mapSet = func(mp map[string]interface{}, value interface{}, key ...string) map[string]interface{} {
		if mp == nil {
			return nil
		}

		if len(key) == 0 {
			return mp
		}

		if len(key) == 1 {
			mp[key[0]] = value
			return mp
		}

		val, ok := mp[key[0]]
		if !ok {
			return mp
		}

		switch v := val.(type) {
		case map[string]interface{}:
			mp[key[0]] = _mapSet(v, value, key[1:]...)
			return mp
		case []interface{}:
			// Check if we can find key in the nested structure
			if len(key) < 2 {
				return mp
			}

			// Check if the key[1] is of type uint, if it is then
			keyNum, err := strconv.Atoi(key[1])
			if err != nil {
				return mp
			}

			if keyNum >= len(v) {
				return mp
			}

			valMapM, ok := v[keyNum].(map[string]interface{})
			if !ok {
				return mp
			}

			v[keyNum] = _mapSet(valMapM, value, key[2:]...)

			mp[key[0]] = v

			return mp
		case []map[string]interface{}:
			// Check if we can find key in the nested structure
			if len(key) < 2 {
				return mp
			}

			// Check if the key[1] is of type uint, if it is then
			keyNum, err := strconv.Atoi(key[1])
			if err != nil {
				return mp
			}

			if keyNum >= len(v) {
				return mp
			}

			v[keyNum] = _mapSet(v[keyNum], value, key[2:]...)

			mp[key[0]] = v

			return mp
		}

		return mp
	}

	_mapSet(mp, value, key...)
}

// RecursiveCastMapStringInterfaceToMapStringInterface will convert a
// map[string]interface{} recursively => map[string]interface{}
func RecursiveCastMapStringInterfaceToMapStringInterface(in map[string]interface{}) map[string]interface{} {
	res := ConvertMapInterfaceMapString(in)
	out, ok := res.(map[string]interface{})
	if !ok {
		fmt.Println("failed to cast")
	}

	return out
}

// ConvertMapInterfaceMapString converts map[interface{}]interface{} => map[string]interface{}
//
// It will also convert []interface{} => []string
func ConvertMapInterfaceMapString(v interface{}) interface{} {
	switch x := v.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v2 := range x {
			switch k2 := k.(type) {
			case string:
				m[k2] = ConvertMapInterfaceMapString(v2)
			default:
				m[fmt.Sprint(k)] = ConvertMapInterfaceMapString(v2)
			}
		}
		v = m

	case []interface{}:
		for i, v2 := range x {
			x[i] = ConvertMapInterfaceMapString(v2)
		}

	case map[string]interface{}:
		for k, v2 := range x {
			x[k] = ConvertMapInterfaceMapString(v2)
		}
	}

	return v
}

// SetOverrideValues returns the value overrides based on current context to install/upgrade helm chart
func SetOverrideValues(ctx *config.Context, mesheryImageVersion string) map[string]interface{} {
	// first initialize all the components' "enabled" field to false
	// this matches to the components listed in install/kubernetes/helm/meshery/values.yaml
	valueOverrides := map[string]interface{}{
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
	}

	// set the "enabled" field to true only for the components listed in the context
	for _, component := range ctx.GetComponents() {
		if _, ok := valueOverrides[component]; ok {
			valueOverrides[component] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	// set the meshery image version
	valueOverrides["image"] = map[string]interface{}{
		"tag": ctx.GetChannel() + "-" + mesheryImageVersion,
	}

	// set the provider
	if ctx.GetProvider() != "" {
		valueOverrides["env"] = map[string]interface{}{
			"PROVIDER": ctx.GetProvider(),
		}
	}

	// disable the operator
	if ctx.GetOperatorStatus() == "disabled" {
		if _, ok := valueOverrides["env"]; !ok {
			valueOverrides["env"] = map[string]interface{}{}
		}
		envOverrides := valueOverrides["env"].(map[string]interface{})
		envOverrides["DISABLE_OPERATOR"] = "'true'"
	}

	return valueOverrides
}

// CheckFileExists checks if the given file exists in system or not
func CheckFileExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("%s does not exist", name)
	}
	return false, errors.Wrap(err, fmt.Sprintf("Failed to read/fetch the file %s", name))
}
