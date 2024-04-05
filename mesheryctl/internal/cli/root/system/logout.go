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

package system

import (
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove authentication for Meshery Server",
	Long: `
Remove authentication for Meshery Server

This command removes the authentication token from the user's filesystem`,
	Args: cobra.MinimumNArgs(0),
	Example: `
// Logout current session with your Meshery Provider.
mesheryctl system logout
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		token, err := mctlCfg.GetTokenForContext(mctlCfg.GetCurrentContextName())
		if err != nil {
			utils.Log.Error(ErrGetCurrentContext(errors.Wrap(err, "failed to find token path for the current context")))
			return nil
		}

		// Replace the content of the token file with empty content
		if err := os.WriteFile(token.GetLocation(), []byte{}, 0666); err != nil {
			log.Error("logout failed: ", err)
			return nil
		}

		log.Println("successfully logged out")
		return nil
	},
}
