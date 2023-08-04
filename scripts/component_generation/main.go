package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
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
const COLUMNRANGE = "!A:AF3" //Update this on addition of new columns

var NameToIndex = map[string]int{ //Update this on addition of new columns
	"modelDisplayName":  0,
	"model":             1,
	"category":          2,
	"subCategory":       3,
	"CRDs":              4,
	"link":              5,
	"hasSchema?":        6,
	"component":         7,
	"shape":             8,
	"primaryColor":      9,
	"secondaryColor":    10,
	"styleOverrides":    11,
	"logoURL":           12,
	"svgColor":          13,
	"svgWhite":          14,
	"svgComplete":       15,
	"genealogy":         16,
	"About Project":     17,
	"Page Subtitle":     18,
	"Docs URL":          19,
	"Standard Blurb":    20,
	"Feature 1":         21,
	"Feature 2":         22,
	"Feature 3":         23,
	"howItWorks":        24,
	"howItWorksDetails": 25,
	"Screenshots":       26,
	"Full Page":         27,
	"Publish?":          28,
}
var (
	AhSearchEndpoint = artifacthub.AhHelmExporterEndpoint

	OutputDirectoryPath     = "../../server/meshmodel"
	ComponentModelsFileName = path.Join(OutputDirectoryPath, "component_models.yaml")
)


var priorityRepos = map[string]bool{"prometheus-community": true, "grafana": true} //Append ahrepos here whose components should be respected and should be used when encountered duplicates

// returns pkgs with sorted pkgs at the front
func sortOnVerified(pkgs []artifacthub.AhPackage) (verified []artifacthub.AhPackage, official []artifacthub.AhPackage, cncf []artifacthub.AhPackage, priority []artifacthub.AhPackage, unverified []artifacthub.AhPackage) {
	for _, pkg := range pkgs {
		if priorityRepos[pkg.Repository] {
			priority = append(priority, pkg)
			continue
		}
		if pkg.CNCF {
			cncf = append(cncf, pkg)
		} else if pkg.Official {
			official = append(official, pkg)
		} else if pkg.VerifiedPublisher {
			verified = append(verified, pkg)
		} else {
			unverified = append(unverified, pkg)
		}
	}
	return
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
	
	modelsFd, err := os.OpenFile(ComponentModelsFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	modelsWriter := Writer{
		file: modelsFd,
	}
	defer modelsFd.Close()
	// move to a new function: getHelmPackages
	pkgs := make([]artifacthub.AhPackage, 0)
	content, err := io.ReadAll(modelsFd)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = yaml.Unmarshal(content, &pkgs)
	if err != nil {
		fmt.Println(err)
	}

	if len(pkgs) == 0 {
		pkgs, err = artifacthub.GetAllAhHelmPackages()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = writeComponentModels(pkgs, &modelsWriter)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	verified, official, cncf, priority, unverified := sortOnVerified(pkgs)
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
		comps   []v1alpha1.ComponentDefinition
		model   string
		helmURL string
	}, 100)
	// Set the range of cells to retrieve.
	rangeString := sheetName + COLUMNRANGE

	// Get the value of the specified cell.
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from sheet: ", err)
		return
	}
	availableModels := make(map[string][]interface{})
	availableComponentsPerModel := make(map[string]map[string]bool)
	for _, val := range resp.Values {
		if len(val) > NameToIndex["model"]+1 {
			key := val[NameToIndex["model"]].(string)
			if key == "" {
				continue
			}
			var compkey string
			if len(val) > NameToIndex["component"]+1 {
				compkey = val[NameToIndex["component"]].(string)
			}
			if compkey == "" {
				availableModels[key] = make([]interface{}, len(val))
				copy(availableModels[key], val)
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
	dp := newdedup()

	executeInStages(StartPipeline, csvChan, spreadsheetChan, dp, priority, cncf, official, verified, unverified)
	time.Sleep(20 * time.Second)
	
	close(spreadsheetChan)
	wg.Wait()
}

// Stages have to run sequentially. The steps within each stage can be concurrent.
// pipeline function should return only after completion
func executeInStages(pipeline func(in chan []artifacthub.AhPackage, csv chan string, spreadsheet chan struct {
	comps   []v1alpha1.ComponentDefinition
	model   string
	helmURL string
}, dp *dedup) error,
	csv chan string,
	spreadsheetChan chan struct {
		comps   []v1alpha1.ComponentDefinition
		model   string
		helmURL string
	}, dp *dedup, pkg ...[]artifacthub.AhPackage) {
	for stageno, p := range pkg {
		input := make(chan []artifacthub.AhPackage)
		go func() {
			for len(p) != 0 {
				x := 50
				if len(p) < x {
					x = len(p)
				}
				input <- p[:x]
				p = p[x:]
			}
			close(input)
		}()
		var wg sync.WaitGroup
		for i := 1; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pipeline(input, csv, spreadsheetChan, dp) //synchronous
				fmt.Println("Pipeline exited for a go routine")
			}()
		}
		wg.Wait()
		fmt.Println("[DEBUG] Completed stage", stageno)
	}
}

type dedup struct {
	m  map[string]bool
	mx sync.Mutex
}

func newdedup() *dedup {
	return &dedup{
		m: make(map[string]bool),
	}
}
func (d *dedup) set(key string) {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.m[key] = true
}
func (d *dedup) check(key string) bool {
	return d.m[key]
}
func StartPipeline(in chan []artifacthub.AhPackage, csv chan string, spreadsheet chan struct {
	comps   []v1alpha1.ComponentDefinition
	model   string
	helmURL string
}, dp *dedup) error {
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
		close(pkgsChan)
	}()
	// writer
	go func() {
		for modelcomps := range compsChan {
			err := writeComponents(modelcomps.comps)
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
	for pkgs := range pkgsChan {
		for _, ap := range pkgs {
			fmt.Printf("[DEBUG] Generating components for: %s with verified status %v\n", ap.Name, ap.VerifiedPublisher)
			comps, err := ap.GenerateComponents()
			if err != nil {
				fmt.Println(err)
				continue
			}
			var newcomps []v1alpha1.ComponentDefinition
			for _, comp := range comps {
				key := fmt.Sprintf("%sMESHERY%s", comp.Kind, comp.APIVersion)
				if !dp.check(key) {
					fmt.Println("SETTING FOR: ", key)
					newcomps = append(newcomps, comp)
					dp.set(key)
				}
			}
			compsCSV <- struct {
				comps []v1alpha1.ComponentDefinition
				model string
			}{
				comps: newcomps,
				model: ap.Name,
			}
			compsChan <- struct {
				comps []v1alpha1.ComponentDefinition
				model string
			}{
				comps: newcomps,
				model: ap.Name,
			}
			spreadsheet <- struct {
				comps   []v1alpha1.ComponentDefinition
				model   string
				helmURL string
			}{
				comps:   newcomps,
				model:   ap.Name,
				helmURL: ap.ChartUrl,
			}
		}

	}
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

func writeComponentModels(models []artifacthub.AhPackage, writer *Writer) error {
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

func writeComponents(cmps []v1alpha1.ComponentDefinition) error {
	// writer.m.Lock()
	// defer writer.m.Unlock()
	for _, comp := range cmps {
		modelPath := filepath.Join(OutputDirectoryPath, comp.Model.Name)
		if _, err := os.Stat(modelPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(modelPath, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Println("created directory ", comp.Model.Name)
		}
		componentPath := filepath.Join(modelPath, comp.Model.Version)
		if _, err := os.Stat(componentPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(componentPath, os.ModePerm)
			if err != nil {
				return err
			}
			fmt.Println("created versioned directory ", comp.Model.Version)
		}
		relationshipsPath := filepath.Join(modelPath, "relationships")
		policiesPath := filepath.Join(modelPath, "policies")
		if _, err := os.Stat(relationshipsPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(relationshipsPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		if _, err := os.Stat(policiesPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(policiesPath, os.ModePerm)
			if err != nil {
				return err
			}
		}
		f, err := os.Create(filepath.Join(componentPath, comp.Kind+".json"))
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
	comps   []v1alpha1.ComponentDefinition
	model   string
	helmURL string
}, am map[string][]interface{}, acpm map[string]map[string]bool) {
	start := time.Now()
	rangeString := sheetName + "!A4:AB4"
	// Get the value of the specified cell.
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
	if err != nil {
		fmt.Println("Unable to retrieve data from sheet: ", err)
		return
	}
	batchSize := 100
	values := make([][]interface{}, 0)
	for entry := range spreadsheet {
		if len(entry.comps) == 0 {
			continue
		}
		for _, comp := range entry.comps {
			if acpm[entry.model][comp.Kind] {
				fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model, comp.Kind)
				continue
			}
			var newValues []interface{}
			if am[entry.model] != nil {
				newValues = make([]interface{}, len(am[entry.model]))
				copy(newValues, am[entry.model])
			} else {
				newValues = make([]interface{}, len(resp.Values[0]))
				copy(newValues, resp.Values[0])
				newValues[NameToIndex["modelDisplayName"]] = entry.model
				newValues[NameToIndex["model"]] = entry.model
			}
			newValues[NameToIndex["component"]] = comp.Kind
			if comp.Schema != "" {
				newValues[NameToIndex["hasSchema?"]] = true
			} else {
				newValues[NameToIndex["hasSchema?"]] = false
			}
			newValues[NameToIndex["link"]] = entry.helmURL
			values = append(values, newValues)
			if acpm[entry.model] == nil {
				acpm[entry.model] = make(map[string]bool)
			}
			acpm[entry.model][comp.Kind] = true
			batchSize--
			fmt.Println("Batch size: ", batchSize)
			if batchSize <= 0 {
				// row := &sheets.ValueRange{
				// 	Values: values,
				// }
				// response2, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
				// values = make([][]interface{}, 0)
				// batchSize = 100
				// if err != nil || response2.HTTPStatusCode != 200 {
				// 	fmt.Println(err)
				// 	continue
				// }
			}
		}
		if am[entry.model] != nil {
			fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model)
			continue
		}
		newValues := make([]interface{}, len(resp.Values[0]))
		copy(newValues, resp.Values[0])
		newValues[NameToIndex["modelDisplayName"]] = entry.model
		newValues[NameToIndex["model"]] = entry.model
		newValues[NameToIndex["CRDs"]] = len(entry.comps)
		newValues[NameToIndex["link"]] = entry.helmURL
		values = append(values, newValues)
		copy(am[entry.model], newValues)
		batchSize--
		fmt.Println("Batch size: ", batchSize)
		if batchSize <= 0 {
			// row := &sheets.ValueRange{
			// 	Values: values,
			// }
			// response2, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
			// values = make([][]interface{}, 0)
			// batchSize = 100
			// if err != nil || response2.HTTPStatusCode != 200 {
			// 	fmt.Println(err)
			// 	continue
			// }
		}
	}
	if len(values) != 0 {
		// row := &sheets.ValueRange{
		// 	Values: values,
		// }
		// response2, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
		// if err != nil || response2.HTTPStatusCode != 200 {
		// 	fmt.Println(err)
		// }
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Time taken by spreadsheet updater in minutes (including the time it required to generate components): %f", elapsed.Minutes())
}

func NewSheetSRV() *sheets.Service {
	ctx := context.Background()
	byt, _ := os.ReadFile("/Users/shabana/Downloads/spreadsheet.json")
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
