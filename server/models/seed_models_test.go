package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/schemas/models/v1beta3/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMesheryCoreContainerComponent_NotAnnotation verifies that the meshery-core
// Container component definition in the shipped models registry is classified with
// isAnnotation: false.
//
// Background: Container components represent workload spec.containers and are linked
// to workload controllers (Deployment, Pod, etc.) via hierarchical-parent-alias relationships.
// Flagging Container as an annotation causes Kanvas to incorrectly hide containers when
// the annotations layer is toggled off (issue #22780).
//
// Note: this guards the component-level metadata.isAnnotation. The same file also
// carries model.metadata.isAnnotation (the embedded parent-model metadata), which is a
// separate path and is already false.
func TestMesheryCoreContainerComponent_NotAnnotation(t *testing.T) {
	containerPath := filepath.Join(mesheryCoreModelDir(t), "components", "Container.json")

	data, err := os.ReadFile(containerPath)
	require.NoError(t, err, "Container.json component definition must exist in meshery-core model")

	var compDef component.ComponentDefinition
	err = json.Unmarshal(data, &compDef)
	require.NoError(t, err, "Container.json must unmarshal into component.ComponentDefinition")

	assert.False(t, compDef.Metadata.IsAnnotation,
		"meshery-core Container component represents infrastructure workload spec.containers and must have isAnnotation: false (issue #22780)")
}
