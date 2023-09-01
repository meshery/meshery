package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

const (
	dumpFile         = "./dump.csv"
	COLUMNRANGE      = "!A:AF3" // Update this on addition of new columns
	OutputDirPath    = "../../server/meshmodel"
	CompModelsFile   = "component_models.yaml"
	CompModelFileURL = ""
	sheetID          = 0
	batchSize        = 50
	// spreadsheetID  = "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw"
)

var (
	NameToIndex = map[string]int{
		// Update this on addition of new columns
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

	priorityRepos = map[string]bool{
		"prometheus-community": true,
		"grafana":              true,
		// Append ahrepos here whose components should be respected and should be used when encountering duplicates
	}

	AhSearchEndpoint = artifacthub.AhHelmExporterEndpoint
	updateFlag       = false
)

type ComponentStruct struct {
	PackageName string                         `yaml:"name"`
	Components  []v1alpha1.ComponentDefinition `yaml:"components"`
}

type Writer struct {
	file *os.File
	m    sync.Mutex
}

type Dedup struct {
	m  map[string]bool
	mx sync.Mutex
}

func newDedup() *Dedup {
	return &Dedup{
		m: make(map[string]bool),
	}
}

func (d *Dedup) set(key string) {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.m[key] = true
}

func (d *Dedup) check(key string) bool {
	d.mx.Lock()
	defer d.mx.Unlock()
	return d.m[key]
}

func main() {
	spreadsheetID := ""
	if len(os.Args) > 1 {
		spreadsheetID = os.Args[1]
	}

	if len(os.Args) > 2 {
		if os.Args[2] == "--update" {
			updateFlag = true
		}
	}

	if _, err := os.Stat(OutputDirPath); err != nil {
		err := os.Mkdir(OutputDirPath, 0744)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	modelsFd, err := os.OpenFile(filepath.Join(OutputDirPath, CompModelsFile), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer modelsFd.Close()

	pkgs, err := loadPackages(modelsFd)
	if err != nil {
		fmt.Println(err)
		return
	}

	if len(pkgs) == 0 {
		pkgs, err = artifacthub.GetAllAhHelmPackages()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = writeComponentModels(modelsFd, pkgs)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	verified, official, cncf, priority, unverified := sortOnVerified(pkgs)

	f, err := os.Create(dumpFile)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer f.Close()
	f.Write([]byte("model,component_count,components\n"))

	srv := NewSheetSRV()

	sheetName, err := getSheetName(srv, spreadsheetID, sheetID)
	if err != nil {
		fmt.Println(err)
		return
	}

	availableModels, availableComponentsPerModel, err := getAvailableModelsAndComponents(srv, spreadsheetID, sheetName)
	if err != nil {
		fmt.Println("Error getting available models & components: %s\n", err)
		return
	}

	if updateFlag {
		err := updateSpreadsheet(srv, spreadsheetID, sheetName, availableModels, availableComponentsPerModel)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	dp := newDedup()

	executeInStages(StartPipeline, f, dp, priority, cncf, official, verified, unverified)
	time.Sleep(20 * time.Second)
}

func loadPackages(file string) ([]artifacthub.AhPackage, error) {
	pkgs := make([]artifacthub.AhPackage, 0)
	var content []byte
	var err error

	if strings.HasPrefix(file, "http://") || strings.HasPrefix(file, "https://") {
		// Fetch the content from a remote URL
		response, err := http.Get(file)
		if err != nil {
			return pkgs, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return pkgs, fmt.Errorf("Failed to fetch content from URL: %s", response.Status)
		}

		content, err = io.ReadAll(response.Body)
		if err != nil {
			return pkgs, err
		}
	} else {
		// Load the content from a local file
		content, err = os.ReadFile(file)
		if err != nil {
			return pkgs, err
		}
	}

	err = yaml.Unmarshal(content, &pkgs)
	if err != nil {
		return pkgs, err
	}
	return pkgs, nil
}

func writeComponentModels(file *os.File, models []artifacthub.AhPackage) error {
	writer := Writer{file: file}
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

func writeComponents(comps []v1alpha1.ComponentDefinition) error {
	for _, comp := range comps {
		modelPath := filepath.Join(OutputDirPath, comp.Model.Name)
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

		// Write the component data to a JSON file
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

func executeInStages(pipeline func(csv *os.File, dp *Dedup, pkgs ...[]artifacthub.AhPackage) error, csv *os.File, dp *Dedup, pkg ...[]artifacthub.AhPackage) {
	for stageno, p := range pkg {
		var wg sync.WaitGroup
		input := make(chan []artifacthub.AhPackage, 10) // Buffer the channel for better concurrency
		// Create goroutines to read from the channel and pass to the pipeline function
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for pkgs := range input {
					pipeline(csv, dp, pkgs)
				}
			}()
		}
		// Split the packages and send them to the channel
		for len(p) > 0 {
			batchSize := 50
			if len(p) < batchSize {
				batchSize = len(p)
			}
			input <- p[:batchSize]
			p = p[batchSize:]
		}
		close(input)
		wg.Wait()
		fmt.Println("[DEBUG] Completed stage", stageno)
	}
}

func getSheetName(srv *sheets.Service, spreadsheetID string, sheetID int) (string, error) {
	var sheetName string
	response1, err := srv.Spreadsheets.Get(spreadsheetID).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		fmt.Println(err)
		return sheetName, err
	}
	for _, v := range response1.Sheets {
		prop := v.Properties
		if prop.SheetId == int64(sheetID) {
			sheetName = prop.Title
			break
		}
	}
	return sheetName, nil
}

func getAvailableModelsAndComponents(srv *sheets.Service, spreadsheetID string, sheetName string) (map[string][]interface{}, map[string]map[string]bool, error) {
	availableModels := make(map[string][]interface{})
	availableComponentsPerModel := make(map[string]map[string]bool)

	rangeString := sheetName + "!A4:AB4"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
	if err != nil {
		return nil, nil, err
	}

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

	return availableModels, availableComponentsPerModel, nil
}

func updateSpreadsheet(srv *sheets.Service, spreadsheetID, sheetName string, availableModels map[string][]interface{}, availableComponentsPerModel map[string]map[string]bool) error {
	batchSize := 100
	values := make([][]interface{}, 0)

	for model, modelValues := range availableModels {
		if componentsMap, ok := availableComponentsPerModel[model]; ok {
			if len(modelValues) > NameToIndex["model"] {
				modelValues[NameToIndex["CRDs"]] = len(componentsMap)
			}
		}

		values = append(values, modelValues)

		batchSize--
		if batchSize <= 0 {
			row := &sheets.ValueRange{
				Values: values,
			}
			_, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
			if err != nil {
				return fmt.Errorf("Error updating spreadsheet: %v", err)
			}
			values = make([][]interface{}, 0)
			batchSize = 100
		}
	}

	if len(values) > 0 {
		row := &sheets.ValueRange{
			Values: values,
		}
		_, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
		if err != nil {
			return fmt.Errorf("Error updating spreadsheet: %v", err)
		}
	}

	return nil
}

func StartPipeline(csv *os.File, dp *Dedup, pkgs ...[]artifacthub.AhPackage) error {
	for _, pkgBatch := range pkgs {
		ahPkgs := make([]artifacthub.AhPackage, 0)

		for _, ap := range pkgBatch {
			fmt.Printf("[DEBUG] Updating package data for: %s\n", ap.Name)
			err := ap.UpdatePackageData()
			if err != nil {
				fmt.Println(err)
				continue
			}
			ahPkgs = append(ahPkgs, ap)
		}

		err := writeComponentModels(csv, ahPkgs)
		if err != nil {
			fmt.Println(err)
		}

		for _, ap := range ahPkgs {
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

			err = writeComponents(newcomps)
			if err != nil {
				fmt.Println(err)
			}

			count := len(newcomps)
			if count > 0 {
				model := ap.Name
				fmt.Println(fmt.Sprintf("[DEBUG] Adding to CSV: %s", model))
				names := "\""
				for _, cmp := range newcomps {
					names += fmt.Sprintf("%s,", cmp.Kind)
				}
				names = strings.TrimSuffix(names, ",")
				names += "\""
				csvEntry := fmt.Sprintf("%s,%d,%s\n", model, count, names)
				_, err := csv.WriteString(csvEntry)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
	return nil
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
