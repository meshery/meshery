package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
	meshkitEvents "github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/logger"
	eventstream "github.com/meshery/meshkit/utils/events"
)

type filterHandlerMockProvider struct {
	*models.DefaultLocalProvider
	persistedEvents int
}

func newFilterHandlerMockProvider() *filterHandlerMockProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()

	return &filterHandlerMockProvider{
		DefaultLocalProvider: base,
	}
}

func (m *filterHandlerMockProvider) PersistEvent(_ meshkitEvents.Event, _ *string) error {
	m.persistedEvents++
	return nil
}

func newFilterHandlerTestContext(t *testing.T) (*Handler, *filterHandlerMockProvider, *models.User) {
	t.Helper()

	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	systemID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create system UUID: %v", err)
	}

	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to create user UUID: %v", err)
	}

	handler := &Handler{
		config: &models.HandlerConfig{
			EventBroadcaster: models.NewBroadcaster("Events"),
		},
		log:          log,
		SystemID:     &systemID,
		EventsBuffer: eventstream.NewEventStreamer(),
	}

	provider := newFilterHandlerMockProvider()
	user := &models.User{ID: userID}

	return handler, provider, user
}

func TestHandleFilterPOST_MalformedRequestBody(t *testing.T) {
	h, provider, user := newFilterHandlerTestContext(t)

	req := httptest.NewRequest(http.MethodPost, "/api/filter", strings.NewReader("{"))
	resp := httptest.NewRecorder()

	h.handleFilterPOST(resp, req, nil, user, provider)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	if provider.persistedEvents != 1 {
		t.Fatalf("expected exactly one persisted event, got %d", provider.persistedEvents)
	}
}

func TestHandleFilterPOST_NullRequestBody(t *testing.T) {
	h, provider, user := newFilterHandlerTestContext(t)

	req := httptest.NewRequest(http.MethodPost, "/api/filter", strings.NewReader("null"))
	resp := httptest.NewRecorder()

	h.handleFilterPOST(resp, req, nil, user, provider)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}

	if provider.persistedEvents != 1 {
		t.Fatalf("expected exactly one persisted event, got %d", provider.persistedEvents)
	}
}
