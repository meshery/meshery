package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
)

func getFixturesDirectory() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("helpers_test.go: Cannot get current working directory")
	}
	currentDirectory := filepath.Dir(filename)
	// get the fixtures file directory
	fixturesDir := filepath.Join(currentDirectory, "fixtures")
	return fixturesDir
}

var fixturesDir = getFixturesDirectory()

func TestBackupConfigFile(t *testing.T) {
	name := "config.yaml"
	configFilePath := filepath.Join(fixturesDir, name)

	backupFileName := "config.bak.yaml"
	backupConfigFilePath := filepath.Join(fixturesDir, backupFileName)

	// creates a config file
	NewGoldenFile(t, name, fixturesDir).Write("mesheryctl")

	BackupConfigFile(configFilePath)

	// check if backup file is present or not
	_, err := os.Stat(backupConfigFilePath)
	if err != nil {
		t.Errorf("BackupConfigFile error = %v", err)
	}
}

func TestStringWithCharset(t *testing.T) {
	// checking the length, since this function returns random strings everytime
	strLength := 10

	gotString := StringWithCharset(strLength)

	if len(gotString) != strLength {
		t.Errorf("StringWithCharset got = %v want = %v", len(gotString), strLength)
	}
}

func TestPrereq(t *testing.T) {
	_, _, err := prereq()
	if err != nil {
		t.Errorf("prereq error = %v", err)
	}
}

func TestSetFileLocation(t *testing.T) {
	originalMesheryFolder := MesheryFolder
	originalDockerComposeFile := DockerComposeFile
	originalAuthConfigFile := AuthConfigFile
	originalDefaultConfigPath := DefaultConfigPath
	defer func() {
		MesheryFolder = originalMesheryFolder
		DockerComposeFile = originalDockerComposeFile
		AuthConfigFile = originalAuthConfigFile
		DefaultConfigPath = originalDefaultConfigPath
	}()
	err := SetFileLocation()
	if err != nil {
		t.Errorf("SetFileLocation error = %v", err)
	}
}

func TestNavigateToBrowser(t *testing.T) {
	// opens up a browser window whenever this test runs
	err := NavigateToBrowser("https://layer5.io")
	if err != nil {
		t.Errorf("NavigateToBrowser error: %v", err)
	}
}

func TestUploadFileWithParams(t *testing.T) {
	fixtureFileName := "listmanifest.api.response.golden" // any arbitrary fixture file
	uploadFilePath := filepath.Join(fixturesDir, "platform", fixtureFileName)
	// set token
	TokenFlag = filepath.Join(fixturesDir, "auth.json")
	// returns *http.Request
	_, err := UploadFileWithParams("https://layer5.io", nil, "meshery", uploadFilePath)

	if err != nil {
		t.Errorf("TestUploadFileWithParams error = %v", err)
	}
}

func TestContentTypeIsHTML(t *testing.T) {
	tests := []struct {
		name           string
		response       *http.Response
		expectedOutput bool
	}{
		{
			name: "correct content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html"},
				},
			},
			expectedOutput: true,
		},
		{
			name: "empty content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{},
				},
			},
			expectedOutput: false,
		},
		{
			name: "incorrect content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"multipart/form-data"},
				},
			},
			expectedOutput: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContentTypeIsHTML(tt.response)
			if got != tt.expectedOutput {
				t.Errorf("ContentTypeIsHTML error = %v want = %v", got, tt.expectedOutput)
			}
		})
	}
}

func TestAskForConfirmation(t *testing.T) {
	tests := []struct {
		name     string
		question string
		input    string
		want     bool
	}{
		{
			name:     "test with input 'yes'",
			question: "question?",
			input:    "yes\n",
			want:     true,
		},
		{
			name:     "test with input 'no'",
			question: "question?",
			input:    "no\n",
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mocking stdio
			// see https://stackoverflow.com/a/64518829
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			_, err = w.WriteString(tt.input)
			if err != nil {
				t.Error(err)
			}
			w.Close()
			stdin := os.Stdin
			defer func() { os.Stdin = stdin }()
			os.Stdin = r

			got := AskForConfirmation(tt.question)
			if got != tt.want {
				t.Errorf("AskForConfirmation got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestCreateConfigFile(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Error(err, "failed to get users home directory")
	}
	originalDefaultConfigPath := DefaultConfigPath
	defer func() { DefaultConfigPath = originalDefaultConfigPath }()
	DefaultConfigPath = filepath.Join(home, "config.yaml")
	err = CreateConfigFile()
	if err != nil {
		t.Errorf("CreateConfigFile error = %v", err)
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr string
	}{
		{
			name: "Correct URL",
			url:  "https://www.layer5.io",
		},
		{
			name:    "Unsupported scheme",
			url:     "mqtt://www.layer5.io",
			wantErr: "mqtt is not a supported protocol",
		},
		{
			name:    "invalid URL",
			url:     "layer5.io",
			wantErr: "parse \"layer5.io\": invalid URI for request",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.url)
			if err != nil {
				// check error message
				if err.Error() != tt.wantErr {
					t.Errorf("ValidateURL error = %v want = %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestReadToken(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    map[string]string
	}{
		{
			name:    "with valid JSON",
			fixture: "readtoken.golden",
			want: map[string]string{
				"message": "meshery",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fixtureFilePath := filepath.Join(fixturesDir, tt.fixture)
			got, err := ReadToken(fixtureFilePath)
			if err != nil {
				t.Errorf("ReadToken error = %v", err)
			}
			// checking map equality
			eq := reflect.DeepEqual(got, tt.want)
			if !eq {
				t.Errorf("ReadToken got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestTruncateID(t *testing.T) {
	id := "1234567890"
	want := "12345678"
	got := TruncateID(id)
	if got != want {
		t.Errorf("TruncateID got = %v want = %v", got, want)
	}
}

func TestPrintToTable(t *testing.T) {
	// mocking Stdout
	// https://stackoverflow.com/a/29339052
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// prints to stdout
	PrintToTable([]string{"firstheader", "secondheader"}, [][]string{{"data1", "data2"}, {"data3", "data4"}})

	w.Close()
	// read from stdout
	out, _ := io.ReadAll(r)

	defer func() { os.Stdout = old }()

	got := string(out)
	want := NewGoldenFile(t, "printToTable.golden", fixturesDir).Load()
	if got != want {
		t.Errorf("PrintToTable got = %v want = %v", got, want)
	}
}

func TestPrintToTableWithFooter(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// prints to stdout
	PrintToTableWithFooter([]string{"firstheader", "secondheader"}, [][]string{{"data1", "data2"}, {"data3", "data4"}}, []string{"footer1", "footer2"})

	w.Close()
	// read from stdout
	out, _ := io.ReadAll(r)

	defer func() { os.Stdout = old }()

	got := string(out)
	want := NewGoldenFile(t, "printToTableWithFooter.golden", fixturesDir).Load()
	if got != want {
		t.Errorf("PrintToTableWithFooter got = %v want = %v", got, want)
	}
}

func TestStringContainedInSlice(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		slice []string
		want  int
	}{
		{
			name:  "test with present string",
			str:   "data2",
			slice: []string{"data1", "data2"},
			want:  1,
		},
		{
			name:  "test with absent string",
			str:   "data3",
			slice: []string{"data1", "data2"},
			want:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringContainedInSlice(tt.str, tt.slice)
			if got != tt.want {
				t.Errorf("StringContainedInSlice got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestStringInSlice(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		slice []string
		want  bool
	}{
		{
			name:  "test with string present",
			str:   "data1",
			slice: []string{"data1", "data2"},
			want:  true,
		},
		{
			name:  "test with string absent",
			str:   "data3",
			slice: []string{"data1", "data2"},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := StringInSlice(tt.str, tt.slice)
			if got != tt.want {
				t.Errorf("StringInSlice got = %v want = %v", got, tt.want)
			}
		})
	}
}

func TestAskForInput(t *testing.T) {
	input := "data1"

	// mocking stdio
	// see https://stackoverflow.com/a/64518829
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.WriteString(input + "\n")
	if err != nil {
		t.Error(err)
	}
	w.Close()
	stdin := os.Stdin
	defer func() { os.Stdin = stdin }()
	os.Stdin = r

	got := AskForInput("Prompt", []string{"data1", "data2"})
	if got != input {
		t.Errorf("AskForInput got = %v want = %v", got, input)
	}
}

func TestParseURLGithub(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		rawRepoOutput string
		pathOutput    string
		expectedError string
	}{
		{
			name:          "test with non-github url",
			url:           "https://layer5.io",
			rawRepoOutput: "https://layer5.io",
			pathOutput:    "",
			expectedError: "only github urls are supported",
		},
		{
			name:          "test with github.com",
			url:           "https://github.com/layer5io/meshery/blob/master/.goreleaser.yml",
			rawRepoOutput: "https://github.com/layer5io/meshery/master",
			pathOutput:    ".goreleaser.yml",
		},
		{
			name:          "test with raw.githubusercontent.com",
			url:           "https://raw.githubusercontent.com/layer5io/meshery/master/.goreleaser.yml",
			rawRepoOutput: "https://raw.githubusercontent.com/layer5io/meshery/master/.goreleaser.yml",
			pathOutput:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRawRepoOutput, gotPathOutput, gotError := ParseURLGithub(tt.url)
			// gotRawOutput != tt.rawRepoOutput || gotPathOutput != tt.pathOutput ||
			if gotRawRepoOutput != tt.rawRepoOutput {
				t.Errorf("ParseURLGithub got = %v, want = %v", gotRawRepoOutput, tt.rawRepoOutput)
			}
			if gotPathOutput != tt.pathOutput {
				t.Errorf("ParseURLGithub got = %v, want = %v", gotPathOutput, tt.pathOutput)
			}
			if gotError != nil {
				if gotError.Error() != tt.expectedError {
					t.Errorf("ParseURLGithub error = %v, want = %v", gotError, tt.expectedError)
				}
			}
		})
	}
}

func TestPrintToTableInStringFormat(t *testing.T) {
	want := NewGoldenFile(t, "PrintToTableInStringFormat.golden", fixturesDir).Load()
	got := PrintToTableInStringFormat([]string{"firstheader", "secondheader"}, [][]string{{"data1", "data2"}, {"data3", "data4"}})
	if got != want {
		t.Errorf("PrintToTableInStringFormat got = %v want = %v", got, want)
	}
}

func TestCreateDefaultSpinner(t *testing.T) {
	// only checking for Suffix and FinalMSG
	got := CreateDefaultSpinner("suffix", "message")
	want := struct {
		Suffix   string
		FinalMsg string
	}{
		Suffix:   " suffix", // The leading space is intentional
		FinalMsg: "message\n",
	}
	if want.Suffix != got.Suffix {
		t.Errorf("CreateDefaultSpinner got = %v want = %v", got.Suffix, want.Suffix)
	}
	if want.FinalMsg != got.FinalMSG {
		t.Errorf("CreateDefaultSpinner got = %v want = %v", got.FinalMSG, want.FinalMsg)
	}
}

func TestContainsStringPrefix(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		str   string
		want  bool
	}{
		{
			name:  "str is present in the slice",
			slice: []string{"data1", "data2"},
			str:   "data2",
			want:  true,
		},
		{
			name:  "str is not present in the slice",
			slice: []string{"data1", "data2"},
			str:   "data3",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsStringPrefix(tt.slice, tt.str)
			if got != tt.want {
				t.Errorf("ContainsStringPrefix got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetOverrideValues(t *testing.T) {
	testChannel := "testChannel"

	tests := []struct {
		name                string
		ctx                 *config.Context
		mesheryImageVersion string
		want                map[string]interface{}
	}{
		{
			name: "Context contains no components and no meshery image version",
			ctx: &config.Context{
				Components: nil,
				Channel:    testChannel,
			},
			mesheryImageVersion: "",
			want: map[string]interface{}{
				"meshery-app-mesh": map[string]interface{}{
					"enabled": false,
				},
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
				"meshery-osm": map[string]interface{}{
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
				"image": map[string]interface{}{
					"tag": testChannel + "-",
				},
			},
		},
		{
			name: "Context contains part of all available components and meshery image version",
			ctx: &config.Context{
				Components: []string{"meshery-istio", "meshery-osm", "meshery-nsm"},
				Channel:    testChannel,
			},
			mesheryImageVersion: "testImageVersion",
			want: map[string]interface{}{
				"meshery-app-mesh": map[string]interface{}{
					"enabled": false,
				},
				"meshery-istio": map[string]interface{}{
					"enabled": true,
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
				"meshery-osm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-nginx-sm": map[string]interface{}{
					"enabled": false,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": false,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-testImageVersion",
				},
			},
		},
		{
			name: "Context contains all available components and meshery image version",
			ctx: &config.Context{
				Components: []string{"meshery-istio", "meshery-linkerd", "meshery-consul", "meshery-kuma",
					"meshery-osm", "meshery-nsm", "meshery-nginx-sm", "meshery-traefik-mesh", "meshery-cilium"},
				Channel: testChannel,
			},
			mesheryImageVersion: "testImageVersion",
			want: map[string]interface{}{
				"meshery-app-mesh": map[string]interface{}{
					"enabled": false,
				},
				"meshery-istio": map[string]interface{}{
					"enabled": true,
				},
				"meshery-cilium": map[string]interface{}{
					"enabled": true,
				},
				"meshery-linkerd": map[string]interface{}{
					"enabled": true,
				},
				"meshery-consul": map[string]interface{}{
					"enabled": true,
				},
				"meshery-kuma": map[string]interface{}{
					"enabled": true,
				},
				"meshery-osm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-nsm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-nginx-sm": map[string]interface{}{
					"enabled": true,
				},
				"meshery-traefik-mesh": map[string]interface{}{
					"enabled": true,
				},
				"image": map[string]interface{}{
					"tag": testChannel + "-testImageVersion",
				},
			},
		},
	}

	for _, tt := range tests {
		got := SetOverrideValues(tt.ctx, tt.mesheryImageVersion)
		eq := reflect.DeepEqual(got, tt.want)
		if !eq {
			t.Errorf("SetOverrideValues %s got = %v want = %v", tt.name, got, tt.want)
		}
	}
}
