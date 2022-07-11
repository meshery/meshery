package mesh

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove a service mesh in the kubernetes cluster",
		Args:  cobra.MinimumNArgs(0),
		Long:  `remove service mesh in the connected kubernetes cluster`,
		Example: `
// Remove a service mesh
mesheryctl mesh remove [mesh adapter name]
		`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			utils.Log.Info("Verifying prerequisites...")
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			if len(args) < 1 {
				meshName, err = validateMesh(mctlCfg, "")
			} else {
				meshName, err = validateMesh(mctlCfg, args[0])
			}
			if err != nil {
				return errors.Wrap(err, "error validating request")
			}

			if err = validateAdapter(mctlCfg, meshName); err != nil {
				return errors.Wrap(err, "adapter not valid")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := utils.CreateDefaultSpinner(fmt.Sprintf("Removing %s", meshName), fmt.Sprintf("\n%s service mesh removed successfully", meshName))
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			s.Start()
			_, err = sendDeployRequest(mctlCfg, meshName, true)
			if err != nil {
				return errors.Wrap(err, "error installing service mesh")
			}
			s.Stop()

			//log.Infof("Verifying Installation")
			//s.Start()
			//_, err = waitForDeployResponse(mctlCfg, meshName)
			//if err != nil {
			//	return errors.Wrap(err, "error verifying installation")
			//}
			//s.Stop()

			return nil
		},
	}
)

func init() {
	removeCmd.Flags().StringVarP(&adapterURL, "adapter", "a", "meshery-istio:10000", "Adapter to use for installation")
	removeCmd.Flags().StringVarP(&namespace, "namespace", "n", "default", "Kubernetes namespace to be used for deploying the validation tests and sample workload")
	removeCmd.Flags().StringVarP(&utils.TokenFlag, "token", "t", "", "Path to token for authenticating to Meshery API")
}
