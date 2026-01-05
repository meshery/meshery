// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package perf

import (
	"fmt"
	"io"
	"net/http"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	deleteAllProfiles bool
)

var linkDocPerfDelete = map[string]string{
	"link":    "![perf-delete-usage](/assets/img/mesheryctl/perf-delete.png)",
	"caption": "Usage of mesheryctl perf delete",
}

var deleteCmd = &cobra.Command{
	Use:   "delete [profile-name]",
	Short: "Delete a Performance profile",
	Long: `Delete Performance profiles by name or delete all profiles
	Documentation for components can be found at https://docs.meshery.io/reference/mesheryctl/perf/delete`,
	Args: cobra.MaximumNArgs(1),
	Example: `
// Delete a specific performance profile
mesheryctl perf delete any-profile

// Delete multiple profiles matching a pattern
mesheryctl perf delete some-profile-XXXXX

// Delete all performance profiles
mesheryctl perf delete --all
	`,
	Annotations: linkDocPerfDelete,
	RunE: func(cmd *cobra.Command, args []string) error {

		// setting up for error formatting
		cmdUsed = "delete"

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		if deleteAllProfiles {
			return deleteAllPerformanceProfiles(mctlCfg)
		}

		if len(args) == 0 {
			return ErrNoProfileName()
		}

		return deletePerformanceProfile(mctlCfg, profileName)
	},
}

func deletePerformanceProfile(mctlCfg *config.MesheryCtlConfig, profileName string) error {
	log.Debug("Fetching performance profile for deletion")

	// Fetch profiles matching the name
	profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), profileName, pageSize, pageNumber-1)
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		return ErrNoProfileFound()
	}

	var profilesToDelete []string
	var profileNamesToDelete []string

	if len(profiles) == 1 {
		// Single profile found
		profilesToDelete = append(profilesToDelete, profiles[0].ID.String())
		profileNamesToDelete = append(profileNamesToDelete, profiles[0].Name)
	} else {
		// Multiple profiles found - ask user to select
		data := profilesToStringArrays(profiles)

		if utils.SilentFlag {
			// If force delete, delete all matching profiles
			for _, profile := range profiles {
				profilesToDelete = append(profilesToDelete, profile.ID.String())
				profileNamesToDelete = append(profileNamesToDelete, profile.Name)
			}
		} else {
			index, err := userPrompt("profile", "Enter index of the profile to delete (or 'all' for all matching profiles)", data)
			if err != nil {
				return ErrorArgumentOverflow()
			}

			if index == -1 { // User selected 'all'
				for _, profile := range profiles {
					profilesToDelete = append(profilesToDelete, profile.ID.String())
					profileNamesToDelete = append(profileNamesToDelete, profile.Name)
				}
			} else {
				profilesToDelete = append(profilesToDelete, profiles[index].ID.String())
				profileNamesToDelete = append(profileNamesToDelete, profiles[index].Name)
			}
		}
	}

	// Confirm deletion unless silent flag is used
	if !utils.SilentFlag {
		var confirmMsg string
		if len(profilesToDelete) == 1 {
			confirmMsg = fmt.Sprintf("Are you sure you want to delete profile '%s'?", profileNamesToDelete[0])
		} else {
			confirmMsg = fmt.Sprintf("Are you sure you want to delete %d profiles?", len(profilesToDelete))
		}

		if !utils.AskForConfirmation(confirmMsg) {
			log.Info("Profile deletion cancelled")
			return nil
		}
	}

	// Delete profiles
	successCount := 0
	for i, profileID := range profilesToDelete {
		err := deleteProfileByID(mctlCfg, profileID, profileNamesToDelete[i])
		if err != nil {
			utils.Log.Error(err)
		} else {
			successCount++
		}
	}

	if successCount > 0 {
		if successCount == 1 {
			log.Info("Performance profile deleted successfully!")
		} else {
			log.Infof("%d performance profiles deleted successfully!", successCount)
		}
	}

	return nil
}

func deleteAllPerformanceProfiles(mctlCfg *config.MesheryCtlConfig) error {
	log.Debug("Fetching all performance profiles for deletion")

	// Fetch all profiles
	profiles, _, err := fetchPerformanceProfiles(mctlCfg.GetBaseMesheryURL(), "", 0, 0)
	if err != nil {
		return err
	}

	if len(profiles) == 0 {
		log.Info("No performance profiles found to delete")
		return nil
	}

	// Confirm deletion unless force flag is used
	if !utils.SilentFlag {
		confirmMsg := fmt.Sprintf("Are you sure you want to delete all %d performance profiles? This action cannot be undone.", len(profiles))
		if !utils.AskForConfirmation(confirmMsg) {
			log.Info("Profile deletion cancelled")
			return nil
		}
	}

	// Delete all profiles
	successCount := 0
	for _, profile := range profiles {
		err := deleteProfileByID(mctlCfg, profile.ID.String(), profile.Name)
		if err != nil {
			fmt.Printf("Failed to delete profile '%s': %v", profile.Name, err)
			utils.Log.Error(err)
		} else {
			successCount++
		}
	}

	if successCount > 0 {
		log.Infof("%d performance profiles deleted successfully!", successCount)
	}

	return nil
}

func deleteProfileByID(mctlCfg *config.MesheryCtlConfig, profileID, profileName string) error {
	log.Debugf("Deleting performance profile: %s (ID: %s)", profileName, profileID)

	req, err := utils.NewRequest("DELETE", mctlCfg.GetBaseMesheryURL()+"/api/user/performance/profiles/"+profileID, nil)
	if err != nil {
		return utils.ErrCreatingRequest(err)
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return err
	}

	defer utils.SafeClose(resp.Body)

	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return ErrFailUnmarshal(err)
		}
		return ErrPerformanceProfileDelete(fmt.Errorf("%s", data))
	}

	log.Debugf("Profile '%s' deleted successfully", profileName)
	return nil
}

func init() {
	deleteCmd.Flags().BoolVar(&deleteAllProfiles, "all", false, "(optional) Delete all performance profiles")
}
