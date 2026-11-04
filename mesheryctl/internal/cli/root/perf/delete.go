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
	"encoding/json"
	"net/url"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	perfprofile "github.com/meshery/schemas/models/v1beta3/performance_profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [profile-name]",
	Short: "Delete a performance profile",
	Long: `Delete a performance profile by its name.
Find more information at: https://docs.meshery.io/reference/mesheryctl/perf/delete`,
	Example: `
// Delete a performance profile (asks for confirmation)
mesheryctl perf delete test-3

// Delete a performance profile without confirmation
mesheryctl perf delete test-3 --yes
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return ErrNoProfileName()
		}
		if len(args) > 1 {
			return ErrorArgumentOverflow()
		}

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return utils.ErrLoadConfig(err)
		}
		baseURL := mctlCfg.GetBaseMesheryURL()

		profiles, err := fetchAllPerformanceProfiles(baseURL, url.QueryEscape(args[0]))
		if err != nil {
			return err
		}

		// The search is a substring match, so prefer profiles named exactly as asked.
		matches := profiles[:0:0]
		for _, p := range profiles {
			if p.Name == args[0] {
				matches = append(matches, p)
			}
		}
		// --yes only authorizes the profile named, never a similarly named one.
		if len(matches) == 0 && !utils.SilentFlag {
			matches = profiles
		}
		if len(matches) == 0 {
			return ErrNoProfileFound()
		}

		index := 0
		if len(matches) > 1 {
			index, err = userPrompt("profile", "Enter index of the profile", profilesToStringArrays(matches))
			if err != nil {
				return ErrUserPrompt(err)
			}
		}
		profile := matches[index]

		if !utils.SilentFlag && !utils.AskForConfirmation("Performance profile '"+profile.Name+"' ("+profile.ID.String()+") will be deleted. Are you sure you want to continue") {
			utils.Log.Info("Delete aborted.")
			return nil
		}

		req, err := utils.NewRequest("DELETE", baseURL+"/api/performance/profiles/"+profile.ID.String(), nil)
		if err != nil {
			return err
		}
		resp, err := utils.MakeRequest(req)
		if err != nil {
			return err
		}
		defer utils.SafeClose(resp.Body)

		utils.Log.Infof("Performance profile '%s' (%s) deleted", profile.Name, profile.ID.String())
		return nil
	},
}

// fetchAllPerformanceProfiles returns every profile matching searchString, across all pages.
func fetchAllPerformanceProfiles(baseURL, searchString string) ([]perfprofile.PerformanceProfile, error) {
	var all []perfprofile.PerformanceProfile
	for page := 0; ; page++ {
		profiles, body, err := fetchPerformanceProfiles(baseURL, searchString, pageSize, page)
		if err != nil {
			return nil, err
		}
		all = append(all, profiles...)

		var res perfprofile.PerformanceProfilePage
		if err := json.Unmarshal(body, &res); err != nil {
			return nil, ErrFailUnmarshal(err)
		}
		if len(profiles) == 0 || len(all) >= res.TotalCount {
			return all, nil
		}
	}
}
