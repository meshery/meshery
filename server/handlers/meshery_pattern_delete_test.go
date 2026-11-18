package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

type mockPatternDeleteProvider struct {
	*models.DefaultLocalProvider
	deleteFn func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error)
}

func (m *mockPatternDeleteProvider) DeleteMesheryPatterns(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
	if m.deleteFn != nil {
		return m.deleteFn(r, patterns)
	}
	return m.DefaultLocalProvider.DeleteMesheryPatterns(r, patterns)
}

func newTestDeletePatternProvider(t *testing.T) *mockPatternDeleteProvider {
	t.Helper()
	local := &models.DefaultLocalProvider{}
	local.Initialize()
	return &mockPatternDeleteProvider{
		DefaultLocalProvider: local,
	}
}

func TestDeleteMultiMesheryPatternsHandler_PayloadTooLarge(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newTestDeletePatternProvider(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	oversizedBody := bytes.Repeat([]byte("a"), maxPatternDeletePayloadSize+10)
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestDeleteMultiMesheryPatternsHandler_ExceedBoundaryByOneByte(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newTestDeletePatternProvider(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	oversizedBody := bytes.Repeat([]byte("a"), maxPatternDeletePayloadSize+1)
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large at limit+1, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestDeleteMultiMesheryPatternsHandler_ExactBoundaryPayload(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newTestDeletePatternProvider(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	raw := `{"patterns":[]}`
	if len(raw) > maxPatternDeletePayloadSize {
		t.Fatalf("raw payload exceeds limit")
	}
	padding := bytes.Repeat([]byte(" "), maxPatternDeletePayloadSize-len(raw))
	exactBoundaryBody := append([]byte(raw), padding...)

	if len(exactBoundaryBody) != maxPatternDeletePayloadSize {
		t.Fatalf("expected exact boundary size %d, got %d", maxPatternDeletePayloadSize, len(exactBoundaryBody))
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", bytes.NewReader(exactBoundaryBody))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code == http.StatusRequestEntityTooLarge {
		t.Fatalf("expected payload at exactly maxPatternDeletePayloadSize to not return 413, got %d", rec.Code)
	}
	// Rejected with 400 because empty patterns array is validated
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty patterns at boundary, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestDeleteMultiMesheryPatternsHandler_EmptyAndWhitespaceBody(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newTestDeletePatternProvider(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	tests := []struct {
		name string
		body string
	}{
		{name: "empty string", body: ""},
		{name: "spaces only", body: "    "},
		{name: "tabs and newlines", body: " \t\r\n "},
		{name: "null literal", body: "null"},
		{name: "null with whitespace", body: "  null  \n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
				t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
			}
		})
	}
}

func TestDeleteMultiMesheryPatternsHandler_MalformedJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	called := false
	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		called = true
		return []byte("{}"), nil
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	malformedJSON := `{"patterns": [{"id":`
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(malformedJSON))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for malformed JSON, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrDecodingCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrDecodingCode, rec.Body.String())
	}
	if called {
		t.Fatalf("provider.DeleteMesheryPatterns was invoked on malformed JSON (silent fallthrough bug detected)")
	}
}

func TestDeleteMultiMesheryPatternsHandler_EmptyPatternsArray(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	called := false
	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		called = true
		return []byte("{}"), nil
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(`{"patterns": []}`))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty patterns array, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
	if called {
		t.Fatalf("provider.DeleteMesheryPatterns was invoked on empty patterns array")
	}
}

func TestDeleteMultiMesheryPatternsHandler_InvalidPatternUUID(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	called := false
	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		called = true
		return []byte("{}"), nil
	}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	tests := []struct {
		name string
		body string
	}{
		{name: "not a uuid", body: `{"patterns": [{"id": "not-a-uuid"}]}`},
		{name: "nil uuid", body: `{"patterns": [{"id": "00000000-0000-0000-0000-000000000000"}]}`},
		{name: "empty id string", body: `{"patterns": [{"id": ""}]}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called = false
			req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request for invalid UUID, got %d (body: %s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), models.ErrInvalidUUIDValueCode) {
				t.Errorf("expected response to contain %s, got: %s", models.ErrInvalidUUIDValueCode, rec.Body.String())
			}
			if called {
				t.Fatalf("provider.DeleteMesheryPatterns was invoked on invalid UUID")
			}
		})
	}
}

func TestDeleteMultiMesheryPatternsHandler_UnauthenticatedUser(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	called := false
	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		called = true
		return []byte("{}"), nil
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(`{"patterns": [{"id": "11111111-1111-1111-1111-111111111111"}]}`))
	rec := httptest.NewRecorder()

	// user is nil
	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401 Unauthorized for nil user, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), "user unauthorized") {
		t.Errorf("expected body to contain 'user unauthorized', got: %s", rec.Body.String())
	}
	if called {
		t.Fatalf("provider.DeleteMesheryPatterns was invoked for unauthenticated user")
	}
}

func TestDeleteMultiMesheryPatternsHandler_Success(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	id1 := uuid.Must(uuid.NewV4()).String()
	id2 := uuid.Must(uuid.NewV4()).String()

	called := false
	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		called = true
		if len(patterns.Patterns) != 2 {
			t.Fatalf("expected 2 patterns in delete call, got %d", len(patterns.Patterns))
		}
		if patterns.Patterns[0].ID != id1 || patterns.Patterns[1].ID != id2 {
			t.Fatalf("pattern IDs mismatch: got %v, want [%s, %s]", patterns.Patterns, id1, id2)
		}
		return []byte(`{"status":"deleted","count":2}`), nil
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	payload := fmt.Sprintf(`{"patterns": [{"id": "%s"}, {"id": "%s"}]}`, id1, id2)
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !called {
		t.Fatalf("expected provider.DeleteMesheryPatterns to be invoked")
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
	}
	if rec.Body.String() != `{"status":"deleted","count":2}` {
		t.Errorf("unexpected body: %s", rec.Body.String())
	}
}

func TestDeleteMultiMesheryPatternsHandler_ProviderError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	id := uuid.Must(uuid.NewV4()).String()

	provider := newTestDeletePatternProvider(t)
	provider.deleteFn = func(r *http.Request, patterns models.MesheryPatternDeleteRequestBody) ([]byte, error) {
		return nil, errors.New("underlying db error")
	}

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	payload := fmt.Sprintf(`{"patterns": [{"id": "%s"}]}`, id)
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", strings.NewReader(payload))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500 Internal Server Error, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrDeletePatternCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrDeletePatternCode, rec.Body.String())
	}
}

func TestDeleteMultiMesheryPatternsHandler_NilLoggerSafety(t *testing.T) {
	h := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{},
		},
		log: nil,
	}

	provider := newTestDeletePatternProvider(t)
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	oversizedBody := bytes.Repeat([]byte("a"), maxPatternDeletePayloadSize+1)
	req := httptest.NewRequest(http.MethodDelete, "/api/patterns/delete", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.DeleteMultiMesheryPatternsHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", rec.Code)
	}
}
