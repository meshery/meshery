package pkg

import (
	"encoding/csv"
	"bytes"
	"encoding/xml"
	"log"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const template string = `---
title: <model-display-name>
subtitle: <Page Subtitle>
integrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh.svg
darkModeIntegrationIcon: ../../../assets/images/service-mesh-icons/aws-app-mesh_white.svg
docURL: <Docs URL>
category: <Category>
subcategory: <Sub-Category>
featureList: [<Feature 1>,<Feature 2>,<Feature 3>]
workingSlides: [
    ../_images/meshmap-visualizer.png,
    ../_images/meshmap-designer.png]
howItWorks: <howItWorks>
howItWorksDetails: howItWorksDetails
published: <Publish>
---
<p>
   <About Project>
</p>
<p>
   <Standard Blurb>
</p>`

// func createEmptyMarkdown(path string) error {
//    file, err := os.Create(path)
//    if err != nil {
//       return err
//    }
//    _, err = file.Write([]byte(template))
//    return err
// }

type TemplateAttributes struct {
	Title                   string
	ModelName               string
	Subtitle                string
	DocURL                  string
	Category                string
	Subcategory             string
	FeatureList             string
	HowItWorks              string
	HowItWorksDetails       string
	AboutProject            string
	StandardBlurb           string
	WorkingSlides           string
	Published               string
	IntegrationIcon         string
	DarkModeIntegrationIcon string
	FullPage                string
	ColorSVG                string
	WhiteSVG                string
}

func (t TemplateAttributes) CreateMarkDown() string {
	// markdown := "---\n"
	// markdown += "title: " + t.Title + "\n"
	// markdown += "subtitle: " + t.Subtitle + "\n"
	// markdown += "integrationIcon: " + t.IntegrationIcon + "\n"
	// markdown += "darkModeIntegrationIcon: " + t.DarkModeIntegrationIcon + "\n"
	// markdown += "docURL: " + t.DocURL + "\n"
	// markdown += "category: " + t.Category + "\n"
	// markdown += "subcategory: " + t.Subcategory + "\n"
	// markdown += "featureList: " + t.FeatureList + "\n"
	// markdown += "workingSlides: " + t.WorkingSlides + "\n"
	// markdown += "howItWorks: " + t.HowItWorks + "\n"
	// markdown += "howItWorksDetails: " + t.HowItWorksDetails + "\n"
	// markdown += "published: " + t.Published + "\n"
	// markdown += "---\n"
	// markdown += t.AboutProject + "\n"
	// markdown += t.StandardBlurb
	markdown := t.FullPage
	markdown = strings.ReplaceAll(markdown, "\r", "\n")
	return markdown
}

// Creates JSON formatted meshmodel attribute item for Meshery docs
func (t TemplateAttributes) CreateJSONItem() string {
	json := "{"
	json += fmt.Sprintf("\"name\":\"%s\"", t.Title)
	// If SVGs exist, then add the paths to json
	if t.ColorSVG != "" {
		json += fmt.Sprintf(",\"color\":\"../images/integration/%s-color.svg\"", t.ModelName)
	}

	if t.WhiteSVG != "" {
		json += fmt.Sprintf(",\"white\":\"../images/integration/%s-white.svg\"", t.ModelName)
	}

	json += fmt.Sprintf(",\"permalink\":\"https://docs.meshery.io/extensibility/integrations/%s\"", FormatName(t.Title))

	json += "}"
	return json
}

func FormatName(input string) string {
	formattedName := strings.ReplaceAll(input, " ", "-")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

const XMLTAG = "<?xml version=\"1.0\" encoding=\"UTF-8\"?><!DOCTYPE svg>"

func WriteToFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
func WriteSVG(path string, svg string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.WriteString(svg)
	if err != nil {
		return err
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

// UpdateSVGString updates the width and height attributes of an SVG file and returns the modified SVG as a string.
func UpdateSVGString(svgStr string, width, height int) (string, error) {
	// Create a reader for the SVG string.
	r := strings.NewReader(svgStr)

	// Create a decoder for the SVG string.
	d := xml.NewDecoder(r)

	// Create a buffer for the modified SVG string.
	var b bytes.Buffer

	// Create an encoder for the buffer.
	e := xml.NewEncoder(&b)

	// Iterate through the tokens in the SVG string.
	for {
		// Read the next token.
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		// If the token is an element name, check if it is an "svg" element.
		if se, ok := t.(xml.StartElement); ok {
			if se.Name.Local == "svg" {
				// Set the width and height attributes to the desired values.
				updatedH := false
				updatedW := false
				xmlnsindex := -1
				for i, a := range se.Attr {
					if a.Name.Local == "width" {
						se.Attr[i].Value = strconv.Itoa(width)
						updatedW = true
					} else if a.Name.Local == "height" {
						se.Attr[i].Value = strconv.Itoa(height)
						updatedH = true
					}
					if a.Name.Local == "xmlns" {
						xmlnsindex = i
					}
				}
				if !updatedH {
					se.Attr = append(se.Attr, xml.Attr{
						Name: xml.Name{
							Local: "height",
						},
						Value: strconv.Itoa(height),
					})
				}
				if !updatedW {
					se.Attr = append(se.Attr, xml.Attr{
						Name: xml.Name{
							Local: "width",
						},
						Value: strconv.Itoa(width),
					})
				}
				if xmlnsindex > -1 {
					se.Attr = append(se.Attr[0:xmlnsindex], se.Attr[xmlnsindex+1:]...)
				}
			} else {
				for i, a := range se.Attr {
					xmlnsindex := -1
					nahbro := 0
					if a.Name.Local == "xmlns" {
						fmt.Println("found at ", i)
						fmt.Println(a.Name)
						xmlnsindex = i
						nahbro++
					}
					if xmlnsindex > -1 {
						se.Attr = append(se.Attr[0:xmlnsindex], se.Attr[xmlnsindex+1:]...)
					}
				}
			}
			t = se
		}
		// Write the modified token to the buffer.
		if err := e.EncodeToken(t); err != nil {
			return "", err
		}
	}

	// Flush the encoder's buffer to the buffer.
	if err := e.Flush(); err != nil {
		return "", err
	}
	var svg string
	if b.String() != "" {
		svg = XMLTAG + b.String()
	}
	return svg, nil
}


func CreateIntegrationDocs(csvIndices CSVIndices) {
	inputFile, err := os.Open("filtered.csv")
	if err != nil {
		log.Println("Error opening input file:", err)
		return
	}
	defer inputFile.Close()

	csvReader := csv.NewReader(inputFile)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Println("Error reading input file:", err)
		return
	}
	GenerateIntegrationDocsSVG(records, csvIndices)
	GenerateIntegrationDocs(records, csvIndices)
}

func GenerateIntegrationDocsSVG(records [][]string, csvIndices CSVIndices) {
	for _, record := range records[1:] {
		svgContent := record[csvIndices.SvgIndex]
		fmt.Println("Creating SVG:", record[csvIndices.NameIndex])
		createFiles("../../docs/assets/img/integrations", ".svg", record[csvIndices.NameIndex], svgContent)
	}
}

func GenerateIntegrationDocs(records [][]string, csvIndices CSVIndices) {

	for _, record := range records[1:] {

		mdContent := generateMDContent(record, record[csvIndices.NameIndex], csvIndices)

		createFiles("../../docs/pages/integrations", ".md", record[csvIndices.NameIndex], mdContent)
	}
}

func createFiles(path, filetype, name, content string) {
	formattedName := formatName(name)

	fullPath := path + "/" + formattedName + filetype

	file, err := os.Create(fullPath)
	fmt.Println("Creating file:", fullPath)
	if err != nil {
		log.Println("Error creating filetype file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		log.Println("Error writing to filetype file:", err)
	}
}

func formatName(input string) string {
	formattedName := strings.ReplaceAll(input, " ", "-")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

func generateMDContent(record []string, name string, csvIndices CSVIndices) string {
	mdContent := ""
	for i := csvIndices.IndexStart; i <= csvIndices.IndexEnd; i++ {
		if record[i] != "" && !strings.Contains(record[i], "https") {
			mdContent += strconv.Itoa(i-(csvIndices.IndexStart-1)) + `.` + ` ` + record[i] + "\n\n"
		}
	}

	formattedName := formatName(name)

	// Remove <p> and </p> tags and replace them with one line gap
	overviewAndFeatures := strings.Replace(mdContent, "<p>", "", -1)
	overviewAndFeatures = strings.Replace(overviewAndFeatures, "</p>", "\n", -1)

	content := `---
layout: enhanced
title: ` + name + `
permalink: extensibility/integrations/` + formattedName + `
type: extensibility
category: integrations
integrations-category: ` + record[2] + `
integrations-subcategory: ` + record[3] + `
display-title: "false"
language: en
list: include
image: /assets/img/integrations/` + formattedName + `.svg
---

<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>


<!-- This needs replaced with the Category property, not the sub-category.
 #### About: ` + record[25] + ` -->

### Overview & Features:

` + overviewAndFeatures

	return content
}

func DeleteTempCreatedFiles() {
	files := []string{"filtered.csv", "integration.csv"}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			log.Printf("Error deleting %s: %v\n", file, err)
		} else {
			log.Printf("%s deleted successfully.\n", file)
		}
	}
}
