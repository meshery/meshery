package workspaces

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	mErrors "github.com/meshery/meshkit/errors"
	"github.com/meshery/schemas/models/v1beta3/workspace"

	"github.com/spf13/cobra"
)

type cmdWorkspaceUpdateFlags struct {
	OrganizationID string `json:"orgId" validate:"required,uuid"`
	Name           string `json:"name" validate:"omitempty"`
	Description    string `json:"description" validate:"omitempty"`
}

var workspaceUpdateFlags cmdWorkspaceUpdateFlags

var updateWorkspaceCmd = &cobra.Command{
	Use:   "update [workspace-id]",
	Short: "Update a workspace",
	Long: `Update a workspace's name and/or description by its ID.
At least one of --name or --description must be provided.
Find more information at: https://docs.meshery.io/reference/references/mesheryctl/workspace/update`,
	Example: `
// Rename a workspace
mesheryctl workspace update [workspace-id] --orgId [orgId] --name [new-name]

// Update a workspace's description
mesheryctl workspace update [workspace-id] --orgId [orgId] --description [new-description]

// Update both
mesheryctl workspace update [workspace-id] --orgId [orgId] --name [new-name] --description [new-description]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mesheryctlflags.ValidateCmdFlags(cmd, &workspaceUpdateFlags)
	},
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl workspace update [workspace-id] --orgId [orgId] [--name NAME] [--description DESCRIPTION]\nRun 'mesheryctl workspace update --help' to see detailed help message"

		if len(args) != 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("please provide exactly one workspace ID\n\n%v", errMsg))
		}
		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %s\n\n%v", args[0], errMsg))
		}
		if workspaceUpdateFlags.Name == "" && workspaceUpdateFlags.Description == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("at least one of --name or --description must be provided\n\n%v", errMsg))
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID := args[0]

		organizationID, err := uuid.FromString(workspaceUpdateFlags.OrganizationID)
		if err != nil {
			return err
		}

		// The server merges partial updates - only fields present in the
		// payload change - but WorkspaceUpdatePayload always requires
		// OrganizationID, so it is included on every request even though
		// only Name/Description are ever user-editable through this command.
		workspaceUpdatePayload := workspace.WorkspaceUpdatePayload{
			OrganizationID: organizationID,
			Name:           workspaceUpdateFlags.Name,
			Description:    workspaceUpdateFlags.Description,
		}
		payloadBytes, err := json.Marshal(workspaceUpdatePayload)
		if err != nil {
			return utils.ErrUnmarshal(err)
		}

		_, err = api.Update(fmt.Sprintf("%s/%s", workspacesApiPath, workspaceID), bytes.NewBuffer(payloadBytes), nil)
		if err != nil {
			if mErrors.GetCode(err) == utils.ErrNotFoundCode {
				utils.Log.Infof("Workspace with ID %s not found", workspaceID)
				return nil
			}
			return err
		}

		utils.Log.Infof("Workspace with ID %s has been updated", workspaceID)
		return nil
	},
}

func init() {
	updateWorkspaceCmd.Flags().StringVarP(&workspaceUpdateFlags.OrganizationID, "orgId", "", "", "(required) organization ID")
	updateWorkspaceCmd.Flags().StringVarP(&workspaceUpdateFlags.Name, "name", "n", "", "New name for the workspace")
	updateWorkspaceCmd.Flags().StringVarP(&workspaceUpdateFlags.Description, "description", "d", "", "New description for the workspace")
}
