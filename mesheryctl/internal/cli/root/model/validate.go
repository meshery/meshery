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

package model

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
	"github.com/meshery/schemas/models/v1beta1/component"
	modelDef "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/spf13/cobra"
)

// ValidationResult represents the result of validating a single entity
type ValidationResult struct {
	FilePath   string   `json:"filePath" yaml:"filePath"`
	EntityType string   `json:"entityType" yaml:"entityType"`
	EntityName string   `json:"entityName" yaml:"entityName"`
	IsValid    bool     `json:"isValid" yaml:"isValid"`
	Errors     []string `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// ValidationSummary holds the overall validation results
type ValidationSummary struct {
	TotalEntities int                `json:"totalEntities" yaml:"totalEntities"`
	TotalValid    int                `json:"totalValid" yaml:"totalValid"`
	TotalInvalid  int                `json:"totalInvalid" yaml:"totalInvalid"`
	Results       []ValidationResult `json:"results" yaml:"results"`
}

var validateModelCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate model definitions",
	Long: `Validate model definitions and their constituent constructs (models, components, relationships)
against their schemas. Supports validating multiple models in one invocation.
Accepts local file paths, directories, and HTTP/HTTPS URLs.
Find more information at: https://docs.meshery.io/reference/mesheryctl/model/validate`,
	Example: `
// Validate a single model definition file
mesheryctl model validate -f model.json

// Validate multiple files
mesheryctl model validate -f model.json -f component.yaml -f relationship.json

// Validate all model definition files in a directory (recursively)
mesheryctl model validate -f ./my-model/

// Validate a model file from a URL
mesheryctl model validate -f https://example.com/model.json

// Output validation results as JSON
mesheryctl model validate -f ./my-model/ -o json

// Output validation results as YAML
mesheryctl model validate -f ./my-model/ -o yaml
	`,
	Args: func(_ *cobra.Command, _ []string) error {
		return nil
	},
	PreRunE: func(cmd *cobra.Command, _ []string) error {
		files, _ := cmd.Flags().GetStringArray("file")
		if len(files) == 0 {
			return fmt.Errorf("at least one file, directory, or URL must be specified with -f/--file\n\nUsage: mesheryctl model validate -f <file|dir|URL> [-f ...]\nRun 'mesheryctl model validate --help' for more information")
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
		fileContents, err := collectModelFiles(files)
		if err != nil {
			return err
		}

		if len(fileContents) == 0 {
			utils.Log.Info("No model definition files found")
			return nil
		}

		// Sort file paths for consistent output
		sortedPaths := make([]string, 0, len(fileContents))
		for p := range fileContents {
			sortedPaths = append(sortedPaths, p)
		}
		sort.Strings(sortedPaths)

		// Validate each file
		summary := ValidationSummary{
			Results: make([]ValidationResult, 0, len(fileContents)),
		}

		for _, path := range sortedPaths {
			result := validateEntityData(path, fileContents[path])
			summary.Results = append(summary.Results, result)
			summary.TotalEntities++
			if result.IsValid {
				summary.TotalValid++
			} else {
				summary.TotalInvalid++
			}
		}

		return displayModelValidationResults(summary, outputFormat)
	},
}

func init() {
	validateModelCmd.Flags().StringArrayP("file", "f", nil, "Path to a model definition file, directory, or URL (can be specified multiple times)")
	validateModelCmd.Flags().StringP("output-format", "o", "", "Output format [json|yaml] (default: table)")
}

// entityHeader is used to detect the entity type from the schemaVersion field
type entityHeader struct {
	SchemaVersion string `json:"schemaVersion" yaml:"schemaVersion"`
}

// collectModelFiles gathers file contents from the given paths (files, directories, or URLs)
func collectModelFiles(paths []string) (map[string][]byte, error) {
	result := make(map[string][]byte)

	for _, path := range paths {
		if utils.IsValidUrl(path) {
			data, err := fetchFileFromURL(path)
			if err != nil {
				return nil, ErrValidateModel(err)
			}
			result[path] = data
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			return nil, ErrValidateModel(fmt.Errorf("cannot access %s: %w", path, err))
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
				// Only include files with a recognized schemaVersion
				if detectEntityType(data) != "" {
					result[filePath] = data
				}
				return nil
			})
			if err != nil {
				return nil, ErrValidateModel(err)
			}
		} else {
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, ErrValidateModel(fmt.Errorf("failed to read file %s: %w", path, err))
			}
			result[path] = data
		}
	}

	return result, nil
}

// fetchFileFromURL retrieves file content from a URL
func fetchFileFromURL(url string) ([]byte, error) {
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

// detectEntityType determines the entity type from the schemaVersion field
func detectEntityType(data []byte) string {
	var header entityHeader
	if err := encoding.Unmarshal(data, &header); err != nil {
		return ""
	}
	switch {
	case strings.HasPrefix(header.SchemaVersion, "models.meshery.io"):
		return "model"
	case strings.HasPrefix(header.SchemaVersion, "components.meshery.io"):
		return "component"
	case strings.HasPrefix(header.SchemaVersion, "relationships.meshery.io"):
		return "relationship"
	default:
		return ""
	}
}

// validateEntityData validates a single entity file and returns the result
func validateEntityData(filePath string, data []byte) ValidationResult {
	entityType := detectEntityType(data)
	if entityType == "" {
		return ValidationResult{
			FilePath:   filePath,
			EntityType: "unknown",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{"unable to determine entity type: missing or unrecognized schemaVersion"},
		}
	}

	switch entityType {
	case "model":
		return validateModelDef(filePath, data)
	case "component":
		return validateComponentDef(filePath, data)
	case "relationship":
		return validateRelationshipDef(filePath, data)
	default:
		return ValidationResult{
			FilePath:   filePath,
			EntityType: entityType,
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{"unsupported entity type: " + entityType},
		}
	}
}

// validateModelDef validates a model definition file
func validateModelDef(filePath string, data []byte) ValidationResult {
	var m modelDef.ModelDefinition
	if err := encoding.Unmarshal(data, &m); err != nil {
		return ValidationResult{
			FilePath:   filePath,
			EntityType: "model",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{fmt.Sprintf("failed to parse model definition: %s", err.Error())},
		}
	}

	var errs []string
	if m.Name == "" {
		errs = append(errs, "name is required")
	}
	if m.SchemaVersion == "" {
		errs = append(errs, "schemaVersion is required")
	}
	if m.Version == "" {
		errs = append(errs, "version is required")
	}

	name := m.Name
	if name == "" {
		name = filepath.Base(filePath)
	}

	return ValidationResult{
		FilePath:   filePath,
		EntityType: "model",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// validateComponentDef validates a component definition file
func validateComponentDef(filePath string, data []byte) ValidationResult {
	var comp component.ComponentDefinition
	if err := encoding.Unmarshal(data, &comp); err != nil {
		return ValidationResult{
			FilePath:   filePath,
			EntityType: "component",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{fmt.Sprintf("failed to parse component definition: %s", err.Error())},
		}
	}

	var errs []string
	if comp.SchemaVersion == "" {
		errs = append(errs, "schemaVersion is required")
	}
	if comp.Component.Kind == "" {
		errs = append(errs, "component.kind is required")
	}
	if comp.Component.Version == "" {
		errs = append(errs, "component.version is required")
	}

	name := comp.DisplayName
	if name == "" {
		name = comp.Component.Kind
	}
	if name == "" {
		name = filepath.Base(filePath)
	}

	return ValidationResult{
		FilePath:   filePath,
		EntityType: "component",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// validateRelationshipDef validates a relationship definition file
func validateRelationshipDef(filePath string, data []byte) ValidationResult {
	var rel relationship.RelationshipDefinition
	if err := encoding.Unmarshal(data, &rel); err != nil {
		return ValidationResult{
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

	return ValidationResult{
		FilePath:   filePath,
		EntityType: "relationship",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// displayModelValidationResults renders the validation results in the specified format
func displayModelValidationResults(summary ValidationSummary, outputFormat string) error {
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
		header := []string{"FILE", "TYPE", "NAME", "STATUS"}
		rows := make([][]string, 0, len(summary.Results))
		for _, result := range summary.Results {
			status := "PASS"
			if !result.IsValid {
				status = "FAIL"
			}
			rows = append(rows, []string{
				result.FilePath,
				result.EntityType,
				result.EntityName,
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

		fmt.Printf("\n%s: %d entities validated, %d passed, %d failed\n",
			utils.BoldString("SUMMARY"),
			summary.TotalEntities,
			summary.TotalValid,
			summary.TotalInvalid,
		)
	}

	if summary.TotalInvalid > 0 {
		return ErrValidateModel(fmt.Errorf("validation failed: %d of %d entities invalid", summary.TotalInvalid, summary.TotalEntities))
	}
	return nil
}
