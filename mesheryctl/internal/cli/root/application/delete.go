package application

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete application file",
	Long:  `delete application file will trigger deletion of the application file`,
	Args:  cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
		}

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/experimental/application/deploy", fileReader)
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

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		log.Infof(string(body))

		return nil
	},
}
