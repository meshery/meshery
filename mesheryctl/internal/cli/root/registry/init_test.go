// # Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const testArtifactHubSource = "https://artifacthub.io/packages/helm/cert-manager/cert-manager"

// validInitFlags returns flags for a scaffold that passes validation.
func validInitFlags(path string) cmdRegistryInitFlags {
	return cmdRegistryInitFlags{
		Source:      testArtifactHubSource,
		Category:    "Security & Compliance",
		SubCategory: "Key Management",
		Description: "Cloud native certificate management for Kubernetes.",
		Path:        path,
	}
}

// setRegistryInitFlags sets the package-level flags RunE reads and restores
// them after the test, so cases do not leak state into one another.
func setRegistryInitFlags(t *testing.T, flags cmdRegistryInitFlags) {
	t.Helper()
	utils.SetupMeshkitLoggerTesting(t, false)
	orig := registryInitFlags
	registryInitFlags = flags
	t.Cleanup(func() { registryInitFlags = orig })
}

func readModelSource(t *testing.T, path string) modelSource {
	t.Helper()
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	var ms modelSource
	require.NoError(t, yaml.Unmarshal(data, &ms))
	return ms
}

func TestValidateModelSourceName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "lowercase with hyphens", input: "cert-manager"},
		{name: "digits", input: "k8s-operator2"},
		{name: "empty", input: "", wantErr: true},
		{name: "uppercase", input: "Cert-Manager", wantErr: true},
		{name: "underscore", input: "cert_manager", wantErr: true},
		{name: "path separator", input: "../cert-manager", wantErr: true},
		{name: "too long", input: strings.Repeat("a", maxNameLength+1), wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateModelSourceName(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDetectSourceKind(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		want    string
		wantErr bool
	}{
		{name: "github https", source: "https://github.com/cert-manager/cert-manager", want: sourceKindGithub},
		{name: "github www", source: "https://www.github.com/cert-manager/cert-manager", want: sourceKindGithub},
		{name: "git scheme", source: "git://github.com/cert-manager/cert-manager/master/deploy/crds", want: sourceKindGithub},
		{name: "artifact hub", source: testArtifactHubSource, want: sourceKindArtifactHub},
		{name: "http manifest", source: "https://example.com/manifests/crds.yaml", want: sourceKindHTTP},
		{name: "unsupported scheme", source: "ftp://example.com/crds.yaml", wantErr: true},
		{name: "not a url", source: "cert-manager", wantErr: true},
		{name: "empty", source: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := detectSourceKind(tt.source)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBuildModelSource(t *testing.T) {
	tests := []struct {
		name           string
		modify         func(f *cmdRegistryInitFlags)
		wantKind       string
		wantRegistrant string
		wantDisplay    string
		wantErr        string
	}{
		{
			name:           "derives kind, registrant and display name",
			wantKind:       sourceKindArtifactHub,
			wantRegistrant: "artifacthub",
			wantDisplay:    "Cert Manager",
		},
		{
			name: "explicit flags override the derived values",
			modify: func(f *cmdRegistryInitFlags) {
				f.SourceKind = "HTTP"
				f.Registrant = "meshery"
				f.DisplayName = "cert-manager"
			},
			wantKind:       sourceKindHTTP,
			wantRegistrant: "meshery",
			wantDisplay:    "cert-manager",
		},
		{
			name:           "http source defaults the registrant to meshery",
			modify:         func(f *cmdRegistryInitFlags) { f.Source = "https://example.com/crds.yaml" },
			wantKind:       sourceKindHTTP,
			wantRegistrant: "meshery",
			wantDisplay:    "Cert Manager",
		},
		{
			name:    "invalid source kind",
			modify:  func(f *cmdRegistryInitFlags) { f.SourceKind = "gitlab" },
			wantErr: "--source-kind",
		},
		{
			name:    "missing description",
			modify:  func(f *cmdRegistryInitFlags) { f.Description = "  " },
			wantErr: "--description",
		},
		{
			name:    "missing category",
			modify:  func(f *cmdRegistryInitFlags) { f.Category = "" },
			wantErr: "--category",
		},
		{
			name:    "description over the schema limit",
			modify:  func(f *cmdRegistryInitFlags) { f.Description = strings.Repeat("d", maxDescriptionLength+1) },
			wantErr: "--description",
		},
		{
			name:           "multibyte description at the schema limit",
			modify:         func(f *cmdRegistryInitFlags) { f.Description = strings.Repeat("証", maxDescriptionLength) },
			wantKind:       sourceKindArtifactHub,
			wantRegistrant: "artifacthub",
			wantDisplay:    "Cert Manager",
		},
		{
			name:    "multibyte description over the schema limit",
			modify:  func(f *cmdRegistryInitFlags) { f.Description = strings.Repeat("証", maxDescriptionLength+1) },
			wantErr: "--description",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := validInitFlags(t.TempDir())
			if tt.modify != nil {
				tt.modify(&flags)
			}
			ms, err := buildModelSource("cert-manager", flags)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, modelSourceSchemaVersion, ms.SchemaVersion)
			assert.Equal(t, "cert-manager", ms.Name)
			assert.Equal(t, tt.wantKind, ms.Source.Kind)
			assert.Equal(t, tt.wantRegistrant, ms.Registrant)
			assert.Equal(t, tt.wantDisplay, ms.DisplayName)
		})
	}
}

func TestRegistryInitCreatesScaffold(t *testing.T) {
	modelsDir := t.TempDir()
	setRegistryInitFlags(t, validInitFlags(modelsDir))

	require.NoError(t, initCmd.RunE(initCmd, []string{"cert-manager"}))

	modelDir := filepath.Join(modelsDir, "cert-manager")
	ms := readModelSource(t, filepath.Join(modelDir, modelSourceFileName))
	assert.Equal(t, modelSource{
		SchemaVersion: modelSourceSchemaVersion,
		Name:          "cert-manager",
		DisplayName:   "Cert Manager",
		Registrant:    "artifacthub",
		Category:      "Security & Compliance",
		SubCategory:   "Key Management",
		Description:   "Cloud native certificate management for Kubernetes.",
		Source:        modelSourceLocation{URL: testArtifactHubSource, Kind: sourceKindArtifactHub},
		Metadata: modelSourceMetadata{
			SvgColor:    "./svg/color.svg",
			SvgWhite:    "./svg/white.svg",
			SvgComplete: "./svg/complete.svg",
		},
	}, ms)

	for _, name := range svgFiles {
		data, err := os.ReadFile(filepath.Join(modelDir, "svg", name))
		require.NoError(t, err, "missing placeholder %s", name)
		assert.Equal(t, placeholderSVG, string(data))
	}
	_, err := os.Stat(filepath.Join(modelDir, "overrides", "components", ".gitkeep"))
	assert.NoError(t, err)

	// website is optional and left out of a new scaffold, and nested fields
	// use the same 2-space indentation as the schema templates.
	raw, err := os.ReadFile(filepath.Join(modelDir, modelSourceFileName))
	require.NoError(t, err)
	assert.NotContains(t, string(raw), "website")
	assert.Contains(t, string(raw), "\nsource:\n  url: "+testArtifactHubSource+"\n")
}

func TestRegistryInitExistingManifest(t *testing.T) {
	tests := []struct {
		name          string
		force         bool
		wantErr       bool
		wantOverwrite bool
	}{
		{name: "refuses to overwrite without --force", wantErr: true},
		{name: "overwrites with --force", force: true, wantOverwrite: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modelsDir := t.TempDir()
			modelDir := filepath.Join(modelsDir, "cert-manager")
			require.NoError(t, os.MkdirAll(filepath.Join(modelDir, "svg"), 0o755))
			manifest := filepath.Join(modelDir, modelSourceFileName)
			require.NoError(t, os.WriteFile(manifest, []byte("existing: true\n"), 0o644))
			customLogo := filepath.Join(modelDir, "svg", "color.svg")
			require.NoError(t, os.WriteFile(customLogo, []byte("<svg>custom</svg>"), 0o644))

			flags := validInitFlags(modelsDir)
			flags.Force = tt.force
			setRegistryInitFlags(t, flags)

			err := initCmd.RunE(initCmd, []string{"cert-manager"})
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "--force")
			} else {
				require.NoError(t, err)
			}

			data, err := os.ReadFile(manifest)
			require.NoError(t, err)
			if tt.wantOverwrite {
				assert.Contains(t, string(data), "schemaVersion: "+modelSourceSchemaVersion)
			} else {
				assert.Equal(t, "existing: true\n", string(data))
			}

			// An existing logo is never replaced, with or without --force.
			logo, err := os.ReadFile(customLogo)
			require.NoError(t, err)
			assert.Equal(t, "<svg>custom</svg>", string(logo))
		})
	}
}

func TestRegistryInitExistingModelWithoutManifest(t *testing.T) {
	// A model that already lives in the repo (for example one merged by PR)
	// gets a manifest added next to its generated version folders.
	modelsDir := t.TempDir()
	modelDir := filepath.Join(modelsDir, "cert-manager")
	versionFile := filepath.Join(modelDir, "v1.13.0", "v1.0.0", "model.json")
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFile), 0o755))
	require.NoError(t, os.WriteFile(versionFile, []byte(`{"name":"cert-manager"}`), 0o644))
	setRegistryInitFlags(t, validInitFlags(modelsDir))

	require.NoError(t, initCmd.RunE(initCmd, []string{"cert-manager"}))

	assert.Equal(t, "cert-manager", readModelSource(t, filepath.Join(modelDir, modelSourceFileName)).Name)
	data, err := os.ReadFile(versionFile)
	require.NoError(t, err)
	assert.Equal(t, `{"name":"cert-manager"}`, string(data), "existing version folders must be left untouched")
}

func TestRegistryInitModelPathIsAFile(t *testing.T) {
	modelsDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(modelsDir, "cert-manager"), []byte("not a dir"), 0o644))
	setRegistryInitFlags(t, validInitFlags(modelsDir))

	err := initCmd.RunE(initCmd, []string{"cert-manager"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not a directory")
}

func TestRegistryInitArgs(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{name: "one valid name", args: []string{"cert-manager"}},
		{name: "no name", args: []string{}, wantErr: true},
		{name: "two names", args: []string{"a", "b"}, wantErr: true},
		{name: "invalid name", args: []string{"Cert_Manager"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := initCmd.Args(initCmd, tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRegistryInitRequiredFlags(t *testing.T) {
	values := map[string]string{
		"source":       testArtifactHubSource,
		"category":     "Security & Compliance",
		"sub-category": "Key Management",
		"description":  "Cloud native certificate management for Kubernetes.",
	}
	for _, missing := range []string{"source", "category", "sub-category", "description"} {
		t.Run(missing, func(t *testing.T) {
			setRegistryInitFlags(t, validInitFlags(t.TempDir()))
			t.Cleanup(func() {
				initCmd.Flags().VisitAll(func(f *pflag.Flag) {
					_ = f.Value.Set(f.DefValue)
					f.Changed = false
				})
				RegistryCmd.SetArgs(nil)
				RegistryCmd.SetOut(nil)
				RegistryCmd.SetErr(nil)
			})

			args := []string{"init", "cert-manager", "--path", t.TempDir()}
			for flag, value := range values {
				if flag != missing {
					args = append(args, "--"+flag, value)
				}
			}
			RegistryCmd.SetArgs(args)
			RegistryCmd.SetOut(io.Discard)
			RegistryCmd.SetErr(io.Discard)

			err := RegistryCmd.Execute()
			require.Error(t, err)
			assert.Contains(t, err.Error(), `required flag(s) "`+missing+`" not set`)
		})
	}
}
