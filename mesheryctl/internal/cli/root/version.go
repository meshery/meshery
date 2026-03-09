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

package root

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var linkDoc = map[string]string{
	"link":    "![version-usage](/assets/img/mesheryctl/version.png)",
	"caption": "Usage of mesheryctl version",
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:         "version",
	Short:       "Show Meshery CLI and Server versions",
	Long:        "Version of Meshery command line client - mesheryctl.",
	Example:     "mesheryctl version",
	Annotations: linkDoc,

	RunE: func(cmd *cobra.Command, args []string) error {
		// --- READ-ONLY CONFIG LOAD (NO VALIDATION) ---
		header := []string{"", "Version", "GitSHA"}
		rows := [][]string{
			{"Client", constants.GetMesheryctlVersion(), constants.GetMesheryctlCommitsha()},
			{"Server", "unavailable", "unavailable"},
		}

		// --- BEST EFFORT SERVER VERSION (NO CONTEXT REQUIREMENT) ---
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.PrintToTable(header, rows, nil)
			return nil
		}

		ctx := mctlCfg.GetCurrentContextSafe()
		if ctx == nil || ctx.Endpoint == "" {
			utils.PrintToTable(header, rows, nil)
			return nil
		}

		req, err := http.NewRequest(
			http.MethodGet,
			fmt.Sprintf("%s/api/system/version", ctx.Endpoint),
			nil,
		)
		if err != nil {
			utils.PrintToTable(header, rows, nil)
			return nil
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			utils.PrintToTable(header, rows, nil)
			return nil
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.PrintToTable(header, rows, nil)
			return nil
		}

		var serverVersion config.Version
		if err := json.Unmarshal(data, &serverVersion); err != nil {
			utils.PrintToTable(header, rows, nil)
			return nil
		}

		rows[1][1] = serverVersion.GetBuild()
		rows[1][2] = serverVersion.GetCommitSHA()

		utils.PrintToTable(header, rows, nil)
		return nil
	},
}
