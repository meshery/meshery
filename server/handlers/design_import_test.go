package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

type designImportTestProvider struct {
	*models.DefaultLocalProvider
	savedPattern  *models.MesheryPattern
	savedResponse []byte
}

func (p *designImportTestProvider) GetProviderToken(*http.Request) (string, error) {
	return "test-token", nil
}

func (p *designImportTestProvider) SaveMesheryPattern(_ string, design *models.MesheryPattern) ([]byte, error) {
	p.savedPattern = design
	return p.savedResponse, nil
}

func (p *designImportTestProvider) PersistEvent(events.Event, string) error {
	return nil
}

func TestConvertFileToDesign_ValidMesheryPattern(t *testing.T) {
	log := newTestLogger(t)
	want := pattern.PatternFile{
		Name:          "valid-design",
		SchemaVersion: "designs.meshery.io/v1beta3",
	}
	source, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("failed to marshal input design: %v", err)
	}

	got, fileType, err := ConvertFileToDesign(FileToImport{
		Data:     source,
		FileName: "valid-design.json",
	}, nil, log)
	if err != nil {
		t.Fatalf("ConvertFileToDesign returned an unexpected error: %v", err)
	}
	if got.Name != want.Name {
		t.Errorf("design name = %q, want %q", got.Name, want.Name)
	}
	if got.SchemaVersion != want.SchemaVersion {
		t.Errorf("schema version = %q, want %q", got.SchemaVersion, want.SchemaVersion)
	}
	if fileType != core.MesheryDesign {
		t.Errorf("file type = %q, want %q", fileType, core.MesheryDesign)
	}
}

func TestMesheryDesignFromParsedFile_UnexpectedType(t *testing.T) {
	design, err := mesheryDesignFromParsedFile("not a pattern.PatternFile")
	if err == nil {
		t.Fatal("expected an error for an unexpected parsed-file type")
	}
	if design.Name != "" || design.SchemaVersion != "" {
		t.Errorf("design = %#v, want the zero value", design)
	}

	wrapped := ErrConvertToDesign(err)
	if got, want := meshkiterrors.GetCode(wrapped), ErrConvertToDesignCode; got != want {
		t.Errorf("structured error code = %q, want %q", got, want)
	}
}

func TestDesignFileImportHandler_SuccessSetsJSONContentType(t *testing.T) {
	source, err := json.Marshal(pattern.PatternFile{
		Name:          "imported-design",
		SchemaVersion: "designs.meshery.io/v1beta3",
	})
	if err != nil {
		t.Fatalf("failed to marshal input design: %v", err)
	}
	requestBody, err := json.Marshal(pattern.MesheryPatternImportFilePayload{
		File:     source,
		FileName: "imported-design.json",
	})
	if err != nil {
		t.Fatalf("failed to marshal import request: %v", err)
	}

	provider := &designImportTestProvider{
		DefaultLocalProvider: &models.DefaultLocalProvider{},
		savedResponse:        []byte(`{"saved":true}`),
	}
	systemID := uuid.Must(uuid.NewV4())
	handler := &Handler{
		config:   &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:      newTestLogger(t),
		SystemID: &systemID,
	}
	req := httptest.NewRequest(http.MethodPost, "/api/pattern/import", bytes.NewReader(requestBody))
	rec := httptest.NewRecorder()

	handler.DesignFileImportHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}
	if got := rec.Body.String(); got != string(provider.savedResponse) {
		t.Errorf("response body = %q, want %q", got, provider.savedResponse)
	}
	if provider.savedPattern == nil {
		t.Fatal("provider did not receive a pattern to save")
	}
	if provider.savedPattern.Name != "imported-design" {
		t.Errorf("saved design name = %q, want imported-design", provider.savedPattern.Name)
	}
	if !bytes.Equal(provider.savedPattern.SourceContent, source) {
		t.Errorf("saved source content = %q, want %q", provider.savedPattern.SourceContent, source)
	}
}

func TestDesignFileImportHandler_PreservesStructuredConversionErrors(t *testing.T) {
	tests := []struct {
		name        string
		body        string
		wantErrCode string
	}{
		{
			name:        "missing import variant",
			body:        `{}`,
			wantErrCode: ErrInvalidImportRequestCode,
		},
		{
			name:        "unparseable design file",
			body:        `{"file":"bm90IGEgZGVzaWdu","fileName":"broken.yaml"}`,
			wantErrCode: ErrConvertToDesignCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := &designImportTestProvider{DefaultLocalProvider: &models.DefaultLocalProvider{}}
			systemID := uuid.Must(uuid.NewV4())
			handler := &Handler{
				config:   &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
				log:      newTestLogger(t),
				SystemID: &systemID,
			}
			req := httptest.NewRequest(http.MethodPost, "/api/pattern/import", strings.NewReader(tt.body))
			rec := httptest.NewRecorder()

			handler.DesignFileImportHandler(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, provider)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if !strings.Contains(rec.Body.String(), tt.wantErrCode) {
				t.Errorf("response body %q does not contain structured error code %q", rec.Body.String(), tt.wantErrCode)
			}
		})
	}
}
