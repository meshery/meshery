package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	patternresource "github.com/meshery/schemas/models/v1beta3/pattern_resource"
)

// TestPatternResourceWireContractMatchesSchemas pins the pattern-resource owner
// to the canonical `userId` declared by the schemas v1beta3
// pattern_resource.MesheryPatternResource construct.
//
// PatternResource is both the built-in provider's GORM model and the body
// RemoteProvider.SaveMesheryPatternResource sends to - and decodes back from -
// meshery-cloud, which reads and writes `userId`. A local rename to `owner`
// broke both directions silently: meshery sent an owner the provider ignored,
// and read back an owner the provider never sent.
func TestPatternResourceWireContractMatchesSchemas(t *testing.T) {
	id := core.Uuid(uuid.Must(uuid.FromString("0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c40")))
	owner := core.Uuid(uuid.Must(uuid.FromString("0195b0ab-1f4d-7a3c-9c1e-3a5f8d2b6c41")))

	b, err := json.Marshal(&PatternResource{ID: &id, UserID: &owner, Name: "demo"})
	if err != nil {
		t.Fatalf("marshal PatternResource: %v", err)
	}
	got := string(b)

	if !strings.Contains(got, `"userId"`) {
		t.Errorf(`marshaled pattern resource is missing the canonical "userId" key: %s`, got)
	}
	if strings.Contains(got, `"owner"`) {
		t.Errorf(`marshaled pattern resource still emits the legacy "owner" key: %s`, got)
	}

	// The provider's reply must round trip back onto UserID.
	var decoded PatternResource
	if err := json.Unmarshal(b, &decoded); err != nil {
		t.Fatalf("unmarshal PatternResource: %v", err)
	}
	if decoded.UserID == nil || *decoded.UserID != owner {
		t.Errorf("UserID = %v, want %v", decoded.UserID, owner)
	}
}

// TestPatternResourceJSONKeysMatchSchemasConstruct compares every JSON key this
// struct emits against the schemas construct it implements. A field renamed on
// either side shows up here rather than as a silently dropped value.
func TestPatternResourceJSONKeysMatchSchemasConstruct(t *testing.T) {
	jsonKeys := func(v interface{}) map[string]bool {
		typ := reflect.TypeOf(v)
		keys := map[string]bool{}
		for i := 0; i < typ.NumField(); i++ {
			tag := typ.Field(i).Tag.Get("json")
			if tag == "" || tag == "-" {
				continue
			}
			keys[strings.Split(tag, ",")[0]] = true
		}
		return keys
	}

	local := jsonKeys(PatternResource{})
	canonical := jsonKeys(patternresource.MesheryPatternResource{})

	for key := range local {
		if !canonical[key] {
			t.Errorf("PatternResource emits %q, which the schemas v1beta3 construct does not declare", key)
		}
	}
}

// TestPatternResourceOwnerColumnIsPinned guards the GORM side of the rename.
// The struct doubles as the built-in provider's model and the schemas Go models
// carry `db:` tags GORM does not read, so the owner column has to be pinned
// explicitly - without the tag GORM would derive `user_id` from the Go field
// name and stop finding the existing `owner` column.
func TestPatternResourceOwnerColumnIsPinned(t *testing.T) {
	field, ok := reflect.TypeOf(PatternResource{}).FieldByName("UserID")
	if !ok {
		t.Fatal("PatternResource has no UserID field")
	}
	if got := field.Tag.Get("gorm"); got != "column:owner" {
		t.Errorf(`UserID gorm tag = %q, want "column:owner"`, got)
	}
}

// TestPatternResourceRoundTripsOwnerThroughMigratedDB is the runtime half of the
// column guard: the tag assertion above proves the annotation is present, not
// that AutoMigrate and the persister agree on it.
//
// The failure this catches is the one documented for PerformanceProfile - gorm
// derives a column from the *field name* via its naming strategy, so renaming
// Owner to UserID renames `owner` to `user_id` unless the explicit column tag
// holds. The persister discards gorm's errors (see GetPatternResources), so a
// drifted column reads back as an empty list with nothing logged.
func TestPatternResourceRoundTripsOwnerThroughMigratedDB(t *testing.T) {
	db := newMigratedDB(t)

	cols, err := db.Migrator().ColumnTypes(&PatternResource{})
	if err != nil {
		t.Fatalf("read migrated columns: %v", err)
	}
	names := map[string]bool{}
	for _, c := range cols {
		names[c.Name()] = true
	}
	if !names["owner"] {
		t.Errorf("migrated pattern_resources has no `owner` column; columns: %v", names)
	}
	if names["user_id"] {
		t.Errorf("migrated pattern_resources gained a `user_id` column; the gorm column tag is not holding: %v", names)
	}

	owner := mustUUID(t)
	persister := &PatternResourcePersister{DB: db}

	saved, err := persister.SavePatternResource(&PatternResource{
		UserID:    &owner,
		Name:      "demo-resource",
		Namespace: "default",
		Type:      "Deployment",
		OAMType:   "workload",
	})
	if err != nil {
		t.Fatalf("SavePatternResource: %v", err)
	}

	byID, err := persister.GetPatternResource(*saved.ID)
	if err != nil {
		t.Fatalf("GetPatternResource: %v", err)
	}
	if byID.UserID == nil || *byID.UserID != owner {
		t.Errorf("owner did not survive the round trip: got %v, want %v", byID.UserID, owner)
	}

	page, err := persister.GetPatternResources("", "", "demo-resource", "default", "Deployment", "workload", 0, 10)
	if err != nil {
		t.Fatalf("GetPatternResources: %v", err)
	}
	if page.TotalCount != 1 || len(page.Resources) != 1 {
		t.Fatalf("listing returned %d resources (totalCount=%d); a drifted column reads back empty", len(page.Resources), page.TotalCount)
	}
	if got := page.Resources[0].UserID; got == nil || *got != owner {
		t.Errorf("listed owner = %v, want %v", got, owner)
	}
}
