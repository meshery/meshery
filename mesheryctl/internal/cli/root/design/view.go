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
		hasComponents := false
		componentsCount := 0
		var components []interface{}

		hasRelationships := false
		relationshipsCount := 0
		var relationships []interface{}

		designName := ""
		sourceType := ""

		if isID {
			// Check components and relationships for summary
			if _, exists := dat["name"]; exists && dat["name"] != nil {
				designName = fmt.Sprintf("%v", dat["name"])
			}

			// Check for source type
			if _, exists := dat["pattern_file"]; exists && dat["pattern_file"] != nil {
				patternFile, ok := dat["pattern_file"].(map[string]interface{})
				if ok && patternFile["source_type"] != nil {
					sourceType = fmt.Sprintf("%v", patternFile["source_type"])
				}

				// Also check for components and relationships within pattern_file
				if pfComponents, ok := patternFile["components"].([]interface{}); ok && len(pfComponents) > 0 {
					hasComponents = true
					componentsCount = len(pfComponents)
					components = pfComponents
				}

				if pfRelationships, ok := patternFile["relationships"].([]interface{}); ok && len(pfRelationships) > 0 {
					hasRelationships = true
					relationshipsCount = len(pfRelationships)
					relationships = pfRelationships
				}

				// If pattern_file is a string, try to parse it for components/relationships
				if pf, ok := dat["pattern_file"].(string); ok {
					var pfObj map[string]interface{}
					if err := json.Unmarshal([]byte(pf), &pfObj); err == nil {
						if pfComponents, ok := pfObj["components"].([]interface{}); ok && len(pfComponents) > 0 {
							hasComponents = true
							componentsCount = len(pfComponents)
							components = pfComponents
						}

						if pfRelationships, ok := pfObj["relationships"].([]interface{}); ok && len(pfRelationships) > 0 {
							hasRelationships = true
							relationshipsCount = len(pfRelationships)
							relationships = pfRelationships
						}
					}
				}
			}

			// Check top-level components
			if _, exists := dat["components"]; exists && dat["components"] != nil {
				if comps, ok := dat["components"].([]interface{}); ok && len(comps) > 0 {
					hasComponents = true
					componentsCount = len(comps)
					components = comps
				}
			}

			// Check top-level relationships
			if _, exists := dat["relationships"]; exists && dat["relationships"] != nil {
				if rels, ok := dat["relationships"].([]interface{}); ok && len(rels) > 0 {
					hasRelationships = true
					relationshipsCount = len(rels)
					relationships = rels
				}
			}

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
			arr := dat["patterns"].([]interface{})
			if len(arr) == 0 {
				utils.Log.Error(ErrDesignNotFound())
				return nil
			}

			// Extract information for summary
			firstMatch := arr[0].(map[string]interface{})
			if _, exists := firstMatch["name"]; exists && firstMatch["name"] != nil {
				designName = fmt.Sprintf("%v", firstMatch["name"])
			}

			// Check for source type and potentially components/relationships in pattern_file
			if _, exists := firstMatch["pattern_file"]; exists && firstMatch["pattern_file"] != nil {
				patternFile, ok := firstMatch["pattern_file"].(map[string]interface{})
				if ok {
					if patternFile["source_type"] != nil {
						sourceType = fmt.Sprintf("%v", patternFile["source_type"])
					}

					// Check for components and relationships within pattern_file
					if pfComponents, ok := patternFile["components"].([]interface{}); ok && len(pfComponents) > 0 {
						hasComponents = true
						componentsCount = len(pfComponents)
						components = pfComponents
					}

					if pfRelationships, ok := patternFile["relationships"].([]interface{}); ok && len(pfRelationships) > 0 {
						hasRelationships = true
						relationshipsCount = len(pfRelationships)
						relationships = pfRelationships
					}
				}

				// If pattern_file is a string, try to parse it for components/relationships
				if pf, ok := firstMatch["pattern_file"].(string); ok {
					var pfObj map[string]interface{}
					if err := json.Unmarshal([]byte(pf), &pfObj); err == nil {
						if pfComponents, ok := pfObj["components"].([]interface{}); ok && len(pfComponents) > 0 {
							hasComponents = true
							componentsCount = len(pfComponents)
							components = pfComponents
						}

						if pfRelationships, ok := pfObj["relationships"].([]interface{}); ok && len(pfRelationships) > 0 {
							hasRelationships = true
							relationshipsCount = len(pfRelationships)
							relationships = pfRelationships
						}
					}
				}
			}

			// Check top-level components
			if _, exists := firstMatch["components"]; exists && firstMatch["components"] != nil {
				if comps, ok := firstMatch["components"].([]interface{}); ok && len(comps) > 0 {
					hasComponents = true
					componentsCount = len(comps)
					components = comps
				}
			}

			// Check top-level relationships
			if _, exists := firstMatch["relationships"]; exists && firstMatch["relationships"] != nil {
				if rels, ok := firstMatch["relationships"].([]interface{}); ok && len(rels) > 0 {
					hasRelationships = true
					relationshipsCount = len(rels)
					relationships = rels
				}
			}

			if !hasComponents {
				utils.Log.Info("No components found for this design.")
			}

			if !hasRelationships {
				utils.Log.Info("No relationships found for this design.")
			}

			if body, err = json.MarshalIndent(arr[0], "", "  "); err != nil {
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
			if hasComponents {
				utils.Log.Info("\nComponents: " + strconv.Itoa(componentsCount))

				if len(components) > 0 {
					utils.Log.Info("Component details:")
					for i, comp := range components {
						c := comp.(map[string]interface{})
						displayName := "Unknown"
						kind := "Unknown"

						if name, ok := c["displayName"]; ok && name != nil {
							displayName = fmt.Sprintf("%v", name)
						}

						if compData, ok := c["component"].(map[string]interface{}); ok && compData != nil {
							if k, ok := compData["kind"]; ok && k != nil {
								kind = fmt.Sprintf("%v", k)
							}
						}

						utils.Log.Info(fmt.Sprintf("  %d. %s (Kind: %s)", i+1, displayName, kind))
					}
				}
			} else {
				if sourceType == "manifest" {
					utils.Log.Info("Components: None (Manifest imports may require additional processing)")
				} else {
					utils.Log.Info("Components: None")
				}
			}

			// Relationships details section
			if hasRelationships {
				utils.Log.Info("\nRelationships: " + strconv.Itoa(relationshipsCount))

				if len(relationships) > 0 {
					utils.Log.Info("Relationship details:")
					for i, rel := range relationships {
						r := rel.(map[string]interface{})
						kind := "Unknown"
						subType := "Unknown"
						status := "Unknown"

						if k, ok := r["kind"]; ok && k != nil {
							kind = fmt.Sprintf("%v", k)
						}

						if st, ok := r["subType"]; ok && st != nil {
							subType = fmt.Sprintf("%v", st)
						}

						if s, ok := r["status"]; ok && s != nil {
							status = fmt.Sprintf("%v", s)
						}

						utils.Log.Info(fmt.Sprintf("  %d. %s relationship of type %s (Status: %s)", i+1, kind, subType, status))

						// Try to extract from/to components if available
						if selectors, ok := r["selectors"].([]interface{}); ok && len(selectors) > 0 {
							for _, selector := range selectors {
								s := selector.(map[string]interface{})
								if allow, ok := s["allow"].(map[string]interface{}); ok {
									if from, ok := allow["from"].([]interface{}); ok && len(from) > 0 {
										for _, f := range from {
											fmap := f.(map[string]interface{})
											if id, ok := fmap["id"].(string); ok {
												utils.Log.Info(fmt.Sprintf("     From component ID: %s", id))
											}
										}
									}

									if to, ok := allow["to"].([]interface{}); ok && len(to) > 0 {
										for _, t := range to {
											tmap := t.(map[string]interface{})
											if id, ok := tmap["id"].(string); ok {
												utils.Log.Info(fmt.Sprintf("     To component ID: %s", id))
											}
										}
									}
								}
							}
						}
					}
				}
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
