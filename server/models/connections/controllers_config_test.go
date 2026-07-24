package connections

import (
	"testing"

	"github.com/meshery/schemas/models/core"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
)

func boolPtr(v bool) *bool    { return &v }
func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }
func modePtr(v controllersconfig.MesheryOperatorConfigDeploymentMode) *controllersconfig.MesheryOperatorConfigDeploymentMode {
	return &v
}

func TestResolveControllersConfigPrecedence(t *testing.T) {
	serverDefault := &controllersconfig.MesheryControllersConfig{
		Operator: &controllersconfig.MesheryOperatorConfig{
			DeploymentMode: modePtr(controllersconfig.Operator),
		},
		Meshsync: &controllersconfig.MeshSyncConfig{
			Replicas:      intPtr(3),
			RedactSecrets: boolPtr(true),
			Version:       strPtr("v1.0.2"),
		},
		Broker: &controllersconfig.MesheryBrokerConfig{
			Replicas: intPtr(3),
		},
	}
	override := &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			Replicas: intPtr(5),
		},
	}

	merged, effective := ResolveControllersConfig(override, serverDefault)

	// Override wins for the field it sets.
	if merged.Meshsync == nil || merged.Meshsync.Replicas == nil || *merged.Meshsync.Replicas != 5 {
		t.Fatalf("expected merged meshsync.replicas=5 (override), got %+v", merged.Meshsync)
	}
	// Server default fills fields the override leaves unset.
	if merged.Meshsync.RedactSecrets == nil || !*merged.Meshsync.RedactSecrets {
		t.Fatalf("expected merged meshsync.redactSecrets=true (server default), got %+v", merged.Meshsync.RedactSecrets)
	}
	if merged.Meshsync.Version == nil || *merged.Meshsync.Version != "v1.0.2" {
		t.Fatalf("expected merged meshsync.version=v1.0.2 (server default), got %+v", merged.Meshsync.Version)
	}
	if merged.Operator == nil || merged.Operator.DeploymentMode == nil || *merged.Operator.DeploymentMode != controllersconfig.Operator {
		t.Fatalf("expected merged operator.deploymentMode=operator (server default), got %+v", merged.Operator)
	}
	// Fields no layer set stay unset in merged (nothing to propagate).
	if merged.Meshsync.BrokerContentDedup != nil {
		t.Fatalf("expected merged meshsync.brokerContentDedup unset, got %v", *merged.Meshsync.BrokerContentDedup)
	}
	// The effective view fills built-in defaults for unset fields.
	if effective.Meshsync.BrokerContentDedup == nil || *effective.Meshsync.BrokerContentDedup {
		t.Fatalf("expected effective meshsync.brokerContentDedup=false (built-in), got %+v", effective.Meshsync.BrokerContentDedup)
	}
	if effective.Broker == nil || effective.Broker.Service == nil || effective.Broker.Service.Type == nil || *effective.Broker.Service.Type != controllersconfig.ClusterIP {
		t.Fatalf("expected effective broker.service.type=ClusterIP (built-in), got %+v", effective.Broker)
	}
	if effective.Broker.Replicas == nil || *effective.Broker.Replicas != 3 {
		t.Fatalf("expected effective broker.replicas=3 (server default), got %+v", effective.Broker.Replicas)
	}
	if effective.SchemaVersion != ControllersConfigSchemaVersion {
		t.Fatalf("expected schemaVersion stamped, got %q", effective.SchemaVersion)
	}
}

func TestResolveControllersConfigNilLayers(t *testing.T) {
	merged, effective := ResolveControllersConfig(nil, nil)
	if merged != nil && (merged.Operator != nil || merged.Meshsync != nil || merged.Broker != nil) {
		t.Fatalf("expected empty merged config for nil layers, got %+v", merged)
	}
	if effective == nil || effective.Meshsync == nil || effective.Meshsync.Replicas == nil || *effective.Meshsync.Replicas != 1 {
		t.Fatalf("expected effective built-in meshsync.replicas=1, got %+v", effective)
	}
	if effective.Operator == nil || effective.Operator.DeploymentMode == nil || string(*effective.Operator.DeploymentMode) != string(MeshsyncDeploymentModeDefault) {
		t.Fatalf("expected effective built-in deploymentMode=%s, got %+v", MeshsyncDeploymentModeDefault, effective.Operator)
	}
}

func TestMergeControllersConfigWatchListAtomic(t *testing.T) {
	base := &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			WatchList: &controllersconfig.MeshSyncWatchList{
				Blacklist: []string{"secrets.v1."},
			},
		},
	}
	override := &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			WatchList: &controllersconfig.MeshSyncWatchList{
				Whitelist: []controllersconfig.MeshSyncWatchedResource{
					{Resource: "pods.v1.", Events: []controllersconfig.MeshSyncWatchedResourceEvents{controllersconfig.ADDED}},
				},
			},
		},
	}

	merged := MergeControllersConfig(override, base)
	wl := merged.Meshsync.WatchList
	if wl == nil || len(wl.Whitelist) != 1 || len(wl.Blacklist) != 0 {
		t.Fatalf("expected watchList replaced atomically by override (whitelist only), got %+v", wl)
	}
}

func TestControllersConfigMetadataRoundTrip(t *testing.T) {
	metadata := core.Map{}
	cfg := &controllersconfig.MesheryControllersConfig{
		Meshsync: &controllersconfig.MeshSyncConfig{
			RedactSecrets:    boolPtr(true),
			OutputNamespaces: []string{"default", "production"},
		},
	}
	if err := SetControllersConfigToMetadata(metadata, cfg); err != nil {
		t.Fatalf("SetControllersConfigToMetadata: %v", err)
	}
	if _, ok := metadata[ControllersConfigMetadataKey]; !ok {
		t.Fatalf("expected metadata key %q to be set", ControllersConfigMetadataKey)
	}

	restored, err := ControllersConfigFromMetadata(metadata)
	if err != nil {
		t.Fatalf("ControllersConfigFromMetadata: %v", err)
	}
	if restored == nil || restored.Meshsync == nil || restored.Meshsync.RedactSecrets == nil || !*restored.Meshsync.RedactSecrets {
		t.Fatalf("expected redactSecrets=true after round trip, got %+v", restored)
	}
	if len(restored.Meshsync.OutputNamespaces) != 2 || restored.Meshsync.OutputNamespaces[0] != "default" {
		t.Fatalf("expected outputNamespaces preserved, got %+v", restored.Meshsync.OutputNamespaces)
	}
	if restored.SchemaVersion != ControllersConfigSchemaVersion {
		t.Fatalf("expected schemaVersion stamped on stored override, got %q", restored.SchemaVersion)
	}

	// An empty document removes the override.
	if err := SetControllersConfigToMetadata(metadata, &controllersconfig.MesheryControllersConfig{}); err != nil {
		t.Fatalf("SetControllersConfigToMetadata (clear): %v", err)
	}
	if _, ok := metadata[ControllersConfigMetadataKey]; ok {
		t.Fatalf("expected metadata key removed for empty document")
	}
	restored, err = ControllersConfigFromMetadata(metadata)
	if err != nil || restored != nil {
		t.Fatalf("expected nil config after clear, got %+v err=%v", restored, err)
	}
}

func TestControllersConfigFromMetadataTolerance(t *testing.T) {
	// Nil metadata inherits.
	cfg, err := ControllersConfigFromMetadata(nil)
	if err != nil || cfg != nil {
		t.Fatalf("expected nil for nil metadata, got %+v err=%v", cfg, err)
	}
	// JSON-string-shaped values are tolerated.
	metadata := core.Map{ControllersConfigMetadataKey: `{"meshsync":{"debugLogging":true}}`}
	cfg, err = ControllersConfigFromMetadata(metadata)
	if err != nil {
		t.Fatalf("expected string value tolerated, got err=%v", err)
	}
	if cfg == nil || cfg.Meshsync == nil || cfg.Meshsync.DebugLogging == nil || !*cfg.Meshsync.DebugLogging {
		t.Fatalf("expected debugLogging=true from string metadata, got %+v", cfg)
	}
	// Garbage is an error, not a panic.
	metadata = core.Map{ControllersConfigMetadataKey: "{not json"}
	if _, err := ControllersConfigFromMetadata(metadata); err == nil {
		t.Fatalf("expected error for malformed metadata")
	}
}

func TestValidateControllersConfig(t *testing.T) {
	lb := controllersconfig.LoadBalancer
	clusterIP := controllersconfig.ClusterIP
	cases := []struct {
		name    string
		cfg     *controllersconfig.MesheryControllersConfig
		wantErr bool
	}{
		{"nil config", nil, false},
		{"valid replicas", &controllersconfig.MesheryControllersConfig{Meshsync: &controllersconfig.MeshSyncConfig{Replicas: intPtr(10)}}, false},
		{"meshsync replicas too high", &controllersconfig.MesheryControllersConfig{Meshsync: &controllersconfig.MeshSyncConfig{Replicas: intPtr(11)}}, true},
		{"meshsync replicas too low", &controllersconfig.MesheryControllersConfig{Meshsync: &controllersconfig.MeshSyncConfig{Replicas: intPtr(0)}}, true},
		{"broker replicas too high", &controllersconfig.MesheryControllersConfig{Broker: &controllersconfig.MesheryBrokerConfig{Replicas: intPtr(11)}}, true},
		{
			"whitelist and blacklist together",
			&controllersconfig.MesheryControllersConfig{Meshsync: &controllersconfig.MeshSyncConfig{WatchList: &controllersconfig.MeshSyncWatchList{
				Whitelist: []controllersconfig.MeshSyncWatchedResource{{Resource: "pods.v1."}},
				Blacklist: []string{"secrets.v1."},
			}}},
			true,
		},
		{
			"whitelist entry missing resource",
			&controllersconfig.MesheryControllersConfig{Meshsync: &controllersconfig.MeshSyncConfig{WatchList: &controllersconfig.MeshSyncWatchList{
				Whitelist: []controllersconfig.MeshSyncWatchedResource{{Resource: ""}},
			}}},
			true,
		},
		{
			"loadBalancerClass without LoadBalancer type",
			&controllersconfig.MesheryControllersConfig{Broker: &controllersconfig.MesheryBrokerConfig{Service: &controllersconfig.MesheryBrokerServiceConfig{
				Type:              &clusterIP,
				LoadBalancerClass: strPtr("metallb"),
			}}},
			true,
		},
		{
			"loadBalancerSourceRanges with LoadBalancer type",
			&controllersconfig.MesheryControllersConfig{Broker: &controllersconfig.MesheryBrokerConfig{Service: &controllersconfig.MesheryBrokerServiceConfig{
				Type:                     &lb,
				LoadBalancerSourceRanges: []string{"10.0.0.0/8"},
			}}},
			false,
		},
		{
			"invalid deployment mode",
			&controllersconfig.MesheryControllersConfig{Operator: &controllersconfig.MesheryOperatorConfig{DeploymentMode: modePtr("sidecar")}},
			true,
		},
		{
			"valid deployment mode",
			&controllersconfig.MesheryControllersConfig{Operator: &controllersconfig.MesheryOperatorConfig{DeploymentMode: modePtr(controllersconfig.Embedded)}},
			false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateControllersConfig(tc.cfg)
			if tc.wantErr && err == nil {
				t.Fatalf("expected validation error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected validation error: %v", err)
			}
		})
	}
}

func TestDeploymentModeFromControllersConfig(t *testing.T) {
	if got := DeploymentModeFromControllersConfig(nil); got != MeshsyncDeploymentModeUndefined {
		t.Fatalf("expected undefined for nil config, got %s", got)
	}
	cfg := &controllersconfig.MesheryControllersConfig{Operator: &controllersconfig.MesheryOperatorConfig{DeploymentMode: modePtr(controllersconfig.Operator)}}
	if got := DeploymentModeFromControllersConfig(cfg); got != MeshsyncDeploymentModeOperator {
		t.Fatalf("expected operator mode, got %s", got)
	}
}
