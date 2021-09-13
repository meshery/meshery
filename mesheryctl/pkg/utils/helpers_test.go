package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func getFixturesDirectory() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("TestBackupConfigFile: Cannot get current working directory")
	}
	currDir := filepath.Dir(filename)
	// get the fixtures file directory
	fixturesDir := filepath.Join(currDir, "fixtures")
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
	actualStringLength := len(gotString)

	if actualStringLength != strLength {
		t.Errorf("StringWithCharset got = %v want = %v", actualStringLength, strLength)
	}
}

func TestPrereq(t *testing.T) {
	_, _, err := prereq()
	if err != nil {
		t.Errorf("prereq error = %v", err)
	}
}

func TestSetFileLocation(t *testing.T) {
	err := SetFileLocation()
	if err != nil {
		t.Errorf("SetFileLocation error = %v", err)
	}
}

// func TestNavigateToBrowser(t *testing.T) {
// 	// opens up a browser window whenever this test runs
// 	err := NavigateToBrowser("https://www.layer5.io")
// 	if err != nil {
// 		t.Errorf("NavigateToBrowser error: %v", err)
// 	}
// }

func TestUploadFileWithParams(t *testing.T) {
	fixtureFileName := "listmanifest.api.response.golden" // any arbitrary fixture file
	uploadFilePath := filepath.Join(fixturesDir, fixtureFileName)

	test := struct {
		uri       string
		params    map[string]string
		paramName string
		path      string
	}{
		uri:       "https://www.layer5.io",
		params:    nil,
		paramName: "meshery",
		path:      uploadFilePath,
	}

	// returns *http.Request
	_, err := UploadFileWithParams(test.uri, test.params, test.paramName, test.path)

	if err != nil {
		t.Errorf("UploadFileWithParams error = %v", err)
	}
}

func TestContentTypeIsHTML(t *testing.T) {
	tests := []struct {
		name     string
		response *http.Response
		want     bool
	}{
		{
			name: "correct content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"text/html"},
				},
			},
			want: true,
		},
		{
			name: "empty content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{},
				},
			},
			want: false,
		},
		{
			name: "incorrect content-type",
			response: &http.Response{
				Header: http.Header{
					"Content-Type": []string{"multipart/form-data"},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContentTypeIsHTML(tt.response)
			if got != tt.want {
				t.Errorf("ContentTypeIsHTML error = %v want = %v", got, tt.want)
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

// func TestCreateConfigFile(t *testing.T) {
// 	err := CreateConfigFile()
// 	if err != nil {
// 		t.Errorf("CreateConfigFile error = %v", err)
// 	}
// }

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
	out, _ := ioutil.ReadAll(r)

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
	out, _ := ioutil.ReadAll(r)

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
			url:           "https://www.layer5.io",
			rawRepoOutput: "https://www.layer5.io",
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
