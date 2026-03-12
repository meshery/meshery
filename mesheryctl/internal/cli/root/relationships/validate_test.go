// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package relationships

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateRelationshipFile_Valid(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "relationship", result.EntityType)
	assert.Equal(t, "edge/mount", result.EntityName)
	assert.Empty(t, result.Errors)
}

func TestValidateRelationshipFile_ValidHierarchical(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "hierarchical",
		"type": "parent",
		"subType": "inventory",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "hierarchical/inventory", result.EntityName)
}

func TestValidateRelationshipFile_ValidSibling(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "sibling",
		"type": "sibling",
		"subType": "network",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.True(t, result.IsValid)
}

func TestValidateRelationshipFile_NotRelationship(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "models.meshery.io/v1beta1",
		"name": "test-model"
	}`)
	result := validateRelationshipFile("not-rel.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors[0], "file is not a relationship definition")
}

func TestValidateRelationshipFile_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)
	result := validateRelationshipFile("bad.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors[0], "file is not a relationship definition")
}

func TestValidateRelationshipFile_InvalidKind(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "invalid-kind",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "invalid kind")
	assert.Contains(t, result.Errors[0], "must be one of [edge, hierarchical, sibling]")
}

func TestValidateRelationshipFile_MissingRequiredFields(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "kind is required")
	assert.Contains(t, result.Errors, "type is required")
	assert.Contains(t, result.Errors, "subType is required")
	assert.Contains(t, result.Errors, "version is required")
}

func TestValidateRelationshipFile_MissingSchemaVersion(t *testing.T) {
	data := []byte(`{
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("relationship.json", data)
	assert.False(t, result.IsValid)
	// Missing schemaVersion means DetectEntityType returns "" so this should
	// fail with "not a relationship definition"
	assert.Contains(t, result.Errors[0], "file is not a relationship definition")
}

func TestValidateRelationshipFile_NameFallbackToFilename(t *testing.T) {
	// When kind and subType are both empty, name should fall back to filename
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"type": "binding",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipFile("my-relationship.json", data)
	assert.False(t, result.IsValid)
	// kind is missing so name derivation yields "" → fallback to filename
	assert.Equal(t, "my-relationship.json", result.EntityName)
}
