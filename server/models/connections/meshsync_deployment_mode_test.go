package connections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMeshsyncDeploymentModeFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected MeshsyncDeploymentMode
	}{
		{name: "empty string", input: "", expected: MeshsyncDeploymentModeUndefined},
		{name: "operator", input: string(MeshsyncDeploymentModeOperator), expected: MeshsyncDeploymentModeOperator},
		{name: "embedded", input: string(MeshsyncDeploymentModeEmbedded), expected: MeshsyncDeploymentModeEmbedded},
		{name: "unknown value", input: "unknown", expected: MeshsyncDeploymentModeUndefined},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MeshsyncDeploymentModeFromString(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMeshsyncDeploymentModeFromMetadata(t *testing.T) {
	tests := []struct {
		name     string
		metadata map[string]any
		expected MeshsyncDeploymentMode
	}{
		{
			name:     "nil metadata",
			metadata: nil,
			expected: MeshsyncDeploymentModeUndefined,
		},
		{
			name:     "no key",
			metadata: map[string]any{},
			expected: MeshsyncDeploymentModeUndefined,
		},
		{
			name:     "wrong type",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: 123},
			expected: MeshsyncDeploymentModeUndefined,
		},
		{
			name:     "empty string",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: ""},
			expected: MeshsyncDeploymentModeUndefined,
		},
		{
			name:     "operator mode string",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)},
			expected: MeshsyncDeploymentModeOperator,
		},
		{
			name:     "embedded mode string",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeEmbedded)},
			expected: MeshsyncDeploymentModeEmbedded,
		},
		{
			name:     "unknown string",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: "something-else"},
			expected: MeshsyncDeploymentModeUndefined,
		},
		{
			name:     "direct MeshsyncDeploymentMode type",
			metadata: map[string]any{MeshsyncDeploymentModeMetadataKey: MeshsyncDeploymentModeOperator},
			expected: MeshsyncDeploymentModeOperator,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MeshsyncDeploymentModeFromMetadata(tt.metadata)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAddMeshsyncDeploymentModeToMetadata(t *testing.T) {
	t.Run("adds mode to empty metadata map", func(t *testing.T) {
		metadata := make(map[string]any)
		SetMeshsyncDeploymentModeToMetadata(metadata, MeshsyncDeploymentModeEmbedded)

		assert.Contains(t, metadata, MeshsyncDeploymentModeMetadataKey)
		assert.Equal(t, MeshsyncDeploymentModeEmbedded, metadata[MeshsyncDeploymentModeMetadataKey])
	})

	t.Run("overwrites existing mode in metadata", func(t *testing.T) {
		metadata := map[string]any{
			MeshsyncDeploymentModeMetadataKey: MeshsyncDeploymentModeOperator,
		}
		SetMeshsyncDeploymentModeToMetadata(metadata, MeshsyncDeploymentModeEmbedded)

		assert.Equal(t, MeshsyncDeploymentModeEmbedded, metadata[MeshsyncDeploymentModeMetadataKey])
	})

	t.Run("nil metadata is a no-op, not a panic", func(t *testing.T) {
		var metadata map[string]any // nil map
		assert.NotPanics(t, func() {
			SetMeshsyncDeploymentModeToMetadata(metadata, MeshsyncDeploymentModeEmbedded)
		})
	})
}
