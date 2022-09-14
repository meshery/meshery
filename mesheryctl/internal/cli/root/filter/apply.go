package filter

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipSave   bool // skip saving a filter
	filterFile string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply filter file",
	Long:  `Apply filter file will trigger deploy of the filter file`,
	Example: `
// Apply WASM filter file (login required)
mesheryctl exp filter apply --file [GitHub Link]

// Apply the file
mesheryctl exp filter apply --file https://github.com/layer5io/wasm-filters/tree/master/http-auth
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		client := &http.Client{}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/experimental/filter/deploy"
		filterURL := mctlCfg.GetBaseMesheryURL() + "/api/experimental/filter"

		// filter name has been passed
		if len(args) > 0 {
			// Merge args to get filter-name
			filterName := strings.Join(args, "%20")

			// search and fetch filters with filter-name
			log.Debug("Fetching filters")

			req, err = utils.NewRequest("GET", filterURL+"?search="+filterName, nil)
			if err != nil {
				return err
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			var response *models.FiltersAPIResponse
			// failsafe (bad api call)
			if resp.StatusCode != 200 {
				return ErrInvalidAPICall(resp.StatusCode)
			}
			defer resp.Body.Close()
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return ErrReadAPIResponse(err)
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return ErrUnmarshal(err)
			}

			index := 0
			if len(response.Filters) == 0 {
				return errors.New("no filters found with the given name")
			} else if len(response.Filters) == 1 {
				filterFile = response.Filters[0].FilterFile
			} else {
				// Multiple filters with same name
				index = multipleFiltersConfirmation(response.Filters)
				filterFile = response.Filters[index].FilterFile
			}
		} else {
			// Method to check if the entered file is a URL or not
			if validURL := govalidator.IsURL(file); !validURL {
				content, err := os.ReadFile(file)
				if err != nil {
					return err
				}
				text := string(content)

				// if --skip-save is not passed we save the filters first
				if !skipSave {
					jsonValues, err := json.Marshal(map[string]interface{}{
						"filter_data": map[string]interface{}{
							"filter_file": text,
						},
						"save": true,
					})
					if err != nil {
						return err
					}
					req, err = utils.NewRequest("POST", filterURL, bytes.NewBuffer(jsonValues))
					if err != nil {
						return err
					}
					resp, err := client.Do(req)
					if err != nil {
						return err
					}
					log.Debug("saved filter file")
					var response []*models.MesheryApplication
					// failsafe (bad api call)
					if resp.StatusCode != 200 {
						return ErrInvalidAPICall(resp.StatusCode)
					}
					defer resp.Body.Close()

					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return ErrReadAPIResponse(err)
					}
					err = json.Unmarshal(body, &response)
					if err != nil {
						return ErrUnmarshal(err)
					}
				}

				// setup filter file
				filterFile = text
			} else {
				var jsonValues []byte
				url, path, err := utils.ParseURLGithub(file)
				if err != nil {
					return err
				}

				log.Debug(url)
				log.Debug(path)

				// save the filter with Github URL
				if !skipSave {
					if path != "" {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"path": path,
							"save": true,
						})
					} else {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"save": true,
						})
					}
				} else { // we don't save the filter
					if path != "" {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"path": path,
							"save": false,
						})
					} else {
						jsonValues, _ = json.Marshal(map[string]interface{}{
							"url":  url,
							"save": false,
						})
					}
				}
				req, err = utils.NewRequest("POST", filterURL, bytes.NewBuffer(jsonValues))
				if err != nil {
					return err
				}

				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				log.Debug("remote hosted filter request success")
				var response []*models.MesheryFilter
				// failsafe (bad api call)
				if resp.StatusCode != 200 {
					return ErrInvalidAPICall(resp.StatusCode)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return ErrReadAPIResponse(err)
				}
				err = json.Unmarshal(body, &response)
				if err != nil {
					return ErrUnmarshal(err)
				}

				// setup filter file here
				filterFile = response[0].FilterFile
			}
		}

		req, err = utils.NewRequest("POST", deployURL, bytes.NewBuffer([]byte(filterFile)))
		if err != nil {
			return err
		}

		res, err := client.Do(req)
		if err != nil {
			return err
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		if res.StatusCode == 200 {
			log.Info("filter successfully deployed")
		}
		log.Info(string(body))
		return nil
	},
}

func multipleFiltersConfirmation(profiles []models.MesheryFilter) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("FilterFile:\n")
		fmt.Printf(a.FilterFile)
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of filter: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		index, err := strconv.Atoi(response)
		if err != nil {
			utils.Log.Info(err)
		}
		if index < 0 || index >= len(profiles) {
			utils.Log.Info("Invalid index")
		} else {
			return index
		}
	}
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "Path to filter file")
	applyCmd.Flags().BoolVarP(&skipSave, "skip-save", "", false, "Skip saving a filter")
}
