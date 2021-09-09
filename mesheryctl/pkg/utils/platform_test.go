package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/constants"
	meshkitkube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// Set absolute path during testing (similar to SetFileLocation in helpers.go)
func SetFileLocationDuringTest(baseDirectory string) error {

	// Get absolute path to current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return errors.New("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	baseDirectory = filepath.Join(currDir, baseDirectory)

	MesheryFolder = filepath.Join(baseDirectory, MesheryFolder)
	manifestsFolder := filepath.Join(MesheryFolder, "manifests")

	err := os.MkdirAll(manifestsFolder, 0755)
	if err != nil {
		return errors.Wrap(err, "Failed to create test directory")
	}

	return nil
}

func TestChangePlatform(t *testing.T) {
	type args struct {
		contextName string
		platform    string
	}

	// Setup path to test config file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}

	currDir := filepath.Dir(filename)
	fixtureDir := currDir + "/fixtures"
	fixture := fixtureDir + "/TestConfig.yaml"

	// Read and write to the test config file
	SetupCustomContextEnv(t, fixture)

	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())

	tests := []struct {
		name    string
		args    args
		wantErr bool
		golden  string
	}{
		{
			name:    "Update platform in gke context (valid context)",
			args:    args{contextName: "gke", platform: "testplatform"},
			wantErr: false,
			golden:  "changeplatform.valid.golden",
		},
		{
			name:    "Update platform in kubernetes context (invalid context)",
			args:    args{contextName: "kubernetes", platform: "testplatform"},
			wantErr: true,
			golden:  "changeplatform.invalid.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := mctlCfg.SetCurrentContext(tt.args.contextName)
			if err != nil {

				if (err != nil) != tt.wantErr {
					t.Fatal("Error setting context", err)
				} else {
					return
				}
			}

			currCtx, err := mctlCfg.GetCurrentContext()
			if err != nil {
				if (err != nil) != tt.wantErr {
					t.Fatal("Error processing context from config", err)
				} else {
					return
				}
			}

			currCtx.SetPlatform(tt.args.platform)

			if err := ChangePlatform(tt.args.contextName, *currCtx); (err != nil) != tt.wantErr {
				t.Errorf("ChangePlatform() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Actual file contents
			actualContent, err := ioutil.ReadFile(fixture)
			if err != nil {
				t.Error(err)
			}

			actualFileContent := string(actualContent)

			// Expected file contents
			testdataDir := currDir + "/testdata"

			golden := NewGoldenFile(t, tt.golden, testdataDir)
			if *update {
				golden.Write(actualFileContent)
			}
			expectedFileContent := golden.Load()

			if expectedFileContent != actualFileContent {
				t.Errorf("expected file content [%v] and actual file content [%v] don't match", expectedFileContent, actualFileContent)
			}

			// Repopulating Expected yaml
			if err := Populate(currDir+"/fixtures/original/TestConfig.yaml", fixture); err != nil {
				t.Error(err, "Could not complete test. Unable to repopulate fixture")
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

	// initialize mock server for handling requests
	StartMockery(t)

	tests := []struct {
		name     string
		manifest Manifest
		golden   string
		wantErr  bool
	}{
		{
			name: "Download test manifest file and check content",
			manifest: Manifest{
				Path: "download-test.yaml",
				Mode: "100644",
				Typ:  "blob",
				SHA:  "d000418ece6776f82cb050a5fd2ecf41f17de49c",
				Size: "752",
				URL:  "https://api.github.com/repos/meshery/meshery/git/blobs/d000418ece6776f82cb050a5fd2ecf41f17de49c",
			},
			golden:  "downloadmanifests.expect.golden",
			wantErr: false,
		},
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	rawManifestsURL := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/latest/install/deployment_yamls/k8s/"
	apiResponse := "Test Content"

	err := SetFileLocationDuringTest("temp")
	if err != nil {
		t.Fatalf("Could not set file location")
	}

	httpmock.RegisterResponder("GET", "https://raw.githubusercontent.com/meshery/meshery/latest/install/deployment_yamls/k8s/download-test.yaml", httpmock.NewStringResponder(200, apiResponse))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			manifests := []Manifest{tt.manifest}

			if err := DownloadManifests(manifests, rawManifestsURL); (err != nil) != tt.wantErr {
				t.Errorf("DownloadManifests() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Actual file contents
			actualContent, err := ioutil.ReadFile(filepath.Join(MesheryFolder, ManifestsFolder, tt.manifest.Path))
			if err != nil {
				t.Error(err)
			}

			actualFileContent := string(actualContent)

			// Expected file contents
			testdataDir := currDir + "/testdata"

			golden := NewGoldenFile(t, tt.golden, testdataDir)
			if *update {
				golden.Write(actualFileContent)
			}
			expectedFileContent := golden.Load()

			if expectedFileContent != actualFileContent {
				t.Errorf("expected file content [%v] and actual file content [%v] don't match", expectedFileContent, actualFileContent)
			}

			err = os.Remove(filepath.Join(MesheryFolder, ManifestsFolder, tt.manifest.Path))
			if err != nil {
				t.Errorf("Could not delete manifest from test folder")
			}

		})
	}

	// stop mock server
	StopMockery(t)
}

func TestDownloadOperatorManifest(t *testing.T) {

	// initialize mock server for handling requests
	StartMockery(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)

	rawManifestsURL := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/latest/install/deployment_yamls/k8s/"
	apiResponse := "Test Content"

	err := SetFileLocationDuringTest("temp")
	if err != nil {
		t.Fatalf("Could not set file location")
	}

	tests := []struct {
		name      string
		urls      []string
		filenames []string
		goldens   []string
		wantErr   bool
	}{
		{
			name: "Download operator manifests with correct URLs",
			urls: []string{
				"https://github.com/layer5io/meshery-operator/blob/master/config/manifests/default.yaml",
				"https://github.com/layer5io/meshery-operator/blob/master/config/samples/meshery_v1alpha1_broker.yaml",
				"https://github.com/layer5io/meshery-operator/blob/master/config/samples/meshery_v1alpha1_meshsync.yaml",
			},
			filenames: []string{
				"default.yaml",
				"meshery_v1alpha1_broker.yaml",
				"meshery_v1alpha1_meshsync.yaml",
			},
			goldens: []string{
				"downloadoperatormanifest.default.expect.golden",
				"downloadoperatormanifest.meshery_v1alpha1_broker.expect.golden",
				"downloadoperatormanifest.default.meshery_v1alpha1_meshsync.golden",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, url := range tt.urls {
				httpmock.RegisterResponder("GET", url, httpmock.NewStringResponder(200, apiResponse))
			}

			if err := DownloadOperatorManifest(); (err != nil) != tt.wantErr {
				t.Errorf("DownloadOperatorManifest() error = %v, wantErr %v", err, tt.wantErr)
			}

			actualFileContents := []string{}

			for _, filename := range tt.filenames {

				// Actual file contents
				actualContent, err := ioutil.ReadFile(filepath.Join(MesheryFolder, ManifestsFolder, filename))
				if err != nil {
					t.Error(err)
				}

				actualFileContent := string(actualContent)
				actualFileContents = append(actualFileContents, actualFileContent)

			}

			// Expected file contents
			testdataDir := currDir + "/testdata/downloadoperatormanifest"

			for i, golden := range tt.goldens {

				golden := NewGoldenFile(t, golden, testdataDir)
				if *update {
					golden.Write(actualFileContents[i])
				}
				expectedFileContent := golden.Load()

				if expectedFileContent != actualFileContents[i] {
					t.Errorf("For file [%v], Expected file content [%v] and actual file content [%v] don't match", tt.filenames[i], expectedFileContent, actualFileContent)
				}

			}

			err = os.Remove(filepath.Join(MesheryFolder, ManifestsFolder, tt.manifest.Path))
			if err != nil {
				t.Errorf("Could not delete manifest from test folder")
			}

		})
	}

	// stop mock server
	StopMockery(t)
}
