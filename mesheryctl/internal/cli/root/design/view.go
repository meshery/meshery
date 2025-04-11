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
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag   bool
	outFormatFlag string
)

var linkDocPatternView = map[string]string{
	"link":    "![pattern-view-usage](/assets/img/mesheryctl/patternView.png)",
	"caption": "Usage of mesheryctl design view",
}

var viewCmd = &cobra.Command{
	Use:   "view design name",
	Short: "Display a design content",
	Long:  `Display the content of a specific design based on name or id`,
	Args:  cobra.MaximumNArgs(1),
	Example: `
// view a design
mesheryctl design view [design-name | ID]
    `,
	Annotations: linkDocPatternView,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		pattern := ""
		isID := false
		// if pattern name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New("-a cannot be used when [design-name|design-id] is specified")
			}
			pattern, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				return ErrPatternInvalidNameOrID(err)
			}
		}
		url := mctlCfg.GetBaseMesheryURL()
		if len(pattern) == 0 {
			if viewAllFlag {
				url += "/api/pattern?populate=pattern_file&pagesize=10000"
			} else {
				return errors.New(utils.DesignViewError("Design name or ID is not specified. Use `-a` to view all designs"))
			}
		} else if isID {
			// if pattern is a valid uuid, then directly fetch the pattern
			url += "/api/pattern/" + pattern
		} else {
			// else search pattern by name
			url += "/api/pattern?populate=pattern_file&search=" + pattern
		}

		req, err := utils.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return utils.ErrReadResponseBody(err)
		}

		var dat map[string]interface{}
		if err = json.Unmarshal(body, &dat); err != nil {
			return utils.ErrUnmarshal(err)
		}

		// Variables to track design structure for summary
		var components []interface{}
		var relationships []interface{}
		designName := ""
		sourceType := ""

		if isID {
			// Extract design info and find components/relationships
			designName, sourceType = extractDesignInfo(dat)
			components, relationships = findComponentsAndRelationships(dat)

			hasComponents := len(components) > 0
			hasRelationships := len(relationships) > 0

			if !hasComponents {
				utils.Log.Info("No components found for this design.")
			}

			if !hasRelationships {
				utils.Log.Info("No relationships found for this design.")
			}

			if body, err = json.MarshalIndent(dat, "", "  "); err != nil {
				return utils.ErrMarshalIndent(err)
			}
		} else if viewAllFlag {
			// only keep the pattern key from the response when viewing all the patterns
			if body, err = json.MarshalIndent(map[string]interface{}{"patterns": dat["patterns"]}, "", "  "); err != nil {
				return utils.ErrMarshalIndent(err)
			}
		} else {
			// use the first match from the result when searching by pattern name
			patternsVal, exists := dat["patterns"]
			if !exists || patternsVal == nil {
				return ErrDesignNotFound()
			}

			patterns, ok := patternsVal.([]interface{})
			if !ok || len(patterns) == 0 {
				return ErrDesignNotFound()
			}

			firstMatchVal := patterns[0]
			firstMatch, ok := firstMatchVal.(map[string]interface{})
			if !ok || firstMatch == nil {
				return ErrDesignNotFound()
			}

			// Extract design info and find components/relationships
			designName, sourceType = extractDesignInfo(firstMatch)
			components, relationships = findComponentsAndRelationships(firstMatch)

			hasComponents := len(components) > 0
			hasRelationships := len(relationships) > 0

			if !hasComponents {
				utils.Log.Info("No components found for this design.")
			}

			if !hasRelationships {
				utils.Log.Info("No relationships found for this design.")
			}

			if body, err = json.MarshalIndent(firstMatch, "", "  "); err != nil {
				return utils.ErrMarshalIndent(err)
			}
		}

		// user may pass flag in lower or upper case but we have to keep it lower
		outFormatFlag = strings.ToLower(outFormatFlag)

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return utils.ErrJSONToYAML(err)
			}
		} else if outFormatFlag != "json" {
			return utils.ErrOutFormatFlag()
		}

		// Display the design content
		utils.Log.Info(string(body))

		// Display summary after the full design output if we're looking at a specific design
		if isID || (!viewAllFlag && len(pattern) > 0) {
			utils.Log.Info("\n----- Design Summary -----")
			if designName != "" {
				utils.Log.Info("Design Name: " + designName)
			}

			if sourceType != "" {
				utils.Log.Info("Source Type: " + sourceType)
			}

			// Explain the structure issue with components/relationships
			utils.Log.Info("\nNote: Components and relationships may appear as null at the top level")
			utils.Log.Info("but are present inside the pattern_file structure.")

			// Components and Relationships details sections
			logResourceCount("Components", len(components), sourceType == "manifest")
			if len(components) > 0 {
				displayComponentDetails(components)
			}
			logResourceCount("Relationships", len(relationships), false)
			if len(relationships) > 0 {
				displayRelationshipDetails(relationships)
			}

			// Provide guidance for manifest imports
			if sourceType == "manifest" {
				utils.Log.Info("\nNote: Manifests are imported as raw Kubernetes resources and may not display")
				utils.Log.Info("visual components or relationships. To create a complete design from this manifest,")
				utils.Log.Info("you can use the Meshery UI Design section to add and connect components.")
			}

			utils.Log.Info("-------------------------")
		}

		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}

// logResourceCount logs the count of a resource type or a message if none are found
func logResourceCount(resourceType string, count int, isManifest bool) {
	if count > 0 {
		utils.Log.Info(fmt.Sprintf("\n%s: %d", resourceType, count))
		return
	}

	// No resources found
	if resourceType == "Components" && isManifest {
		utils.Log.Info(fmt.Sprintf("%s: None (Manifest imports may require additional processing)", resourceType))
	} else {
		utils.Log.Info(fmt.Sprintf("%s: None", resourceType))
	}
}

// extractDesignInfo extracts design name and source type from a design map
func extractDesignInfo(design map[string]interface{}) (string, string) {
	designName := ""
	sourceType := ""

	if design == nil {
		return designName, sourceType
	}

	// Extract design name if exists
	if nameVal, exists := design["name"]; exists && nameVal != nil {
		designName = utils.SafeCastToString(nameVal)
	}

	// Extract source type from pattern_file if exists
	if patternFileVal, exists := design["pattern_file"]; exists && patternFileVal != nil {
		// Check if pattern_file is a map
		if patternFile, ok := patternFileVal.(map[string]interface{}); ok && patternFile != nil {
			if sourceTypeVal, exists := patternFile["source_type"]; exists && sourceTypeVal != nil {
				sourceType = utils.SafeCastToString(sourceTypeVal)
			}
		}
	}

	return designName, sourceType
}

// findComponentsAndRelationships extracts components and relationships from a design map
func findComponentsAndRelationships(design map[string]interface{}) ([]interface{}, []interface{}) {
	var components []interface{}
	var relationships []interface{}

	if design == nil {
		return components, relationships
	}

	// Helper function to extract array data from a map if it exists and is non-empty
	extractArray := func(data map[string]interface{}, key string) []interface{} {
		if val, exists := data[key]; exists && val != nil {
			if arr, ok := val.([]interface{}); ok && len(arr) > 0 {
				return arr
			}
		}
		return nil
	}

	// Check top-level components and relationships
	if arr := extractArray(design, "components"); arr != nil {
		components = arr
	}

	if arr := extractArray(design, "relationships"); arr != nil {
		relationships = arr
	}

	// Check for components and relationships within pattern_file
	if patternFileVal, exists := design["pattern_file"]; exists && patternFileVal != nil {
		// Try pattern_file as map
		if patternFile, ok := patternFileVal.(map[string]interface{}); ok && patternFile != nil {
			// Only overwrite components if we haven't found any yet
			if len(components) == 0 {
				if arr := extractArray(patternFile, "components"); arr != nil {
					components = arr
				}
			}

			// Only overwrite relationships if we haven't found any yet
			if len(relationships) == 0 {
				if arr := extractArray(patternFile, "relationships"); arr != nil {
					relationships = arr
				}
			}
		}

		// Try pattern_file as string and parse JSON if we're still missing data
		if len(components) == 0 || len(relationships) == 0 {
			if pfStr, ok := patternFileVal.(string); ok && pfStr != "" {
				var pfObj map[string]interface{}
				if err := json.Unmarshal([]byte(pfStr), &pfObj); err == nil && pfObj != nil {
					// Only extract components if we haven't found any yet
					if len(components) == 0 {
						if arr := extractArray(pfObj, "components"); arr != nil {
							components = arr
						}
					}

					// Only extract relationships if we haven't found any yet
					if len(relationships) == 0 {
						if arr := extractArray(pfObj, "relationships"); arr != nil {
							relationships = arr
						}
					}
				}
			}
		}
	}

	return components, relationships
}

// displayComponentDetails formats and logs component details
func displayComponentDetails(components []interface{}) {
	if len(components) == 0 {
		return
	}

	utils.Log.Info("Component details:")
	for i, compVal := range components {
		comp, ok := compVal.(map[string]interface{})
		if !ok || comp == nil {
			utils.Log.Info(fmt.Sprintf("  %d. Unknown component (invalid format)", i+1))
			continue
		}

		displayName := "Unknown"
		kind := "Unknown"

		if nameVal, exists := comp["displayName"]; exists && nameVal != nil {
			displayName = utils.SafeCastToString(nameVal)
		}

		if compDataVal, exists := comp["component"]; exists && compDataVal != nil {
			if compData, ok := compDataVal.(map[string]interface{}); ok && compData != nil {
				if kindVal, exists := compData["kind"]; exists && kindVal != nil {
					kind = utils.SafeCastToString(kindVal)
				}
			}
		}

		utils.Log.Info(fmt.Sprintf("  %d. %s (Kind: %s)", i+1, displayName, kind))
	}
}

// displayRelationshipDetails formats and logs relationship details
func displayRelationshipDetails(relationships []interface{}) {
	if len(relationships) == 0 {
		return
	}

	utils.Log.Info("Relationship details:")

	// Helper function to safely get strings from map with fallback
	getStringFromMap := func(r map[string]interface{}, key string) string {
		val, exists := r[key]
		if !exists || val == nil {
			return "Unknown"
		}
		return utils.SafeCastToString(val)
	}

	// Helper to log component IDs
	logComponentIDs := func(components []interface{}, prefix string) {
		for _, item := range components {
			component, ok := item.(map[string]interface{})
			if !ok || component == nil {
				continue
			}

			if id := utils.SafeCastToString(component["id"]); id != "" {
				utils.Log.Info(fmt.Sprintf("     %s component ID: %s", prefix, id))
			}
		}
	}

	for i, relVal := range relationships {
		r, ok := relVal.(map[string]interface{})
		if !ok || r == nil {
			utils.Log.Info(fmt.Sprintf("  %d. Unknown relationship (invalid format)", i+1))
			continue
		}

		kind := getStringFromMap(r, "kind")
		subType := getStringFromMap(r, "subType")
		status := getStringFromMap(r, "status")

		utils.Log.Info(fmt.Sprintf("  %d. %s relationship of type %s (Status: %s)",
			i+1, kind, subType, status))

		// Process selectors
		if selectors, ok := r["selectors"].([]interface{}); ok && len(selectors) > 0 {
			for _, selectorVal := range selectors {
				selector, _ := selectorVal.(map[string]interface{})
				if allow, ok := selector["allow"].(map[string]interface{}); ok {

					// Process from components
					if from, ok := allow["from"].([]interface{}); ok {
						logComponentIDs(from, "From")
					}

					// Process to components
					if to, ok := allow["to"].([]interface{}); ok {
						logComponentIDs(to, "To")
					}
				}
			}
		}
	}
}
