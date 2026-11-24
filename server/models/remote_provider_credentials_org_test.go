package models

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/httputil"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta2/organization"
)

const (
	testCredsUserID  = "user-credentials-org"
	testCredsToken   = "test-token"
	testSoleOrgID    = "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	testOtherOrgID   = "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	testSelectedOrg  = "cccccccc-cccc-cccc-cccc-cccccccccccc"
	credsOrgEndpoint = "/api/identity/orgs"
	credsAPIEndpoint = "/api/integrations/credentials"
)

func credentialsOrgCapabilities() Capabilities {
	return Capabilities{
		{Feature: PersistCredentials, Endpoint: credsAPIEndpoint},
		{Feature: PersistOrganizations, Endpoint: credsOrgEndpoint},
		{Feature: SyncPrefs, Endpoint: "/user/preferences"},
	}
}

func orgsPageJSON(t *testing.T, ids ...string) []byte {
	t.Helper()
	page := OrganizationsPage{
		TotalCount: len(ids),
		Page:       0,
		PageSize:   uint64(len(ids)),
	}
	for _, id := range ids {
		orgID := uuid.Must(uuid.FromString(id))
		page.Organizations = append(page.Organizations, &organization.Organization{
			ID:          orgID,
			Name:        "org-" + id[:8],
			Description: "desc-" + id[:8],
		})
	}
	body, err := json.Marshal(page)
	if err != nil {
		t.Fatalf("marshal organizations page: %v", err)
	}
	return body
}

func credentialsPageJSON(t *testing.T) []byte {
	t.Helper()
	body, err := json.Marshal(&CredentialsPage{Page: 0, PageSize: 10, TotalCount: 0, Credentials: nil})
	if err != nil {
		t.Fatalf("marshal credentials page: %v", err)
	}
	return body
}

func newCredentialsOrgRequest(t *testing.T, pref *Preference) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, credsAPIEndpoint, nil)
	req.AddCookie(&http.Cookie{Name: TokenCookieName, Value: testCredsToken})
	if pref != nil {
		req = req.WithContext(context.WithValue(req.Context(), PerfObjCtxKey, pref))
	}
	return req
}

func newCredentialsOrgProvider(t *testing.T, serverURL string) *RemoteProvider {
	t.Helper()
	provider := newTestRemoteProvider(t, serverURL)
	provider.Capabilities = credentialsOrgCapabilities()
	provider.SessionPreferencePersister = newTestSessionPreferencePersister(t)
	provider.syncChan = make(chan *userSession, 1)
	return provider
}

func TestGetUserCredentialsSendsSelectedOrganizationHeader(t *testing.T) {
	var credsCalls, orgsCalls int
	var gotOrg string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			credsCalls++
			gotOrg = r.Header.Get(CurrentOrgIDHeader)
			if r.Header.Get("Authorization") == "" {
				t.Error("expected Authorization on the credentials request")
			}
			_, _ = w.Write(credentialsPageJSON(t))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			orgsCalls++
			t.Error("organizations endpoint should not be called when a selection is already saved")
			http.Error(w, "unexpected", http.StatusInternalServerError)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}
	page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
	if err != nil {
		t.Fatalf("GetUserCredentials: %v", err)
	}
	if page == nil {
		t.Fatal("expected a credentials page, got nil")
	}
	if credsCalls != 1 {
		t.Fatalf("credentials calls = %d, want 1", credsCalls)
	}
	if orgsCalls != 0 {
		t.Fatalf("organizations calls = %d, want 0", orgsCalls)
	}
	if gotOrg != testSelectedOrg {
		t.Errorf("X-Current-Orgid = %q, want %q", gotOrg, testSelectedOrg)
	}
}

func TestGetUserCredentialsAutoSelectsAndPersistsSoleOrganization(t *testing.T) {
	var credsCalls, orgsCalls int
	var gotOrg string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			orgsCalls++
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			credsCalls++
			gotOrg = r.Header.Get(CurrentOrgIDHeader)
			_, _ = w.Write(credentialsPageJSON(t))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := NewDefaultPreference()
	page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
	if err != nil {
		t.Fatalf("GetUserCredentials: %v", err)
	}
	if page == nil {
		t.Fatal("expected a credentials page, got nil")
	}
	if orgsCalls != 1 {
		t.Fatalf("organizations calls = %d, want 1", orgsCalls)
	}
	if credsCalls != 1 {
		t.Fatalf("credentials calls = %d, want 1", credsCalls)
	}
	if gotOrg != testSoleOrgID {
		t.Errorf("X-Current-Orgid = %q, want the auto-selected %q", gotOrg, testSoleOrgID)
	}
	if pref.SelectedOrganizationID != testSoleOrgID {
		t.Errorf("in-memory selection = %q, want %q", pref.SelectedOrganizationID, testSoleOrgID)
	}
	stored, err := provider.ReadFromPersister(testCredsUserID)
	if err != nil {
		t.Fatalf("ReadFromPersister: %v", err)
	}
	if stored.SelectedOrganizationID != testSoleOrgID {
		t.Errorf("persisted selection = %q, want %q", stored.SelectedOrganizationID, testSoleOrgID)
	}
	select {
	case sess := <-provider.syncChan:
		if sess == nil || sess.session == nil || sess.session.SelectedOrganizationID != testSoleOrgID {
			t.Errorf("preference sync did not carry the auto-selected organization")
		}
		if sess.token != testCredsToken {
			t.Errorf("preference sync token = %q, want %q", sess.token, testCredsToken)
		}
	default:
		t.Error("expected a preference sync to be enqueued for the auto-selected organization")
	}
}

func TestGetUserCredentialsRequiresSelectionWhenSeveralOrganizations(t *testing.T) {
	var credsCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID, testOtherOrgID))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			credsCalls++
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(credentialsPageJSON(t))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, NewDefaultPreference()), testCredsUserID, 0, 10, "", "")
	if page != nil {
		t.Fatal("expected no credentials page when selection is required")
	}
	if credsCalls != 0 {
		t.Fatalf("credentials endpoint was called %d times; want 0 when no organization is selected", credsCalls)
	}
	sel, ok := AsOrganizationSelectionRequired(err)
	if !ok {
		t.Fatalf("error = %v, want organization_selection_required with choices", err)
	}
	if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusConflict {
		t.Errorf("provider status = %d, ok=%v, want 409", status, ok)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "select an organization") {
		t.Errorf("error = %q, want it to tell the UI to select an organization", err)
	}
	body := sel.Body()
	if body.Code != OrganizationSelectionRequiredCode {
		t.Errorf("body code = %q, want %q", body.Code, OrganizationSelectionRequiredCode)
	}
	if len(body.Choices) != 2 {
		t.Fatalf("choices = %#v, want 2 organizations", body.Choices)
	}
	gotIDs := map[string]OrganizationSelectionChoice{}
	for _, c := range body.Choices {
		gotIDs[c.ID.String()] = c
	}
	for _, id := range []string{testSoleOrgID, testOtherOrgID} {
		c, ok := gotIDs[id]
		if !ok {
			t.Fatalf("choices missing %s: %#v", id, body.Choices)
		}
		if c.Name != "org-"+id[:8] || c.Description != "desc-"+id[:8] {
			t.Errorf("choice %s = {%q, %q}, want name and description", id, c.Name, c.Description)
		}
	}
	stored, err := provider.ReadFromPersister(testCredsUserID)
	if err != nil {
		t.Fatalf("ReadFromPersister: %v", err)
	}
	if stored != nil && stored.SelectedOrganizationID != "" {
		t.Errorf("selection was persisted unexpectedly: %q", stored.SelectedOrganizationID)
	}
}

func TestGetUserCredentialsNoQualifyingOrganization(t *testing.T) {
	var credsCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t))
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			credsCalls++
			_, _ = w.Write(credentialsPageJSON(t))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, NewDefaultPreference()), testCredsUserID, 0, 10, "", "")
	if page != nil {
		t.Fatal("expected no credentials page when no organization qualifies")
	}
	if credsCalls != 0 {
		t.Fatalf("credentials endpoint was called %d times; want 0", credsCalls)
	}
	if errors.GetCode(err) != ErrNoQualifyingOrganizationCode {
		t.Fatalf("error code = %q, want %q (err=%v)", errors.GetCode(err), ErrNoQualifyingOrganizationCode, err)
	}
	if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusForbidden {
		t.Errorf("provider status = %d, ok=%v, want 403", status, ok)
	}
	if !strings.Contains(strings.ToLower(err.Error()), "no qualifying organization") {
		t.Errorf("error = %q, want a no qualifying organization message", err)
	}
}

func TestGetUserCredentialsDefensiveCloudForbiddenAndNotFound(t *testing.T) {
	cases := []struct {
		name       string
		remoteCode int
		wantCode   string
		wantStatus int
		wantMsg    string
	}{
		{
			name:       "cloud 403",
			remoteCode: http.StatusForbidden,
			wantCode:   ErrNoQualifyingOrganizationCode,
			wantStatus: http.StatusForbidden,
			wantMsg:    "no qualifying organization",
		},
		{
			name:       "cloud 404",
			remoteCode: http.StatusNotFound,
			wantCode:   ErrOrganizationMembershipCode,
			wantStatus: http.StatusNotFound,
			wantMsg:    "not a member of selected org",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint) {
					if r.Header.Get(CurrentOrgIDHeader) != testSelectedOrg {
						t.Errorf("X-Current-Orgid = %q, want %q", r.Header.Get(CurrentOrgIDHeader), testSelectedOrg)
					}
					w.WriteHeader(tc.remoteCode)
					_, _ = io.WriteString(w, `{"error":"scope refused"}`)
					return
				}
				http.NotFound(w, r)
			}))
			defer server.Close()

			provider := newCredentialsOrgProvider(t, server.URL)
			pref := &Preference{SelectedOrganizationID: testSelectedOrg}
			page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
			if page != nil {
				t.Fatal("expected no credentials page on a defensive cloud refusal")
			}
			if errors.GetCode(err) != tc.wantCode {
				t.Fatalf("error code = %q, want %q (err=%v)", errors.GetCode(err), tc.wantCode, err)
			}
			if status, ok := httputil.ProviderStatusCode(err); !ok || status != tc.wantStatus {
				t.Errorf("provider status = %d, ok=%v, want %d", status, ok, tc.wantStatus)
			}
			if !strings.Contains(strings.ToLower(err.Error()), tc.wantMsg) {
				t.Errorf("error = %q, want it to contain %q", err, tc.wantMsg)
			}
		})
	}
}

func TestGetUserCredentialsSurfacesCloudOrganizationSelectionRequired(t *testing.T) {
	choiceA := uuid.Must(uuid.FromString(testSoleOrgID))
	choiceB := uuid.Must(uuid.FromString(testOtherOrgID))
	cloudBody, err := json.Marshal(OrganizationSelectionRequiredError{
		Code:    OrganizationSelectionRequiredCode,
		Message: "select an organization",
		Choices: []OrganizationSelectionChoice{
			{ID: choiceA, Name: "Org A", Description: "first"},
			{ID: choiceB, Name: "Org B", Description: "second"},
		},
	})
	if err != nil {
		t.Fatalf("marshal 409 body: %v", err)
	}

	// Simulate a Cloud that still answers 403 (v1.0.384) as well as 409,
	// so the body is honored in both cases when it carries choices.
	for _, remoteCode := range []int{http.StatusConflict, http.StatusForbidden} {
		t.Run(http.StatusText(remoteCode), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint) {
					if r.Header.Get(CurrentOrgIDHeader) != testSelectedOrg {
						t.Errorf("X-Current-Orgid = %q, want %q", r.Header.Get(CurrentOrgIDHeader), testSelectedOrg)
					}
					w.WriteHeader(remoteCode)
					_, _ = w.Write(cloudBody)
					return
				}
				http.NotFound(w, r)
			}))
			defer server.Close()

			provider := newCredentialsOrgProvider(t, server.URL)
			pref := &Preference{SelectedOrganizationID: testSelectedOrg}
			page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
			if page != nil {
				t.Fatal("expected no credentials page on organization_selection_required")
			}
			sel, ok := AsOrganizationSelectionRequired(err)
			if !ok {
				t.Fatalf("error = %v, want organization_selection_required", err)
			}
			if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusConflict {
				t.Errorf("provider status = %d, ok=%v, want 409", status, ok)
			}
			body := sel.Body()
			if body.Code != OrganizationSelectionRequiredCode {
				t.Errorf("body code = %q, want %q", body.Code, OrganizationSelectionRequiredCode)
			}
			if len(body.Choices) != 2 || body.Choices[0].Name != "Org A" || body.Choices[1].Description != "second" {
				t.Errorf("choices = %#v, want Cloud's org choices", body.Choices)
			}
		})
	}
}

func TestGetUserCredentialsReadsPersistedSelectionWhenPreferenceIsAbsentFromContext(t *testing.T) {
	var gotOrg string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint) {
			gotOrg = r.Header.Get(CurrentOrgIDHeader)
			_, _ = w.Write(credentialsPageJSON(t))
			return
		}
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint) {
			t.Error("organizations endpoint should not be called when a selection is already persisted")
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	if err := provider.WriteToPersister(testCredsUserID, &Preference{SelectedOrganizationID: testSelectedOrg}); err != nil {
		t.Fatalf("WriteToPersister: %v", err)
	}
	if _, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, nil), testCredsUserID, 0, 10, "", ""); err != nil {
		t.Fatalf("GetUserCredentials: %v", err)
	}
	if gotOrg != testSelectedOrg {
		t.Errorf("X-Current-Orgid = %q, want persisted %q", gotOrg, testSelectedOrg)
	}
}

func TestGetUserCredentialsDoesNotPromoteForbiddenWithoutContractCode(t *testing.T) {
	// A 403 body that happens to include a "choices" array but is not the
	// organization_selection_required contract must stay a 403.
	body := `{"choices":[{"id":"` + testSoleOrgID + `","name":"A","description":""},{"id":"` + testOtherOrgID + `","name":"B","description":""}]}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint) {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, body)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}
	_, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
	if _, ok := AsOrganizationSelectionRequired(err); ok {
		t.Fatalf("403 without organization_selection_required was promoted to 409: %v", err)
	}
	if errors.GetCode(err) != ErrNoQualifyingOrganizationCode {
		t.Fatalf("error code = %q, want %q (err=%v)", errors.GetCode(err), ErrNoQualifyingOrganizationCode, err)
	}
	if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusForbidden {
		t.Errorf("provider status = %d, ok=%v, want 403", status, ok)
	}
}

func TestGetUserCredentialsFillsChoicesWhenCloud409BodyIsUnusable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"error":"not the contract body"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID, testOtherOrgID))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}
	page, err := provider.GetUserCredentials(newCredentialsOrgRequest(t, pref), testCredsUserID, 0, 10, "", "")
	if page != nil {
		t.Fatal("expected no credentials page on a 409")
	}
	sel, ok := AsOrganizationSelectionRequired(err)
	if !ok {
		t.Fatalf("error = %v, want organization_selection_required", err)
	}
	if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusConflict {
		t.Errorf("provider status = %d, ok=%v, want 409", status, ok)
	}
	if len(sel.Choices) != 2 {
		t.Fatalf("choices = %#v, want the organizations list as a fallback", sel.Choices)
	}
}

func testCredential(t *testing.T) *Credential {
	t.Helper()
	userID := uuid.Must(uuid.FromString("dddddddd-dddd-dddd-dddd-dddddddddddd"))
	credID := uuid.Must(uuid.FromString("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"))
	return &Credential{ID: credID, UserId: userID, Name: "cred"}
}

func TestCredentialMutationsSendSelectedOrganizationHeader(t *testing.T) {
	cred := testCredential(t)
	userID := cred.UserId.String()
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}

	type verb struct {
		name   string
		method string
		call   func(t *testing.T, provider *RemoteProvider, req *http.Request) error
	}
	verbs := []verb{
		{
			name:   "save",
			method: http.MethodPost,
			call: func(t *testing.T, provider *RemoteProvider, _ *http.Request) error {
				t.Helper()
				_, err := provider.SaveUserCredential(testCredsToken, cred)
				return err
			},
		},
		{
			name:   "update",
			method: http.MethodPut,
			call: func(t *testing.T, provider *RemoteProvider, req *http.Request) error {
				t.Helper()
				_, err := provider.UpdateUserCredential(req, cred)
				return err
			},
		},
		{
			name:   "delete",
			method: http.MethodDelete,
			call: func(t *testing.T, provider *RemoteProvider, req *http.Request) error {
				t.Helper()
				_, err := provider.DeleteUserCredential(req, cred.ID)
				return err
			},
		},
		{
			name:   "getByID",
			method: http.MethodGet,
			call: func(t *testing.T, provider *RemoteProvider, _ *http.Request) error {
				t.Helper()
				_, _, err := provider.GetCredentialByID(testCredsToken, userID, cred.ID)
				return err
			},
		},
	}

	for _, tc := range verbs {
		t.Run(tc.name, func(t *testing.T) {
			var credsCalls, orgsCalls int
			var gotOrg string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
					if r.Method != tc.method {
						t.Errorf("method = %s, want %s", r.Method, tc.method)
					}
					credsCalls++
					gotOrg = r.Header.Get(CurrentOrgIDHeader)
					if r.Method == http.MethodPost {
						w.WriteHeader(http.StatusCreated)
					} else {
						w.WriteHeader(http.StatusOK)
					}
					_, _ = io.WriteString(w, `{"id":"`+cred.ID.String()+`","name":"cred"}`)
				case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
					orgsCalls++
					http.Error(w, "organizations should not be listed when a selection exists", http.StatusInternalServerError)
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			provider := newCredentialsOrgProvider(t, server.URL)
			if err := provider.WriteToPersister(userID, pref); err != nil {
				t.Fatalf("seed persister: %v", err)
			}
			if err := tc.call(t, provider, newCredentialsOrgRequest(t, pref)); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			if credsCalls != 1 {
				t.Errorf("credentials calls = %d, want 1", credsCalls)
			}
			if orgsCalls != 0 {
				t.Errorf("organizations calls = %d, want 0", orgsCalls)
			}
			if gotOrg != testSelectedOrg {
				t.Errorf("X-Current-Orgid = %q, want %q", gotOrg, testSelectedOrg)
			}
		})
	}
}

func TestCredentialMutationsRequireSelectionWhenSeveralOrganizations(t *testing.T) {
	cred := testCredential(t)
	userID := cred.UserId.String()
	verbs := []struct {
		name string
		call func(t *testing.T, provider *RemoteProvider, req *http.Request) error
	}{
		{"save", func(t *testing.T, provider *RemoteProvider, _ *http.Request) error {
			t.Helper()
			_, err := provider.SaveUserCredential(testCredsToken, cred)
			return err
		}},
		{"update", func(t *testing.T, provider *RemoteProvider, req *http.Request) error {
			t.Helper()
			_, err := provider.UpdateUserCredential(req, cred)
			return err
		}},
		{"delete", func(t *testing.T, provider *RemoteProvider, req *http.Request) error {
			t.Helper()
			_, err := provider.DeleteUserCredential(req, cred.ID)
			return err
		}},
		{"getByID", func(t *testing.T, provider *RemoteProvider, _ *http.Request) error {
			t.Helper()
			_, _, err := provider.GetCredentialByID(testCredsToken, userID, cred.ID)
			return err
		}},
	}

	for _, tc := range verbs {
		t.Run(tc.name, func(t *testing.T) {
			var credsCalls int
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch {
				case strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
					credsCalls++
					http.Error(w, "credentials must not be called", http.StatusInternalServerError)
				case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
					_, _ = w.Write(orgsPageJSON(t, testSoleOrgID, testOtherOrgID))
				default:
					http.NotFound(w, r)
				}
			}))
			defer server.Close()

			provider := newCredentialsOrgProvider(t, server.URL)
			if err := provider.WriteToPersister(userID, NewDefaultPreference()); err != nil {
				t.Fatalf("seed persister: %v", err)
			}
			req := newCredentialsOrgRequest(t, NewDefaultPreference())
			req = req.WithContext(context.WithValue(req.Context(), UserCtxKey, &User{ID: cred.UserId}))
			err := tc.call(t, provider, req)
			if credsCalls != 0 {
				t.Errorf("credentials calls = %d, want 0", credsCalls)
			}
			sel, ok := AsOrganizationSelectionRequired(err)
			if !ok {
				t.Fatalf("error = %v, want organization_selection_required", err)
			}
			if len(sel.Choices) != 2 {
				t.Errorf("choices = %#v, want two organizations", sel.Choices)
			}
			if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusConflict {
				t.Errorf("provider status = %d, ok=%v, want 409", status, ok)
			}
		})
	}
}

func TestCredentialMutationsAutoSelectSoleOrganization(t *testing.T) {
	cred := testCredential(t)
	userID := cred.UserId.String()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, credsAPIEndpoint):
			if r.Header.Get(CurrentOrgIDHeader) != testSoleOrgID {
				t.Errorf("X-Current-Orgid = %q, want auto-selected %q", r.Header.Get(CurrentOrgIDHeader), testSoleOrgID)
			}
			w.WriteHeader(http.StatusCreated)
			_, _ = io.WriteString(w, `{"id":"`+cred.ID.String()+`","name":"cred"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	if _, err := provider.SaveUserCredential(testCredsToken, cred); err != nil {
		t.Fatalf("SaveUserCredential: %v", err)
	}
	stored, err := provider.ReadFromPersister(userID)
	if err != nil {
		t.Fatalf("ReadFromPersister: %v", err)
	}
	if stored == nil || stored.SelectedOrganizationID != testSoleOrgID {
		t.Fatalf("persisted selection = %#v, want %q", stored, testSoleOrgID)
	}
	select {
	case session := <-provider.syncChan:
		if session == nil || session.token != testCredsToken {
			t.Errorf("sync token = %#v, want %q", session, testCredsToken)
		}
		if session.session == nil || session.session.SelectedOrganizationID != testSoleOrgID {
			t.Errorf("sync selection = %#v, want %q", session.session, testSoleOrgID)
		}
	default:
		t.Fatal("expected a preference sync after auto-select")
	}
}

func TestGetCredentialByIDDefensiveCloudNotFoundStaysResourceNotFoundWhenStillAMember(t *testing.T) {
	cred := testCredential(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint+"/"):
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"error":"credential not found"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t, testSelectedOrg))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}
	if err := provider.WriteToPersister(cred.UserId.String(), pref); err != nil {
		t.Fatalf("seed persister: %v", err)
	}
	_, status, err := provider.GetCredentialByID(testCredsToken, cred.UserId.String(), cred.ID)
	if err == nil {
		t.Fatal("expected a not-found error")
	}
	if _, ok := AsOrganizationSelectionRequired(err); ok {
		t.Fatalf("resource 404 was promoted to 409: %v", err)
	}
	if errors.GetCode(err) == ErrOrganizationMembershipCode {
		t.Fatalf("resource 404 was mapped to membership: %v", err)
	}
	if status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", status)
	}
}

func TestGetCredentialByIDMapsCloudNotFoundToMembershipWhenOrgIsNotTheCallers(t *testing.T) {
	cred := testCredential(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint+"/"):
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(w, `{"error":"not a member"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID, testOtherOrgID))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	pref := &Preference{SelectedOrganizationID: testSelectedOrg}
	if err := provider.WriteToPersister(cred.UserId.String(), pref); err != nil {
		t.Fatalf("seed persister: %v", err)
	}
	_, status, err := provider.GetCredentialByID(testCredsToken, cred.UserId.String(), cred.ID)
	if errors.GetCode(err) != ErrOrganizationMembershipCode {
		t.Fatalf("error code = %q, want %q (err=%v)", errors.GetCode(err), ErrOrganizationMembershipCode, err)
	}
	if status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", status)
	}
}

func TestSaveUserCredentialDoesNotTreatResourceConflictAsSelectionRequired(t *testing.T) {
	cred := testCredential(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, credsAPIEndpoint) {
			w.WriteHeader(http.StatusConflict)
			_, _ = io.WriteString(w, `{"error":"credential name already exists"}`)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	if err := provider.WriteToPersister(cred.UserId.String(), &Preference{SelectedOrganizationID: testSelectedOrg}); err != nil {
		t.Fatalf("seed persister: %v", err)
	}
	_, err := provider.SaveUserCredential(testCredsToken, cred)
	if _, ok := AsOrganizationSelectionRequired(err); ok {
		t.Fatalf("duplicate-name 409 was promoted to organization_selection_required: %v", err)
	}
	if err == nil {
		t.Fatal("expected a conflict error")
	}
	if status, ok := httputil.ProviderStatusCode(err); !ok || status != http.StatusConflict {
		t.Errorf("provider status = %d, ok=%v, want 409", status, ok)
	}
}

func TestGetCredentialByIDWithoutUserIDDoesNotRequireSelection(t *testing.T) {
	// Token-only callers such as K8sContextFromConnection cannot consult a
	// saved selection. Listing orgs and 409-ing would break multi-org users.
	cred := testCredential(t)
	var credsCalls, orgsCalls int
	var gotOrg string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint+"/"):
			credsCalls++
			gotOrg = r.Header.Get(CurrentOrgIDHeader)
			_, _ = io.WriteString(w, `{"id":"`+cred.ID.String()+`","name":"cred"}`)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsOrgEndpoint):
			orgsCalls++
			_, _ = w.Write(orgsPageJSON(t, testSoleOrgID, testOtherOrgID))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	got, status, err := provider.GetCredentialByID(testCredsToken, "", cred.ID)
	if err != nil {
		t.Fatalf("GetCredentialByID: %v", err)
	}
	if status != http.StatusOK || got == nil || got.ID != cred.ID {
		t.Fatalf("got status=%d cred=%#v, want 200", status, got)
	}
	if credsCalls != 1 {
		t.Errorf("credentials calls = %d, want 1", credsCalls)
	}
	if orgsCalls != 0 {
		t.Errorf("organizations calls = %d, want 0", orgsCalls)
	}
	if gotOrg != "" {
		t.Errorf("X-Current-Orgid = %q, want empty on the token-only path", gotOrg)
	}
}

func TestGetCredentialByIDPlainForbiddenStaysResourceError(t *testing.T) {
	cred := testCredential(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, credsAPIEndpoint+"/") {
			w.WriteHeader(http.StatusForbidden)
			_, _ = io.WriteString(w, `{"error":"not allowed to read this credential"}`)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	provider := newCredentialsOrgProvider(t, server.URL)
	if err := provider.WriteToPersister(cred.UserId.String(), &Preference{SelectedOrganizationID: testSelectedOrg}); err != nil {
		t.Fatalf("seed persister: %v", err)
	}
	_, status, err := provider.GetCredentialByID(testCredsToken, cred.UserId.String(), cred.ID)
	if err == nil {
		t.Fatal("expected a forbidden error")
	}
	if errors.GetCode(err) == ErrNoQualifyingOrganizationCode {
		t.Fatalf("resource 403 was mapped to no qualifying organization: %v", err)
	}
	if _, ok := AsOrganizationSelectionRequired(err); ok {
		t.Fatalf("resource 403 was promoted to 409: %v", err)
	}
	if status != http.StatusForbidden {
		t.Errorf("status = %d, want 403", status)
	}
}
