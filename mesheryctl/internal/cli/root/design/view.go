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
	"net/url"
	"strconv"
	"strings"

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

		switch {
		case isID:
			urlPath := getDesignViewUrlPath(design, true, false)
			apiResponse, err := api.Fetch[map[string]interface{}](urlPath)
			if err != nil {
				return err
			}
			designData = apiResponse

		case designViewFlagsProvided.All:
			urlPath := getDesignViewUrlPath("", false, true)
			apiResponse, err := api.Fetch[map[string]interface{}](urlPath)
			if err != nil {
				return err
			}
			designData = map[string]interface{}{"patterns": (*apiResponse)["patterns"]}

		default:
			allPatterns, err := fetchDesignSearch(design, api.Fetch[map[string]interface{}])
			if err != nil {
				return err
			}
			selected, err := selectDesignByName(allPatterns, design)
			if err != nil {
				return err
			}
			designData = selected
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

func getDesignViewUrlPath(design string, isID bool, viewAll bool) string {
	baseApiPath := "api/pattern"
	if isID {
		return fmt.Sprintf("%s/%s", baseApiPath, url.PathEscape(design))
	}

	queryParams := url.Values{}
	queryParams.Add("populate", "pattern_file")

	if viewAll {
		queryParams.Add("pagesize", "10000")
	}

	if !viewAll {
		queryParams.Add("search", design)
	}

	return fmt.Sprintf("%s?%s", baseApiPath, queryParams.Encode())
}

const designSearchPageSize = 10000

var runDesignSelectPrompt = utils.RunSelectPrompt

// extractInt extracts an integer from a map, handling JSON numbers that unmarshal as float64, int, int64, or json.Number.
func extractInt(m map[string]interface{}, keys ...string) (int, bool) {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				return int(val), true
			case int:
				return val, true
			case int64:
				return int(val), true
			case json.Number:
				if n, err := val.Int64(); err == nil {
					return int(n), true
				}
			}
		}
	}
	return 0, false
}

// fetchDesignSearch collects every page so that selection can consider exact
// matches and duplicates even when the provider caps the requested page size.
func fetchDesignSearch(name string, fetch func(string) (*map[string]interface{}, error)) ([]interface{}, error) {
	var allPatterns []interface{}
	var initialPageSize int
	hasTotalCount := false
	totalCount := 0

	for pageNum := 0; ; pageNum++ {
		values := url.Values{}
		values.Add("page", strconv.Itoa(pageNum))
		values.Add("pagesize", strconv.Itoa(designSearchPageSize))
		values.Add("populate", "pattern_file")
		values.Add("search", name)

		query := fmt.Sprintf("api/pattern?%s", values.Encode())
		page, err := fetch(query)
		if err != nil {
			return nil, err
		}
		if page == nil {
			return nil, ErrDesignInvalidApiResponse("empty design search response")
		}

		rawPatterns, ok := (*page)["patterns"]
		if !ok {
			return nil, ErrDesignInvalidApiResponse("'patterns' field missing")
		}
		patternsArr, ok := rawPatterns.([]interface{})
		if !ok {
			return nil, ErrDesignInvalidApiResponse("'patterns' field is not of expected type")
		}

		if tot, ok := extractInt(*page, "totalCount", "total_count"); ok {
			hasTotalCount = true
			totalCount = tot
		}
		curPageSize, hasPageSize := extractInt(*page, "pageSize", "page_size")

		if pageNum == 0 {
			if hasTotalCount && totalCount == 0 {
				return []interface{}{}, nil
			}
			if hasPageSize {
				initialPageSize = curPageSize
			}
		} else if hasPageSize && curPageSize != initialPageSize {
			return nil, ErrDesignInvalidApiResponse("design search page size changed")
		}

		allPatterns = append(allPatterns, patternsArr...)

		if hasTotalCount && len(allPatterns) >= totalCount {
			return allPatterns, nil
		}
		if len(patternsArr) == 0 || (hasPageSize && curPageSize == 0) {
			if hasTotalCount && len(allPatterns) < totalCount {
				return nil, ErrDesignInvalidApiResponse("incomplete design search results")
			}
			return allPatterns, nil
		}
	}
}

// selectDesignByName prefers case-sensitive whole-name matches, then falls
// back to case-insensitive whole-name matches. Substring results are ignored.
func selectDesignByName(patterns []interface{}, name string) (map[string]interface{}, error) {
	var exact, folded []map[string]interface{}
	for _, p := range patterns {
		m, ok := p.(map[string]interface{})
		if !ok || m == nil {
			continue
		}
		patternName, ok := m["name"].(string)
		if !ok {
			continue
		}
		if patternName == name {
			exact = append(exact, m)
		} else if strings.EqualFold(patternName, name) {
			folded = append(folded, m)
		}
	}

	matches := exact
	if len(matches) == 0 {
		matches = folded
	}
	if len(matches) == 0 {
		return nil, ErrDesignNotFound(name)
	}
	if len(matches) == 1 {
		return matches[0], nil
	}

	items := make([]string, len(matches))
	for i, m := range matches {
		pName, _ := m["name"].(string)
		pID, _ := m["id"].(string)
		if pID == "" {
			pID = "unknown"
		}
		items[i] = fmt.Sprintf("Name: %s, ID: %s", pName, pID)
	}

	idx, err := runDesignSelectPrompt("Select design", items)
	if err != nil {
		return nil, err
	}
	if idx < 0 || idx >= len(matches) {
		return nil, fmt.Errorf("invalid design selection: %d", idx)
	}
	return matches[idx], nil
}

func init() {
	viewCmd.Flags().BoolVarP(&designViewFlagsProvided.All, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&designViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
