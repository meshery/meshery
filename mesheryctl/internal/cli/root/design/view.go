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
	"fmt"
	"net/url"

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

		urlPath := getDesignViewUrlPath(design, isID, designViewFlagsProvided.All)

		apiResponse, err := api.Fetch[map[string]interface{}](urlPath)
		if err != nil {
			return err
		}

		var designData interface{}

		if isID { // Only one design will be returned when searching by ID
			designData = apiResponse
		}

		if designViewFlagsProvided.All { // all designs will be returned when --all flag is provided
			designData = map[string]interface{}{"patterns": (*apiResponse)["patterns"]}
		}

		if !isID && !designViewFlagsProvided.All { // only the first match will be returned when searching by name
			patterns, ok := (*apiResponse)["patterns"]
			if !ok {
				return ErrDesignInvalidApiResponse("'patterns' field missing")
			}
			arr, ok := patterns.([]interface{})
			if !ok {
				return ErrDesignInvalidApiResponse("'patterns' field is not of expected type")
			}
			if len(arr) == 0 {
				return ErrDesignNotFound()
			}
			designData = arr[0]
		}

		outputFormatFactory := display.OutputFormatterFactory[interface{}]{}
		outputFormatter, err := outputFormatFactory.New(designViewFlagsProvided.OutputFormat, designData)
		if err != nil {
			return err
		}

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

func init() {
	viewCmd.Flags().BoolVarP(&designViewFlagsProvided.All, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&designViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
