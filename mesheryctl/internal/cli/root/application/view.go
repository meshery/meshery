package application

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	viewAllFlag   bool
	outFormatFlag string
)

var viewCmd = &cobra.Command{
	Use:   "view <application name>",
	Short: "Display application(s)",
	Long:  `Displays the contents of a specific application based on name or id`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		application := ""
		isID := false
		// if application name/id available
		if len(args) > 0 {
			if viewAllFlag {
				return errors.New("-a cannot be used when [application-name|application-id] is specified")
			}
			application = args[0]
			// check if the application argument is a valid uuid v4 string
			isID, err = regexp.MatchString("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$", application)
			if err != nil {
				return err
			}
		}
		url := mctlCfg.GetBaseMesheryURL()
		if len(application) == 0 {
			if viewAllFlag {
				url += "/api/experimental/application?page_size=10000"
			} else {
				return errors.New("[application-name|application-id] not specified, use -a to view all applications")
			}
		} else if isID {
			// if application is a valid uuid, then directly fetch the application
			url += "/api/experimental/application/" + application
		} else {
			// else search application by name
			url += "/api/experimental/application?search=" + application
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}

		err = utils.AddAuthDetails(req, tokenPath)
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
		body, err := ioutil.ReadAll(res.Body)
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
			// only keep the application key from the response when viewing all the applications
			if body, err = json.MarshalIndent(map[string]interface{}{"applications": dat["applications"]}, "", "  "); err != nil {
				return err
			}
		} else {
			// use the first match from the result when searching by application name
			arr := dat["applications"].([]interface{})
			if len(arr) == 0 {
				log.Infof("application with name: %s not found", application)
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
		log.Info(string(body))
		return nil
	},
}

func init() {
	viewCmd.Flags().BoolVarP(&viewAllFlag, "all", "a", false, "(optional) view all applications available")
	viewCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
}
