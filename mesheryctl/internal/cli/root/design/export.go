// Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package design

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	encoding "github.com/layer5io/meshkit/encoding"
	meshkitutils "github.com/layer5io/meshkit/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	designType string
	outputDir  string
)

type UserProfile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var exportCmd = &cobra.Command{
	Use:   "export [pattern-name | ID]",
	Short: "Export a design from Meshery",
	Long: `The 'export' command allows you to export a specific design from your Meshery server.
You can specify the design by its name or ID and optionally define the type of design.
The command also supports specifying an output directory where the exported design will be saved.
By default, the exported design will be saved in the current directory. The different types of design
type allowed are oci, original, and current. The default design type is current.`,
	Example: `
	# Export a design with a specific ID
	mesheryctl design export [pattern-name | ID]
	
	# Export a design with a specific ID and type
	mesheryctl design export [pattern-name | ID] --type [design-type]
	
	# Export a design and save it to a specific directory
	mesheryctl design export [pattern-name | ID] --output ./designs
	
	# Export a design with a specific type and save it to a directory
	mesheryctl design export [pattern-name | ID] --type [design-type] --output ./exports
	`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		patternNameOrID := strings.Join(args, " ")
		design, isID, err := utils.ValidId(mctlCfg.GetBaseMesheryURL(), patternNameOrID, "pattern")
		if err != nil {
			utils.Log.Error(err)
			return nil
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		if !isID {
			if design, err = fetchPatternIDByName(baseUrl, design); err != nil {
				utils.Log.Error(err)
				return nil
			}
		}

		designType, _ := cmd.Flags().GetString("type")
		if designType == "" {
			designType = "current"
		}

		if err := exportDesign(baseUrl, design, designType); err != nil {
			utils.Log.Error(err)
			return nil
		}

		return nil
	},
}

func fetchPatternIDByName(baseUrl, patternName string) (string, error) {
	patternUrl := fmt.Sprintf("%s/api/pattern?search=%s", baseUrl, url.QueryEscape(patternName))
	req, err := utils.NewRequest(http.MethodGet, patternUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := utils.MakeRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", models.ErrDoRequest(err, resp.Request.Method, patternUrl)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrReadFromBody(err)
	}
	var response struct {
		TotalCount int                     `json:"total_count"`
		Patterns   []models.MesheryPattern `json:"patterns"`
	}
	if err := encoding.Unmarshal(buf, &response); err != nil {
		return "", err
	}

	if response.TotalCount == 0 {
		return "", ErrDesignNotFound()
	} else if response.TotalCount == 1 {
		return response.Patterns[0].ID.String(), nil
	}

	selectedPattern := selectPatternPrompt(response.Patterns, baseUrl)
	return selectedPattern.ID.String(), nil
}

func exportDesign(baseUrl, design, designType string) error {
	dataURL := fmt.Sprintf("%s/api/pattern/%s", baseUrl, design)
	pattern, err := fetchPatternData(dataURL)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/pattern/download/%s", baseUrl, design)
	switch designType {
	case "oci":
		url += "?oci=true"
	case "original":
		url += fmt.Sprintf("/%s", pattern.Type.String)
	}

	resp, err := makeRequest(http.MethodGet, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.ErrDoRequest(err, resp.Request.Method, url)
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return ErrReadFromBody(err)
	}

	filename := generateFilename(pattern.Name, design, designType)
	outputFilePath := filepath.Join(outputDir, filename)
	outputFilePath = getUniqueFilename(outputFilePath)

	if err = os.MkdirAll(filepath.Dir(outputFilePath), 0755); err != nil {
		return models.ErrMakeDir(err, outputFilePath)
	}

	return writeToFile(outputFilePath, buf)
}

func fetchPatternData(dataURL string) (*models.MesheryPattern, error) {
	resp, err := makeRequest(http.MethodGet, dataURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, models.ErrDoRequest(err, resp.Request.Method, dataURL)
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, ErrReadFromBody(err)
	}

	var pattern models.MesheryPattern
	if err = encoding.Unmarshal(buf.Bytes(), &pattern); err != nil {
		return nil, err
	}

	return &pattern, nil
}

func makeRequest(method, url string) (*http.Response, error) {
	req, err := utils.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return utils.MakeRequest(req)
}

func generateFilename(patternName, design, designType string) string {
	lastPartOfID := strings.Split(design, "-")[len(strings.Split(design, "-"))-1]
	filename := fmt.Sprintf("%s_%s", patternName, lastPartOfID)
	switch designType {
	case "oci":
	case "original":
		filename += ".tar.gz"
	default:
		filename += ".yaml"
	}
	return filename
}

func writeToFile(outputFilePath string, buf *bytes.Buffer) error {
	err := meshkitutils.WriteToFile(outputFilePath, buf.String())
	if err != nil {
		return err
	}
	utils.Log.Info("Design exported successfully to ", outputFilePath)
	return nil
}

func centerAlign(text string, width int) string {
	if len(text) >= width {
		return text
	}
	leftPadding := (width - len(text)) / 2
	rightPadding := width - len(text) - leftPadding
	return strings.Repeat(" ", leftPadding) + text + strings.Repeat(" ", rightPadding)
}

func getOwnerName(ownerID string, baseURL string) (string, error) {
	url := fmt.Sprintf("%s/api/user/profile/%s", baseURL, ownerID)
	resp, err := makeRequest(http.MethodGet, url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", models.ErrDoRequest(err, resp.Request.Method, url)
	}

	var userProfile models.User
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", ErrReadFromBody(err)
	}

	if err := encoding.Unmarshal([]byte(body), &userProfile); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s", userProfile.FirstName, userProfile.LastName), nil
}

func selectPatternPrompt(patterns []models.MesheryPattern, baseURL string) models.MesheryPattern {
	columns := []string{"Design Name", "Created At", "Updated At", "Type", "Owner", "Pattern ID"}
	widths := []int{20, 20, 20, 20, 20, 10}

	headingParts := make([]string, len(columns))
	for i, col := range columns {
		headingParts[i] = centerAlign(col, widths[i])
	}
	heading := strings.Join(headingParts, " ")

	patternInfos := make([]string, len(patterns)+1)
	patternInfos[0] = heading
	for i, pattern := range patterns {
		createdAt := "N/A"
		updatedAt := "N/A"
		owner := "N/A"

		if pattern.CreatedAt != nil {
			createdAt = pattern.CreatedAt.Format("2006-01-02 15:04:05")
		}
		if pattern.UpdatedAt != nil {
			updatedAt = pattern.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		if pattern.UserID != nil {
			ownerName, err := getOwnerName(*pattern.UserID, baseURL)
			if err == nil {
				owner = ownerName
			}
		}
		id := pattern.ID.String()
		lastFourID := id[len(id)-4:]

		rowParts := []string{
			centerAlign(pattern.Name, widths[0]),
			centerAlign(createdAt, widths[1]),
			centerAlign(updatedAt, widths[2]),
			centerAlign(pattern.Type.String, widths[3]),
			centerAlign(owner, widths[4]),
			centerAlign(lastFourID, widths[5]),
		}
		patternInfos[i+1] = strings.Join(rowParts, " ")
	}

	prompt := promptui.Select{
		Label: "Select a design",
		Items: patternInfos,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . | cyan }}{{ end }}",
			Inactive: "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . }}{{ end }}",
			Selected: "{{ if eq . \"" + heading + "\" }}{{ . | bold }}{{ else }}{{ . | green }}{{ end }}",
		},
		StartInSearchMode: false,
	}

	for {
		i, _, err := prompt.Run()
		if err == nil {
			if i == 0 {
				continue
			}
			return patterns[i-1]
		}
	}
}

func getUniqueFilename(filename string) string {
	base, ext := filepath.Split(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for i := 1; ; i++ {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		filename = fmt.Sprintf("%s(%d)%s", base, i, ext)
	}
	return filename
}

func init() {
	exportCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	exportCmd.Flags().StringVarP(&designType, "type", "", "", "Specify the design type to export")
	exportCmd.Flags().StringVarP(&outputDir, "output", "o", "", "Specify the output directory to save the design")
}
