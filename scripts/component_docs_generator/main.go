package main

import (
	"log"
	"os"

	helpers "github.com/layer5io/component_docs_generator/helpers"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: main <CSV_URL>")
	}

	csvURL := os.Args[1]

	helpers.GetIntegrationDocsCSVFile(csvURL)
	csvIndices, _ := helpers.ReadIndexFromJSONFile()
	helpers.FilterRecordsByPublishFlag(csvIndices)
	helpers.CreateIntegrationDocs(csvIndices)
	helpers.DeleteTempCreatedFiles()
}
