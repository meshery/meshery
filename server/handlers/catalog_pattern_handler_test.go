package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
)

type catalogPatternMockProvider struct {
	*models.DefaultLocalProvider
	publishResp   []byte
	publishErr    error
	unpublishResp []byte
	unpublishErr  error
}

func newCatalogPatternMockProvider() *catalogPatternMockProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &catalogPatternMockProvider{
		DefaultLocalProvider: base,
	}
}

func (m *catalogPatternMockProvider) PublishCatalogPattern(_ *http.Request, _ *models.MesheryCatalogPatternRequestBody) ([]byte, error) {
	return m.publishResp, m.publishErr
}

func (m *catalogPatternMockProvider) UnPublishCatalogPattern(_ *http.Request, _ *models.MesheryCatalogPatternRequestBody) ([]byte, error) {
	return m.unpublishResp, m.unpublishErr
}

func (m *catalogPatternMockProvider) PersistEvent(_ events.Event, _ *string) error {
	return nil
}

func newCatalogPatternHandlerForTest(t *testing.T) (*Handler, *catalogPatternMockProvider, *models.User) {
	t.Helper()

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	systemID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create system id: %v", err)
	}
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create user id: %v", err)
	}

	provider := newCatalogPatternMockProvider()

	handler := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("events-test"),
			PatternChannel:   models.NewBroadcaster("pattern-test"),
		},
		log:      log,
		SystemID: &systemID,
	}

	return handler, provider, &models.User{ID: userID}
}

func catalogPatternRequestBody(t *testing.T) string {
	t.Helper()

	patternID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create pattern id: %v", err)
	}

	return fmt.Sprintf(`{"id":"%s","catalog_data":{"test":"value"}}`, patternID.String())
}

func TestPublishCatalogPatternHandler_UnmarshalErrorReturns500(t *testing.T) {
	handler, provider, user := newCatalogPatternHandlerForTest(t)
	provider.publishResp = []byte("{")

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/publish", strings.NewReader(catalogPatternRequestBody(t)))
	rec := httptest.NewRecorder()

	handler.PublishCatalogPatternHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error failed to publish catalog design") {
		t.Fatalf("expected response body to include publish error text, got %q", rec.Body.String())
	}
}

func TestPublishCatalogPatternHandler_NullResponseReturns500(t *testing.T) {
	handler, provider, user := newCatalogPatternHandlerForTest(t)
	provider.publishResp = []byte("null")

	req := httptest.NewRequest(http.MethodPost, "/api/pattern/catalog/publish", strings.NewReader(catalogPatternRequestBody(t)))
	rec := httptest.NewRecorder()

	handler.PublishCatalogPatternHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error failed to publish catalog design") {
		t.Fatalf("expected response body to include publish error text, got %q", rec.Body.String())
	}
}

func TestUnPublishCatalogPatternHandler_UnmarshalErrorReturns500(t *testing.T) {
	handler, provider, user := newCatalogPatternHandlerForTest(t)
	provider.unpublishResp = []byte("{")

	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/catalog/unpublish", strings.NewReader(catalogPatternRequestBody(t)))
	rec := httptest.NewRecorder()

	handler.UnPublishCatalogPatternHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error failed to unpublish catalog design") {
		t.Fatalf("expected response body to include unpublish error text, got %q", rec.Body.String())
	}
}

func TestUnPublishCatalogPatternHandler_NullResponseReturns500(t *testing.T) {
	handler, provider, user := newCatalogPatternHandlerForTest(t)
	provider.unpublishResp = []byte("null")

	req := httptest.NewRequest(http.MethodDelete, "/api/pattern/catalog/unpublish", strings.NewReader(catalogPatternRequestBody(t)))
	rec := httptest.NewRecorder()

	handler.UnPublishCatalogPatternHandler(rec, req, nil, user, provider)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "Error failed to unpublish catalog design") {
		t.Fatalf("expected response body to include unpublish error text, got %q", rec.Body.String())
	}
}
