// This file hosts the MesheryControllersConfig helpers: how the server layers
// the Meshery Operator / MeshSync / Broker configuration for a Kubernetes
// connection. The wire contract lives in
// github.com/meshery/schemas/models/v1alpha1/controllers_config; the layering
// (per-connection override -> server-wide default -> built-in default) is
// Meshery domain logic and therefore lives here, next to the
// MeshsyncDeploymentMode helpers it complements.
package connections

import (
	"encoding/json"

	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
)

// ControllersConfigMetadataKey is the key under which a connection's
// controllers-configuration override is stored on the connection's metadata
// map. Kept as snake_case for consistency with the sibling
// meshsync_deployment_mode entry persisted on the same map.
const ControllersConfigMetadataKey = "controllers_config"

// ControllersConfigSchemaVersion is the schema version stamped onto every
// controllers-configuration document the server stores or returns.
const ControllersConfigSchemaVersion = "controllers.meshery.io/v1alpha1"

// ControllersConfigFromMetadata extracts the per-connection controllers
// configuration override stored on a connection's metadata map. A missing key
// or nil metadata yields (nil, nil): the connection inherits entirely. The
// stored value is a JSON-shaped map (or a JSON string, for tolerance) and is
// round-tripped through encoding/json into the typed schema model.
func ControllersConfigFromMetadata(metadata core.Map) (*controllersconfig.MesheryControllersConfig, error) {
	if metadata == nil {
		return nil, nil
	}
	raw, exists := metadata[ControllersConfigMetadataKey]
	if !exists || raw == nil {
		return nil, nil
	}

	var encoded []byte
	var err error
	switch v := raw.(type) {
	case string:
		encoded = []byte(v)
	default:
		encoded, err = json.Marshal(v)
		if err != nil {
			return nil, ErrControllersConfigMetadata(err)
		}
	}

	cfg := &controllersconfig.MesheryControllersConfig{}
	if err := json.Unmarshal(encoded, cfg); err != nil {
		return nil, ErrControllersConfigMetadata(err)
	}
	return cfg, nil
}

// SetControllersConfigToMetadata writes (or overwrites) the controllers
// configuration override on a connection's metadata map, stored as a
// JSON-shaped map so it serializes identically through every provider path.
// Passing a nil or empty config removes the override entirely. A nil metadata
// map is a no-op (writing to a nil map would panic); callers that need to set
// an override on a fresh connection must initialise the map first.
func SetControllersConfigToMetadata(metadata core.Map, cfg *controllersconfig.MesheryControllersConfig) error {
	if metadata == nil {
		return nil
	}
	if cfg == nil || (cfg.Operator == nil && cfg.Meshsync == nil && cfg.Broker == nil) {
		delete(metadata, ControllersConfigMetadataKey)
		return nil
	}
	cfg.SchemaVersion = ControllersConfigSchemaVersion

	encoded, err := json.Marshal(cfg)
	if err != nil {
		return ErrControllersConfigMetadata(err)
	}
	asMap := map[string]interface{}{}
	if err := json.Unmarshal(encoded, &asMap); err != nil {
		return ErrControllersConfigMetadata(err)
	}
	metadata[ControllersConfigMetadataKey] = asMap
	return nil
}

// MergeControllersConfig overlays the override onto the base, producing a
// document that contains only explicitly-set fields from either layer, with
// the override winning per field. Scalars merge per leaf. Collections
// (watchList, outputNamespaces, outputResources, service annotations and
// source ranges) merge atomically: a layer that sets one replaces the lower
// layer's value entirely, because element-wise merging of a whitelist or a
// namespace filter would produce a scope neither layer asked for.
func MergeControllersConfig(override, base *controllersconfig.MesheryControllersConfig) *controllersconfig.MesheryControllersConfig {
	if override == nil && base == nil {
		return nil
	}
	merged := &controllersconfig.MesheryControllersConfig{SchemaVersion: ControllersConfigSchemaVersion}
	merged.Operator = mergeOperatorConfig(sectionOperator(override), sectionOperator(base))
	merged.Meshsync = mergeMeshSyncConfig(sectionMeshsync(override), sectionMeshsync(base))
	merged.Broker = mergeBrokerConfig(sectionBroker(override), sectionBroker(base))
	return merged
}

func sectionOperator(c *controllersconfig.MesheryControllersConfig) *controllersconfig.MesheryOperatorConfig {
	if c == nil {
		return nil
	}
	return c.Operator
}

func sectionMeshsync(c *controllersconfig.MesheryControllersConfig) *controllersconfig.MeshSyncConfig {
	if c == nil {
		return nil
	}
	return c.Meshsync
}

func sectionBroker(c *controllersconfig.MesheryControllersConfig) *controllersconfig.MesheryBrokerConfig {
	if c == nil {
		return nil
	}
	return c.Broker
}

func mergeOperatorConfig(override, base *controllersconfig.MesheryOperatorConfig) *controllersconfig.MesheryOperatorConfig {
	if override == nil && base == nil {
		return nil
	}
	merged := &controllersconfig.MesheryOperatorConfig{}
	if base != nil {
		merged.DeploymentMode = base.DeploymentMode
		merged.Version = base.Version
	}
	if override != nil {
		if override.DeploymentMode != nil {
			merged.DeploymentMode = override.DeploymentMode
		}
		if override.Version != nil {
			merged.Version = override.Version
		}
	}
	return merged
}

func mergeMeshSyncConfig(override, base *controllersconfig.MeshSyncConfig) *controllersconfig.MeshSyncConfig {
	if override == nil && base == nil {
		return nil
	}
	merged := &controllersconfig.MeshSyncConfig{}
	if base != nil {
		*merged = *base
	}
	if override != nil {
		if override.Version != nil {
			merged.Version = override.Version
		}
		if override.Replicas != nil {
			merged.Replicas = override.Replicas
		}
		if override.WatchList != nil {
			merged.WatchList = override.WatchList
		}
		if override.OutputNamespaces != nil {
			merged.OutputNamespaces = override.OutputNamespaces
		}
		if override.OutputResources != nil {
			merged.OutputResources = override.OutputResources
		}
		if override.RedactSecrets != nil {
			merged.RedactSecrets = override.RedactSecrets
		}
		if override.BrokerContentDedup != nil {
			merged.BrokerContentDedup = override.BrokerContentDedup
		}
		if override.DebugLogging != nil {
			merged.DebugLogging = override.DebugLogging
		}
	}
	return merged
}

func mergeBrokerConfig(override, base *controllersconfig.MesheryBrokerConfig) *controllersconfig.MesheryBrokerConfig {
	if override == nil && base == nil {
		return nil
	}
	merged := &controllersconfig.MesheryBrokerConfig{}
	if base != nil {
		merged.Version = base.Version
		merged.Replicas = base.Replicas
		merged.Service = base.Service
	}
	if override != nil {
		if override.Version != nil {
			merged.Version = override.Version
		}
		if override.Replicas != nil {
			merged.Replicas = override.Replicas
		}
		if override.Service != nil {
			merged.Service = mergeBrokerServiceConfig(override.Service, merged.Service)
		}
	}
	return merged
}

func mergeBrokerServiceConfig(override, base *controllersconfig.MesheryBrokerServiceConfig) *controllersconfig.MesheryBrokerServiceConfig {
	if override == nil && base == nil {
		return nil
	}
	merged := &controllersconfig.MesheryBrokerServiceConfig{}
	if base != nil {
		*merged = *base
	}
	if override != nil {
		if override.Type != nil {
			merged.Type = override.Type
		}
		if override.Annotations != nil {
			merged.Annotations = override.Annotations
		}
		if override.LoadBalancerClass != nil {
			merged.LoadBalancerClass = override.LoadBalancerClass
		}
		if override.LoadBalancerSourceRanges != nil {
			merged.LoadBalancerSourceRanges = override.LoadBalancerSourceRanges
		}
		if override.ExternalEndpointOverride != nil {
			merged.ExternalEndpointOverride = override.ExternalEndpointOverride
		}
	}
	return merged
}

// BuiltInControllersConfig returns the built-in defaults: the values that
// apply when neither a per-connection override nor a server-wide default sets
// a field. These mirror the defaults baked into MeshSync, Meshery Operator,
// and Meshery Broker themselves.
func BuiltInControllersConfig() *controllersconfig.MesheryControllersConfig {
	defaultMode := controllersconfig.MesheryOperatorConfigDeploymentMode(MeshsyncDeploymentModeDefault)
	one := 1
	off := false
	clusterIP := controllersconfig.ClusterIP
	return &controllersconfig.MesheryControllersConfig{
		SchemaVersion: ControllersConfigSchemaVersion,
		Operator: &controllersconfig.MesheryOperatorConfig{
			DeploymentMode: &defaultMode,
		},
		Meshsync: &controllersconfig.MeshSyncConfig{
			Replicas:           &one,
			RedactSecrets:      &off,
			BrokerContentDedup: &off,
			DebugLogging:       &off,
		},
		Broker: &controllersconfig.MesheryBrokerConfig{
			Replicas: &one,
			Service: &controllersconfig.MesheryBrokerServiceConfig{
				Type: &clusterIP,
			},
		},
	}
}

// ResolveControllersConfig applies the precedence chain: per-connection
// override -> server-wide default -> built-in default. It returns two
// documents:
//
//   - merged: only the fields explicitly set at the override or default
//     layer, override winning. This is what the server propagates to a
//     cluster; unset fields are left untouched so the controllers' own
//     defaults apply (and previously-propagated values are withdrawn).
//   - effective: merged overlaid onto the built-in defaults. This is the
//     complete resolved view surfaced to clients.
func ResolveControllersConfig(override, serverDefault *controllersconfig.MesheryControllersConfig) (merged, effective *controllersconfig.MesheryControllersConfig) {
	merged = MergeControllersConfig(override, serverDefault)
	effective = MergeControllersConfig(merged, BuiltInControllersConfig())
	if effective != nil {
		effective.SchemaVersion = ControllersConfigSchemaVersion
	}
	return merged, effective
}

// DeploymentModeFromControllersConfig extracts the deployment mode carried by
// a controllers configuration document, collapsing to
// MeshsyncDeploymentModeUndefined when the document does not set one.
func DeploymentModeFromControllersConfig(cfg *controllersconfig.MesheryControllersConfig) MeshsyncDeploymentMode {
	if cfg == nil || cfg.Operator == nil || cfg.Operator.DeploymentMode == nil {
		return MeshsyncDeploymentModeUndefined
	}
	return MeshsyncDeploymentModeFromString(string(*cfg.Operator.DeploymentMode))
}

// ValidateControllersConfig enforces the guardrails that the schema cannot
// fully express plus the ones it does (so the server rejects invalid
// documents regardless of transport): replica ranges, watch-list mutual
// exclusion, and broker service coherence.
func ValidateControllersConfig(cfg *controllersconfig.MesheryControllersConfig) error {
	if cfg == nil {
		return nil
	}
	if ms := cfg.Meshsync; ms != nil {
		if ms.Replicas != nil && (*ms.Replicas < 1 || *ms.Replicas > 10) {
			return ErrControllersConfigInvalid("meshsync.replicas must be between 1 and 10")
		}
		if wl := ms.WatchList; wl != nil {
			if len(wl.Whitelist) > 0 && len(wl.Blacklist) > 0 {
				return ErrControllersConfigInvalid("meshsync.watchList: whitelist and blacklist are mutually exclusive; set at most one")
			}
			for _, entry := range wl.Whitelist {
				if entry.Resource == "" {
					return ErrControllersConfigInvalid("meshsync.watchList.whitelist entries must set resource")
				}
			}
		}
	}
	if br := cfg.Broker; br != nil {
		if br.Replicas != nil && (*br.Replicas < 1 || *br.Replicas > 10) {
			return ErrControllersConfigInvalid("broker.replicas must be between 1 and 10")
		}
		if svc := br.Service; svc != nil {
			if svc.Type != nil {
				switch *svc.Type {
				case controllersconfig.ClusterIP, controllersconfig.NodePort, controllersconfig.LoadBalancer:
				default:
					return ErrControllersConfigInvalid("broker.service.type must be one of ClusterIP, NodePort, LoadBalancer")
				}
			}
			isLoadBalancer := svc.Type != nil && *svc.Type == controllersconfig.LoadBalancer
			if !isLoadBalancer && svc.LoadBalancerClass != nil {
				return ErrControllersConfigInvalid("broker.service.loadBalancerClass is only valid when broker.service.type is LoadBalancer")
			}
			if !isLoadBalancer && len(svc.LoadBalancerSourceRanges) > 0 {
				return ErrControllersConfigInvalid("broker.service.loadBalancerSourceRanges is only valid when broker.service.type is LoadBalancer")
			}
		}
	}
	if op := cfg.Operator; op != nil && op.DeploymentMode != nil {
		mode := MeshsyncDeploymentModeFromString(string(*op.DeploymentMode))
		if mode == MeshsyncDeploymentModeUndefined {
			return ErrControllersConfigInvalid("operator.deploymentMode must be either \"operator\" or \"embedded\"")
		}
	}
	return nil
}
