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
	"path/filepath"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/validation"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate relationship definitions",
	Long: `Validate relationship definition files for basic structural correctness
and required fields. Supports validating multiple relationships in one invocation.
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

		// Collect all files to validate (only relationship definitions from directories)
		fileContents, err := validation.CollectFiles(files, validation.RelationshipOnlyFilter)
		if err != nil {
			return ErrValidateRelationship(err)
		}

		if len(fileContents) == 0 {
			utils.Log.Info("No relationship definition files found")
			return nil
		}

		sortedPaths := validation.SortedPaths(fileContents)

		// Validate each file
		summary := validation.Summary{
			Results: make([]validation.Result, 0, len(fileContents)),
		}

		for _, path := range sortedPaths {
			result := validateRelationshipFile(path, fileContents[path])
			summary.Results = append(summary.Results, result)
			summary.TotalEntities++
			if result.IsValid {
				summary.TotalValid++
			} else {
				summary.TotalInvalid++
			}
		}

		if err := validation.DisplayResults(summary, outputFormat, "relationships"); err != nil {
			return err
		}
		if summary.TotalInvalid > 0 {
			return ErrValidateRelationship(fmt.Errorf("validation failed: %d of %d relationships invalid", summary.TotalInvalid, summary.TotalEntities))
		}
		return nil
	},
}

func init() {
	validateCmd.Flags().StringArrayP("file", "f", nil, "Path to a relationship definition file, directory, or URL (can be specified multiple times)")
	validateCmd.Flags().StringP("output-format", "o", "", "Output format [json|yaml] (default: table)")
}

// validateRelationshipFile validates a single relationship definition file
func validateRelationshipFile(filePath string, data []byte) validation.Result {
	entityType := validation.DetectEntityType(data)
	if entityType != "relationship" {
		return validation.Result{
			FilePath:   filePath,
			EntityType: "relationship",
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{"file is not a relationship definition: missing or unrecognized schemaVersion (expected relationships.meshery.io/*)"},
		}
	}

	var rel relationship.RelationshipDefinition
	if err := encoding.Unmarshal(data, &rel); err != nil {
		return validation.Result{
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

	return validation.Result{
		FilePath:   filePath,
		EntityType: "relationship",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}
