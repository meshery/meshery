package handlers

import (
	"encoding/json"
	"testing"
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
