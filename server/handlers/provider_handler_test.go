package handlers

import (
	"encoding/json"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshsync/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestHandler returns a minimal Handler wired with the given providers and provider override.
func newTestHandler(t *testing.T, providers map[string]models.Provider, providerOverride string) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	return &Handler{
		config: &models.HandlerConfig{
			Providers:          providers,
			ProviderCookieName: "meshery-provider",
		},
		Provider: providerOverride,
		log:      log,
	}
}

func TestProviderHandler_NoneRedirectsThroughLogin(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/api/provider?provider=None", nil)
	rec := httptest.NewRecorder()

	h.ProviderHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/user/login?provider=None" {
		t.Errorf("expected redirect to /user/login?provider=None, got %s", loc)
	}

	// Verify provider cookie is set
	var found bool
	for _, c := range resp.Cookies() {
		if c.Name == "meshery-provider" && c.Value == "None" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected meshery-provider cookie to be set to None")
	}
}

func TestProviderUIHandler_NoneEnvVarRedirectsThroughLogin(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "None")

	req := httptest.NewRequest(http.MethodGet, "/provider", nil)
	rec := httptest.NewRecorder()

	h.ProviderUIHandler(rec, req)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/user/login" {
		t.Errorf("expected redirect to /user/login, got %s", loc)
	}
}

func TestLoginHandler_NoneProviderRedirectsHome(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None", nil)
	rec := httptest.NewRecorder()

	h.LoginHandler(rec, req, local, false)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/" {
		t.Errorf("expected redirect to /, got %s", loc)
	}
}

func TestLoginHandler_NoneProviderDeepLink(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	// ref is base64-raw-url-encoded "/extension/meshmap"
	req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None&ref=L2V4dGVuc2lvbi9tZXNobWFw", nil)
	rec := httptest.NewRecorder()

	h.LoginHandler(rec, req, local, false)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	if resp.StatusCode != http.StatusFound {
		t.Fatalf("expected 302 Found, got %d", resp.StatusCode)
	}

	loc := resp.Header.Get("Location")
	if loc != "/extension/meshmap" {
		t.Errorf("expected redirect to /extension/meshmap, got %s", loc)
	}
}

func TestLoginHandler_NoneProviderRejectsOpenRedirect(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	cases := []struct {
		name string
		ref  string
	}{
		{"absolute URL", base64.RawURLEncoding.EncodeToString([]byte("https://evil.example/phish"))},
		{"protocol-relative", base64.RawURLEncoding.EncodeToString([]byte("//evil.example/phish"))},
		{"javascript scheme", base64.RawURLEncoding.EncodeToString([]byte("javascript:alert(1)"))},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/user/login?provider=None&ref="+tc.ref, nil)
			rec := httptest.NewRecorder()

			h.LoginHandler(rec, req, local, false)

			resp := rec.Result()
			t.Cleanup(func() {
				if err := resp.Body.Close(); err != nil {
					t.Errorf("failed to close response body: %v", err)
				}
			})

			loc := resp.Header.Get("Location")
			if loc != "/" {
				t.Errorf("expected redirect to / for unsafe ref, got %s", loc)
			}
		})
	}
}

func TestLoginHandler_NoneProviderFromMiddleware(t *testing.T) {
	local := &models.DefaultLocalProvider{}
	local.Initialize()

	providers := map[string]models.Provider{local.Name(): local}
	h := newTestHandler(t, providers, "")

	req := httptest.NewRequest(http.MethodGet, "/api/some/endpoint", nil)
	rec := httptest.NewRecorder()

	// fromMiddleWare=true: simulates AuthMiddleware intercepting an API request
	h.LoginHandler(rec, req, local, true)

	resp := rec.Result()
	t.Cleanup(func() {
		if err := resp.Body.Close(); err != nil {
			t.Errorf("failed to close response body: %v", err)
		}
	})

	// Should be a no-op pass-through (200 with no body), not a redirect
	if resp.StatusCode == http.StatusFound {
		t.Fatalf("expected no redirect from middleware path, got 302 to %s", resp.Header.Get("Location"))
	}
}

func TestDeleteMeshSyncResource(t *testing.T) {
	handler := newTestHandler(t, map[string]models.Provider{}, "")

		tests := []struct {
			name                 string
			migrateResourceTable bool
			expectedStatus       int
			expectedContentType  string
			expectedBody         string
		}{
			{
				name: "given resource table migrated when DeleteMeshSyncResource then return status 200 and deleted true",
				migrateResourceTable: true,
				expectedStatus:       http.StatusOK,
				expectedContentType:  "application/json",
				expectedBody:         `"deleted":true`,
			},
			{
				name: "given resource table not migrated when DeleteMeshSyncResource then return status 500 and failed to delete",
				migrateResourceTable: false,
				expectedStatus:       http.StatusInternalServerError,
				expectedContentType:  "text/plain",
				expectedBody:         "Failed to Delete",
			},
		}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
			if err != nil {
				t.Fatalf("failed to create in-memory database: %v", err)
			}

			if tt.migrateResourceTable {
				if err := db.AutoMigrate(&model.KubernetesResource{}); err != nil {
					t.Fatalf("failed to migrate kubernetes resource table: %v", err)
				}
			}

			provider := &models.DefaultLocalProvider{}
			provider.Initialize()
			provider.GenericPersister = &database.Handler{DB: db}

			req := httptest.NewRequest(http.MethodDelete, "/api/system/meshsync/resources/resource-1", nil)
			req = mux.SetURLVars(req, map[string]string{"id": "resource-1"})
			rw := httptest.NewRecorder()

			handler.DeleteMeshSyncResource(rw, req, nil, nil, provider)

			if rw.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rw.Code)
			}

			if contentType := rw.Header().Get("Content-Type"); !strings.HasPrefix(contentType, tt.expectedContentType) {
				t.Fatalf("expected content type prefix %q, got %q", tt.expectedContentType, contentType)
			}

			if tt.expectedContentType == "application/json" {
				var response struct {
					Deleted bool `json:"deleted"`
				}

				if err := json.NewDecoder(rw.Body).Decode(&response); err != nil {
					t.Fatalf("failed to decode json response: %v", err)
				}

				if !response.Deleted {
					t.Fatal("expected 'deleted' field to be true")
				}
			} else {
				if body := rw.Body.String(); !strings.Contains(body, tt.expectedBody) {
					t.Fatalf("expected body to contain %q, got %q", tt.expectedBody, body)
				}
			}
		})
	}
}
