package connections

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// TestBuildMesheryConnectionPayloadEmitsCanonicalCamelCaseKeys locks in
// the canonical camelCase wire-format contract for the four
// server-* JSONB metadata keys flipped in Phase 5c. The downstream
// Remote Provider DAO (meshery-cloud) accepts both spellings during
// the rollout window, but the upstream emitter must emit only the
// canonical form. Regressing this back to snake_case would be a wire
// contract violation per the identifier-naming overhaul.
func TestBuildMesheryConnectionPayloadEmitsCanonicalCamelCaseKeys(t *testing.T) {
	viper.Set("INSTANCE_ID", "instance-123")
	viper.Set("BUILD", "v0.0.0-test")
	viper.Set("COMMITSHA", "deadbeef")
	t.Cleanup(func() {
		viper.Set("INSTANCE_ID", "")
		viper.Set("BUILD", "")
		viper.Set("COMMITSHA", "")
	})

	payload := BuildMesheryConnectionPayload("https://example.test", nil)

	assert.Equal(t, "instance-123", payload.MetaData["serverId"])
	assert.Equal(t, "v0.0.0-test", payload.MetaData["serverVersion"])
	assert.Equal(t, "deadbeef", payload.MetaData["serverBuildSha"])
	assert.Equal(t, "https://example.test", payload.MetaData["serverLocation"])

	for _, legacy := range []string{"server_id", "server_version", "server_build_sha", "server_location"} {
		_, present := payload.MetaData[legacy]
		assert.Falsef(t, present, "legacy snake_case key %q must not be emitted", legacy)
	}

	assert.Equal(t, "meshery", payload.Kind)
	assert.Equal(t, "platform", payload.Type)
	assert.Equal(t, "management", payload.SubType)
	assert.Equal(t, CONNECTED, payload.Status)
}

func TestMergePayloadOntoExisting(t *testing.T) {
	t.Run("nil safety", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MergePayloadOntoExisting(nil, &Connection{Name: "cluster-1"})
			MergePayloadOntoExisting(&ConnectionPayload{Name: "cluster-2"}, nil)
			MergePayloadOntoExisting(nil, nil)
		})
	})

	t.Run("preserves existing fields when omitted in payload", func(t *testing.T) {
		existing := &Connection{
			Name:           "persisted-cluster",
			Kind:           "kubernetes",
			ConnectionType: "platform",
			SubType:        "k8s",
			Status:         CONNECTED,
		}

		// Issue #8716: Partial PUT update sending only status must not wipe name or kind
		payload := &ConnectionPayload{
			Status: DISCONNECTED,
		}

		MergePayloadOntoExisting(payload, existing)

		assert.Equal(t, "persisted-cluster", payload.Name)
		assert.Equal(t, "kubernetes", payload.Kind)
		assert.Equal(t, "platform", payload.Type)
		assert.Equal(t, "k8s", payload.SubType)
		assert.Equal(t, DISCONNECTED, payload.Status)
	})

	t.Run("does not overwrite explicitly supplied fields", func(t *testing.T) {
		existing := &Connection{
			Name:           "persisted-cluster",
			Kind:           "kubernetes",
			ConnectionType: "platform",
			SubType:        "k8s",
			Status:         CONNECTED,
		}

		payload := &ConnectionPayload{
			Name:    "renamed-cluster",
			Kind:    "custom-kind",
			Type:    "custom-type",
			SubType: "custom-sub",
			Status:  DISCONNECTED,
		}

		MergePayloadOntoExisting(payload, existing)

		assert.Equal(t, "renamed-cluster", payload.Name)
		assert.Equal(t, "custom-kind", payload.Kind)
		assert.Equal(t, "custom-type", payload.Type)
		assert.Equal(t, "custom-sub", payload.SubType)
		assert.Equal(t, DISCONNECTED, payload.Status)
	})
}
