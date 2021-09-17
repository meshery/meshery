package utils

import (
	"encoding/json"
	"flag"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/constants"
)

var update = flag.Bool("update", false, "update golden files")

var versionTests = []struct {
	v        string // input
	expected [3]int // expected
}{
	{"Client Version: v1.10.1", [3]int{1, 10, 1}},
	{"client Version: v1.10.1", [3]int{1, 10, 1}},
	{"Client Version - v1.10.1", [3]int{1, 10, 1}},
	{"v1.10.1", [3]int{1, 10, 1}},
	{"Client Version: v1.10.1 beta:2", [3]int{1, 10, 1}},
	{"Client Version: v2.1348.1", [3]int{2, 1348, 1}},
}

func TestParseKubectlShortVersion(t *testing.T) {
	for _, tt := range versionTests {
		actual, err := parseKubectlShortVersion(tt.v)
		if err != nil {
			t.Fatalf("Unexpected error while parsing kubectl short version: %v", err)
		}
		if actual != tt.expected {
			t.Fatalf("Expected to get %v but got %v", tt.expected, actual)
		}
	}
}

func TestParseKubectlShortVersionIncorrectVersion(t *testing.T) {
	_, err := parseKubectlShortVersion("Not really a version")
	if err == nil {
		t.Fatalf("Expected to get an error")
	}
}
func TestMapSet(t *testing.T) {
	type args struct {
		mp    map[string]interface{}
		value interface{}
		key   []string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{
			name: "insert at root level",
			args: args{
				mp:    map[string]interface{}{},
				value: 1234,
				key:   []string{"k1"},
			},
			want: map[string]interface{}{
				"k1": 1234,
			},
		},
		{
			name: "insert in nested map",
			args: args{
				mp: map[string]interface{}{
					"k1": 123,
					"k2": map[string]interface{}{
						"nk1": map[string]interface{}{},
					},
				},
				value: []string{"val1", "val2"},
				key:   []string{"k2", "nk1", "nk1k1"},
			},
			want: map[string]interface{}{
				"k1": 123,
				"k2": map[string]interface{}{
					"nk1": map[string]interface{}{
						"nk1k1": []string{"val1", "val2"},
					},
				},
			},
		},
		{
			name: "insert in map when key is invalid",
			args: args{
				mp:    map[string]interface{}{},
				value: []string{"val1", "val2"},
				key:   []string{"k2", "nk1", "nk1k1"},
			},
			want: map[string]interface{}{},
		},
		{
			name: "insert in nested map with map slice",
			args: args{
				mp: map[string]interface{}{
					"k1": []map[string]interface{}{
						{
							"0nk1": 1234,
						},
						{
							"1nk1": 12345,
						},
					},
				},
				value: 123456,
				key:   []string{"k1", "1", "1nk1"},
			},
			want: map[string]interface{}{
				"k1": []map[string]interface{}{
					{
						"0nk1": 1234,
					},
					{
						"1nk1": 123456,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if MapSet(tt.args.mp, tt.args.value, tt.args.key...); !reflect.DeepEqual(tt.args.mp, tt.want) {
				t.Errorf("MapSet() = %v, want %v", tt.args.mp, tt.want)
			}
		})
	}
}

func TestMapGet(t *testing.T) {
	type args struct {
		mp  map[string]interface{}
		key []string
	}
	tests := []struct {
		name  string
		args  args
		want  interface{}
		want1 bool
	}{
		{
			name: "when value exists at root level",
			args: args{
				mp: map[string]interface{}{
					"k1": 1234,
				},
				key: []string{"k1"},
			},
			want:  1234,
			want1: true,
		},
		{
			name: "when value exists inside nested map",
			args: args{
				mp: map[string]interface{}{
					"k1": map[string]interface{}{
						"k2": 1234,
					},
				},
				key: []string{"k1", "k2"},
			},
			want:  1234,
			want1: true,
		},
		{
			name: "when value exists inside nested map slice",
			args: args{
				mp: map[string]interface{}{
					"k1": map[string]interface{}{
						"k2": []map[string]interface{}{
							{
								"k3": []string{"val"},
							},
						},
					},
				},
				key: []string{"k1", "k2", "0", "k3"},
			},
			want:  []string{"val"},
			want1: true,
		},
		{
			name: "when value doesn't exist",
			args: args{
				mp: map[string]interface{}{
					"k1": 1234,
				},
				key: []string{"k2"},
			},
			want:  nil,
			want1: false,
		},
		{
			name: "when value isn't deep within the map",
			args: args{
				mp: map[string]interface{}{
					"k1": 1234,
					"k2": map[string]interface{}{
						"k3": map[string]interface{}{},
					},
				},
				key: []string{"k2", "k3", "k4"},
			},
			want:  nil,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := MapGet(tt.args.mp, tt.args.key...)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapGet() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("MapGet() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTransformYAML(t *testing.T) {
	type args struct {
		yamlByt   []byte
		keys      []string
		transform func(interface{}) (interface{}, error)
	}

	sampleYAML := []byte(`apiVersion: v1
kind: Service
metadata:
  name: svc1
spec:`)

	sampleYAML2 := []byte(`apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 1
  template:
    spec:
      containers:
      - env:
        - name: EVENT
          value: mesheryLocal
        - name: PROVIDER_BASE_URLS
          value: https://meshery.layer5.io
        image: layer5/meshery:stable-v0.5.37`)

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "Modify nested string field",
			args: args{
				yamlByt: []byte(string(sampleYAML)), // not a mistake - creating a copy
				keys:    []string{"metadata", "name"},
				transform: func(i interface{}) (interface{}, error) {
					val, ok := i.(string)
					if !ok {
						return nil, fmt.Errorf("expected type to be string, got: %T", val)
					}

					return "svc1-changed", nil
				},
			},
			want: []byte(`apiVersion: v1
kind: Service
metadata:
  name: svc1-changed
spec: null
`),
			wantErr: false,
		},
		{
			name: "Modify field which does not exist",
			args: args{
				yamlByt: []byte(string(sampleYAML)), // not a mistake - creating a copy
				keys:    []string{"metadata", "name", "t1"},
				transform: func(i interface{}) (interface{}, error) {
					val, ok := i.(string)
					if !ok {
						return nil, fmt.Errorf("expected type to be string, got: %T", val)
					}

					return "svc1-changed", nil
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Modify field at root level",
			args: args{
				yamlByt: []byte(string(sampleYAML)), // not a mistake - creating a copy
				keys:    []string{"spec"},
				transform: func(i interface{}) (interface{}, error) {
					return map[string]interface{}{
						"selector": map[string]interface{}{
							"app": "dep",
						},
					}, nil
				},
			},
			want: []byte(`apiVersion: v1
kind: Service
metadata:
  name: svc1
spec:
  selector:
    app: dep
`),
			wantErr: false,
		},
		{
			name: "Modify field at root level",
			args: args{
				yamlByt: []byte(string(sampleYAML2)), // not a mistake - creating a copy
				keys:    []string{"spec", "template", "spec", "containers", "0", "env"},
				transform: func(i interface{}) (interface{}, error) {
					envVarI, ok := i.([]interface{})
					if !ok {
						return i, fmt.Errorf("unexpected data type")
					}

					return append(envVarI, map[string]interface{}{
						"name":  "TEMP",
						"value": "tempvalue",
					}), nil
				},
			},
			want: []byte(`apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 1
  template:
    spec:
      containers:
      - env:
        - name: EVENT
          value: mesheryLocal
        - name: PROVIDER_BASE_URLS
          value: https://meshery.layer5.io
        - name: TEMP
          value: tempvalue
        image: layer5/meshery:stable-v0.5.37
`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TransformYAML(tt.args.yamlByt, tt.args.transform, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TransformYAML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TransformYAML() = %+#v, want %+#v", got, tt.want)
			}
		})
	}
}
func TestGetManifestTreeURL(t *testing.T) {
	// initialize mock server for handling requests
	StartMockery(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		name         string
		version      string
		url          string
		expectOutput string
		fixture      string
		expectErr    bool
	}{
		{
			name:         "Test getting manifest tree url",
			expectErr:    false,
			expectOutput: "manifesturl.expect.golden",
			fixture:      "manifesturl.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/v0.5.45?recursive=1",
			version:      "v0.5.45",
		},
		{
			name:         "Test getting manifest tree url with wrong version",
			expectErr:    true,
			expectOutput: "manifesturl.out.err.golden",
			fixture:      "manifesturl.err.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/v0000000?recursive=1",
			version:      "v0000000",
		},
		{
			name:         "Test getting manifest tree url with empty version",
			expectErr:    true,
			expectOutput: "manifesturlempty.out.err.golden",
			fixture:      "manifesturl.empty.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/?recursive=1",
			version:      "",
		},
		{
			name:         "Test getting manifest tree url with wrong version",
			expectErr:    false,
			expectOutput: "manifesturl.out.latest.golden",
			fixture:      "manifesturl.latest.api.response.golden",
			url:          "https://api.github.com/repos/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/git/trees/" + "master" + "?recursive=1",
			version:      "master",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := NewGoldenFile(t, tt.expectOutput, testdataDir)
			// mock response
			apiResponse := NewGoldenFile(t, tt.fixture, fixturesDir).Load()
			httpmock.RegisterResponder("GET", tt.url,
				httpmock.NewStringResponder(200, apiResponse))
			actualurl, err := GetManifestTreeURL(tt.version)
			if err != nil {
				if tt.expectErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					Equals(t, expectedResponse, err.Error())
					return
				}
			}
			if *update {
				golden.Write(actualurl)
			}
			expectedResponse := golden.Load()
			Equals(t, expectedResponse, actualurl)
		})
	}
	// stop mock server
	StopMockery(t)
}

func TestListManifests(t *testing.T) {
	// initialize mock server for handling requests
	StartMockery(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []struct {
		name         string
		url          string
		expectOutput string
		fixture      string
		expectErr    bool
	}{
		{
			name:         "Test listing manifests",
			expectErr:    false,
			expectOutput: "listmanifest.expect.golden",
			fixture:      "listmanifest.api.response.golden",
			url:          "https://api.github.com/repos/meshery/meshery/git/trees/47c634a49e6d143a54d734437a26ad233146ddf5",
		},
		{
			name:         "Test listing manifests with wrong url",
			expectErr:    true,
			expectOutput: "listmanifest.expect.err.golden",
			fixture:      "listmanifest.api.err.response.golden",
			url:          "https://api.github.com/repos/meshery/meshery/git/trees/gibberish",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata")
			golden := NewGoldenFile(t, tt.expectOutput, testdataDir)
			// mock response
			apiResponse := NewGoldenFile(t, tt.fixture, fixturesDir).Load()
			httpmock.RegisterResponder("GET", tt.url,
				httpmock.NewStringResponder(200, apiResponse))
			manifests, err := ListManifests(tt.url)
			if err != nil {
				if tt.expectErr {
					if *update {
						golden.Write(err.Error())
					}
					expectedResponse := golden.Load()
					Equals(t, expectedResponse, err.Error())
					return
				}
			}
			manifestsactual, err := json.Marshal(&manifests)
			if err != nil {
				t.Error("Could not unmarshall manifests from response")
			}
			if *update {
				golden.Write(string(manifestsactual))
			}
			expectedResponse := golden.Load()
			Equals(t, expectedResponse, string(manifestsactual))
		})
	}
	// stop mock server
	StopMockery(t)
}
