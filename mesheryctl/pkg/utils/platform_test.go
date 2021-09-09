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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func init() {

	// Get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("Not able to get current working directory")
	}

	basePath := filepath.Dir(filename)
	SetFileLocationTesting(basePath)

	// Create required directories for testing, if not already existing
	manifestsFolder := filepath.Join(MesheryFolder, "manifests")

	err := os.MkdirAll(manifestsFolder, os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create test directory")
	}

}

func TestChangePlatform(t *testing.T) {
	type args struct {
		contextName string
		platform    string
	}

	currDir := GetBasePath(t)
	fixtureDir := currDir + "/fixtures/platform"
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
					// handles the case when an invalid context was intentionally supplied
					return
				}
			}

			currCtx, err := mctlCfg.GetCurrentContext()
			if err != nil {
				t.Fatal("Error processing context from config: ", err)
			}

			currCtx.SetPlatform(tt.args.platform)

			if err := ChangePlatform(tt.args.contextName, *currCtx); (err != nil) != tt.wantErr {
				t.Errorf("ChangePlatform() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Actual file contents
			actualContent, err := ioutil.ReadFile(fixture)
			if err != nil {
				t.Error("Error reading actual file contents: ", err)
			}

			actualFileContent := string(actualContent)

			// Expected file contents
			testdataDir := currDir + "/testdata/platform"

			golden := NewGoldenFile(t, tt.golden, testdataDir)
			if *update {
				golden.Write(actualFileContent)
			}
			expectedFileContent := golden.Load()

			if expectedFileContent != actualFileContent {
				t.Errorf("Expected file content \n[%v]\n and actual file content \n[%v]\n don't match", expectedFileContent, actualFileContent)
			}

			// Repopulating Expected yaml
			if err := Populate(currDir+"/fixtures/platform/original/TestConfig.yaml", fixture); err != nil {
				t.Fatal(err, "Could not complete test. Unable to repopulate fixture")
			}

		})
	}
}

func TestChangeConfigEndpoint(t *testing.T) {

	// Setup path to test config file
	currDir := GetBasePath(t)
	testConfigPath := currDir + "/fixtures/platform/TestChangeEndpointConfig.yaml"

	SetupCustomContextEnv(t, testConfigPath)

	mctlCfg, _ := config.GetMesheryCtl(viper.GetViper())

	tests := []struct {
		name            string
		ctxName         string
		endpointAddress string
		golden          string
		wantErr         bool
	}{
		{
			name:            "ChangeConfigEndpoint with platform docker",
			ctxName:         "local",
			endpointAddress: "http://localhost:55555",
			golden:          "changeconfigendpoint.expect.docker.golden",
			wantErr:         false,
		},
		{
			name:            "ChangeConfigEndpoint with platform kubernetes",
			ctxName:         "local2",
			endpointAddress: "http://localhost:44444",
			golden:          "changeconfigendpoint.expect.kubernetes.golden",
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			err := mctlCfg.SetCurrentContext(tt.ctxName)
			if err != nil {
				t.Fatal("error setting context", err)
			}

			currCtx, err := mctlCfg.GetCurrentContext()
			if err != nil {
				t.Fatal("error processing context from config", err)
			}

			currCtx.SetEndpoint(tt.endpointAddress)

			if err := ChangeConfigEndpoint(tt.ctxName, currCtx); (err != nil) != tt.wantErr {
				t.Errorf("ChangeConfigEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Actual file contents
			actualContent, err := ioutil.ReadFile(testConfigPath)
			if err != nil {
				t.Error("Error reading actual file contents: ", err)
			}

			actualFileContent := string(actualContent)

			// Expected file contents
			testdataDir := currDir + "/testdata/platform"

			golden := NewGoldenFile(t, tt.golden, testdataDir)
			if *update {
				golden.Write(actualFileContent)
			}
			expectedFileContent := golden.Load()

			if expectedFileContent != actualFileContent {
				t.Errorf("Expected file content \n[%v]\n and actual file content \n[%v]\n don't match", expectedFileContent, actualFileContent)
			}

			// Repopulating Expected yaml
			if err := Populate(currDir+"/fixtures/platform/original/TestChangeEndpointConfig.yaml", testConfigPath); err != nil {
				t.Fatal(err, "Could not complete test. Unable to repopulate fixture")
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
			file, _ := ioutil.ReadFile("fixtures/platform/" + tt.fixture)
			manifest := Manifest{}
			_ = json.Unmarshal([]byte(file), &manifest)

			tt.args.manifest = manifest

			got := GetManifestURL(tt.args.manifest, tt.args.rawManifestsURL)

			// Read golden file
			golden := NewGoldenFile(t, tt.golden, "testdata/platform/")
			if *update {
				golden.Write(got)
			}

			want := golden.Load()

			if got != want {
				t.Fatalf("GetManifestURL() = %v, want %v", got, want)
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
				t.Fatalf("Error in creating kubernetes client in GetPods() test")
			}

			// List the pods in the MesheryNamespace
			got, err := GetPods(client, tt.namespace)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GetPods() error = %v, wantErr %v", err, tt.wantErr)
			}

			// check return type
			if fmt.Sprintf("%T", got) != "*v1.PodList" {
				t.Errorf("GetPods() error: want type = *v1.PodList, got type %T", got)
			}

			// check non nil
			if got == nil {
				t.Fatalf("GetPods() error: got nil PodList")
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
				t.Fatalf("IsPodRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDownloadManifests(t *testing.T) {

	// initialize mock server for handling requests
	StartMockery(t)

	currDir := GetBasePath(t)

	rawManifestsURL := "https://raw.githubusercontent.com/" + constants.GetMesheryGitHubOrg() + "/" + constants.GetMesheryGitHubRepo() + "/latest/install/deployment_yamls/k8s/"
	apiResponse := "Test Content"

	tests := []struct {
		name     string
		url      string
		manifest Manifest
		golden   string
		wantErr  bool
	}{
		{
			name: "Download test manifest file and check content",
			url:  "https://raw.githubusercontent.com/meshery/meshery/latest/install/deployment_yamls/k8s/download-test.yaml",
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			httpmock.RegisterResponder("GET", tt.url, httpmock.NewStringResponder(200, apiResponse))

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
			testdataDir := currDir + "/testdata/platform"

			golden := NewGoldenFile(t, tt.golden, testdataDir)
			if *update {
				golden.Write(actualFileContent)
			}
			expectedFileContent := golden.Load()

			if expectedFileContent != actualFileContent {
				t.Errorf("Expected file content [%v] and actual file content [%v] don't match", expectedFileContent, actualFileContent)
			}

			err = os.Remove(filepath.Join(MesheryFolder, ManifestsFolder, tt.manifest.Path))
			if err != nil {
				t.Fatalf("Could not delete manifest from test folder")
			}

		})
	}

	// stop mock server
	StopMockery(t)
}

func TestDownloadOperatorManifest(t *testing.T) {

	// initialize mock server for handling requests
	StartMockery(t)

	currDir := GetBasePath(t)

	apiResponse := "Test Content"

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
				"https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/manifests/default.yaml",
				"https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_broker.yaml",
				"https://raw.githubusercontent.com/layer5io/meshery-operator/master/config/samples/meshery_v1alpha1_meshsync.yaml",
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

			actualFileContents := make([]string, len(tt.filenames))

			for i, filename := range tt.filenames {

				// Actual file contents
				actualContent, err := ioutil.ReadFile(filepath.Join(MesheryFolder, ManifestsFolder, filename))
				if err != nil {
					t.Error(err)
				}

				actualFileContent := string(actualContent)
				actualFileContents[i] = actualFileContent

			}

			// Expected file contents
			testdataDir := currDir + "/testdata/platform/downloadoperatormanifest"

			for i, golden := range tt.goldens {

				golden := NewGoldenFile(t, golden, testdataDir)
				if *update {
					golden.Write(actualFileContents[i])
				}
				expectedFileContent := golden.Load()

				if expectedFileContent != actualFileContents[i] {
					t.Errorf("For file [%v], Expected file content [%v] and actual file content [%v] don't match", tt.filenames[i], expectedFileContent, actualFileContents[i])
				}

				err := os.Remove(filepath.Join(MesheryFolder, ManifestsFolder, tt.filenames[i]))
				if err != nil {
					t.Errorf("Could not delete operator manifest [%v] from test folder", tt.filenames[i])
				}

			}

		})
	}

	// stop mock server
	StopMockery(t)
}

func TestGetManifestTreeURL(t *testing.T) {
	// initialize mock server for handling requests
	StartMockery(t)

	// get current directory
	currDir := GetBasePath(t)

	fixturesDir := filepath.Join(currDir, "fixtures/platform")
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
			testdataDir := filepath.Join(currDir, "testdata/platform")
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
	currDir := GetBasePath(t)

	fixturesDir := filepath.Join(currDir, "fixtures/platform")
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
			testdataDir := filepath.Join(currDir, "testdata/platform")
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
