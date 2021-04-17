package pattern

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
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
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		client := &http.Client{}
		var req *http.Request
		if strings.Contains(file, "github.com") || strings.Contains(file, "raw.githubusercontent.com") {
			owner, repo, path, err := parseGitHubURL(file)
			if err != nil {
				return err
			}
			reqURL := fmt.Sprintf("%s/api/experimental/patternfile/import/github.com/%s/%s?path=%s", mctlCfg.GetBaseMesheryURL(), owner, repo, path)
			req, err = http.NewRequest("GET", reqURL, nil)
			if err != nil {
				return err
			}
		} else {
			fileReader, err := os.Open(file)
			if err != nil {
				return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
			}
			req, err = http.NewRequest("POST", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile/deploy", fileReader)
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

		log.Infof(string(body))

		return nil
	},
}

func parseGitHubURL(url string) (string, string, string, error) {
	ogPath := strings.Replace(url, "https://", "", 1)
	idx := strings.Index(ogPath, "/")
	if idx == -1 {
		return "", "", "", errors.New("invalid GitHub URL")
	}
	host := ogPath[:idx]
	paths := strings.Split(ogPath, "/")
	if (host == "github.com" && len(paths) < 6) || (host == "raw.githubcontent.com" && len(paths) < 5) {
		return "", "", "", errors.New(fmt.Sprintf("failed to retrieve file from URL: %s", url))
	}
	if host == "github.com" {
		return paths[1], paths[2], strings.Join(paths[5:], "/"), nil
	} else if host == "raw.githubusercontent.com" {
		return paths[1], paths[2], strings.Join(paths[4:], "/"), nil
	} else {
		return "", "", "", errors.New("invalid GitHub URL")
	}
}
