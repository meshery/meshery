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
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	// modelSourceSchemaVersion is the schema the scaffold conforms to
	// (meshery/schemas, constructs/v1alpha1/model_source).
	modelSourceSchemaVersion = "modelsources.meshery.io/v1alpha1"
	modelSourceFileName      = "model-source.yaml"

	sourceKindGithub      = "github"
	sourceKindArtifactHub = "artifacthub"
	sourceKindHTTP        = "http"

	registryInitDirPerm  = 0o755
	registryInitFilePerm = 0o644
)

// Field limits from the v1alpha1 model_source schema.
const (
	maxNameLength        = 100
	maxDisplayNameLength = 100
	maxRegistrantLength  = 100
	maxCategoryLength    = 100
	maxDescriptionLength = 2000
	maxSourceURLLength   = 2048
)

var modelSourceNamePattern = regexp.MustCompile(`^[a-z0-9-]+$`)

// placeholderSVG is written for each logo until the contributor replaces it.
const placeholderSVG = `<svg xmlns="http://www.w3.org/2000/svg" width="32" height="32" viewBox="0 0 32 32"><rect width="32" height="32" rx="4" fill="#00B39F"/></svg>
`

// svgFiles maps each placeholder logo to the metadata field that references it.
var svgFiles = []string{"color.svg", "white.svg", "complete.svg"}

// modelSource mirrors the required and scaffolded fields of the v1alpha1
// model_source schema.
// TODO: switch to the generated types in github.com/meshery/schemas/models/v1alpha1/model_source
// and validate against the published schema once meshery/schemas#1293 is released.
type modelSource struct {
	SchemaVersion string              `yaml:"schemaVersion"`
	Name          string              `yaml:"name"`
	DisplayName   string              `yaml:"displayName"`
	Registrant    string              `yaml:"registrant"`
	Category      string              `yaml:"category"`
	SubCategory   string              `yaml:"subCategory"`
	Description   string              `yaml:"description"`
	Source        modelSourceLocation `yaml:"source"`
	Metadata      modelSourceMetadata `yaml:"metadata"`
}

type modelSourceLocation struct {
	URL  string `yaml:"url"`
	Kind string `yaml:"kind"`
}

type modelSourceMetadata struct {
	SvgColor    string `yaml:"svgColor"`
	SvgWhite    string `yaml:"svgWhite"`
	SvgComplete string `yaml:"svgComplete"`
}

type cmdRegistryInitFlags struct {
	Source      string
	SourceKind  string
	DisplayName string
	Registrant  string
	Category    string
	SubCategory string
	Description string
	Path        string
	Force       bool
}

var registryInitFlags cmdRegistryInitFlags

var initCmd = &cobra.Command{
	Use:   "init <model-name>",
	Short: "Scaffold a model-source.yaml for a new model",
	Long: `Scaffold models/<model-name>/model-source.yaml, the canonical generation inputs for a model,
along with placeholder SVGs and an overrides/ directory, so a new model can be added by pull request.

Run it from the root of a meshery/meshery checkout so the model lands in models/.
Find more information at: https://docs.meshery.io/reference/mesheryctl/registry/init`,
	Example: `
// Scaffold a model whose source is an Artifact Hub package
mesheryctl registry init cert-manager --source https://artifacthub.io/packages/helm/cert-manager/cert-manager --category "Security & Compliance" --sub-category "Key Management" --description "Cloud native certificate management for Kubernetes."

// Scaffold a model from a GitHub repository, with an explicit registrant
mesheryctl registry init my-operator --source git://github.com/example/my-operator/main/deploy/crds --registrant meshery --category Provisioning --sub-category "Automation & Configuration" --description "Operator for my workloads."

// Overwrite an existing model-source.yaml (SVGs and overrides are kept)
mesheryctl registry init cert-manager --source https://artifacthub.io/packages/helm/cert-manager/cert-manager --category "Security & Compliance" --sub-category "Key Management" --description "Cloud native certificate management for Kubernetes." --force
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("exactly one model name is required\n\nUsage: mesheryctl registry init <model-name> --source <url> [flags]"))
		}
		return validateModelSourceName(args[0])
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ms, err := buildModelSource(args[0], registryInitFlags)
		if err != nil {
			return err
		}

		modelDir := filepath.Join(registryInitFlags.Path, ms.Name)
		written, err := scaffoldModelSource(modelDir, ms, registryInitFlags.Force)
		if err != nil {
			return err
		}

		utils.Log.Infof("Scaffolded model %q in %s", ms.Name, modelDir)
		for _, file := range written {
			utils.Log.Infof("  %s", file)
		}
		utils.Log.Info(registryInitNextSteps(modelDir))
		return nil
	},
}

func init() {
	initCmd.Flags().StringVarP(&registryInitFlags.Source, "source", "s", "", "(required) upstream URL of the model source: a GitHub repository, an Artifact Hub package, or an HTTP manifest or chart")
	initCmd.Flags().StringVar(&registryInitFlags.SourceKind, "source-kind", "", "kind of source [github|artifacthub|http] (default: detected from --source)")
	initCmd.Flags().StringVar(&registryInitFlags.DisplayName, "display-name", "", "human-readable model name (default: derived from the model name)")
	initCmd.Flags().StringVar(&registryInitFlags.Registrant, "registrant", "", "registrant of the model (default: derived from the source kind)")
	initCmd.Flags().StringVar(&registryInitFlags.Category, "category", "", "(required) model category")
	initCmd.Flags().StringVar(&registryInitFlags.SubCategory, "sub-category", "", "(required) model subcategory")
	initCmd.Flags().StringVar(&registryInitFlags.Description, "description", "", "(required) summary description of the model")
	initCmd.Flags().StringVarP(&registryInitFlags.Path, "path", "p", "models", "directory the model folder is created under")
	initCmd.Flags().BoolVar(&registryInitFlags.Force, "force", false, "overwrite an existing model-source.yaml (SVGs, overrides and version folders are never touched)")

	_ = initCmd.MarkFlagRequired("source")
	_ = initCmd.MarkFlagRequired("category")
	_ = initCmd.MarkFlagRequired("sub-category")
	_ = initCmd.MarkFlagRequired("description")
}

// validateModelSourceName applies the schema's rules for the name field.
func validateModelSourceName(name string) error {
	if name == "" || len(name) > maxNameLength || !modelSourceNamePattern.MatchString(name) {
		return ErrRegistryInit(fmt.Errorf("invalid model name %q: use 1-%d lowercase letters, digits or hyphens (for example cert-manager)", name, maxNameLength))
	}
	return nil
}

// buildModelSource assembles and validates the manifest from the command's
// arguments and flags, filling in derived defaults.
func buildModelSource(name string, flags cmdRegistryInitFlags) (modelSource, error) {
	source := strings.TrimSpace(flags.Source)

	kind := strings.ToLower(strings.TrimSpace(flags.SourceKind))
	if kind == "" {
		detected, err := detectSourceKind(source)
		if err != nil {
			return modelSource{}, err
		}
		kind = detected
	} else if !isValidSourceKind(kind) {
		return modelSource{}, ErrRegistryInit(fmt.Errorf("invalid --source-kind %q: must be one of %s, %s or %s", flags.SourceKind, sourceKindGithub, sourceKindArtifactHub, sourceKindHTTP))
	}

	displayName := strings.TrimSpace(flags.DisplayName)
	if displayName == "" {
		displayName = displayNameFromModelName(name)
	}

	registrant := strings.TrimSpace(flags.Registrant)
	if registrant == "" {
		registrant = defaultRegistrant(kind)
	}

	ms := modelSource{
		SchemaVersion: modelSourceSchemaVersion,
		Name:          name,
		DisplayName:   displayName,
		Registrant:    registrant,
		Category:      strings.TrimSpace(flags.Category),
		SubCategory:   strings.TrimSpace(flags.SubCategory),
		Description:   strings.TrimSpace(flags.Description),
		Source: modelSourceLocation{
			URL:  source,
			Kind: kind,
		},
		Metadata: modelSourceMetadata{
			SvgColor:    "./svg/color.svg",
			SvgWhite:    "./svg/white.svg",
			SvgComplete: "./svg/complete.svg",
		},
	}
	if err := validateModelSource(ms); err != nil {
		return modelSource{}, err
	}
	return ms, nil
}

// validateModelSource checks the schema's required fields and length limits,
// so the command never writes a manifest that would fail validation.
func validateModelSource(ms modelSource) error {
	checks := []struct {
		field string
		value string
		max   int
	}{
		{"--display-name", ms.DisplayName, maxDisplayNameLength},
		{"--registrant", ms.Registrant, maxRegistrantLength},
		{"--category", ms.Category, maxCategoryLength},
		{"--sub-category", ms.SubCategory, maxCategoryLength},
		{"--description", ms.Description, maxDescriptionLength},
		{"--source", ms.Source.URL, maxSourceURLLength},
	}
	for _, c := range checks {
		if c.value == "" {
			return ErrRegistryInit(fmt.Errorf("%s must not be empty", c.field))
		}
		if len(c.value) > c.max {
			return ErrRegistryInit(fmt.Errorf("%s must be at most %d characters", c.field, c.max))
		}
	}
	return nil
}

func isValidSourceKind(kind string) bool {
	switch kind {
	case sourceKindGithub, sourceKindArtifactHub, sourceKindHTTP:
		return true
	}
	return false
}

// detectSourceKind infers source.kind from the source URL.
func detectSourceKind(source string) (string, error) {
	if source == "" {
		return "", ErrRegistryInit(fmt.Errorf("--source must not be empty"))
	}
	u, err := url.Parse(source)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", ErrRegistryInit(fmt.Errorf("--source %q is not a valid URL", source))
	}
	host := strings.TrimPrefix(strings.ToLower(u.Hostname()), "www.")
	switch {
	case u.Scheme == "git" || host == "github.com":
		return sourceKindGithub, nil
	case host == "artifacthub.io":
		return sourceKindArtifactHub, nil
	case u.Scheme == "http" || u.Scheme == "https":
		return sourceKindHTTP, nil
	}
	return "", ErrRegistryInit(fmt.Errorf("cannot detect the source kind of %q; pass --source-kind (%s, %s or %s)", source, sourceKindGithub, sourceKindArtifactHub, sourceKindHTTP))
}

// defaultRegistrant derives the registrant from the source kind; --registrant overrides it.
func defaultRegistrant(kind string) string {
	switch kind {
	case sourceKindGithub:
		return "github"
	case sourceKindArtifactHub:
		return "artifacthub"
	}
	return "meshery"
}

// displayNameFromModelName turns "cert-manager" into "Cert Manager".
func displayNameFromModelName(name string) string {
	words := strings.Split(name, "-")
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(strings.Fields(strings.Join(words, " ")), " ")
}

// scaffoldModelSource writes the manifest and any missing placeholders under
// modelDir and returns the files it wrote. An existing model-source.yaml is
// only replaced when force is set; existing SVGs, overrides and version
// folders are never modified.
func scaffoldModelSource(modelDir string, ms modelSource, force bool) ([]string, error) {
	if info, err := os.Stat(modelDir); err == nil && !info.IsDir() {
		return nil, ErrRegistryInit(fmt.Errorf("%s exists and is not a directory", modelDir))
	} else if err != nil && !os.IsNotExist(err) {
		return nil, ErrRegistryInit(fmt.Errorf("cannot use %s as a model folder: %w", modelDir, err))
	}

	manifestPath := filepath.Join(modelDir, modelSourceFileName)
	if _, err := os.Stat(manifestPath); err == nil && !force {
		return nil, ErrRegistryInit(fmt.Errorf("%s already exists; pass --force to overwrite it", manifestPath))
	}

	content, err := marshalModelSource(ms)
	if err != nil {
		return nil, ErrRegistryInit(err)
	}

	svgDir := filepath.Join(modelDir, "svg")
	overridesDir := filepath.Join(modelDir, "overrides", "components")
	for _, dir := range []string{svgDir, overridesDir} {
		if err := os.MkdirAll(dir, registryInitDirPerm); err != nil {
			return nil, ErrRegistryInit(err)
		}
	}

	if err := os.WriteFile(manifestPath, content, registryInitFilePerm); err != nil {
		return nil, ErrRegistryInit(err)
	}
	written := []string{manifestPath}

	placeholders := map[string]string{filepath.Join(overridesDir, ".gitkeep"): ""}
	for _, name := range svgFiles {
		placeholders[filepath.Join(svgDir, name)] = placeholderSVG
	}
	for _, path := range sortedKeys(placeholders) {
		created, err := writeIfMissing(path, placeholders[path])
		if err != nil {
			return nil, ErrRegistryInit(err)
		}
		if created {
			written = append(written, path)
		}
	}
	return written, nil
}

// marshalModelSource renders the manifest with 2-space indentation, matching
// the model_source templates in meshery/schemas.
func marshalModelSource(ms modelSource) ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(ms); err != nil {
		return nil, err
	}
	if err := enc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// writeIfMissing creates path with content unless it already exists.
func writeIfMissing(path, content string) (bool, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, registryInitFilePerm)
	if os.IsExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if _, err := f.WriteString(content); err != nil {
		_ = f.Close()
		return false, err
	}
	return true, f.Close()
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func registryInitNextSteps(modelDir string) string {
	return fmt.Sprintf(`
Next steps:
1. Review %s and add any optional fields (version, metadata, publish, website)
2. Replace the placeholder logos in %s
3. Add per-component overrides under %s if needed
4. Make sure the model is under models/ at the root of a meshery/meshery checkout, then open a pull request`,
		filepath.Join(modelDir, modelSourceFileName),
		filepath.Join(modelDir, "svg"),
		filepath.Join(modelDir, "overrides", "components"))
}
