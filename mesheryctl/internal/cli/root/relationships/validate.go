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

package relationships

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/spf13/cobra"
)

// RelValidationResult represents the result of validating a single relationship
type RelValidationResult struct {
	FilePath   string   `json:"filePath" yaml:"filePath"`
	EntityType string   `json:"entityType" yaml:"entityType"`
	EntityName string   `json:"entityName" yaml:"entityName"`
	IsValid    bool     `json:"isValid" yaml:"isValid"`
	Errors     []string `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// RelValidationSummary holds the overall validation results
type RelValidationSummary struct {
	TotalEntities int                   `json:"totalEntities" yaml:"totalEntities"`
	TotalValid    int                   `json:"totalValid" yaml:"totalValid"`
	TotalInvalid  int                   `json:"totalInvalid" yaml:"totalInvalid"`
	Results       []RelValidationResult `json:"results" yaml:"results"`
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate relationship definitions",
	Long: `Validate relationship definition files against the relationship schema.
Supports validating multiple relationships in one invocation.
Accepts local file paths, directories (recursively searched), and HTTP/HTTPS URLs.
Find more information at: https://docs.meshery.io/reference/mesheryctl/relationship/validate`,
	Example: `
// Validate a single relationship definition file
mesheryctl relationship validate -f relationship.json

// Validate multiple relationship files
mesheryctl relationship validate -f rel1.yaml -f rel2.yaml

// Validate all relationship files in a directory (recursively)
mesheryctl relationship validate -f ./relationships/

// Validate a relationship file from a URL
mesheryctl relationship validate -f https://example.com/relationship.json

// Output validation results as JSON
mesheryctl relationship validate -f ./relationships/ -o json

// Output validation results as YAML
mesheryctl relationship validate -f ./relationships/ -o yaml
	`,
	Args: func(_ *cobra.Command, _ []string) error {
		return nil
	},
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		files, _ := cmd.Flags().GetStringArray("file")
		if len(files) == 0 {
			return fmt.Errorf("at least one file, directory, or URL must be specified with -f/--file\n\nUsage: mesheryctl relationship validate -f <file|dir|URL> [-f ...]\nRun 'mesheryctl relationship validate --help' for more information")
		}
		outputFormat, _ := cmd.Flags().GetString("output-format")
		if outputFormat != "" {
			if err := display.ValidateOutputFormat(outputFormat); err != nil {
				return err
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		files, _ := cmd.Flags().GetStringArray("file")
		outputFormat, _ := cmd.Flags().GetString("output-format")

		// Collect all files to validate
		fileContents, err := collectRelationshipFiles(files)
		if err != nil {
			return err
		}

		if len(fileContents) == 0 {
			utils.Log.Info("No relationship definition files found")
			return nil
		}

		// Sort file paths for consistent output
		sortedPaths := make([]string, 0, len(fileContents))
		for p := range fileContents {
			sortedPaths = append(sortedPaths, p)
		}
		sort.Strings(sortedPaths)

		// Validate each file
		summary := RelValidationSummary{
			Results: make([]RelValidationResult, 0, len(fileContents)),
		}

		for _, path := range sortedPaths {
			result := validateRelationship(path, fileContents[path])
			summary.Results = append(summary.Results, result)
			summary.TotalEntities++
			if result.IsValid {
				summary.TotalValid++
			} else {
				summary.TotalInvalid++
			}
		}

		return displayRelValidationResults(summary, outputFormat)
	},
}

func init() {
	validateCmd.Flags().StringArrayP("file", "f", nil, "Path to a relationship definition file, directory, or URL (can be specified multiple times)")
	validateCmd.Flags().StringP("output-format", "o", "", "Output format [json|yaml] (default: table)")
}

// relSchemaHeader is used to detect entity type from schemaVersion
type relSchemaHeader struct {
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
}

// isRelationshipFile checks if the data represents a relationship definition
func isRelationshipFile(data []byte) bool {
	var header relSchemaHeader
	if err := encoding.Unmarshal(data, &header); err != nil {
		return false
	}
	return strings.HasPrefix(header.SchemaVersion, "relationships.meshery.io")
}

// collectRelationshipFiles gathers relationship file contents from the given paths
func collectRelationshipFiles(paths []string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	for _, path := range paths {
		if utils.IsValidUrl(path) {
			data, err := fetchRelFileFromURL(path)
			if err != nil {
				return nil, ErrValidateRelationship(err)
			}
			result[path] = data
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, ErrValidateRelationship(fmt.Errorf("cannot access %s: %w", path, err))
		}

		if info.IsDir() {
			err := filepath.Walk(path, func(filePath string, fi os.FileInfo, walkErr error) error {
				if walkErr != nil {
					return walkErr
				}
				if fi.IsDir() {
					return nil
				}
				ext := strings.ToLower(filepath.Ext(filePath))
				if ext != ".json" && ext != ".yaml" && ext != ".yml" {
					return nil
				}
				data, readErr := os.ReadFile(filePath)
				if readErr != nil {
					return nil // skip unreadable files during directory walk
				}
				// Only include relationship definition files
				if isRelationshipFile(data) {
					result[filePath] = data
				}
				return nil
			})
			if err != nil {
				return nil, ErrValidateRelationship(err)
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, ErrValidateRelationship(fmt.Errorf("failed to read file %s: %w", path, err))
			}
			// Always include explicitly specified files
			result[path] = data
		}
	}

	return result, nil
}

// fetchRelFileFromURL retrieves file content from a URL
func fetchRelFileFromURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: HTTP %d", url, resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// validateRelationship validates a single relationship definition file
func validateRelationship(filePath string, data []byte) RelValidationResult {
	if !isRelationshipFile(data) {
		return RelValidationResult{
			FilePath:   filePath,
			EntityType: "relationship",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{"file is not a relationship definition: missing or unrecognized schemaVersion (expected relationships.meshery.io/*)"},
		}
	}

	var rel relationship.RelationshipDefinition
	if err := encoding.Unmarshal(data, &rel); err != nil {
		return RelValidationResult{
			FilePath:   filePath,
			EntityType: "relationship",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{fmt.Sprintf("failed to parse relationship definition: %s", err.Error())},
		}
	}

	var errs []string
	if rel.SchemaVersion == "" {
		errs = append(errs, "schemaVersion is required")
	}
	if rel.Kind == "" {
		errs = append(errs, "kind is required")
	} else {
		validKinds := map[relationship.RelationshipDefinitionKind]bool{
			relationship.Edge:         true,
			relationship.Hierarchical: true,
			relationship.Sibling:      true,
		}
		if !validKinds[rel.Kind] {
			errs = append(errs, fmt.Sprintf("invalid kind %q: must be one of [edge, hierarchical, sibling]", rel.Kind))
		}
	}
	if rel.RelationshipType == "" {
		errs = append(errs, "type is required")
	}
	if rel.SubType == "" {
		errs = append(errs, "subType is required")
	}
	if rel.Version == "" {
		errs = append(errs, "version is required")
	}

	name := string(rel.Kind)
	if rel.SubType != "" {
		name = fmt.Sprintf("%s/%s", rel.Kind, rel.SubType)
	}
	if name == "/" || name == "" {
		name = filepath.Base(filePath)
	}

	return RelValidationResult{
		FilePath:   filePath,
		EntityType: "relationship",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// displayRelValidationResults renders the validation results in the specified format
func displayRelValidationResults(summary RelValidationSummary, outputFormat string) error {
	switch outputFormat {
	case "json":
		formatter := display.NewJSONOutputFormatter(summary)
		if err := formatter.Display(); err != nil {
			return err
		}
	case "yaml":
		formatter := display.NewYAMLOutputFormatter(summary)
		if err := formatter.Display(); err != nil {
			return err
		}
	default:
		header := []string{"FILE", "NAME", "KIND", "STATUS"}
		rows := make([][]string, 0, len(summary.Results))
		for _, result := range summary.Results {
			status := "PASS"
			if !result.IsValid {
				status = "FAIL"
			}
			rows = append(rows, []string{
				result.FilePath,
				result.EntityName,
				"relationship",
				status,
			})
		}
		utils.PrintToTable(header, rows, nil)

		// Show errors for failed entities
		for _, result := range summary.Results {
			if !result.IsValid {
				fmt.Printf("\n  %s %s:\n", utils.BoldString("Errors in"), result.FilePath)
				for _, e := range result.Errors {
					fmt.Printf("    - %s\n", e)
				}
			}
		}

		fmt.Printf("\n%s: %d relationships validated, %d passed, %d failed\n",
			utils.BoldString("SUMMARY"),
			summary.TotalEntities,
			summary.TotalValid,
			summary.TotalInvalid,
		)
	}

	if summary.TotalInvalid > 0 {
		return ErrValidateRelationship(fmt.Errorf("validation failed: %d of %d relationships invalid", summary.TotalInvalid, summary.TotalEntities))
	}
	return nil
}
