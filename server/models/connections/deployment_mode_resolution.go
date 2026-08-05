// This file owns the single definition of "which MeshSync deployment mode is
// this Kubernetes connection running, and why".
//
// Historically `connection.metadata.meshsync_deployment_mode` carried two
// different facts under one key:
//
//  1. the user's explicit per-connection choice (written by the Connection
//     Wizard's MeshSync Mode step and by the kubeconfig import picker), and
//  2. a materialized cache of the resolved mode (rewritten by the controllers
//     editor on every save) that the pre-layered consumers - the connection
//     state machine, the header status chips, the kubeconfig flows - read.
//
// Because every Kubernetes connection has the key written at registration,
// treating it as fact (1) made it outrank every other layer, so a server-wide
// default could never reach an existing connection. Separating the two facts is
// the fix:
//
//   - `controllers_config.operator.deploymentMode` is the ONLY store of the
//     explicit per-connection choice. Every entry point writes it through
//     SetDeploymentModeOverride.
//   - `meshsync_deployment_mode` is ONLY the materialized cache, written through
//     MaterializeMeshsyncDeploymentMode and read by legacy consumers that want
//     "what is this connection running right now".
//
// The cache still participates in resolution, but strictly *below* the layered
// document, as a compatibility floor for connections registered before the
// layered document existed. It can therefore no longer mask a server-wide
// default.
package connections

import (
	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
)

// DeploymentModeLayer names the layer that supplied a resolved deployment mode.
// Reported so callers (and, through them, the user) can say where an effective
// mode came from instead of presenting it as an unexplained fact.
type DeploymentModeLayer string

const (
	// DeploymentModeLayerLayeredConfig is the layered controllers-configuration
	// document: the per-connection override merged over the server-wide
	// default. This is where every mode set from the UI lands.
	DeploymentModeLayerLayeredConfig DeploymentModeLayer = "layeredConfig"
	// DeploymentModeLayerLegacyMetadata is the connection's materialized
	// meshsync_deployment_mode entry, consulted only when no layer of the
	// controllers-configuration document sets a mode. It preserves the choice
	// of connections registered before the layered document existed.
	DeploymentModeLayerLegacyMetadata DeploymentModeLayer = "legacyConnectionMetadata"
	// DeploymentModeLayerServerEnvDefault is MESHSYNC_DEFAULT_DEPLOYMENT_MODE.
	DeploymentModeLayerServerEnvDefault DeploymentModeLayer = "serverEnvDefault"
	// DeploymentModeLayerBuiltIn is the compiled-in default.
	DeploymentModeLayerBuiltIn DeploymentModeLayer = "builtIn"
)

// ResolvedDeploymentMode is a deployment mode together with the layer that
// supplied it.
type ResolvedDeploymentMode struct {
	Mode  MeshsyncDeploymentMode
	Layer DeploymentModeLayer
}

// ResolveDeploymentMode decides the deployment mode a connection runs in.
//
// merged is the already-layered controllers configuration for the connection
// (per-connection override merged over the server-wide default, i.e. the
// `merged` document from ResolveControllersConfig). It outranks everything
// else, so an explicit per-connection override still wins and an inherited
// value genuinely follows the server-wide default.
//
// metadata's meshsync_deployment_mode entry is the compatibility floor
// described in this file's package comment: it applies only when no layer of
// the controllers-configuration document sets a mode.
//
// serverEnvDefault is MESHSYNC_DEFAULT_DEPLOYMENT_MODE; when it too is
// undefined the compiled-in default applies. The result is never
// MeshsyncDeploymentModeUndefined.
func ResolveDeploymentMode(
	merged *controllersconfig.MesheryControllersConfig,
	metadata core.Map,
	serverEnvDefault MeshsyncDeploymentMode,
) ResolvedDeploymentMode {
	if mode := DeploymentModeFromControllersConfig(merged); mode != MeshsyncDeploymentModeUndefined {
		return ResolvedDeploymentMode{Mode: mode, Layer: DeploymentModeLayerLayeredConfig}
	}
	if mode := MeshsyncDeploymentModeFromMetadata(metadata); mode != MeshsyncDeploymentModeUndefined {
		return ResolvedDeploymentMode{Mode: mode, Layer: DeploymentModeLayerLegacyMetadata}
	}
	if serverEnvDefault != MeshsyncDeploymentModeUndefined {
		return ResolvedDeploymentMode{Mode: serverEnvDefault, Layer: DeploymentModeLayerServerEnvDefault}
	}
	return ResolvedDeploymentMode{Mode: MeshsyncDeploymentModeDefault, Layer: DeploymentModeLayerBuiltIn}
}

// SetDeploymentModeOverride records an explicit per-connection deployment-mode
// choice in the one place that stores it: the connection's layered
// controllers-configuration override. Every entry point that lets a user pick a
// mode (the Connection Wizard's MeshSync Mode step, the kubeconfig import
// picker, the controllers editor) must go through here, so the editor and the
// wizard can never disagree about what the connection is set to.
//
// Passing MeshsyncDeploymentModeUndefined clears the override, returning the
// connection to inheriting the server-wide default. A nil metadata map is a
// no-op, mirroring the other metadata writers in this package.
//
// Writing the override does not refresh the materialized cache; callers that
// persist the connection must also call MaterializeMeshsyncDeploymentMode with
// the newly resolved mode so legacy readers stay in step.
func SetDeploymentModeOverride(metadata core.Map, mode MeshsyncDeploymentMode) error {
	if metadata == nil {
		return nil
	}
	override, err := ControllersConfigFromMetadata(metadata)
	if err != nil {
		return err
	}
	if override == nil {
		override = &controllersconfig.MesheryControllersConfig{}
	}

	if mode == MeshsyncDeploymentModeUndefined {
		if override.Operator == nil {
			return nil
		}
		override.Operator.DeploymentMode = nil
		if override.Operator.Version == nil {
			// An operator section that sets nothing is not an override.
			override.Operator = nil
		}
	} else {
		if override.Operator == nil {
			override.Operator = &controllersconfig.MesheryOperatorConfig{}
		}
		value := controllersconfig.MesheryOperatorConfigDeploymentMode(mode)
		override.Operator.DeploymentMode = &value
	}

	return SetControllersConfigToMetadata(metadata, override)
}
