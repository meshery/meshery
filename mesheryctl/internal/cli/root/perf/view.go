package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

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
		proName := args[0]
		var req *http.Request
		url := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles"

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
			return errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", proName)
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
		//s := dat["profiles"].([]interface{})

		for _, i := range dat["profiles"].([]interface{}) {
			t := i.(map[string]interface{})["name"]
			if proName == t {
				fmt.Printf("name: %v\n", i.(map[string]interface{})["name"])
				fmt.Printf("endpoint: %v\n", i.(map[string]interface{})["endpoints"])
				fmt.Printf("load_generators %v\n", i.(map[string]interface{})["load_generators"])
				fmt.Printf("Test run duration %v\n", i.(map[string]interface{})["duration"])
				return nil
			}
			log.Fatalf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", proName)
		}
		return nil
	},
}
