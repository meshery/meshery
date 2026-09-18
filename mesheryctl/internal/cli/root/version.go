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

package root

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/printer"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type versionInfo struct {
	Version string `json:"version" yaml:"version"`
	GitSHA  string `json:"gitSha" yaml:"gitSha"`
}

type versionOutput struct {
	Client versionInfo `json:"client" yaml:"client"`
	Server versionInfo `json:"server" yaml:"server"`
}

// stringOutput make human string output for version command
func (v versionOutput) stringOutput() string {
	clientSHA := formatShortSHA(v.Client.GitSHA)

	serverSHA := v.Server.GitSHA
	if serverSHA != "unavailable" {
		serverSHA = formatShortSHA(serverSHA)
	}

	out := fmt.Sprintf("Client Version: %s", v.Client.Version)
	if clientSHA != "" && clientSHA != "unavailable" {
		out += fmt.Sprintf(" (%s)", clientSHA)
	}

	out += fmt.Sprintf("\nServer Version: %s", v.Server.Version)
	if serverSHA != "" && serverSHA != "unavailable" {
		out += fmt.Sprintf(" (%s)", serverSHA)
	}

	return out + "\n"
}

// formatShortSHA make commitSHA shorter for human output
func formatShortSHA(sha string) string {
	if len(sha) > 7 {
		return sha[:7]
	}
	return sha
}

var (
	// Mesheryctl config - holds config handler
	mctlCfg *config.MesheryCtlConfig

	outputFormat string
)

var linkDoc = map[string]string{
	"link":    "![version-usage](../../images/version.png)",
	"caption": "Usage of mesheryctl version",
}

func init() {
	printer.AddFormatFlag(versionCmd, &outputFormat)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show Meshery CLI and Server versions",
	Long:  `Version of Meshery command line client - mesheryctl.`,
	Example: `
// To view the current version and SHA of release binary of mesheryctl client 
mesheryctl version
	`,
	Annotations: linkDoc,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			// get the currCtx
			utils.Log.Error(ErrProcessingConfig(err))
			userResponse := false
			userResponse = utils.AskForConfirmation("Looks like you are using an outdated config file. Do you want to generate a new config file?")
			if userResponse {
				utils.BackupConfigFile(utils.DefaultConfigPath)
				// Create config file if not present in meshery folder
				err = utils.CreateConfigFile()
				if err != nil {
					utils.Log.Error(ErrCreatingConfigFile)
				}

				// Add Token to context file
				err = config.AddTokenToConfig(utils.TemplateToken, utils.DefaultConfigPath)
				if err != nil {
					utils.Log.Error(ErrAddingTokenToConfig)
				}

				// Add Context to context file
				err = config.AddContextToConfig("local", utils.TemplateContext, utils.DefaultConfigPath, true, false)
				if err != nil {
					utils.Log.Error(ErrAddingContextToConfig)
				}

				utils.Log.Info(
					fmt.Sprintf("Default config file created at %s",
						utils.DefaultConfigPath,
					))

				mctlCfg, err = config.GetMesheryCtl(viper.GetViper())
				if err != nil {
					utils.Log.Error(ErrUnmarshallingConfigFile)
				}
				currCtx, err := mctlCfg.GetCurrentContext()
				if err != nil {
					return err
				}
				err = currCtx.ValidateVersion()
				if err != nil {
					return err
				}
				return nil
			}
			return models.ErrUnmarshal(errors.New("invalid config file, encountered error processing json"), "meshconfig")
		}
		currCtx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = currCtx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		activePrinter, err := printer.New(outputFormat, cmd.OutOrStdout())

		if err != nil {
			return err
		}

		build := constants.GetMesheryctlVersion()
		commitsha := constants.GetMesheryctlCommitsha()

		isHumanOutput := outputFormat == "" || outputFormat == "table" || outputFormat == "string"

		if isHumanOutput {
			defer utils.CheckMesheryctlClientVersion(build)
		}

		out := versionOutput{
			Client: versionInfo{Version: build, GitSHA: commitsha},
			Server: versionInfo{Version: "unavailable", GitSHA: "unavailable"},
		}

		if mctlCfg != nil {
			if srv, err := fetchServerVersion(mctlCfg.GetBaseMesheryURL()); err == nil {
				out.Server = srv
			} else if isHumanOutput {
				utils.Log.Warn(ErrConnectingToServer(err))
			}
		}

		return activePrinter.Print(out, func(w io.Writer) error {
			_, err := fmt.Fprint(w, out.stringOutput())
			return err
		})
	},
}

// fetchServerVersion helper function to fetch version info from api
func fetchServerVersion(baseURL string) (versionInfo, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(fmt.Sprintf("%s/api/system/version", baseURL))
	if err != nil {
		return versionInfo{}, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return versionInfo{}, fmt.Errorf("server returned status: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return versionInfo{}, err
	}

	var raw config.Version
	if err := json.Unmarshal(data, &raw); err != nil {
		return versionInfo{}, err
	}

	return versionInfo{
		Version: raw.GetBuild(),
		GitSHA:  raw.GetCommitSHA(),
	}, nil
}
