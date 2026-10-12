package connections

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	mesheryctllogger "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/logger"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func init() {
	if utils.Log == nil {
		utils.Log = mesheryctllogger.GetMeshkitLogger(logrus.InfoLevel)
	}
}

// TestConnectionIDForContext exercises the connection-id resolution that backs
// the `connection create` -> `connection view`/`connection delete` handoff. The
// id must be resolved by an EXACT context-name match; a multi-context kubeconfig
// registers every context in one response, so an arbitrary fallback could return
// a different context's id than the one the user selected (the "arbitrary
// conn-id" review fix). When the server does not echo the requested context the
// function must return "" rather than guess.
func TestConnectionIDForContext(t *testing.T) {
	const (
		minikubeID = "aaaaaaaa-1111-1111-1111-111111111111"
		otherID    = "bbbbbbbb-2222-2222-2222-222222222222"
	)

	tests := []struct {
		name     string
		response saveK8sContextResponse
		cname    string
		want     string
	}{
		{
			name: "exact match in registeredContexts returns its id (positive path)",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "exact match in connectedContexts returns its id",
			response: saveK8sContextResponse{
				ConnectedContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "multi-context response returns the requested context, not the first",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "prod-cluster", ConnectionID: otherID},
					{Name: "minikube", ConnectionID: minikubeID},
				},
			},
			cname: "minikube",
			want:  minikubeID,
		},
		{
			name: "requested context not echoed returns empty, never an arbitrary id",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "prod-cluster", ConnectionID: otherID},
					{Name: "staging", ConnectionID: "cccccccc-3333-3333-3333-333333333333"},
				},
			},
			cname: "minikube",
			want:  "",
		},
		{
			name: "name matches but id empty returns empty",
			response: saveK8sContextResponse{
				RegisteredContexts: []registeredK8sContext{
					{Name: "minikube", ConnectionID: ""},
				},
			},
			cname: "minikube",
			want:  "",
		},
		{
			name:     "empty response returns empty",
			response: saveK8sContextResponse{},
			cname:    "minikube",
			want:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := connectionIDForContext(tt.response, tt.cname); got != tt.want {
				t.Fatalf("connectionIDForContext(%q) = %q, want %q", tt.cname, got, tt.want)
			}
		})
	}
}

// TestSaveK8sContextResponseWireContract pins the camelCase wire tags the CLI
// parses off the server's /api/system/kubernetes response. If a tag drifts to
// snake_case the id resolution silently returns "" and the create->view/delete
// handoff regresses without a compile error, so assert the contract directly.
func TestSaveK8sContextResponseWireContract(t *testing.T) {
	const body = `{
		"connectedContexts": [
			{"name": "minikube", "connectionId": "aaaaaaaa-1111-1111-1111-111111111111"}
		],
		"registeredContexts": [
			{"name": "prod-cluster", "connectionId": "bbbbbbbb-2222-2222-2222-222222222222"}
		]
	}`

	var response saveK8sContextResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		t.Fatalf("failed to unmarshal server response: %v", err)
	}

	if got, want := connectionIDForContext(response, "minikube"), "aaaaaaaa-1111-1111-1111-111111111111"; got != want {
		t.Fatalf("connectionIDForContext(minikube) = %q, want %q", got, want)
	}
	if got, want := connectionIDForContext(response, "prod-cluster"), "bbbbbbbb-2222-2222-2222-222222222222"; got != want {
		t.Fatalf("connectionIDForContext(prod-cluster) = %q, want %q", got, want)
	}
}

// TestCreateConnectionFlags guards the flag registrations on `connection create`.
func TestCreateConnectionFlags(t *testing.T) {
	fileFlag := createConnectionCmd.Flags().Lookup("file")
	if fileFlag == nil {
		t.Fatal("expected --file flag to be registered on createConnectionCmd")
	}
	if fileFlag.Shorthand != "f" {
		t.Errorf("expected shorthand 'f' for --file, got %q", fileFlag.Shorthand)
	}

	contextFlagObj := createConnectionCmd.Flags().Lookup("context")
	if contextFlagObj == nil {
		t.Fatal("expected --context flag to be registered on createConnectionCmd")
	}
	if contextFlagObj.Shorthand != "c" {
		t.Errorf("expected shorthand 'c' for --context, got %q", contextFlagObj.Shorthand)
	}

	typeFlag := createConnectionCmd.Flags().Lookup("type")
	if typeFlag == nil {
		t.Fatal("expected --type flag to be registered on createConnectionCmd")
	}
	if typeFlag.Shorthand != "t" {
		t.Errorf("expected shorthand 't' for --type, got %q", typeFlag.Shorthand)
	}
}

// TestCreateConnectionArgsValidation exercises validation for connection create arguments and flags.
func TestCreateConnectionArgsValidation(t *testing.T) {
	// Preserve global state for cleanup
	origType := connectionType
	origPath := kubeconfigPath
	defer func() {
		connectionType = origType
		kubeconfigPath = origPath
	}()

	tests := []struct {
		name       string
		connType   string
		filePath   string
		wantErr    bool
		errContain string
	}{
		{
			name:       "neither file nor type provided errors",
			connType:   "",
			filePath:   "",
			wantErr:    true,
			errContain: "either --file or --type is required",
		},
		{
			name:       "unsupported type provided errors",
			connType:   "invalid-type",
			filePath:   "",
			wantErr:    true,
			errContain: "Invalid connection type",
		},
		{
			name:       "file combined with cloud provider type errors",
			connType:   "aks",
			filePath:   "/path/to/kubeconfig",
			wantErr:    true,
			errContain: "--file cannot be used with --type aks",
		},
		{
			name:     "valid file without type succeeds",
			connType: "",
			filePath: "/path/to/kubeconfig",
			wantErr:  false,
		},
		{
			name:     "valid type minikube succeeds",
			connType: "minikube",
			filePath: "",
			wantErr:  false,
		},
		{
			name:     "valid type kubernetes succeeds",
			connType: "kubernetes",
			filePath: "",
			wantErr:  false,
		},
		{
			name:     "valid type kubernetes with file succeeds",
			connType: "kubernetes",
			filePath: "/path/to/kubeconfig",
			wantErr:  false,
		},
		{
			name:     "valid cloud type eks succeeds",
			connType: "eks",
			filePath: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connectionType = tt.connType
			kubeconfigPath = tt.filePath
			err := createConnectionCmd.Args(createConnectionCmd, []string{})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				if !strings.Contains(err.Error(), tt.errContain) {
					t.Errorf("error %q does not contain %q", err.Error(), tt.errContain)
				}
			} else if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

// TestSetContextSendsSelectedContextsAndContextName verifies that setContext transmits both
// selectedContexts (JSON array required by server) and contextName (for backward compatibility),
// and parses the resulting connection ID.
func TestSetContextSendsSelectedContextsAndContextName(t *testing.T) {
	var receivedContextName string
	var receivedSelectedContexts string
	var receivedFileName string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/system/kubernetes" {
			http.NotFound(w, r)
			return
		}
		err := r.ParseMultipartForm(1 << 20)
		if err != nil {
			t.Errorf("failed to parse multipart form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		receivedContextName = r.FormValue("contextName")
		receivedSelectedContexts = r.FormValue("selectedContexts")

		_, header, err := r.FormFile("k8sfile")
		if err == nil && header != nil {
			receivedFileName = header.Filename
		}

		resp := saveK8sContextResponse{
			RegisteredContexts: []registeredK8sContext{
				{Name: "cluster-beta", ConnectionID: "conn-123-abc"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	utils.SetupContextEnv(t)
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatalf("failed to get mesheryctl config: %v", err)
	}
	currCtx, err := mctlCfg.GetCurrentContext()
	if err != nil {
		t.Fatalf("failed to get current context: %v", err)
	}
	currCtx.Endpoint = ts.URL
	if err := config.UpdateContextInConfig(currCtx, mctlCfg.CurrentContext); err != nil {
		t.Fatalf("failed to update context in config: %v", err)
	}

	tmpDir := t.TempDir()
	tokenFile := filepath.Join(tmpDir, "auth.json")
	if err := os.WriteFile(tokenFile, []byte(`{"token": "test-token", "meshery-provider": "Local"}`), 0600); err != nil {
		t.Fatalf("failed to write dummy token file: %v", err)
	}
	origTokenFlag := utils.TokenFlag
	utils.TokenFlag = tokenFile
	defer func() { utils.TokenFlag = origTokenFlag }()

	dummyConfigFile := filepath.Join(tmpDir, "test-kubeconfig.yaml")
	if err := os.WriteFile(dummyConfigFile, []byte("apiVersion: v1\nkind: Config\n"), 0644); err != nil {
		t.Fatalf("failed to write dummy config file: %v", err)
	}

	connID, err := setContext(dummyConfigFile, "cluster-beta")
	if err != nil {
		t.Fatalf("setContext failed: %v", err)
	}

	if connID != "conn-123-abc" {
		t.Errorf("setContext returned connectionID %q, want %q", connID, "conn-123-abc")
	}
	if receivedContextName != "cluster-beta" {
		t.Errorf("received contextName = %q, want %q", receivedContextName, "cluster-beta")
	}
	if receivedSelectedContexts != `["cluster-beta"]` {
		t.Errorf("received selectedContexts = %q, want %q", receivedSelectedContexts, `["cluster-beta"]`)
	}
	if receivedFileName != "test-kubeconfig.yaml" {
		t.Errorf("received file = %q, want %q", receivedFileName, "test-kubeconfig.yaml")
	}
}

// TestCreateKubeconfigConnection_EdgeCases verifies edge case error handling for generic kubeconfigs.
func TestCreateKubeconfigConnection_EdgeCases(t *testing.T) {
	t.Run("nonexistent file returns read error", func(t *testing.T) {
		err := createKubeconfigConnection("/nonexistent/path/to/kubeconfig")
		if err == nil {
			t.Fatal("expected error for nonexistent file, got nil")
		}
		if !strings.Contains(err.Error(), "Unable to read kubeconfig file") {
			t.Errorf("error %q does not contain expected message", err.Error())
		}
	})

	t.Run("malformed file returns read error", func(t *testing.T) {
		tmpDir := t.TempDir()
		malformedPath := filepath.Join(tmpDir, "malformed-kubeconfig.yaml")
		if err := os.WriteFile(malformedPath, []byte("this is not yaml: ["), 0644); err != nil {
			t.Fatalf("failed to write malformed file: %v", err)
		}
		err := createKubeconfigConnection(malformedPath)
		if err == nil {
			t.Fatal("expected error for malformed file, got nil")
		}
		if !strings.Contains(err.Error(), "Unable to read kubeconfig file") {
			t.Errorf("error %q does not contain expected message", err.Error())
		}
	})

	t.Run("kubeconfig with zero contexts returns error", func(t *testing.T) {
		tmpDir := t.TempDir()
		emptyCtxPath := filepath.Join(tmpDir, "zero-ctx-kubeconfig.yaml")
		emptyConfig := "apiVersion: v1\nkind: Config\nclusters: []\ncontexts: []\nusers: []\n"
		if err := os.WriteFile(emptyCtxPath, []byte(emptyConfig), 0644); err != nil {
			t.Fatalf("failed to write empty contexts config: %v", err)
		}
		err := createKubeconfigConnection(emptyCtxPath)
		if err == nil {
			t.Fatal("expected error for zero contexts kubeconfig, got nil")
		}
		if !strings.Contains(err.Error(), "no contexts found in") {
			t.Errorf("error %q does not contain expected 'no contexts found'", err.Error())
		}
	})

	t.Run("requested nonexistent context returns error before writing config", func(t *testing.T) {
		tmpDir := t.TempDir()
		validPath := filepath.Join(tmpDir, "valid-kubeconfig.yaml")
		validConfig := "apiVersion: v1\nkind: Config\nclusters:\n- cluster:\n    server: https://127.0.0.1:6443\n  name: c1\ncontexts:\n- context:\n    cluster: c1\n    user: u1\n  name: ctx-actual\nusers:\n- name: u1\n  user:\n    token: t1\n"
		if err := os.WriteFile(validPath, []byte(validConfig), 0644); err != nil {
			t.Fatalf("failed to write valid config: %v", err)
		}

		// Ensure utils.ConfigPath points to a known existing file that must NOT be overwritten
		existingConfig := filepath.Join(tmpDir, "existing-config.yaml")
		if err := os.WriteFile(existingConfig, []byte("sentinel-content"), 0600); err != nil {
			t.Fatalf("failed to write existing config: %v", err)
		}
		origConfigPath := utils.ConfigPath
		utils.ConfigPath = existingConfig
		defer func() { utils.ConfigPath = origConfigPath }()

		origContextFlag := contextFlag
		contextFlag = "ctx-nonexistent"
		defer func() { contextFlag = origContextFlag }()

		err := createKubeconfigConnection(validPath)
		if err == nil {
			t.Fatal("expected error for nonexistent requested context, got nil")
		}
		if !strings.Contains(err.Error(), `context "ctx-nonexistent" not found in kubeconfig`) {
			t.Errorf("error %q does not contain expected context not found message", err.Error())
		}

		// Assert sentinel file was preserved and not overwritten
		content, err := os.ReadFile(existingConfig)
		if err != nil {
			t.Fatalf("failed to read existing config: %v", err)
		}
		if string(content) != "sentinel-content" {
			t.Errorf("existing config was overwritten: %q", string(content))
		}
	})

	t.Run("writeKubeconfigSafely creates file with 0600 permissions atomically", func(t *testing.T) {
		tmpDir := t.TempDir()
		destFile := filepath.Join(tmpDir, "out-config.yaml")
		cfg := clientcmdapi.NewConfig()
		cfg.Clusters["test"] = &clientcmdapi.Cluster{Server: "https://localhost:6443"}
		cfg.Contexts["test"] = &clientcmdapi.Context{Cluster: "test"}
		cfg.CurrentContext = "test"

		err := writeKubeconfigSafely(cfg, destFile)
		if err != nil {
			t.Fatalf("writeKubeconfigSafely failed: %v", err)
		}

		info, err := os.Stat(destFile)
		if err != nil {
			t.Fatalf("failed to stat destination file: %v", err)
		}
		if perm := info.Mode().Perm(); perm != 0600 {
			t.Errorf("expected file permissions 0600, got %o", perm)
		}

		// Verify no leftover temp files exist in tmpDir
		entries, err := os.ReadDir(tmpDir)
		if err != nil {
			t.Fatalf("failed to read dir: %v", err)
		}
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), "kubeconfig-") && strings.HasSuffix(e.Name(), ".tmp") {
				t.Errorf("leftover temporary file found: %s", e.Name())
			}
		}
	})
}
