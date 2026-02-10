package utils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
	"github.com/meshery/meshkit/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
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
	ExpectedError        error `default:"nil"`
	IsOutputGolden       bool  `default:"true"`
}

type GoldenFile struct {
	t    *testing.T
	name string
	dir  string
}

func NewGoldenFile(t *testing.T, name string, directory string) *GoldenFile {
	return &GoldenFile{t: t, name: name, dir: directory}
}

func InitTestEnvironment(t *testing.T) *TestHelper {
	SetupContextEnv(t)
	StartMockery(t)
	testContext := NewTestHelper(t)
	return testContext
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
	defer func() { _ = source.Close() }()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() { _ = destination.Close() }()
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
			_ = conn.Close()
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
	cleaned := StripAnsiEscapeCodes(data)
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
func StripAnsiEscapeCodes(text string) string {
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

// AssertMeshkitErrorsEqual compares  relevant fields of two meshkit errors
func AssertMeshkitErrorsEqual(t *testing.T, got, expected error) {
	t.Helper()
	assert.Equal(t, reflect.TypeOf(got), reflect.TypeOf(expected), "error type mismatch")
	assert.Equal(t, errors.GetCode(got), errors.GetCode(expected), "error code mismatch")
	assert.Equal(t, errors.GetLDescription(got), errors.GetLDescription(expected), "long description mismatch")
}

type MesheryListCommandTest struct {
	Name             string
	Args             []string
	URL              string
	Fixture          string
	ExpectedResponse string
	ExpectError      bool
	ExpectedError    error `default:"nil"`
	IsOutputGolden   bool  `default:"true"`
}

func GetToken(t *testing.T) string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	return filepath.Join(currDir, "fixtures", "token.golden")
}

func InvokeMesheryctlTestListCommand(t *testing.T, updateGoldenFile *bool, cmd *cobra.Command, tests []MesheryListCommandTest, commandDir string, commandName string) {
	testContext := InitTestEnvironment(t)

	fixturesDir := filepath.Join(commandDir, "fixtures")

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer ResetCommandFlags(cmd, t)

			apiResponse := NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

			TokenFlag = GetToken(t)

			httpmock.RegisterResponder("GET", testContext.BaseURL+tt.URL,
				httpmock.NewStringResponder(200, apiResponse))

			testdataDir := filepath.Join(commandDir, "testdata")
			golden := NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			var buf bytes.Buffer

			// Properly save and restore stdout using defer
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Ensure stdout is always restored
			defer func() {
				os.Stdout = originalStdout
			}()

			_ = SetupMeshkitLoggerTesting(t, false)

			cmd.SetArgs(tt.Args)
			cmd.SetOut(originalStdout)
			err := cmd.Execute()

			// Close write end before reading
			_ = w.Close()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// Keep this check to see if output is golden file during transition
					if tt.IsOutputGolden {
						// write it in file
						if *updateGoldenFile {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						Equals(t, expectedResponse, err.Error())
						return
					}
					AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				t.Fatal(err)
			}

			_, errCopy := io.Copy(&buf, r)
			if errCopy != nil {
				t.Fatal(errCopy)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			actualResponse := buf.String()

			if *updateGoldenFile {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := CleanStringFromHandlePagination(expectedResponse)

			Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Logf("List %s test", commandName)
	}

	StopMockery(t)
}

type MesheryCommandTest struct {
	Name             string
	Args             []string
	HttpMethod       string
	HttpStatusCode   int
	URL              string
	Fixture          string
	ExpectedResponse string
	ExpectError      bool
	IsOutputGolden   bool  `default:"true"`
	ExpectedError    error `default:"nil"`
}

func InvokeMesheryctlTestCommand(t *testing.T, updateGoldenFile *bool, cmd *cobra.Command, tests []MesheryCommandTest, commandDir string, commandName string) {
	testContext := InitTestEnvironment(t)

	fixturesDir := filepath.Join(commandDir, "fixtures")

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer ResetCommandFlags(cmd, t)

			if tt.Fixture != "" {
				apiResponse := NewGoldenFile(t, tt.Fixture, fixturesDir).Load()

				TokenFlag = GetToken(t)

				url := testContext.BaseURL + tt.URL
				httpMethod := tt.HttpMethod

				if tt.HttpStatusCode < 0 {
					httpmock.RegisterResponder(httpMethod, url,
						func(req *http.Request) (*http.Response, error) {
							return nil, &net.OpError{Op: "dial", Net: "tcp", Addr: nil, Err: net.ErrClosed}
						})
				} else {
					httpmock.RegisterResponder(httpMethod, url,
						httpmock.NewStringResponder(tt.HttpStatusCode, apiResponse))
				}

			}

			testdataDir := filepath.Join(commandDir, "testdata")
			golden := NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			originalStdout := os.Stdout
			b := SetupMeshkitLoggerTesting(t, false)
			defer func() {
				os.Stdout = originalStdout
			}()
			cmd.SetArgs(tt.Args)
			cmd.SetOut(b)
			err := cmd.Execute()
			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// Keep this check to see if output is golden file during transition
					if tt.IsOutputGolden {

						// write it in file
						if *updateGoldenFile {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						Equals(t, expectedResponse, err.Error())
						return
					}
					AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return

				}
				t.Fatal(err)

			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			actualResponse := b.String()

			if *updateGoldenFile {
				golden.Write(actualResponse)
			}

			expectedResponse := golden.Load()

			cleanedActualResponse := CleanStringFromHandlePagination(actualResponse)
			cleanedExpectedResponse := CleanStringFromHandlePagination(expectedResponse)

			Equals(t, cleanedExpectedResponse, cleanedActualResponse)
		})
		t.Logf("Test '%s' executed", tt.Name)
	}
	StopMockery(t)
}

type MesheryMultiURLCommamdTest struct {
	Name             string
	Args             []string
	URLs             []MockURL
	ExpectedResponse string
	Token            string
	ExpectError      bool
	IsOutputGolden   bool  `default:"true"`
	ExpectedError    error `default:"nil"`
}

func RunMesheryctlMultiURLTests(t *testing.T, updateGoldenFile *bool, cmd *cobra.Command, tests []MesheryMultiURLCommamdTest, commandDir string, commandName string, resetVariables func()) {
	_ = InitTestEnvironment(t)

	fixturesDir := filepath.Join(commandDir, "fixtures")

	// Run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer resetVariables()
			defer ResetCommandFlags(cmd, t)

			if tt.Token != "" {
				TokenFlag = tt.Token
			} else {
				TokenFlag = GetToken(t)
			}

			for _, mock := range tt.URLs {
				apiResponse := NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			testdataDir := filepath.Join(commandDir, "testdata")
			golden := NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			// Properly save and restore stdout using defer
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Ensure stdout is always restored
			defer func() {
				os.Stdout = originalStdout
			}()

			Log = SetupMeshkitLogger("mesheryctl", false, w)

			cmd.SetArgs(tt.Args)
			cmd.SetOut(w)
			err := cmd.Execute()

			// Close write end before reading
			_ = w.Close()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// Keep this check to see if output is golden file during transition
					if tt.IsOutputGolden {
						// write it in file
						if *updateGoldenFile {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						Equals(t, expectedResponse, err.Error())
						return
					}
					AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				// Unexpected error - fail immediately
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			var buf bytes.Buffer
			_, errCopy := io.Copy(&buf, r)
			if errCopy != nil {
				t.Fatal(errCopy)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			actualResponse := buf.String()

			if *updateGoldenFile {
				golden.Write(actualResponse)
			}

			expectedResponse := golden.Load()

			cleanedActualResponse := CleanStringFromHandlePagination(actualResponse)
			cleanedExpectedResponse := CleanStringFromHandlePagination(expectedResponse)

			Equals(t, cleanedExpectedResponse, cleanedActualResponse)
		})
		t.Logf("Test '%s' executed", tt.Name)
	}
	StopMockery(t)
}

func RunMesheryctlMultipleURLsListTests(t *testing.T, updateGoldenFile *bool, cmd *cobra.Command, tests []MesheryMultiURLCommamdTest, commandDir string, commandName string, resetVariables func()) {
	_ = InitTestEnvironment(t)

	fixturesDir := filepath.Join(commandDir, "fixtures")

	// run tests
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			defer resetVariables()
			defer ResetCommandFlags(cmd, t)

			if tt.Token != "" {
				TokenFlag = tt.Token
			} else {
				TokenFlag = GetToken(t)
			}

			for _, mock := range tt.URLs {
				apiResponse := NewGoldenFile(t, mock.Response, fixturesDir).Load()
				httpmock.RegisterResponder(mock.Method, mock.URL,
					httpmock.NewStringResponder(mock.ResponseCode, apiResponse))
			}

			testdataDir := filepath.Join(commandDir, "testdata")
			golden := NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			var buf bytes.Buffer

			// Properly save and restore stdout using defer
			originalStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Ensure stdout is always restored
			defer func() {
				os.Stdout = originalStdout
			}()

			_ = SetupMeshkitLoggerTesting(t, false)

			cmd.SetArgs(tt.Args)
			cmd.SetOut(originalStdout)
			err := cmd.Execute()

			// Close write end before reading
			w.Close()

			if err != nil {
				// if we're supposed to get an error
				if tt.ExpectError {
					// Keep this check to see if output is golden file during transition
					if tt.IsOutputGolden {
						// write it in file
						if *updateGoldenFile {
							golden.Write(err.Error())
						}
						expectedResponse := golden.Load()

						Equals(t, expectedResponse, err.Error())
						return
					}
					AssertMeshkitErrorsEqual(t, err, tt.ExpectedError)
					return
				}
				// Unexpected error - fail immediately
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.ExpectError {
				t.Fatalf("expected an error but command succeeded")
			}

			_, errCopy := io.Copy(&buf, r)
			if errCopy != nil {
				t.Fatal(errCopy)
			}

			actualResponse := buf.String()

			if *updateGoldenFile {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()

			cleanedActualResponse := CleanStringFromHandlePagination(actualResponse)
			cleanedExceptedResponse := CleanStringFromHandlePagination(expectedResponse)

			Equals(t, cleanedExceptedResponse, cleanedActualResponse)
		})
		t.Logf("List %s test", commandName)
	}

	StopMockery(t)
}

func ResetCommandFlags(c *cobra.Command, t *testing.T) {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			t.Fatalf("failed to reset flag %q: %v", f.Name, err)
		}
	})
	for _, sub := range c.Commands() {
		ResetCommandFlags(sub, t)
	}
}
