package utils

import (
	"bufio"
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
	docsBaseURL    = "https://docs.meshery.io/"
	rootUsageURL   = docsBaseURL + "reference/mesheryctl"
	perfUsageURL   = docsBaseURL + "reference/mesheryctl/perf"
	systemUsageURL = docsBaseURL + "reference/mesheryctl/system"
	meshUsageURL   = docsBaseURL + "reference/mesheryctl/mesh"

	// Meshery Server Location
	EndpointProtocol = "http"
)

type cmdType string

const (
	cmdRoot   cmdType = "root"
	cmdPerf   cmdType = "perf"
	cmdMesh   cmdType = "mesh"
	cmdSystem cmdType = "system"
)

var (
	// ResetFlag indicates if a reset is required
	ResetFlag bool
	// SkipResetFlag indicates if fetching the updated manifest files is required
	SkipResetFlag bool
	// MesheryEndpoint is the default URL in which Meshery is exposed
	MesheryEndpoint = "http://localhost:9081"
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
	// ViperCompose is an instance of viper for docker-compose
	ViperCompose = viper.New()
	// ViperDocker is an instance of viper for the meshconfig file when the platform is docker
	ViperDocker = viper.New()
	// ViperK8s is an instance of viper for the meshconfig file when the platform is kubernetes
	ViperK8s = viper.New()
	// SilentFlag skips waiting for user input and proceeds with default options
	SilentFlag bool
	// PlatformFlag sets the platform for the initial config file
	PlatformFlag string
)

// ListOfAdapters returns the list of adapters available
var ListOfAdapters = []string{"meshery-istio", "meshery-linkerd", "meshery-consul", "meshery-nsm", "meshery-kuma", "meshery-cpx", "meshery-osm", "meshery-traefik-mesh"}

// TemplateContext is the template context provided when creating a config file
var TemplateContext = config.Context{
	Endpoint: EndpointProtocol + "://localhost:9081",
	Token:    "Default",
	Platform: "docker",
	Adapters: ListOfAdapters,
	Channel:  "stable",
	Version:  "latest",
}

// TemplateToken is the template token provided when creating a config file
var TemplateToken = config.Token{
	Name:     "Default",
	Location: AuthConfigFile,
}

type cryptoSource struct{}

func (s cryptoSource) Seed(seed int64) {}

// Int63 to generate high security rand through crypto
func (s cryptoSource) Int63() int64 {
	return int64(s.Uint64() & ^uint64(1<<63))
}

func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		log.Fatal(err)
	}
	return v
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

// TODO: Use the same DownloadFile function from MeshKit instead of the function below
// and change all it's occurrences

// DownloadFile from url and save to configured file location
func DownloadFile(filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrapf(err, "failed to make GET request to %s", url)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", filepath)
	}
	defer func() {
		_ = out.Close()
	}()
	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
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

// NavigateToBroswer naviagtes to the endpoint displaying Meshery UI in the broswer, based on the host operating system.
func NavigateToBrowser(endpoint string) error {
	//check for os of host machine
	if runtime.GOOS == "windows" {
		// Meshery running on Windows host
		err := exec.Command("rundll32", "url.dll,FileProtocolHandler", endpoint).Start()
		if err != nil {
			return errors.Wrap(err, SystemError("failed to exec command"))
		}
	} else if runtime.GOOS == "linux" {
		// Meshery running on Linux host
		_, err := exec.LookPath("xdg-open")
		if err != nil {
			return errors.Wrap(err, SystemError("failed to exec command"))
			//find out what to do here!
		}
		err = exec.Command("xdg-open", endpoint).Start()
		if err != nil {
			return errors.Wrap(err, SystemError("failed to exec command"))
		}
	} else {
		// Assume Meshery running on MacOS host
		err := exec.Command("open", endpoint).Start()
		if err != nil {
			return errors.Wrap(err, SystemError("failed to exec command"))
		}
	}

	return nil
}

// UploadFileWithParams returns a request configured to upload files with other values
func UploadFileWithParams(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
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

	request, err := http.NewRequest("POST", uri, body)
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

// ReadToken returns a map of the token passed in
func ReadToken(filepath string) (map[string]string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		err = errors.Wrap(err, "could not read token:")
		return nil, err
	}
	var tokenObj map[string]string
	if err := json.Unmarshal(file, &tokenObj); err != nil {
		err = errors.Wrap(err, "token file invalid:")
		return nil, err
	}
	return tokenObj, nil
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

func CreateDefaultSpinner(suffix string, finalMsg string) *spinner.Spinner {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	s.Suffix = " " + suffix
	s.FinalMSG = finalMsg + "\n"
	return s
}
