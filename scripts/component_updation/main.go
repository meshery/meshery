package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/layer5io/component_scraper/pkg"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

var (
	ColumnNamesToExtract = []string{"Project Name", "Helm Chart", "Category", "Sub-Category", "Shape", "Primary Color", "Secondary Color", "Logo URL"}
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
		for _, entry := range entries {
			f, err := os.OpenFile(filepath.Join(dirpath, entry.Name()), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				return err
			}
			byt, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			var component v1alpha1.ComponentDefinition
			err = json.Unmarshal(byt, &component)
			if err != nil {
				return err
			}
			if component.Metadata.Metadata == nil {
				component.Metadata.Metadata = make(map[string]interface{})
			}
			for key, value := range changeFields {
				component.Metadata.Metadata[key] = value
			}
			byt, err = json.Marshal(component)
			if err != nil {
				return err
			}
			_, err = f.Write(byt)
			if err != nil {
				return err
			}
			f.Close()
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
