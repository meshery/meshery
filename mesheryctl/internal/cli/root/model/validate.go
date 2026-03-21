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
	"path/filepath"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/validation"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/encoding"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
	modelDef "github.com/meshery/schemas/models/v1beta1/model"
	"github.com/spf13/cobra"
)

var validateModelCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate model definitions",
	Long: `Validate model definitions and their constituent constructs (models, components, relationships)
for basic structural correctness and required fields. Supports validating multiple models in one invocation.
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
		fileContents, err := validation.CollectFiles(files, validation.AnyEntityFilter)
		if err != nil {
			return ErrValidateModel(err)
		}

		if len(fileContents) == 0 {
			utils.Log.Info("No model definition files found")
			return nil
		}

		sortedPaths := validation.SortedPaths(fileContents)

		// Validate each file
		summary := validation.Summary{
			Results: make([]validation.Result, 0, len(fileContents)),
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

		if err := validation.DisplayResults(summary, outputFormat, "entities"); err != nil {
			return err
		}
		if summary.TotalInvalid > 0 {
			return ErrValidateModel(fmt.Errorf("validation failed: %d of %d entities invalid", summary.TotalInvalid, summary.TotalEntities))
		}
		return nil
	},
}

func init() {
	validateModelCmd.Flags().StringArrayP("file", "f", nil, "Path to a model definition file, directory, or URL (can be specified multiple times)")
	validateModelCmd.Flags().StringP("output-format", "o", "", "Output format [json|yaml] (default: table)")
}

// validateEntityData validates a single entity file and returns the result
func validateEntityData(filePath string, data []byte) validation.Result {
	entityType := validation.DetectEntityType(data)
	if entityType == "" {
		return validation.Result{
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
		return validation.Result{
			FilePath:   filePath,
			EntityType: entityType,
			EntityName: filepath.Base(filePath),
			IsValid:    false,
			Errors:     []string{"unsupported entity type: " + entityType},
		}
	}
}

// validateModelDef validates a model definition file
func validateModelDef(filePath string, data []byte) validation.Result {
	var m modelDef.ModelDefinition
	if err := encoding.Unmarshal(data, &m); err != nil {
		return validation.Result{
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

	return validation.Result{
		FilePath:   filePath,
		EntityType: "model",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// validateComponentDef validates a component definition file
func validateComponentDef(filePath string, data []byte) validation.Result {
	var comp component.ComponentDefinition
	if err := encoding.Unmarshal(data, &comp); err != nil {
		return validation.Result{
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

	return validation.Result{
		FilePath:   filePath,
		EntityType: "component",
		EntityName: name,
		IsValid:    len(errs) == 0,
		Errors:     errs,
	}
}

// validateRelationshipDef validates a relationship definition file
func validateRelationshipDef(filePath string, data []byte) validation.Result {
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

