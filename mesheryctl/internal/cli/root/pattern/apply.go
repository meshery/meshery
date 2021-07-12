package pattern

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply pattern file",
	Long:  `Apply pattern file will trigger deploy of the pattern file`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		client := &http.Client{}

		// set default tokenpath for perf command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		reqURL := mctlCfg.GetBaseMesheryURL() + "/api/experimental/pattern"

		// Method to check if the entered file is a URL or not
		if validURL := govalidator.IsURL(file); !validURL {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			text := string(content)
			jsonValues, err := json.Marshal(map[string]interface{}{
				"pattern_data": map[string]interface{}{
					"pattern_file": text,
				},
				"save": true,
			})
			if err != nil {
				return err
			}
			req, err = http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonValues))
			if err != nil {
				return err
			}
		} else {
			// add protocol if missing
			if !strings.Contains(file, "https") {
				file = "https://" + file
			}
			// change github url to raw.github url
			url, err := utils.ParseURLGithub(file)
			if err != nil {
				return err
			}

			log.Debug(url)

			jsonValues, err := json.Marshal(map[string]interface{}{
				"url":  url,
				"save": true,
			})
			if err != nil {
				return err
			}
			req, err = http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonValues))
			if err != nil {
				return err
			}
		}

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode == 200 {
			log.Info("pattern successfully applied")
		} else {
			log.Info(string(body))
		}
		return nil
	},
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
	_ = applyCmd.MarkFlagRequired("file")
}
