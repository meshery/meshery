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
	"strings"

	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	encoding "github.com/meshery/meshkit/encoding"
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
mesheryctl design delete [design-name]

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
					// Not a file, URL or ID: look the target up as a design name.
					pattern, err = findDesignIDByExactName(mctlCfg.GetBaseMesheryURL(), target)
					if err != nil {
						return err
					}
					isID = true
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
			patternFileByt, err = fetchDesignFromGithub(targetFile)
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

// findDesignIDByExactName returns the ID of the one design named exactly name.
// The server's search is a case-insensitive substring match, so results are
// filtered down to whole-name matches before anything is deleted: an exact
// match wins, and only without one does a match that differs in case count.
// Two designs that match equally well are reported as ambiguous instead of
// picking one.
func findDesignIDByExactName(baseURL, name string) (string, error) {
	searchURL := fmt.Sprintf("%s/api/pattern?search=%s&pagesize=10000", baseURL, url.QueryEscape(name))
	req, err := utils.NewRequest(http.MethodGet, searchURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := utils.MakeRequest(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to look up design %q by name: %w", name, utils.ErrResponseStatus(resp.StatusCode))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", utils.ErrReadFromBody(err)
	}
	var response models.PatternsAPIResponse
	if err := encoding.Unmarshal(body, &response); err != nil {
		return "", utils.ErrUnmarshal(err)
	}

	// The server's search ignores case, so fall back to a case-insensitive
	// match when no name matches exactly.
	var ids, foldIDs []string
	for _, p := range response.Patterns {
		switch {
		case p.Name == name:
			ids = append(ids, utils.FormatUUID(p.ID, ""))
		case strings.EqualFold(p.Name, name):
			foldIDs = append(foldIDs, utils.FormatUUID(p.ID, ""))
		}
	}
	if len(ids) == 0 {
		ids = foldIDs
	}
	switch len(ids) {
	case 0:
		return "", ErrDesignNotFound(name)
	case 1:
		return ids[0], nil
	default:
		return "", ErrAmbiguousDesignName(name, ids)
	}
}

func init() {
	deleteCmd.Flags().StringVarP(&file, "file", "f", "", "Path to design file")
}
