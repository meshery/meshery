package utils

import (
	"bufio"
	"bytes"
	"context"
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

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
)

const (
	dockerComposeWebURL         = "https://api.github.com/repos/docker/compose/releases/latest"
	defaultDockerComposeVersion = "1.24.1/docker-compose"
	dockerComposeBinaryURL      = "https://github.com/docker/compose/releases/download/"
	dockerComposeBinary         = "/usr/local/bin/docker-compose"

	// Usage URLs
	docsBaseURL = "https://docs.meshery.io/"

	rootUsageURL   = docsBaseURL + "guides/mesheryctl/#global-commands-and-flags"
	perfUsageURL   = docsBaseURL + "guides/mesheryctl/#performance-management"
	systemUsageURL = docsBaseURL + "guides/mesheryctl/#meshery-lifecycle-management"
	meshUsageURL   = docsBaseURL + "guides/mesheryctl/#service-mesh-lifecycle-management"
)

const (

	// Repo Details
	mesheryGitHubOrg  string = "layer5io"
	mesheryGitHubRepo string = "meshery"
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
	// ServiceAccount is the name of a Kubernetes manifest file required to setup Meshery
	// check https://github.com/layer5io/meshery/tree/master/install/deployment_yamls/k8s
	ServiceAccount = "service-account.yaml"
	// ViperCompose is an instance of viper for docker-compose
	ViperCompose = viper.New()
	// SilentFlag skips waiting for user input and proceeds with default options
	SilentFlag bool
)

// ListOfAdapters returns the list of adapters available
var ListOfAdapters = []string{"meshery-istio", "meshery-linkerd", "meshery-consul", "meshery-nsm", "meshery-kuma", "meshery-cpx", "meshery-osm", "meshery-traefik-mesh"}

// TemplateContext is the template context provided when creating a config file
var TemplateContext = config.Context{
	Endpoint: "http://localhost:9081",
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

// GetMesheryGitHubOrg retrieves the name of the GitHub organization under which the Meshery repository resides.
func GetMesheryGitHubOrg() string {
	return mesheryGitHubOrg
}

// GetMesheryGitHubRepo retrieves the name of the Meshery repository
func GetMesheryGitHubRepo() string {
	return mesheryGitHubRepo
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

//PreReqCheck prerequisites check
func PreReqCheck(subcommand string, focusedContext string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}
	currCtx, err := mctlCfg.SetCurrentContext(focusedContext)
	if err != nil {
		return err
	}

	if currCtx.Platform == "docker" {
		//Check whether docker daemon is running or not
		if err := exec.Command("docker", "ps").Run(); err != nil {
			log.Info("Docker is not running.")
			//No auto installation of docker for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "Please start Docker. Run `mesheryctl system %s` once Docker is started.", subcommand)
			}
			err = startdockerdaemon(subcommand)
			if err != nil {
				return errors.Wrapf(err, "failed to start Docker.")
			}
		}
		//Check for installed docker-compose on client system
		if err := exec.Command("docker-compose", "-v").Run(); err != nil {
			log.Info("Docker-Compose is not installed")
			//No auto installation of Docker-compose for windows
			if runtime.GOOS == "windows" {
				return errors.Wrapf(err, "please install docker-compose. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
			}
			err = installprereq()
			if err != nil {
				return errors.Wrapf(err, "failed to install prerequisites. Run `mesheryctl system %s` after docker-compose is installed.", subcommand)
			}
		}
	} else if currCtx.Platform == "kubernetes" {
		client, err := meshkitkube.New([]byte(""))

		if err != nil {
			return errors.Wrapf(err, "failed to create new client")
		}

		podInterface := client.KubeClient.CoreV1().Pods("")
		_, err = podInterface.List(context.TODO(), v1.ListOptions{})

		if err != nil {
			log.Info("Kubernetes unreachable.")
			return errors.Wrap(err, "Kubernetes is not available. Verify Kubernetes is up, reachable, and a valid cert / token is available.")
		}
	} else {
		return errors.New(fmt.Sprintf("%v platform not supported", currCtx.Platform))
	}
	return nil
}

func startdockerdaemon(subcommand string) error {
	userResponse := false
	// read user input on whether to start Docker daemon or not.
	if SilentFlag {
		userResponse = true
	} else {
		userResponse = AskForConfirmation("Start Docker now")
	}
	if userResponse != true {
		return errors.Errorf("Please start Docker, then run the command `mesheryctl system %s`", subcommand)
	}

	log.Info("Attempting to start Docker...")
	// once user gaves permission, start docker daemon on linux/macOS
	if runtime.GOOS == "linux" {
		if err := exec.Command("sudo", "service", "docker", "start").Run(); err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
	} else {
		// Assuming we are on macOS, try to start Docker from default path
		cmd := exec.Command("/Applications/Docker.app/Contents/MacOS/Docker")
		err := cmd.Start()
		if err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
		// wait for few seconds for docker to start
		err = exec.Command("sleep", "30").Run()
		if err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
		// check whether docker started successfully or not, throw an error message otherwise
		if err := exec.Command("docker", "ps").Run(); err != nil {
			return errors.Wrapf(err, "please start Docker then run the command `mesheryctl system %s`", subcommand)
		}
	}
	log.Info("Prerequisite Docker started.")
	return nil
}

func installprereq() error {
	log.Info("Attempting Docker-Compose installation...")
	ostype, osarch, err := prereq()
	if err != nil {
		return errors.Wrap(err, "failed to get prerequisites")
	}

	osdetails := strings.TrimRight(string(ostype), "\r\n") + "-" + strings.TrimRight(string(osarch), "\r\n")

	dockerComposeBinaryURL := dockerComposeBinaryURL
	//checks for the latest docker-compose
	resp, err := http.Get(dockerComposeWebURL)
	if err != nil {
		dockerComposeBinaryURL = dockerComposeBinaryURL + defaultDockerComposeVersion
	} else {
		var dat map[string]interface{}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "failed to read response body")
		}
		if err := json.Unmarshal(body, &dat); err != nil {
			return errors.Wrap(err, "failed to unmarshal json into object")
		}
		num := dat["tag_name"]
		dockerComposeBinaryURL = fmt.Sprintf(dockerComposeBinaryURL+"%v/docker-compose", num)
	}
	dockerComposeBinaryURL = dockerComposeBinaryURL + "-" + osdetails
	if err := DownloadFile(dockerComposeBinary, dockerComposeBinaryURL); err != nil {
		return errors.Wrapf(err, "failed to download %s from %s", dockerComposeBinary, dockerComposeBinaryURL)
	}
	if err := exec.Command("chmod", "+x", dockerComposeBinary).Run(); err != nil {
		return errors.Wrap(err, "failed to execute command")
	}
	log.Info("Prerequisite Docker Compose is installed.")
	return nil
}

// IsMesheryRunning checks if the meshery server containers are up and running
func IsMesheryRunning(currPlatform string) (bool, error) {
	switch currPlatform {
	case "docker":
		{
			op, err := exec.Command("docker-compose", "-f", DockerComposeFile, "ps").Output()
			if err != nil {
				return false, err
			}
			return strings.Contains(string(op), "meshery"), nil
		}
	case "kubernetes":
		{
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				return false, errors.Wrap(err, "failed to create new client")
			}

			podInterface := client.KubeClient.CoreV1().Pods(MesheryNamespace)
			_, err = podInterface.List(context.TODO(), v1.ListOptions{})

			if err != nil {
				return false, err
			}

			return true, err
		}
	}

	return false, nil
}

// AddAuthDetails Adds authentication cookies to the request
func AddAuthDetails(req *http.Request, filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		err = errors.Wrap(err, "could not read token:")
		return err
	}
	var tokenObj map[string]string
	if err := json.Unmarshal(file, &tokenObj); err != nil {
		err = errors.Wrap(err, "token file invalid:")
		return err
	}
	req.AddCookie(&http.Cookie{
		Name:     tokenName,
		Value:    tokenObj[tokenName],
		HttpOnly: true,
	})
	req.AddCookie(&http.Cookie{
		Name:     providerName,
		Value:    tokenObj[providerName],
		HttpOnly: true,
	})
	return nil
}

// UpdateAuthDetails checks gets the token (old/refreshed) from meshery server and writes it back to the config file
func UpdateAuthDetails(filepath string) error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	// TODO: get this from the global config
	req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/gettoken", bytes.NewBuffer([]byte("")))
	if err != nil {
		err = errors.Wrap(err, "error Creating the request :")
		return err
	}
	if err := AddAuthDetails(req, filepath); err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	defer SafeClose(resp.Body)

	if err != nil {
		err = errors.Wrap(err, "error dispatching there request :")
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "error reading body :")
		return err
	}

	if ContentTypeIsHTML(resp) {
		return errors.New("invalid body")
	}

	return ioutil.WriteFile(filepath, data, os.ModePerm)
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

// RootError returns a formatted error message with a link to 'root' command usage page at
// in addition to the error message
func RootError(msg string) string {
	return formatError(msg, cmdRoot)
}

// PerfError returns a formatted error message with a link to 'perf' command usage page at
// in addition to the error message
func PerfError(msg string) string {
	return formatError(msg, cmdPerf)
}

// SystemError returns a formatted error message with a link to 'system' command usage page
// in addition to the error message
func SystemError(msg string) string {
	return formatError(msg, cmdSystem)
}

// MeshError returns a formatted error message with a link to 'mesh' command usage page in addition to the error message
//func MeshError(msg string) string {
//	return formatError(msg, cmdMesh)
//}

// formatError returns a formatted error message with a link to the meshery command URL
func formatError(msg string, cmd cmdType) string {
	switch cmd {
	case cmdRoot:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, rootUsageURL)
	case cmdPerf:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, perfUsageURL)
	case cmdMesh:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, meshUsageURL)
	case cmdSystem:
		return fmt.Sprintf("%s\nSee %s for usage details\n", msg, systemUsageURL)
	}
	return fmt.Sprintf("%s\n", msg)
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

// AddTokenToConfig adds token passed to it to mesheryctl config file
func AddTokenToConfig(token config.Token, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	if mctlCfg.Tokens == nil {
		mctlCfg.Tokens = []config.Token{}
	}

	for i := range mctlCfg.Tokens {
		if mctlCfg.Tokens[i].Name == token.Name {
			return errors.New("error adding token: a token with same name already exists")
		}
	}

	mctlCfg.Tokens = append(mctlCfg.Tokens, token)

	viper.Set("contexts", mctlCfg.Contexts)
	viper.Set("current-context", mctlCfg.CurrentContext)
	viper.Set("tokens", mctlCfg.Tokens)

	err = viper.WriteConfig()
	if err != nil {
		return err
	}

	return nil
}

// AddContextToConfig adds context passed to it to mesheryctl config file
func AddContextToConfig(contextName string, context config.Context, configPath string, set bool) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return err
	}

	viper.SetConfigFile(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return errors.Wrap(err, "error processing config")
	}

	if mctlCfg.Contexts == nil {
		mctlCfg.Contexts = map[string]config.Context{}
	}

	_, exists := mctlCfg.Contexts[contextName]
	if exists {
		return errors.New("error adding context: a context with same name already exists")
	}

	mctlCfg.Contexts[contextName] = context
	if set {
		mctlCfg.CurrentContext = contextName
	}

	viper.Set("contexts", mctlCfg.Contexts)
	viper.Set("current-context", mctlCfg.CurrentContext)
	viper.Set("tokens", mctlCfg.Tokens)

	err = viper.WriteConfig()
	if err != nil {
		return err
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
