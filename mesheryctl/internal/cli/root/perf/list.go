package perf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"

	term "github.com/nsf/termbox-go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	page         uint
	totalPage    uint
	totalResults uint
	limitResults uint = 10
)

func reset() {
	term.Sync() // cosmestic purpose
}

var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List Performance profiles and Results of profiles",
	Long:    `List all the available performance profiles and results of a performance profile`,
	Args:    cobra.MaximumNArgs(1),
	Example: "mesheryctl perf list \nmesheryctl perf list [profile-id]",
	RunE: func(cmd *cobra.Command, args []string) error {
		page = 0 //Set the page number for API response to be zero

		// Get viper instance used for context
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return errors.Wrap(err, "error processing config")
		}
		if len(args) == 0 {
			// initialize termbox
			err := term.Init()
			if err != nil {
				panic(err)
			}
			// mainProfileLoop outputs profiles with pagination
		mainProfileLoop:
			for {
				data, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL() + "/api/user/performance/profiles")
				if err != nil {
					return err
				} else if len(data) > 0 {
					log.Info(fmt.Sprintf("Page %d out of %d | Total Results: %d", page, totalPage, totalResults))
					utils.PrintToTable([]string{"Name", "ID", "RESULTS", "LAST-RUN"}, data)
					if page == totalPage || len(data) == 0 {
						fmt.Printf("End of the results.\n")
						break mainProfileLoop
					}
				} else if len(data) == 0 {
					fmt.Printf("End of the results.")
					break mainProfileLoop
				}
				// askProfileLoop ask for keypress to output profiles on next page
				fmt.Println("Press Spacebar to advance, Ctrl+C to stop.")
			askProfileLoop:
				for {
					switch ev := term.PollEvent(); ev.Type {
					case term.EventKey:
						switch ev.Key {
						case term.KeySpace:
							break askProfileLoop
						case term.KeyCtrlC:
							break mainProfileLoop
						}
					}
				}
				reset() // make screen clear
			}
			// unhide cursor as termbox will hide it.
			// term.Close()
			fmt.Println("\x1b[?25h")
			return nil
		}
		// Output results of a performance profile
		profileID := args[0]
		// initialize termbox
		err = term.Init()
		if err != nil {
			panic(err)
		}
		// mainResultloop outputs results of a profile with pagination
	mainResultloop:
		for {
			data, err := fetchPerformanceProfileResults(mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID+"/results", profileID)
			if err != nil {
				return err
			} else if len(data) > 0 {
				log.Info(fmt.Sprintf("Page %d out of %d | Total Results: %d", page, totalResults/limitResults+1, totalResults))
				utils.PrintToTable([]string{"NAME", "MESH", "START-TIME", "QPS", "DURATION", "P50", "P99.9"}, data)
				if page == totalPage || len(data) == 0 {
					fmt.Printf("End of the results.")
					break
				}
			} else if len(data) == 0 {
				fmt.Printf("End of the results.")
				break
			}
			fmt.Println("Press Spacebar to advance, Ctrl+C to stop.")
			// askResultLoop ask for keypress to output results on next page
		askResultLoop:
			for {
				switch ev := term.PollEvent(); ev.Type {
				case term.EventKey:
					switch ev.Key {
					case term.KeySpace:
						break askResultLoop
					case term.KeyCtrlC:
						break mainResultloop
					}
				}
			}
			reset() //Clear screen
		}
		// unhide cursor as termbox will hide it.
		fmt.Println("\x1b[?25h")
		return nil
	},
}

// Fetch all the profiles
func fetchPerformanceProfiles(url string) ([][]string, error) {
	client := &http.Client{}
	var response *models.PerformanceProfilesAPIResponse
	tempURL := fmt.Sprintf("%s?page_size=%d&page=%d", url, limitResults, page)
	req, err := http.NewRequest("GET", tempURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, utils.PerfError("Failed to fetch performance results"))
	}
	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// failsafe for the case when a valid uuid v4 is not an id of any pattern (bad api call)
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Response Status Code %d, possible Server Error", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}
	var data [][]string

	for _, profile := range response.Profiles {
		lastRun := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(profile.LastRun.Month()), profile.LastRun.Day(), profile.LastRun.Year(), profile.LastRun.Hour(), profile.LastRun.Minute(), profile.LastRun.Second())
		data = append(data, []string{profile.Name, profile.ID.String(), fmt.Sprintf("%d", profile.TotalResults), lastRun})
	}

	//increase the page count and set totalPage and totalResults
	page++
	totalResults = response.TotalCount
	if response.TotalCount%limitResults == 0 {
		totalPage = response.TotalCount / limitResults
	} else {
		totalPage = response.TotalCount/limitResults + 1
	}

	return data, nil
}

// Fetch results for a specific profile
func fetchPerformanceProfileResults(url string, profileID string) ([][]string, error) {
	client := &http.Client{}
	var response *models.PerformanceResultsAPIResponse
	tempURL := fmt.Sprintf("%s?pageSize=%d&page=%d", url, limitResults, page)
	req, err := http.NewRequest("GET", tempURL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, utils.PerfError("Failed to fetch performance results"))
	}
	err = utils.AddAuthDetails(req, tokenPath)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("Performance profile `%s` not found. Please verify profile name and try again. Use `mesheryctl perf list` to see a list of performance profiles.", profileID)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal response body")
	}

	var data [][]string

	// append data for single profile
	for _, result := range response.Results {
		serviceMesh := "No Mesh"
		if result.Mesh != "" {
			serviceMesh = result.Mesh
		}
		startTime := fmt.Sprintf("%d-%d-%d %d:%d:%d", int(result.TestStartTime.Month()), result.TestStartTime.Day(), result.TestStartTime.Year(), result.TestStartTime.Hour(), result.TestStartTime.Minute(), result.TestStartTime.Second())
		p50 := result.RunnerResults.DurationHistogram.Percentiles[0].Value
		p99_9 := result.RunnerResults.DurationHistogram.Percentiles[len(result.RunnerResults.DurationHistogram.Percentiles)-1].Value
		timeDuration := strings.SplitAfterN(strconv.FormatUint(uint64(result.RunnerResults.Duration), 10), "", 5)
		duration := fmt.Sprintf("%s%s.%s%ss", timeDuration[0], timeDuration[1], timeDuration[2], timeDuration[3])
		data = append(data, []string{result.Name, serviceMesh, startTime, fmt.Sprintf("%f", result.RunnerResults.QPS), duration, fmt.Sprintf("%f", p50), fmt.Sprintf("%f", p99_9)})
	}
	//increase the page count and set totalPage and totalResults
	page++
	totalResults = response.TotalCount
	if response.TotalCount%limitResults == 0 {
		totalPage = response.TotalCount / limitResults
	} else {
		totalPage = response.TotalCount/limitResults + 1
	}
	return data, nil
}
