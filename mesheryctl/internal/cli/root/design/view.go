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
	"strconv"
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
			utils.Log.Error(err)
			return nil
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
				utils.Log.Error(ErrPatternInvalidNameOrID(err))
				return nil
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
			utils.Log.Error(err)
			return nil
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			utils.Log.Error(utils.ErrReadResponseBody(err))
			return nil
		}

		var dat map[string]interface{}
		if err = json.Unmarshal(body, &dat); err != nil {
			utils.Log.Error(utils.ErrUnmarshal(err))
			return nil
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
				utils.Log.Error(utils.ErrMarshalIndent(err))
				return nil
			}
		} else if viewAllFlag {
			// only keep the pattern key from the response when viewing all the patterns
			if body, err = json.MarshalIndent(map[string]interface{}{"patterns": dat["patterns"]}, "", "  "); err != nil {
				utils.Log.Error(utils.ErrMarshalIndent(err))
				return nil
			}
		} else {
			// use the first match from the result when searching by pattern name
			patternsVal, exists := dat["patterns"]
			if !exists || patternsVal == nil {
				utils.Log.Error(ErrDesignNotFound())
				return nil
			}

			patterns, ok := patternsVal.([]interface{})
			if !ok || len(patterns) == 0 {
				utils.Log.Error(ErrDesignNotFound())
				return nil
			}

			firstMatchVal := patterns[0]
			firstMatch, ok := firstMatchVal.(map[string]interface{})
			if !ok || firstMatch == nil {
				utils.Log.Error(ErrDesignNotFound())
				return nil
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
				utils.Log.Error(utils.ErrMarshalIndent(err))
				return nil
			}
		}

		// user may pass flag in lower or upper case but we have to keep it lower
		outFormatFlag = strings.ToLower(outFormatFlag)

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				utils.Log.Error(utils.ErrJSONToYAML(err))
				return nil
			}
		} else if outFormatFlag != "json" {
			utils.Log.Error(utils.ErrOutFormatFlag())
			return nil
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

			// Components details section
			hasComponents := len(components) > 0
			componentsCount := len(components)

			if hasComponents {
				utils.Log.Info("\nComponents: " + strconv.Itoa(componentsCount))
				displayComponentDetails(components)
			} else {
				if sourceType == "manifest" {
					utils.Log.Info("Components: None (Manifest imports may require additional processing)")
				} else {
					utils.Log.Info("Components: None")
				}
			}

			// Relationships details section
			hasRelationships := len(relationships) > 0
			relationshipsCount := len(relationships)

			if hasRelationships {
				utils.Log.Info("\nRelationships: " + strconv.Itoa(relationshipsCount))
				displayRelationshipDetails(relationships)
			} else {
				if sourceType == "manifest" {
					utils.Log.Info("Relationships: None (Manifest imports typically don't include relationship data)")
				} else {
					utils.Log.Info("Relationships: None")
				}
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

// extractDesignInfo extracts design name and source type from a design map
func extractDesignInfo(design map[string]interface{}) (string, string) {
	designName := ""
	sourceType := ""

	if design == nil {
		return designName, sourceType
	}

	// Extract design name if exists
	if nameVal, exists := design["name"]; exists && nameVal != nil {
		if name, ok := nameVal.(string); ok {
			designName = name
		} else {
			designName = fmt.Sprintf("%v", nameVal)
		}
	}

	// Extract source type from pattern_file if exists
	if patternFileVal, exists := design["pattern_file"]; exists && patternFileVal != nil {
		// Check if pattern_file is a map
		if patternFile, ok := patternFileVal.(map[string]interface{}); ok && patternFile != nil {
			if sourceTypeVal, exists := patternFile["source_type"]; exists && sourceTypeVal != nil {
				if st, ok := sourceTypeVal.(string); ok {
					sourceType = st
				} else {
					sourceType = fmt.Sprintf("%v", sourceTypeVal)
				}
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

	// Check top-level components
	if compsVal, exists := design["components"]; exists && compsVal != nil {
		if comps, ok := compsVal.([]interface{}); ok && len(comps) > 0 {
			components = comps
		}
	}

	// Check top-level relationships
	if relsVal, exists := design["relationships"]; exists && relsVal != nil {
		if rels, ok := relsVal.([]interface{}); ok && len(rels) > 0 {
			relationships = rels
		}
	}

	// Check for components and relationships within pattern_file
	if patternFileVal, exists := design["pattern_file"]; exists && patternFileVal != nil {
		// Try pattern_file as map
		if patternFile, ok := patternFileVal.(map[string]interface{}); ok && patternFile != nil {
			// Check for components in pattern_file map
			if pfCompsVal, exists := patternFile["components"]; exists && pfCompsVal != nil {
				if pfComponents, ok := pfCompsVal.([]interface{}); ok && len(pfComponents) > 0 {
					components = pfComponents
				}
			}

			// Check for relationships in pattern_file map
			if pfRelsVal, exists := patternFile["relationships"]; exists && pfRelsVal != nil {
				if pfRelationships, ok := pfRelsVal.([]interface{}); ok && len(pfRelationships) > 0 {
					relationships = pfRelationships
				}
			}
		}

		// Try pattern_file as string and parse JSON
		if len(components) == 0 || len(relationships) == 0 {
			if pfStr, ok := patternFileVal.(string); ok && pfStr != "" {
				var pfObj map[string]interface{}
				if err := json.Unmarshal([]byte(pfStr), &pfObj); err == nil && pfObj != nil {
					// Check for components in parsed pattern_file
					if len(components) == 0 {
						if pfCompsVal, exists := pfObj["components"]; exists && pfCompsVal != nil {
							if pfComponents, ok := pfCompsVal.([]interface{}); ok && len(pfComponents) > 0 {
								components = pfComponents
							}
						}
					}

					// Check for relationships in parsed pattern_file
					if len(relationships) == 0 {
						if pfRelsVal, exists := pfObj["relationships"]; exists && pfRelsVal != nil {
							if pfRelationships, ok := pfRelsVal.([]interface{}); ok && len(pfRelationships) > 0 {
								relationships = pfRelationships
							}
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
	if components == nil || len(components) == 0 {
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
			if name, ok := nameVal.(string); ok {
				displayName = name
			} else {
				displayName = fmt.Sprintf("%v", nameVal)
			}
		}

		if compDataVal, exists := comp["component"]; exists && compDataVal != nil {
			if compData, ok := compDataVal.(map[string]interface{}); ok && compData != nil {
				if kindVal, exists := compData["kind"]; exists && kindVal != nil {
					if k, ok := kindVal.(string); ok {
						kind = k
					} else {
						kind = fmt.Sprintf("%v", kindVal)
					}
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
