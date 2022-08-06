package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var offboardCmd = &cobra.Command{
	Use:   "offboard",
	Short: "Offboard application",
	Long:  `Offboard application will trigger undeploy of application`,
	Args:  cobra.MinimumNArgs(0),
	Example: `
// Offboard application by providing file path
mesheryctl app offboard -f [filepath]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			const errMsg=`Usage: mesheryctl app offboard -f [filepath]`
			return fmt.Errorf("no file path provided \n\n%v", errMsg)
		} 
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/application/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern"

		// Read file
		fileReader, err := os.Open(file)
		if err != nil {
			return errors.New(utils.AppError(fmt.Sprintf("failed to read file %s", file)))
		}

		client := &http.Client{}
		req, err := utils.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/application/deploy", fileReader)
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode == 200 {
			utils.Log.Info("app successfully offboarded")
		}
		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	offboardCmd.Flags().StringVarP(&file, "file", "f", "", "Path to app file")
}
