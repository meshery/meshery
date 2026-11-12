package models

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/schemas/models/v1beta1/workspace"
)

func TestRemoteProviderGetWorkspaces_ReturnsEmptyPageForNoContent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces" {
			http.NotFound(w, r)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistWorkspaces, Endpoint: "/workspaces"}}

	data, err := provider.GetWorkspaces("token", "2", "10", "", "", "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	var page workspace.WorkspacePage
	if err := json.Unmarshal(data, &page); err != nil {
		t.Fatalf("expected valid JSON workspace page, got %v", err)
	}

	if page.Page != 2 {
		t.Fatalf("expected page 2, got %d", page.Page)
	}
	if page.PageSize != 0 {
		t.Fatalf("expected page size 0, got %d", page.PageSize)
	}
	if page.TotalCount != 0 {
		t.Fatalf("expected total count 0, got %d", page.TotalCount)
	}
	if len(page.Workspaces) != 0 {
		t.Fatalf("expected no workspaces, got %d", len(page.Workspaces))
	}
}

func TestRemoteProviderGetWorkspaces_PreservesJSONBodyForOK(t *testing.T) {
	expected := `{"page":1,"page_size":1,"total_count":1,"workspaces":[{"name":"demo"}]}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/workspaces" {
			http.NotFound(w, r)
			return
		}

		_, _ = io.WriteString(w, expected)
	}))
	defer server.Close()

	provider := newTestRemoteProvider(t, server.URL)
	provider.Capabilities = Capabilities{{Feature: PersistWorkspaces, Endpoint: "/workspaces"}}

	data, err := provider.GetWorkspaces("token", "1", "10", "", "", "", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if string(data) != expected {
		t.Fatalf("expected response body %s, got %s", expected, string(data))
	}
}

func TestRemoteProviderGetViewsOfWorkspace_ReturnsErrorWhenProviderUnreachable(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://127.0.0.1:1")
	provider.Capabilities = Capabilities{{Feature: PersistWorkspaces, Endpoint: "/workspaces"}}

	req := httptest.NewRequest(http.MethodGet, "http://meshery.local/workspaces/workspace-id/views", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "token"})

	data, err := provider.GetViewsOfWorkspace(req, "workspace-id", "", "", "", "", "")
	if err == nil {
		t.Fatal("expected error when remote provider is unreachable")
	}
	if data != nil {
		t.Fatalf("expected nil data on unreachable provider, got %q", string(data))
	}
}

func TestRemoteProviderGetTeamsOfWorkspace_ReturnsErrorWhenProviderUnreachable(t *testing.T) {
	provider := newTestRemoteProvider(t, "http://127.0.0.1:1")
	provider.Capabilities = Capabilities{{Feature: PersistWorkspaces, Endpoint: "/workspaces"}}

	req := httptest.NewRequest(http.MethodGet, "http://meshery.local/workspaces/workspace-id/teams", nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: "token"})

	data, err := provider.GetTeamsOfWorkspace(req, "workspace-id", "", "", "", "", "")
	if err == nil {
		t.Fatal("expected error when remote provider is unreachable")
	}
	if data != nil {
		t.Fatalf("expected nil data on unreachable provider, got %q", string(data))
	}
}
