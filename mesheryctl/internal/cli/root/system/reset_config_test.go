package system

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/viper"
)

// TestResetContextInConfigWritesActiveConfig verifies that system reset resets
// the context in the meshconfig given by --config and leaves the default
// meshconfig untouched.
func TestResetContextInConfigWritesActiveConfig(t *testing.T) {
	customPath := filepath.Join(t.TempDir(), "custom.yaml")
	custom := `contexts:
    prod:
        endpoint: http://meshery.example.com:9081
        token: default
        platform: docker
        channel: edge
        version: v0.8.0
current-context: prod
tokens:
    - location: auth.json
      name: default
`
	if err := os.WriteFile(customPath, []byte(custom), 0o600); err != nil {
		t.Fatal(err)
	}

	defaultPath := utils.CopyMeshconfigFixture(t, utils.SharedTestConfigPath(t))
	defaultBefore, err := os.ReadFile(defaultPath)
	if err != nil {
		t.Fatal(err)
	}

	prevCfgFile, prevDefaultPath := utils.CfgFile, utils.DefaultConfigPath
	utils.CfgFile, utils.DefaultConfigPath = customPath, defaultPath
	t.Cleanup(func() { utils.CfgFile, utils.DefaultConfigPath = prevCfgFile, prevDefaultPath })
	viper.Reset()
	t.Cleanup(viper.Reset)

	if err := resetContextInConfig("prod", "docker"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v := viper.New()
	v.SetConfigFile(customPath)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	mctlCfg, err := config.GetMesheryCtl(v)
	if err != nil {
		t.Fatal(err)
	}
	got, ok := mctlCfg.Contexts["prod"]
	if !ok {
		t.Fatalf("context prod missing from %s", customPath)
	}
	if got.Endpoint != utils.TemplateContext.Endpoint || got.Channel != utils.TemplateContext.Channel ||
		got.Version != utils.TemplateContext.Version || got.Platform != "docker" {
		t.Errorf("expected prod in %s reset to the default context on docker, got %+v", customPath, got)
	}

	defaultAfter, err := os.ReadFile(defaultPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(defaultBefore, defaultAfter) {
		t.Errorf("expected the default meshconfig to be unchanged, got:\n%s", defaultAfter)
	}
}
