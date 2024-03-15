// Copyright 2023 Layer5, Inc.
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

package system

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"github.com/layer5io/meshery/mesheryctl/internal/cli/root/config"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils/walker"
	"github.com/manifoldco/promptui"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	oras "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
)

var (
	// flag used to specify the page number in list command
	pageNumberFlag int
	// flag used to specify format of output of view {model-name} command
	outFormatFlag string

	// Maximum number of rows to be displayed in a page
	maxRowsPerPage = 25

	// Color for the whiteboard printer
	whiteBoardPrinter = color.New(color.FgHiBlack, color.BgWhite, color.Bold)
)

// represents the mesheryctl exp model list subcommand.
var listModelCmd = &cobra.Command{
	Use:   "list",
	Short: "list registered models",
	Long:  "list name of all registered models",
	Example: `
// View list of models
mesheryctl exp model list

// View list of models with specified page number (25 models per page)
mesheryctl exp model list --page 2
    `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return errors.New(utils.SystemModelSubError("this command takes no arguments\n", "list"))
		}
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		var url string
		if cmd.Flags().Changed("page") {
			url = fmt.Sprintf("%s/api/meshmodels/models?page=%d", baseUrl, pageNumberFlag)
		} else {
			url = fmt.Sprintf("%s/api/meshmodels/models?pagesize=all", baseUrl)
		}
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"Model", "Category", "Version"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
			}
		}

		if len(rows) == 0 {
			// if no model is found
			// fmt.Println("No model(s) found")
			whiteBoardPrinter.Println("No model(s) found")
			return nil
		}

		if cmd.Flag("count").Value.String() == "true" {
			// fmt.Println("Total number of models: ", len(rows))
			whiteBoardPrinter.Println("Total number of models: ", len(rows))
			return nil
		}

		if cmd.Flags().Changed("page") {
			utils.PrintToTable(header, rows)
		} else {
			startIndex := 0
			endIndex := min(len(rows), startIndex+maxRowsPerPage)
			for {
				// Clear the entire terminal screen
				utils.ClearLine()

				// Print number of models and current page number
				whiteBoardPrinter.Print("Total number of models: ", len(rows))
				fmt.Println()
				whiteBoardPrinter.Print("Page: ", startIndex/maxRowsPerPage+1)
				fmt.Println()

				whiteBoardPrinter.Println("Press Enter or ↓ to continue, Esc or Ctrl+C (Ctrl+Cmd for OS user) to exit")

				utils.PrintToTable(header, rows[startIndex:endIndex])
				keysEvents, err := keyboard.GetKeys(10)
				if err != nil {
					return err
				}

				defer func() {
					_ = keyboard.Close()
				}()

				event := <-keysEvents
				if event.Err != nil {
					utils.Log.Error(fmt.Errorf("unable to capture keyboard events"))
					break
				}

				if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
					break
				}

				if event.Key == keyboard.KeyEnter || event.Key == keyboard.KeyArrowDown {
					startIndex += maxRowsPerPage
					endIndex = min(len(rows), startIndex+maxRowsPerPage)
				}

				if startIndex >= len(rows) {
					break
				}
			}
		}

		return nil
	},
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// represents the mesheryctl exp model view [model-name] subcommand.
var viewModelCmd = &cobra.Command{
	Use:   "view",
	Short: "view model",
	Long:  "view a model queried by its name",
	Example: `
// View current provider
mesheryctl exp model view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model view [model-name]\nRun 'mesheryctl exp model view --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("model name isn't specified\n\n%v", errMsg)
		} else if len(args) > 1 {
			return fmt.Errorf("too many arguments\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		model := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/models/%s?pagesize=all", baseUrl, model)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		var selectedModel v1alpha1.Model

		if modelsResponse.Count == 0 {
			fmt.Println("No model(s) found for the given name ", model)
			return nil
		} else if modelsResponse.Count == 1 {
			selectedModel = modelsResponse.Models[0]
		} else {
			selectedModel = selectModelPrompt(modelsResponse.Models)
		}

		var output []byte

		// user may pass flag in lower or upper case but we have to keep it lower
		// in order to make it consistent while checking output format
		outFormatFlag = strings.ToLower(outFormatFlag)
		if outFormatFlag == "yaml" {
			if output, err = yaml.Marshal(selectedModel); err != nil {
				return errors.Wrap(err, "failed to format output in YAML")
			}
			fmt.Print(string(output))
		} else if outFormatFlag == "json" {
			return outputJson(selectedModel)
		} else {
			return errors.New("output-format choice invalid, use [json|yaml]")
		}

		return nil
	},
}

// represents the mesheryctl exp model search [query-text] subcommand.
var searchModelCmd = &cobra.Command{
	Use:   "search",
	Short: "search models",
	Long:  "search a models by search string",
	Example: `
// View current provider
mesheryctl exp model search [query-text]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model search [query-text]\nRun 'mesheryctl exp model search --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("Search term is missing\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}

		baseUrl := mctlCfg.GetBaseMesheryURL()
		queryText := args[0]

		url := fmt.Sprintf("%s/api/meshmodels/models?search=%s&pagesize=all", baseUrl, queryText)
		req, err := utils.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		resp, err := utils.MakeRequest(req)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		// defers the closing of the response body after its use, ensuring that the resources are properly released.
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		modelsResponse := &models.MeshmodelsAPIResponse{}
		err = json.Unmarshal(data, modelsResponse)
		if err != nil {
			utils.Log.Error(err)
			return err
		}

		header := []string{"Model", "Category", "Version"}
		rows := [][]string{}

		for _, model := range modelsResponse.Models {
			if len(model.DisplayName) > 0 {
				rows = append(rows, []string{model.Name, model.Category.Name, model.Version})
			}
		}

		if len(rows) == 0 {
			// if no model is found
			// fmt.Println("No model(s) found")
			whiteBoardPrinter.Println("No model(s) found")
			return nil
		} else {
			utils.PrintToTable(header, rows)
		}

		return nil
	},
}

var importModelCmd = &cobra.Command{
	Use:   "import",
	Short: "import model",
	Long:  "import a model from the registry",
	Example: `

// Import a model
mesheryctl exp model import --file <URI>
	`,
	Args: func(_ *cobra.Command, args []string) error {
		const errMsg = "Usage: mesheryctl exp model import <URI/path to Model>\nRun 'mesheryctl exp model import --help' to see detailed help message"
		if len(args) == 0 {
			return fmt.Errorf("Path/URI term is missing\n\n%v", errMsg)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Here we have to call the model import API /api/meshmodel/components/register with POST
		// The body of the request will be the json of type below
		// type MeshModelRegistrantData struct {
		// 	Host       Host                 `json:"host"`
		// 	EntityType types.CapabilityType `json:"entityType"`
		// 	Entity     []byte               `json:"entity"` //This will be type converted to appropriate entity on server based on passed entity type
		// }

		// Check if file path is provided
		if validURL := govalidator.IsURL(args[0]); !validURL {
			if _, err := os.Stat(args[0]); os.IsNotExist(err) {
				return errors.New("file path does not exist")
			}
			// get the directory path
			dirPath := args[0]
			// read the files in the directory
			err := processDirectory(dirPath)
			if err != nil {
				fmt.Println("Error processing directory:", err)
				return err
			}
		} else if isGithubURL := strings.Contains(args[0], "github.com"); isGithubURL {
			// Here we will be reading the content of the file from the URL

			// Parse the URL
			parsedURL, err := url.Parse(args[0])
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return err
			}

			// Extract parts
			pathParts := strings.Split(parsedURL.Path, "/")

			githubInfo := walker.NewGithub()
			githubInfo.Owner(pathParts[1])
			githubInfo.Repo(pathParts[2])
			githubInfo.Branch(pathParts[4])
			githubInfo.Root(strings.Join(pathParts[5:], "/"))
			FileInterceptor := func(file walker.GithubContentAPI) error {

				// Download the file content
				contentResponse, err := http.Get(file.DownloadURL)
				if err != nil {
					return err
				}
				defer contentResponse.Body.Close()

				// Read the file content directly into a []byte
				fileContent, err := ioutil.ReadAll(contentResponse.Body)
				if err != nil {
					return err
				}

				// Register the file content as a component
				if err := registerComponent(file.Name, fileContent); err != nil {
					return err
				}

				return nil
			}

			githubInfo.RegisterFileInterceptor(FileInterceptor)

			if err := githubInfo.Walk(); err != nil {
				fmt.Println("Error occurred during traversal:", err)
				return err
			}
		} else {
			// If the user wants to import the model from the oci registry

			// Here extract the registry address and tag
			// docker.io/repo/container:tag
			reg := strings.Split(args[0], ":")
			regAddress := reg[0]
			tag := reg[1]
			// 0. Create an OCI layout store
			fs, err := file.New("./tmp/")
			if err != nil {
				panic(err)
			}
			defer fs.Close()
			// 1. Connect to a remote repository
			ctx := context.Background()
			repo, err := remote.NewRepository(regAddress)
			if err != nil {
				panic(err)
			}

			// 2. Copy from the remote repository to the OCI layout store
			_, err = oras.Copy(ctx, repo, tag, fs, tag, oras.DefaultCopyOptions)
			if err != nil {
				panic(err)
			}

			// fmt.Println("manifest pulled:", manifestDescriptor.Digest, manifestDescriptor.MediaType)

			// read file content and send it to the server
			dirPath := "./tmp"

			err = processDirectory(dirPath)
			if err != nil {
				fmt.Println("Error processing directory:", err)
				return err
			}
		}
		return nil
	},
}

// ModelCmd represents the mesheryctl exp model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "View list of models and detail of models",
	Long:  "View list of models and detailed information of a specific model",
	Example: `
// To view list of components
mesheryctl exp model list

// To view a specific model
mesheryctl exp model view [model-name]
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		//Check prerequisite

		mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			return err
		}
		err = utils.IsServerRunning(mctlCfg.GetBaseMesheryURL())
		if err != nil {
			return err
		}
		ctx, err := mctlCfg.GetCurrentContext()
		if err != nil {
			return err
		}
		err = ctx.ValidateVersion()
		if err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		if ok := utils.IsValidSubcommand(availableSubcommands, args[0]); !ok {
			return errors.New(utils.SystemModelSubError(fmt.Sprintf("'%s' is an invalid subcommand. Please provide required options from [view]. Use 'mesheryctl exp model --help' to display usage guide.\n", args[0]), "model"))
		}
		_, err := config.GetMesheryCtl(viper.GetViper())
		if err != nil {
			log.Fatalln(err, "error processing config")
		}
		err = viewProviderCmd.RunE(cmd, args)
		if err != nil {
			return err
		}
		err = cmd.Usage()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	listModelCmd.Flags().IntVarP(&pageNumberFlag, "page", "p", 1, "(optional) List next set of models with --page (default = 1)")
	listModelCmd.Flags().BoolP("count", "c", false, "(optional) Get the number of models in total")
	viewModelCmd.Flags().StringVarP(&outFormatFlag, "output-format", "o", "yaml", "(optional) format to display in [json|yaml]")
	availableSubcommands = []*cobra.Command{listModelCmd, viewModelCmd, searchModelCmd, importModelCmd}
	ModelCmd.AddCommand(availableSubcommands...)
}

// selectModelPrompt lets user to select a model if models are more than one
func selectModelPrompt(models []v1alpha1.Model) v1alpha1.Model {
	modelArray := []v1alpha1.Model{}
	modelNames := []string{}

	modelArray = append(modelArray, models...)

	for _, model := range modelArray {
		modelName := fmt.Sprintf("%s, version: %s", model.DisplayName, model.Version)
		modelNames = append(modelNames, modelName)
	}

	prompt := promptui.Select{
		Label: "Select a model",
		Items: modelNames,
	}

	for {
		i, _, err := prompt.Run()
		if err != nil {
			continue
		}

		return modelArray[i]
	}
}

func outputJson(model v1alpha1.Model) error {
	if err = prettifyJson(model); err != nil {
		// if prettifyJson return error, marshal output in conventional way using json.MarshalIndent
		// but it doesn't convert unicode to its corresponding HTML string (it is default behavior)
		// e.g unicode representation of '&' will be printed as '\u0026'
		if output, err := json.MarshalIndent(model, "", "  "); err != nil {
			return errors.Wrap(err, "failed to format output in JSON")
		} else {
			fmt.Print(string(output))
		}
	}
	return nil
}

// prettifyJson takes a v1alpha1.Model struct as input, marshals it into a nicely formatted JSON representation,
// and prints it to standard output with proper indentation and without escaping HTML entities.
func prettifyJson(model v1alpha1.Model) error {
	// Create a new JSON encoder that writes to the standard output (os.Stdout).
	enc := json.NewEncoder(os.Stdout)
	// Configure the JSON encoder settings.
	// SetEscapeHTML(false) prevents special characters like '<', '>', and '&' from being escaped to their HTML entities.
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	// Any errors during the encoding process will be returned as an error.
	return enc.Encode(model)
}


func processDirectory(dirPath string) error {
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %s", dirPath, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			subdirPath := filepath.Join(dirPath, entry.Name())
			err := processDirectory(subdirPath)
			if err != nil {
				fmt.Printf("Error processing directory %s: %s\n", subdirPath, err)
				continue
			}
		} else {
			fileContent, err := ioutil.ReadFile(filepath.Join(dirPath, entry.Name()))
			if err != nil {
				fmt.Printf("Error reading file %s: %s\n", entry.Name(), err)
				continue
			}

			err = registerComponent(entry.Name(), fileContent)
			if err != nil {
				fmt.Println("Error registering component:", entry.Name(), err)
			}
		}
	}

	return nil
}

func registerComponent(fileName string, content []byte) error {

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		return err
	}
	baseURL := mctlCfg.GetBaseMesheryURL()

	registrantData := &registry.MeshModelRegistrantData{
		Host: registry.Host{
			Hostname: "localhost",
		},
		EntityType: "component",
		Entity:     content,
	}

	url := fmt.Sprintf("%s/api/meshmodel/components/register", baseURL)

	// Marshal the registrant data into JSON
	requestBody, err := json.Marshal(registrantData)
	if err != nil {
		return err
	}

	// Create a new HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Model imported successfully:", fileName)
	} else {
		fmt.Println("Error registering component:", fileName, resp.Status)
	}

	return nil
}