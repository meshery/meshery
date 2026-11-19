package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// credentialSpyProvider captures the Credential passed to SaveUserCredential
// and UpdateUserCredential so tests can assert that the authenticated
// user's ID replaces any client-supplied `userId` in the request body.
type credentialSpyProvider struct {
	*models.DefaultLocalProvider
	observedSave       atomic.Pointer[models.Credential]
	observedUpdate     atomic.Pointer[models.Credential]
	observedDelete     atomic.Pointer[uuid.UUID]
	observedDeleteUser atomic.Pointer[models.User]
	observedGet        atomic.Pointer[uuid.UUID]
	observedList       atomic.Pointer[[2]int] // page, pageSize
	stubbedCredential  *models.Credential
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

func (m *credentialSpyProvider) DeleteUserCredential(req *http.Request, id uuid.UUID) (*models.Credential, error) {
	m.observedDelete.Store(&id)
	if req != nil {
		if user, ok := req.Context().Value(models.UserCtxKey).(*models.User); ok && user != nil {
			m.observedDeleteUser.Store(user)
		}
	}
	return nil, nil
}

func (m *credentialSpyProvider) GetCredentialByID(_ string, id uuid.UUID) (*models.Credential, int, error) {
	m.observedGet.Store(&id)
	if m.stubbedCredential != nil {
		return m.stubbedCredential, http.StatusOK, nil
	}
	return &models.Credential{ID: id}, http.StatusOK, nil
}

func (m *credentialSpyProvider) GetUserCredentials(_ *http.Request, _ string, page, pageSize int, _, _ string) (*models.CredentialsPage, error) {
	m.observedList.Store(&[2]int{page, pageSize})
	return &models.CredentialsPage{Page: page, PageSize: pageSize}, nil
}

// TestGetUserCredentials_PageSizeQueryParamContract pins the list page size to
// the canonical camelCase `pageSize` that schemas' generated getUserCredentials
// sends, while keeping the legacy `page_size` spelling working. The handler
// read only `page_size`, so the Credentials table's page sizes 20 and 25 were
// ignored and every page held the default 10.
func TestGetUserCredentials_PageSizeQueryParamContract(t *testing.T) {
	for _, tc := range []struct {
		name         string
		query        string
		wantPageSize int
	}{
		{"canonical camelCase", "page=1&pageSize=20", 20},
		{"legacy snake_case", "page=1&page_size=20", 20},
		{"camelCase wins when both are sent", "page=1&pageSize=20&page_size=25", 20},
		{"default when neither is sent", "page=1", 10},
		{"capped at 25", "page=1&pageSize=50", 25},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newCredentialSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials?"+tc.query, nil)
			rec := httptest.NewRecorder()

			h.GetUserCredentials(rec, req, nil, &models.User{}, p)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			got := p.observedList.Load()
			if got == nil {
				t.Fatal("provider GetUserCredentials was not invoked")
			}
			if got[0] != 1 {
				t.Errorf("page = %d, want 1", got[0])
			}
			if got[1] != tc.wantPageSize {
				t.Errorf("pageSize = %d, want %d", got[1], tc.wantPageSize)
			}
		})
	}
}

// TestDeleteUserCredential_QueryParamContract pins the delete param to the
// canonical camelCase `credentialId` that schemas' generated
// deleteUserCredential sends, while keeping the legacy `credential_id`
// spelling working. The handler read only `credential_id`, so the UI had to
// hand-roll a local RTK endpoint to talk to it — the generated one resolved a
// nil UUID and deleted nothing.
func TestDeleteUserCredential_QueryParamContract(t *testing.T) {
	id := uuid.Must(uuid.FromString("55555555-5555-5555-5555-555555555555"))

	for _, tc := range []struct {
		name  string
		query string
	}{
		{"canonical camelCase", "credentialId=" + id.String()},
		{"legacy snake_case", "credential_id=" + id.String()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newCredentialSpyProvider()

			req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials?"+tc.query, nil)
			rec := httptest.NewRecorder()

			h.DeleteUserCredential(rec, req, nil, &models.User{}, p)

			if rec.Code != http.StatusOK {
				t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			got := p.observedDelete.Load()
			if got == nil {
				t.Fatal("provider DeleteUserCredential was not invoked")
			}
			if *got != id {
				t.Errorf("deleted id = %v, want %v", *got, id)
			}
		})
	}
}

// TestDeleteUserCredential_RejectsMissingOrInvalidID pins that a missing or
// malformed id is a 400 and never reaches the provider. uuid.FromStringOrNil
// used to coerce both to the zero UUID, so the handler reported a successful
// delete having removed nothing.
func TestDeleteUserCredential_RejectsMissingOrInvalidID(t *testing.T) {
	for _, tc := range []struct {
		name  string
		query string
	}{
		{"absent", ""},
		{"empty", "credentialId="},
		{"malformed", "credentialId=not-a-uuid"},
		{"explicit nil uuid", "credentialId=00000000-0000-0000-0000-000000000000"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newCredentialSpyProvider()

			req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials?"+tc.query, nil)
			rec := httptest.NewRecorder()

			h.DeleteUserCredential(rec, req, nil, &models.User{}, p)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			if p.observedDelete.Load() != nil {
				t.Error("provider DeleteUserCredential was invoked for an unusable id")
			}
			if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
		})
	}
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

func TestSaveUserCredentialAcceptsModelProviderSecretShapes(t *testing.T) {
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}
	tests := []struct {
		name       string
		body       string
		wantType   string
		wantSecret map[string]interface{}
	}{
		{
			name:     "api key provider",
			body:     `{"name":"OpenAI API Key","type":"openai","secret":{"apiKey":"test-key"}}`,
			wantType: "openai",
			wantSecret: map[string]interface{}{
				"apiKey": "test-key",
			},
		},
		{
			name:     "aws bedrock multi field credentials",
			body:     `{"name":"AWS Bedrock","type":"aws-bedrock","secret":{"accessKeyId":"AKIA_TEST","secretAccessKey":"secret","sessionToken":"optional"}}`,
			wantType: "aws-bedrock",
			wantSecret: map[string]interface{}{
				"accessKeyId":     "AKIA_TEST",
				"secretAccessKey": "secret",
				"sessionToken":    "optional",
			},
		},
		{
			name:     "vertex service account json",
			body:     `{"name":"Vertex AI","type":"vertex-ai","secret":{"serviceAccountJson":{"project_id":"meshery-test","client_email":"svc@example.com"}}}`,
			wantType: "vertex-ai",
			wantSecret: map[string]interface{}{
				"serviceAccountJson": map[string]interface{}{
					"project_id":   "meshery-test",
					"client_email": "svc@example.com",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newCredentialSpyProvider()

			req := httptest.NewRequest(http.MethodPost, "/api/integrations/credentials", bytes.NewBufferString(tc.body))
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.SaveUserCredential(rec, req, nil, authUser, p)

			if rec.Code != http.StatusCreated {
				t.Fatalf("expected 201, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			saved := p.observedSave.Load()
			if saved == nil {
				t.Fatal("provider SaveUserCredential was not invoked")
			}
			if saved.Type != tc.wantType {
				t.Fatalf("got type %q, want %q", saved.Type, tc.wantType)
			}
			if !reflect.DeepEqual(map[string]interface{}(saved.Secret), tc.wantSecret) {
				t.Fatalf("got secret %#v, want %#v", saved.Secret, tc.wantSecret)
			}
		})
	}
}

func TestUpdateUserCredentialAcceptsModelProviderSecretShape(t *testing.T) {
	credentialID := uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))
	authUser := &models.User{ID: uuid.Must(uuid.FromString("11111111-1111-1111-1111-111111111111"))}
	body := `{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","name":"AWS Bedrock","type":"aws-bedrock","secret":{"accessKeyId":"AKIA_TEST","secretAccessKey":"secret"}}`

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	h.UpdateUserCredential(rec, req, nil, authUser, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	updated := p.observedUpdate.Load()
	if updated == nil {
		t.Fatal("provider UpdateUserCredential was not invoked")
	}
	if updated.ID != credentialID || updated.Type != "aws-bedrock" {
		t.Fatalf("unexpected updated credential: id=%s type=%q", updated.ID.String(), updated.Type)
	}
	wantSecret := map[string]interface{}{
		"accessKeyId":     "AKIA_TEST",
		"secretAccessKey": "secret",
	}
	if !reflect.DeepEqual(map[string]interface{}(updated.Secret), wantSecret) {
		t.Fatalf("got secret %#v, want %#v", updated.Secret, wantSecret)
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

// TestGetUserCredentialByID_RejectsMissingOrInvalidID pins that a missing,
// malformed, or nil credential ID path parameter returns 400 Bad Request
// with ErrInvalidUUIDValueCode and never reaches the provider.
func TestGetUserCredentialByID_RejectsMissingOrInvalidID(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   string
	}{
		{"absent", ""},
		{"malformed", "not-a-uuid"},
		{"explicit nil uuid", "00000000-0000-0000-0000-000000000000"},
		{"too short", "1234"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			p := newCredentialSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+tc.id, nil)
			if tc.id != "" {
				req = mux.SetURLVars(req, map[string]string{"credentialID": tc.id})
			}
			rec := httptest.NewRecorder()

			h.GetUserCredentialByID(rec, req, nil, &models.User{}, p)

			if rec.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d (body=%q)", rec.Code, rec.Body.String())
			}
			if p.observedGet.Load() != nil {
				t.Error("provider GetCredentialByID was invoked for an unusable id")
			}
			if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != models.ErrInvalidUUIDValueCode {
				t.Errorf("expected code %q, got %q", models.ErrInvalidUUIDValueCode, decoded.Code)
			}
		})
	}
}

// TestGetUserCredentialByID_ValidID verifies that a valid UUID successfully
// queries the provider and returns 200 OK for the owning user.
func TestGetUserCredentialByID_ValidID(t *testing.T) {
	validID := uuid.Must(uuid.FromString("11111111-2222-3333-4444-555555555555"))
	authUser := &models.User{ID: uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))}
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()
	p.stubbedCredential = &models.Credential{
		ID:     validID,
		UserId: authUser.ID,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+validID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": validID.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, authUser, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	observed := p.observedGet.Load()
	if observed == nil {
		t.Fatalf("provider GetCredentialByID was not invoked")
	}
	if *observed != validID {
		t.Fatalf("provider queried with %v, want %v", *observed, validID)
	}
}

// TestGetUserCredentialByID_UnownedCredentialReturns404 verifies that requesting
// a legacy unowned credential (UserId == uuid.Nil) returns 404 Not Found (fail-closed).
func TestGetUserCredentialByID_UnownedCredentialReturns404(t *testing.T) {
	credID := uuid.Must(uuid.FromString("11111111-2222-3333-4444-555555555555"))
	authUser := &models.User{ID: uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))}

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()
	p.stubbedCredential = &models.Credential{
		ID:     credID,
		UserId: uuid.Nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+credID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": credID.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, authUser, p)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body=%q)", rec.Code, rec.Body.String())
	}
}

// TestGetUserCredentialByID_NilUserReturns404 verifies that calling GetUserCredentialByID
// when the caller is unauthenticated/unresolved (user is nil) returns 404 Not Found (fail-closed).
func TestGetUserCredentialByID_NilUserReturns404(t *testing.T) {
	credID := uuid.Must(uuid.FromString("11111111-2222-3333-4444-555555555555"))

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()
	p.stubbedCredential = &models.Credential{
		ID:     credID,
		UserId: uuid.Must(uuid.NewV4()),
	}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+credID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": credID.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, nil, p)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body=%q)", rec.Code, rec.Body.String())
	}
}

// TestGetUserCredentialByID_ForeignUserReturns404 verifies that requesting
// a credential belonging to another user returns 404 Not Found.
func TestGetUserCredentialByID_ForeignUserReturns404(t *testing.T) {
	credID := uuid.Must(uuid.FromString("11111111-2222-3333-4444-555555555555"))
	authUser := &models.User{ID: uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))}
	foreignUser := uuid.Must(uuid.FromString("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"))

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()
	p.stubbedCredential = &models.Credential{
		ID:     credID,
		UserId: foreignUser,
	}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+credID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": credID.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, authUser, p)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d (body=%q)", rec.Code, rec.Body.String())
	}
}

// TestGetUserCredentialByID_OwnerReturns200 verifies that the owner of a
// credential receives 200 OK and the credential body.
func TestGetUserCredentialByID_OwnerReturns200(t *testing.T) {
	credID := uuid.Must(uuid.FromString("11111111-2222-3333-4444-555555555555"))
	owner := &models.User{ID: uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))}

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()
	p.stubbedCredential = &models.Credential{
		ID:     credID,
		UserId: owner.ID,
		Name:   "my-cred",
	}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+credID.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": credID.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, owner, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
}

// TestDeleteUserCredential_PropagatesAuthenticatedUser verifies that
// DeleteUserCredential passes the authenticated user to the provider via request context.
func TestDeleteUserCredential_PropagatesAuthenticatedUser(t *testing.T) {
	id := uuid.Must(uuid.FromString("55555555-5555-5555-5555-555555555555"))
	authUser := &models.User{ID: uuid.Must(uuid.FromString("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"))}

	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials?credentialId="+id.String(), nil)
	rec := httptest.NewRecorder()

	h.DeleteUserCredential(rec, req, nil, authUser, p)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%q)", rec.Code, rec.Body.String())
	}
	receivedUser := p.observedDeleteUser.Load()
	if receivedUser == nil {
		t.Fatal("provider DeleteUserCredential was not passed authenticated user in context")
	}
	if receivedUser.ID != authUser.ID {
		t.Fatalf("received user ID %v, want %v", receivedUser.ID, authUser.ID)
	}
}

func TestSaveUserCredential_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/credentials", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	h.SaveUserCredential(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestSaveUserCredential_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/credentials", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	h.SaveUserCredential(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetUserCredentials_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials", nil)
	rec := httptest.NewRecorder()

	h.GetUserCredentials(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestGetUserCredentials_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials", nil)
	rec := httptest.NewRecorder()

	h.GetUserCredentials(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestUpdateUserCredential_NilUser_ReturnsUnauthorized(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	p := newCredentialSpyProvider()

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	h.UpdateUserCredential(rec, req, nil, nil, p)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 Unauthorized for nil user, got %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshaling error response: %v", err)
	}
	if body["error"] != "user unauthorized" {
		t.Fatalf("expected 'user unauthorized', got %q", body["error"])
	}
}

func TestUpdateUserCredential_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	h.UpdateUserCredential(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestGetUserCredentialByID_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	id := uuid.Must(uuid.NewV4())

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+id.String(), nil)
	req = mux.SetURLVars(req, map[string]string{"credentialID": id.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, nil, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}

func TestDeleteUserCredential_NilProvider_ReturnsInternalServerError(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	id := uuid.Must(uuid.NewV4())
	authUser := &models.User{ID: uuid.Must(uuid.NewV4())}

	req := httptest.NewRequest(http.MethodDelete, "/api/integrations/credentials?credentialId="+id.String(), nil)
	rec := httptest.NewRecorder()

	h.DeleteUserCredential(rec, req, nil, authUser, nil)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 InternalServerError for nil provider, got %d", rec.Code)
	}
	decoded := decodeErrorBody(t, rec.Body.Bytes())
	if decoded.Code != ErrNilProviderCode {
		t.Fatalf("expected %s, got %s", ErrNilProviderCode, decoded.Code)
	}
}
