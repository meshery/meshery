package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/constants"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TestHelper struct {
	Version string
	BaseURL string
}

type MockURL struct {
	// method such as GET or POST
	Method string
	// url to mock the request
	URL string
	// response for the request
	Response string
	// response code
	ResponseCode int
}

func NewTestHelper(_ *testing.T) *TestHelper {
	return &TestHelper{
		Version: "v0.5.10",
		BaseURL: MesheryEndpoint,
	}
}

type CmdTestInput struct {
	Name                 string
	Args                 []string
	ExpectedResponse     string
	ExpectedResponseYaml string
	ExpectError          bool
	ErrorStringContains  []string
}

type GoldenFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewGoldenFile(t *testing.T, name string, directory string) *GoldenFile {
	return &GoldenFile{t: t, name: name, dir: directory}
}

// equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

// Path to the current file
func GetBasePath(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	return filepath.Dir(filename)
}

// Load a Golden file
func (tf *GoldenFile) Load() string {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := os.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}
	// ensuring that the newline characters in the content are consistent and match the expected newline representation
	normalizedContent := strings.ReplaceAll(string(content), "\r\n", "\n")
	return normalizedContent
}

// Load a Golden file
func (tf *GoldenFile) LoadByte() []byte {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := os.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return content
}

// write a Golden file
func (tf *GoldenFile) Write(content string) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.WriteFile(path, []byte(content), 0755)
			if err != nil {
				fmt.Printf("Unable to write file: %v", err)
			}
			return
		}
		tf.t.Fatal(err)
	}

	err = os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

// write a Golden file
func (tf *GoldenFile) WriteInByte(content []byte) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

// use default context /pkg/utils/TestConfig.yaml
func SetupContextEnv(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Error("unable to locate meshery directory")
	}
	viper.Reset()
	viper.SetConfigFile(path + "/../../../../pkg/utils/TestConfig.yaml")
	DefaultConfigPath = path + "/../../../../pkg/utils/TestConfig.yaml"
	//fmt.Println(viper.ConfigFileUsed())
	err = viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// setup logrus formatter and return the buffer in which commands output is to be set.
func SetupLogrusGrabTesting(_ *testing.T, _ bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	logrus.SetOutput(b)
	SetupLogrusFormatter()
	return b
}

// setup meshkit logger for testing and return the buffer in which commands output is to be set.
func SetupMeshkitLoggerTesting(_ *testing.T, verbose bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	Log = SetupMeshkitLogger("mesheryctl", verbose, b)
	return b
}

// setup custom context with SetupCustomContextEnv
func SetupCustomContextEnv(t *testing.T, pathToContext string) {
	viper.Reset()
	ViperCompose = viper.New()
	ViperMeshconfig = viper.New()

	viper.SetConfigFile(pathToContext)
	DefaultConfigPath = pathToContext
	//fmt.Println(viper.ConfigFileUsed())
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// Start mock HTTP client to mock requests
func StartMockery(t *testing.T) {
	// activate http mocking
	httpmock.Activate()

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := NewGoldenFile(t, "validate.version.github.golden", fixturesDir).Load()

	// For validate version requests
	url1 := "https://github.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/releases/tag/" + "v0.5.54"
	httpmock.RegisterResponder("GET", url1,
		httpmock.NewStringResponder(200, apiResponse))
}

// stop HTTP mock client
func StopMockery(_ *testing.T) {
	httpmock.DeactivateAndReset()
}

// Set file location for testing stuff
func SetFileLocationTesting(dir string) {
	MesheryFolder = filepath.Join(dir, "fixtures", MesheryFolder)
	DockerComposeFile = filepath.Join(MesheryFolder, DockerComposeFile)
	AuthConfigFile = filepath.Join(MesheryFolder, AuthConfigFile)
}

func Populate(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

func StartMockMesheryServer(t *testing.T) error {
	serverAddr := strings.TrimPrefix(MesheryEndpoint, "http://")
	l, err := net.Listen("tcp", serverAddr)
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				// Log error and break the loop if it's a permanent error
				if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
					t.Logf("Failed to accept connection: %v", err)
					break
				}
				continue
			}
			// Close the connection to verify IsServerRunning() in auth.go
			conn.Close()
		}
	}()

	// Give the server some time to start
	time.Sleep(100 * time.Millisecond)
	return nil
}

// The HandlePagination function add special characters that is not
// handle properly in test. This function will remove undesired characters
// and spaces to ensure excepted versus actual result match when using http.MockURL
func CleanStringFromHandlePagination(data string) string {
	cleaned := stripAnsiEscapeCodes(data)
	cleaned = formatToTabs(cleaned)
	return cleaned
}

// removeANSICodes removes ANSI escape codes from a string.
//
// Parameters:
//
//	text - The input string that may contain ANSI escape sequences.
//
// Returns:
//
//	A string with the ANSI escape codes removed.
func stripAnsiEscapeCodes(text string) string {
	ansi := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return ansi.ReplaceAllString(text, "")
}

// formatToTabs replaces multiple spaces with tabs and trims spaces
//
// Parameters:
//
//	s - The input string containing columns separated by multiple spaces.
//
// Returns:
//
//	A string where multiple spaces are replaced with a single tab between columns, and leading/trailing spaces are removed.
func formatToTabs(data string) string {
	s := strings.TrimSpace(data)

	// Replace multiple spaces with a single tab
	re := regexp.MustCompile(`\s{2,}`) // Match 2 or more spaces
	s = re.ReplaceAllString(s, "\t")

	return s
}

type MesheryListCommamdTest struct {
	Name             string
	Args             []string
	URL              string
	Fixture          string
	ExpectedResponse string
	ExpectError      bool
}

func InvokeMesheryctlTestListCommand(t *testing.T, updateGoldenFile *bool, cmd *cobra.Command, tests []MesheryListCommamdTest, commandDir string, commadName string) {
	// setup current context
	SetupContextEnv(t)

	//initialize mock server for handling requests
	StartMockery(t)

	// create a test helper
	testContext := NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	apiToken := filepath.Join(currDir, "fixtures", "token.golden")
	fixturesDir := filepath.Join(commandDir, "fixtures")

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			apiResponse := NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			TokenFlag = apiToken

			httpmock.RegisterResponder("GET", testContext.BaseURL+tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(commandDir, "testdata")
			golden := NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Grab console prints
			rescueStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			_ = SetupMeshkitLoggerTesting(t, false)
			cmd.SetArgs(tt.Args)
			cmd.SetOut(rescueStdout)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// write it in file
					if *updateGoldenFile {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()

					Equals(t, expectedResponse, err.Error())
					return
				}
				t.Fatal(err)
			}

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = rescueStdout

			actualResponse := string(out)

			if *updateGoldenFile {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := CleanStringFromHandlePagination(expectedResponse)

			Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Logf("List %s test", commadName)
	}

	StopMockery(t)
}
