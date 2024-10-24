package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

// Helper function to convert a string to UPPER_SNAKE_CASE
func toSnakeCase(str string) string {
	snake := strings.ReplaceAll(str, " ", "_")
	snake = strings.ToUpper(snake)
	return snake
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <input_csv>")
		return
	}

	inputFile := os.Args[1]
	outputFile := "permission_constants.js"

	// Open the input CSV file
	f, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer f.Close()

	// Read the CSV file
	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	// Check if the CSV file has at least 2 rows (metadata + header)
	if len(records) < 2 {
		fmt.Println("CSV file does not contain enough rows")
		return
	}

	// Extract header and ensure it has 'Function', 'Key ID', and 'Inserted' columns
	header := records[1] // header is the second row
	functionIdx, keyIDIdx, insertedIdx := -1, -1, -1
	for i, column := range header {
		if column == "Function" {
			functionIdx = i
		} else if column == "Key ID" {
			keyIDIdx = i
		} else if column == "Inserted" {
			insertedIdx = i
		}
	}

	// Validate that 'Function', 'Key ID', and 'Inserted' columns were found
	if functionIdx == -1 || keyIDIdx == -1 || insertedIdx == -1 {
		fmt.Println("CSV does not contain 'Function', 'Key ID', or 'Inserted' columns")
		return
	}

	// Open the output JavaScript file
	output, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		return
	}
	defer output.Close()

	// Start writing the keys.js content
	fmt.Fprintln(output, "export const keys = {")

	// Iterate over the data rows (starting from the third row)
	for _, row := range records[2:] {
		function := row[functionIdx]
		keyID := row[keyIDIdx]
		inserted := row[insertedIdx]

		// Skip records with an empty Key ID or Inserted field
		if keyID == "" || inserted == "" {
			continue
		}

		// Convert the function to UPPER_SNAKE_CASE
		snakeCaseFunction := toSnakeCase(function)

		// Write the object to the output file
		fmt.Fprintf(output, "  \"%s\": { subject: \"%s\", action: \"%s\" },\n", snakeCaseFunction, function, keyID)
	}

	// Close the keys object
	fmt.Fprintln(output, "};")

	fmt.Printf("Successfully generated %s\n", outputFile)
}