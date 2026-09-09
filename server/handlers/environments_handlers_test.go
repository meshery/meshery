package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/v1beta3/environment"
)

// TestEnvironmentPayloadWire_UnmarshalJSON exercises the dual-accept body
// contract for POST /api/environments. The wrapper must route both
// `organizationId` (canonical) and `organization_id` (legacy) onto the
// underlying schemas-generated OrgID field, with canonical taking
// precedence when both are present. Mirrors
// TestWorkspacePayloadWire_UnmarshalJSON.
func TestEnvironmentPayloadWire_UnmarshalJSON(t *testing.T) {
	const (
		canonicalUUID = "11111111-1111-1111-1111-111111111111"
		legacyUUID    = "22222222-2222-2222-2222-222222222222"
	)

	cases := []struct {
		name    string
		body    string
		wantOrg string
	}{
		{
			name:    "canonical organizationId only",
			body:    `{"name":"env","organizationId":"` + canonicalUUID + `"}`,
			wantOrg: canonicalUUID,
		},
		{
			name:    "legacy organization_id only",
			body:    `{"name":"env","organization_id":"` + legacyUUID + `"}`,
			wantOrg: legacyUUID,
		},
		{
			name:    "both supplied, canonical wins",
			body:    `{"name":"env","organizationId":"` + canonicalUUID + `","organization_id":"` + legacyUUID + `"}`,
			wantOrg: canonicalUUID,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var wire environmentPayloadWire
			if err := json.Unmarshal([]byte(tc.body), &wire); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if got := wire.OrgID.String(); got != tc.wantOrg {
				t.Fatalf("OrgID = %q, want %q", got, tc.wantOrg)
			}
			if wire.Name != "env" {
				t.Fatalf("Name = %q, want %q", wire.Name, "env")
			}
		})
	}
}

// TestEnvironmentPayloadWire_MarshalsCanonicalCamelCase is the regression
// test for the bug where environment creation failed for every Layer5 Cloud
// user: the handler previously wrapped the deprecated v1beta1
// environment.EnvironmentPayload (json tag `organization_id`). Unmarshaling
// a client's camelCase `organizationId` into that wrapper succeeded, but
// remote_provider.go re-marshals the unwrapped EnvironmentPayload to build
// the outbound request to the remote provider - and that second marshal
// used the v1beta1 struct's own snake_case tag, silently downgrading the
// org id back to `organization_id` on the wire to Layer5 Cloud. Layer5
// Cloud's endpoint expects camelCase, so the org id was never populated,
// which violated a NOT NULL/FK constraint on the environments table and
// produced a 500 (surfaced to the UI as a misleading 404 "unable to get
// result", from the handler always mapping remote errors to
// ErrGetResult/http.StatusNotFound).
//
// This asserts the fix: EnvironmentPayload now comes from v1beta3, whose
// native tag is camelCase, so the outbound marshal is correct regardless
// of which spelling the client sent.
func TestEnvironmentPayloadWire_MarshalsCanonicalCamelCase(t *testing.T) {
	const orgUUID = "33333333-3333-3333-3333-333333333333"

	cases := []struct {
		name string
		body string
	}{
		{name: "client sent canonical organizationId", body: `{"name":"env","organizationId":"` + orgUUID + `"}`},
		{name: "client sent legacy organization_id", body: `{"name":"env","organization_id":"` + orgUUID + `"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var wire environmentPayloadWire
			if err := json.Unmarshal([]byte(tc.body), &wire); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}

			// Simulate remote_provider.go's SaveEnvironment/UpdateEnvironment,
			// which re-marshals the unwrapped EnvironmentPayload to build the
			// outbound request to the remote provider.
			out, err := json.Marshal(&wire.EnvironmentPayload)
			if err != nil {
				t.Fatalf("unexpected marshal error: %v", err)
			}

			var onWire map[string]interface{}
			if err := json.Unmarshal(out, &onWire); err != nil {
				t.Fatalf("unexpected re-unmarshal error: %v", err)
			}

			if _, present := onWire["organization_id"]; present {
				t.Fatalf("outbound payload regressed to snake_case organization_id: %s", out)
			}
			if got, _ := onWire["organizationId"].(string); got != orgUUID {
				t.Fatalf("outbound payload organizationId = %q, want %q (full body: %s)", got, orgUUID, out)
			}
		})
	}
}

// environmentFailingProvider embeds DefaultLocalProvider and fails every
// environment call with a caller-supplied error, so the tests below can assert
// what the handler does with a provider failure.
type environmentFailingProvider struct {
	*models.DefaultLocalProvider
	err error
}

func newEnvironmentFailingProvider(err error) *environmentFailingProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &environmentFailingProvider{DefaultLocalProvider: base, err: err}
}

func (m *environmentFailingProvider) SaveEnvironment(_ *http.Request, _ *environment.EnvironmentPayload, _ string, _ bool) ([]byte, error) {
	return nil, m.err
}

func (m *environmentFailingProvider) GetEnvironments(_, _, _, _, _, _, _ string) ([]byte, error) {
	return nil, m.err
}

func (m *environmentFailingProvider) DeleteEnvironment(_ *http.Request, _ string) ([]byte, error) {
	return nil, m.err
}

// decodeErrorBody parses the JSON error envelope every non-2xx response
// carries (see docs/content/en/project/contributing/error-contract.md).
func decodeErrorBody(t *testing.T, body []byte) struct {
	Error                string   `json:"error"`
	Code                 string   `json:"code"`
	ProbableCause        []string `json:"probableCause"`
	SuggestedRemediation []string `json:"suggestedRemediation"`
} {
	t.Helper()
	var decoded struct {
		Error                string   `json:"error"`
		Code                 string   `json:"code"`
		ProbableCause        []string `json:"probableCause"`
		SuggestedRemediation []string `json:"suggestedRemediation"`
	}
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("error body did not parse as JSON: %v (body=%q)", err, body)
	}
	return decoded
}

// TestSaveEnvironmentHandler_PropagatesProviderStatus is the regression test
// for the bug that showed an Org Admin a success message for an environment
// that was never created.
//
// The handler used to answer every provider failure with
// writeMeshkitError(w, ErrGetResult(err), http.StatusNotFound): a hardcoded 404
// carrying meshery-server-1033 ("unable to get result", probable cause "Result
// Identifier provided is not valid"). A remote-provider 403 therefore reached
// the browser as a 404 describing the performance-results subsystem, and the
// UI did not recognise it as a failure at all.
//
// The contract asserted here: the provider's real status reaches the client,
// and the error code names the operation that actually failed.
func TestSaveEnvironmentHandler_PropagatesProviderStatus(t *testing.T) {
	cases := []struct {
		name         string
		providerErr  error
		wantStatus   int
		wantCodeIs   string
		wantCauseHas string
	}{
		{
			name:        "provider 403 surfaces as 403, not 404",
			providerErr: models.ErrPost(errors.New("failed to save the environment"), "Environment", http.StatusForbidden),
			wantStatus:  http.StatusForbidden,
			wantCodeIs:  ErrSaveEnvironmentCode,
		},
		{
			name:        "provider 409 surfaces as 409",
			providerErr: models.ErrPost(errors.New("environment already exists"), "Environment", http.StatusConflict),
			wantStatus:  http.StatusConflict,
			wantCodeIs:  ErrSaveEnvironmentCode,
		},
		{
			name:        "provider 500 surfaces as 500",
			providerErr: models.ErrPost(errors.New("internal error"), "Environment", http.StatusInternalServerError),
			wantStatus:  http.StatusInternalServerError,
			wantCodeIs:  ErrSaveEnvironmentCode,
		},
		{
			name:        "unreachable provider defaults to 502, never 404",
			providerErr: models.ErrUnreachableRemoteProvider(errors.New("dial tcp: connection refused")),
			wantStatus:  http.StatusBadGateway,
			wantCodeIs:  ErrSaveEnvironmentCode,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newEnvironmentFailingProvider(tc.providerErr)

			body := `{"name":"prod","description":"","organizationId":"11111111-1111-1111-1111-111111111111"}`
			req := httptest.NewRequest(http.MethodPost, "/api/environments", strings.NewReader(body))
			rec := httptest.NewRecorder()

			h.SaveEnvironment(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d (body=%q)", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if rec.Code == http.StatusCreated {
				t.Fatal("a failed create must never answer 201 - that is what made the UI show success")
			}

			decoded := decodeErrorBody(t, rec.Body.Bytes())
			if decoded.Code != tc.wantCodeIs {
				t.Errorf("code = %q, want %q", decoded.Code, tc.wantCodeIs)
			}
			if decoded.Code == ErrGetResultCode {
				t.Errorf("code regressed to the performance-results code %s", ErrGetResultCode)
			}
			if strings.Contains(decoded.Error, "unable to get result") {
				t.Errorf("message regressed to the results wording: %q", decoded.Error)
			}
			if len(decoded.SuggestedRemediation) == 0 {
				t.Error("expected environment-specific remediation on the wire")
			}
		})
	}
}

// TestGetEnvironmentsHandler_PropagatesProviderStatus covers the read path with
// the same contract: a provider 403 must not be reported as a 404.
func TestGetEnvironmentsHandler_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newEnvironmentFailingProvider(
		models.ErrFetch(errors.New("forbidden"), "Environments", http.StatusForbidden),
	)

	req := httptest.NewRequest(http.MethodGet, "/api/environments?orgId=11111111-1111-1111-1111-111111111111", nil)
	req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
	rec := httptest.NewRecorder()

	h.GetEnvironments(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrGetEnvironmentsCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrGetEnvironmentsCode)
	}
}

// TestDeleteEnvironmentHandler_PropagatesProviderStatus covers the delete path.
func TestDeleteEnvironmentHandler_PropagatesProviderStatus(t *testing.T) {
	h := newTestHandler(t, map[string]models.Provider{}, "")
	provider := newEnvironmentFailingProvider(
		models.ErrFetch(errors.New("forbidden"), "Environment", http.StatusForbidden),
	)

	req := httptest.NewRequest(http.MethodDelete, "/api/environments/env-1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "env-1"})
	rec := httptest.NewRecorder()

	h.DeleteEnvironmentHandler(rec, req, nil, nil, provider)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d (body=%q)", rec.Code, http.StatusForbidden, rec.Body.String())
	}
	if decoded := decodeErrorBody(t, rec.Body.Bytes()); decoded.Code != ErrDeleteEnvironmentCode {
		t.Errorf("code = %q, want %q", decoded.Code, ErrDeleteEnvironmentCode)
	}
}
