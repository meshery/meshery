package models

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestMesheryPatternMarshalsSchemasOwnerContract pins the design-owner wire
// contract to the canonical schemas v1beta3 design.MesheryPattern shape: the
// server must emit "userId" (and "user" when the owner profile is joined) and
// must never emit the legacy "owner" key that the schema-generated UI client no
// longer reads. Emitting "owner" is exactly what fed the design Info modal
// "Owner: undefined undefined" bug, so this is its regression guard.
func TestMesheryPatternMarshalsSchemasOwnerContract(t *testing.T) {
	id := LocalProviderUserID
	p := MesheryPattern{
		Name:   "demo",
		UserID: &id,
		User:   LocalProviderContentUser(),
	}

	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal MesheryPattern: %v", err)
	}
	got := string(b)

	if !strings.Contains(got, `"userId"`) {
		t.Errorf(`marshaled pattern is missing the "userId" key: %s`, got)
	}
	if !strings.Contains(got, `"user"`) {
		t.Errorf(`marshaled pattern is missing the joined "user" key: %s`, got)
	}
	if strings.Contains(got, `"owner"`) {
		t.Errorf(`marshaled pattern still emits the legacy "owner" key: %s`, got)
	}

	// The embedded owner profile must carry a resolvable name so the UI can
	// render the owner without a separate lookup and never shows "undefined".
	var decoded struct {
		UserID string `json:"userId"`
		User   struct {
			ID        string `json:"id"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		} `json:"user"`
	}
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal marshaled pattern: %v", err)
	}
	if decoded.UserID != LocalProviderUserID.String() {
		t.Errorf("userId = %q, want %q", decoded.UserID, LocalProviderUserID.String())
	}
	if decoded.User.FirstName == "" || decoded.User.LastName == "" {
		t.Errorf("embedded user name is empty: %+v", decoded.User)
	}
}

// TestMesheryPatternOmitsOwnerWhenUnset ensures the userId/user keys are omitted
// (omitempty) when no owner is populated, rather than emitting a null owner or
// resurrecting the legacy "owner" key.
func TestMesheryPatternOmitsOwnerWhenUnset(t *testing.T) {
	b, err := json.Marshal(MesheryPattern{Name: "demo"})
	if err != nil {
		t.Fatalf("marshal MesheryPattern: %v", err)
	}
	got := string(b)

	if strings.Contains(got, `"owner"`) {
		t.Errorf(`unset pattern must not emit the legacy "owner" key: %s`, got)
	}
	if strings.Contains(got, `"userId"`) {
		t.Errorf(`unset pattern must omit "userId" (omitempty): %s`, got)
	}
	if strings.Contains(got, `"user"`) {
		t.Errorf(`unset pattern must omit "user" (omitempty): %s`, got)
	}
}

// TestStampLocalProviderOwner verifies the persister stamps the built-in
// provider's single user as the design owner on read. UserID/User are gorm:"-"
// and never loaded from the DB, so the stamp is what makes a locally stored
// design resolve to the local "meshery" user in the UI.
func TestStampLocalProviderOwner(t *testing.T) {
	p := &MesheryPattern{Name: "demo"}
	stampLocalProviderOwner(p)

	if p.UserID == nil || *p.UserID != LocalProviderUserID {
		t.Fatalf("UserID = %v, want %v", p.UserID, LocalProviderUserID)
	}
	if p.User == nil {
		t.Fatal("User profile was not stamped")
	}
	if p.User.ID != LocalProviderUserID {
		t.Errorf("User.ID = %v, want %v", p.User.ID, LocalProviderUserID)
	}
	if p.User.FirstName != "Meshery" || p.User.LastName != "Meshery" {
		t.Errorf("User name = %q %q, want %q %q", p.User.FirstName, p.User.LastName, "Meshery", "Meshery")
	}

	// A nil receiver must be a safe no-op, not a panic.
	stampLocalProviderOwner(nil)
}

// TestStampLocalProviderOwnerSkipsPublished pins that seeded, community-authored
// catalog content is never attributed to the built-in provider's user.
// SeedContent imports docs/data/catalog with Visibility: Published; stamping
// those would both misattribute third-party work and - because the local
// current-user id is the same LocalProviderUserID - hand the local user the
// owner-gated edit and visibility affordances over the whole catalog.
func TestStampLocalProviderOwnerSkipsPublished(t *testing.T) {
	p := &MesheryPattern{Name: "seeded catalog design", Visibility: Published}
	stampLocalProviderOwner(p)

	if p.UserID != nil {
		t.Errorf("published design was stamped with UserID %v, want none", p.UserID)
	}
	if p.User != nil {
		t.Errorf("published design was stamped with a user profile %+v, want none", p.User)
	}

	// The design must not appear owned by the local current user, which is what
	// the UI's canEditDesign / canChangeVisibility gates compare against.
	b, err := json.Marshal(p)
	if err != nil {
		t.Fatalf("marshal MesheryPattern: %v", err)
	}
	if strings.Contains(string(b), LocalProviderUserID.String()) {
		t.Errorf("published design leaks the local user id: %s", b)
	}
}

// TestStampLocalProviderOwnerStampsPrivateAndPublic verifies the owner is still
// stamped for every visibility the local user genuinely owns.
func TestStampLocalProviderOwnerStampsPrivateAndPublic(t *testing.T) {
	for _, visibility := range []string{"", Private, Public} {
		p := &MesheryPattern{Name: "demo", Visibility: visibility}
		stampLocalProviderOwner(p)

		if p.UserID == nil || *p.UserID != LocalProviderUserID {
			t.Errorf("visibility %q: UserID = %v, want %v", visibility, p.UserID, LocalProviderUserID)
		}
		if p.User == nil {
			t.Errorf("visibility %q: user profile was not stamped", visibility)
		}
	}
}

// TestStampLocalProviderFilterOwner pins the same owner rule for filters. The
// canonical wire key for a filter's owner is "owner" (schemas v1beta3
// filter.MesheryFilter, and what meshery-cloud emits), not the design contract's
// "userId" - a locally invented key is what left the filter Info modal unable to
// resolve an owner at all.
func TestStampLocalProviderFilterOwner(t *testing.T) {
	f := &MesheryFilter{Name: "demo"}
	stampLocalProviderFilterOwner(f)

	if f.Owner == nil || *f.Owner != LocalProviderUserID.String() {
		t.Fatalf("Owner = %v, want %v", f.Owner, LocalProviderUserID)
	}

	b, err := json.Marshal(f)
	if err != nil {
		t.Fatalf("marshal MesheryFilter: %v", err)
	}
	got := string(b)
	if !strings.Contains(got, `"owner"`) {
		t.Errorf(`marshaled filter is missing the "owner" key: %s`, got)
	}
	if strings.Contains(got, `"userId"`) {
		t.Errorf(`marshaled filter must not invent a local "userId" key: %s`, got)
	}

	published := &MesheryFilter{Name: "seeded", Visibility: Published}
	stampLocalProviderFilterOwner(published)
	if published.Owner != nil {
		t.Errorf("published filter was stamped with Owner %v, want none", published.Owner)
	}

	// A nil receiver must be a safe no-op, not a panic.
	stampLocalProviderFilterOwner(nil)
}

// TestLocalProviderUserIdentityIsSingleSourced pins that the owner profile
// joined onto content resources is the same identity /api/user reports. The two
// used to be hand-copied, so renaming the local user in one place silently
// desynchronised the other with no compile or test failure.
func TestLocalProviderUserIdentityIsSingleSourced(t *testing.T) {
	u := LocalProviderUser()
	c := LocalProviderContentUser()

	if u.ID != c.ID {
		t.Errorf("id: /api/user = %v, content owner = %v", u.ID, c.ID)
	}
	if u.FirstName != c.FirstName || u.LastName != c.LastName {
		t.Errorf("name: /api/user = %q %q, content owner = %q %q", u.FirstName, u.LastName, c.FirstName, c.LastName)
	}
	if u.Email != c.Email {
		t.Errorf("email: /api/user = %q, content owner = %q", u.Email, c.Email)
	}
}
