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
	"io"
	neturl "net/url"
	"slices"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag      bool
	outputFormatFlag string
)

var validOutputFormat = []string{"json", "yaml"}

var linkDocPatternView = map[string]string{
	"link":    "![pattern-view-usage](/assets/img/mesheryctl/patternView.png)",
	"caption": "Usage of mesheryctl design view",
}

var viewCmd = &cobra.Command{
	Use:   "view design name",
	Short: "Display a design content",
	Long:  `Display the content of a specific design based on name or id`,
	Args: func(cmd *cobra.Command, args []string) error {
		// Validate output-format
		if !slices.Contains(validOutputFormat, strings.ToLower(outputFormatFlag)) {
			return utils.ErrInvalidArgument(errors.New("output-format choice is invalid, use [json|yaml]"))
		}
		return nil
	},

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
				return utils.ErrInvalidNameOrID(err)
			}
		}

		url := mctlCfg.GetBaseMesheryURL()
		if len(pattern) == 0 {
			if viewAllFlag {
				url += "/api/pattern?populate=pattern_file&pagesize=10000"
			} else {
				return ErrDesignNameOrIDNotSpecified()
			}
		} else if isID {
			// if pattern is a valid uuid, then directly fetch the pattern
			url += "/api/pattern/" + pattern
		} else {
			// else search pattern by name
			escapedPattern := neturl.QueryEscape(pattern)
			url += "/api/pattern?populate=pattern_file&search=" + escapedPattern
		}

		req, err := utils.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}

		defer func() { _ = res.Body.Close() }()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		var dat map[string]interface{}
		if err = json.Unmarshal(body, &dat); err != nil {
			return utils.ErrUnmarshal(err)
		}

		var output any

		if isID {
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
			arr := dat["patterns"].([]interface{})
			if len(arr) == 0 {
				return ErrDesignNotFound()
			}
			if body, err = json.MarshalIndent(arr[0], "", "  "); err != nil {
				return utils.ErrMarshalIndent(err)
			}
		}

		if outputFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return utils.ErrJSONToYAML(err)
			}
		} else if outputFormatFlag != "json" {
			return utils.ErrOutFormatFlag()
		}
		utils.Log.Info(string(body))
		return nil

		outputFormatterFactory := display.OutputFormatterFactory[any]{}
		outputFormatter, err := outputFormatterFactory.New(outputFormatFlag, output)

		if err != nil {
			return err
		}

		return outputFormatter.Display()
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all designs available")
	viewCmd.Flags().StringVarP(&outputFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
