// Copyright Meshery Authors
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

package design

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models/pattern/core"
	patternutils "github.com/meshery/meshery/server/models/pattern/utils"
	schemaCore "github.com/meshery/schemas/models/core"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/spf13/cobra"
)

// githubDownloadTimeout bounds the download of a design file from GitHub.
const githubDownloadTimeout = 30 * time.Second

var (
	availableSubcommands []*cobra.Command
	file                 string
)

// DesignCmd represents the root command for design commands
var DesignCmd = &cobra.Command{
	Use:   "design",
	Short: "Manage cloud native designs",
	Long: `Manage cloud and cloud native infrastructure using predefined designs.
Find more information at: https://docs.meshery.io/reference/mesheryctl/design`,
	Example: `
// Apply design file:
mesheryctl design apply --file [path to design file | URL of the file]

// Delete design file:
mesheryctl design delete --file [path to design file]

// View design file:
mesheryctl design view [design-name | ID]

// List all designs:
mesheryctl design list

// Import a design:
mesheryctl design import -f [file/URL] -n [name]

// Export a design:
mesheryctl design export [design-name | ID]
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			suggestions := make([]string, 0)
			for _, subcmd := range availableSubcommands {
				if strings.HasPrefix(subcmd.Name(), args[0]) {
					suggestions = append(suggestions, subcmd.Name())
				}
			}
			if len(suggestions) > 0 {
				return ErrInvalidCommand(args[0], suggestions)
			}
			return ErrInvalidCommand(args[0], []string{})
		}
		return nil
	},
}

func init() {
	DesignCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, deleteCmd, viewCmd, listCmd, importCmd, deployDesignCmd, exportCmd, designUndeployCmd, evaluateCmd}
	DesignCmd.AddCommand(availableSubcommands...)
}

// ValidSourceTypes contains the schema-driven canonical source types from github.com/meshery/schemas/models/core.
var ValidSourceTypes = []schemaCore.IaCFileTypes{
	schemaCore.HelmChart,
	schemaCore.DockerCompose,
	schemaCore.K8sManifest,
	schemaCore.K8sKustomize,
	schemaCore.MesheryDesign,
}

// sourceTypesMap maps canonical schema values, shorthand aliases, and legacy titles
// (all case-insensitive) to canonical schemaCore.IaCFileTypes.
var sourceTypesMap = map[string]schemaCore.IaCFileTypes{
	// Canonical schema values
	string(schemaCore.HelmChart):     schemaCore.HelmChart,
	string(schemaCore.DockerCompose): schemaCore.DockerCompose,
	string(schemaCore.K8sManifest):   schemaCore.K8sManifest,
	string(schemaCore.K8sKustomize):  schemaCore.K8sKustomize,
	string(schemaCore.MesheryDesign): schemaCore.MesheryDesign,

	// Shorthand aliases
	"helm":      schemaCore.HelmChart,
	"compose":   schemaCore.DockerCompose,
	"docker":    schemaCore.DockerCompose,
	"manifest":  schemaCore.K8sManifest,
	"k8s":       schemaCore.K8sManifest,
	"kustomize": schemaCore.K8sKustomize,
	"design":    schemaCore.MesheryDesign,

	// Human-readable titles (for backward compatibility)
	"helm chart":          schemaCore.HelmChart,
	"docker compose":      schemaCore.DockerCompose,
	"kubernetes manifest": schemaCore.K8sManifest,
	"k8s kustomize":       schemaCore.K8sKustomize,
	"meshery design":      schemaCore.MesheryDesign,
}

// GetDesignSourceTypes returns the slice of canonical source type strings.
func GetDesignSourceTypes() []string {
	types := make([]string, 0, len(ValidSourceTypes))
	for _, t := range ValidSourceTypes {
		types = append(types, string(t))
	}
	return types
}

// ValidateAndNormalizeSourceType validates and normalizes a source type string into its canonical schemaCore.IaCFileTypes.
func ValidateAndNormalizeSourceType(sType string) (schemaCore.IaCFileTypes, error) {
	clean := strings.ToLower(strings.TrimSpace(sType))
	if clean == "" {
		return "", ErrInValidSource(sType, GetDesignSourceTypes())
	}
	if canonical, ok := sourceTypesMap[clean]; ok {
		return canonical, nil
	}
	return "", ErrInValidSource(sType, GetDesignSourceTypes())
}

func getDesignSourceTypes() ([]string, error) {
	return GetDesignSourceTypes(), nil
}

func retrieveProvidedSourceType(sType string, validDesignSourceTypes []string) (string, error) {
	canonical, err := ValidateAndNormalizeSourceType(sType)
	if err == nil {
		return string(canonical), nil
	}
	clean := strings.ToLower(strings.TrimSpace(sType))
	for _, v := range validDesignSourceTypes {
		if strings.EqualFold(clean, v) {
			return clean, nil
		}
	}
	return "", ErrInValidSource(sType, validDesignSourceTypes)
}

// fetchDesignByID fetches a design from the server by ID and parses it into a PatternFile.
func fetchDesignByID(baseUrl, designID string) (pattern.PatternFile, error) {
	dataURL := fmt.Sprintf("%s/api/pattern/%s", baseUrl, designID)
	patternData, err := fetchPatternData(dataURL)
	if err != nil {
		return pattern.PatternFile{}, err
	}

	pf, err := core.NewPatternFile([]byte(patternData.PatternFile))
	if err != nil {
		return pattern.PatternFile{}, ErrParseDesignFile(err)
	}
	// core.NewPatternFile returns a v1beta3/design.PatternFile, but the
	// evaluation engine (EvaluationRequest/EvaluationResponse) is a documented
	// v1beta1/pattern carve-out, so bridge it back to v1beta1.
	bridged, err := patternutils.PatternV1beta3ToV1beta1(&pf)
	if err != nil {
		return pattern.PatternFile{}, ErrParseDesignFile(err)
	}
	if bridged == nil {
		return pattern.PatternFile{}, ErrParseDesignFile(errors.New("bridged pattern file is nil"))
	}
	return *bridged, nil
}

// fetchDesignFromGithub downloads a design file from a github.com or
// raw.githubusercontent.com URL. The server has no endpoint that reads a design
// from a URL without saving it, so the file is fetched here. It uses a plain
// request, because utils.NewRequest would send the Meshery token to GitHub.
func fetchDesignFromGithub(fileURL string) ([]byte, error) {
	repoURL, filePath, err := utils.ParseURLGithub(fileURL)
	if err != nil {
		return nil, utils.ErrParseGithubFile(err, fileURL)
	}

	rawURL := repoURL
	if filePath != "" {
		rawURL = strings.Replace(repoURL, "https://github.com/", "https://raw.githubusercontent.com/", 1) + "/" + filePath
	} else {
		rawURL = fileURL
	}

	client := &http.Client{Timeout: githubDownloadTimeout}
	resp, err := client.Get(rawURL)
	if err != nil {
		return nil, utils.ErrFailRequest(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, utils.ErrResponseStatus(resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrReadFromBody(err)
	}
	return content, nil
}

func isGithubURL(fileURL string) bool {
	u, err := url.Parse(fileURL)
	if err != nil {
		return false
	}
	return u.Hostname() == "github.com" || u.Hostname() == "raw.githubusercontent.com"
}

// designFileName derives a clean file name from a file path or URL, stripping query parameters.
func designFileName(filePathOrURL string) string {
	if utils.IsValidUrl(filePathOrURL) {
		if u, err := url.Parse(filePathOrURL); err == nil && u.Path != "" {
			base := path.Base(u.Path)
			if base != "." && base != "/" && base != "" {
				return base
			}
		}
		return "design.yaml"
	}
	return filepath.Base(filePathOrURL)
}

// fetchDesignFile downloads a design file from a URL. If the URL points to GitHub,
// it is normalized and fetched via fetchDesignFromGithub. For other HTTP/HTTPS URLs,
// it is fetched directly.
func fetchDesignFile(fileURL string) ([]byte, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return nil, utils.ErrFailRequest(err)
	}
	scheme := strings.ToLower(u.Scheme)
	if scheme != "http" && scheme != "https" {
		return nil, ErrInvalidDesignFile(fmt.Errorf("unsupported URL scheme %q: only http and https are supported", u.Scheme))
	}

	if isGithubURL(fileURL) {
		return fetchDesignFromGithub(fileURL)
	}

	client := &http.Client{Timeout: githubDownloadTimeout}
	resp, err := client.Get(fileURL)
	if err != nil {
		return nil, utils.ErrFailRequest(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, utils.ErrResponseStatus(resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.ErrReadFromBody(err)
	}
	return content, nil
}
