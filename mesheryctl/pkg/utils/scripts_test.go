package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

func TestGKEConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		config    GKEConfig
		wantError bool
	}{
		{
			name: "valid config",
			config: GKEConfig{
				ConfigPath: "/tmp/config",
				SAName:     "test",
				Namespace:  "default",
			},
			wantError: false,
		},
		{
			name: "empty configPath",
			config: GKEConfig{
				ConfigPath: "",
				SAName:     "test",
				Namespace:  "default",
			},
			wantError: true,
		},
		{
			name: "empty SAName",
			config: GKEConfig{
				ConfigPath: "/tmp/config",
				SAName:     "",
				Namespace:  "default",
			},
			wantError: true,
		},
		{
			name: "empty namespace",
			config: GKEConfig{
				ConfigPath: "/tmp/config",
				SAName:     "test",
				Namespace:  "",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.wantError {
				t.Errorf("validate() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestInvalidConfigurations(t *testing.T) {
	invalidCfg := &GKEConfig{
		ConfigPath: "",
		SAName:     "",
		Namespace:  "",
	}

	if err := invalidCfg.validate(); err == nil {
		t.Error("Expected validation to fail for empty config")
	}
}

func TestGenerateConfigGKERequiresArgs(t *testing.T) {
	err := GenerateConfigGKE("", "", "")
	if err == nil {
		t.Fatal("expected error for empty arguments")
	}
}

func TestGenerateConfigGKEFailsWithoutCluster(t *testing.T) {
	t.Setenv("KUBECONFIG", filepath.Join(t.TempDir(), "does-not-exist"))

	err := GenerateConfigGKE(filepath.Join(t.TempDir(), "out.yaml"), "test-sa", "default")
	if err == nil {
		t.Fatal("expected error when no Kubernetes cluster is available")
	}
}

func TestClusterEndpointFromConfig(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		cfg := clientcmdapi.NewConfig()
		cfg.CurrentContext = "ctx-1"
		cfg.Contexts["ctx-1"] = &clientcmdapi.Context{Cluster: "cluster-a"}
		cfg.Clusters["cluster-a"] = &clientcmdapi.Cluster{Server: "https://example.invalid:6443"}

		name, endpoint, err := clusterEndpointFromConfig(*cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if name != "cluster-a" {
			t.Errorf("cluster name = %q, want cluster-a", name)
		}
		if endpoint != "https://example.invalid:6443" {
			t.Errorf("endpoint = %q, want https://example.invalid:6443", endpoint)
		}
	})

	t.Run("missing current context", func(t *testing.T) {
		cfg := clientcmdapi.NewConfig()
		if _, _, err := clusterEndpointFromConfig(*cfg); err == nil {
			t.Fatal("expected error when current context is empty")
		}
	})

	t.Run("missing cluster", func(t *testing.T) {
		cfg := clientcmdapi.NewConfig()
		cfg.CurrentContext = "ctx-1"
		cfg.Contexts["ctx-1"] = &clientcmdapi.Context{Cluster: "missing"}
		if _, _, err := clusterEndpointFromConfig(*cfg); err == nil {
			t.Fatal("expected error when cluster is missing")
		}
	})
}

func TestWriteKubeconfig(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "kubeconfig.yaml")
	cfg := &GKEConfig{
		ConfigPath: out,
		SAName:     "sa-test",
		Namespace:  "default",
	}

	ca := []byte("fake-ca-cert")
	if err := cfg.writeKubeconfig("my-cluster", "https://127.0.0.1:6443", "fake-token", ca); err != nil {
		t.Fatalf("writeKubeconfig: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read kubeconfig: %v", err)
	}
	content := string(data)
	for _, want := range []string{"my-cluster", "https://127.0.0.1:6443", "fake-token", "sa-test-default-my-cluster"} {
		if !strings.Contains(content, want) {
			t.Errorf("kubeconfig missing %q; content:\n%s", want, content)
		}
	}

	caPath := filepath.Join(dir, "ca.crt")
	gotCA, err := os.ReadFile(caPath)
	if err != nil {
		t.Fatalf("read ca.crt: %v", err)
	}
	if string(gotCA) != string(ca) {
		t.Errorf("ca.crt = %q, want %q", gotCA, ca)
	}
}

func TestTokenSecretName(t *testing.T) {
	cfg := &GKEConfig{SAName: "sa-meshery-abc"}
	if got := cfg.tokenSecretName(); got != "sa-meshery-abc-token" {
		t.Errorf("tokenSecretName() = %q, want sa-meshery-abc-token", got)
	}
}
