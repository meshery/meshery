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

package filter

import (
	"errors"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/server/models"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type filterViewFlags struct {
	viewAllFlag  bool
	outputFormat string
	save         bool
}

var filterViewFlagsProvided filterViewFlags

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View filter(s)",
	Long:  `Displays the contents of a specific filter based on name or id`,
	Example: `
// View the specified WASM filter
// A unique prefix of the name or ID can also be provided. If the prefix is not unique, the first match will be returned.
mesheryctl filter view "[filter-name | ID]"

// View all filter files
mesheryctl filter view --all

// View all filter files in json
mesheryctl filter view --all --output-format json

// View all filter files in json and save it to a file
mesheryctl filter view --all --output-format json -s

//View multi-word named filter files. Multi-word filter names should be enclosed in quotes
mesheryctl filter view "filter name"
        `,
	Args: cobra.ArbitraryArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// Validate output-format
		return display.ValidateOutputFormat(filterViewFlagsProvided.outputFormat)
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		// for formatting errors
		subCmdUsed := cmd.Use

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}

		filter := ""
		isID := false
		// if filter name/id available
		if len(args) > 0 {
			if filterViewFlagsProvided.viewAllFlag {
				return ErrViewAllWithName(subCmdUsed)
			}

			filterArg, err := parseQuotedArg(args, subCmdUsed)
			if err != nil {
				return err
			}

			filter, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), filterArg, "filter")
			if err != nil {
				return utils.ErrInvalidNameOrID(err)
			}
		}

		urlString := ""
		if len(filter) == 0 {
			if filterViewFlagsProvided.viewAllFlag {
				urlString = "api/filter?pagesize=10000"
			} else {
				return utils.ErrInvalidNameOrID(errors.New(errFilterNameOrIDNotProvided))
			}
		} else if isID {
			// if filter is a valid uuid, then directly fetch the filter
			urlString = "api/filter/" + filter
		} else {
			// else search filter by name
			urlString = "api/filter?search=" + url.QueryEscape(filter)
		}

		var selectedFilter *models.MesheryFilter
		var filterPage *models.MesheryFilterPage
		if isID {
			// Fetch single filter
			filter, err := api.Fetch[models.MesheryFilter](urlString)
			if err != nil {
				return err
			}
			selectedFilter = filter
		} else {
			page, err := api.Fetch[models.MesheryFilterPage](urlString)
			if err != nil {
				return err
			}
			filterPage = page

			if !filterViewFlagsProvided.viewAllFlag {
				if len(filterPage.Filters) == 0 {
					utils.Log.Info(fmt.Sprintf("filter with name: %q not found", filter))
					return nil
				}
				selectedFilter = page.Filters[0]
			}
		}

		outputFormatterFactory := display.OutputFormatterFactory[any]{}
		var data any

		if filterViewFlagsProvided.viewAllFlag {
			data = filterPage
		} else {
			data = selectedFilter
		}

		outputFormatter, err := outputFormatterFactory.New(filterViewFlagsProvided.outputFormat, data)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if filterViewFlagsProvided.save {
			err := saveToFile(
				filterViewFlagsProvided.outputFormat,
				outputFormatter,
				selectedFilter,
				filterViewFlagsProvided.viewAllFlag,
			)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func saveToFile(
	outputFormat string,
	outputFormatter display.OutputFormatter[any],
	selectedFilter *models.MesheryFilter,
	isAll bool,
) error {
	var fileName string

	if isAll {
		fileName = "filters_all"
	} else {
		shortID := selectedFilter.ID.String()[:8]
		sanitizer := strings.NewReplacer("/", "_", " ", "_")
		sanitizedName := sanitizer.Replace(selectedFilter.Name)
		fileName = fmt.Sprintf("filter_%s_%s", sanitizedName, shortID)
	}

	file := filepath.Join(utils.MesheryFolder, fileName)

	outputFormatterSaverFactory := display.OutputFormatterSaverFactory[any]{}
	outputFormatterSaver, err := outputFormatterSaverFactory.New(outputFormat, outputFormatter)
	if err != nil {
		return err
	}
	outputFormatterSaver = outputFormatterSaver.WithFilePath(file)
	return outputFormatterSaver.Save()
}

// Check if the argument starts and ends with double quotes
func parseQuotedArg(args []string, subCmdUsed string) (string, error) {
	fullArg := strings.Join(args, " ")

	if strings.HasPrefix(fullArg, "\"") && strings.HasSuffix(fullArg, "\"") {
		return strings.Trim(fullArg, "\""), nil
	} else if len(args) == 1 {
		return args[0], nil
	}
	return "", ErrMultiWordFilterName(subCmdUsed)
}

func init() {
	viewCmd.Flags().BoolVarP(&filterViewFlagsProvided.viewAllFlag, "all", "a", false, "(optional) view all filters available")
	viewCmd.Flags().StringVarP(&filterViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewCmd.Flags().BoolVarP(&filterViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}
