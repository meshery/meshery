package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
)

func newMeshOpsTestLocalProvider(t *testing.T) models.Provider {
	t.Helper()
	prefPersister, err := models.NewMapPreferencePersister()
	if err != nil {
		t.Fatalf("failed to create preference persister: %v", err)
	}
	return &models.DefaultLocalProvider{
		MapPreferencePersister: prefPersister,
	}
}

func newMeshOpsTestHandler(t *testing.T) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	return &Handler{
		config: &models.HandlerConfig{
			AdapterTracker: newMockAdaptersTracker(),
		},
		log: log,
	}
}

func TestMeshAdapterConfigHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/manage", nil)
	rec := httptest.NewRecorder()

	h.MeshAdapterConfigHandler(rec, req, &models.Preference{}, nil, newMeshOpsTestLocalProvider(t))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for nil user, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "user unauthorized") {
		t.Fatalf("expected response body to contain %q, got %q", "user unauthorized", rec.Body.String())
	}
}

func TestMeshAdapterConfigHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/manage", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	h.MeshAdapterConfigHandler(rec, req, &models.Preference{}, user, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for nil provider, got %d", rec.Code)
	}
}

func TestMeshAdapterConfigHandler_NilPrefObj_DoesNotPanic(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodDelete, "/api/system/adapter/manage?adapter=nonexistent:1000", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	// Should not panic even if prefObj is nil
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MeshAdapterConfigHandler panicked with nil prefObj: %v", r)
		}
	}()

	h.MeshAdapterConfigHandler(rec, req, nil, user, newMeshOpsTestLocalProvider(t))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for nonexistent adapter, got %d", rec.Code)
	}
}

func TestMeshOpsHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/operation", nil)
	rec := httptest.NewRecorder()

	h.MeshOpsHandler(rec, req, &models.Preference{}, nil, newMeshOpsTestLocalProvider(t))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for nil user, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "user unauthorized") {
		t.Fatalf("expected response body to contain %q, got %q", "user unauthorized", rec.Body.String())
	}
}

func TestMeshOpsHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/operation", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	h.MeshOpsHandler(rec, req, &models.Preference{}, user, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for nil provider, got %d", rec.Code)
	}
}

func TestMeshOpsHandler_NilPrefObj_DoesNotPanic(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	form := url.Values{}
	form.Set("adapter", "meshery-istio:10000")
	req := httptest.NewRequest(http.MethodPost, "/api/system/adapter/operation", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MeshOpsHandler panicked with nil prefObj: %v", r)
		}
	}()

	h.MeshOpsHandler(rec, req, nil, user, newMeshOpsTestLocalProvider(t))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request when adapter not found, got %d", rec.Code)
	}
}

func TestAdapterPingHandler_NilPrefObj_DoesNotPanic(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/adapter/manage?adapter=meshery-istio:10000", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AdapterPingHandler panicked with nil prefObj: %v", r)
		}
	}()

	h.AdapterPingHandler(rec, req, nil, nil, nil)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 for unconfigured adapter, got %d", rec.Code)
	}
}

func TestAdaptersHandler_NilConfig_DoesNotPanic(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	h := &Handler{
		config: nil,
		log:    log,
	}
	req := httptest.NewRequest(http.MethodGet, "/api/system/adapters", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AdaptersHandler panicked with nil config: %v", r)
		}
	}()

	h.AdaptersHandler(rec, req, nil, nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}
}

func TestSessionSyncHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)
	rec := httptest.NewRecorder()

	h.SessionSyncHandler(rec, req, &models.Preference{}, nil, newMeshOpsTestLocalProvider(t))

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for nil user, got %d", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "user unauthorized") {
		t.Fatalf("expected response body to contain %q, got %q", "user unauthorized", rec.Body.String())
	}
}

func TestSessionSyncHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	h.SessionSyncHandler(rec, req, &models.Preference{}, user, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 for nil provider, got %d", rec.Code)
	}
}

func TestSessionSyncHandler_NilPrefObj_DoesNotPanic(t *testing.T) {
	h := newMeshOpsTestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SessionSyncHandler panicked with nil prefObj: %v", r)
		}
	}()

	h.SessionSyncHandler(rec, req, nil, user, newMeshOpsTestLocalProvider(t))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}
}

func TestSessionSyncHandler_NilConfig_DoesNotPanic(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	h := &Handler{
		config: nil,
		log:    log,
	}
	req := httptest.NewRequest(http.MethodGet, "/api/system/sync", nil)
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	existingAdapter := &models.Adapter{Location: "test-adapter:1000", Name: "test"}
	pref := &models.Preference{
		MeshAdapters: []*models.Adapter{existingAdapter},
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("SessionSyncHandler panicked with nil config: %v", r)
		}
	}()

	h.SessionSyncHandler(rec, req, pref, user, newMeshOpsTestLocalProvider(t))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", rec.Code)
	}
	if len(pref.MeshAdapters) != 1 || pref.MeshAdapters[0].Location != "test-adapter:1000" {
		t.Fatalf("expected existing adapters to be preserved when config is nil, got %+v", pref.MeshAdapters)
	}
}
