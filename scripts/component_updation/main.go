package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/layer5io/component_scraper/pkg"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

var (
	ColumnNamesToExtract = []string{"Project Name", "Helm Chart", "Category", "Sub-Category", "Shape", "Primary Color", "Secondary Color", "Logo URL", "SVG_Color", "SVG_White"}
	PrimaryColumnName    = "Helm Chart"
	OutputPath           = "../../server/output"
)

func main() {

	url := os.Args[1]
	if url == "" {
		log.Fatal("provide a valid URL")
		return
	}
	filep, err := pkg.DownloadCSV(url)
	if err != nil {
		log.Fatal(err)
		return
	}
	file, err := os.Open(filep)
	if err != nil {
		log.Fatal(err)
		return
	}
	csvReader := csv.NewReader(file)
	output, err := pkg.GetEntries(csvReader, ColumnNamesToExtract)
	if err != nil {
		log.Fatal(err)
		return
	}
	file.Close()
	os.Remove(file.Name())

	err = pkg.PopulateEntries(OutputPath, output, PrimaryColumnName, func(dirpath string, changeFields map[string]string) error {
		entries, err := os.ReadDir(dirpath)
		if err != nil {
			return err
		}
		for _, versionentry := range entries {
			if versionentry.IsDir() {
				entries, err := os.ReadDir(filepath.Join(dirpath, versionentry.Name()))
				if err != nil {
					return err
				}
				for _, entry := range entries {
					path, err := filepath.Abs(filepath.Join(dirpath, versionentry.Name(), entry.Name()))
					if err != nil {
						return err
					}
					byt, err := os.ReadFile(path)
					if err != nil {
						return err
					}
					var component v1alpha1.ComponentDefinition
					err = json.Unmarshal(byt, &component)
					if err != nil {
						return err
					}
					if component.Metadata == nil {
						component.Metadata = make(map[string]interface{})
					}
					for key, value := range changeFields {
						if key == "Category" {
							component.Model.Category = value
						} else if key == "Sub-Category" {
							component.Model.SubCategory = value
						} else if isInColumnNames(key) {
							component.Metadata[key] = value
						}
					}
					byt, err = json.Marshal(component)
					if err != nil {
						return err
					}
					err = os.WriteFile(filepath.Join(dirpath, versionentry.Name(), entry.Name()), byt, 0777)
					if err != nil {
						return err
					}
				}
			}

		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
func isInColumnNames(key string) bool {
	for _, n := range ColumnNamesToExtract {
		if n == key {
			return true
		}
	}
	return false
}
