package perf

import (
	"fmt"
	"io"
	"net/http"

	"github.com/ghodss/yaml"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [profile-name]",
	Short: "Delete performance profile",
	Long:  `Delete a performance profile by name`,
	Args:  cobra.ExactArgs(1),
	Example: `
// Delete a performance profile by name
mesheryctl perf delete profile-name
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// setting up for error formatting
		cmdUsed = "delete"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return errors.Wrap(err, utils.PerfError("failed to get configuration"))
		}

		profileName := args[0]

		// First search for profiles matching the name
		profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), profileName, pageSize, 0)
		if err != nil {
			utils.Log.Error(err)
			return errors.Wrap(err, utils.PerfError("failed to search for performance profiles"))
		}

		if len(profiles) == 0 {
			utils.Log.Info("No Performance Profiles found with name: ", profileName)
			return nil
		}

		var profileID string
		var foundProfile bool

		// If multiple profiles match, find the exact match
		for _, profile := range profiles {
			if profile.Name == profileName {
				profileID = profile.ID.String()
				foundProfile = true
				break
			}
		}

		// If exact match not found, but we have results, prompt user to select
		if !foundProfile && len(profiles) > 0 {
			// Print available profiles
			data := profilesToStringArrays(profiles)
			utils.PrintToTable([]string{"Name", "ID", "RESULTS", "Load-Generator", "Last-Run"}, data)

			if !utils.SilentFlag {
				// Prompt user to select profile to delete
				index, err := userPrompt("profile", "Enter index of the profile to delete", data)
				if err != nil {
					return errors.Wrap(err, utils.PerfError("failed to get user selection"))
				}
				profileID = profiles[index].ID.String()
			} 
		}

		// Perform delete if we found a profile
		if profileID != "" {
			// Confirm deletion
			if !utils.SilentFlag {
				userResponse := utils.AskForConfirmation("Are you sure you want to delete this performance profile?")
				if !userResponse {
					return nil
				}
			}

			// Delete the profile
			response, err := deletePerformanceProfile(mctlCfg.GetBaseMesheryURL(), profileID)
			if err != nil {
				utils.Log.Error(err)
				return errors.Wrap(err, utils.PerfError("failed to delete performance profile"))
			}

			if outputFormatFlag != "" {
				if outputFormatFlag == "yaml" {
					body, err := yaml.JSONToYAML(response)
					if err != nil {
						return errors.Wrap(err, utils.PerfError("failed to parse response to yaml"))
					}
					utils.Log.Info(string(body))
				} else if outputFormatFlag == "json" {
					utils.Log.Info(string(response))
				} else {
					utils.Log.Error(ErrInvalidOutputChoice())
				}
			} else {
				utils.Log.Info("Performance profile deleted successfully")
			}
			return nil
		}

		utils.Log.Info("No Performance Profile found with name: ", profileName)
		return nil
	},
}

// deletePerformanceProfile deletes a performance profile with the given ID
func deletePerformanceProfile(baseURL, profileID string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/user/performance/profiles/%s", baseURL, profileID)

	req, err := utils.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, utils.PerfError("failed to read response body"))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to delete performance profile. Status code: %d, Response: %s", resp.StatusCode, string(body))
	}

	return body, nil
}
