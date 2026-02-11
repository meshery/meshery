package workspaces

import (
	"fmt"

	"github.com/meshery/meshery/mesheryctl/internal/cli/pkg/api"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"

	"github.com/spf13/cobra"
)

var deleteWorkspaceCmd = &cobra.Command{
	Use:   "delete [workspace-id]",
	Short: "Delete a workspace under an organization",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return utils.ErrInvalidArgument(fmt.Errorf("[ Workspace ID ] not specified"))
		}

		if !utils.IsUUID(args[0]) {
			return utils.ErrInvalidUUID(fmt.Errorf("invalid workspace ID: %q", args[0]))
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		workspaceID := args[0]

		_, err := api.Delete(fmt.Sprintf("%s/%s", workspacesApiPath, workspaceID))
		if err != nil {
			return err
		}

		utils.Log.Infof("Workspace with ID %s has been deleted", workspaceID)
		return nil
	},
}
