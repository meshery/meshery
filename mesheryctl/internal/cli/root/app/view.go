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

package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag   bool
	outFormatFlag string
)

var linkDocAppView = map[string]string{
	"link":    "![app-view-usage](/assets/img/mesheryctl/app-view.png)",
	"caption": "Usage of mesheryctl app view",
}

var viewCmd = &cobra.Command{
	Use:   "view application name",
	Short: "Display application(s)",
	Long:  `Displays the contents of a specific application based on name or id`,
	Example: `
// View applictaions with name
mesheryctl app view [app-name]

// View applications with id
mesheryctl app view [app-id]

// View all applications
mesheryctl app view --all
	`,
	Annotations: linkDocAppView,
	Args:        cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		var baseUrl = mctlCfg.GetBaseMesheryURL()
		application := ""
		isID := false
		applicationID := ""
		// if application name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New("-a cannot be used when [application-name|application-id] is specified")
			}
			applicationID, isID, err = utils.ValidId(baseUrl, args[0], "application")
			if err != nil {
				return err
			}
		}
		var req *http.Request
		// url := mctlCfg.GetBaseMesheryURL()
		var response *models.ApplicationsAPIResponse
		// Merge args to get app-name
		application = strings.Join(args, "%20")
		if len(application) == 0 {
			if viewAllFlag {
				baseUrl += "/api/application?pagesize=10000"
			} else {
				return errors.New("[application-name|application-id] not specified, use -a to view all applications")
			}
		} else if isID {
			// if application is a valid uuid, then directly fetch the application
			baseUrl += "/api/application/" + applicationID
		} else {
			// else search application by name
			baseUrl += "/api/application?search=" + application
		}

		req, err = utils.NewRequest("GET", baseUrl, nil)
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
			return errors.Wrap(err, "failed to unmarshal response body")
		}
		if isID {
			if body, err = json.MarshalIndent(dat, "", "  "); err != nil {
				return err
			}
		} else if viewAllFlag {
			if body, err = json.MarshalIndent(map[string]interface{}{"applications": dat["applications"]}, "", "  "); err != nil {
				return err
			}
		} else {
			if err = json.Unmarshal(body, &response); err != nil {
				return errors.Wrap(err, "failed to unmarshal response body")
			}
			if response.TotalCount == 0 {
				return errors.New("application does not exit. Please get an app name and try again. Use `mesheryctl app list` to see a list of applications")
			}
			// Manage more than one apps with similar name
			for _, app := range response.Applications {
				if response.Applications == nil {
					return errors.New("application name not provide. Please get an app name and try again. Use `mesheryctl app list` to see a list of applications")
				}
				body, err = json.MarshalIndent(&app, "", "  ")
				if err != nil {
					return err
				}
				if outFormatFlag == "json" {
					utils.Log.Info(string(body))
					continue
				}
				if outFormatFlag == "yaml" {
					if body, err = yaml.JSONToYAML(body); err != nil {
						return errors.Wrap(err, "failed to convert json to yaml")
					}
					utils.Log.Info(string(body))
					continue
				}
			}
		}

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return errors.Wrap(err, "failed to convert json to yaml")
			}
		} else if outFormatFlag != "json" {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}
		if viewAllFlag || isID {
			utils.Log.Info(string(body))
		}
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all applications available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
