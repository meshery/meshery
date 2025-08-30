package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Integration sheet details
	SpreadsheetID = "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw"
	ModelsGID     = 0 // Usually the first sheet in Google Sheets
	BaseURL       = "https://docs.google.com/spreadsheets/d/" + SpreadsheetID + "/export?format=csv&gid="
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: generate_test_csv <output_directory>")
		os.Exit(1)
	}

	outputDir := os.Args[1]

	// Create output directory
	err := os.MkdirAll(outputDir, 0755)
	if err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		os.Exit(1)
	}

	// Download and extract first three complete model rows from integration sheet
	headers, modelRows, err := getFirstThreeModels()
	if err != nil {
		fmt.Printf("Error getting models from integration sheet: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Downloaded %d complete model rows from integration sheet\n", len(modelRows))

	// Generate all required CSV files
	err = createTestModelCSV(filepath.Join(outputDir, "models.csv"), headers, modelRows)
	if err != nil {
		fmt.Printf("Error generating models CSV: %v\n", err)
		os.Exit(1)
	}

	err = createMinimalComponentsCSV(filepath.Join(outputDir, "components.csv"))
	if err != nil {
		fmt.Printf("Error generating components CSV: %v\n", err)
		os.Exit(1)
	}

	err = createMinimalRelationshipsCSV(filepath.Join(outputDir, "relationships.csv"))
	if err != nil {
		fmt.Printf("Error generating relationships CSV: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully generated test CSV files in: %s\n", outputDir)
}

func getFirstThreeModels() ([]string, [][]string, error) {
	url := BaseURL + fmt.Sprintf("%d", ModelsGID)
	fmt.Printf("Downloading models sheet from: %s\n", url)

	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; TestBot/1.0)")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to download models sheet: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed to download models sheet: HTTP %d", resp.StatusCode)
	}

	// Parse CSV to extract first three complete model rows
	reader := csv.NewReader(resp.Body)

	// Read header row
	headers, err := reader.Read()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read CSV header: %v", err)
	}

	var modelRows [][]string
	rowCount := 0

	for {
		record, err := reader.Read()
		if err != nil {
			break // EOF or error
		}

		// Skip empty rows or rows with empty model names
		if len(record) > 1 && strings.TrimSpace(record[1]) != "" && strings.TrimSpace(record[1]) != "model" {
			modelRows = append(modelRows, record)
			rowCount++
			if rowCount >= 3 {
				break // We only need first 3 models
			}
		}
	}

	if len(modelRows) == 0 {
		return nil, nil, fmt.Errorf("no models found in the sheet")
	}

	// Ensure we have at least 3 models, duplicate the last one if necessary
	for len(modelRows) < 3 {
		modelRows = append(modelRows, modelRows[len(modelRows)-1])
	}

	return headers, modelRows[:3], nil // Return headers and exactly 3 model rows
}

func createTestModelCSV(outputPath string, headers []string, modelRows [][]string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create models CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// Write the actual model data from integration sheet
	for _, modelRow := range modelRows {
		if err := writer.Write(modelRow); err != nil {
			return fmt.Errorf("failed to write model record: %v", err)
		}
	}

	return nil
}

func createMinimalComponentsCSV(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create components CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write minimal headers for components
	headers := []string{
		"registrant", "model", "component", "description", "shape",
		"primaryColor", "secondaryColor", "svgColor", "svgWhite", "svgComplete",
		"schema", "docs", "styleOverrides", "styles", "shapePolygonPoints",
		"defaultData", "capabilities", "logoURL", "genealogy", "isAnnotation", "version", "status",
	}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// Write minimal component data (just headers, no actual data)
	return nil
}

func createMinimalRelationshipsCSV(outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create relationships CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write minimal headers for relationships
	headers := []string{"kind", "type", "description", "subType", "evaluator", "selector"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// Write minimal relationship data (just headers, no actual data)
	return nil
}
