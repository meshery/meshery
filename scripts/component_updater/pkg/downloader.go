package pkg

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadCSV(url string) (string, error) {
	file, err := os.CreateTemp("./", "*.csv")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	byt, _ := io.ReadAll(resp.Body)
	file.WriteString(string(byt))
	byt, _ = io.ReadAll(file)
	path, _ := filepath.Abs(file.Name())
	return path, nil
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

type CSVIndices struct {
	FlagIndex  int `json:"flag-index"`
	SvgIndex   int `json:"svg-index"`
	IndexStart int `json:"index-start"`
	IndexEnd   int `json:"index-end"`
	NameIndex  int `json:"name-index"`
}
