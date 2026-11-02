package models

import (
	"testing"

	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/errors"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
)

// bootChartVersion stands in for the Meshery Server release stamped into the
// deployment config at boot. It is a published, at-or-above-floor version so
// that pinning it against publishedCharts is the identity: these tests are
// about how `operator.version` layers, not about chart-version pinning, which
// operator_chart_pinning_test.go covers.
const bootChartVersion = "v1.0.64"

func operatorVersionConfig(version string, mode connections.MeshsyncDeploymentMode) *controllersconfig.MesheryControllersConfig {
	cfg := &controllersconfig.MesheryControllersConfig{
		SchemaVersion: connections.ControllersConfigSchemaVersion,
		Operator:      &controllersconfig.MesheryOperatorConfig{},
	}
	if version != "" {
		cfg.Operator.Version = &version
	}
	if mode != connections.MeshsyncDeploymentModeUndefined {
		value := controllersconfig.MesheryOperatorConfigDeploymentMode(mode)
		cfg.Operator.DeploymentMode = &value
	}
	return cfg
}

func newTestControllersHelper(t *testing.T) *MesheryControllersHelper {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	mch := NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{
			MesheryReleaseVersion: bootChartVersion,
			HelmChartRepo:         "https://example.invalid/charts",
			GetHelmOverrides:      func(bool) map[string]interface{} { return map[string]interface{}{} },
		},
		nil,
		nil,
		nil,
		nil,
	)
	// Resolve against a fixed catalogue: chart-version pinning must never make
	// a unit test depend on the network or on what is published today.
	mch.chartVersions = func(string, string) ([]string, error) { return publishedCharts, nil }
	return mch
}

// TestB2OperatorVersionSelectsTheOperatorChartVersion covers defect B2: the
// form offered operator.version, the server stored and validated it, and
// nothing ever read it. The resolved value must now decide the Helm chart
// version the Meshery Operator is deployed at, and leaving it unset must
// reproduce the previous behavior exactly.
func TestB2OperatorVersionSelectsTheOperatorChartVersion(t *testing.T) {
	tests := []struct {
		name   string
		config *controllersconfig.MesheryControllersConfig
		want   string
	}{
		{
			name:   "no configuration at all keeps the boot-time chart version",
			config: nil,
			want:   bootChartVersion,
		},
		{
			name:   "no operator section keeps the boot-time chart version",
			config: &controllersconfig.MesheryControllersConfig{SchemaVersion: connections.ControllersConfigSchemaVersion},
			want:   bootChartVersion,
		},
		{
			name:   "operator section without a version keeps the boot-time chart version",
			config: operatorVersionConfig("", connections.MeshsyncDeploymentModeOperator),
			want:   bootChartVersion,
		},
		{
			name:   "a blank version is not an override",
			config: operatorVersionConfig("   ", connections.MeshsyncDeploymentModeOperator),
			want:   bootChartVersion,
		},
		{
			name:   "an explicit version becomes the chart version",
			config: operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator),
			want:   "v0.7.9",
		},
		{
			name:   "surrounding whitespace is trimmed",
			config: operatorVersionConfig("  v0.7.9  ", connections.MeshsyncDeploymentModeOperator),
			want:   "v0.7.9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mch := newTestControllersHelper(t)
			mch.SetControllersConfig(tt.config)

			got := mch.operatorDeploymentConfig().MesheryReleaseVersion
			if got != tt.want {
				t.Fatalf("operator chart version = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestB2OperatorDeploymentConfigKeepsTheServerWideFields asserts that threading
// operator.version replaces the chart version and nothing else: the chart
// repository and the Helm overrides (which carry the installed-adapter set)
// still come from the server-wide deployment config built at boot.
func TestB2OperatorDeploymentConfigKeepsTheServerWideFields(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator))

	depConfig := mch.operatorDeploymentConfig()
	if depConfig.HelmChartRepo != "https://example.invalid/charts" {
		t.Fatalf("chart repo = %q, want the server-wide value", depConfig.HelmChartRepo)
	}
	if depConfig.GetHelmOverrides == nil {
		t.Fatal("expected the server-wide Helm overrides function to be preserved")
	}
	if mch.oprDepConfig.MesheryReleaseVersion != bootChartVersion {
		t.Fatalf("server-wide deployment config was mutated: %q", mch.oprDepConfig.MesheryReleaseVersion)
	}
}

// TestB2ReconcileOperatorChartVersionNoOps pins the cases where a
// chart-version reconcile must not touch the cluster. Each of these returns
// before any Kubernetes client is built, so an unreachable context is safe.
func TestB2ReconcileOperatorChartVersionNoOps(t *testing.T) {
	k8sctx := K8sContext{ID: "ctx-1", Name: "test-cluster"}

	tests := []struct {
		name    string
		arrange func(mch *MesheryControllersHelper) *OperatorTracker
		want    string
	}{
		{
			name: "embedded mode: no operator is installed, so the field is inert",
			arrange: func(mch *MesheryControllersHelper) *OperatorTracker {
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeEmbedded)
				mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeEmbedded))
				return NewOperatorTracker(false)
			},
			want: "v0.7.9",
		},
		{
			name: "operator disabled server-wide",
			arrange: func(mch *MesheryControllersHelper) *OperatorTracker {
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
				mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator))
				return NewOperatorTracker(true)
			},
			want: "v0.7.9",
		},
		{
			name: "operator explicitly undeployed for this context is not resurrected",
			arrange: func(mch *MesheryControllersHelper) *OperatorTracker {
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
				mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator))
				tracker := NewOperatorTracker(false)
				tracker.Undeployed(k8sctx.ID, true)
				return tracker
			},
			want: "v0.7.9",
		},
		{
			name: "unchanged version: the attached handler already carries it",
			arrange: func(mch *MesheryControllersHelper) *OperatorTracker {
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
				mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator))
				mch.attachedOperatorChartVersion = "v0.7.9"
				return NewOperatorTracker(false)
			},
			want: "v0.7.9",
		},
		{
			name: "inheriting after attaching at the boot version is not a change",
			arrange: func(mch *MesheryControllersHelper) *OperatorTracker {
				mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
				mch.SetControllersConfig(operatorVersionConfig("", connections.MeshsyncDeploymentModeOperator))
				mch.attachedOperatorChartVersion = bootChartVersion
				return NewOperatorTracker(false)
			},
			want: bootChartVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mch := newTestControllersHelper(t)
			tracker := tt.arrange(mch)

			chartVersion, redeployed, err := mch.ReconcileOperatorChartVersion(k8sctx, tracker)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if redeployed {
				t.Fatal("expected no redeploy")
			}
			if chartVersion != tt.want {
				t.Fatalf("chart version = %q, want %q", chartVersion, tt.want)
			}
		})
	}
}

// TestB2ReconcileOperatorChartVersionRedeploysOnChange asserts the positive
// path reaches the cluster rather than silently no-opping: an operator-mode
// connection whose resolved chart version differs from the attached handler's
// re-attaches and re-runs the Helm release. The test context has no reachable
// cluster, so the attempt surfaces as a structured MeshKit error - which is
// itself the assertion that the guards let it through instead of swallowing
// the change.
func TestB2ReconcileOperatorChartVersionRedeploysOnChange(t *testing.T) {
	mch := newTestControllersHelper(t)
	mch.SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator)
	mch.SetControllersConfig(operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeOperator))
	mch.attachedOperatorChartVersion = bootChartVersion

	chartVersion, redeployed, err := mch.ReconcileOperatorChartVersion(
		K8sContext{ID: "ctx-1", Name: "test-cluster"},
		NewOperatorTracker(false),
	)
	if redeployed {
		t.Fatal("expected no successful redeploy against an unreachable cluster")
	}
	if chartVersion != "v0.7.9" {
		t.Fatalf("chart version = %q, want the configured version", chartVersion)
	}
	if err == nil {
		t.Fatal("expected the failed redeploy attempt to be reported, not swallowed")
	}
	switch code := errors.GetCode(err); code {
	case ErrReconcileOperatorChartVersionCode, ErrOperatorHandlerNotAttachedCode:
	default:
		t.Fatalf("expected a structured MeshKit error carrying this package's code, got %q from %v", code, err)
	}
}

// TestB2OperatorVersionResolvesThroughTheLayeredDocument asserts the field is
// layered like every other: a per-connection override wins, an unset override
// inherits the server-wide default, and neither layer setting it leaves the
// chart version tracking the Meshery Server release.
func TestB2OperatorVersionResolvesThroughTheLayeredDocument(t *testing.T) {
	tests := []struct {
		name          string
		override      *controllersconfig.MesheryControllersConfig
		serverDefault *controllersconfig.MesheryControllersConfig
		want          string
	}{
		{
			name:          "override wins over the server-wide default",
			override:      operatorVersionConfig("v0.7.9", connections.MeshsyncDeploymentModeUndefined),
			serverDefault: operatorVersionConfig("v0.8.1", connections.MeshsyncDeploymentModeUndefined),
			want:          "v0.7.9",
		},
		{
			name:          "no override inherits the server-wide default",
			override:      nil,
			serverDefault: operatorVersionConfig("v0.8.1", connections.MeshsyncDeploymentModeUndefined),
			want:          "v0.8.1",
		},
		{
			name:          "neither layer set: chart version tracks the server release",
			override:      nil,
			serverDefault: nil,
			want:          bootChartVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merged, _ := connections.ResolveControllersConfig(tt.override, tt.serverDefault)
			mch := newTestControllersHelper(t)
			mch.SetControllersConfig(merged)

			got := mch.operatorDeploymentConfig().MesheryReleaseVersion
			if got != tt.want {
				t.Fatalf("operator chart version = %q, want %q", got, tt.want)
			}
		})
	}
}
