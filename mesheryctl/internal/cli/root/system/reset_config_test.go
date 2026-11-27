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

// TestResetContextInConfigKeepsCurrentContext verifies that resetting a
// context picked with --context does not make it the current context.
func TestResetContextInConfigKeepsCurrentContext(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	meshconfig := `contexts:
    local:
        endpoint: http://localhost:9081
        token: default
        platform: docker
        channel: stable
        version: latest
    staging:
        endpoint: http://staging.example.com:9081
        token: default
        platform: kubernetes
        channel: edge
        version: v0.8.0
current-context: local
tokens:
    - location: auth.json
      name: default
`
	if err := os.WriteFile(configPath, []byte(meshconfig), 0o600); err != nil {
		t.Fatal(err)
	}

	prevCfgFile, prevDefaultPath := utils.CfgFile, utils.DefaultConfigPath
	utils.CfgFile, utils.DefaultConfigPath = configPath, configPath
	t.Cleanup(func() { utils.CfgFile, utils.DefaultConfigPath = prevCfgFile, prevDefaultPath })
	viper.Reset()
	t.Cleanup(viper.Reset)

	if err := resetContextInConfig("staging", "kubernetes"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	mctlCfg, err := config.GetMesheryCtl(v)
	if err != nil {
		t.Fatal(err)
	}
	if mctlCfg.CurrentContext != "local" {
		t.Errorf("expected current-context to stay local, got %s", mctlCfg.CurrentContext)
	}
	staging := mctlCfg.Contexts["staging"]
	if staging.Endpoint != utils.TemplateContext.Endpoint || staging.Version != utils.TemplateContext.Version ||
		staging.Platform != "kubernetes" {
		t.Errorf("expected staging reset to the default context on kubernetes, got %+v", staging)
	}
	if local := mctlCfg.Contexts["local"]; local.Endpoint != "http://localhost:9081" || local.Platform != "docker" {
		t.Errorf("expected local to be unchanged, got %+v", local)
	}
}

// TestResetMesheryConfigRefreshesConfigStruct exercises the actual production
// Cobra command execution (SystemCmd.Execute) and verifies via fetchManifestsFunc
// that manifest fetching receives the refreshed configuration with default channel,
// default version, default endpoint, preserved platform, and correct selected context.
func TestResetMesheryConfigRefreshesConfigStruct(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantContext  string
		wantPlatform string
	}{
		{
			name:         "normal reset on current context",
			args:         []string{"reset", "-y"},
			wantContext:  "local",
			wantPlatform: "docker",
		},
		{
			name:         "reset with temporary context flag",
			args:         []string{"reset", "-y", "-c", "staging"},
			wantContext:  "staging",
			wantPlatform: "kubernetes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configPath := filepath.Join(t.TempDir(), "config.yaml")
			meshconfig := `contexts:
    local:
        endpoint: http://custom-local:9081
        token: default
        platform: docker
        channel: edge
        version: v0.9.999-stale-pre-reset
    staging:
        endpoint: http://custom-staging:9081
        token: default
        platform: kubernetes
        channel: edge
        version: v0.9.999-stale-pre-reset
current-context: local
tokens:
    - location: auth.json
      name: default
`
			if err := os.WriteFile(configPath, []byte(meshconfig), 0o600); err != nil {
				t.Fatal(err)
			}

			prevCfgFile, prevDefaultPath := utils.CfgFile, utils.DefaultConfigPath
			utils.CfgFile, utils.DefaultConfigPath = configPath, configPath
			t.Cleanup(func() { utils.CfgFile, utils.DefaultConfigPath = prevCfgFile, prevDefaultPath })

			viper.Reset()
			t.Cleanup(viper.Reset)

			viper.SetConfigFile(configPath)
			if err := viper.ReadInConfig(); err != nil {
				t.Fatal(err)
			}

			// Intercept manifest fetching to observe the exact config passed to it
			var capturedCfg *config.MesheryCtlConfig
			oldFetch := fetchManifestsFunc
			fetchManifestsFunc = func(cfg *config.MesheryCtlConfig) error {
				capturedCfg = cfg
				return nil
			}
			t.Cleanup(func() {
				fetchManifestsFunc = oldFetch
				BreakupFunc()
			})

			// Execute actual production command path
			b := utils.SetupMeshkitLoggerTesting(t, false)
			SystemCmd.SetOut(b)
			SystemCmd.SetArgs(tt.args)
			err := SystemCmd.Execute()
			if err != nil {
				t.Fatalf("unexpected execution error: %v", err)
			}

			if capturedCfg == nil {
				t.Fatal("expected fetchManifestsFunc to be called with refreshed configuration, got nil")
			}

			if capturedCfg.GetCurrentContextName() != tt.wantContext {
				t.Errorf("expected selected context %q, got %q", tt.wantContext, capturedCfg.GetCurrentContextName())
			}

			capturedCtx, err := capturedCfg.GetCurrentContext()
			if err != nil {
				t.Fatalf("failed to retrieve current context from captured config: %v", err)
			}

			if capturedCtx.GetChannel() != utils.TemplateContext.Channel {
				t.Errorf("expected channel %q, got %q", utils.TemplateContext.Channel, capturedCtx.GetChannel())
			}
			if capturedCtx.GetVersion() != utils.TemplateContext.Version {
				t.Errorf("expected version %q, got %q", utils.TemplateContext.Version, capturedCtx.GetVersion())
			}
			if capturedCtx.GetEndpoint() != utils.TemplateContext.Endpoint {
				t.Errorf("expected endpoint %q, got %q", utils.TemplateContext.Endpoint, capturedCtx.GetEndpoint())
			}
			if capturedCtx.GetPlatform() != tt.wantPlatform {
				t.Errorf("expected platform %q, got %q", tt.wantPlatform, capturedCtx.GetPlatform())
			}

			// For the temporary context case, independently reload the config file
			// using a fresh Viper instance and verify persisted current-context remains unchanged
			vCheck := viper.New()
			vCheck.SetConfigFile(configPath)
			if err := vCheck.ReadInConfig(); err != nil {
				t.Fatal(err)
			}
			persistedCfg, err := config.GetMesheryCtl(vCheck)
			if err != nil {
				t.Fatal(err)
			}
			if persistedCfg.CurrentContext != "local" {
				t.Errorf("expected persisted current-context on disk to remain 'local', got %q", persistedCfg.CurrentContext)
			}
		})
	}
}
