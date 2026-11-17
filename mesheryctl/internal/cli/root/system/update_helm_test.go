package system

import (
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	c "github.com/meshery/meshery/mesheryctl/pkg/constants"
	meshkitkube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/spf13/viper"
)

// TestUpgradeMesheryHelmChartsPassesCallbackURL verifies that system update on
// Kubernetes upgrades the Helm charts with the callback and provider URLs
// configured for Meshery, as system start installs them.
func TestUpgradeMesheryHelmChartsPassesCallbackURL(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set(c.CallbackURLENV, "https://meshery.example.com/api/user/token")
	viper.Set(c.ProviderURLsENV, "https://provider.example.com")

	type helmCall struct {
		version     string
		action      meshkitkube.HelmChartAction
		callbackURL string
		providerURL string
	}
	var calls []helmCall
	prev := applyHelmChartsFn
	applyHelmChartsFn = func(_ *meshkitkube.Client, _ *config.Context, version string, _ bool, act meshkitkube.HelmChartAction, callbackURL, providerURL string) error {
		calls = append(calls, helmCall{version, act, callbackURL, providerURL})
		return nil
	}
	t.Cleanup(func() { applyHelmChartsFn = prev })

	ctx := &config.Context{Platform: "kubernetes", Channel: "stable", Version: "latest"}
	if err := upgradeMesheryHelmCharts(nil, ctx, "latest"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := helmCall{"latest", meshkitkube.UPGRADE, "https://meshery.example.com/api/user/token", "https://provider.example.com"}
	if len(calls) != 1 || calls[0] != want {
		t.Errorf("expected one Helm call %+v, got %+v", want, calls)
	}
}
