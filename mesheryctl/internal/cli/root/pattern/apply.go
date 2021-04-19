package pattern

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	//"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
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
		client := &http.Client{}
		var req *http.Request
		var err error
		if strings.Contains(file, "github.com") || strings.Contains(file, "raw.githubusercontent.com") {
			req, err = applyImportGitHubFileHandler(file)
		} else {
			req, err = applyLocalFileHandler(file)
		}

		if err != nil {
			return err
		}

		err = utils.AddAuthDetails(req, tokenPath)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		fmt.Println(res.StatusCode)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		log.Infof(string(body))

		return nil
	},
}

func applyImportGitHubFileHandler(file string) (*http.Request, error) {
	owner, repo, path, err := utils.ParseGitHubURL(file)
	if err != nil {
		return nil, err
	}
	reqURL, err := utils.ConstructURL("/api/experimental/patternfile/import/github.com", owner, repo)
	if err != nil {
		return nil, err
	}
	return http.NewRequest("GET", reqURL+"?path="+path, nil)
}

func applyLocalFileHandler(file string) (*http.Request, error) {
	fileReader, err := os.Open(file)
	if err != nil {
		return nil, errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
	}
	reqURL, err := utils.ConstructURL("/api/experimental/paternfile/deploy")
	if err != nil {
		return nil, err
	}
	return http.NewRequest("POST", reqURL, fileReader)
}
