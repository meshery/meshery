// Copyright 2023 Layer5, Inc.
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

package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

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

var viewCmd = &cobra.Command{
	Use:   "view [filter-name | ID]",
	Short: "View filter(s)",
	Long:  `Displays the contents of a specific filter based on name or id`,
	Example: `
// View the specified WASM filter
// A unique prefix of the name or ID can also be provided. If the prefix is not unique, the first match will be returned.
mesheryctl filter view [filter-name | ID]	

// View all filter files
mesheryctl filter view --all
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		var baseUrl = mctlCfg.GetBaseMesheryURL()
		filter := ""
		isID := false
		// if filter name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New(utils.FilterViewError("--all cannot be used when filter name or ID is specified\nUse 'mesheryctl filter view --help' to display usage guide\n"))
			}
			filter, isID, err = utils.ValidId(baseUrl, args[0], "filter")
			if err != nil {
				return errors.New("invalid filter name or ID. " + err.Error())
			}
		}

		if len(filter) == 0 {
			if viewAllFlag {
				baseUrl += "/api/filter?pagesize=10000"
			} else {
				return errors.New(utils.FilterViewError("filter-name or ID not specified, use -a to view all filters\nUse 'mesheryctl filter view --help' to display usage guide\n"))
			}
		} else if isID {
			// if filter is a valid uuid, then directly fetch the filter
			baseUrl += "/api/filter/" + filter
		} else {
			// else search filter by name
			baseUrl += "/api/filter?search=" + url.QueryEscape(filter)
		}

		req, err := utils.NewRequest("GET", baseUrl, nil)
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
			return err
		}

		var dat map[string]interface{}
		if err = json.Unmarshal(body, &dat); err != nil {
			return ErrUnmarshal(err)
		}

		if isID {
			if body, err = json.MarshalIndent(dat, "", "  "); err != nil {
				return err
			}
		} else if viewAllFlag {
			// only keep the filter key from the response when viewing all the filters
			if body, err = json.MarshalIndent(map[string]interface{}{"filters": dat["filters"]}, "", "  "); err != nil {
				return err
			}
		} else {
			// use the first match from the result when searching by filter name
			arr := dat["filters"].([]interface{})
			if len(arr) == 0 {
				utils.Log.Info(fmt.Sprintf("filter with name: %s not found", filter))
				return nil
			}
			if body, err = json.MarshalIndent(arr[0], "", "  "); err != nil {
				return err
			}
		}

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return errors.Wrap(err, "failed to convert json to yaml")
			}
		} else if outFormatFlag != "json" {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}
		utils.Log.Info(string(body))
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all filters available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
