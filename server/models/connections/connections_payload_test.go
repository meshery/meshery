package connections

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/core"
	modelv1beta1 "github.com/meshery/schemas/models/v1beta1/model"
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

			payload := &ConnectionPayload{Name: "cluster-2"}
			MergePayloadOntoExisting(payload, nil)
			assert.Equal(t, "cluster-2", payload.Name)

			MergePayloadOntoExisting(nil, nil)
		})
	})

	t.Run("preserves existing fields when omitted in payload", func(t *testing.T) {
		credID := core.Uuid(uuid.Must(uuid.NewV4()))
		existingMeta := map[string]interface{}{"env": "production"}
		desc := "production monitoring instance"
		targetURL := "https://prometheus.internal:9090"
		modelRef := &modelv1beta1.ModelReference{Name: "prometheus-model"}
		connSchema := core.Map{"type": "object", "title": "PromConn"}
		credSchema := core.Map{"type": "object", "title": "PromCred"}

		existing := &Connection{
			Name:             "persisted-cluster",
			Kind:             "prometheus",
			ConnectionType:   "observability",
			SubType:          "monitoring",
			Status:           CONNECTED,
			CredentialID:     &credID,
			Metadata:         existingMeta,
			Description:      &desc,
			Url:              &targetURL,
			ModelReference:   modelRef,
			ConnectionSchema: connSchema,
			CredentialSchema: credSchema,
		}

		// Issue #23907: Partial PUT update sending only status must not wipe
		// description, url, model reference, schemas or legacy connection fields.
		payload := &ConnectionPayload{
			Status: DISCONNECTED,
		}

		MergePayloadOntoExisting(payload, existing)

		assert.Equal(t, "persisted-cluster", payload.Name)
		assert.Equal(t, "prometheus", payload.Kind)
		assert.Equal(t, "observability", payload.Type)
		assert.Equal(t, "monitoring", payload.SubType)
		assert.Equal(t, DISCONNECTED, payload.Status)
		assert.Equal(t, &credID, payload.CredentialID)
		assert.Equal(t, existingMeta, payload.MetaData)
		assert.Equal(t, &desc, payload.Description)
		assert.Equal(t, &targetURL, payload.Url)
		assert.Equal(t, modelRef, payload.ModelReference)
		assert.Equal(t, connSchema, payload.ConnectionSchema)
		assert.Equal(t, credSchema, payload.CredentialSchema)
	})

	t.Run("does not overwrite explicitly supplied fields", func(t *testing.T) {
		credID1 := core.Uuid(uuid.Must(uuid.NewV4()))
		credID2 := core.Uuid(uuid.Must(uuid.NewV4()))
		meta1 := map[string]interface{}{"env": "production"}
		meta2 := map[string]interface{}{"env": "staging"}
		desc1 := "old description"
		desc2 := "new description"
		url1 := "https://old.url"
		url2 := "https://new.url"
		modelRef1 := &modelv1beta1.ModelReference{Name: "old-model"}
		modelRef2 := &modelv1beta1.ModelReference{Name: "new-model"}
		connSchema1 := core.Map{"version": "1.0"}
		connSchema2 := core.Map{"version": "2.0"}
		credSchema1 := core.Map{"auth": "token"}
		credSchema2 := core.Map{"auth": "basic"}

		existing := &Connection{
			Name:             "persisted-cluster",
			Kind:             "prometheus",
			ConnectionType:   "observability",
			SubType:          "monitoring",
			Status:           CONNECTED,
			CredentialID:     &credID1,
			Metadata:         meta1,
			Description:      &desc1,
			Url:              &url1,
			ModelReference:   modelRef1,
			ConnectionSchema: connSchema1,
			CredentialSchema: credSchema1,
		}

		payload := &ConnectionPayload{
			Name:             "renamed-cluster",
			Kind:             "grafana",
			Type:             "telemetry",
			SubType:          "dashboards",
			Status:           DISCONNECTED,
			CredentialID:     &credID2,
			MetaData:         meta2,
			Description:      &desc2,
			Url:              &url2,
			ModelReference:   modelRef2,
			ConnectionSchema: connSchema2,
			CredentialSchema: credSchema2,
		}

		MergePayloadOntoExisting(payload, existing)

		assert.Equal(t, "renamed-cluster", payload.Name)
		assert.Equal(t, "grafana", payload.Kind)
		assert.Equal(t, "telemetry", payload.Type)
		assert.Equal(t, "dashboards", payload.SubType)
		assert.Equal(t, DISCONNECTED, payload.Status)
		assert.Equal(t, &credID2, payload.CredentialID)
		assert.Equal(t, meta2, payload.MetaData)
		assert.Equal(t, &desc2, payload.Description)
		assert.Equal(t, &url2, payload.Url)
		assert.Equal(t, modelRef2, payload.ModelReference)
		assert.Equal(t, connSchema2, payload.ConnectionSchema)
		assert.Equal(t, credSchema2, payload.CredentialSchema)
	})
}
