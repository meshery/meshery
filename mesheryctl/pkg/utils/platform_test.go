package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/constants"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

func TestChangePlatform(t *testing.T) {
	type args struct {
		currCtx string
		ctx     config.Context
	}

	// Setup path to test config file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testConfigPath := currDir + "/fixtures/TestConfig.yaml"

	// Read and write to the test config file
	SetupCustomContextEnv(t, testConfigPath)

	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())

	err := mctlCfg.SetCurrentContext("local")
	if err != nil {
		t.Error("error setting context", err)
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		t.Error("error processing context from config", err)
	}

	currCtx.SetPlatform("testplatform")

	tests := []struct {
		name    string
		args    args
		wantErr bool
		golden  string
	}{
		{
			name:    "Update platform in gke context (valid context)",
			args:    args{currCtx: "testplatform"},
			wantErr: false,
			golden:  "testchangeplatform.valid.golden",
		},
		{
			name:    "Update platform in kubernetes context (invalid context)",
			args:    args{currCtx: "testplatform"},
			wantErr: true,
			golden:  "testchangeplatform.invalid.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := ChangePlatform(tt.args.currCtx, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ChangePlatform() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualResponse := b.String()
			// Expected response
			testdataDir := filepath.Join(currDir, "testdata/context")
			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)
			if *update {
				golden.Write(actualResponse)
			}
			expectedResponse := golden.Load()
			if expectedResponse != actualResponse {
				t.Errorf("expected response %v and actual response %v don't match", expectedResponse, actualResponse)
			}
		})
	}
}

func TestGetManifestURL(t *testing.T) {
	type args struct {
		manifest        Manifest
		rawManifestsURL string
	}

	rawMfestsURL := "https://testurl.com/k8s/"

	tests := []struct {
		name    string
		args    args
		fixture string
		golden  string
	}{
		{
			name: "Get Manifest URL (blob)",
			args: args{
				rawManifestsURL: rawMfestsURL,
			},
			fixture: "getmanifesturl.blob.golden",
			golden:  "getmanifesturl.expect.blob.golden",
		},
		{
			name: "Get Manifest URL (non-blob)",
			args: args{
				rawManifestsURL: rawMfestsURL,
			},
			fixture: "getmanifesturl.other.golden",
			golden:  "getmanifesturl.expect.other.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Read fixture
			file, _ := ioutil.ReadFile("fixtures/" + tt.fixture)
			manifest := Manifest{}
			_ = json.Unmarshal([]byte(file), &manifest)

			tt.args.manifest = manifest

			got := GetManifestURL(tt.args.manifest, tt.args.rawManifestsURL)

			// Read golden file
			file, _ = ioutil.ReadFile("testdata/" + tt.golden)
			want := string(file)

			if got != want {
				t.Errorf("GetManifestURL() = %v, want %v", got, want)
			}
		})
	}
}

func TestGetPods(t *testing.T) {

	tests := []struct {
		name      string
		namespace string
		wantErr   bool
	}{
		{
			name:      "GetPods() with MesheryNamespace",
			namespace: MesheryNamespace,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// create an kubernetes client
			client, err := meshkitkube.New([]byte(""))

			if err != nil {
				t.Errorf("Error in creating kubernetes client in GetPods() test")
				return
			}

			// List the pods in the MesheryNamespace
			got, err := GetPods(client, tt.namespace)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// check return type
			if fmt.Sprintf("%T", got) != "*v1.PodList" {
				t.Errorf("GetPods() error: want type = *v1.PodList, got type %T", got)
			}

			// check non nil
			if got == nil {
				t.Errorf("GetPods() error: got nil PodList")
			}

		})
	}
}

func TestIsPodRequired(t *testing.T) {
	type args struct {
		requiredPods []string
		pod          string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Pod '" + "test'" + " is required",
			args: args{
				requiredPods: []string{"test", "test1", "test2"},
				pod:          "test",
			},
			want: true,
		},
		{
			name: "Pod '" + "test'" + " is not required",
			args: args{
				requiredPods: []string{"test1", "test2"},
				pod:          "test",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPodRequired(tt.args.requiredPods, tt.args.pod); got != tt.want {
				t.Errorf("IsPodRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChangeConfigEndpoint(t *testing.T) {
	type args struct {
		currCtx string
		ctx     *config.Context
	}

	// Setup path to test config file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testConfigPath := currDir + "/fixtures/TestConfig.yaml"

	SetupCustomContextEnv(t, testConfigPath)

	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())

	err := mctlCfg.SetCurrentContext("local")
	if err != nil {
		t.Error("error setting context", err)
	}

	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		t.Error("error processing context from config", err)
	}

	endpoint_address := "http://localhost:55555"
	currCtx.SetEndpoint(endpoint_address)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ChangeConfigEndpoint1",
			args: args{
				currCtx: mctlCfg.CurrentContext,
				ctx:     currCtx,
			},
			wantErr: false,
		},
		// Can also test with context gke, should give no change since platform isnt docker or kubernetes
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := ChangeConfigEndpoint(tt.args.currCtx, tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("ChangeConfigEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDownloadManifests(t *testing.T) {
	type args struct {
		manifestArr     []Manifest
		rawManifestsURL string
	}

	// initialize mock server for handling requests
	StartMockery(t)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			wantErr: false,
		},
	}

	err := SetFileLocation()
	if err != nil {
		t.Fatalf("Could not set file location")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rawManifestsURL := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/latest/install/deployment_yamls/k8s/"

			apiResponse := "Test Content"
			manifest := Manifest{
				Path: "download-test.yaml",
				Mode: "100644",
				Typ:  "blob",
				SHA:  "d000418ece6776f82cb050a5fd2ecf41f17de49c",
				Size: "752",
				URL:  "https://api.github.com/repos/meshery/meshery/git/blobs/d000418ece6776f82cb050a5fd2ecf41f17de49c",
			}
			// httpmock.RegisterResponder("GET", GetManifestURL(manifest, rawManifestsURL), httpmock.NewStringResponder(200, apiResponse))
			httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/meshery/meshery/latest/install/deployment_yamls/k8s/download-test.yaml", httpmock.NewStringResponder(200, apiResponse))

			manifests := []Manifest{manifest}

			if err := DownloadManifests(manifests, rawManifestsURL); (err != nil) != tt.wantErr {
				t.Errorf("DownloadManifests() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Todo: check file contents and delete

		})
	}

	// stop mock server
	StopMockery(t)
}
