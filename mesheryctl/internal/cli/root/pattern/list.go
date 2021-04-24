package pattern

import (
	"encoding/json"
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
<<<<<<< HEAD
=======
	"net/http"
>>>>>>> 306424b1... list command v1
=======
	"io/ioutil"
	"net/http"
	"time"
>>>>>>> 2c41dd5c... wip list command
=======
=======
>>>>>>> ac869cc8... add patterns list command
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
<<<<<<< HEAD
>>>>>>> d1661977... add patterns list command
=======
	"net/http"
>>>>>>> 3927d252... list command v1
=======
>>>>>>> ac869cc8... add patterns list command
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
	token   string
	verbose bool
)

var listCmd = &cobra.Command{
	Use:  "list",
	Long: "Display list of all available pattern files",
	Args: cobra.MinimumNArgs(0),
<<<<<<< HEAD
=======
	token string
=======
	token   string
	allflag bool
>>>>>>> 2c41dd5c... wip list command
=======
	token string
>>>>>>> 3927d252... list command v1
=======
	token   string
	allflag bool
>>>>>>> ac869cc8... add patterns list command
)

var listCmd = &cobra.Command{
<<<<<<< HEAD
	Use:   "list",
	Short: "list pattern files",
<<<<<<< HEAD
<<<<<<< HEAD
	Long:  "List available pattern files",
	Args:  cobra.MinimumNArgs(1),
>>>>>>> 306424b1... list command v1
=======
	Long:  "Display list of all available pattern files",
	Args:  cobra.MinimumNArgs(0),
>>>>>>> 2c41dd5c... wip list command
=======
	Long:  "List available pattern files",
	Args:  cobra.MinimumNArgs(1),
>>>>>>> 3927d252... list command v1
=======
	Long: "Display list of all available pattern files",
	Args: cobra.MinimumNArgs(0),
>>>>>>> ac869cc8... add patterns list command
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
		var response models.PatternsAPIResponse

		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile", nil)
		if err != nil {
			return err
		}
		err = utils.AddAuthDetails(req, token)
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
		err = json.Unmarshal(body, &response)
<<<<<<< HEAD
=======
		var Patterns models.PatternsApiResponse
=======
=======
>>>>>>> ac869cc8... add patterns list command
		var response models.PatternsApiResponse
>>>>>>> 2c41dd5c... wip list command

		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile", nil)
		if err != nil {
			return err
		}
		err = utils.AddAuthDetails(req, token)
		if err != nil {
			return err
		}
<<<<<<< HEAD
		client.Do(req)
		err = json.NewDecoder(req.Body).Decode(&Patterns)
>>>>>>> 306424b1... list command v1
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
		if err != nil {
			return err
		}

<<<<<<< HEAD
<<<<<<< HEAD
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
		tokenObj, err := utils.ReadToken(token)
		if err != nil {
			return err
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		if verbose == true {
			if provider == "None" {
				for _, v := range response.Patterns {
					PatternID := v.ID.String()
					PatterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
				}
				utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
				return nil
			}
			for _, v := range response.Patterns {
				PatternID := utils.TruncateID(v.ID.String())
				var UserID string
				if v.UserID != nil {
					UserID = *v.UserID
				} else {
					UserID = "null"
				}
				PatterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternID, UserID, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data)

			return nil
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Patterns {
				PatterName := fmt.Sprintf("%s", strings.Trim(v.Name, filepath.Ext(v.Name)))
				PatternID := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternID, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Display full length user and pattern file identifiers")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "path to token")
	_ = listCmd.MarkFlagRequired("token")
}
<<<<<<< HEAD
=======
	},
}
>>>>>>> 306424b1... list command v1
=======
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		json.Unmarshal(body, &response)
		tokenObj, err := utils.ReadToken(token)
		if err != nil {
			return err
		}
		provider := tokenObj["meshery-provider"]
		var data [][]string

		if allflag == true {
			if provider == "None" {
				for _, v := range response.Patterns {
					PatternId := v.ID.String()
					PatterName := v.Name
					CreatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year(), v.CreatedAt.Hour(), v.CreatedAt.Minute(), v.CreatedAt.Second())
					UpdatedAt := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year(), v.UpdatedAt.Hour(), v.UpdatedAt.Minute(), v.UpdatedAt.Second())
					data = append(data, []string{PatternId, PatterName, CreatedAt, UpdatedAt})
				}
				utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
				return nil
			}

			for _, v := range response.Patterns {
				PatternId := utils.TruncateID(v.ID.String())
				var UserId string
				if v.UserID != nil {
					UserId = *v.UserID
				} else {
					UserId = "null"
				}
				PatterName := v.Name
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternId, UserId, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "USER ID", "NAME", "CREATED", "UPDATED"}, data)

			return nil
		}

		// Check if messhery provider is set
		if provider == "None" {
			for _, v := range response.Patterns {
				PatterName := fmt.Sprintf("%s", strings.Trim(v.Name, filepath.Ext(v.Name)))
				PatternId := utils.TruncateID(v.ID.String())
				CreatedAt := fmt.Sprintf("%d-%d-%d", int(v.CreatedAt.Month()), v.CreatedAt.Day(), v.CreatedAt.Year())
				UpdatedAt := fmt.Sprintf("%d-%d-%d", int(v.UpdatedAt.Month()), v.UpdatedAt.Day(), v.UpdatedAt.Year())
				data = append(data, []string{PatternId, PatterName, CreatedAt, UpdatedAt})
			}
			utils.PrintToTable([]string{"PATTERN ID", "NAME", "CREATED", "UPDATED"}, data)
		}

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&allflag, "all", "a", false, "Display full length user and pattern file identifiers")
	listCmd.Flags().StringVarP(&token, "token", "t", "", "path to token")
	listCmd.MarkFlagRequired("token")
}
<<<<<<< HEAD
>>>>>>> 2c41dd5c... wip list command
=======
		var Patterns models.PatternsApiResponse

		client := &http.Client{}
		req, err := http.NewRequest("GET", mctlCfg.GetBaseMesheryURL()+"/api/experimental/patternfile/", nil)
		if err != nil {
			return err
		}
		err = utils.AddAuthDetails(req, "./auth.json")
		if err != nil {
			return err
		}
		client.Do(req)
		err = json.NewDecoder(req.Body).Decode(&Patterns)
		if err != nil {
			return err
		}

	},
}
>>>>>>> 3927d252... list command v1
=======
>>>>>>> ac869cc8... add patterns list command
=======
>>>>>>> 9af444385e81347cabaf254dcad319009aef2028
