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
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete filter file",
	Long:  `delete filter file will trigger deletion of the filter file`,
	Example: `
// Delete the specified WASM filter file using name or ID
mesheryctl exp filter delete [filter-name | ID]

// Delete using the file name
mesheryctl exp filter delete test-wasm
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		client := &http.Client{}
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		filter := ""
		isID := false
		if len(args) > 0 {
			filter, isID, err = utils.Valid(args[0], "filter")
			if err != nil {
				return err
			}
		}

		// Delete the filter using the id
		if isID {
			err := utils.DeleteConfiguration(filter, "filter")
			if err != nil {
				return errors.Wrap(err, utils.FilterError(fmt.Sprintf("failed to delete filter %s", args[0])))
			}
			utils.Log.Info("Filter ", args[0], " deleted successfully")
			return nil
		}

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
		}

		req, err := utils.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/filter/deploy", fileReader)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		if res.StatusCode != 200 {
			return ErrInvalidAPICall(res.StatusCode)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return errors.New("Error processing response body. " + err.Error())
		}

		utils.Log.Info(string(body))

		return nil
	},
}
