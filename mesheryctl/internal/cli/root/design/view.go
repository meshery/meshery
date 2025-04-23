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
	verboseFlag   bool
)

var linkDocPatternView = map[string]string{
	"link":    "![pattern-view-usage](/assets/img/mesheryctl/patternView.png)",
	"caption": "Usage of mesheryctl design view",
}

var viewCmd = &cobra.Command{
	Use:   "view design-name",
	Short: "Display design content",
	Long:  `Display the content of a specific design based on name or id`,
	Args:  cobra.MaximumNArgs(1),
	Example: `
// View a design
mesheryctl design view [design-name | ID]

// View a design with verbosity increased
mesheryctl design view [design-name | ID] --verbose

// View a design displaying all file content
mesheryctl design view --all
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
		// TODO: Currently using "pattern" in API endpoints and responses for backward compatibility.
		// Future versions will align these with "design" terminology
		if len(pattern) == 0 {
			if viewAllFlag {
				url += "/api/pattern?populate=pattern_file&pagesize=10000"
			} else {
				return errors.New(utils.DesignViewError("Please specify a design name or ID. Alternatively, use the `-a` flag to view all available designs."))
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
			designName, sourceType, components, relationships = extractDesignDetails(dat)
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

			designName, sourceType, components, relationships = extractDesignDetails(firstMatch)
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

		// Show summary only with verbose flag
		if verboseFlag && (isID || (!viewAllFlag && len(pattern) > 0)) {
			displayDesignSummary(designName, sourceType, components, relationships)

			if sourceType == "manifest" {
				utils.Log.Info("\nTip: Manifest imports require processing in Meshery UI or Kanvas to display visual components")
			}
		} else if !verboseFlag && (isID || (!viewAllFlag && len(pattern) > 0)) {
			// Add helpful note about verbose flag
			utils.Log.Info("\nTip: Use '--verbose' or '-v' flag to view detailed design summary")
		}

		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "(optional) display detailed design information")
}

// extractDesignDetails extracts design name, source type, components, and relationships
func extractDesignDetails(design map[string]interface{}) (string, string, []interface{}, []interface{}) {
	if design == nil {
		return "", "", nil, nil
	}

	// Extract design name
	designName := utils.SafeCastToString(design["name"])

	// Extract source type from design structure
	sourceType := ""
	if designFileVal, exists := design["pattern_file"]; exists && designFileVal != nil {
		if designFile, ok := designFileVal.(map[string]interface{}); ok {
			sourceType = utils.SafeCastToString(designFile["source_type"])
		}
	}

	// Extract components and relationships using the helper function
	components := utils.ExtractDesignElements(design, "components")
	relationships := utils.ExtractDesignElements(design, "relationships")

	return designName, sourceType, components, relationships
}

// Update displayDesignSummary to handle all logging in one place
func displayDesignSummary(designName, sourceType string, components, relationships []interface{}) {
	if !verboseFlag {
		return
	}

	// Build complete summary
	summaryLines := []string{
		"\n=== Design Summary ===",
		fmt.Sprintf("Name: %s", designName),
	}

	if sourceType != "" {
		summaryLines = append(summaryLines, fmt.Sprintf("Type: %s", sourceType))
	}

	summaryLines = append(summaryLines, "", "Components:")
	if len(components) > 0 {
		for i, comp := range components {
			compMap, ok := comp.(map[string]interface{})
			if !ok {
				continue
			}

			name := utils.SafeCastToString(compMap["displayName"])
			if name == "" {
				name = "Unnamed"
			}

			kind := "Unknown"
			if compData, ok := compMap["component"].(map[string]interface{}); ok {
				kind = utils.SafeCastToString(compData["kind"])
			}

			summaryLines = append(summaryLines, fmt.Sprintf("  %d. %s (%s)", i+1, name, kind))
		}
	} else {
		summaryLines = append(summaryLines, "  No components found in this design")
	}

	// Relationships section
	summaryLines = append(summaryLines, "", "Relationships:")
	if len(relationships) > 0 {
		for i, rel := range relationships {
			relMap, ok := rel.(map[string]interface{})
			if !ok {
				continue
			}

			kind := utils.SafeCastToString(relMap["kind"])
			subType := utils.SafeCastToString(relMap["subType"])
			status := utils.SafeCastToString(relMap["status"])

			summaryLines = append(summaryLines, fmt.Sprintf("  %d. %s %s (%s)",
				i+1, kind, subType, status))
		}
	} else {
		summaryLines = append(summaryLines, "  No relationships found in this design")
	}

	// Add notes section if needed
	if sourceType == "manifest" {
		summaryLines = append(summaryLines, "", "Note: Manifest imports require processing in Meshery UI or Kanvas to display visual components")
	}

	if len(components) > 0 || len(relationships) > 0 {
		summaryLines = append(summaryLines, "", "Note: Design elements are nested in the structure")
	}

	summaryLines = append(summaryLines, "\n===========================")

	// Display entire summary at once
	utils.Log.Info(strings.Join(summaryLines, "\n"))
}
