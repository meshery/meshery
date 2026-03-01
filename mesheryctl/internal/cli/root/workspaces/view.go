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

package workspaces

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/manifoldco/promptui"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type workspaceViewFlags struct {
	outputFormat string
	save         bool
}

var workspaceViewFlagsProvided workspaceViewFlags

var viewWorkspaceCmd = &cobra.Command{
	Use:   "view",
	Short: "View a workspace",
	Long: `View a workspace by its ID or name.
Find more information at: https://docs.meshery.io/reference/mesheryctl/exp/workspace/view`,
	Example: `
// View details of a specific workspace in default format (yaml)
mesheryctl exp workspace view [workspace-name|workspace-id]

// View details of a specific workspace in JSON format
mesheryctl exp workspace view [workspace-name|workspace-id] --output-format json

// View details of a specific workspace and save it to a file
mesheryctl exp workspace view [workspace-name|workspace-id] --output-format json --save
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp workspace view [workspace-name|workspace-id]\nRun 'mesheryctl exp workspace view --help' to see detailed help message"
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("workspace name or ID isn't specified\n\n%v", errMsg))
		}
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("too many arguments\n\n%v", errMsg))
		}
		return display.ValidateOutputFormat(workspaceViewFlagsProvided.outputFormat)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceNameOrID := args[0]
		var selectedWorkspace *models.Workspace

		if isWorkspaceArgumentUUID(workspaceNameOrID) {
			fetchedWorkspace, err := fetchWorkspaceByID(workspaceNameOrID)
			if err != nil {
				return err
			}
			selectedWorkspace = fetchedWorkspace
		} else {
			fetchedWorkspace, err := fetchWorkspaceByName(workspaceNameOrID)
			if err != nil {
				return err
			}
			if fetchedWorkspace == nil {
				fmt.Println("No workspace(s) found for the given name: ", workspaceNameOrID)
				return nil
			}
			selectedWorkspace = fetchedWorkspace
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return utils.ErrRetrieveHomeDir(errors.Wrap(err, "failed to determine user home directory"))
		}

		outputFormatterFactory := display.OutputFormatterFactory[models.Workspace]{}
		outputFormatter, err := outputFormatterFactory.New(workspaceViewFlagsProvided.outputFormat, *selectedWorkspace)
		if err != nil {
			return err
		}

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if workspaceViewFlagsProvided.save {
			workspaceString := strings.ReplaceAll(selectedWorkspace.Name, " ", "_")
			if workspaceString == "" {
				workspaceString = selectedWorkspace.ID.String()
			}
			fileName := fmt.Sprintf("workspace_%s.%s", workspaceString, strings.ToLower(workspaceViewFlagsProvided.outputFormat))
			file := filepath.Join(homeDir, ".meshery", fileName)

			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[models.Workspace]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(workspaceViewFlagsProvided.outputFormat, outputFormatter)
			if err != nil {
				return err
			}
			outputFormatterSaver = outputFormatterSaver.WithFilePath(file)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func selectWorkspacePrompt(workspacesList []models.Workspace) (*models.Workspace, error) {
	workspaceNames := []string{}

	for _, w := range workspacesList {
		workspaceName := fmt.Sprintf("ID: %s, Name: %s, OrgID: %s", w.ID.String(), w.Name, w.OrganizationID.String())
		workspaceNames = append(workspaceNames, workspaceName)
	}

	prompt := promptui.Select{
		Label: "Select workspace",
		Items: workspaceNames,
	}

	i, _, err := prompt.Run()
	if err != nil {
		if errors.Is(err, promptui.ErrInterrupt) {
			return nil, errors.New("workspace selection cancelled")
		}
		return nil, err
	}
	return &workspacesList[i], nil
}

func isWorkspaceArgumentUUID(arg string) bool {
	_, err := uuid.FromString(arg)
	return err == nil
}

func fetchWorkspaceByID(workspaceID string) (*models.Workspace, error) {
	urlPath := fmt.Sprintf("%s/%s", workspacesApiPath, workspaceID)
	fetchedWorkspace, err := api.Fetch[models.Workspace](urlPath)
	if err != nil {
		return nil, err
	}
	return fetchedWorkspace, nil
}

func fetchWorkspaceByName(workspaceName string) (*models.Workspace, error) {
	viewUrlValue := url.Values{}
	viewUrlValue.Add("search", workspaceName)
	viewUrlValue.Add("pagesize", "all")

	urlPath := fmt.Sprintf("%s?%s", workspacesApiPath, viewUrlValue.Encode())

	workspacesResponse, err := api.Fetch[models.WorkspacePage](urlPath)
	if err != nil {
		return nil, err
	}

	if workspacesResponse.TotalCount == 0 {
		return nil, nil
	}

	if workspacesResponse.TotalCount > 1 {
		return selectWorkspacePrompt(workspacesResponse.Workspaces)
	}

	return &workspacesResponse.Workspaces[0], nil
}

func init() {
	viewWorkspaceCmd.Flags().StringVarP(&workspaceViewFlagsProvided.outputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewWorkspaceCmd.Flags().BoolVarP(&workspaceViewFlagsProvided.save, "save", "s", false, "(optional) save output as a JSON/YAML file")
}