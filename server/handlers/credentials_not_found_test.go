package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models"
)

// credentialErrProvider fails credential updates and deletes with err.
type credentialErrProvider struct {
	*models.DefaultLocalProvider
	err error
}

func (p *credentialErrProvider) UpdateUserCredential(_ *http.Request, _ *models.Credential) (*models.Credential, error) {
	return nil, p.err
}

func (p *credentialErrProvider) DeleteUserCredential(_ *http.Request, _ uuid.UUID) (*models.Credential, error) {
	return nil, p.err
}

func newCredentialErrProvider(err error) *credentialErrProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &credentialErrProvider{DefaultLocalProvider: base, err: err}
}

var credentialErrCases = []struct {
	name     string
	err      error
	wantCode int
}{
	{"local provider not found", models.ErrResultNotFound(fmt.Errorf("credential not found")), http.StatusNotFound},
	{"remote provider 404", models.ErrFetch(fmt.Errorf("failed to update the credential"), "not found", http.StatusNotFound), http.StatusNotFound},
	{"other failure", fmt.Errorf("database is locked"), http.StatusInternalServerError},
}

func TestUpdateUserCredential_NotFoundIs404(t *testing.T) {
	for _, tc := range credentialErrCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			body := `{"id":"abcdefab-abcd-abcd-abcd-abcdefabcdef","name":"renamed","type":"token"}`
			req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials", bytes.NewBufferString(body))
			rec := httptest.NewRecorder()

			h.UpdateUserCredential(rec, req, nil, &models.User{ID: uuid.Must(uuid.NewV4())}, newCredentialErrProvider(tc.err))

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}
		})
	}
}

func TestDeleteUserCredential_NotFoundIs404(t *testing.T) {
	for _, tc := range credentialErrCases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials?credentialId="+uuid.Must(uuid.NewV4()).String(), nil)
			rec := httptest.NewRecorder()

			h.DeleteUserCredential(rec, req, nil, &models.User{}, newCredentialErrProvider(tc.err))

			if rec.Code != tc.wantCode {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantCode, rec.Body.String())
			}
		})
	}
}
