package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/layer5io/meshkit/encoding"
	"github.com/layer5io/meshkit/generators"
	"github.com/layer5io/meshkit/generators/models"

	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/csv"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/connection"
	"github.com/meshery/schemas/models/v1beta1/model"
)

var (
	shouldRegisterMod = "publishToSites"
)

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

func (m *ModelCSV) UpdateModelDefinition(modelDef *model.ModelDefinition) error {
	metadata := modelDef.Metadata
	if metadata == nil {
		metadata = &model.ModelDefinition_Metadata{}
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
func (mcv *ModelCSV) CreateModelDefinition(version, defVersion string) model.ModelDefinition {
	status := entity.Ignored
	if strings.ToLower(mcv.PublishToRegistry) == "true" {
		status = entity.Enabled
	}
	var catname category.CategoryDefinition
	catname.Name = mcv.Category
	registrant := createNewRegistrant(mcv.Registrant)

	model := model.ModelDefinition{
		Category:    catname,
		Description: mcv.Description,
		DisplayName: mcv.ModelDisplayName,

		SchemaVersion: v1beta1.ModelSchemaVersion,
		Name:          mcv.Model,
		Status:        model.ModelDefinitionStatus(status),
		Registrant:    registrant,
		SubCategory:   mcv.SubCategory,
		Model: model.Model{
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
		Log.Info("Downloading CSV from: ", sheetURL, csvPath)
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
			Log.Info(fmt.Sprintf("Reading registrant [%s] model [%s]", data.Registrant, data.Model))
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

func (m ModelCSV) CreateMarkDownForMDStyle(componentsMetadata string) string {
	formattedName := utils.FormatName(m.Model)

	var template string = `---
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
	markdown := fmt.Sprintf(template,
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
		m.Feature1,
		m.Feature2,
		m.Feature3,
		m.HowItWorks,
		m.HowItWorksDetails,
	)

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
func AssignDefaultsForCompDefs(componentDef *component.ComponentDefinition, modelDef *model.ModelDefinition) {
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
func GenerateComponentsFromPkg(pkg models.Package, compDirPath string, defVersion string, modelDef model.ModelDefinition) (int, int, error) {
	comps, err := pkg.GenerateComponents()
	if err != nil {
		return 0, 0, err
	}
	lengthOfComps := len(comps)
	for _, comp := range comps {
		comp.Version = defVersion
		if modelDef.Metadata == nil {
			modelDef.Metadata = &model.ModelDefinition_Metadata{}
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
