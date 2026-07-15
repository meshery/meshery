package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

func TestInstallExtensionHandler_RemoteProviderReturnsNotImplemented(t *testing.T) {
	remote := &models.RemoteProvider{
		ProviderProperties: models.ProviderProperties{
			ProviderType: models.RemoteProviderType,
		},
	}
	h := newTestHandler(t, map[string]models.Provider{}, "")
	req := httptest.NewRequest(http.MethodPost, "/api/provider/extension/install", strings.NewReader(`{"extType":"navigator"}`))
	rec := httptest.NewRecorder()
	userID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate user id: %v", err)
	}

	h.InstallExtensionHandler(rec, req, nil, &models.User{ID: userID}, remote)

	if rec.Code != http.StatusNotImplemented {
		t.Fatalf("expected %d, got %d", http.StatusNotImplemented, rec.Code)
	}
}
