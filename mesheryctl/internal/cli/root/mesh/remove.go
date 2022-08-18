package mesh

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [mesh-name]",
	Short: "Remove a service mesh",
	Long:  "Remove a service mesh from the connected Kubernetes cluster",
	Example: `  // Remove Linkerd:
  mesheryctl mesh remove linkerd --namespace linkerd-ns`,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := utils.CreateDefaultSpinner(
			fmt.Sprintf("Removing %s", meshName),
			fmt.Sprintf("\n%s removed successfully", meshName),
		)

		s.Start()
		_, err = sendOperationRequest(mctlCfg, strings.ToLower(meshName), true)
		if err != nil {
			return errors.Wrap(err, "error installing service mesh")
		}
		s.Stop()

		return nil
	},
}

func init() {
	removeCmd.Flags().StringVarP(
		&namespace, "namespace", "n", "default",
		"Kubernetes namespace where the mesh is deployed",
	)
}
