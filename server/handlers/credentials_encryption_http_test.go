package handlers

// The API-level half of credential encryption at rest.
//
// server/models pins that the persisted row is ciphertext. This file pins what
// a client receives over the wire, which is the claim that lets
// ui/utils/credentialSecret.ts stay a plaintext-only resolver: the provider read
// path decrypts before the handler serializes, so no ciphertext envelope ever
// leaves the server and the UI never has to recognize one.
//
// It drives the real credential handlers over a real DefaultLocalProvider and a
// real datastore - no spy provider - because the encrypt-on-write /
// decrypt-on-read pairing only exists once those three are wired together.

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/schemas/models/core"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	// grafanaAPIKeyFixture is the secret material under test. It must reach the
	// client and must not reach the datastore.
	grafanaAPIKeyFixture = "eyJrIjoiU3VwZXJTZWNyZXRHcmFmYW5hS2V5In0"

	// ciphertextMarker is the reserved property the server writes around a
	// sealed secret. Its appearance in a response body is the failure this file
	// exists to catch.
	ciphertextMarker = "__mesheryEncryptedSecret"
)

// newCredentialAPIHarness returns the real credential handlers wired to a real
// local provider over an in-memory datastore, plus the authenticated user the
// handlers bind credentials to.
func newCredentialAPIHarness(t *testing.T) (*Handler, *models.DefaultLocalProvider, *models.User) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("opening in-memory database: %v", err)
	}
	if err := db.AutoMigrate(&models.Credential{}); err != nil {
		t.Fatalf("migrating credential schema: %v", err)
	}

	provider := &models.DefaultLocalProvider{GenericPersister: &database.Handler{DB: db}}
	user := &models.User{ID: uuid.Must(uuid.NewV4())}

	return newTestHandler(t, map[string]models.Provider{}, ""), provider, user
}

// postCredential saves a credential through the HTTP handler and returns the id
// the row landed under.
func postCredential(t *testing.T, h *Handler, p models.Provider, user *models.User, body string) core.Uuid {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/integrations/credentials", bytes.NewBufferString(body))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "tok"))
	rec := httptest.NewRecorder()

	h.SaveUserCredential(rec, req, nil, user, p)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST credential: status = %d, want 201 (body=%q)", rec.Code, rec.Body.String())
	}

	page := getCredentials(t, h, p, user)
	if len(page.Credentials) == 0 {
		t.Fatal("POST credential: nothing was persisted")
	}
	return page.Credentials[len(page.Credentials)-1].ID
}

// getCredentials lists credentials through the HTTP handler, decoding the same
// body a client receives.
func getCredentials(t *testing.T, h *Handler, p models.Provider, user *models.User) models.CredentialsPage {
	t.Helper()

	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials?order=created_at+asc", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "tok"))
	rec := httptest.NewRecorder()

	h.GetUserCredentials(rec, req, nil, user, p)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET credentials: status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}

	var page models.CredentialsPage
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatalf("GET credentials: decoding response: %v (body=%q)", err, rec.Body.String())
	}
	assertNoCiphertext(t, "GET /api/integrations/credentials", rec.Body.String())
	return page
}

// assertNoCiphertext fails if a response body carries any part of the server's
// at-rest envelope.
func assertNoCiphertext(t *testing.T, what, body string) {
	t.Helper()

	for _, needle := range []string{ciphertextMarker, "meshery.enc.v1:"} {
		if strings.Contains(body, needle) {
			t.Errorf("%s leaked %q to the client; the read path did not decrypt:\n%s", what, needle, body)
		}
	}
}

// TestCredentialAPIReturnsPlaintextWhileTheRowIsCiphertext is the end-to-end
// contract in one test: what the datastore holds and what the client receives
// are different, and only the datastore side is sealed.
func TestCredentialAPIReturnsPlaintextWhileTheRowIsCiphertext(t *testing.T) {
	h, provider, user := newCredentialAPIHarness(t)

	id := postCredential(t, h, provider, user,
		`{"name":"prod-grafana","type":"grafana","secret":{"grafanaURL":"https://grafana.example","grafanaAPIKey":"`+grafanaAPIKeyFixture+`"}}`)

	// The datastore side.
	var storedSecret string
	if err := provider.GetGenericPersister().
		Raw("SELECT secret FROM credentials WHERE id = ?", id).Scan(&storedSecret).Error; err != nil {
		t.Fatalf("reading the stored secret column: %v", err)
	}
	if strings.Contains(storedSecret, grafanaAPIKeyFixture) {
		t.Errorf("the credential's API key was persisted in the clear: %s", storedSecret)
	}
	if !strings.Contains(storedSecret, ciphertextMarker) {
		t.Errorf("the persisted secret is not a Meshery envelope: %s", storedSecret)
	}

	// The client side: the list endpoint the UI's credentials page calls.
	page := getCredentials(t, h, provider, user)
	if len(page.Credentials) != 1 {
		t.Fatalf("listed %d credentials, want 1", len(page.Credentials))
	}
	if got := page.Credentials[0].Secret["grafanaAPIKey"]; got != grafanaAPIKeyFixture {
		t.Errorf("listed grafanaAPIKey = %v, want the plaintext key", got)
	}

	// The client side: the single-credential endpoint a caller about to
	// authenticate uses.
	req := httptest.NewRequest(http.MethodGet, "/api/integrations/credentials/"+id.String(), nil)
	req = mux.SetURLVars(req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "tok")),
		map[string]string{"credentialID": id.String()})
	rec := httptest.NewRecorder()

	h.GetUserCredentialByID(rec, req, nil, user, provider)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET credential by id: status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}
	assertNoCiphertext(t, "GET /api/integrations/credentials/{credentialID}", rec.Body.String())

	var fetched models.Credential
	if err := json.Unmarshal(rec.Body.Bytes(), &fetched); err != nil {
		t.Fatalf("decoding the single-credential response: %v", err)
	}
	if got := fetched.Secret["grafanaAPIKey"]; got != grafanaAPIKeyFixture {
		t.Errorf("fetched grafanaAPIKey = %v, want the plaintext key", got)
	}
}

// TestCredentialAPIServesAPreEncryptionRowUnchanged is the upgrade path seen
// from the client: a row written before encryption shipped is plaintext in the
// datastore and must still serialize into the same response it always did, with
// no migration and no user action.
func TestCredentialAPIServesAPreEncryptionRowUnchanged(t *testing.T) {
	h, provider, user := newCredentialAPIHarness(t)

	legacy := models.Credential{
		ID:     uuid.Must(uuid.NewV4()),
		Name:   "legacy-prometheus",
		Type:   "prometheus",
		UserId: user.ID,
		Secret: core.Map{"prometheusURL": "https://prom.example", "secret": "pre-upgrade-token"},
	}
	if err := provider.GetGenericPersister().Table("credentials").Create(&legacy).Error; err != nil {
		t.Fatalf("seeding a pre-encryption row: %v", err)
	}

	postCredential(t, h, provider, user,
		`{"name":"prod-grafana","type":"grafana","secret":{"grafanaAPIKey":"`+grafanaAPIKeyFixture+`"}}`)

	page := getCredentials(t, h, provider, user)
	if len(page.Credentials) != 2 {
		t.Fatalf("listed %d credentials, want 2", len(page.Credentials))
	}

	secrets := map[string]map[string]interface{}{}
	for _, cred := range page.Credentials {
		secrets[cred.Name] = cred.Secret
	}
	if got := secrets["legacy-prometheus"]["secret"]; got != "pre-upgrade-token" {
		t.Errorf("pre-encryption credential listed with secret = %v, want the untouched plaintext", got)
	}
	if got := secrets["prod-grafana"]["grafanaAPIKey"]; got != grafanaAPIKeyFixture {
		t.Errorf("encrypted credential listed with grafanaAPIKey = %v, want the plaintext key", got)
	}
}

// TestCredentialAPIRenameKeepsTheStoredSecret is the adjacent regression: a
// rename carries no `secret`, the handler unmarshals into an empty non-nil map,
// and GORM's Updates() does not treat that as a zero value - so the rename used
// to wipe the credential's secret. Driven over HTTP because the empty map is
// something the handler creates, not something a caller passes.
func TestCredentialAPIRenameKeepsTheStoredSecret(t *testing.T) {
	h, provider, user := newCredentialAPIHarness(t)

	id := postCredential(t, h, provider, user,
		`{"name":"prod-grafana","type":"grafana","secret":{"grafanaAPIKey":"`+grafanaAPIKeyFixture+`"}}`)

	req := httptest.NewRequest(http.MethodPut, "/api/integrations/credentials",
		bytes.NewBufferString(`{"id":"`+id.String()+`","name":"prod-grafana-renamed","type":"grafana"}`))
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "tok"))
	rec := httptest.NewRecorder()

	h.UpdateUserCredential(rec, req, nil, user, provider)
	if rec.Code != http.StatusOK {
		t.Fatalf("PUT credential: status = %d, want 200 (body=%q)", rec.Code, rec.Body.String())
	}

	page := getCredentials(t, h, provider, user)
	if len(page.Credentials) != 1 {
		t.Fatalf("listed %d credentials, want 1", len(page.Credentials))
	}
	if page.Credentials[0].Name != "prod-grafana-renamed" {
		t.Errorf("credential name = %q, want the renamed one", page.Credentials[0].Name)
	}
	if got := page.Credentials[0].Secret["grafanaAPIKey"]; got != grafanaAPIKeyFixture {
		t.Errorf("renaming the credential left grafanaAPIKey = %v, want the stored plaintext key", got)
	}
}
