package perf

import (
	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Performance profiles",
	Long:    `List all the available performance profiles`,
	Example: "mesheryctl perf list",
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		performanceProfiles := mctlCfg.PerfProfiles
		var data [][]string
		for _, profile := range performanceProfiles {
			data = append(data, []string{profile.Name})
		}
		utils.PrintToTable([]string{"NAME", "RESULTS", "LAST-RUN"}, data)

	},
}
