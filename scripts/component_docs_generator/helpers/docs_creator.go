package helpers

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

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
layout: default
title: ` + name + `
permalink: integrations/` + formattedName + `
type: installation
category: integrations
display-title: "false"
language: en
list: include
image: /assets/img/integrations/` + formattedName + `.svg
---

<h1>{{ page.title }} <img src="{{ page.image }}" style="width: 35px; height: 35px;" /></h1>


<!-- This needs replaced with the Category property, not the sub-category.
 #### Category: ` + record[1] + ` -->

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
