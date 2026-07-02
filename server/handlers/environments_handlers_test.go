package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

// environmentSpyProvider embeds DefaultLocalProvider and records the orgID
// that GetEnvironments / GetEnvironmentByID are invoked with. It allows the
// handler tests below to verify that the handler extracted the query
// parameter correctly.
type environmentSpyProvider struct {
	*models.DefaultLocalProvider
	observedOrgID string
	called        bool
}

func newEnvironmentSpyProvider() *environmentSpyProvider {
	base := &models.DefaultLocalProvider{}
	base.Initialize()
	return &environmentSpyProvider{DefaultLocalProvider: base}
}

func (m *environmentSpyProvider) GetEnvironments(_, _, _, _, _, _, orgID string) ([]byte, error) {
	m.called = true
	m.observedOrgID = orgID
	return []byte(`{"environments":[]}`), nil
}

func (m *environmentSpyProvider) GetEnvironmentByID(_ *http.Request, _, orgID string) ([]byte, error) {
	m.called = true
	m.observedOrgID = orgID
	return []byte(`{}`), nil
}

// TestGetEnvironmentsHandler_RequiresOrgId asserts that the handler returns
// 400 when orgId is missing and 200 when it is provided.
func TestGetEnvironmentsHandler_RequiresOrgId(t *testing.T) {
	cases := []struct {
		name         string
		rawQuery     string
		wantStatus   int
		wantOrgID    string
		wantProvider bool
	}{
		{
			name:         "orgId present returns 200",
			rawQuery:     "orgId=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "missing orgId returns 400",
			rawQuery:     "",
			wantStatus:   http.StatusBadRequest,
			wantProvider: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newEnvironmentSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/environments?"+tc.rawQuery, nil)
			req = req.WithContext(context.WithValue(req.Context(), models.TokenCtxKey, "test-token"))
			rec := httptest.NewRecorder()

			h.GetEnvironments(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d (body=%q)", tc.wantStatus, rec.Code, rec.Body.String())
			}

			if provider.called != tc.wantProvider {
				t.Fatalf("provider called=%v, want %v", provider.called, tc.wantProvider)
			}

			if tc.wantProvider && provider.observedOrgID != tc.wantOrgID {
				t.Fatalf("provider received orgID=%q, want %q", provider.observedOrgID, tc.wantOrgID)
			}

			if tc.wantStatus == http.StatusBadRequest {
				if !strings.Contains(rec.Body.String(), "orgId") {
					t.Errorf("expected 400 body to mention orgId, got %q", rec.Body.String())
				}
			}
		})
	}
}

// TestGetEnvironmentByIDHandler_RequiresOrgId mirrors the coverage above
// for the single-environment endpoint.
func TestGetEnvironmentByIDHandler_RequiresOrgId(t *testing.T) {
	cases := []struct {
		name         string
		rawQuery     string
		wantStatus   int
		wantOrgID    string
		wantProvider bool
	}{
		{
			name:         "orgId present returns 200",
			rawQuery:     "orgId=abc",
			wantStatus:   http.StatusOK,
			wantOrgID:    "abc",
			wantProvider: true,
		},
		{
			name:         "missing orgId returns 400",
			rawQuery:     "",
			wantStatus:   http.StatusBadRequest,
			wantProvider: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			h := newTestHandler(t, map[string]models.Provider{}, "")
			provider := newEnvironmentSpyProvider()

			req := httptest.NewRequest(http.MethodGet, "/api/environments/env-1?"+tc.rawQuery, nil)
			req = mux.SetURLVars(req, map[string]string{"id": "env-1"})
			rec := httptest.NewRecorder()

			h.GetEnvironmentByIDHandler(rec, req, nil, nil, provider)

			if rec.Code != tc.wantStatus {
				t.Fatalf("expected status %d, got %d (body=%q)", tc.wantStatus, rec.Code, rec.Body.String())
			}

			if provider.called != tc.wantProvider {
				t.Fatalf("provider called=%v, want %v", provider.called, tc.wantProvider)
			}

			if tc.wantProvider && provider.observedOrgID != tc.wantOrgID {
				t.Fatalf("provider received orgID=%q, want %q", provider.observedOrgID, tc.wantOrgID)
			}

			if tc.wantStatus == http.StatusBadRequest {
				if !strings.Contains(rec.Body.String(), "orgId") {
					t.Errorf("expected 400 body to mention orgId, got %q", rec.Body.String())
				}
			}
		})
	}
}

// TestEnvironmentPayloadWire_UnmarshalJSON exercises the dual-accept body
// contract for POST /api/environments. The schemas-generated struct tags
// OrgId as json:"organization_id", but canonical in-repo consumers now emit
// `organizationId`. This wrapper must intercept both spellings during the
// Phase 2 deprecation window. Canonical must win when both are supplied.
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
		{
			name:    "neither supplied, OrgId is zero UUID",
			body:    `{"name":"env"}`,
			wantOrg: "00000000-0000-0000-0000-000000000000",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var wire environmentPayloadWire
			if err := json.Unmarshal([]byte(tc.body), &wire); err != nil {
				t.Fatalf("unexpected unmarshal error: %v", err)
			}
			if got := wire.OrgId.String(); got != tc.wantOrg {
				t.Fatalf("OrgId = %q, want %q", got, tc.wantOrg)
			}
			if wire.Name != "env" {
				t.Fatalf("Name = %q, want %q", wire.Name, "env")
			}
		})
	}
}

// TestEnvironmentPayloadWire_UnmarshalJSON_StaleReceiverReset verifies that
// reusing an environmentPayloadWire receiver zeroes OrgId when neither
// spelling is present, preventing stale data leaking between unmarshal calls.
func TestEnvironmentPayloadWire_UnmarshalJSON_StaleReceiverReset(t *testing.T) {
	const seedUUID = "11111111-1111-1111-1111-111111111111"

	var wire environmentPayloadWire
	// First unmarshal seeds OrgId.
	if err := json.Unmarshal([]byte(`{"name":"first","organizationId":"`+seedUUID+`"}`), &wire); err != nil {
		t.Fatalf("first unmarshal error: %v", err)
	}
	if wire.OrgId.String() != seedUUID {
		t.Fatalf("setup failed: OrgId = %q, want %q", wire.OrgId.String(), seedUUID)
	}

	// Second unmarshal omits OrgId — receiver must be zeroed, not stale.
	if err := json.Unmarshal([]byte(`{"name":"second"}`), &wire); err != nil {
		t.Fatalf("second unmarshal error: %v", err)
	}
	const zeroUUID = "00000000-0000-0000-0000-000000000000"
	if wire.OrgId.String() != zeroUUID {
		t.Fatalf("stale OrgId not zeroed: got %q, want %q", wire.OrgId.String(), zeroUUID)
	}
}