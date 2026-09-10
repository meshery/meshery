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
	mesheryctllogger "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/logger"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
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
	ExpectedError        error
	IsOutputGolden       bool
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
	viper.Set("LOG_LEVEL", int(logrus.InfoLevel))
	testContext := NewTestHelper(t)
	return testContext
}

// Equals fails the test unless exp and act are deeply equal, reporting the
// difference through tb.
//
// Never report through os.Stdout. The golden-file helpers below swap os.Stdout
// for a pipe to capture command output and close the write end before
// comparing, so anything printed to os.Stdout from here goes to a closed pipe
// and is lost. `fmt.Printf` also detaches the message from the failing test
// under `go test -json`, which is what CI runs. Either one alone made every
// golden failure undiagnosable - a bare FAIL with no expected/got.
//
// Formatting matters as much as the destination: `%#v` renders a multi-line
// golden as a single escaped Go-syntax blob, so even a message that surfaced
// said only that something differed, never what. Strings are therefore diffed
// rather than dumped; see stringDiff.
func Equals(tb testing.TB, exp, act interface{}) {
	tb.Helper()
	if reflect.DeepEqual(exp, act) {
		return
	}

	expStr, expOK := exp.(string)
	actStr, actOK := act.(string)
	if expOK && actOK {
		tb.Fatalf("golden mismatch\n%s", stringDiff(expStr, actStr))
	}

	tb.Fatalf("not equal\n\texpected: %#v\n\tactual:   %#v", exp, act)
}

// stringDiff renders a line-oriented difference between two strings, naming the
// first line that differs and making otherwise-invisible causes visible:
// trailing whitespace, a missing or extra trailing newline, and CR/LF drift are
// the usual reasons a golden file "looks identical" but is not.
func stringDiff(exp, act string) string {
	var b strings.Builder

	expLines := strings.Split(exp, "\n")
	actLines := strings.Split(act, "\n")

	// Compare presence as well as content. For "a" versus "a\n", Split yields
	// ["a"] and ["a", ""], so comparing only the (defaulted) empty strings at
	// index 1 finds no difference and the diagnostic silently omits the very
	// thing that differs - a missing or extra final newline.
	firstDiff := -1
	for i := 0; i < len(expLines) || i < len(actLines); i++ {
		haveE, haveA := i < len(expLines), i < len(actLines)
		var e, a string
		if haveE {
			e = expLines[i]
		}
		if haveA {
			a = actLines[i]
		}
		if haveE != haveA || e != a {
			firstDiff = i
			break
		}
	}

	if firstDiff >= 0 {
		fmt.Fprintf(&b, "\tfirst difference at line %d:\n", firstDiff+1)
		var e, a string
		haveE, haveA := firstDiff < len(expLines), firstDiff < len(actLines)
		if haveE {
			e = expLines[firstDiff]
		}
		if haveA {
			a = actLines[firstDiff]
		}
		// A single-line golden (rendered JSON, for instance) can be many
		// kilobytes wide. Dumping both in full names the line but not the
		// difference, so window around the first differing column instead.
		if haveE && haveA && (len(e) > diffLineWidth || len(a) > diffLineWidth) {
			col := 0
			for col < len(e) && col < len(a) && e[col] == a[col] {
				col++
			}
			fmt.Fprintf(&b, "\t  first differing column %d of %d/%d\n", col+1, len(e), len(a))
			fmt.Fprintf(&b, "\t  expected: %s\n", window(e, col))
			fmt.Fprintf(&b, "\t  actual:   %s\n", window(a, col))
		} else {
			if haveE {
				fmt.Fprintf(&b, "\t  expected: %s\n", visible(e))
			} else {
				fmt.Fprintf(&b, "\t  expected: <no line %d; expected ends at line %d>\n", firstDiff+1, len(expLines))
			}
			if haveA {
				fmt.Fprintf(&b, "\t  actual:   %s\n", visible(a))
			} else {
				fmt.Fprintf(&b, "\t  actual:   <no line %d; actual ends at line %d>\n", firstDiff+1, len(actLines))
			}
		}
	}

	fmt.Fprintf(&b, "\tline counts: expected %d, actual %d\n", len(expLines), len(actLines))
	fmt.Fprintf(&b, "\tbyte counts: expected %d, actual %d\n", len(exp), len(act))

	b.WriteString("\n\t--- expected ---\n")
	writeNumbered(&b, expLines)
	b.WriteString("\n\t--- actual ---\n")
	writeNumbered(&b, actLines)

	if strings.TrimSpace(exp) == strings.TrimSpace(act) {
		b.WriteString("\n\tNOTE: the two differ only in leading/trailing whitespace.\n")
	}
	return b.String()
}

// diffLineWidth is the point past which a differing line is windowed around
// the first differing column rather than printed whole.
const diffLineWidth = 200

// window renders a slice of s centred on col, with a caret marking the column.
func window(s string, col int) string {
	const pad = 60
	start := col - pad
	if start < 0 {
		start = 0
	}
	end := col + pad
	if end > len(s) {
		end = len(s)
	}
	prefix, suffix := "", ""
	if start > 0 {
		prefix = "..."
	}
	if end < len(s) {
		suffix = "..."
	}
	seg := visible(s[start:end])
	return fmt.Sprintf("%s%s%s   [at column %d]", prefix, seg, suffix, col+1)
}

func writeNumbered(b *strings.Builder, lines []string) {
	for i, l := range lines {
		fmt.Fprintf(b, "\t%4d | %s\n", i+1, visible(l))
	}
}

// visible marks up characters that are otherwise impossible to see in a diff.
func visible(s string) string {
	r := strings.NewReplacer("\r", "␍", "\t", "→")
	out := r.Replace(s)
	if out != strings.TrimRight(out, " ") {
		out = strings.TrimRight(out, " ") + strings.Repeat("·", len(out)-len(strings.TrimRight(out, " ")))
	}
	return out
}

// GetBasePath returns the directory containing this file
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

// LoadByte loads a golden file and returns its raw bytes
func (tf *GoldenFile) LoadByte() []byte {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	content, err := os.ReadFile(path)
	if err != nil {
		tf.t.Fatalf("could not read file %s: %v", tf.name, err)
	}

	return content
}

// Write writes content to the golden file, creating it if it does not already exist
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

// WriteInByte writes raw bytes to the golden file
func (tf *GoldenFile) WriteInByte(content []byte) {
	tf.t.Helper()
	path := filepath.Join(tf.dir, tf.name)
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		tf.t.Fatalf("could not write %s: %v", tf.name, err)
	}
}

// SharedTestConfigPath resolves pkg/utils/TestConfig.yaml from the working
// directory of a test package four levels below the mesheryctl module root,
// which is where every caller of SetupContextEnv lives.
func SharedTestConfigPath(t *testing.T) string {
	t.Helper()
	path, err := os.Getwd()
	if err != nil {
		t.Fatalf("unable to locate meshery directory: %v", err)
	}
	return filepath.Join(path, "..", "..", "..", "..", "pkg", "utils", "TestConfig.yaml")
}

// CopyMeshconfigFixture copies a meshconfig fixture into a directory private to
// t and returns the copy's path.
//
// pkg/utils/TestConfig.yaml is loaded by every mesheryctl test package, and
// `go test ./mesheryctl/...` runs those packages as concurrent processes. It
// therefore has to stay read-only on disk: a command under test persists the
// active meshconfig through viper.WriteConfig, which truncates the file before
// it rewrites it, and a sibling package reading inside that window parses an
// empty document with no error at all. The resulting config has an empty
// current-context, and the first request made with it reaches
// MesheryCtlConfig.GetBaseMesheryURL, which calls Log.Fatal and exits the whole
// test binary - a package-level FAIL with no failing test named. Pointing each
// test at its own copy keeps writes off the shared fixture.
func CopyMeshconfigFixture(t *testing.T, src string) string {
	t.Helper()
	dst := filepath.Join(t.TempDir(), "config.yaml")
	if err := Populate(src, dst); err != nil {
		t.Fatalf("unable to copy meshconfig fixture %v: %v", src, err)
	}
	return dst
}

// SetupContextEnv sets up the test context using the default pkg/utils/TestConfig.yaml configuration
func SetupContextEnv(t *testing.T) {
	configPath := CopyMeshconfigFixture(t, SharedTestConfigPath(t))
	viper.Reset()
	viper.SetConfigFile(configPath)
	DefaultConfigPath = configPath
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// SetupLogrusGrabTesting sets up the logrus formatter and returns the buffer that command output is written to.
func SetupLogrusGrabTesting(_ *testing.T, _ bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	logrus.SetOutput(b)
	SetupLogrusFormatter()
	return b
}

// SetupMeshkitLoggerTesting sets up the Meshkit logger for testing and returns the buffer that command output is written to.
func SetupMeshkitLoggerTesting(_ *testing.T, verbose bool) *bytes.Buffer {
	b := bytes.NewBufferString("")
	logLevel := logrus.InfoLevel
	if verbose {
		logLevel = logrus.DebugLevel
	}
	logger := mesheryctllogger.GetMeshkitLogger(logLevel)
	logger.SetLevel(logLevel)
	logger.UpdateLogOutput(b)
	Log = logger
	return b
}

// SetupCustomContextEnv sets up the test context using the configuration file at pathToContext
func SetupCustomContextEnv(t *testing.T, pathToContext string) {
	viper.Reset()
	ViperCompose = viper.New()
	ViperMeshconfig = viper.New()

	viper.SetConfigFile(pathToContext)
	DefaultConfigPath = pathToContext
	CfgFile = pathToContext
	err := viper.ReadInConfig()
	if err != nil {
		t.Errorf("unable to read configuration from %v, %v", viper.ConfigFileUsed(), err.Error())
	}

	_, err = config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Error("error processing config", err)
	}
}

// StartMockery activates HTTP mocking so requests made during the test are intercepted
func StartMockery(t *testing.T) {
	// activate http mocking
	httpmock.Activate()
}

// StopMockery deactivates HTTP mocking and resets it
func StopMockery(_ *testing.T) {
	httpmock.DeactivateAndReset()
}

// SetFileLocationTesting points MesheryFolder, DockerComposeFile, and AuthConfigFile at fixtures under dir, for use in tests
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

// CleanStringFromHandlePagination removes undesired characters and spaces added by the
// HandlePagination function so that expected and actual results match in tests using MockURL
func CleanStringFromHandlePagination(data string) string {
	cleaned := StripAnsiEscapeCodes(data)
	cleaned = formatToTabs(cleaned)
	return cleaned
}

// StripAnsiEscapeCodes removes ANSI escape codes from a string.
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

// AssertMeshkitErrorsEqual compares relevant fields of two meshkit errors
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
	ExpectedError    error
	IsOutputGolden   bool
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
	IsOutputGolden   bool
	ExpectedError    error
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
	IsOutputGolden   bool
	ExpectedError    error
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

			buf := SetupMeshkitLoggerTesting(t, false)

			cmd.SetArgs(tt.Args)
			cmd.SetOut(buf)
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
				// Unexpected error - fail immediately
				t.Fatalf("unexpected error: %v", err)
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
			assert.NoError(t, w.Close())

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
		f.Changed = false
	})
	c.PersistentFlags().VisitAll(func(f *pflag.Flag) {
		if err := f.Value.Set(f.DefValue); err != nil {
			t.Fatalf("failed to reset persistent flag %q: %v", f.Name, err)
		}
		f.Changed = false
	})
	for _, sub := range c.Commands() {
		ResetCommandFlags(sub, t)
	}
}
