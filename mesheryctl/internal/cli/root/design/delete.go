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

package design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete design file",
	Long: `delete design file will trigger deletion of the design file.
	When a positional target is passed, it is resolved with the following precedence:
	1. Existing local file path
	2. Remote URL
	3. Design ID or name lookup
	Find more information at: https://docs.meshery.io/reference/mesheryctl/design/delete`,
	Args: func(cmd *cobra.Command, args []string) error {
		if cmd.Flags().Changed("file") && strings.TrimSpace(file) == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("file path cannot be empty"))
		}
		if len(args) == 0 && file == "" {
			if err := cmd.Usage(); err != nil {
				return err
			}
			return utils.ErrInvalidArgument(fmt.Errorf("please provide a design name, ID, or file path"))
		}
		if len(args) > 0 && strings.TrimSpace(args[0]) == "" {
			return utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty"))
		}
		if len(args) > 1 {
			return utils.ErrInvalidArgument(fmt.Errorf("expected at most one positional argument, got %d", len(args)))
		}
		if len(args) > 0 && file != "" {
			return utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a positional argument and the --file flag"))
		}
		return nil
	},
	Example: `
// delete a design file
mesheryctl design delete [file | URL]

// delete a design by ID
mesheryctl design delete [ID]

// delete a design by name
mesheryctl design delete [name]

// delete a design using --file
mesheryctl design delete -f [file | URL]
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		var req *http.Request
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}

		targetFile := file
		pattern := ""
		isID := false
		if len(args) > 0 {
			target := args[0]
			if _, statErr := os.Stat(target); statErr == nil {
				targetFile = target
			} else if utils.IsValidUrl(target) {
				targetFile = target
			} else {
				pattern, isID, err = utils.ValidId(mctlCfg.GetBaseMesheryURL(), target, "pattern")
				if err != nil {
					if errors.GetCode(err) == utils.ErrNotFoundCode {
						return ErrDesignNotFound(target)
					}
					if errors.GetCode(err) == utils.ErrAmbiguousIDCode {
						return err
					}
					// If the server is unreachable or another network error occurs,
					// fall through to treat target as a file path so the descriptive
					// file-read error is surfaced instead of an obscure ID error.
					targetFile = target
				} else if !isID {
					resolvedID, nameErr := resolveDesignByName(mctlCfg.GetBaseMesheryURL(), target)
					if nameErr != nil {
						if errors.GetCode(nameErr) == ErrDesignsNotFoundCode || errors.GetCode(nameErr) == utils.ErrNotFoundCode {
							return ErrDesignNotFound(target)
						}
						if errors.GetCode(nameErr) == utils.ErrAmbiguousIDCode {
							return nameErr
						}
						// ValidId already reached this server moments ago, so a name
						// lookup failure is a real lookup error. Return it directly
						// instead of masking it with a misleading file-path error.
						return nameErr
					} else {
						pattern = resolvedID
						isID = true
					}
				}
			}
		}

		// Delete the pattern using the id
		if isID {
			err := utils.DeleteConfiguration(mctlCfg.GetBaseMesheryURL(), pattern, "pattern")
			if err != nil {
				if errors.GetCode(err) == utils.ErrNotFoundCode {
					return ErrDesignNotFound(args[0])
				}
				return ErrDeleteDesign(err, args[0])
			}
			utils.Log.Info("Design ", args[0], " deleted")
			return nil
		}
		deployURL := mctlCfg.GetBaseMesheryURL() + "/api/pattern/deploy"

		var patternFileByt []byte
		// If file path not a valid URL, treat it like a local file path
		if !utils.IsValidUrl(targetFile) {
			content, err := os.ReadFile(targetFile)
			if err != nil {
				return utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, targetFile))
			}

			patternFileByt = content
		} else {
			// Else treat it like a URL
			patternFileByt, err = fetchDesignFile(targetFile)
			if err != nil {
				return err
			}
		}

		// The server decodes this endpoint's body as a JSON
		// MesheryPatternFileDeployPayload, so send the design the same way
		// design deploy does rather than as raw file content.
		payloadBytes, err := json.Marshal(models.MesheryPatternFileDeployPayload{
			PatternFile: string(patternFileByt),
		})
		if err != nil {
			return utils.ErrMarshal(err)
		}

		req, err = utils.NewRequest("DELETE", deployURL, bytes.NewBuffer(payloadBytes))
		if err != nil {
			return utils.ErrCreatingRequest(err)
		}

		res, err := utils.MakeRequest(req)
		if err != nil {
			return utils.ErrRequestResponse(err)
		}

		defer func() { _ = res.Body.Close() }()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return utils.ErrReadFromBody(err)
		}

		utils.Log.Info(string(body))

		return nil
	},
}

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}

// resolveDesignByName queries the server for designs matching the provided name.
// It walks every page of results so an exact match on a later page is found and
// a second exact match on a later page still triggers the ambiguity guard.
// If exactly one exact name match is found, its ID is returned.
// If no design matches, ErrDesignNotFound is returned.
// If multiple designs share the exact same name, ErrAmbiguousID is returned.
func resolveDesignByName(baseURL, name string) (string, error) {
	const pageSize = 10000
	const maxPages = 50

	matchedIDs := []string{}
	seen := map[string]bool{}

	var totalCount int64
	var hasTotalCount bool
	var seenRows int64
	completed := false

	for page := 0; page < maxPages; page++ {
		searchURL := fmt.Sprintf(
			"%s/api/pattern?search=%s&pagesize=%d&page=%d",
			baseURL, url.QueryEscape(name), pageSize, page,
		)
		req, err := utils.NewRequest(http.MethodGet, searchURL, nil)
		if err != nil {
			return "", err
		}
		res, err := utils.MakeRequest(req)
		if err != nil {
			return "", err
		}

		body, err := io.ReadAll(res.Body)
		_ = res.Body.Close()
		if err != nil {
			return "", utils.ErrReadResponseBody(err)
		}

		// A non-200 response is a lookup failure, not a "not found" result.
		if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
			return "", utils.ErrRequestResponse(fmt.Errorf(
				"design name lookup failed with status %d: %s", res.StatusCode, string(body),
			))
		}

		var dat map[string]interface{}
		if err := json.Unmarshal(body, &dat); err != nil {
			return "", utils.ErrUnmarshal(err)
		}

		patterns, ok := dat["patterns"].([]interface{})
		if !ok || len(patterns) == 0 {
			if page == 0 {
				return "", ErrDesignNotFound(name)
			}
			completed = true
			break
		}

		seenRows += int64(len(patterns))
		for _, p := range patterns {
			pmap, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			// Keep the strict equality: substring matches do not count.
			if pmap["name"] == name {
				if id, ok := pmap["id"].(string); ok && id != "" && !seen[id] {
					seen[id] = true
					matchedIDs = append(matchedIDs, id)
				}
			}
		}

		if count, ok := parsePatternTotalCount(dat); ok {
			totalCount, hasTotalCount = count, true
		}
		if hasTotalCount {
			if seenRows >= totalCount || len(patterns) == 0 {
				completed = true
				break
			}
		} else if len(patterns) < pageSize {
			completed = true
			break
		}
	}

	// If the server reported a total but we did not see that many rows, the
	// lookup is incomplete: error rather than guess at uniqueness.
	if hasTotalCount && seenRows < totalCount {
		return "", utils.ErrIDLookupIncomplete(totalCount, seenRows)
	}
	if !completed {
		return "", utils.ErrIDLookupPageLimit(seenRows)
	}

	if len(matchedIDs) == 0 {
		return "", ErrDesignNotFound(name)
	}
	if len(matchedIDs) > 1 {
		return "", utils.ErrAmbiguousID(name, matchedIDs)
	}
	return matchedIDs[0], nil
}

// parsePatternTotalCount reads the total item count from a paginated search
// response, accepting both the current "totalCount" and legacy "total_count".
func parsePatternTotalCount(dat map[string]interface{}) (int64, bool) {
	for _, key := range []string{"totalCount", "total_count"} {
		if v, ok := dat[key]; ok && v != nil {
			switch val := v.(type) {
			case float64:
				return int64(val), true
			case int:
				return int64(val), true
			case int64:
				return val, true
			case string:
				if parsed, err := strconv.ParseInt(val, 10, 64); err == nil {
					return parsed, true
				}
			}
		}
	}
	return 0, false
}
