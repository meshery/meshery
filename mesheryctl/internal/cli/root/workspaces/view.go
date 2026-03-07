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
	"github.com/meshery/meshery/server/models"
	"github.com/spf13/cobra"
)

type workspaceViewFlags struct {
	OutputFormat string `json:"output-format" validate:"required,oneof=json yaml"`
	Save         bool   `json:"save" validate:"boolean"`
	OrgID        string `json:"orgId"`
}

var workspaceViewFlagsProvided workspaceViewFlags

var viewWorkspaceCmd = &cobra.Command{
	Use:   "view [workspace-name|workspace-id]",
	Short: "View a workspace",
	Long: `View a workspace by its ID or name.
Find more information at: https://docs.meshery.io/reference/mesheryctl/exp/workspace/view`,
	Example: `
// View details of a specific workspace by ID
mesheryctl exp workspace view [workspace-id]

// View details of a specific workspace by name (requires --orgId)
mesheryctl exp workspace view [workspace-name] --orgId [orgId]

// View details of a specific workspace in JSON format
mesheryctl exp workspace view [workspace-id] --output-format json

// View details of a specific workspace and save it to a file
mesheryctl exp workspace view [workspace-id] --output-format json --save
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		flagValidator, ok := cmd.Context().Value(mesheryctlflags.FlagValidatorKey).(*mesheryctlflags.FlagValidator)
		if !ok || flagValidator == nil {
			return utils.ErrCommandContextMissing("flags-validator")
		}
		err := flagValidator.Validate(workspaceViewFlagsProvided)
		if err != nil {
			return utils.ErrFlagsInvalid(err)
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp workspace view [workspace-name|workspace-id]\nRun 'mesheryctl exp workspace view --help' to see detailed help message"
		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace name or ID\n\n%v", errMsg))
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceNameOrID := args[0]
		orgID := workspaceViewFlagsProvided.OrgID

		var selectedWorkspace *models.Workspace

		if utils.IsUUID(workspaceNameOrID) {
			urlPath := fmt.Sprintf("%s/%s", workspacesApiPath, workspaceNameOrID)
			fetchedWorkspace, err := api.Fetch[models.Workspace](urlPath)
			if err != nil {
				return err
			}
			selectedWorkspace = fetchedWorkspace
		} else {
			if orgID == "" {
				return utils.ErrInvalidArgument(fmt.Errorf("--orgId is required when searching by name\n\nUsage: mesheryctl exp workspace view [workspace-name] --orgId [orgId]"))
			}

			viewUrlValue := url.Values{}
			viewUrlValue.Add("search", workspaceNameOrID)
			viewUrlValue.Add("orgID", orgID)
			viewUrlValue.Add("pagesize", "all")

			urlPath := fmt.Sprintf("%s?%s", workspacesApiPath, viewUrlValue.Encode())

			workspacesResponse, err := api.Fetch[models.WorkspacePage](urlPath)
			if err != nil {
				return err
			}

			if workspacesResponse.TotalCount == 0 {
				fmt.Println("No workspace(s) found with the name:", workspaceNameOrID)
				return nil
			}

			if workspacesResponse.TotalCount == 1 {
				selectedWorkspace = &workspacesResponse.Workspaces[0]
			} else {
				workspaceNames := make([]string, len(workspacesResponse.Workspaces))
				for i, w := range workspacesResponse.Workspaces {
					workspaceNames[i] = fmt.Sprintf("ID: %s, Name: %s, OrgID: %s", w.ID.String(), w.Name, w.OrganizationID.String())
				}
				i, err := utils.RunSelectPrompt("Select workspace", workspaceNames)
				if err != nil {
					return err
				}
				selectedWorkspace = &workspacesResponse.Workspaces[i]
			}
		}

		outputFormatterFactory := display.OutputFormatterFactory[models.Workspace]{}
		outputFormatter, err := outputFormatterFactory.New(workspaceViewFlagsProvided.OutputFormat, *selectedWorkspace)
		if err != nil {
			return err
		}

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
			outputFormatterSaverFactory := display.OutputFormatterSaverFactory[models.Workspace]{}
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
	viewWorkspaceCmd.Flags().StringVarP(&workspaceViewFlagsProvided.OrgID, "orgId", "", "", "(optional) organization ID to search workspace by name")
}