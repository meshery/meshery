package app

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var offboardCmd = &cobra.Command{
	Use:   "offboard",
	Short: "Offboard application",
	Long:  `Offboard application will trigger deletion of the application file`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
	Offboard application by providing file path
	mesheryctl app offboard -f <filepath>
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		// set default tokenpath for app offboard command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.SystemError(fmt.Sprintf("failed to read file %s", file)))
		}

		client := &http.Client{}
		req, err := http.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/application/deploy", fileReader)
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

func init() {
	offboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to app file")
}
