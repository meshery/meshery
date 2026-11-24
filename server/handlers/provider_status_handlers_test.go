package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
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

type connectionsFailingProvider struct {
	*models.DefaultLocalProvider
	err error
}

func newConnectionsFailingProvider(err error) *connectionsFailingProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &connectionsFailingProvider{DefaultLocalProvider: base, err: err}
}

func (p *connectionsFailingProvider) GetConnections(_ *http.Request, _ string, _, _ int, _, _, _ string, _, _, _ []string, _ string) (*connections.ConnectionPage, error) {
	if p.err != nil {
		return nil, p.err
	}
	return &connections.ConnectionPage{}, nil
}

func TestGetConnectionsPreservesProviderStatus(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{
			name:       "tagged provider error",
			err:        models.ErrFetch(errors.New("service unavailable"), "Connections", http.StatusServiceUnavailable),
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "generic error",
			err:        errors.New("database unavailable"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "success",
			wantStatus: http.StatusOK,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newConnectionsFailingProvider(test.err)
			req := httptest.NewRequest(http.MethodGet, "/api/integrations/connections", nil)
			rec := httptest.NewRecorder()

			h.GetConnections(rec, req, nil, &models.User{}, provider)

			if rec.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, test.wantStatus, rec.Body.String())
			}
			if test.err != nil {
				if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetConnectionsCode {
					t.Errorf("code = %q, want %q", decoded.Code, ErrGetConnectionsCode)
				}
			}
		})
	}
}
