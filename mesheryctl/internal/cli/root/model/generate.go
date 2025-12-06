package model

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"
	"time"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	meshkitutils "github.com/meshery/meshkit/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type ModelGenerator interface {
	Generate() error
}

type UrlModelGenerator struct {
	TemplateFile string
	Url          string
	SkipRegister bool
}

type CsvModelGenerator struct {
	ModelFile        string
	ComponentFile    string
	RelationshipFile string
	SkipRegister     bool
}

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate models from a file",
	Long: `Generate models by specifying the directory, file, or URL. You can also provide a template JSON file and registrant name
Documentation for models generate can be found at https://docs.meshery.io/reference/mesheryctl/model/generate`,
	Example: ` 
// Generate a model from a CSV file(s)
mesheryctl model generate --f [path-to-csv-drectory]

// Generate a model from a Uri baesd on a JSON template
mesheryctl model generate --f [URL] -t [path-to-template.json]

// Generate a model from a Uri baesd on a JSON template skipping registration
mesheryctl model generate --f [URL] -t [path-to-template.json] -r
	`,
	Args: func(cmd *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl model generate [ file | filePath | URL ]\nRun 'mesheryctl model generate --help' to see detailed help message"
		file, _ := cmd.Flags().GetString("file")
		if file == "" && len(args) == 0 {
			return fmt.Errorf("[ file | filepath | URL ] isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var path string
		file, _ := cmd.Flags().GetString("file")
		if file != "" {
			path = file
		} else {
			path = args[0]
		}
		isUrl := utils.IsValidUrl(path)

		register, _ := cmd.Flags().GetBool("register")
		template, _ := cmd.Flags().GetString("template")

		// Path is a url
		if isUrl {
			// If template is provided, use URL model generator (existing behavior)
			if template != "" {
				urlModelGenerator := &UrlModelGenerator{
					TemplateFile: template,
					Url:          path,
					SkipRegister: register,
				}
				return urlModelGenerator.Generate()
			}

			// Otherwise, download URL content and treat as CSV
			utils.Log.Info("Downloading file from URL: ", path)

			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "mesheryctl-model-")
			if err != nil {
				return fmt.Errorf("failed to create temporary directory: %w", err)
			}
			defer os.RemoveAll(tempDir) // Clean up temp directory

			// Determine filename from URL
			u, err := url.Parse(path)
			if err != nil {
				return fmt.Errorf("failed to parse URL %q: %w", path, err)
			}

			filename := pathpkg.Base(u.Path)
			if filename == "." || filename == "/" || filename == "" {
				filename = "model.csv"
			}
			tempFilePath := filepath.Join(tempDir, filename)

			// Download the file using simple HTTP GET
			err = downloadFileSimple(tempFilePath, path)
			if err != nil {
				// Try meshkit utils as fallback
				err = meshkitutils.DownloadFile(tempFilePath, path)
				if err != nil {
					return fmt.Errorf("failed to download file from URL: %w", err)
				}
			}

			utils.Log.Info("Downloaded file to: ", tempFilePath)

			// Use the downloaded file path for CSV processing
			path = tempFilePath
		}

		// Path is a file or directory (including downloaded URL content)
		err := meshkitRegistryUtils.SetLogger(true)
		if err != nil {
			utils.Log.Info("Error setting logger: ", err)
		}

		modelcsvpath, componentcsvpath, relationshipcsvpath, err := meshkitRegistryUtils.GetCsv(path)

		if err != nil {
			return err
		}

		csvModelGenerator := &CsvModelGenerator{
			ModelFile:        modelcsvpath,
			ComponentFile:    componentcsvpath,
			RelationshipFile: relationshipcsvpath,
			SkipRegister:     register,
		}

		return csvModelGenerator.Generate()
	},
}

func init() {
	generateModelCmd.Flags().SetNormalizeFunc(func(f *pflag.FlagSet, name string) pflag.NormalizedName {
		return pflag.NormalizedName(strings.ToLower(name))
	})

	generateModelCmd.Flags().StringP("file", "f", "", "Specify path to the file or directory")
	generateModelCmd.Flags().StringP("template", "t", "", "Specify path to the template JSON file")
	generateModelCmd.Flags().BoolP("register", "r", false, "Skip registration of the model")

}

func (u *UrlModelGenerator) Generate() error {
	utils.Log.Info("Generating model from URL: ", u.Url)

	fileData, err := os.ReadFile(u.TemplateFile)
	if err != nil {
		return utils.ErrFileRead(err)
	}

	err = registerModel(fileData, nil, nil, "", "url", u.Url, !u.SkipRegister)
	if err != nil {
		return err
	}

	locationForModel := utils.MesheryFolder + "/models"
	utils.Log.Info("Model can be accessed from ", locationForModel)

	return nil
}

func (c *CsvModelGenerator) Generate() error {
	utils.Log.Info("Generating model from CSV files")

	var modelData, componentData, relationshipData []byte
	var err error

	filePaths := []struct {
		path string
		data *[]byte
	}{
		{c.ModelFile, &modelData},
		{c.ComponentFile, &componentData},
		{c.RelationshipFile, &relationshipData},
	}

	for _, f := range filePaths {
		*f.data, err = os.ReadFile(f.path)
		if err != nil {
			return utils.ErrFileRead(err)
		}
	}

	err = registerModel(modelData, componentData, relationshipData, "model.csv", "csv", "", !c.SkipRegister)
	if err != nil {
		return err
	}

	locationForModel := utils.MesheryFolder + "/models"
	utils.Log.Info("Model can be accessed from ", locationForModel)

	locationForLogs := utils.MesheryFolder + "/logs/registry"
	utils.Log.Info("Logs for the csv generation can be accessed ", locationForLogs)

	return nil
}

// downloadFileSimple downloads a file from a URL using standard HTTP GET
func downloadFileSimple(filepath string, urlStr string) error {
	client := http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		return fmt.Errorf("http GET request to %s failed: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad http status for %s: %s", urlStr, resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filepath, err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filepath, err)
	}

	return nil
}
