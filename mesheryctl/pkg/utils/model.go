package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	mutils "github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	"github.com/layer5io/meshkit/utils/csv"
)

var (
	shouldRegisterMod        = "publishToSites"
	shouldRegisterToRegsitry = "publishToRegistry"
)

type ModelCSV struct {
	ModelDisplayName   string `json:"modelDisplayName"`
	Model              string `json:"model"`
	Registrant         string `json:"registrant"`
	Category           string `json:"category"`
	SubCategory        string `json:"subCategory"`
	Description        string `json:"description"`
	SourceURL          string `json:"sourceURL"`
	Website            string `json:"website"`
	Docs               string `json:"docs"`
	Shape              string `json:"shape"`
	PrimaryColor       string `json:"primaryColor"`
	SecondaryColor     string `json:"secondaryColor"`
	StyleOverrides     string `json:"styleOverrides"`
	Styles             string `json:"styles"`
	ShapePolygonPoints string `json:"shapePolygonPoints"`
	DefaultData        string `json:"defaultData"`
	Capabilities       string `json:"capabilities"`
	LogoURL            string `json:"logoURL"`
	SVGColor           string `json:"svgColor"`
	SVGWhite           string `json:"svgWhite"`
	SVGComplete        string `json:"svgComplete"`
	PublishToRegistry  string `json:"publishToRegistry"`
	IsAnnotation       string `json:"isAnnotation"`
	AboutProject       string `json:"aboutProject"`
	PageSubtTitle      string `json:"pageSubtitle"`
	DocsURL            string `json:"docsURL"`
	StandardBlurb      string `json:"standardBlurb"`
	Feature1           string `json:"feature1"`
	Feature2           string `json:"feature2"`
	Feature3           string `json:"feature3"`
	HowItWorks         string `json:"howItWorks"`
	HowItWorksDetails  string `json:"howItWorksDetails"`
	Screenshots        string `json:"screenshots"`
	FullPage           string `json:"fullPage"`
	PublishToSites     string `json:"publishToSites"`
}

func (mcv *ModelCSV) CreateModelDefinition(version string) v1alpha1.Model {
	var isModelPublished, isAnnotation bool
	if strings.ToLower(mcv.PublishToRegistry) == "true" {
		isModelPublished = true
	}
	if strings.ToLower(mcv.IsAnnotation) == "true" {
		isAnnotation = true
	}

	model := v1alpha1.Model{
		Name:        mcv.Model,
		Version:     version,
		DisplayName: mcv.ModelDisplayName,
		Metadata: map[string]interface{}{
			"isAnnotation": isAnnotation,
			"svgColor":     mcv.SVGColor,
			"svgWhite":     mcv.SVGWhite,
			"svgComplete":  mcv.SVGComplete,
			"subCategory":  mcv.SubCategory, // move as first class attribute
			"published":    isModelPublished,
		},
		Category: v1alpha1.Category{
			Name: mcv.Category,
		},
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
	fmt.Println("Downloading CSV from:", sheetURL)
	dirPath := filepath.Join(utils.GetHome(), ".meshery", "content")
	_ = os.MkdirAll(dirPath, 0755)
	csvPath := filepath.Join(dirPath, "models.csv")
	err := utils.DownloadFile(csvPath, sheetURL)
	if err != nil {
		return nil, err
	}

	return &ModelCSVHelper{
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: sheetURL,
		Models:         []ModelCSV{},
		CSVPath:        csvPath,
		Title:          spreadsheetName,
	}, nil
}

func (mch *ModelCSVHelper) ParseModelsSheet(parseForDocs bool) error {
	ch := make(chan ModelCSV, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[ModelCSV](mch.CSVPath, rowIndex, nil, func(columns []string, currentRow []string) bool {
		index := 0

		if parseForDocs {
			index = GetIndexForRegisterCol(columns, shouldRegisterMod)
		} else {
			index = GetIndexForRegisterCol(columns, shouldRegisterToRegsitry)
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
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			errorChan <- err
		}
	}()
	for {
		select {

		case data := <-ch:
			data.Registrant = mutils.ReplaceSpacesAndConvertToLowercase(data.Registrant)
			mch.Models = append(mch.Models, data)
			fmt.Printf("Reading Modal: %s from Registrant: %s\n", data.Model, data.Registrant)
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
layout: enhanced
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
display-title: "false"
---
<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>

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
		m.AboutProject,
		m.StandardBlurb,
	)

	markdown = strings.ReplaceAll(markdown, "\r", "\n")

	return markdown
}

func (m ModelCSVHelper) Cleanup() error {
	// remove csv file
	fmt.Println("Removing CSV file:", m.CSVPath)
	err := os.Remove(m.CSVPath)
	if err != nil {
		return err
	}

	return nil
}
