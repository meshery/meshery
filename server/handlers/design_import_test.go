package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

func TestFetchFileFromURL_RejectsSSRF(t *testing.T) {
	restrictedURLs := []struct {
		name string
		url  string
	}{
		{"IPv4 Loopback", "http://127.0.0.1:9081/api"},
		{"Localhost hostname", "http://localhost:9081/api"},
		{"IPv6 Loopback", "http://[::1]:9081/api"},
		{"AWS/GCP IMDS", "http://169.254.169.254/latest/meta-data/"},
		{"IPv4-mapped IPv6 IMDS", "http://[::ffff:169.254.169.254]/latest/meta-data/"},
		{"GCP metadata hostname", "http://metadata.google.internal/computeMetadata/v1/"},
		{"AWS instance-data hostname", "http://instance-data/latest/meta-data/"},
		{"Private Class A", "http://10.0.0.1/design.yaml"},
		{"Private Class B", "http://172.16.0.1/design.yaml"},
		{"Private Class C", "http://192.168.1.1/design.yaml"},
		{"File scheme", "file:///etc/passwd"},
		{"FTP scheme", "ftp://example.com/design.yaml"},
		{"Empty host", "http:///design.yaml"},
	}

	for _, tc := range restrictedURLs {
		t.Run(tc.name, func(t *testing.T) {
			_, err := fetchFileFromURL(tc.url)
			if err == nil {
				t.Fatalf("fetchFileFromURL(%q) should have failed with SSRF error, but succeeded", tc.url)
			}

			// Ensure error code is models.ErrDoRequestCode
			code := errors.GetCode(err)
			if code != models.ErrDoRequestCode {
				t.Errorf("expected error code %q, got %q (err: %v)", models.ErrDoRequestCode, code, err)
			}
		})
	}
}

func TestResolveImportVariant_URLImportSSRF(t *testing.T) {
	const restrictedURL = "http://169.254.169.254/latest/meta-data/"
	designName := "Malicious Design"

	body := pattern.MesheryPatternImportRequestBody{}
	if err := body.FromMesheryPatternImportURLPayload(pattern.MesheryPatternImportURLPayload{
		Name: &designName,
		Url:  restrictedURL,
	}); err != nil {
		t.Fatalf("failed to construct import request body: %v", err)
	}

	variant, err := resolveImportVariant(body)
	if err == nil {
		t.Fatal("expected resolveImportVariant to reject restricted URL, but it succeeded")
	}

	// Verify the variant still carries the URL so the downstream event logs the source
	if variant.File.FileName != restrictedURL {
		t.Errorf("expected variant.File.FileName to preserve %q, got %q", restrictedURL, variant.File.FileName)
	}
	if variant.Name != designName {
		t.Errorf("expected variant.Name to preserve %q, got %q", designName, variant.Name)
	}

	if !strings.Contains(err.Error(), "restricted") && !strings.Contains(err.Error(), "forbidden") && !strings.Contains(err.Error(), "unable to make") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestDesignFileImportHandler_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/api/pattern/import", strings.NewReader("{}"))
	rec := httptest.NewRecorder()

	h.DesignFileImportHandler(rec, req, nil, nil, &models.DefaultLocalProvider{})

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
}

func TestDesignFileImportHandler_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := &Handler{}
	req := httptest.NewRequest(http.MethodPost, "/api/pattern/import", strings.NewReader("{}"))
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	h.DesignFileImportHandler(rec, req, nil, user, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
}

func TestDesignFileImportHandler_NilSystemID_DoesNotPanic(t *testing.T) {
	h := &Handler{
		SystemID: nil,
		config:   &models.HandlerConfig{EventBroadcaster: &models.Broadcast{}},
	}
	req := httptest.NewRequest(http.MethodPost, "/api/pattern/import", strings.NewReader("invalid json"))
	rec := httptest.NewRecorder()
	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	provider := &models.DefaultLocalProvider{}

	h.DesignFileImportHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 BadRequest for invalid body, got %d", rec.Code)
	}
}
