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

package system

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// represents the `mesheryctl system model list` subcommand.
var listModelCmd = &cobra.Command{
	Use:   "list",
	Short: "list models",
	Long:  "list name of all registered models",
	Example: `
// View current provider
mesheryctl system model list
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite
		hcOptions := &HealthCheckOptions{
			IsPreRunE:  true,
			PrintLogs:  false,
			Subcommand: cmd.Use,
		}
		hc, err := NewHealthChecker(hcOptions)
		if err != nil {
			return ErrHealthCheckFailed(err)
		}
		// execute healthchecks
		err = hc.RunPreflightHealthChecks()
		if err != nil {
			cmd.SilenceUsage = true
			return err
		}
		cfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		ctx, err := cfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()

		url := fmt.Sprintf("%s/api/meshmodels/models?pagesize=all", baseUrl)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(ErrGettingRequestContext(err))
			return err
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			utils.Log.Error(ErrConnectingToServer(err))
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(ErrInvalidAPIResponse(err))
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(ErrUnmarshallingAPIData(err))
			return err
		}

		header := []string{"Model Name"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.DisplayName})
			}
		}

		if len(rows) == 0 {
			// if no model is found
			utils.Log.Info("No models found")
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}

// providerCmd represents the `mesheryctl system model` command
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Generate and update components",
	Long:  `Generates and updates components in binary form`,
	Example: `
// To view list of components
mesheryctl system model list

// To view a specific model
mesheryctl system model view [model-name]

// To update model
mesheryctl system model update --source [path-to-spreadsheet] [--output-dir] [relative path persist assets/icons]

// To generate model
mesheryctl system model gen
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New(utils.SystemModelSubError("please specify a flag or subcommand. Use 'mesheryctl system model --help' to display user guide.\n", "model"))
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl system model --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		err = viewProviderCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	availableSubcommands = []*cobra.Command{listModelCmd}
	modelCmd.AddCommand(availableSubcommands...)
}
