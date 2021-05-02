package pattern

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	//"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	//"github.com/spf13/viper"
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
		url, path, err := utils.ParseURL(file)
		if err != nil {
			return err
		}

		reqURL, err := utils.ConstructURL("/api/experimental/pattern")
		if err != nil {
			return err
		}
		if path == "" && url == file {
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
		} else {
			jsonValues, err := json.Marshal(map[string]interface{}{
				"url":  url,
				"path": path,
				"save": true,
			})
			if err != nil {
				return err
			}
			req, err = http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonValues))
		}

		if err != nil {
			return err
		}

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		fmt.Println(res)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		log.Info(string(body))
		return nil
	},
}
