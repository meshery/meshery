package pattern

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/constants"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	skipSave    bool // skip saving a pattern
	patternFile string
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply pattern file",
	Long:  `Apply pattern file will trigger deploy of the pattern file`,
	Example: `
	// apply a pattern file
	mesheryctl pattern apply -f <file | URL>

	// deploy a saved pattern
	mesheryctl pattern apply <pattern-name>
	`,
	Args: cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		var err error
		client := &http.Client{}

		// set default tokenpath for perf apply command.
		if tokenPath == "" {
			tokenPath = constants.GetCurrentAuthToken()
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}

		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/experimental/pattern/deploy"
		patternURL := mctlCfg.GetBaseMesheryURL() + "/api/experimental/pattern"

		// pattern name has been passed
		if len(args) > 0 {
			// Merge args to get pattern-name
			patternName := strings.Join(args, "%20")

			// search and fetch patterns with pattern-name
			log.Debug("Fetching patterns")

			req, err = http.NewRequest("GET", patternURL+"?search="+patternName, nil)
			if err != nil {
				return err
			}

			err = utils.AddAuthDetails(req, tokenPath)
			if err != nil {
				return errors.New("authentication token not found. please supply a valid user token with the --token (or -t) flag")
			}

			resp, err := client.Do(req)
			if err != nil {
				return err
			}

			var response *models.PatternsAPIResponse
			// failsafe (bad api call)
			if resp.StatusCode != 200 {
				return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return errors.Wrap(err, utils.PerfError("failed to read response body"))
			}
			err = json.Unmarshal(body, &response)
			if err != nil {
				return errors.Wrap(err, "failed to unmarshal response body")
			}

			index := 0
			if len(response.Patterns) == 0 {
				return errors.New("no patterns found with the given name")
			} else if len(response.Patterns) == 1 {
				patternFile = response.Patterns[0].PatternFile
			} else {
				// Multiple patterns with same name
				index = multiplePatternsConfirmation(response.Patterns)
				patternFile = response.Patterns[index].PatternFile
			}
		} else {
			// Method to check if the entered file is a URL or not
			if validURL := govalidator.IsURL(file); !validURL {
				content, err := ioutil.ReadFile(file)
				if err != nil {
					return err
				}
				text := string(content)

				// if --skip-save is not passed we save the pattern first
				if !skipSave {
					jsonValues, err := json.Marshal(map[string]interface{}{
						"pattern_data": map[string]interface{}{
							"pattern_file": text,
						},
						"save": true,
					})
					if err != nil {
						return err
					}
					req, err = http.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
					if err != nil {
						return err
					}
					err = utils.AddAuthDetails(req, tokenPath)
					if err != nil {
						return err
					}
					resp, err := client.Do(req)
					if err != nil {
						return err
					}
					log.Debug("saved pattern file")
					var response *models.MesheryPattern
					// failsafe (bad api call)
					if resp.StatusCode != 200 {
						return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
					}
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return errors.Wrap(err, utils.PerfError("failed to read response body"))
					}
					err = json.Unmarshal(body, &response)
					if err != nil {
						return errors.Wrap(err, "failed to unmarshal response body")
					}
				}

				// setup pattern file
				patternFile = text
			} else {
				var jsonValues []byte
				url, err := utils.ParseURLGithub(file)
				if err != nil {
					return err
				}
				// save the pattern with Github URL
				if !skipSave {
					jsonValues, _ = json.Marshal(map[string]interface{}{
						"url":  url,
						"save": true,
					})
				} else { // we don't save the pattern
					jsonValues, _ = json.Marshal(map[string]interface{}{
						"url":  url,
						"save": false,
					})
				}
				req, err = http.NewRequest("POST", patternURL, bytes.NewBuffer(jsonValues))
				if err != nil {
					return err
				}
				err = utils.AddAuthDetails(req, tokenPath)
				if err != nil {
					return err
				}
				resp, err := client.Do(req)
				if err != nil {
					return err
				}
				log.Debug("remote hosted pattern request success")
				var response []*models.MesheryPattern
				// failsafe (bad api call)
				if resp.StatusCode != 200 {
					return errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return errors.Wrap(err, utils.PerfError("failed to read response body"))
				}
				err = json.Unmarshal(body, &response)
				if err != nil {
					return errors.Wrap(err, "failed to unmarshal response body")
				}

				// setup pattern file here
				patternFile = response[0].PatternFile
			}
		}

		req, err = http.NewRequest("POST", deployURL, bytes.NewBuffer([]byte(patternFile)))
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

		if res.StatusCode == 200 {
			log.Info("pattern successfully applied")
		} else {
			log.Info(string(body))
		}
		return nil
	},
}

func multiplePatternsConfirmation(profiles []models.MesheryPattern) int {
	reader := bufio.NewReader(os.Stdin)

	for index, a := range profiles {
		fmt.Printf("Index: %v\n", index)
		fmt.Printf("Name: %v\n", a.Name)
		fmt.Printf("ID: %s\n", a.ID.String())
		fmt.Printf("PatternFile:\n")
		fmt.Printf(a.PatternFile)
		fmt.Println("---------------------")
	}

	for {
		fmt.Printf("Enter the index of profile: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		response = strings.ToLower(strings.TrimSpace(response))
		index, err := strconv.Atoi(response)
		if err != nil {
			log.Info(err)
		}
		if index < 0 || index >= len(profiles) {
			log.Info("Invalid index")
		} else {
			return index
		}
	}
}

func init() {
	applyCmd.Flags().StringVarP(&file, "file", "f", "", "Path to pattern file")
	applyCmd.Flags().BoolVarP(&skipSave, "skip-save", "", false, "Skip saving a pattern")
}
