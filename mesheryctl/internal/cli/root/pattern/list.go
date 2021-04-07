package pattern

import (
	"encoding/json"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	token string
)
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list pattern files",
	Long:  "List available pattern files",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		var Patterns models.PatternsApiResponse

		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile/", nil)
		if err != nil {
			return err
		}
		err = utils.AddAuthDetails(req, "./auth.json")
		if err != nil {
			return err
		}
		client.Do(req)
		err = json.NewDecoder(req.Body).Decode(&Patterns)
		if err != nil {
			return err
		}

	},
}
