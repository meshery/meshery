package helpers

import (
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"os"
)

func GetIntegrationDocsCSVFile(csvURL string) {

	destinationFile := "integration.csv"

	resp, err := http.Get(csvURL)
	if err != nil {
		log.Println("Error while fetching the CSV:", err)
		return
	}
	defer resp.Body.Close()

	csvContent, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the CSV:", err)
		return
	}

	err = os.WriteFile(destinationFile, csvContent, 0644)
	if err != nil {
		log.Println("Error while saving the CSV to the destination file:", err)
		return
	}

	log.Println("CSV file downloaded and saved as", destinationFile)
}

func FilterRecordsByPublishFlag() {
	inputFile, err := os.Open("integration.csv")
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

	outputFile, err := os.Create("filtered.csv")
	if err != nil {
		log.Println("Error creating the output file:", err)
		return
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	header := records[0]
	if err := csvWriter.Write(header); err != nil {
		log.Println("Error writing the header to the output file:", err)
	}

	publishColumnIndex := 32 // AG Column index is 32

	writeFilteredDataToCSVFile(publishColumnIndex, records, csvWriter)

}

func writeFilteredDataToCSVFile(publishColumnIndex int, records [][]string, csvWriter *csv.Writer) {
	for _, record := range records[1:] {
		agValue := record[publishColumnIndex]

		if agValue == "TRUE" {
			if err := csvWriter.Write(record); err != nil {
				log.Println("Error writing to the output file:", err)
			}
		}
	}
}
