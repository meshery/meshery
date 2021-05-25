package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var outFormatFlag string

type result struct {
	Name          string
	Endpoint      string
	QPS           int
	Duration      string
	Loadgenerator string
}

var viewCmd = &cobra.Command{
	Use:     "view",
	Short:   "view perf profile",
	Long:    `See the configuration of your performance profile`,
	Example: "mesheryctl perf view [ performance test profile name ]",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		proName := args[0]
		var req *http.Request
		url := mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles?search=" + proName
		var response *models.PerformanceProfilePage
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
			return errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", proName)
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, utils.PerfError("failed to read response body"))
		}

		if err = json.Unmarshal(data, &response); err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		var a result
		if response.TotalCount == 0 {
			return errors.New("profile does not exit. Please get a profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles")
		}
		for _, profile := range response.Profiles {
			if response.Profiles == nil {
				return errors.New("profile name not provide. Please get a profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles")
			}
			a = result{
				Name:          profile.Name,
				Endpoint:      profile.Endpoints[0],
				QPS:           profile.QPS,
				Duration:      profile.Duration,
				Loadgenerator: profile.LoadGenerators[0],
			}
			if outFormatFlag == "json" {
				if data, err = json.MarshalIndent(&a, "", "  "); err != nil {
					return err
				}
				fmt.Println(string(data))
			} else if outFormatFlag == "" {
				fmt.Printf("name: %v\n", a.Name)
				fmt.Printf("endpoint: %v\n", a.Endpoint)
				fmt.Printf("load_generators %v\n", a.Loadgenerator)
				fmt.Printf("Test run duration %v\n", a.Duration)
				fmt.Println("#####################")
			} else {
				return errors.New("output-format choice invalid, use json")
			}

		}
		return nil
	},
}

func init() {
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "", "(optional) format to display in json")
	viewCmd.Flags().StringVarP(&tokenPath, "token", "t", utils.AuthConfigFile, "(required) Path to meshery auth config")
	_ = viewCmd.MarkFlagRequired("token")
}
