package helpers

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
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

	outputFile, err := os.Create(destinationFile)
	if err != nil {
		log.Println("Error while creating the destination file:", err)
		return
	}
	defer outputFile.Close()

	reader := bufio.NewReader(resp.Body)

	_, err = io.Copy(outputFile, reader)
	if err != nil {
		log.Println("Error while copying the CSV content:", err)
		return
	}

	log.Println("CSV file downloaded and saved as", destinationFile)
}

func ReadIndexFromJSONFile() (CSVIndices, error) {
	var config CSVIndices

	file, err := os.Open("index.json")
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return config, err
	}
	defer file.Close()

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return config, err
	}

	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return config, err
	}
	return config, nil
}

func FilterRecordsByPublishFlag(csvIndices CSVIndices) {
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

	if len(records) > 1 {
		header := records[1]
		if err := csvWriter.Write(header); err != nil {
			log.Println("Error writing the header to the output file:", err)
		}

		publishColumnIndex := csvIndices.FlagIndex // Publish Flag Index
		writeFilteredDataToCSVFile(publishColumnIndex, records, csvWriter)
	} else {
		log.Println("No data records found in the input CSV file.")
		return
	}

}

func writeFilteredDataToCSVFile(publishColumnIndex int, records [][]string, csvWriter *csv.Writer) {
	if len(records) < 3 {
		log.Println("Not enough records to process.")
		return
	}

	for _, record := range records[2:] {
		agValue := record[publishColumnIndex]

		if agValue == "TRUE" {
			if err := csvWriter.Write(record); err != nil {
				log.Println("Error writing to the output file:", err)
			}
		}
	}
}

type CSVIndices struct {
	FlagIndex  int `json:"flag-index"`
	SvgIndex   int `json:"svg-index"`
	IndexStart int `json:"index-start"`
	IndexEnd   int `json:"index-end"`
	NameIndex  int `json:"name-index"`
}
