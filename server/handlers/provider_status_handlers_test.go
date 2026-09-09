package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
)

// keysOrgsFailingProvider scripts GetUsersKeys and GetOrganizations to fail so
// each handler's error path can be exercised in isolation. Like the sibling
// GetUserByID and GetEnvironments coverage, these pin the two rules that the
// misreporting bug violated: a provider error with no status of its own must
// fall back to 502 (never a fabricated 404), and a provider error that DOES
// carry a status must reach the client verbatim.
type keysOrgsFailingProvider struct {
	*models.DefaultLocalProvider
	err error
}

func newKeysOrgsFailingProvider(err error) *keysOrgsFailingProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &keysOrgsFailingProvider{DefaultLocalProvider: base, err: err}
}

func (p *keysOrgsFailingProvider) GetUsersKeys(_, _, _, _, _, _ string, _ string) ([]byte, error) {
	return nil, p.err
}

func (p *keysOrgsFailingProvider) GetOrganizations(_, _, _, _, _, _ string) ([]byte, error) {
	return nil, p.err
}

// withTokenContext attaches a token so the handlers clear their token guard and
// reach the provider call under test.
func withTokenContext(req *http.Request) *http.Request {
	return req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
}

func TestGetUsersKeys_ProviderErrorReturnsFallbackStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newKeysOrgsFailingProvider(errors.New("remote blew up"))

	req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/users/keys", nil))
	rec := httptest.NewRecorder()

	h.GetUsersKeys(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadGateway, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetUsersKeysCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetUsersKeysCode)
	}
}

func TestGetUsersKeys_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newKeysOrgsFailingProvider(models.ErrFetch(errors.New("forbidden"), "Keys", http.StatusForbidden))

	req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/users/keys", nil))
	rec := httptest.NewRecorder()

	h.GetUsersKeys(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetUsersKeysCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetUsersKeysCode)
	}
}

func TestGetOrganizations_ProviderErrorReturnsFallbackStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newKeysOrgsFailingProvider(errors.New("remote blew up"))

	req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil))
	rec := httptest.NewRecorder()

	h.GetOrganizations(rec, req, nil, nil, provider)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusBadGateway, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetOrganizationsCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetOrganizationsCode)
	}
}

func TestGetOrganizations_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newKeysOrgsFailingProvider(models.ErrFetch(errors.New("forbidden"), "Organizations", http.StatusForbidden))

	req := withTokenContext(httptest.NewRequest(http.MethodGet, "/api/identity/orgs", nil))
	rec := httptest.NewRecorder()

	h.GetOrganizations(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetOrganizationsCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetOrganizationsCode)
	}
}
