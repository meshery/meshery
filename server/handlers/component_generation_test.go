package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/meshery/meshery/server/models"
)

func TestMeshModelGenerationHandler_PayloadTooLarge(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	oversizedBody := bytes.Repeat([]byte("a"), maxComponentGenerationPayloadSize+10)
	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestMeshModelGenerationHandler_ExceedBoundaryByOneByte(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Body of exactly maxComponentGenerationPayloadSize + 1 byte must be rejected with 413
	oversizedBody := bytes.Repeat([]byte("a"), maxComponentGenerationPayloadSize+1)
	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413 Request Entity Too Large at limit+1, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestMeshModelGenerationHandler_ExactBoundaryPayload(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	// Construct a payload padded to exactly maxComponentGenerationPayloadSize with whitespace
	raw := `{"data":[]}`
	if len(raw) > maxComponentGenerationPayloadSize {
		t.Fatalf("raw payload exceeds limit")
	}
	padding := bytes.Repeat([]byte(" "), maxComponentGenerationPayloadSize-len(raw))
	exactBoundaryBody := append([]byte(raw), padding...)

	if len(exactBoundaryBody) != maxComponentGenerationPayloadSize {
		t.Fatalf("expected exact boundary size %d, got %d", maxComponentGenerationPayloadSize, len(exactBoundaryBody))
	}

	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", bytes.NewReader(exactBoundaryBody))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	// MaxBytesReader allows exactly maxComponentGenerationPayloadSize bytes, so it must NOT return 413
	if rec.Code == http.StatusRequestEntityTooLarge {
		t.Fatalf("expected payload at exactly maxComponentGenerationPayloadSize to not return 413, got %d", rec.Code)
	}
	// It parses the body successfully, and then rejects the empty data array with 400 Bad Request
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty data array at boundary, got %d (body: %s)", rec.Code, rec.Body.String())
	}
}

func TestMeshModelGenerationHandler_EmptyAndWhitespaceBody(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

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
			req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			h.MeshModelGenerationHandler(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400 Bad Request, got %d (body: %s)", rec.Code, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
				t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
			}
		})
	}
}

func TestMeshModelGenerationHandler_MalformedJSON(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	malformedJSON := `{"data": [{"name": "test"`
	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", strings.NewReader(malformedJSON))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for malformed JSON, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestMeshModelGenerationHandler_EmptyDataArray(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")

	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", strings.NewReader(`{"data": []}`))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400 Bad Request for empty data array, got %d (body: %s)", rec.Code, rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), ErrRequestBodyCode) {
		t.Errorf("expected response to contain %s, got: %s", ErrRequestBodyCode, rec.Body.String())
	}
}

func TestMeshModelGenerationHandler_NilLoggerSafety(t *testing.T) {
	// Ensure nil logger doesn't cause a panic on error paths
	h := &Handler{
		config: &models.HandlerConfig{
			Providers: map[string]models.Provider{},
		},
		log: nil,
	}

	oversizedBody := bytes.Repeat([]byte("a"), maxComponentGenerationPayloadSize+1)
	req := httptest.NewRequest(http.MethodPost, "/api/meshmodels/generate", bytes.NewReader(oversizedBody))
	rec := httptest.NewRecorder()

	h.MeshModelGenerationHandler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d", rec.Code)
	}
}
