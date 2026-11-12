package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/database"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
)

func newTestPreferenceFixture(t *testing.T, migrate bool) (*Handler, *models.DefaultLocalProvider) {
	t.Helper()

	db, err := database.New(database.Options{
		Engine:   database.SQLITE,
		Filename: filepath.Join(t.TempDir(), "meshery.db"),
	})
	if err != nil {
		t.Fatalf("open database: %v", err)
	}
	if migrate {
		if err := models.AutoMigrateSystemTables(&db); err != nil {
			t.Fatalf("AutoMigrateSystemTables: %v", err)
		}
	}

	systemID := uuid.Must(uuid.NewV4())
	h := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	provider := &models.DefaultLocalProvider{
		TestProfilesPersister: &models.TestProfilesPersister{DB: &db},
		GenericPersister:      &db,
	}
	return h, provider
}

type mockTestPreferenceProvider struct {
	models.DefaultLocalProvider
	storeErr  error
	fetchErr  error
	getObj    *perfprofile.PerformanceTestConfig
	getErr    error
	deleteErr error
}

func (m *mockTestPreferenceProvider) SMPTestConfigStore(_ *http.Request, _ *perfprofile.PerformanceTestConfig) (string, error) {
	return "", m.storeErr
}

func (m *mockTestPreferenceProvider) SMPTestConfigFetch(_ *http.Request, _, _, _, _ string) ([]byte, error) {
	return nil, m.fetchErr
}

func (m *mockTestPreferenceProvider) SMPTestConfigGet(_ *http.Request, _ string) (*perfprofile.PerformanceTestConfig, error) {
	return m.getObj, m.getErr
}

func (m *mockTestPreferenceProvider) SMPTestConfigDelete(_ *http.Request, _ string) error {
	return m.deleteErr
}

func TestUserTestPreferenceMethodNotAllowed(t *testing.T) {
	h, provider := newTestPreferenceFixture(t, true)
	req := httptest.NewRequest(http.MethodPut, "/api/user/prefs/perf", nil)
	rec := httptest.NewRecorder()

	h.UserTestPreferenceHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d. body: %s", rec.Code, http.StatusMethodNotAllowed, rec.Body.String())
	}
}

func TestUserTestPreferenceInvalidUUIDsReturnBadRequest(t *testing.T) {
	h, provider := newTestPreferenceFixture(t, true)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	for _, badUUID := range []string{"not-a-uuid", uuid.Nil.String()} {
		getReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid="+badUUID, nil)
		getRec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(getRec, getReq, nil, user, provider)
		if getRec.Code != http.StatusBadRequest {
			t.Errorf("GET uuid=%q status = %d, want %d. body: %s", badUUID, getRec.Code, http.StatusBadRequest, getRec.Body.String())
		}

		delReq := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid="+badUUID, nil)
		delRec := httptest.NewRecorder()
		h.UserTestPreferenceHandler(delRec, delReq, nil, user, provider)
		if delRec.Code != http.StatusBadRequest {
			t.Errorf("DELETE uuid=%q status = %d, want %d. body: %s", badUUID, delRec.Code, http.StatusBadRequest, delRec.Body.String())
		}
	}
}

func TestUserTestPreferenceLocalNotFoundReturns404(t *testing.T) {
	h, provider := newTestPreferenceFixture(t, true)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	missingID := uuid.Must(uuid.NewV4()).String()

	getReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid="+missingID, nil)
	getRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(getRec, getReq, nil, user, provider)
	if getRec.Code != http.StatusNotFound {
		t.Fatalf("GET missing status = %d, want %d. body: %s", getRec.Code, http.StatusNotFound, getRec.Body.String())
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid="+missingID, nil)
	delRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(delRec, delReq, nil, user, provider)
	if delRec.Code != http.StatusNotFound {
		t.Fatalf("DELETE missing status = %d, want %d. body: %s", delRec.Code, http.StatusNotFound, delRec.Body.String())
	}
}

func TestUserTestPreferenceLocalCRUDAndDefaultPagination(t *testing.T) {
	h, provider := newTestPreferenceFixture(t, true)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	payload := `{
		"name": "smoke-test",
		"duration": "10s",
		"clients": [{
			"loadGenerator": "fortio",
			"protocol": "http",
			"connections": 1,
			"rps": 5,
			"endpointUrls": ["http://localhost:8080"]
		}]
	}`
	postReq := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", bytes.NewBufferString(payload))
	postRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(postRec, postReq, nil, user, provider)
	if postRec.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want %d. body: %s", postRec.Code, http.StatusOK, postRec.Body.String())
	}

	var stored struct {
		TestUUID string `json:"test_uuid"`
	}
	if err := json.Unmarshal(postRec.Body.Bytes(), &stored); err != nil || stored.TestUUID == "" {
		t.Fatalf("unexpected POST response %s (err=%v)", postRec.Body.String(), err)
	}

	// Omitted page/pagesize must use defaults rather than 500.
	listReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf", nil)
	listRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(listRec, listReq, nil, user, provider)
	if listRec.Code != http.StatusOK {
		t.Fatalf("GET list status = %d, want %d. body: %s", listRec.Code, http.StatusOK, listRec.Body.String())
	}
	if ct := listRec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("GET list Content-Type = %q, want application/json", ct)
	}
	var page models.UserTestProfiles
	if err := json.Unmarshal(listRec.Body.Bytes(), &page); err != nil {
		t.Fatalf("unmarshal list page: %v", err)
	}
	if page.TotalCount != 1 || len(page.TestConfigs) != 1 {
		t.Fatalf("list totalCount=%d len=%d, want 1", page.TotalCount, len(page.TestConfigs))
	}

	// Bad pagination parameters must return 400 Bad Request.
	badPageReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?page=abc", nil)
	badPageRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(badPageRec, badPageReq, nil, user, provider)
	if badPageRec.Code != http.StatusBadRequest {
		t.Fatalf("GET invalid page status = %d, want %d. body: %s", badPageRec.Code, http.StatusBadRequest, badPageRec.Body.String())
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid="+stored.TestUUID, nil)
	getRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(getRec, getReq, nil, user, provider)
	if getRec.Code != http.StatusOK {
		t.Fatalf("GET single status = %d, want %d. body: %s", getRec.Code, http.StatusOK, getRec.Body.String())
	}
	if ct := getRec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("GET single Content-Type = %q, want application/json", ct)
	}

	delReq := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid="+stored.TestUUID, nil)
	delRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(delRec, delReq, nil, user, provider)
	if delRec.Code != http.StatusOK {
		t.Fatalf("DELETE status = %d, want %d. body: %s", delRec.Code, http.StatusOK, delRec.Body.String())
	}

	// Second delete and get must return 404.
	delAgainRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(delAgainRec, delReq, nil, user, provider)
	if delAgainRec.Code != http.StatusNotFound {
		t.Fatalf("second DELETE status = %d, want %d. body: %s", delAgainRec.Code, http.StatusNotFound, delAgainRec.Body.String())
	}

	getAfterRec := httptest.NewRecorder()
	h.UserTestPreferenceHandler(getAfterRec, getReq, nil, user, provider)
	if getAfterRec.Code != http.StatusNotFound {
		t.Fatalf("GET after DELETE status = %d, want %d. body: %s", getAfterRec.Code, http.StatusNotFound, getAfterRec.Body.String())
	}
}

func TestUserTestPreferenceHandlersPropagateProviderStatus(t *testing.T) {
	cases := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{name: "bad request", err: httputil.WithProviderStatus(errors.New("bad request"), http.StatusBadRequest), wantStatus: http.StatusBadRequest},
		{name: "unauthorized", err: httputil.WithProviderStatus(errors.New("unauthorized"), http.StatusUnauthorized), wantStatus: http.StatusUnauthorized},
		{name: "forbidden", err: httputil.WithProviderStatus(errors.New("forbidden"), http.StatusForbidden), wantStatus: http.StatusForbidden},
		{name: "not found", err: httputil.WithProviderStatus(errors.New("not found"), http.StatusNotFound), wantStatus: http.StatusNotFound},
		{name: "conflict", err: httputil.WithProviderStatus(errors.New("conflict"), http.StatusConflict), wantStatus: http.StatusConflict},
		{name: "untagged defaults to 500", err: errors.New("boom"), wantStatus: http.StatusInternalServerError},
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	validID := uuid.Must(uuid.NewV4()).String()
	validBody := `{"name":"remote-cfg","duration":"5s","clients":[{"loadGenerator":"fortio","protocol":"http","connections":1,"rps":1,"endpointUrls":["http://localhost:8080"]}]}`

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			systemID := uuid.Must(uuid.NewV4())
			h := &Handler{
				config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
				log:      newTestLogger(t),
				SystemID: &systemID,
			}
			provider := &mockTestPreferenceProvider{
				storeErr:  tc.err,
				fetchErr:  tc.err,
				getErr:    tc.err,
				deleteErr: tc.err,
			}

			postReq := httptest.NewRequest(http.MethodPost, "/api/user/prefs/perf", bytes.NewBufferString(validBody))
			postRec := httptest.NewRecorder()
			h.UserTestPreferenceStore(postRec, postReq, nil, user, provider)
			if postRec.Code != tc.wantStatus {
				t.Errorf("Store status = %d, want %d. body: %s", postRec.Code, tc.wantStatus, postRec.Body.String())
			}

			listReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf", nil)
			listRec := httptest.NewRecorder()
			h.UserTestPreferenceGet(listRec, listReq, nil, user, provider)
			if listRec.Code != tc.wantStatus {
				t.Errorf("Fetch status = %d, want %d. body: %s", listRec.Code, tc.wantStatus, listRec.Body.String())
			}

			getReq := httptest.NewRequest(http.MethodGet, "/api/user/prefs/perf?uuid="+validID, nil)
			getRec := httptest.NewRecorder()
			h.UserTestPreferenceGet(getRec, getReq, nil, user, provider)
			if getRec.Code != tc.wantStatus {
				t.Errorf("Get status = %d, want %d. body: %s", getRec.Code, tc.wantStatus, getRec.Body.String())
			}

			delReq := httptest.NewRequest(http.MethodDelete, "/api/user/prefs/perf?uuid="+validID, nil)
			delRec := httptest.NewRecorder()
			h.UserTestPreferenceDelete(delRec, delReq, nil, user, provider)
			if delRec.Code != tc.wantStatus {
				t.Errorf("Delete status = %d, want %d. body: %s", delRec.Code, tc.wantStatus, delRec.Body.String())
			}
		})
	}
}
