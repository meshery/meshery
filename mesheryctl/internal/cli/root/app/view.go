package app

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag   bool
	outFormatFlag string
)

var viewCmd = &cobra.Command{
	Use:   "view application name",
	Short: "Display application(s)",
	Long:  `Displays the contents of a specific application based on name or id`,
	Example: `
// View applictaions with name
mesheryctl app view [app-name]

// View applications with id
mesheryctl app view [app-id]

// View all applications
mesheryctl app view --all

! Refer below image link for usage
* Usage of mesheryctl app view
# ![app-view-usage](/assets/img/mesheryctl/app-view.png)
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		application := ""
		isID := false
		applicationID := ""
		// if application name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New("-a cannot be used when [application-name|application-id] is specified")
			}
			applicationID, isID, err = utils.Valid(args[0], "application")
			if err != nil {
				return err
			}
		}
		var req *http.Request
		url := mctlCfg.GetBaseMesheryURL()
		var response *models.ApplicationsAPIResponse
		// Merge args to get app-name
		application = strings.Join(args, "%20")
		if len(application) == 0 {
			if viewAllFlag {
				url += "/api/application?page_size=10000"
			} else {
				return errors.New("[application-name|application-id] not specified, use -a to view all applications")
			}
		} else if isID {
			// if application is a valid uuid, then directly fetch the application
			url += "/api/application/" + applicationID
		} else {
			// else search application by name
			url += "/api/application?search=" + application
		}

		client := &http.Client{}
		req, err = utils.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}
		if res.StatusCode != 200 {
			// failsafe for the case when a valid uuid v4 is not an id of any application (bad api call)
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
			if body, err = json.MarshalIndent(map[string]interface{}{"applications": dat["applications"]}, "", "  "); err != nil {
				return err
			}
		} else {
			if err = json.Unmarshal(body, &response); err != nil {
				return errors.Wrap(err, "failed to unmarshal response body")
			}
			if response.TotalCount == 0 {
				return errors.New("application does not exit. Please get an app name and try again. Use `mesheryctl app list` to see a list of applications")
			}
			// Manage more than one apps with similar name
			for _, app := range response.Applications {
				if response.Applications == nil {
					return errors.New("application name not provide. Please get an app name and try again. Use `mesheryctl app list` to see a list of applications")
				}
				body, err = json.MarshalIndent(&app, "", "  ")
				if err != nil {
					return err
				}
				if outFormatFlag == "json" {
					utils.Log.Info(string(body))
					continue
				}
				if outFormatFlag == "yaml" {
					if body, err = yaml.JSONToYAML(body); err != nil {
						return errors.Wrap(err, "failed to convert json to yaml")
					}
					utils.Log.Info(string(body))
					continue
				}
			}
		}

		if outFormatFlag == "yaml" {
			if body, err = yaml.JSONToYAML(body); err != nil {
				return errors.Wrap(err, "failed to convert json to yaml")
			}
		} else if outFormatFlag != "json" {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}
		if viewAllFlag || isID {
			utils.Log.Info(string(body))
		}
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all applications available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
