package connections

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestGetContexts_RequestError(t *testing.T) {
	utils.SetupContextEnv(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("POST", "http://localhost:9081/api/system/kubernetes/contexts",
		httpmock.NewErrorResponder(fmt.Errorf("connection refused")))

	tmpFile := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(tmpFile, []byte("apiVersion: v1\nkind: Config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := getContexts(tmpFile)
	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestGetContexts_NonOKStatus(t *testing.T) {
	utils.SetupContextEnv(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	httpmock.RegisterResponder("POST", "http://localhost:9081/api/system/kubernetes/contexts",
		httpmock.NewStringResponder(http.StatusInternalServerError, `internal error`))

	tmpFile := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(tmpFile, []byte("apiVersion: v1\nkind: Config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := getContexts(tmpFile)
	if err == nil {
		t.Fatal("expected error for 500 status but got nil")
	}
}

func TestGetContexts_Success(t *testing.T) {
	utils.SetupContextEnv(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	body := `[{"name": "context-1"}, {"name": "context-2"}]`
	httpmock.RegisterResponder("POST", "http://localhost:9081/api/system/kubernetes/contexts",
		httpmock.NewStringResponder(http.StatusOK, body))

	tmpFile := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(tmpFile, []byte("apiVersion: v1\nkind: Config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	names, err := getContexts(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 contexts, got %d", len(names))
	}
	if names[0] != "context-1" || names[1] != "context-2" {
		t.Fatalf("unexpected context names: %v", names)
	}
}

func TestGetContexts_InvalidContextName(t *testing.T) {
	utils.SetupContextEnv(t)
	httpmock.Activate()
	t.Cleanup(httpmock.DeactivateAndReset)

	body := `[{"name": 123}]`
	httpmock.RegisterResponder("POST", "http://localhost:9081/api/system/kubernetes/contexts",
		httpmock.NewStringResponder(http.StatusOK, body))

	tmpFile := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(tmpFile, []byte("apiVersion: v1\nkind: Config\n"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := getContexts(tmpFile)
	if err == nil {
		t.Fatal("expected error for invalid context name but got nil")
	}
}
