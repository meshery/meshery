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
	"math"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type designViewFlags struct {
	OutputFormat string
	All          bool
}

var designViewFlagsProvided designViewFlags

var linkDocPatternView = map[string]string{
	"link":    "![pattern-view-usage](../../../images/patternView.png)",
	"caption": "Usage of mesheryctl design view",
}

var viewCmd = &cobra.Command{
	Use:   "view [design-name | ID]",
	Short: "Display design content",
	Long: `Display the content of a specific design based on name or id.
	Find more information at: https://docs.meshery.io/reference/mesheryctl/design/view`,
	Args: cobra.MaximumNArgs(1),
	Example: `
// view a design
mesheryctl design view [design-name | ID]
	`,
	Annotations: linkDocPatternView,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 && designViewFlagsProvided.All {
			return utils.ErrInvalidArgument(errors.New("--all flag cannot be used when [design-name|design-id] is specified"))
		}

		if len(args) == 0 && !designViewFlagsProvided.All {
			return ErrDesignNameOrIDNotSpecified()
		}

		// Reject the format before any request, so a name shared by several
		// designs does not prompt for a design that cannot be displayed.
		if _, err := (&display.OutputFormatterFactory[interface{}]{}).New(designViewFlagsProvided.OutputFormat, nil); err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		design := ""
		isID := false

		if len(args) > 0 { // if design name or ID is provided
			design, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), args[0], "pattern")
			if err != nil {
				return utils.ErrInvalidNameOrID(err)
			}
		}

		var designData interface{}

		if isID || designViewFlagsProvided.All {
			apiResponse, err := api.Fetch[map[string]interface{}](getDesignViewUrlPath(design, isID))
			if err != nil {
				return err
			}
			if isID { // Only one design will be returned when searching by ID
				designData = apiResponse
			} else { // all designs will be returned when --all flag is provided
				designData = map[string]interface{}{"patterns": (*apiResponse)["patterns"]}
			}
		} else {
			patterns, err := fetchDesignSearch(design, api.Fetch[map[string]interface{}])
			if err != nil {
				return err
			}
			designData, err = selectDesignByName(patterns, design)
			if err != nil {
				return err
			}
		}

		outputFormatFactory := display.OutputFormatterFactory[interface{}]{}
		outputFormatter, err := outputFormatFactory.New(designViewFlagsProvided.OutputFormat, designData)
		if err != nil {
			return err
		}
		outputFormatter = outputFormatter.WithOutput(cmd.OutOrStdout())

		return outputFormatter.Display()
	},
}

func getDesignViewUrlPath(design string, isID bool) string {
	baseApiPath := "api/pattern"
	if isID {
		return fmt.Sprintf("%s/%s", baseApiPath, url.PathEscape(design))
	}

	queryParams := url.Values{}
	queryParams.Add("populate", "pattern_file")
	queryParams.Add("pagesize", "10000")

	return fmt.Sprintf("%s?%s", baseApiPath, queryParams.Encode())
}

// fetchDesignSearch collects every page of a name search. The server's search
// is a substring match and pages its results, so the design carrying the whole
// name can be on any page. Designs are kept as returned so that the output of
// design view is unchanged.
func fetchDesignSearch(name string, fetch func(string) (*map[string]interface{}, error)) ([]map[string]interface{}, error) {
	var designs []map[string]interface{}
	seen := map[string]bool{}
	for page := 0; ; page++ {
		queryParams := url.Values{}
		queryParams.Add("populate", "pattern_file")
		queryParams.Add("search", name)
		queryParams.Add("pagesize", "10000")
		queryParams.Add("page", strconv.Itoa(page))

		response, err := fetch("api/pattern?" + queryParams.Encode())
		if err != nil {
			return nil, err
		}
		if response == nil {
			return nil, ErrDesignInvalidApiResponse("empty response")
		}
		patterns, ok := (*response)["patterns"]
		if !ok {
			return nil, ErrDesignInvalidApiResponse("'patterns' field missing")
		}
		arr, ok := patterns.([]interface{})
		if !ok {
			return nil, ErrDesignInvalidApiResponse("'patterns' field is not of expected type")
		}

		added := 0
		for _, item := range arr {
			design, ok := item.(map[string]interface{})
			if !ok {
				return nil, ErrDesignInvalidApiResponse("design entry is not an object")
			}
			if id, _ := design["id"].(string); id != "" {
				if seen[id] {
					continue
				}
				seen[id] = true
			}
			designs = append(designs, design)
			added++
		}

		total, ok, err := designSearchTotal(*response)
		if err != nil {
			return nil, err
		}
		// Without a total the response cannot say whether more pages exist.
		if !ok || len(designs) >= total {
			return designs, nil
		}
		// Each further page must add designs, or the loop would never end.
		if added == 0 {
			return nil, ErrDesignInvalidApiResponse(fmt.Sprintf("search for %q returned %d of %d designs", name, len(designs), total))
		}
	}
}

// designSearchTotal reads the number of designs matching a search, from the
// current "totalCount" field or the legacy "total_count" field. A total that
// is present but not a whole number is an error: treating it as missing would
// stop after the first page and could miss the design with the name.
func designSearchTotal(response map[string]interface{}) (int, bool, error) {
	for _, key := range []string{"totalCount", "total_count"} {
		value, ok := response[key]
		if !ok || value == nil {
			continue
		}
		switch n := value.(type) {
		case float64:
			if n == math.Trunc(n) {
				return int(n), true, nil
			}
		case int:
			return n, true, nil
		case int64:
			return int(n), true, nil
		case json.Number:
			if i, err := n.Int64(); err == nil {
				return int(i), true, nil
			}
		}
		return 0, false, ErrDesignInvalidApiResponse(fmt.Sprintf("'%s' field is not a whole number: %v", key, value))
	}
	return 0, false, nil
}

// runDesignSelectPrompt draws the prompt on stderr: stdout carries the design,
// and `design view <name> -o json > design.json` must stay valid JSON.
var runDesignSelectPrompt = func(label string, items []string) (int, error) {
	prompt := promptui.Select{Label: label, Items: items, Stdout: os.Stderr}
	i, _, err := prompt.Run()
	if err != nil {
		return 0, utils.ErrPromptCancelled()
	}
	return i, nil
}

// selectDesignByName prefers case-sensitive whole-name matches, then falls
// back to case-insensitive whole-name matches. Designs whose name only
// contains name are ignored. Several designs with the name are offered for
// selection instead of picking one.
func selectDesignByName(designs []map[string]interface{}, name string) (map[string]interface{}, error) {
	var exact, folded []map[string]interface{}
	for _, design := range designs {
		designName, _ := design["name"].(string)
		if designName == name {
			exact = append(exact, design)
		} else if strings.EqualFold(designName, name) {
			folded = append(folded, design)
		}
	}
	matches := exact
	if len(matches) == 0 {
		matches = folded
	}
	switch len(matches) {
	case 0:
		return nil, ErrDesignNotFound(name)
	case 1:
		return matches[0], nil
	}
	items := make([]string, len(matches))
	for i, design := range matches {
		items[i] = fmt.Sprintf("Name: %v, ID: %v", design["name"], design["id"])
	}
	i, err := runDesignSelectPrompt("Select design", items)
	if err != nil {
		return nil, err
	}
	if i < 0 || i >= len(matches) {
		return nil, fmt.Errorf("invalid design selection: %d", i)
	}
	return matches[i], nil
}

func init() {
	viewCmd.Flags().BoolVarP(&designViewFlagsProvided.All, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&designViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
