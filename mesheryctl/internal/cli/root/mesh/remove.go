package mesh

import (
	"fmt"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	removeCmd = &cobra.Command{
		Use:   "remove",
		Short: "remove a service mesh in the kubernetes cluster",
		Args:  cobra.MinimumNArgs(1),
		Long:  `remove service mesh in the connected kubernetes cluster`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Verifying prerequisites...")
			if len(args) < 1 {
				return errors.New("no service mesh specified")
			}

			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			if err = validateAdapter(mctlCfg, tokenPath, args[0]); err != nil {
				return errors.Wrap(err, "adapter not valid")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			s := utils.CreateDefaultSpinner(fmt.Sprintf("Removing %s", args[0]), fmt.Sprintf("\n%s service mesh removed successfully", args[0]))
			mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
			if err != nil {
				return errors.Wrap(err, "error processing config")
			}

			s.Start()
			_, err = sendDeployRequest(mctlCfg, args[0], true)
			if err != nil {
				return errors.Wrap(err, "error installing service mesh")
			}
			s.Stop()

			//log.Infof("Verifying Installation")
			//s.Start()
			//_, err = waitForDeployResponse(mctlCfg, args[0])
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
	removeCmd.Flags().StringVarP(&tokenPath, "tokenPath", "t", "", "Path to token for authenticating to Meshery API")
	_ = removeCmd.MarkFlagRequired("tokenPath")
}
