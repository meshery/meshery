package filter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag   bool
	outFormatFlag string
)

var viewCmd = &cobra.Command{
	Use:   "view [filter name]",
	Short: "Display filters(s)",
	Long:  `Displays the contents of a specific filter based on name or id`,
	Example: `
// View the specified WASM filter file
mesheryctl exp filter view [filter-name | ID]	

// View using filter name
mesheryctl exp filter view test-wasm
	`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		filter := ""
		isID := false
		// if filter name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New("-a cannot be used when [filter-name|filter-id] is specified")
			}
			filter = args[0]
			// check if the filter argument is a valid uuid v4 string
			isID, err = regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$", filter)
			if err != nil {
				return err
			}
		}
		url := mctlCfg.GetBaseMesheryURL()
		if len(filter) == 0 {
			if viewAllFlag {
				url += "/api/filter?page_size=10000"
			} else {
				return errors.New("[filter-name|filter-id] not specified, use -a to view all filters")
			}
		} else if isID {
			// if filter is a valid uuid, then directly fetch the filter
			url += "/api/filter/" + filter
		} else {
			// else search filter by name
			url += "/api/filter?search=" + filter
		}

		client := &http.Client{}
		req, err := utils.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			// failsafe for the case when a valid uuid v4 is not an id of any filter (bad api call)
			return errors.Errorf("Response Status Code %d, possible invalid ID", res.StatusCode)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		var dat map[string]interface{}
		if err = json.Unmarshal(body, &dat); err != nil {
			return errors.Wrap(err, "failed to unmarshal response body")
		}

		if isID {
			if body, err = json.MarshalIndent(dat, "", "  "); err != nil {
				return err
			}
		} else if viewAllFlag {
			// only keep the filter key from the response when viewing all the filters
			if body, err = json.MarshalIndent(map[string]interface{}{"filters": dat["filters"]}, "", "  "); err != nil {
				return err
			}
		} else {
			// use the first match from the result when searching by filter name
			arr := dat["filters"].([]interface{})
			if len(arr) == 0 {
				utils.Log.Info(fmt.Sprintf("filter with name: %s not found", filter))
				return nil
			}
			if body, err = json.MarshalIndent(arr[0], "", "  "); err != nil {
				return err
			}
		}

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return errors.Wrap(err, "failed to convert json to yaml")
			}
		} else if outFormatFlag != "json" {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}
		utils.Log.Info(string(body))
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all filters available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
