package models

import (
	"encoding/json"
	"strings"
	"testing"
)

const testSelectedOrgID = "0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40"

// TestPreferenceMarshalsCanonicalSelectedOrganizationId pins the selected
// organization to the canonical camelCase `selectedOrganizationId` that the
// schemas v1beta1 user.Preference construct declares and meshery-cloud reads.
//
// Meshery Server spelled it `selectedOrganizationID` (all-caps ID, which the
// identifier-naming contract forbids). Every other key on this struct is
// camelCase, so the odd one silently broke the round trip through a remote
// provider: executePrefSync PUT a key meshery-cloud ignores, and the provider's
// reply carried a key this struct ignored. The selected organization therefore
// never persisted to Layer5 Cloud.
func TestPreferenceMarshalsCanonicalSelectedOrganizationId(t *testing.T) {
	b, err := json.Marshal(&Preference{SelectedOrganizationID: testSelectedOrgID})
	if err != nil {
		t.Fatalf("marshal Preference: %v", err)
	}
	got := string(b)

	if !strings.Contains(got, `"selectedOrganizationId"`) {
		t.Errorf(`marshaled preference is missing the canonical "selectedOrganizationId" key: %s`, got)
	}
	if strings.Contains(got, `"selectedOrganizationID"`) {
		t.Errorf(`marshaled preference still emits the legacy all-caps key: %s`, got)
	}
}

// TestPreferenceUnmarshalsCanonicalAndLegacySelectedOrganizationId pins that a
// provider reply on the canonical key is read (the bug), and that preferences
// already persisted under the legacy key stay readable after the rename.
func TestPreferenceUnmarshalsCanonicalAndLegacySelectedOrganizationId(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{"canonical", `{"selectedOrganizationId":"` + testSelectedOrgID + `"}`},
		{"legacy", `{"selectedOrganizationID":"` + testSelectedOrgID + `"}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var p Preference
			if err := json.Unmarshal([]byte(tc.body), &p); err != nil {
				t.Fatalf("unmarshal Preference: %v", err)
			}
			if p.SelectedOrganizationID != testSelectedOrgID {
				t.Errorf("SelectedOrganizationID = %q, want %q", p.SelectedOrganizationID, testSelectedOrgID)
			}
		})
	}
}

// TestPreferenceCanonicalKeyWinsOverLegacy pins precedence when a payload
// carries both: the canonical key is authoritative, so a stale legacy value
// cannot shadow a fresh canonical one.
func TestPreferenceCanonicalKeyWinsOverLegacy(t *testing.T) {
	body := `{"selectedOrganizationId":"` + testSelectedOrgID + `","selectedOrganizationID":"stale"}`

	var p Preference
	if err := json.Unmarshal([]byte(body), &p); err != nil {
		t.Fatalf("unmarshal Preference: %v", err)
	}
	if p.SelectedOrganizationID != testSelectedOrgID {
		t.Errorf("SelectedOrganizationID = %q, want the canonical %q", p.SelectedOrganizationID, testSelectedOrgID)
	}
}

// TestPreferenceLegacyKeyUpdatesAlreadySelectedOrganization pins the decision on
// payload key presence rather than on the destination being empty.
//
// UserPrefsHandler decodes the request body onto the *models.Preference the
// session middleware already populated from the persister, so a user who has an
// organization selected reaches UnmarshalJSON with the field non-empty. Gating
// the legacy fallback on that field being empty silently discarded the update
// for exactly the out-of-repo clients (Meshery UI extensions, Kanvas) the
// fallback exists to serve - the handler still answered 200 and re-persisted the
// old value.
func TestPreferenceLegacyKeyUpdatesAlreadySelectedOrganization(t *testing.T) {
	const previouslySelectedOrgID = "0195b0ab-1f4d-7a3c-9c1e-000000000001"

	p := &Preference{SelectedOrganizationID: previouslySelectedOrgID}
	body := `{"selectedOrganizationID":"` + testSelectedOrgID + `"}`

	if err := json.Unmarshal([]byte(body), p); err != nil {
		t.Fatalf("unmarshal Preference: %v", err)
	}
	if p.SelectedOrganizationID != testSelectedOrgID {
		t.Errorf("SelectedOrganizationID = %q, want the legacy-keyed update %q", p.SelectedOrganizationID, testSelectedOrgID)
	}
}

// TestPreferenceCanonicalKeyWinsOverLegacyOnMerge pins that presence-based
// precedence still lets the canonical key win when both are sent onto an
// already-populated Preference.
func TestPreferenceCanonicalKeyWinsOverLegacyOnMerge(t *testing.T) {
	p := &Preference{SelectedOrganizationID: "0195b0ab-1f4d-7a3c-9c1e-000000000001"}
	body := `{"selectedOrganizationId":"` + testSelectedOrgID + `","selectedOrganizationID":"stale"}`

	if err := json.Unmarshal([]byte(body), p); err != nil {
		t.Fatalf("unmarshal Preference: %v", err)
	}
	if p.SelectedOrganizationID != testSelectedOrgID {
		t.Errorf("SelectedOrganizationID = %q, want the canonical %q", p.SelectedOrganizationID, testSelectedOrgID)
	}
}

// TestPreferenceOmittedOrganizationKeyPreservesSelection pins that a body
// carrying neither spelling leaves the persisted selection alone - the merge
// path decodes partial bodies (an anonymous-stats toggle, say) onto the full
// persisted preference.
func TestPreferenceOmittedOrganizationKeyPreservesSelection(t *testing.T) {
	p := &Preference{SelectedOrganizationID: testSelectedOrgID}

	if err := json.Unmarshal([]byte(`{"anonymousUsageStats":true}`), p); err != nil {
		t.Fatalf("unmarshal Preference: %v", err)
	}
	if p.SelectedOrganizationID != testSelectedOrgID {
		t.Errorf("SelectedOrganizationID = %q, want the persisted %q preserved", p.SelectedOrganizationID, testSelectedOrgID)
	}
	if !p.AnonymousUsageStats {
		t.Error("AnonymousUsageStats was not applied from the partial body")
	}
}

// TestPreferenceRoundTripsEveryField guards the custom UnmarshalJSON: aliasing
// the struct to avoid infinite recursion is easy to get subtly wrong, and a
// dropped field would be invisible to the focused tests above.
func TestPreferenceRoundTripsEveryField(t *testing.T) {
	want := &Preference{
		AnonymousUsageStats:               true,
		AnonymousPerfResults:              true,
		DashboardPreferences:              map[string]interface{}{"theme": "dark"},
		SelectedOrganizationID:            testSelectedOrgID,
		SelectedWorkspaceForOrganizations: map[string]string{testSelectedOrgID: "ws-1"},
		UsersExtensionPreferences:         map[string]interface{}{"catalogContent": true},
		RemoteProviderPreferences:         map[string]interface{}{"k": "v"},
		LoadTestPreferences:               &LoadTestPreferences{ConcurrentRequests: 2, QueriesPerSecond: 10},
	}

	b, err := json.Marshal(want)
	if err != nil {
		t.Fatalf("marshal Preference: %v", err)
	}

	var got Preference
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal Preference: %v", err)
	}

	if got.SelectedOrganizationID != want.SelectedOrganizationID {
		t.Errorf("SelectedOrganizationID = %q, want %q", got.SelectedOrganizationID, want.SelectedOrganizationID)
	}
	if !got.AnonymousUsageStats || !got.AnonymousPerfResults {
		t.Errorf("anonymous opt-ins were dropped: %+v", got)
	}
	if got.DashboardPreferences["theme"] != "dark" {
		t.Errorf("DashboardPreferences = %v, want theme=dark", got.DashboardPreferences)
	}
	if got.SelectedWorkspaceForOrganizations[testSelectedOrgID] != "ws-1" {
		t.Errorf("SelectedWorkspaceForOrganizations = %v", got.SelectedWorkspaceForOrganizations)
	}
	if got.UsersExtensionPreferences["catalogContent"] != true {
		t.Errorf("UsersExtensionPreferences = %v", got.UsersExtensionPreferences)
	}
	if got.RemoteProviderPreferences["k"] != "v" {
		t.Errorf("RemoteProviderPreferences = %v", got.RemoteProviderPreferences)
	}
	if got.LoadTestPreferences == nil ||
		got.LoadTestPreferences.ConcurrentRequests != 2 ||
		got.LoadTestPreferences.QueriesPerSecond != 10 {
		t.Errorf("LoadTestPreferences = %+v", got.LoadTestPreferences)
	}
}
