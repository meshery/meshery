package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	tuitls "github.com/layer5io/meshkit/utils"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/models/meshmodel/entity"
	mutils "github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
	v1beta1Model "github.com/meshery/schemas/models/v1beta1/model"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ModelDefinition struct {
	ModelName         string `json:"model_name"`
	Registrant        string `json:"registrant"`
	PackageName       string `json:"packageName"`
	SourceURL         string `json:"source_url"`
	Category          string `json:"category"`
	SubCategory       string `json:"subCategory"`
	Shape             string `json:"shape"`
	PrimaryColor      string `json:"primaryColor"`
	SecondaryColor    string `json:"secondaryColor"`
	SVGColor          string `json:"svgColor"`
	SVGWhite          string `json:"svgWhite"`
	SVGComplete       string `json:"svgComplete"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	PublishToRegistry string `json:"publishToRegistry"`
	IsAnnotation      string `json:"isAnnotation"`
	Capabilities      string `json:"capabilities"`
}

type ConnectionDefinition struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Models      []ModelDefinition `json:"connections"`
}

var (
	logDirPath   = filepath.Join(mutils.GetHome(), ".meshery", "logs", "models")
	logFile      *os.File
	errorLogFile *os.File
	cwd          string
	outputDir    string
	registryDir  string
	defVersion   = "v1.0.0"
)

var generateModelCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a model from a Model definition",
	Long:  "Generate a model by referencing a Connection definition file",

	PreRunE: func(cmd *cobra.Command, args []string) error {
		err := os.MkdirAll(logDirPath, 0755)
		if err != nil {
			return fmt.Errorf("error creating log directory: %w", err)
		}

		utils.Log.SetLevel(logrus.DebugLevel)
		logFilePath := filepath.Join(logDirPath, "model-generation.log")
		logFile, err = os.Create(logFilePath)
		if err != nil {
			return err
		}

		utils.LogError.SetLevel(logrus.ErrorLevel)
		errorLogFilePath := filepath.Join(logDirPath, "model-generation-errors.log")
		errorLogFile, err = os.Create(errorLogFilePath)
		if err != nil {
			return err
		}

		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ = os.Getwd()
		registryDir = filepath.Join(cwd, outputDir)
		file, err := cmd.Flags().GetString("file")
		if err != nil {
			return err
		}

		// Read and parse the JSON file
		content, err := ioutil.ReadFile(file)
		if err != nil {
			utils.LogError.Error(errors.New("error reading the input file"))
			return err
		}
		var connectionDef ConnectionDefinition
		if err := json.Unmarshal(content, &connectionDef); err != nil {
			utils.LogError.Error(errors.New("error unmarshalling the input file"))
			return err
		}

		if len(connectionDef.Models) == 0 {
			utils.LogError.Error(errors.New("No model found"))
			return err
		}

		// Process models concurrently
		processModelsConcurrently(connectionDef)
		return nil
	},
}

func processModelsConcurrently(connectionDef ConnectionDefinition) {
	defer func() {
		utils.Log.UpdateLogOutput(os.Stdout)
		utils.LogError.UpdateLogOutput(os.Stdout)
	}()

	var wg sync.WaitGroup
	concurrencyLimit := make(chan struct{}, 20) // Limit to 20 concurrent goroutines

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	multiErrorWriter := io.MultiWriter(os.Stdout, errorLogFile)

	utils.Log.UpdateLogOutput(multiWriter)
	utils.LogError.UpdateLogOutput(multiErrorWriter)

	for _, modelDef := range connectionDef.Models {
		wg.Add(1)
		concurrencyLimit <- struct{}{}

		go func(model ModelDefinition) {
			defer wg.Done()
			defer func() { <-concurrencyLimit }()

			utils.Log.Info("Generating model - ", model.ModelName)
			generator, err := generators.NewGenerator(model.Registrant, model.SourceURL, model.PackageName)
			if err != nil {
				utils.LogError.Error(errors.New(fmt.Sprintf("error generating model: %s", model.ModelName)))
			}

			pkg, err := generator.GetPackage()
			if err != nil {
				utils.LogError.Error(errors.New(fmt.Sprintf("error generating model: %s", model.ModelName)))
			}

			version := pkg.GetVersion()
			modelDirPath, compDirPath, err := createVersionedDirectory(version, model.PackageName)
			if err != nil {
				utils.LogError.Error(errors.New(fmt.Sprintf("error generating model: %s", model.ModelName)))
			}

			modelDef, alreadyExsit, err := writeModelDefToFileSystem(&model, version, modelDirPath)
			if err != nil {
				utils.LogError.Error(err)
			}

			if alreadyExsit {
				utils.Log.Info("Model already exists: " + model.ModelName)
			}

			comps, err := pkg.GenerateComponents()
			if err != nil {
				utils.LogError.Error(err)
			}
			lengthOfComps := len(comps)

			for _, comp := range comps {
				comp.Version = defVersion
				// Assign the component status corresponding to modelDef status.
				// i.e. If modelDef is enabled comps are also "enabled". Ultimately all individual comps itself will have ability to control their status.
				// The status "enabled" indicates that the component will be registered inside the registry.
				if modelDef.Metadata == nil {
					modelDef.Metadata = &v1beta1Model.ModelDefinition_Metadata{}
				}
				if modelDef.Metadata.AdditionalProperties == nil {
					modelDef.Metadata.AdditionalProperties = make(map[string]interface{})
				}

				if comp.Model.Metadata.AdditionalProperties != nil {
					modelDef.Metadata.AdditionalProperties["source_uri"] = comp.Model.Metadata.AdditionalProperties["source_uri"]
				}
				comp.Model = *modelDef

				utils.AssignDefaultsForCompDefs(&comp, modelDef)
				compAlreadyExist, err := comp.WriteComponentDefinition(compDirPath, "json")
				if compAlreadyExist {
					lengthOfComps--
				}
				if err != nil {
					utils.LogError.Error(err)
				}
			}

			utils.Log.Info(" extracted ", lengthOfComps, " components for ", model.ModelName, " (", model.PackageName, ")")

		}(modelDef)
	}

	wg.Wait()
	utils.Log.Info("Model generation complete")
}

func writeModelDefToFileSystem(model *ModelDefinition, version, modelDir string) (*v1beta1Model.ModelDefinition, bool, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	filePath := filepath.Join(modelDir, "model.json")
	tmpFilePath := filepath.Join(modelDir, "tmp_model.json")

	// Ensure the temporary file is removed regardless of what happens
	defer func() {
		_ = os.Remove(tmpFilePath)
	}()

	// Check if the file exists

	if _, err := os.Stat(filePath); err == nil {
		existingData, err := os.ReadFile(filePath)
		if err != nil {
			goto NewGen
		}

		err = modelDef.WriteModelDefinition(tmpFilePath, "json")
		if err != nil {
			goto NewGen
		}

		newData, err := os.ReadFile(tmpFilePath)
		if err != nil {
			goto NewGen
		}

		// Compare the existing and new data
		if bytes.Equal(existingData, newData) {
			var oldModelDef v1beta1Model.ModelDefinition
			err = encoding.Unmarshal(existingData, &oldModelDef)
			if err != nil {
				goto NewGen
			}
			// If they are the same, return without changes
			return &oldModelDef, true, nil
		}
	}
NewGen:
	// Write the model definition to the actual file if it's new or different
	err := modelDef.WriteModelDefinition(filePath, "json")
	if err != nil {
		return nil, false, err
	}

	return &modelDef, false, nil
}

func (modelDef *ModelDefinition) CreateModelDefinition(version, defVersion string) model.ModelDefinition {

	status := entity.Ignored
	if strings.ToLower(modelDef.PublishToRegistry) == "true" {
		status = entity.Enabled
	}
	// Creating the Category and Registrant from Metadata
	var catname category.CategoryDefinition
	catname.Name = modelDef.Category
	registrant := createNewRegistrant(modelDef.Registrant)

	// Constructing the ModelDefinition with values from connectionDefinition
	models := model.ModelDefinition{
		Category:    catname,
		Description: modelDef.Description,
		DisplayName: modelDef.ModelName,

		SchemaVersion: v1beta1.ModelSchemaVersion,
		Name:          modelDef.PackageName,
		Status:        model.ModelDefinitionStatus(status),
		Registrant:    registrant,
		SubCategory:   modelDef.SubCategory,
		Model: model.Model{
			Version: version,
		},
		Version: defVersion,
	}
	err := modelDef.updateModelDefinition(&models)
	if err != nil {
	}
	return models
}

var modelMetadataValues = []string{
	"PrimaryColor", "SecondaryColor", "SVGColor", "SVGWhite", "SVGComplete", "styleOverrides", "capabilities", "isAnnotation", "Shape",
}

func (m *ModelDefinition) updateModelDefinition(modelDef *model.ModelDefinition) error {
	metadatas := modelDef.Metadata
	if metadatas == nil {
		metadatas = &model.ModelDefinition_Metadata{}
	}
	if metadatas.AdditionalProperties == nil {
		metadatas.AdditionalProperties = make(map[string]interface{})
	}
	modelMetadata, err := tuitls.MarshalAndUnmarshal[ModelDefinition, map[string]interface{}](*m)
	if err != nil {
		return err
	}
	for _, key := range modelMetadataValues {
		switch key {
		case "PrimaryColor":
			if m.PrimaryColor != "" {
				metadatas.PrimaryColor = &m.PrimaryColor
			}
		case "SecondaryColor":
			if m.SecondaryColor != "" {
				metadatas.SecondaryColor = &m.SecondaryColor
			}
		case "SVGColor", "SVGWhite", "SVGComplete":
			var svg string
			if key == "SVGColor" {
				svg = m.SVGColor
			} else if key == "SVGWhite" {
				svg = m.SVGWhite
			} else {
				svg = m.SVGComplete
			}
			// Attempt to update the SVG string
			updatedSvg, err := tuitls.UpdateSVGString(svg, 20, 20, false)
			if err == nil {
				if key == "SVGColor" {
					metadatas.SvgColor = updatedSvg
				} else if key == "SVGWhite" {
					metadatas.SvgWhite = updatedSvg
				} else {
					metadatas.SvgComplete = &updatedSvg
				}
			} else {
				// If SVG update fails, use the original SVG value
				metadatas.AdditionalProperties[key] = svg
			}
		case "capabilities":
			var capabilities []capability.Capability
			if m.Capabilities != "" {
				err := encoding.Unmarshal([]byte(m.Capabilities), &capabilities)
				if err != nil {
					return err
				}
			}
			metadatas.Capabilities = &capabilities
		case "isAnnotation":
			isAnnotation := false
			if strings.ToLower(m.IsAnnotation) == "true" {
				isAnnotation = true
			}
			metadatas.IsAnnotation = &isAnnotation
		default:
			// For keys that do not have a direct mapping, store them in AdditionalProperties
			metadatas.AdditionalProperties[key] = modelMetadata[key]
		}
	}
	modelDef.Metadata = metadatas
	return nil
}

func createNewRegistrant(registrantName string) connection.Connection {
	kind := tuitls.ReplaceSpacesAndConvertToLowercase(registrantName)
	switch kind {
	case "artifacthub":
		registrantName = "Artifact Hub"
	case "github":
		registrantName = "Github"
	case "meshery":
		registrantName = "meshery"
	case "kubernetes":
		registrantName = "Kubernetes"
	}
	newRegistrant := connection.Connection{
		Name:   registrantName,
		Status: connection.Discovered,
		Type:   "registry",
		Kind:   kind,
	}
	return newRegistrant
}

func createVersionedDirectory(version, modelName string) (string, string, error) {
	modelDirPath := filepath.Join(registryDir, modelName, version, defVersion)
	err := mutils.CreateDirectory(modelDirPath)
	if err != nil {
		return "", "", err
	}

	compDirPath := filepath.Join(modelDirPath, "components")
	err = mutils.CreateDirectory(compDirPath)
	return modelDirPath, compDirPath, err
}

func init() {
	generateModelCmd.Flags().StringP("file", "f", "", "Path to the Connectionss definition file")
	generateModelCmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "../server/meshmodel", "location to output generated models, defaults to ../server/meshmodels")
	ModelCmd.AddCommand(generateModelCmd)
}
