package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

// credentialSpyProvider captures the Credential passed to SaveUserCredential
// and UpdateUserCredential so tests can assert that the authenticated
// user's ID replaces any client-supplied `userId` in the request body.
type credentialSpyProvider struct {
	*models.DefaultLocalProvider
	observedSave   atomic.Pointer[models.Credential]
	observedUpdate atomic.Pointer[models.Credential]
}

func newCredentialSpyProvider() *credentialSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &credentialSpyProvider{DefaultLocalProvider: base}
}

func (m *credentialSpyProvider) SaveUserCredential(_ string, c *models.Credential) (*models.Credential, error) {
	cp := *c
	m.observedSave.Store(&cp)
	return c, nil
}

func (m *credentialSpyProvider) UpdateUserCredential(_ *http.Request, c *models.Credential) (*models.Credential, error) {
	cp := *c
	m.observedUpdate.Store(&cp)
	return c, nil
}

// TestSaveUserCredential_ClientSuppliedUserIdCannotOverride verifies the
// authorization-hardening fix: a client posting a body with a `userId`
// pointing at some other user must NOT be able to create a credential
// under that other user's account. The handler binds `credential.UserId`
// to the authenticated user's ID after unmarshal precisely so this
// overwrite attempt fails safely.
func TestSaveUserCredential_ClientSuppliedUserIdCannotOverride(t *testing.T) {
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}
	attacker := uuid.Must(uuid.FromString("22222222-2222-2222-2222-222222222222"))
	body := fmt.Sprintf(`{"userId":%q,"name":"poisoned","type":"token","secret":{"k":"v"}}`, attacker)

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/credentials", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.SaveUserCredential(rec, req, nil, authUser, p)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	saved := p.observedSave.Load()
	if saved == nil {
		t.Fatalf("provider SaveUserCredential was not invoked")
	}
	if saved.UserId != authUser.ID {
		t.Fatalf("credential.UserId was %v, want authenticated user %v — client-supplied userId leaked through", saved.UserId, authUser.ID)
	}
	if saved.UserId == attacker {
		t.Fatalf("credential.UserId matches attacker-supplied value %v — authorization bypass", attacker)
	}
}

// TestUpdateUserCredential_ClientSuppliedUserIdCannotOverride verifies
// the same hardening on the update path: a client cannot redirect an
// update onto another user's credential by supplying a foreign
// `userId` in the body.
func TestUpdateUserCredential_ClientSuppliedUserIdCannotOverride(t *testing.T) {
	authUser := &models.User{ID: uuid.Must(uuid.FromString("33333333-3333-3333-3333-333333333333"))}
	attacker := uuid.Must(uuid.FromString("44444444-4444-4444-4444-444444444444"))
	body := fmt.Sprintf(`{"userId":%q,"id":"abcdefab-abcd-abcd-abcd-abcdefabcdef","name":"renamed","type":"token","secret":{}}`, attacker)

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.UpdateUserCredential(rec, req, nil, authUser, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	updated := p.observedUpdate.Load()
	if updated == nil {
		t.Fatalf("provider UpdateUserCredential was not invoked")
	}
	if updated.UserId != authUser.ID {
		t.Fatalf("credential.UserId was %v, want authenticated user %v — client-supplied userId leaked through", updated.UserId, authUser.ID)
	}
	if updated.UserId == attacker {
		t.Fatalf("credential.UserId matches attacker-supplied value %v — authorization bypass", attacker)
	}
}
