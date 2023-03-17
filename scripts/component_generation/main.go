package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"time"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/artifacthub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gopkg.in/yaml.v3"
)

const dumpFile = "./dump.csv"

var (
	AhSearchEndpoint = artifacthub.AhHelmExporterEndpoint

	OutputDirectoryPath     = "../../server/meshmodel/components"
	ComponentsFileName      = "components.yaml"
	ComponentModelsFileName = "component_models.yaml"
)

type ComponentModel struct {
	SystemName string `yaml:"systemName"`
	AhRepo     string `yaml:"ahRepo"`
	Version    string `yaml:"version"`
	RepoUrl    string `yaml:"repoUrl"`
}

func convertPackagesToCompModels(pkgs []artifacthub.AhPackage) []ComponentModel {
	models := make([]ComponentModel, 0)
	for _, ap := range pkgs {
		models = append(models, ComponentModel{
			SystemName: ap.Name,
			Version:    ap.Version,
			AhRepo:     ap.Repository,
			RepoUrl:    ap.RepoUrl,
		})
	}
	return models
}

func convertCompModelsToPackages(models []ComponentModel) []artifacthub.AhPackage {
	pkgs := make([]artifacthub.AhPackage, 0)
	for _, cm := range models {
		pkgs = append(pkgs, artifacthub.AhPackage{
			Name:       cm.SystemName,
			Version:    cm.Version,
			Repository: cm.AhRepo,
			RepoUrl:    cm.RepoUrl,
		})
	}
	return pkgs
}

func main() {
	if len(os.Args) > 1 {
		spreadsheetID = os.Args[1]
	}
	if _, err := os.Stat(OutputDirectoryPath); err != nil {
		err := os.Mkdir(OutputDirectoryPath, 0744)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	compsFd, err := os.OpenFile(filepath.Join(OutputDirectoryPath, ComponentsFileName), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = SplitYamlIntoFiles(compsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	compsWriter := Writer{
		file: compsFd,
	}
	modelsFd, err := os.OpenFile(filepath.Join(OutputDirectoryPath, ComponentModelsFileName), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	modelsWriter := Writer{
		file: modelsFd,
	}
	defer modelsFd.Close()
	defer compsFd.Close()
	// move to a new function: getHelmPackages
	content := make([]byte, 0)
	models := make([]ComponentModel, 0)
	pkgs := make([]artifacthub.AhPackage, 0)
	content, err = io.ReadAll(modelsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(content, &models)
	if err != nil {
		fmt.Println(err)
	}
	pkgs = convertCompModelsToPackages(models)
	if pkgs == nil || len(pkgs) == 0 {
		pkgs, err = artifacthub.GetAllAhHelmPackages()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = writeComponentModels(convertPackagesToCompModels(pkgs), &modelsWriter)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	inputChan := make(chan []artifacthub.AhPackage)
	csvChan := make(chan string, 50)
	f, err := os.Create(dumpFile)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
	}

	f.Write([]byte("model,component_count,components\n"))
	go func() {
		for entry := range csvChan {
			f.Write([]byte(entry))
		}
	}()
	srv := NewSheetSRV()
	// Convert sheet ID to sheet name.
	response1, err := srv.Spreadsheets.Get(spreadsheetID).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		fmt.Println(err)
		return
	}
	sheetName := ""
	for _, v := range response1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetID) {
			sheetName = prop.Title
			break
		}
	}
	spreadsheetChan := make(chan struct {
		comps []v1alpha1.ComponentDefinition
		model string
	}, 100)
	// Set the range of cells to retrieve.
	rangeString := sheetName + "!B:G" //modelname column

	// Get the value of the specified cell.
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from sheet: ", err)
		return
	}
	availableModels := make(map[string]bool)
	availableComponentsPerModel := make(map[string]map[string]bool)
	for _, val := range resp.Values {
		if len(val) > 1 {
			key := val[0].(string)
			if key == "" {
				continue
			}
			var compkey string
			if len(val) > 5 {
				compkey = val[5].(string)
			}
			if compkey == "" {
				availableModels[key] = true
				continue
			}
			if availableComponentsPerModel[key] == nil {
				availableComponentsPerModel[key] = make(map[string]bool)
			}
			availableComponentsPerModel[key][compkey] = true
		}
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		Spreadsheet(srv, sheetName, spreadsheetChan, availableModels, availableComponentsPerModel)
	}()

	for i := 0; i <= 10; i++ {
		StartPipeline(inputChan, csvChan, &compsWriter, spreadsheetChan)
	}
	inputChan <- pkgs[:10]
	for len(pkgs) != 0 {
		if len(pkgs) < 50 {
			inputChan <- pkgs
			time.Sleep(10 * time.Second)
			return
		}
		inputChan <- pkgs[:50]
		pkgs = pkgs[50:]
	}
	time.Sleep(20 * time.Second)

	// split files
	err = SplitYamlIntoFiles(compsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = os.Remove(filepath.Join(OutputDirectoryPath, ComponentsFileName))
	if err != nil {
		fmt.Println(err)
		return
	}
	close(spreadsheetChan)
	wg.Wait()
}

func StartPipeline(in chan []artifacthub.AhPackage, csv chan string, writer *Writer, spreadsheet chan struct {
	comps []v1alpha1.ComponentDefinition
	model string
}) error {
	pkgsChan := make(chan []artifacthub.AhPackage)
	compsChan := make(chan struct {
		comps []v1alpha1.ComponentDefinition
		model string
	})
	compsCSV := make(chan struct {
		comps []v1alpha1.ComponentDefinition
		model string
	})
	// updating pacakge data
	go func() {
		for pkgs := range in {
			ahPkgs := make([]artifacthub.AhPackage, 0)
			for _, ap := range pkgs {
				fmt.Println("[DEBUG] Updating package data for: ", ap.Name)
				err := ap.UpdatePackageData()
				if err != nil {
					fmt.Println(err)
					continue
				}
				ahPkgs = append(ahPkgs, ap)
			}
			pkgsChan <- ahPkgs
		}
	}()
	// generation of components
	go func() {
		for pkgs := range pkgsChan {
			for _, ap := range pkgs {
				fmt.Println("[DEBUG] Generating components for: ", ap.Name)
				comps, err := ap.GenerateComponents()
				if err != nil {
					fmt.Println(err)
					continue
				}
				compsCSV <- struct {
					comps []v1alpha1.ComponentDefinition
					model string
				}{
					comps: comps,
					model: ap.Name,
				}
				spreadsheet <- struct {
					comps []v1alpha1.ComponentDefinition
					model string
				}{
					comps: comps,
					model: ap.Name,
				}
			}

		}
	}()
	// writer
	go func() {
		for modelcomps := range compsChan {
			err := writeComponents(modelcomps.comps, writer)
			if err != nil {
				fmt.Println(err)
			}
			compsCSV <- struct {
				comps []v1alpha1.ComponentDefinition
				model string
			}{
				comps: modelcomps.comps,
				model: modelcomps.model,
			}
		}
	}()
	if _, err := os.Stat(dumpFile); os.IsExist(err) {
		// If file exists, delete it
		err := os.Remove(dumpFile)
		if err != nil {
			fmt.Printf("Error deleting file: %s\n", err)
		}
	}

	go func() {
		for comps := range compsCSV {
			count := len(comps.comps)
			names := "\""
			for _, cmp := range comps.comps {
				names += fmt.Sprintf("%s,", cmp.Kind)
			}
			names = strings.TrimSuffix(names, ",")
			names += "\""
			if count > 0 {
				model := comps.model
				fmt.Println(fmt.Sprintf("[DEBUG]Adding to CSV: %s", model))
				csv <- fmt.Sprintf("%s,%d,%s\n", model, count, names)
			}
		}
	}()

	return nil
}

type ComponentStruct struct {
	PackageName string                         `yaml:"name"`
	Components  []v1alpha1.ComponentDefinition `yaml:"components"`
}

func wrapComponentsInWritableStruct(comps []v1alpha1.ComponentDefinition, pkgName string) ComponentStruct {
	return ComponentStruct{
		PackageName: pkgName,
		Components:  comps,
	}

}

type Writer struct {
	file *os.File
	m    sync.Mutex
}

func writeComponentModels(models []ComponentModel, writer *Writer) error {
	writer.m.Lock()
	defer writer.m.Unlock()
	val, err := yaml.Marshal(models)
	if err != nil {
		return err
	}
	_, err = writer.file.Write(val)
	if err != nil {
		return err
	}
	return nil
}

func writeComponents(cmps []v1alpha1.ComponentDefinition, writer *Writer) error {
	// writer.m.Lock()
	// defer writer.m.Unlock()
	for _, comp := range cmps {
		path := filepath.Join(OutputDirectoryPath, comp.Model.Name)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Println("created directory ", comp.Model.Name)
		}
		path = filepath.Join(path, comp.Model.Version)
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(path, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Println("created versioned directory ", comp.Model.Version)
		}
		f, err := os.Create(filepath.Join(path, comp.Kind+".json"))
		if err != nil {
			return err
		}
		byt, err := json.Marshal(comp)
		if err != nil {
			return err
		}
		_, err = f.Write(byt)
		if err != nil {
			return err
		}
	}
	return nil
	// compsToWrite := make([]ComponentStruct, 0)
	// content, err := io.ReadAll(writer.file)
	// if err != nil {
	// 	return err
	// }
	// // 1. the file is empty
	// if string(content) == "" {
	// 	for _, cs := range cmps {
	// 		if len(cs.Components) != 0 {
	// 			compsToWrite = append(compsToWrite, cs)
	// 		}
	// 	}
	// }
	// // 2. the file already has some components in it
	// if string(content) != "" {
	// 	err = yaml.Unmarshal(content, &compsToWrite)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	for _, cs := range cmps {
	// 		if len(cs.Components) != 0 {
	// 			fmt.Println("[DEBUG] Writing components for package: ", cs.PackageName)
	// 			compsToWrite = append(compsToWrite, cs)
	// 		}
	// 	}
	// }
	// if len(compsToWrite) == 0 {
	// 	return nil
	// }
	// // write components
	// val, err := yaml.Marshal(compsToWrite)
	// if err != nil {
	// 	return err
	// }
	// _, err = writer.file.Write(val)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

// this function should take in the file descriptor for a yaml
// file that contains array of items and split that into multiple files
// with each item having certain number of items
// TODO: Refactor
func SplitYamlIntoFiles(file *os.File) error {
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var list []ComponentStruct
	err = yaml.Unmarshal(fileContent, &list)
	if err != nil {
		return err
	}

	fmt.Println(list)
	result := make([]([]ComponentStruct), 0)
	dummy := make([]ComponentStruct, 0)
	for i, comp := range list {
		dummy = append(dummy, comp)
		if (i+1)%30 == 0 || i+1 == len(list) {
			result = append(result, dummy)
			dummy = make([]ComponentStruct, 0)
		}
	}
	for i, fileContent := range result {
		file, err := os.Create(fmt.Sprintf("%s/components%d.yaml", OutputDirectoryPath, i+1))
		if err != nil {
			return err
		}
		out, err := yaml.Marshal(fileContent)
		if err != nil {
			return err
		}
		_, err = file.Write(out)
		if err != nil {
			return err
		}
	}

	return nil
}

var spreadsheetID = "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw"

const sheetID = 0

func Spreadsheet(srv *sheets.Service, sheetName string, spreadsheet chan struct {
	comps []v1alpha1.ComponentDefinition
	model string
}, am map[string]bool, acpm map[string]map[string]bool) {
	start := time.Now()
	rangeString := sheetName + "!A4:AB4"
	for entry := range spreadsheet {
		if len(entry.comps) == 0 {
			continue
		}
		values := [][]interface{}{}
		for _, comp := range entry.comps {
			if acpm[entry.model][comp.Kind] {
				fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model, comp.Kind)
				continue
			}
			// Get the value of the specified cell.
			resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
			if err != nil {
				fmt.Println("Unable to retrieve data from sheet: ", err)
				return
			}
			newValues := resp.Values[0]
			newValues[6] = comp.Kind
			newValues[1] = entry.model
			newValues[0] = entry.model
			values = append(values, newValues)
			time.Sleep(1 * time.Second) //To keep rps under 1rps
		}
		if am[entry.model] {
			fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model)
			continue
		}
		// Get the value of the specified cell.
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
		if err != nil {
			fmt.Println("Unable to retrieve data from sheet: ", err)
			return
		}
		newValues := resp.Values[0]
		newValues[1] = entry.model
		newValues[0] = entry.model
		newValues[4] = len(entry.comps)
		values = append(values, newValues)
		row := &sheets.ValueRange{
			Values: values,
		}
		// // srv.Spreadsheets.Values.Get(spreadsheetID, sheetName)
		// // srv.Spreadsheets.Values.Update()
		response2, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
		if err != nil || response2.HTTPStatusCode != 200 {
			fmt.Println(err)
			continue
		}
		time.Sleep(1 * time.Second) //To keep rps under 1rps
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Time taken by spreadsheet updater (including the time it required to generate components): %f", elapsed.Minutes())
}

func NewSheetSRV() *sheets.Service {
	ctx := context.Background()
	byt, _ := base64.StdEncoding.DecodeString(os.Getenv("CRED"))
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(byt, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Println("ERR2", err)
		return nil
	}
	// create client with config and context
	client := config.Client(ctx)
	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("ERR3", err)
		return nil
	}
	return srv
}
