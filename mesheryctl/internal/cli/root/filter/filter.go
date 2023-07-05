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

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	availableSubcommands []*cobra.Command
	file                 string
)

// FilterCmd represents the root command for filter commands
var FilterCmd = &cobra.Command{
	Use:   "filter",
	Short: "Cloud Native Filter Management",
	Long:  ``,
	Example: `
// Base command for WASM filters (experimental feature)
mesheryctl exp filter [subcommands]	
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.FilterError(fmt.Sprintf("'%s' is a invalid command.  Use 'mesheryctl exp filter --help' to display usage guide.\n", args[0])))
		}
		return nil
	},
}

var importFilterCmd = &cobra.Command{
	Use:   "import [URI]",
	Short: "Import a WASM filter",
	Long:  "Import a WASM filter from a URI (http/s) or local filesystem path",
	Args:  cobra.ExactArgs(1),
	RunE:  importFilter,
}

func importFilter(cmd *cobra.Command, args []string) error {
	uri := args[0]
	configPath, _ := cmd.Flags().GetString("config")

	// Perform the import logic using the provided URI and configPath
	err := importFilterFromURI(uri, configPath)
	if err != nil {
		return err
	}

	fmt.Println("WASM Filter imported")
	return nil
}

// saveFilter saves the filter content from the given reader to the specified file path
func saveFilter(reader io.Reader, filePath string) error {
	// Read the filter content
	filterContent, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	// Write the filter content to the file
	err = os.WriteFile(filePath, filterContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

// saveFilterFromPath saves the filter content from the file at the given path to the specified file path
func saveFilterFromPath(sourcePath, destinationPath string) error {
	// Read the filter content from the source file
	filterContent, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destinationPath, filterContent, 0644)
	if err != nil {
		return err
	}

	return nil
}

func importFilterFromURI(uri, configPath string) error {
	if uri[:7] == "http://" || uri[:8] == "https://" {
		resp, err := http.Get(uri)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		err = saveFilter(resp.Body, configPath)
		if err != nil {
			return err
		}
	} else {
		err := saveFilterFromPath(uri, configPath)
		if err != nil {
			return err
		}
	}

	return nil
}

func init() {
	FilterCmd.PersistentFlags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token file default from current context")

	availableSubcommands = []*cobra.Command{applyCmd, viewCmd, deleteCmd, listCmd}
	FilterCmd.AddCommand(availableSubcommands...)

	importFilterCmd.Flags().StringP("config", "c", "", "Configuration file path")

	FilterCmd.AddCommand(importFilterCmd)
}
