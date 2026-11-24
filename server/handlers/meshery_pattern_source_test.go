package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

type mockPatternSourceProvider struct {
	*models.DefaultLocalProvider
	token         string
	sourceContent []byte
	sourceErr     error
}

func newMockPatternSourceProvider(token string, content []byte, err error) *mockPatternSourceProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &mockPatternSourceProvider{
		DefaultLocalProvider: base,
		token:                token,
		sourceContent:        content,
		sourceErr:            err,
	}
}

func (m *mockPatternSourceProvider) GetProviderToken(r *http.Request) (string, error) {
	return m.token, nil
}

func (m *mockPatternSourceProvider) GetDesignSourceContent(token, designID string) ([]byte, error) {
	if m.sourceErr != nil {
		return nil, m.sourceErr
	}
	return m.sourceContent, nil
}

func TestDownloadMesheryPatternHandler_InvalidPathUUID_ReturnsBadRequest(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "empty id", id: ""},
		{name: "malformed id", id: "not-a-uuid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

			req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+tc.id, nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id})
			rec := httptest.NewRecorder()

			h.DownloadMesheryPatternHandler(rec, req, nil, authUser, newMockPatternSourceProvider("test-token", nil, nil))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
				t.Errorf("expected response to contain %q, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
			}
		})
	}
}

func TestGetMesheryPatternSourceHandler_InvalidPathUUID_ReturnsBadRequest(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{name: "empty id", id: ""},
		{name: "malformed id", id: "not-a-uuid"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

			req := httptest.NewRequest(http.MethodGet, "/api/pattern/download/"+tc.id+"/k8smanifest", nil)
			req = mux.SetURLVars(req, map[string]string{"id": tc.id, "sourcetype": "k8smanifest"})
			rec := httptest.NewRecorder()

			h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, newMockPatternSourceProvider("test-token", nil, nil))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
				t.Errorf("expected response to contain %q, got %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
			}
		})
	}
}

func TestGetMesheryPatternSourceHandler_Success_Yaml(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
	sampleYAML := []byte("apiVersion: v1\nkind: Pod\nmetadata:\n  name: test-pod\n")

	mockProvider := newMockPatternSourceProvider("test-token", sampleYAML, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/pattern/download/%s/k8smanifest", patID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "k8smanifest"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, mockProvider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/x-yaml" {
		t.Errorf("expected Content-Type application/x-yaml, got %s", contentType)
	}
	if !bytes.Equal(rec.Body.Bytes(), sampleYAML) {
		t.Errorf("expected body %q, got %q", string(sampleYAML), rec.Body.String())
	}
}

func TestGetMesheryPatternSourceHandler_Success_HelmChart(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}
	sampleTar := []byte("fake-tar-binary-content")

	mockProvider := newMockPatternSourceProvider("test-token", sampleTar, nil)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/pattern/download/%s/helmchart", patID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "helmchart"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, mockProvider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if contentType := rec.Header().Get("Content-Type"); contentType != "application/x-tar" {
		t.Errorf("expected Content-Type application/x-tar, got %s", contentType)
	}
	if !bytes.Equal(rec.Body.Bytes(), sampleTar) {
		t.Errorf("expected body %q, got %q", string(sampleTar), rec.Body.String())
	}
}

func TestGetMesheryPatternSourceHandler_NotFound_Returns404(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	notFoundErr := models.ErrResultNotFound(fmt.Errorf("design not found"))
	mockProvider := newMockPatternSourceProvider("test-token", nil, notFoundErr)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/pattern/download/%s/k8smanifest", patID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "k8smanifest"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, mockProvider)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404 Not Found, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestGetMesheryPatternSourceHandler_InternalServerError_Returns500(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	patID := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	dbErr := models.ErrDBRead(fmt.Errorf("database disk error"))
	mockProvider := newMockPatternSourceProvider("test-token", nil, dbErr)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/pattern/download/%s/k8smanifest", patID.String()), nil)
	req = mux.SetURLVars(req, map[string]string{"id": patID.String(), "sourcetype": "k8smanifest"})
	rec := httptest.NewRecorder()

	h.GetMesheryPatternSourceHandler(rec, req, nil, authUser, mockProvider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}
