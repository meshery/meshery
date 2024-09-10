package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/layer5io/meshkit/encoding"

	"github.com/layer5io/meshkit/models/meshmodel/entity"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/csv"
	"github.com/meshery/schemas/models/v1alpha1/capability"
	"github.com/meshery/schemas/models/v1beta1"
	"github.com/meshery/schemas/models/v1beta1/category"
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

func NewModelCSVHelper(sheetURL, spreadsheetName string, spreadsheetID int64) (*ModelCSVHelper, error) {
	sheetURL = sheetURL + "/pub?output=csv" + "&gid=" + strconv.FormatInt(spreadsheetID, 10)
	Log.Info("Downloading CSV from: ", sheetURL)
	dirPath := filepath.Join(utils.GetHome(), ".meshery", "content")
	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return nil, utils.ErrCreateDir(err, dirPath)
	}
	csvPath := filepath.Join(dirPath, "models.csv")
	err = utils.DownloadFile(csvPath, sheetURL)
	if err != nil {
		return nil, utils.ErrReadingRemoteFile(err)
	}

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
		`../_images/meshmap-visualizer.png`,
		`../_images/meshmap-designer.png`,
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
