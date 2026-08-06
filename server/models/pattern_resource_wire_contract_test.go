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
