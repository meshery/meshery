package perf

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "view perf profile",
	Long:  `See the configuration of your performance profile`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		pattern := args[0]
		var req *http.Request
		url := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles/" + pattern

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return errors.Wrapf(err, utils.PerfError("Failed to invoke performance test"))
		}
		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode != 200 {
			// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
			return errors.Errorf("Response Status Code %d, possible invalid ID", resp.StatusCode)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}
		var dat map[string]interface{}
		if err = json.Unmarshal(data, &dat); err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}
		log.Info(string(data))
		return nil
	},
}
