package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spf13/cobra"
)

type mockCloser struct {
	closeFunc func() error
}

func (m mockCloser) Close() error {
	return m.closeFunc()

}

func TestSafeClose(t *testing.T) {

	log := logrus.New()
	hook := test.NewGlobal()
	log.AddHook(hook)

	// testcases for SafeClose(co io.Closer)
	t.Run("SafeClose", func(t *testing.T) {
		// define a io.Closer for testing
		expectedErr := errors.New("close error")
		mc := &mockCloser{
			closeFunc: func() error {
				return expectedErr
			},
		}
		SafeClose(mc)

		if len(hook.Entries) != 1 {
			t.Fatal("expected 1 log entry")
		}
	})
}

func TestBackupConfigFile(t *testing.T) {
	cfgFile := "testdata/config.yaml"
	tmpFile, err := os.CreateTemp("", "config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := tmpFile.Write(data); err != nil {
		t.Fatal(err)
	}
	BackupConfigFile(tmpFile.Name())
	if _, err := os.Stat("/tmp/config.bak.yaml"); os.IsNotExist(err) {
		t.Errorf("BackupConfigFile failed: backup file does not exist")
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

func TestTruncateID(t *testing.T) {
	id := "1234567890"
	want := "12345678"
	got := TruncateID(id)
	if got != want {
		t.Errorf("TruncateID got = %v want = %v", got, want)
	}
}

func TestPrintToTable(t *testing.T) {
	tests := []struct {
		name   string
		header []string
		data   [][]string
		want   string
	}{
		{
			name:   "Single row table",
			header: []string{"Header1", "Header2"},
			data:   [][]string{{"Data1", "Data2"}},
			want:   "HEADER1\tHEADER2 \nData1  \tData2  \t\n",
		},
		{
			name:   "Multiple rows table",
			header: []string{"Header1", "Header2"},
			data:   [][]string{{"Data1", "Data2"}, {"Data3", "Data4"}},
			want:   "HEADER1\tHEADER2 \nData1  \tData2  \t\nData3  \tData4  \t\n",
		},
		{
			name:   "Empty table",
			header: []string{"Header1", "Header2"},
			data:   [][]string{},
			want:   "HEADER1\tHEADER2 \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintToTable(tt.header, tt.data)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			_, err := io.Copy(&buf, r)
			if err != nil {
				t.Fatal("failed to read from buffer")
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("\nPrintToTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrintToTableWithFooter(t *testing.T) {
	tests := []struct {
		name     string
		header   []string
		data     [][]string
		footer   []string
		expected string
	}{
		{
			name:   "Basic table with footer",
			header: []string{"HEADER1", "HEADER2"},
			data: [][]string{
				{"Data1", "Data2"},
				{"Data3", "Data4"},
			},
			footer: []string{"FOOTER1", "FOOTER2"},
			expected: "HEADER1\tHEADER2 \n" +
				"Data1  \tData2  \t\n" +
				"Data3  \tData4  \t\n" +
				"\n" +
				"  FOOTER1  FOOTER2  \n\n",
		},
		{
			name:   "Empty data with footer",
			header: []string{"HEADER1", "HEADER2"},
			data:   [][]string{},
			footer: []string{"FOOTER1", "FOOTER2"},
			expected: "HEADER1\tHEADER2 \n" +
				"\n" +
				"  FOOTER1  FOOTER2  \n\n",
		},
		{
			name:   "No footer",
			header: []string{"HEADER1", "HEADER2"},
			data: [][]string{
				{"Data1", "Data2"},
				{"Data3", "Data4"},
			},
			footer: []string{},
			expected: "HEADER1\tHEADER2 \n" +
				"Data1  \tData2  \t\n" +
				"Data3  \tData4  \t\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			PrintToTableWithFooter(tt.header, tt.data, tt.footer)

			w.Close()
			out, _ := io.ReadAll(r)
			os.Stdout = old

			if got := string(out); got != tt.expected {
				t.Errorf("PrintToTableWithFooter() = %v, want %v", got, tt.expected)
			}
		})
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
	tests := []struct {
		name     string
		header   []string
		data     [][]string
		expected string
	}{
		{
			name:   "Normal case with multiple rows and columns",
			header: []string{"Header1", "Header2"},
			data: [][]string{
				{"Row1Col1", "Row1Col2"},
				{"Row2Col1", "Row2Col2"},
			},
			expected: "HEADER1 \tHEADER2  \n" +
				"Row1Col1\tRow1Col2\t\n" +
				"Row2Col1\tRow2Col2\t\n",
		},
		{
			name:     "Empty data",
			header:   []string{"Header1", "Header2"},
			data:     [][]string{},
			expected: "HEADER1\tHEADER2 \n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PrintToTableInStringFormat(tt.header, tt.data)
			if got != tt.expected {
				t.Errorf("PrintToTableInStringFormat() = %v, want %v", got, tt.expected)
			}
		})
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

func TestTransformYAML(t *testing.T) {
	tests := []struct {
		name          string
		yamlData      []byte
		transform     func(interface{}) (interface{}, error)
		keys          []string
		expectedData  []byte
		expectedError bool
	}{
		{
			name: "valid transformation with nested keys",
			yamlData: []byte(`
parent:
  child: value
`),
			transform: func(val interface{}) (interface{}, error) {
				return "new_value", nil
			},
			keys: []string{"parent", "child"},
			expectedData: []byte(`
parent:
  child: new_value
`),
			expectedError: false,
		},
		{
			name: "invalid path in keys",
			yamlData: []byte(`
parent:
  child: value
`),
			transform: func(val interface{}) (interface{}, error) {
				return "new_value", nil
			},
			keys:          []string{"parent", "nonexistent"},
			expectedData:  nil,
			expectedError: true,
		},
		{
			name: "error during transformation",
			yamlData: []byte(`
parent:
  child: value
`),
			transform: func(val interface{}) (interface{}, error) {
				return nil, fmt.Errorf("transformation error")
			},
			keys:          []string{"parent", "child"},
			expectedData:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, err := TransformYAML(tt.yamlData, tt.transform, tt.keys...)
			if (err != nil) != tt.expectedError {
				t.Errorf("TransformYAML() error = %v, expectedError %v", err, tt.expectedError)
				return
			}
			if !tt.expectedError && strings.TrimSpace(string(gotData)) != strings.TrimSpace(string(tt.expectedData)) {
				t.Errorf("TransformYAML() got = %s, want %s", gotData, tt.expectedData)
			}
		})
	}
}

func TestMapGet(t *testing.T) {
	tests := []struct {
		name     string
		mp       map[string]interface{}
		key      []string
		expected interface{}
		ok       bool
	}{
		{
			name:     "Nil Map",
			mp:       nil,
			key:      []string{"key"},
			expected: nil,
			ok:       false,
		},
		{
			name:     "Empty Key",
			mp:       map[string]interface{}{"key": "value"},
			key:      []string{},
			expected: map[string]interface{}{"key": "value"},
			ok:       true,
		},
		{
			name:     "Single Level Key",
			mp:       map[string]interface{}{"key": "value"},
			key:      []string{"key"},
			expected: "value",
			ok:       true,
		},
		{
			name: "Nested Map",
			mp: map[string]interface{}{
				"key": map[string]interface{}{
					"nestedKey": "nestedValue",
				},
			},
			key:      []string{"key", "nestedKey"},
			expected: "nestedValue",
			ok:       true,
		},
		{
			name: "Array in Map",
			mp: map[string]interface{}{
				"key": []interface{}{
					map[string]interface{}{"nestedKey": "nestedValue"},
				},
			},
			key:      []string{"key", "0", "nestedKey"},
			expected: "nestedValue",
			ok:       true,
		},
		{
			name:     "Invalid Key",
			mp:       map[string]interface{}{},
			key:      []string{"key1"},
			expected: nil,
			ok:       false,
		},
		{
			name: "Array in Map with Nested Map",
			mp: map[string]interface{}{
				"key1": []map[string]interface{}{
					{"nestedKey1": "nestedValue1"},
					{"nestedKey2": "nestedValue2"},
				},
			},
			key:      []string{"key1", "1", "nestedKey2"},
			expected: "nestedValue2",
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := MapGet(tt.mp, tt.key...)
			if !reflect.DeepEqual(got, tt.expected) || ok != tt.ok {
				t.Errorf("MapGet() = %v, %v; want %v, %v", got, ok, tt.expected, tt.ok)
			}
		})
	}
}
func TestMapSet(t *testing.T) {
	tests := []struct {
		name     string
		mp       map[string]interface{}
		value    interface{}
		keys     []string
		expected map[string]interface{}
	}{
		{
			name:  "Basic Key-Value Pair",
			mp:    map[string]interface{}{},
			value: "value1",
			keys:  []string{"key1"},
			expected: map[string]interface{}{
				"key1": "value1",
			},
		},
		{
			name: "Nested Map",
			mp: map[string]interface{}{
				"key1": map[string]interface{}{
					"key2": "value2",
				},
			},
			value: "newValue",
			keys:  []string{"key1", "key2"},
			expected: map[string]interface{}{
				"key1": map[string]interface{}{
					"key2": "newValue",
				},
			},
		},
		{
			name: "Array in Map",
			mp: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{
						"key2": "value2",
					},
				},
			},
			value: "newValue",
			keys:  []string{"key1", "0", "key2"},
			expected: map[string]interface{}{
				"key1": []interface{}{
					map[string]interface{}{
						"key2": "newValue",
					},
				},
			},
		},
		{
			name:     "Invalid Key",
			mp:       map[string]interface{}{},
			value:    "value1",
			keys:     []string{},
			expected: map[string]interface{}{},
		},
		{
			name: "Array in Map with Nested Map",
			mp: map[string]interface{}{
				"key1": []map[string]interface{}{
					{"nestedKey1": "nestedValue1"},
					{"nestedKey2": "nestedValue2"},
				},
			},
			value: "newValue",
			keys:  []string{"key1", "1", "nestedKey2"},

			expected: map[string]interface{}{
				"key1": []map[string]interface{}{
					{"nestedKey1": "nestedValue1"},
					{"nestedKey2": "newValue"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MapSet(tt.mp, tt.value, tt.keys...)
			if !reflect.DeepEqual(tt.mp, tt.expected) {
				t.Errorf("MapSet() = %v, want %v", tt.mp, tt.expected)
			}
		})
	}
}

func TestRecursiveCastMapStringInterfaceToMapStringInterface(t *testing.T) {
	input := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
	}
	expected := map[string]interface{}{
		"key1": "value1",
		"key2": 2,
	}

	got := RecursiveCastMapStringInterfaceToMapStringInterface(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("RecursiveCastMapStringInterfaceToMapStringInterface() = %v, want %v", got, expected)
	}
}

func TestConvertMapInterfaceMapString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "Convert map[interface{}]interface{} to map[string]interface{}",
			input: map[interface{}]interface{}{
				"key1": "value1",
				2:      "value2",
				"key3": map[interface{}]interface{}{
					"nestedKey1": "nestedValue1",
					3:            "nestedValue2",
				},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"2":    "value2",
				"key3": map[string]interface{}{
					"nestedKey1": "nestedValue1",
					"3":          "nestedValue2",
				},
			},
		},
		{
			name:     "Convert []interface{} to []string",
			input:    []interface{}{"value1", "value2", "value3"},
			expected: []interface{}{"value1", "value2", "value3"},
		},
		{
			name: "Nested map[string]interface{} with map[interface{}]interface{}",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": map[interface{}]interface{}{
					"nestedKey1": "nestedValue1",
					3:            "nestedValue2",
				},
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{
					"nestedKey1": "nestedValue1",
					"3":          "nestedValue2",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertMapInterfaceMapString(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ConvertMapInterfaceMapString() = %v, want %v", result, tt.expected)
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
				Components: []string{"meshery-istio", "meshery-nsm"},
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
					"meshery-nsm", "meshery-nginx-sm", "meshery-traefik-mesh", "meshery-cilium"},
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
		got := SetOverrideValues(tt.ctx, tt.mesheryImageVersion, "", "")
		eq := reflect.DeepEqual(got, tt.want)
		if !eq {
			t.Errorf("SetOverrideValues %s got = %v want = %v", tt.name, got, tt.want)
		}
	}
}

func TestCheckFileExists(t *testing.T) {
	t.Run("File exists", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		exists, err := CheckFileExists(tmpFile.Name())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !exists {
			t.Errorf("expected file to exist, but it does not")
		}
	})

	t.Run("Other error", func(t *testing.T) {
		// Simulate an error by using an invalid file path
		invalidFilePath := string([]byte{0x00})

		exists, err := CheckFileExists(invalidFilePath)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if exists {
			t.Errorf("expected file to not exist, but it does")
		}
		if !strings.Contains(err.Error(), "Failed to read/fetch the file") {
			t.Errorf("expected wrapped error, got %v", err)
		}
	})
}

func TestContains(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		col      []string
		expected int
	}{
		{
			name:     "Key present in slice",
			key:      "test",
			col:      []string{"foo", "bar", "test"},
			expected: 2,
		},
		{
			name:     "Key not present in slice",
			key:      "missing",
			col:      []string{"foo", "bar", "test"},
			expected: -1,
		},
		{
			name:     "Empty slice",
			key:      "test",
			col:      []string{},
			expected: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.key, tt.col); got != tt.expected {
				t.Errorf("Contains() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFindInSlice(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		items         []string
		expectedIndex int
		expectedFound bool
	}{
		{
			name:          "key found at index 0",
			key:           "apple",
			items:         []string{"apple", "banana", "cherry"},
			expectedIndex: 0,
			expectedFound: true,
		},
		{
			name:          "key not found",
			key:           "date",
			items:         []string{"apple", "banana", "cherry"},
			expectedIndex: -1,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, found := FindInSlice(tt.key, tt.items)
			if index != tt.expectedIndex || found != tt.expectedFound {
				t.Errorf("FindInSlice(%v, %v) = (%v, %v), want (%v, %v)", tt.key, tt.items, index, found, tt.expectedIndex, tt.expectedFound)
			}
		})
	}
}

func TestGetPageQueryParameter(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		setFlag  bool
		expected string
	}{
		{
			name:     "Page flag not set",
			page:     1,
			setFlag:  false,
			expected: "pagesize=all",
		},
		{
			name:     "Page flag set to 2",
			page:     2,
			setFlag:  true,
			expected: "page=2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := &cobra.Command{}
			if tt.setFlag {
				cmd.Flags().Int("page", tt.page, "page number")
				err := cmd.Flags().Set("page", fmt.Sprintf("%d", tt.page))
				if err != nil {
					t.Fatalf("failed to set flag: %v", err)
				}
			}

			result := GetPageQueryParameter(cmd, tt.page)
			if result != tt.expected {
				t.Errorf("GetPageQueryParameter() = %v, want %v", result, tt.expected)
			}
		})
	}
}
