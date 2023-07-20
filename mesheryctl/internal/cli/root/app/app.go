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
	"fmt"
	"io"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
	validSourceTypes     []string
)

// AppCmd represents the root command for app commands
var AppCmd = &cobra.Command{
	Use:   "app",
	Short: "Cloud Native Apps Management",
	Long:  `All apps operations: import, list, view, onboard and offboard`,
	Example: `
// Base command
mesheryctl app [subcommand]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			availableSubCmds := []string{"onboard", "offboard", "list", "import", "view"}

			suggestedCmd := utils.FindClosestArg(args[0], availableSubCmds)
			if suggestedCmd != "" && suggestedCmd[0] == args[0][0] {
				return errors.New(utils.AppError(fmt.Sprintf("'%s' is a invalid command for '%s'. Did you mean this?\n\t%s\n", args[0], cmd.CalledAs(), suggestedCmd)))
			}
			return errors.New(utils.AppError(fmt.Sprintf("'%s' is a invalid command for '%s'. Use 'mesheryctl perf --help' to display usage guide.\n", args[0], cmd.CalledAs())))
		}
		return nil
	},
}

func init() {
	AppCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{onboardCmd, viewCmd, offboardCmd, listCmd, importCmd}
	AppCmd.AddCommand(availableSubcommands...)
}

func getSourceTypes() error {
	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}
	validTypesURL := mctlCfg.GetBaseMesheryURL() + "/api/application/types"
	req, err := utils.NewRequest("GET", validTypesURL, nil)
	if err != nil {
		return err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response []*models.ApplicationSourceTypesAPIResponse

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Log.Debug("failed to read response body")
		return errors.Wrap(err, "couldn't read response from server. Please try again after some time")
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		utils.Log.Debug("failed to unmarshal JSON response body")
		return errors.Wrap(err, "couldn't process response received from server")
	}

	for _, apiResponse := range response {
		validSourceTypes = append(validSourceTypes, apiResponse.ApplicationType)
	}

	return nil
}

func isValidSource(sType string) bool {
	for _, validType := range validSourceTypes {
		if validType == sType {
			return true
		}
	}

	return false
}
