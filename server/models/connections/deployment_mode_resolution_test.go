package connections

import (
	"testing"

	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func modeDoc(mode MeshsyncDeploymentMode) *controllersconfig.MesheryControllersConfig {
	value := controllersconfig.MesheryOperatorConfigDeploymentMode(mode)
	return &controllersconfig.MesheryControllersConfig{
		Operator: &controllersconfig.MesheryOperatorConfig{DeploymentMode: &value},
	}
}

// TestB3ServerWideDeploymentModeReachesInheritingConnection pins the defect
// this precedence exists to fix: every Kubernetes connection carries a
// materialized meshsync_deployment_mode entry from registration, so ranking
// that entry above the layered document made a server-wide default change
// unreachable for every existing connection.
func TestB3ServerWideDeploymentModeReachesInheritingConnection(t *testing.T) {
	t.Run("server-wide default wins over the materialized cache", func(t *testing.T) {
		// The connection inherits (no per-connection override), so the merged
		// document carries only the server-wide default.
		merged := modeDoc(MeshsyncDeploymentModeOperator)
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeEmbedded)}

		resolved := ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeEmbedded)

		assert.Equal(t, MeshsyncDeploymentModeOperator, resolved.Mode)
		assert.Equal(t, DeploymentModeLayerLayeredConfig, resolved.Layer)
	})

	t.Run("per-connection override still wins over the server-wide default", func(t *testing.T) {
		override := modeDoc(MeshsyncDeploymentModeEmbedded)
		serverDefault := modeDoc(MeshsyncDeploymentModeOperator)
		merged, _ := ResolveControllersConfig(override, serverDefault)
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		resolved := ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeOperator)

		assert.Equal(t, MeshsyncDeploymentModeEmbedded, resolved.Mode)
		assert.Equal(t, DeploymentModeLayerLayeredConfig, resolved.Layer)
	})

	t.Run("materialized cache is the compatibility floor when no layer sets a mode", func(t *testing.T) {
		// A connection registered before the layered document existed: its
		// recorded mode must survive rather than snapping to the env default.
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		resolved := ResolveDeploymentMode(nil, metadata, MeshsyncDeploymentModeEmbedded)

		assert.Equal(t, MeshsyncDeploymentModeOperator, resolved.Mode)
		assert.Equal(t, DeploymentModeLayerLegacyMetadata, resolved.Layer)
	})

	t.Run("a document without an operator section does not shadow the cache", func(t *testing.T) {
		merged := &controllersconfig.MesheryControllersConfig{
			Operator: &controllersconfig.MesheryOperatorConfig{},
		}
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		resolved := ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeEmbedded)

		assert.Equal(t, MeshsyncDeploymentModeOperator, resolved.Mode)
		assert.Equal(t, DeploymentModeLayerLegacyMetadata, resolved.Layer)
	})

	t.Run("falls through to the server env default, then the built-in default", func(t *testing.T) {
		fromEnv := ResolveDeploymentMode(nil, core.Map{}, MeshsyncDeploymentModeOperator)
		assert.Equal(t, MeshsyncDeploymentModeOperator, fromEnv.Mode)
		assert.Equal(t, DeploymentModeLayerServerEnvDefault, fromEnv.Layer)

		builtIn := ResolveDeploymentMode(nil, nil, MeshsyncDeploymentModeUndefined)
		assert.Equal(t, MeshsyncDeploymentModeDefault, builtIn.Mode)
		assert.Equal(t, DeploymentModeLayerBuiltIn, builtIn.Layer)
	})

	t.Run("never resolves to undefined", func(t *testing.T) {
		resolved := ResolveDeploymentMode(nil, nil, MeshsyncDeploymentModeUndefined)
		assert.NotEqual(t, MeshsyncDeploymentModeUndefined, resolved.Mode)
	})
}

// TestB4WizardAndEditorShareOneDeploymentModeStore pins the convergence of the
// Connection Wizard's MeshSync Mode step and the controllers editor: both write
// the explicit choice through SetDeploymentModeOverride, so the editor can never
// report an "Override" chip for a mode the connection is not running.
func TestB4WizardAndEditorShareOneDeploymentModeStore(t *testing.T) {
	t.Run("the wizard's choice is readable as a controllers-config override", func(t *testing.T) {
		metadata := core.Map{}

		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		require.NotNil(t, override)
		assert.Equal(t, MeshsyncDeploymentModeOperator, DeploymentModeFromControllersConfig(override))
		assert.Equal(t, ControllersConfigSchemaVersion, override.SchemaVersion)
	})

	t.Run("the wizard's choice outranks a contrary server-wide default", func(t *testing.T) {
		metadata := core.Map{}
		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))
		MaterializeMeshsyncDeploymentMode(metadata, MeshsyncDeploymentModeOperator)

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		merged, _ := ResolveControllersConfig(override, modeDoc(MeshsyncDeploymentModeEmbedded))

		resolved := ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeEmbedded)
		assert.Equal(t, MeshsyncDeploymentModeOperator, resolved.Mode)
		assert.Equal(t, DeploymentModeLayerLayeredConfig, resolved.Layer)
	})

	t.Run("writing the mode preserves the rest of the override document", func(t *testing.T) {
		replicas := 3
		metadata := core.Map{}
		require.NoError(t, SetControllersConfigToMetadata(metadata, &controllersconfig.MesheryControllersConfig{
			Meshsync: &controllersconfig.MeshSyncConfig{Replicas: &replicas},
		}))

		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeEmbedded))

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		require.NotNil(t, override.Meshsync)
		require.NotNil(t, override.Meshsync.Replicas)
		assert.Equal(t, 3, *override.Meshsync.Replicas)
		assert.Equal(t, MeshsyncDeploymentModeEmbedded, DeploymentModeFromControllersConfig(override))
	})

	t.Run("writing the mode preserves a sibling operator field", func(t *testing.T) {
		version := "v0.7.0"
		metadata := core.Map{}
		require.NoError(t, SetControllersConfigToMetadata(metadata, &controllersconfig.MesheryControllersConfig{
			Operator: &controllersconfig.MesheryOperatorConfig{Version: &version},
		}))

		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		require.NotNil(t, override.Operator)
		require.NotNil(t, override.Operator.Version)
		assert.Equal(t, version, *override.Operator.Version)
		assert.Equal(t, MeshsyncDeploymentModeOperator, DeploymentModeFromControllersConfig(override))
	})

	t.Run("clearing the mode returns the connection to inheriting", func(t *testing.T) {
		metadata := core.Map{}
		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))

		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeUndefined))

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		assert.Equal(t, MeshsyncDeploymentModeUndefined, DeploymentModeFromControllersConfig(override))
		assert.NotContains(t, metadata, ControllersConfigMetadataKey,
			"an override that sets nothing must not linger as an empty document")

		merged, _ := ResolveControllersConfig(override, modeDoc(MeshsyncDeploymentModeEmbedded))
		assert.Equal(t, MeshsyncDeploymentModeEmbedded, ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeOperator).Mode)
	})

	t.Run("clearing the mode keeps a sibling operator field", func(t *testing.T) {
		version := "v0.7.0"
		metadata := core.Map{}
		require.NoError(t, SetControllersConfigToMetadata(metadata, &controllersconfig.MesheryControllersConfig{
			Operator: &controllersconfig.MesheryOperatorConfig{Version: &version},
		}))
		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))

		require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeUndefined))

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		require.NotNil(t, override.Operator)
		require.NotNil(t, override.Operator.Version)
		assert.Equal(t, version, *override.Operator.Version)
		assert.Nil(t, override.Operator.DeploymentMode)
	})

	t.Run("a malformed stored override surfaces as an error, not a silent overwrite", func(t *testing.T) {
		metadata := core.Map{ControllersConfigMetadataKey: "{not json"}

		err := SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator)

		require.Error(t, err)
		assert.Equal(t, "{not json", metadata[ControllersConfigMetadataKey])
	})

	t.Run("nil metadata is a no-op, not a panic", func(t *testing.T) {
		var metadata core.Map
		assert.NotPanics(t, func() {
			assert.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeOperator))
		})
	})
}

// TestB1EffectiveConfigReportsTheModeTheConnectionRuns pins the fact the
// controllers editor decides field applicability from. Meshery Operator manages
// MeshSync and Meshery Broker, so in embedded mode most of the document reaches
// nothing; a client that renders those fields as live (or as inert) is right
// only if the effective document names the mode the apply path actually
// resolves - including the materialized cache and the server env default, which
// plain layer merging never sees.
func TestB1EffectiveConfigReportsTheModeTheConnectionRuns(t *testing.T) {
	effectiveMode := func(cfg *controllersconfig.MesheryControllersConfig) MeshsyncDeploymentMode {
		require.NotNil(t, cfg)
		require.NotNil(t, cfg.Operator)
		require.NotNil(t, cfg.Operator.DeploymentMode)
		return MeshsyncDeploymentModeFromString(string(*cfg.Operator.DeploymentMode))
	}

	t.Run("mode carried only by the materialized cache is reported, not the built-in default", func(t *testing.T) {
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		_, effective, resolved := ResolveConnectionControllersConfig(nil, nil, metadata, MeshsyncDeploymentModeUndefined)

		assert.Equal(t, MeshsyncDeploymentModeOperator, effectiveMode(effective))
		assert.Equal(t, DeploymentModeLayerLegacyMetadata, resolved.Layer)
	})

	t.Run("mode carried only by the server env default is reported", func(t *testing.T) {
		_, effective, resolved := ResolveConnectionControllersConfig(nil, nil, core.Map{}, MeshsyncDeploymentModeOperator)

		assert.Equal(t, MeshsyncDeploymentModeOperator, effectiveMode(effective))
		assert.Equal(t, DeploymentModeLayerServerEnvDefault, resolved.Layer)
	})

	t.Run("a per-connection override still wins", func(t *testing.T) {
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		_, effective, resolved := ResolveConnectionControllersConfig(
			modeDoc(MeshsyncDeploymentModeEmbedded), modeDoc(MeshsyncDeploymentModeOperator),
			metadata, MeshsyncDeploymentModeOperator,
		)

		assert.Equal(t, MeshsyncDeploymentModeEmbedded, effectiveMode(effective))
		assert.Equal(t, DeploymentModeLayerLayeredConfig, resolved.Layer)
	})

	t.Run("the merged document is left carrying only explicitly-set fields", func(t *testing.T) {
		// merged is what propagates to the cluster: stamping the resolved mode
		// onto it would turn an inherited mode into an explicit one and defeat
		// the withdrawal semantics of every other field.
		metadata := core.Map{MeshsyncDeploymentModeMetadataKey: string(MeshsyncDeploymentModeOperator)}

		merged, _, _ := ResolveConnectionControllersConfig(nil, nil, metadata, MeshsyncDeploymentModeOperator)

		assert.Equal(t, MeshsyncDeploymentModeUndefined, DeploymentModeFromControllersConfig(merged))
	})

	t.Run("everything else still resolves through the layer chain", func(t *testing.T) {
		version := "v0.8.1"
		override := &controllersconfig.MesheryControllersConfig{
			Operator: &controllersconfig.MesheryOperatorConfig{Version: &version},
		}

		_, effective, _ := ResolveConnectionControllersConfig(override, nil, core.Map{}, MeshsyncDeploymentModeEmbedded)

		require.NotNil(t, effective.Operator.Version)
		assert.Equal(t, version, *effective.Operator.Version)
		require.NotNil(t, effective.Meshsync)
		require.NotNil(t, effective.Meshsync.Replicas)
		assert.Equal(t, 1, *effective.Meshsync.Replicas, "built-in defaults still fill the rest of the document")
	})
}

// TestShouldRecordDeploymentModeOverride pins the rule that keeps a server-wide
// default reachable. The Connection Wizard's mode picker is pre-selected rather
// than empty, so registration receives a mode on every import and cannot tell a
// deliberate choice from an untouched default. Recording both as an override
// pinned every new connection, and a pinned connection ignores the server-wide
// default - which silently defeated
// TestB3ServerWideDeploymentModeReachesInheritingConnection for every connection
// created through the UI.
func TestShouldRecordDeploymentModeOverride(t *testing.T) {
	t.Run("a mode matching the inherited one is not an override", func(t *testing.T) {
		assert.False(t, ShouldRecordDeploymentModeOverride(
			MeshsyncDeploymentModeEmbedded, MeshsyncDeploymentModeEmbedded),
			"the wizard's pre-selected default must leave the connection inheriting")
	})

	t.Run("a mode diverging from the inherited one is an override", func(t *testing.T) {
		assert.True(t, ShouldRecordDeploymentModeOverride(
			MeshsyncDeploymentModeOperator, MeshsyncDeploymentModeEmbedded),
			"a deliberate divergence must be pinned")
	})

	t.Run("no requested mode is never an override", func(t *testing.T) {
		assert.False(t, ShouldRecordDeploymentModeOverride(
			MeshsyncDeploymentModeUndefined, MeshsyncDeploymentModeEmbedded))
	})

	// The end-to-end consequence: a connection registered with the inherited
	// mode must still follow a later change to the server-wide default.
	t.Run("a connection registered with the inherited mode follows a later default change", func(t *testing.T) {
		metadata := core.Map{}
		inherited := MeshsyncDeploymentModeEmbedded

		// Registration with the wizard's pre-selected default.
		if ShouldRecordDeploymentModeOverride(MeshsyncDeploymentModeEmbedded, inherited) {
			require.NoError(t, SetDeploymentModeOverride(metadata, MeshsyncDeploymentModeEmbedded))
		}
		MaterializeMeshsyncDeploymentMode(metadata, MeshsyncDeploymentModeEmbedded)

		override, err := ControllersConfigFromMetadata(metadata)
		require.NoError(t, err)
		assert.Equal(t, MeshsyncDeploymentModeUndefined, DeploymentModeFromControllersConfig(override),
			"registration must not pin an override for the inherited mode")

		// The admin later sets the server-wide default to operator.
		merged, _ := ResolveControllersConfig(override, modeDoc(MeshsyncDeploymentModeOperator))
		resolved := ResolveDeploymentMode(merged, metadata, MeshsyncDeploymentModeEmbedded)
		assert.Equal(t, MeshsyncDeploymentModeOperator, resolved.Mode,
			"the connection must follow the new server-wide default")
		assert.Equal(t, DeploymentModeLayerLayeredConfig, resolved.Layer)
	})
}
