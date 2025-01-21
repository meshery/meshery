package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/generators/models"
	meshkitutils "github.com/layer5io/meshkit/utils"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
	"google.golang.org/api/sheets/v4"

	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/csv"
	"github.com/layer5io/meshkit/utils/store"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	_model "github.com/meshery/schemas/models/v1beta1/model"
)

var modelToCompGenerateTracker = store.NewGenericThreadSafeStore[compGenerateTracker]()

type compGenerateTracker struct {
	totalComps int
	version    string
}

var (
	GoogleSpreadSheetURL     = "https://docs.google.com/spreadsheets/d/"
	totalAggregateComponents int
	totalAggregateModel      int
	logDirPath               = filepath.Join(utils.GetHome(), ".meshery", "logs", "registry")
)
var (
	shouldRegisterMod = "publishToSites"
)
var (
	artifactHubCount        = 0
	artifactHubRateLimit    = 100
	artifactHubRateLimitDur = 5 * time.Minute
	artifactHubMutex        sync.Mutex
)
var defVersion = "v1.0.0"

type ModelCSV struct {
	Registrant         string `json:"registrant" csv:"registrant"`
	ModelDisplayName   string `json:"modelDisplayName" csv:"modelDisplayName"`
	Model              string `json:"model" csv:"model"`
	Category           string `json:"category" csv:"category"`
	SubCategory        string `json:"subCategory" csv:"subCategory"`
	Description        string `json:"description" csv:"description"`
	SourceURL          string `json:"sourceURL" csv:"sourceURL"`
	Website            string `json:"website" csv:"website"`
	Docs               string `json:"docs" csv:"docs"`
	Shape              string `json:"shape" csv:"shape"`
	PrimaryColor       string `json:"primaryColor" csv:"primaryColor"`
	SecondaryColor     string `json:"secondaryColor" csv:"secondaryColor"`
	StyleOverrides     string `json:"styleOverrides" csv:"styleOverrides"`
	Styles             string `json:"styles" csv:"styles"`
	ShapePolygonPoints string `json:"shapePolygonPoints" csv:"shapePolygonPoints"`
	DefaultData        string `json:"defaultData" csv:"defaultData"`
	Capabilities       string `json:"capabilities" csv:"capabilities"`
	LogoURL            string `json:"logoURL" csv:"logoURL"`
	SVGColor           string `json:"svgColor" csv:"svgColor"`
	SVGWhite           string `json:"svgWhite" csv:"svgWhite"`
	SVGComplete        string `json:"svgComplete" csv:"svgComplete"`
	IsAnnotation       string `json:"isAnnotation" csv:"isAnnotation"`
	PublishToRegistry  string `json:"publishToRegistry" csv:"publishToRegistry"`
	AboutProject       string `json:"aboutProject" csv:"-"`
	PageSubtTitle      string `json:"pageSubtitle" csv:"-"`
	DocsURL            string `json:"docsURL" csv:"-"`
	StandardBlurb      string `json:"standardBlurb" csv:"-"`
	Feature1           string `json:"feature1" csv:"-"`
	Feature2           string `json:"feature2" csv:"-"`
	Feature3           string `json:"feature3" csv:"-"`
	HowItWorks         string `json:"howItWorks" csv:"-"`
	HowItWorksDetails  string `json:"howItWorksDetails" csv:"-"`
	Screenshots        string `json:"screenshots" csv:"-"`
	FullPage           string `json:"fullPage" csv:"-"`
	PublishToSites     string `json:"publishToSites" csv:"-"`
}

var modelMetadataValues = []string{
	"primaryColor", "secondaryColor", "svgColor", "svgWhite", "svgComplete", "styleOverrides", "capabilities", "isAnnotation", "shape",
}

func (m *ModelCSV) UpdateModelDefinition(modelDef *_model.ModelDefinition) error {
	metadata := modelDef.Metadata
	if metadata == nil {
		metadata = &_model.ModelDefinition_Metadata{}
	}
	if metadata.AdditionalProperties == nil {
		metadata.AdditionalProperties = make(map[string]interface{})
	}
	modelMetadata, err := utils.MarshalAndUnmarshal[ModelCSV, map[string]interface{}](*m)
	if err != nil {
		return err
	}
	for _, key := range modelMetadataValues {
		switch key {
		case "primaryColor":
			if m.PrimaryColor != "" {
				metadata.PrimaryColor = &m.PrimaryColor
			}
		case "secondaryColor":
			if m.SecondaryColor != "" {
				metadata.SecondaryColor = &m.SecondaryColor
			}
		case "svgColor", "svgWhite", "svgComplete":
			var svg string
			if key == "svgColor" {
				svg = m.SVGColor
			} else if key == "svgWhite" {
				svg = m.SVGWhite
			} else {
				svg = m.SVGComplete
			}
			// Attempt to update the SVG string
			updatedSvg, err := utils.UpdateSVGString(svg, SVG_WIDTH, SVG_HEIGHT, false)
			if err == nil {
				if key == "svgColor" {
					metadata.SvgColor = updatedSvg
				} else if key == "svgWhite" {
					metadata.SvgWhite = updatedSvg
				} else {
					metadata.SvgComplete = &updatedSvg
				}
			} else {
				// If SVG update fails, use the original SVG value
				metadata.AdditionalProperties[key] = svg
			}
		case "capabilities":
			var capabilities []capability.Capability
			if m.Capabilities != "" {
				err := encoding.Unmarshal([]byte(m.Capabilities), &capabilities)
				if err != nil {
					return err
				}
			}
			metadata.Capabilities = &capabilities
		case "isAnnotation":
			isAnnotation := false
			if strings.ToLower(m.IsAnnotation) == "true" {
				isAnnotation = true
			}
			metadata.IsAnnotation = &isAnnotation
		default:
			// For keys that do not have a direct mapping, store them in AdditionalProperties
			metadata.AdditionalProperties[key] = modelMetadata[key]
		}
	}
	modelDef.Metadata = metadata
	return nil
}
func (mcv *ModelCSV) CreateModelDefinition(version, defVersion string) _model.ModelDefinition {
	status := entity.Ignored
	if strings.ToLower(mcv.PublishToRegistry) == "true" {
		status = entity.Enabled
	}
	var catname category.CategoryDefinition
	catname.Name = mcv.Category
	registrant := createNewRegistrant(mcv.Registrant)

	model := _model.ModelDefinition{
		Category:    catname,
		Description: mcv.Description,
		DisplayName: mcv.ModelDisplayName,

		SchemaVersion: v1beta1.ModelSchemaVersion,
		Name:          mcv.Model,
		Status:        _model.ModelDefinitionStatus(status),
		Registrant:    registrant,
		SubCategory:   mcv.SubCategory,
		Model: _model.Model{
			Version: version,
		},
		Version: defVersion,
	}
	err := mcv.UpdateModelDefinition(&model)
	if err != nil {
		Log.Error(err)
	}
	return model
}

type ModelCSVHelper struct {
	SpreadsheetID  int64
	SpreadsheetURL string
	Title          string
	CSVPath        string
	Models         []ModelCSV
}

func NewModelCSVHelper(sheetURL, spreadsheetName string, spreadsheetID int64, localCsvPath string) (*ModelCSVHelper, error) {
	var csvPath string
	// Download the CSV file from the spreadsheet URL
	if localCsvPath == "" {
		// Set the directory path for storing the downloaded CSV
		dirPath := filepath.Join(utils.GetHome(), ".meshery", "content")
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			return nil, utils.ErrCreateDir(err, dirPath)
		}

		// Set the CSV file path
		csvPath = filepath.Join(dirPath, "models.csv")
		sheetURL = sheetURL + "/pub?output=csv" + "&gid=" + strconv.FormatInt(spreadsheetID, 10)
		Log.Info("Downloading CSV from: ", sheetURL)
		err = utils.DownloadFile(csvPath, sheetURL)
		if err != nil {
			return nil, utils.ErrReadingRemoteFile(err)
		}

	} else {
		csvPath = localCsvPath
	}

	// Return the initialized ModelCSVHelper with the appropriate CSV path
	return &ModelCSVHelper{
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: sheetURL,
		Models:         []ModelCSV{},
		CSVPath:        csvPath,
		Title:          spreadsheetName,
	}, nil
}

func (mch *ModelCSVHelper) ParseModelsSheet(parseForDocs bool, modelName string) error {
	ch := make(chan ModelCSV, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[ModelCSV](mch.CSVPath, rowIndex, nil, func(columns []string, currentRow []string) bool {
		index := 0
		if parseForDocs {
			index = GetIndexForRegisterCol(columns, shouldRegisterMod)
		} else {
			// Generation of models should not consider publishedToRegistry column value.
			// Generation should happen for all models, while during registration "published" attribute should be respected.
			return true
		}
		if index != -1 && index < len(currentRow) {
			shouldRegister := currentRow[index]
			return strings.ToLower(shouldRegister) == "true"
		}
		return false
	})
	if err != nil {
		return ErrFileRead(err)
	}

	go func() {
		Log.Info("Parsing Models...")
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			errorChan <- err
		}
	}()
	for {
		select {

		case data := <-ch:
			if modelName != "" && data.Model != modelName {
				continue
			}
			mch.Models = append(mch.Models, data)
		case err := <-errorChan:
			return ErrFileRead(err)

		case <-csvReader.Context.Done():
			return nil
		}
	}
}

// template := `---
// title: <model-display-name>
// subtitle: <Page Subtitle>
// integrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh.svg
// darkModeIntegrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh_white.svg
// docURL: <Docs URL>
// category: <Category>
// subcategory: <Sub-Category>
// featureList: [<Feature 1>,<Feature 2>,<Feature 3>]
// workingSlides: [
//
//	../_images/meshmap-visualizer.png,
//	../_images/meshmap-designer.png]
//
// howItWorks: <howItWorks>
// howItWorksDetails: howItWorksDetails
// published: <Publish>
// ---
// <p>
//
//	<About Project>
//
// </p>
// <p>
//
//	<Standard Blurb>
//
// </p>`
func (m ModelCSV) CreateMarkDownForMDXStyle(componentsMetadata string) string {
	formattedName := utils.FormatName(m.Model)
	var template string = `---
title: %s
subtitle: %s
integrationIcon: icons/color/%s-color.svg
darkModeIntegrationIcon: icons/white/%s-white.svg
docURL: %s
description: %s
category: %s
subcategory: %s
registrant: %s
components: %v
featureList: [
  "%s",
  "%s",
  "%s"
]
workingSlides: [
  %s,
  %s
]
howItWorks: "%s"
howItWorksDetails: "%s"
published: %s
---
<p>
%s
</p>
%s
`
	markdown := fmt.Sprintf(template,
		m.ModelDisplayName,
		m.PageSubtTitle,
		formattedName,
		formattedName,
		m.DocsURL,
		m.Description,
		m.Category,
		m.SubCategory,
		m.Registrant,
		componentsMetadata,
		m.Feature1,
		m.Feature2,
		m.Feature3,
		`../_images/kanvas-visualizer.png`,
		`../_images/kanvas-designer.png`,
		m.HowItWorks,
		m.HowItWorksDetails,
		m.PublishToSites,
		m.AboutProject,
		m.StandardBlurb,
	)
	markdown = strings.ReplaceAll(markdown, "\r", "\n")
	return markdown
}
func createNewRegistrant(registrantName string) connection.Connection {
	kind := utils.ReplaceSpacesAndConvertToLowercase(registrantName)
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

// Creates JSON formatted meshmodel attribute item for JSON Style docs
func (m ModelCSV) CreateJSONItem(iconDir string) string {
	formattedModelName := utils.FormatName(m.Model)
	json := "{"
	json += fmt.Sprintf("\"name\":\"%s\"", m.Model)
	// If SVGs exist, then add the paths to json
	if m.SVGColor != "" {
		json += fmt.Sprintf(",\"color\":\"%s/icons/color/%s-color.svg\"", iconDir, formattedModelName)
	}

	if m.SVGWhite != "" {
		json += fmt.Sprintf(",\"white\":\"%s/icons/white/%s-white.svg\"", iconDir, formattedModelName)
	}

	json += fmt.Sprintf(",\"permalink\":\"%s\"", m.DocsURL)

	json += "}"
	return json
}

func (m ModelCSV) CreateMarkDownForMDStyle(componentsMetadata, relationshipMetadata string, componentsCount, relationshipsCount int, outputFor string) string {
	formattedName := utils.FormatName(m.Model)
	var template, markdown string

	if outputFor == "mesherydocs" {
		template = `---
layout: integration
title: %s
subtitle: %s
image: /assets/img/integrations/%s/icons/color/%s-color.svg
permalink: extensibility/integrations/%s
docURL: %s
description: %s
integrations-category: %s
integrations-subcategory: %s
registrant: %s
components: %v
components-count: %d
relationships: %v
relationship-count: %d
featureList: [
  "%s",
  "%s",
  "%s"
]
howItWorks: "%s"
howItWorksDetails: "%s"
language: en
list: include
type: extensibility
category: integrations
---
`
		markdown = fmt.Sprintf(template,
			m.ModelDisplayName,
			m.PageSubtTitle,
			formattedName,
			formattedName,
			formattedName,
			m.DocsURL,
			m.Description,
			m.Category,
			m.SubCategory,
			m.Registrant,
			componentsMetadata,
			componentsCount,
			relationshipMetadata,
			relationshipsCount,
			m.Feature1,
			m.Feature2,
			m.Feature3,
			m.HowItWorks,
			m.HowItWorksDetails,
		)
	} else if outputFor == "mesheryio" {
		template = `---
layout: single-page-model
item-type: model
name: %s
subtitle: %s
colorIcon: /assets/images/integration/%s/icons/color/%s-color.svg
whiteIcon: /assets/images/integration/%s/icons/white/%s-white.svg
docURL: %s
description: %s
category: %s
subcategory: %s
registrant: %s
components: %v
componentsCount: %d
relationships: %v
relationshipsCount: %d
featureList: [
  "%s",
  "%s",
  "%s"
]
howItWorks: "%s"
howItWorksDetails: "%s"
---
`
		markdown = fmt.Sprintf(template,
			m.ModelDisplayName,
			m.PageSubtTitle,
			formattedName,
			formattedName,
			formattedName,
			formattedName,
			m.DocsURL,
			m.Description,
			m.Category,
			m.SubCategory,
			m.Registrant,
			componentsMetadata,
			componentsCount,
			relationshipMetadata,
			relationshipsCount,
			m.Feature1,
			m.Feature2,
			m.Feature3,
			m.HowItWorks,
			m.HowItWorksDetails,
		)
	}

	markdown = strings.ReplaceAll(markdown, "\r", "\n")

	return markdown
}

func (m ModelCSVHelper) Cleanup() error {
	// remove csv file
	Log.Info("Removing CSV file: ", m.CSVPath)
	err := os.Remove(m.CSVPath)
	if err != nil {
		return err
	}

	return nil
}
func AssignDefaultsForCompDefs(componentDef *component.ComponentDefinition, modelDef *_model.ModelDefinition) {
	// Assign the status from the model to the component
	compStatus := component.ComponentDefinitionStatus(modelDef.Status)
	componentDef.Status = &compStatus

	// Initialize AdditionalProperties and Styles if nil
	if componentDef.Metadata.AdditionalProperties == nil {
		componentDef.Metadata.AdditionalProperties = make(map[string]interface{})
	}
	if componentDef.Styles == nil {
		componentDef.Styles = &component.Styles{}
	}

	// Use reflection to map model metadata to component styles
	stylesValue := reflect.ValueOf(componentDef.Styles).Elem()

	// Iterate through modelDef.Metadata
	if modelDef.Metadata != nil {
		if modelDef.Metadata.AdditionalProperties["styleOverrides"] != nil {
			styleOverrides, ok := modelDef.Metadata.AdditionalProperties["styleOverrides"].(string)
			if ok {
				err := encoding.Unmarshal([]byte(styleOverrides), &componentDef.Styles)
				if err != nil {
					LogError.Error(err)
				}
			}
		}
		if (modelDef.Metadata.Capabilities) != nil {
			componentDef.Capabilities = modelDef.Metadata.Capabilities
		}
		if modelDef.Metadata.PrimaryColor != nil {
			componentDef.Styles.PrimaryColor = *modelDef.Metadata.PrimaryColor
		}
		if modelDef.Metadata.SecondaryColor != nil {
			componentDef.Styles.SecondaryColor = modelDef.Metadata.SecondaryColor
		}
		if modelDef.Metadata.SvgColor != "" {
			componentDef.Styles.SvgColor = modelDef.Metadata.SvgColor
		}
		if modelDef.Metadata.SvgComplete != nil {
			componentDef.Styles.SvgComplete = *modelDef.Metadata.SvgComplete
		}
		if modelDef.Metadata.SvgWhite != "" {
			componentDef.Styles.SvgWhite = modelDef.Metadata.SvgWhite
		}

		// Iterate through AdditionalProperties and assign appropriately
		for k, v := range modelDef.Metadata.AdditionalProperties {
			if k == "styleOverrides" {
				continue
			}
			// Check if the field exists in Styles
			if field := stylesValue.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(k, name)
			}); field.IsValid() && field.CanSet() {
				switch field.Kind() {
				case reflect.Ptr:
					ptrType := field.Type().Elem()
					val := reflect.New(ptrType).Elem()

					if val.Kind() == reflect.String {
						val.SetString(v.(string))
					} else if val.Kind() == reflect.Float32 {
						val.SetFloat(v.(float64))
					} else if val.Kind() == reflect.Int {
						val.SetInt(int64(v.(int)))
					} else {
						val.Set(reflect.ValueOf(v))
					}

					field.Set(val.Addr())
				case reflect.String:
					field.SetString(v.(string))
				case reflect.Float32:
					field.SetFloat(v.(float64))
				case reflect.Int:
					field.SetInt(int64(v.(int)))
				default:
					field.Set(reflect.ValueOf(v))
				}
			} else {
				componentDef.Metadata.AdditionalProperties[k] = v
			}
		}
	}
}
func GenerateComponentsFromPkg(pkg models.Package, compDirPath string, defVersion string, modelDef _model.ModelDefinition) (int, int, error) {
	comps, err := pkg.GenerateComponents()
	if err != nil {
		return 0, 0, err
	}
	lengthOfComps := len(comps)
	for _, comp := range comps {
		comp.Version = defVersion
		if modelDef.Metadata == nil {
			modelDef.Metadata = &_model.ModelDefinition_Metadata{}
		}
		if modelDef.Metadata.AdditionalProperties == nil {
			modelDef.Metadata.AdditionalProperties = make(map[string]interface{})
		}
		if comp.Model.Metadata.AdditionalProperties != nil {
			modelDef.Metadata.AdditionalProperties["source_uri"] = comp.Model.Metadata.AdditionalProperties["source_uri"]
		}
		comp.Model = modelDef

		AssignDefaultsForCompDefs(&comp, &modelDef)
		alreadyExists, err := comp.WriteComponentDefinition(compDirPath, "json")
		if err != nil {
			return 0, 0, err
		}
		if alreadyExists {
			lengthOfComps--
		}
	}
	return len(comps), lengthOfComps, nil
}
func GenerateModels(registrant string, sourceURl string, modelName string) (models.Package, string, error) {
	generator, err := generators.NewGenerator(registrant, sourceURl, modelName)
	if err != nil {
		return nil, "", err
	}
	pkg, err := generator.GetPackage()
	if err != nil {
		return nil, "", err
	}
	version := pkg.GetVersion()
	return pkg, version, nil
}

func RateLimitArtifactHub() {
	artifactHubMutex.Lock()
	defer artifactHubMutex.Unlock()

	if artifactHubCount > 0 && artifactHubCount%artifactHubRateLimit == 0 {
		Log.Info("Rate limit reached for Artifact Hub. Waiting for ", artifactHubRateLimitDur)
		time.Sleep(artifactHubRateLimitDur)
	}
	artifactHubCount++
}
func createVersionedDirectoryForModelAndComp(version, modelName, storeLocation string) (string, string, error) {
	modelDirPath := filepath.Join(storeLocation, modelName, version, defVersion)
	err := utils.CreateDirectory(modelDirPath)
	if err != nil {
		return "", "", err
	}

	compDirPath := filepath.Join(modelDirPath, "components")
	err = utils.CreateDirectory(compDirPath)
	return modelDirPath, compDirPath, err
}

func writeModelDefToFileSystem(model *ModelCSV, version, modelDefPath string) (*_model.ModelDefinition, bool, error) {
	modelDef := model.CreateModelDefinition(version, defVersion)
	filePath := filepath.Join(modelDefPath, "model.json")
	tmpFilePath := filepath.Join(modelDefPath, "tmp_model.json")

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
			var oldModelDef _model.ModelDefinition
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
func parseModelSheet(url string, modelName string, sheetGID int64, modelCSVFilePath string) (*ModelCSVHelper, error) {
	modelCSVHelper, err := NewModelCSVHelper(url, "Models", sheetGID, modelCSVFilePath)
	if err != nil {
		return nil, err
	}

	err = modelCSVHelper.ParseModelsSheet(false, modelName)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return modelCSVHelper, nil
}
func parseComponentSheet(url string, modelName string, componentSpredsheetGID int64, componentCSVFilePath string) (*ComponentCSVHelper, error) {
	compCSVHelper, err := NewComponentCSVHelper(url, "Components", componentSpredsheetGID, componentCSVFilePath)
	if err != nil {
		return nil, err
	}

	err = compCSVHelper.ParseComponentsSheet(modelName)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return compCSVHelper, nil
}
func parseRelationshipSheet(url string, modelName string, relationshipSpredsheetGID int64, relationshipCSVFilePath string) (*RelationshipCSVHelper, error) {
	relationshipCSVHelper, err := NewRelationshipCSVHelper(url, "Relationships", relationshipSpredsheetGID, relationshipCSVFilePath)
	if err != nil {
		return nil, err
	}
	err = relationshipCSVHelper.ParseRelationshipsSheet(modelName)
	if err != nil {
		return nil, ErrGenerateModel(err, "unable to start model generation")
	}
	return relationshipCSVHelper, nil
}
func logModelGenerationSummary(modelToCompGenerateTracker *store.GenerticThreadSafeStore[compGenerateTracker]) {
	for key, val := range modelToCompGenerateTracker.GetAllPairs() {
		Log.Info(fmt.Sprintf("Generated %d components for model [%s] %s", val.totalComps, key, val.version))
		totalAggregateComponents += val.totalComps
		if val.totalComps > 0 {
			totalAggregateModel++
		}
	}

	Log.Info(fmt.Sprintf("-----------------------------\n-----------------------------\nGenerated %d models and %d components", totalAggregateModel, totalAggregateComponents))
}

// This function serves dual purposes: it is invoked either via the UI generation or a spreadsheet. Based on the invocation source, the logger configuration is dynamically set to either output solely to the terminal or to a multi-writer.

// Steps Involved:

//  1. Sheet Parsing: Parse the three sheets and prepare their respective csvHelper instances.
//  2. Registrant-Based Processing:
//     2.1 Meshery: Use the Meshery components sheet as the source.
//     2.2 GitHub/ArtifactHub: Perform additional steps, such as determining the registrant type via the URL and generating a package. This package serves as the basis for component generation (refer to the MeshKit function for details).
//     2.3 File Writing: Write the generated components to their respective paths in the format {model-name/{release-version-in-repo}/defversion/}.
//  3. Handling Relationships: Relationships marked with * as the version are written to all versions, with the filenames updated in the sheet.
//
// If their is ever an error with the writing of file back to spreadsheet of column mismatch just update utils.ComponentCSV struct.
func InvokeGenerationFromSheet(wg *sync.WaitGroup, path string, modelsheetID, componentSheetID int64, spreadsheeetID string, modelName string, modelCSVFilePath, componentCSVFilePath, spreadsheeetCred, relationshipCSVFilePath string, relationshipSheetID int64, srv *sheets.Service) error {
	weightedSem := semaphore.NewWeighted(20)
	url := GoogleSpreadSheetURL + spreadsheeetID
	totalAvailableModels := 0
	spreadsheeetChan := make(chan SpreadsheetData)
	relationshipUpdateChan := make(chan RelationshipCSV)
	defer func() {
		logModelGenerationSummary(modelToCompGenerateTracker)

		// Log.UpdateLogOutput(os.Stdout)
		Log.Info(fmt.Sprintf("Summary: %d models, %d components generated.", totalAggregateModel, totalAggregateComponents))

		Log.Info("See ", logDirPath, " for detailed logs.")

		totalAggregateModel = 0
		totalAggregateComponents = 0
	}()
	modelCSVHelper, err := parseModelSheet(url, modelName, modelsheetID, modelCSVFilePath)
	if err != nil {
		return err
	}

	componentCSVHelper, err := parseComponentSheet(url, modelName, componentSheetID, componentCSVFilePath)
	if err != nil {
		return err
	}
	relationshipCSVHelper, err := parseRelationshipSheet(url, modelName, relationshipSheetID, relationshipCSVFilePath)
	if err != nil {
		return err
	}
	var wgForSpreadsheetUpdate sync.WaitGroup
	wgForSpreadsheetUpdate.Add(1)

	go func() {
		ProcessModelToComponentsMap(componentCSVHelper.Components)
		VerifyandUpdateSpreadsheet(spreadsheeetCred, &wgForSpreadsheetUpdate, srv, spreadsheeetChan, spreadsheeetID, modelCSVFilePath, componentCSVFilePath)
	}()

	var wgForRelationshipUpdates sync.WaitGroup
	wgForRelationshipUpdates.Add(1)
	go func() {
		defer wgForRelationshipUpdates.Done()
		for updatedRelationship := range relationshipUpdateChan {
			// Collect the updated relationships
			relationshipCSVHelper.UpdatedRelationships = append(relationshipCSVHelper.UpdatedRelationships, updatedRelationship)
		}
	}()
	// Iterate models from the spreadsheet
	for _, model := range modelCSVHelper.Models {
		if modelName != "" && modelName != model.Model {
			continue
		}
		totalAvailableModels++
		ctx := context.Background()

		err := weightedSem.Acquire(ctx, 1)
		if err != nil {
			break
		}
		wg.Add(1)
		go func(model ModelCSV) {
			defer func() {
				wg.Done()
				weightedSem.Release(1)
			}()
			var err error

			if utils.ReplaceSpacesAndConvertToLowercase(model.Registrant) == "meshery" {
				err = GenerateDefsForCoreRegistrant(model, componentCSVHelper, path, modelName)
				if err != nil {
					LogError.Error(err)
				}
				return
			}

			generator, err := generators.NewGenerator(model.Registrant, model.SourceURL, model.Model)
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				LogError.Error(err)
				return
			}

			if utils.ReplaceSpacesAndConvertToLowercase(model.Registrant) == "artifacthub" {
				RateLimitArtifactHub()
			}

			pkg, err := generator.GetPackage()
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				LogError.Error(err)
				return
			}

			version := pkg.GetVersion()

			comps, err := pkg.GenerateComponents()
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				LogError.Error(err)
				return
			}
			lengthOfComps := len(comps)
			if lengthOfComps == 0 {
				err = ErrGenerateModel(fmt.Errorf("no components found for model"), model.Model)
				LogError.Error(err)
				return
			}
			modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model, path)
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				LogError.Error(err)
				return
			}
			modelDef, alreadyExist, err := writeModelDefToFileSystem(&model, version, modelDirPath)
			if err != nil {
				err = ErrGenerateModel(err, model.Model)
				LogError.Error(err)
				return
			}
			if alreadyExist {
				totalAvailableModels--
			}
			for _, comp := range comps {
				comp.Version = defVersion
				// Assign the component status corresponding to model status.
				// i.e., If model is enabled, comps are also "enabled". Ultimately, all individual comps will have the ability to control their status.
				// The status "enabled" indicates that the component will be registered inside the registry.
				if modelDef.Metadata == nil {
					modelDef.Metadata = &_model.ModelDefinition_Metadata{}
				}
				if modelDef.Metadata.AdditionalProperties == nil {
					modelDef.Metadata.AdditionalProperties = make(map[string]interface{})
				}

				if comp.Model.Metadata.AdditionalProperties != nil {
					modelDef.Metadata.AdditionalProperties["source_uri"] = comp.Model.Metadata.AdditionalProperties["source_uri"]
				}
				comp.Model = *modelDef

				AssignDefaultsForCompDefs(&comp, modelDef)
				compAlreadyExist, err := comp.WriteComponentDefinition(compDirPath, "json")
				if compAlreadyExist {
					lengthOfComps--
				}
				if err != nil {
					err = ErrGenerateModel(err, model.Model)
					LogError.Error(err)
					return
				}
			}
			if !alreadyExist {
				if len(comps) == 0 {
					err = ErrGenerateModel(fmt.Errorf("no components found for model"), model.Model)
					LogError.Error(err)

				} else {
					log.Info("Current model: ", model.Model)
					log.Info("Extracted ", lengthOfComps, " components for ", model.ModelDisplayName, " (", model.Model, ")")
				}
			} else {
				if len(comps) > 0 {
					log.Info("Model already exists: ", model.Model)
				} else {
					err = ErrGenerateModel(fmt.Errorf("no components found for model"), model.Model)
					LogError.Error(err)
				}
			}
			spreadsheeetChan <- SpreadsheetData{
				Model:      &model,
				Components: comps,
			}

			modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
				totalComps: lengthOfComps,
				version:    version,
			})
		}(model)
	}

	wg.Wait()
	close(spreadsheeetChan)
	ProcessRelationships(relationshipCSVHelper, relationshipUpdateChan, path)
	close(relationshipUpdateChan)
	wgForRelationshipUpdates.Wait()
	err = relationshipCSVHelper.UpdateRelationshipSheet(srv, spreadsheeetCred, spreadsheeetID, relationshipCSVFilePath)
	if err == nil {
		Log.Info("Updated relationship to sheet.")
	}
	wgForSpreadsheetUpdate.Wait()
	return nil
}

// For registrants eg: meshery, whose components needs to be directly created by referencing meshery/schemas repo.
// the sourceURL contains the path of models component definitions
func GenerateDefsForCoreRegistrant(model ModelCSV, ComponentCSVHelper *ComponentCSVHelper, path string, modelName string) error {
	var version string
	parts := strings.Split(model.SourceURL, "/")
	// Assuming the URL is always of the format "protocol://github.com/owner/repo/tree/definitions/{model-name}/version/components"
	// We know the version is the 7th element (0-indexed) in the split URL
	if len(parts) >= 8 {
		version = parts[8] // Fetch the version from the expected position
	} else {
		return fmt.Errorf("invalid SourceURL format: %s", model.SourceURL)
	}
	isModelPublished, _ := strconv.ParseBool(model.PublishToRegistry)
	var compDefComps []component.ComponentDefinition
	alreadyExist := false
	actualCompCount := 0
	defer func() {
		modelToCompGenerateTracker.Set(model.Model, compGenerateTracker{
			totalComps: len(compDefComps) - actualCompCount,
			version:    version,
		})
	}()
	status := entity.Ignored
	if isModelPublished {
		status = entity.Enabled
	}
	_status := component.ComponentDefinitionStatus(status)
	modelDirPath, compDirPath, err := createVersionedDirectoryForModelAndComp(version, model.Model, path)
	if err != nil {
		err = ErrGenerateModel(err, model.Model)
		return err
	}
	modelDef, alreadyExists, err := writeModelDefToFileSystem(&model, version, modelDirPath)
	if err != nil {
		return ErrGenerateModel(err, model.Model)
	}
	isModelPublishToSite, _ := strconv.ParseBool(model.PublishToSites)
	alreadyExist = alreadyExists

	for registrant, models := range ComponentCSVHelper.Components {
		if registrant != "meshery" {
			continue
		}
		for _, comps := range models {
			for _, comp := range comps {
				if comp.Model != model.Model {
					continue
				}
				var componentDef component.ComponentDefinition
				componentDef, err = comp.CreateComponentDefinition(isModelPublishToSite, "v1.0.0")
				if err != nil {
					if LogError != nil {
						LogError.Error(ErrUpdateComponent(err, modelName, comp.Component))
						continue
					}
					Log.Error(ErrUpdateComponent(err, modelName, comp.Component))
					continue
				}
				componentDef.Status = &_status
				componentDef.Model = *modelDef
				alreadyExists, err = componentDef.WriteComponentDefinition(compDirPath, "json")
				if err != nil {
					err = ErrGenerateComponent(err, comp.Model, componentDef.DisplayName)
					if LogError != nil {
						LogError.Error(ErrUpdateComponent(err, modelName, comp.Component))
						continue
					}
					Log.Error(ErrUpdateComponent(err, modelName, comp.Component))
					continue
				}
				if alreadyExists {
					actualCompCount++
				}
				compDefComps = append(compDefComps, componentDef)
			}
		}
	}

	if !alreadyExist {
		if len(compDefComps) == 0 {
			err = ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model)
			return err
		} else if len(compDefComps)-actualCompCount == 0 {
			Log.Info("Current model: ", model.Model)
			Log.Info(" no change in components for ", model.ModelDisplayName, " (", model.Model, ")")
		} else {
			Log.Info("Current model: ", model.Model)
			Log.Info(" extracted ", len(compDefComps)-actualCompCount, " components for ", model.ModelDisplayName, " (", model.Model, ")")
		}
	} else {
		if len(compDefComps) > 0 {
			if len(compDefComps)-actualCompCount == 0 {
				Log.Info("Model already exists: ", model.Model)
			} else {
				Log.Info("Current model: ", model.Model)
				Log.Info(" extracted ", len(compDefComps)-actualCompCount, " components for ", model.ModelDisplayName, " (", model.Model, ")")
			}
		} else {
			err = ErrGenerateModel(fmt.Errorf("no components found for model "), model.Model)
			return err
		}
	}
	if err != nil {
		return ErrGenerateModel(err, model.Model)
	}
	return nil
}
func SetLogger(ismultiWriter bool) error {
	logDirPath := filepath.Join(meshkitutils.GetHome(), ".meshery", "logs", "registry")
	err := os.MkdirAll(logDirPath, 0755)
	if err != nil {
		return err
	}
	logFilePath := filepath.Join(logDirPath, "model-generation.log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		return err
	}

	logErrorFilePath := filepath.Join(logDirPath, "registry-errors.log")
	errorLogFile, err := os.Create(logErrorFilePath)
	if err != nil {
		return err
	}
	if Log == nil {
		Log = SetupMeshkitLogger("mesheryctl", true, logFile)
		LogError = SetupMeshkitLogger("mesheryctl-error", true, errorLogFile)
	}
	Log.SetLevel(logrus.DebugLevel)
	LogError.SetLevel(logrus.ErrorLevel)
	multiWriter := io.MultiWriter(logFile)
	multiErrorWriter := io.MultiWriter(errorLogFile)
	if ismultiWriter {
		multiWriter = io.MultiWriter(os.Stdout, logFile)
		multiErrorWriter = io.MultiWriter(os.Stdout, errorLogFile)
	}

	Log.UpdateLogOutput(multiWriter)
	LogError.UpdateLogOutput(multiErrorWriter)
	return nil
}
