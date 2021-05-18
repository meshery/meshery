package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outFormatFlag string

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "view perf profile",
	Long:    `See the configuration of your performance profile`,
	Example: "mesheryctl perf view [ performance test profile name ]",
	Args:    cobra.MaximumNArgs(1),
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
			map2 := make(map[string]interface{})
			for k, v := range i.(map[string]interface{}) {
				if k == "name" || k == "endpoints" || k == "qps" || k == "duration" || k == "load_generators" {
					map2[k] = v
				}
			}
			if outFormatFlag == "json" && proName == t {
				// create a second map to copy the informations we want to
				if data, err = json.MarshalIndent(map2, "", "  "); err != nil {
					return err
				}
				log.Info(string(data))
			} else if proName == t {
				fmt.Printf("name: %v\n", map2["name"])
				fmt.Printf("endpoint: %v\n", map2["endpoints"])
				fmt.Printf("load_generators %v\n", map2["load_generators"])
				fmt.Printf("Test run duration %v\n", map2["duration"])
			}
			if outFormatFlag != "json" && outFormatFlag != "" {
				return errors.New("output-format choice invalid, use json")
			}
		}
		return nil
	},
}

func init() {
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "", "(optional) format to display in json")
}
