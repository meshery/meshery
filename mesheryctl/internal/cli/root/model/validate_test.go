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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateEntityData_UnknownType(t *testing.T) {
	data := []byte(`{"name":"test"}`)
	result := validateEntityData("test.json", data)
	assert.False(t, result.IsValid)
	assert.Equal(t, "unknown", result.EntityType)
	assert.Equal(t, "test.json", result.EntityName)
	assert.Contains(t, result.Errors[0], "unable to determine entity type")
}

func TestValidateModelDef_Valid(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "models.meshery.io/v1beta1",
		"name": "test-model",
		"version": "v0.1.0"
	}`)
	result := validateModelDef("model.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "model", result.EntityType)
	assert.Equal(t, "test-model", result.EntityName)
	assert.Empty(t, result.Errors)
}

func TestValidateModelDef_MissingName(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "models.meshery.io/v1beta1",
		"version": "v0.1.0"
	}`)
	result := validateModelDef("model.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "name is required")
}

func TestValidateModelDef_MissingSchemaVersion(t *testing.T) {
	data := []byte(`{
		"name": "test-model",
		"version": "v0.1.0"
	}`)
	result := validateModelDef("model.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "schemaVersion is required")
}

func TestValidateModelDef_MissingVersion(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "models.meshery.io/v1beta1",
		"name": "test-model"
	}`)
	result := validateModelDef("model.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "version is required")
}

func TestValidateModelDef_MissingAllFields(t *testing.T) {
	data := []byte(`{}`)
	result := validateModelDef("model.json", data)
	assert.False(t, result.IsValid)
	assert.Len(t, result.Errors, 3)
	assert.Contains(t, result.Errors, "name is required")
	assert.Contains(t, result.Errors, "schemaVersion is required")
	assert.Contains(t, result.Errors, "version is required")
	// When name is empty, falls back to filename
	assert.Equal(t, "model.json", result.EntityName)
}

func TestValidateModelDef_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)
	result := validateModelDef("model.json", data)
	assert.False(t, result.IsValid)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "failed to parse model definition")
}

func TestValidateComponentDef_Valid(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"displayName": "TestComponent",
		"component": {
			"kind": "Deployment",
			"version": "apps/v1"
		}
	}`)
	result := validateComponentDef("component.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "component", result.EntityType)
	assert.Equal(t, "TestComponent", result.EntityName)
	assert.Empty(t, result.Errors)
}

func TestValidateComponentDef_MissingKind(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"component": {
			"version": "apps/v1"
		}
	}`)
	result := validateComponentDef("component.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "component.kind is required")
}

func TestValidateComponentDef_MissingComponentVersion(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"component": {
			"kind": "Deployment"
		}
	}`)
	result := validateComponentDef("component.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "component.version is required")
}

func TestValidateComponentDef_FallbackNameToKind(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"component": {
			"kind": "Deployment",
			"version": "apps/v1"
		}
	}`)
	result := validateComponentDef("component.json", data)
	assert.True(t, result.IsValid)
	// When displayName is empty, falls back to component.kind
	assert.Equal(t, "Deployment", result.EntityName)
}

func TestValidateComponentDef_FallbackNameToFilename(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"component": {
			"version": "apps/v1"
		}
	}`)
	result := validateComponentDef("component.json", data)
	// When both displayName and kind are empty, falls back to filename
	assert.Equal(t, "component.json", result.EntityName)
}

func TestValidateComponentDef_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)
	result := validateComponentDef("component.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors[0], "failed to parse component definition")
}

func TestValidateRelationshipDef_Valid(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "relationship", result.EntityType)
	assert.Equal(t, "edge/mount", result.EntityName)
	assert.Empty(t, result.Errors)
}

func TestValidateRelationshipDef_ValidHierarchical(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "hierarchical",
		"type": "parent",
		"subType": "inventory",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "hierarchical/inventory", result.EntityName)
}

func TestValidateRelationshipDef_ValidSibling(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "sibling",
		"type": "sibling",
		"subType": "network",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.True(t, result.IsValid)
}

func TestValidateRelationshipDef_InvalidKind(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "invalid-kind",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0], "invalid kind")
	assert.Contains(t, result.Errors[0], "must be one of [edge, hierarchical, sibling]")
}

func TestValidateRelationshipDef_MissingKind(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "kind is required")
}

func TestValidateRelationshipDef_MissingType(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "type is required")
}

func TestValidateRelationshipDef_MissingSubType(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"type": "binding",
		"version": "v0.0.1"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "subType is required")
}

func TestValidateRelationshipDef_MissingVersion(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"type": "binding",
		"subType": "mount"
	}`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors, "version is required")
}

func TestValidateRelationshipDef_InvalidJSON(t *testing.T) {
	data := []byte(`not json`)
	result := validateRelationshipDef("relationship.json", data)
	assert.False(t, result.IsValid)
	assert.Contains(t, result.Errors[0], "failed to parse relationship definition")
}

func TestValidateEntityData_RoutesToModel(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "models.meshery.io/v1beta1",
		"name": "test-model",
		"version": "v0.1.0"
	}`)
	result := validateEntityData("model.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "model", result.EntityType)
}

func TestValidateEntityData_RoutesToComponent(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "components.meshery.io/v1beta1",
		"displayName": "TestComp",
		"component": {
			"kind": "Deployment",
			"version": "apps/v1"
		}
	}`)
	result := validateEntityData("component.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "component", result.EntityType)
}

func TestValidateEntityData_RoutesToRelationship(t *testing.T) {
	data := []byte(`{
		"schemaVersion": "relationships.meshery.io/v1alpha3",
		"kind": "edge",
		"type": "binding",
		"subType": "mount",
		"version": "v0.0.1"
	}`)
	result := validateEntityData("relationship.json", data)
	assert.True(t, result.IsValid)
	assert.Equal(t, "relationship", result.EntityType)
}
