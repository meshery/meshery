package workspaces

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/display"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/schemas/models/v1beta1/workspace"
	"github.com/spf13/cobra"
)

type workspaceViewFlags struct {
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
	OrgID        string `json:"orgId" validate:"required,uuid"`
}

var workspaceViewFlagsProvided workspaceViewFlags

func formatWorkspaceLabel(rows []workspace.AvailableWorkspace) []string {
	labels := []string{}
	for _, w := range rows {
		labels = append(labels, fmt.Sprintf("%s (ID: %s)", w.Name, w.ID.String()))
	}
	return labels
}

var viewWorkspaceCmd = &cobra.Command{
	Use:   "view [workspace-name|workspace-id]",
	Short: "View a workspace",
	Long: `View a workspace by its ID or name.
Find more information at: https://docs.meshery.io/reference/mesheryctl/workspace/view`,
	Example: `
// View details of a specific workspace by ID
mesheryctl workspace view [workspace-id] --orgId [orgId]

// View details of a specific workspace by name
mesheryctl workspace view [workspace-name] --orgId [orgId]

// View details of a specific workspace in JSON format
mesheryctl workspace view [workspace-id] --orgId [orgId] --output-format json

// View details of a specific workspace and save it to a file
mesheryctl workspace view [workspace-id] --orgId [orgId] --output-format json --save
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &workspaceViewFlagsProvided)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl workspace view [workspace-name|workspace-id]\nRun 'mesheryctl workspace view --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace name or ID\n\n%v", errMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceNameOrID := args[0]
		selectedWorkspace := new(workspace.Workspace)

		var urlPath string
		var displayData display.DisplayDataAsync

		if utils.IsUUID(workspaceNameOrID) {
			query := url.Values{}
			query.Set("orgID", workspaceViewFlagsProvided.OrgID)
			urlPath = fmt.Sprintf("%s/%s?%s", workspacesApiPath, url.PathEscape(workspaceNameOrID), query.Encode())
			fetchedWorkspace, err := api.Fetch[workspace.Workspace](urlPath)
			if err != nil {
				return err
			}
			*selectedWorkspace = *fetchedWorkspace
		} else {
			selectedAvailableWorkspace := new(workspace.AvailableWorkspace)
			viewUrlValue := url.Values{}
			viewUrlValue.Add("orgID", workspaceViewFlagsProvided.OrgID)
			viewUrlValue.Add("search", workspaceNameOrID)

			urlPath = fmt.Sprintf("%s?%s", workspacesApiPath, viewUrlValue.Encode())
			displayData = display.DisplayDataAsync{
				UrlPath:    urlPath,
				SearchTerm: workspaceNameOrID,
			}

			err := display.PromptAsyncPagination(
				displayData,
				formatWorkspaceLabel,
				func(data *workspace.WorkspacePage) ([]workspace.AvailableWorkspace, int64) {
					return data.Workspaces, int64(data.TotalCount)
				},
				selectedAvailableWorkspace,
			)
			if err != nil {
				return err
			}

			query := url.Values{}
			query.Set("orgID", workspaceViewFlagsProvided.OrgID)
			urlPath = fmt.Sprintf("%s/%s?%s", workspacesApiPath, url.PathEscape(selectedAvailableWorkspace.ID.String()), query.Encode())
			fetchedWorkspace, err := api.Fetch[workspace.Workspace](urlPath)
			if err != nil {
				return err
			}
			*selectedWorkspace = *fetchedWorkspace
		}

		outputFormatterFactory := display.OutputFormatterFactory[workspace.Workspace]{}
		outputFormatter, err := outputFormatterFactory.New(workspaceViewFlagsProvided.OutputFormat, *selectedWorkspace)
		if err != nil {
			return err
		}

		outputFormatter = outputFormatter.WithOutput(cmd.OutOrStdout())

		err = outputFormatter.Display()
		if err != nil {
			return err
		}

		if workspaceViewFlagsProvided.Save {
			workspaceString := strings.ReplaceAll(selectedWorkspace.Name, " ", "_")
			if workspaceString == "" {
				workspaceString = selectedWorkspace.ID.String()
			}
			fileName := filepath.Join(utils.MesheryFolder, fmt.Sprintf("workspace_%s.%s", workspaceString, workspaceViewFlagsProvided.OutputFormat))
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[workspace.Workspace]{}
			outputFormatterSaver, err := outputFormatterSaverFactory.New(workspaceViewFlagsProvided.OutputFormat, outputFormatter)
			if err != nil {
				return err
			}
			outputFormatterSaver = outputFormatterSaver.WithFilePath(fileName)
			err = outputFormatterSaver.Save()
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	viewWorkspaceCmd.Flags().StringVarP(&workspaceViewFlagsProvided.OutputFormat, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	viewWorkspaceCmd.Flags().BoolVarP(&workspaceViewFlagsProvided.Save, "save", "s", false, "(optional) save output as a JSON/YAML file")
	viewWorkspaceCmd.Flags().StringVarP(&workspaceViewFlagsProvided.OrgID, "orgId", "", "", "(required) organization ID")
}
