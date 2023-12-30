package main

import (
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

	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/artifacthub"
	"gopkg.in/yaml.v3"
)

const dumpFile = "./dump.csv"
const COLUMNRANGE = "!A:AH3" //Update this on addition of new columns
const APPENDRANGE = "!A4:V4"
const sheetID = 0

var spreadsheetID = "1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw"

type componentWrapper struct {
	comps   []v1alpha1.ComponentDefinition
	model   string
	helmURL string
}

type pipelineFunc func(in chan artifacthub.AhPackage, csv chan componentWrapper, spreadsheet chan componentWrapper, dp *dedup) error

type ComponentStruct struct {
	PackageName string                         `yaml:"name"`
	Components  []v1alpha1.ComponentDefinition `yaml:"components"`
}

var NameToIndex = map[string]int{ //Update this on addition of new columns
	"modelDisplayName":   0,
	"model":              1,
	"category":           2,
	"subCategory":        3,
	"CRDs":               4,
	"link":               5,
	"hasSchema?":         6,
	"component":          7,
	"shape":              8,
	"primaryColor":       9,
	"secondaryColor":     10,
	"styleOverrides":     11,
	"styles":             12,
	"shapePolygonPoints": 13,
	"defaultData":        14,
	"logoURL":            15,
	"svgColor":           16,
	"svgWhite":           17,
	"svgComplete":        18,
	"genealogy":          19,
	"isAnnotation":       20,
	"PublishToRegistry":  21,
	"isModelAnnotation":  22,
	"About Project":      24,
	"Page Subtitle":      25,
	"Docs URL":           26,
	"Standard Blurb":     27,
	"Feature 1":          28,
	"Feature 2":          29,
	"Feature 3":          30,
	"howItWorks":         31,
	"howItWorksDetails":  32,
	"Screenshots":        33,
	"Full Page":          34,
}
var (
	AhSearchEndpoint = artifacthub.AhHelmExporterEndpoint

	OutputDirectoryPath     = "../../server/meshmodel"
	ComponentModelsFileName = path.Join(OutputDirectoryPath, "component_models.yaml")
)

var priorityRepos = map[string]bool{"prometheus-community": true, "grafana": true} //Append ahrepos here whose components should be respected and should be used when encountered duplicates

var log, _ = logger.New("mesheryctl", logger.Options{
	Format:     logger.SyslogLogFormat,
	DebugLevel: false,
})

// returns pkgs with sorted pkgs at the front
func sortOnVerified(pkgs []artifacthub.AhPackage) (verified []artifacthub.AhPackage, official []artifacthub.AhPackage, cncf []artifacthub.AhPackage, priority []artifacthub.AhPackage, unverified []artifacthub.AhPackage) {
	for _, pkg := range pkgs {
		if priorityRepos[pkg.Repository] {
			priority = append(priority, pkg)
		} else if pkg.CNCF {
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
			log.Error(err)
			return
		}
	}

	modelsFd, err := os.OpenFile(ComponentModelsFileName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
		return
	}

	defer modelsFd.Close()
	// move to a new function: getHelmPackages
	pkgs := make([]artifacthub.AhPackage, 0)
	content, err := io.ReadAll(modelsFd)
	if err != nil {
		log.Error(err)
		return
	}
	err = yaml.Unmarshal(content, &pkgs)
	if err != nil {
		log.Error(err)
	}

	if len(pkgs) == 0 {
		pkgs, err = artifacthub.GetAllAhHelmPackages()
		if err != nil {
			log.Error(err)
			return
		}
		err = writeComponentModels(pkgs, modelsFd)
		if err != nil {
			log.Error(err)
			return
		}
	}
	verified, official, cncf, priority, unverified := sortOnVerified(pkgs)
	compsCSV := make(chan componentWrapper, 50)
	go dumpCsv(compsCSV)

	srv := NewSheetSRV()
	// Convert sheet ID to sheet name.
	response1, err := srv.Spreadsheets.Get(spreadsheetID).Fields("sheets(properties(sheetId,title))").Do()
	if err != nil || response1.HTTPStatusCode != 200 {
		log.Error(err)
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
	spreadsheetChan := make(chan componentWrapper, 100)
	// Set the range of cells to retrieve.
	rangeString := sheetName + COLUMNRANGE

	// Get the value of the specified cell.
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
	if err != nil {
		log.Error(ErrorFailedRetreiving(err))
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

	updater := spreadsheetUpdater{
		spreadsheetID:               spreadsheetID,
		sheetName:                   sheetName,
		spreadsheetChan:             spreadsheetChan,
		availableModels:             availableModels,
		availableComponentsPerModel: availableComponentsPerModel,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		updater.update(srv)
	}()
	dp := newdedup()

	executeInStages(StartPipeline, compsCSV, spreadsheetChan, dp, priority, cncf, official, verified, unverified)
	time.Sleep(20 * time.Second)

	close(spreadsheetChan)
	wg.Wait()
}

func dumpCsv(compsCSV chan componentWrapper) {
	f, err := os.Create(dumpFile)
	if err != nil {
		log.Error(err)
		return
	}
	f.Write([]byte("model,component_count,components\n"))
	for comps := range compsCSV {
		count := len(comps.comps)
		names := "\""
		for _, cmp := range comps.comps {
			names += fmt.Sprintf("%s,", cmp.Kind)
		}
		names = strings.TrimSuffix(names, ",")
		names += "\""
		if count > 0 {
			f.Write([]byte(fmt.Sprintf("%s,%d,%s\n", comps.model, count, names)))
		}
	}
}

// Stages have to run sequentially. The steps within each stage can be concurrent.
// pipeline function should return only after completion
func executeInStages(pipeline pipelineFunc, compsCSV chan componentWrapper,
	spreadsheetChan chan componentWrapper, dp *dedup, pkgs ...[]artifacthub.AhPackage) {
	for stageno, pkg := range pkgs {
		input := make(chan artifacthub.AhPackage)
		var wg sync.WaitGroup
		for i := 1; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				pipeline(input, compsCSV, spreadsheetChan, dp) //synchronous
				log.Debug("Pipeline exited for a go routine")
			}()
		}
		for _, p := range pkg {
			input <- p
		}
		wg.Wait()
		log.Debug("Completed stage", stageno)
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

func StartPipeline(pkgChan chan artifacthub.AhPackage, compsCSV chan componentWrapper, spreadsheet chan componentWrapper, dp *dedup) error {
	for pkg := range pkgChan {
		if err := pkg.UpdatePackageData(); err != nil {
			log.Error(err)
			continue
		}
		log.Debug(fmt.Sprintf("Generating components for: %s with verified status %v", pkg.Name, pkg.VerifiedPublisher))
		comps, err := pkg.GenerateComponents()
		if err != nil {
			log.Error(err)
			continue
		}
		var newcomps []v1alpha1.ComponentDefinition
		for _, comp := range comps {
			key := fmt.Sprintf("%sMESHERY%s", comp.Kind, comp.APIVersion)
			if !dp.check(key) {
				log.Debug("SETTING FOR: ", key)
				newcomps = append(newcomps, comp)
				dp.set(key)
			}
		}
		if err := writeComponents(comps); err != nil {
			log.Error(err)
		}
		compsCSV <- componentWrapper{
			comps: newcomps,
			model: pkg.Name,
		}
		spreadsheet <- componentWrapper{
			comps:   newcomps,
			model:   pkg.Name,
			helmURL: pkg.ChartUrl,
		}
	}
	return nil
}

func writeComponentModels(models []artifacthub.AhPackage, file *os.File) error {
	val, err := yaml.Marshal(models)
	if err != nil {
		return err
	}
	if _, err := file.Write(val); err != nil {
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
			log.Debug("created directory ", comp.Model.Name)
		}
		componentPath := filepath.Join(modelPath, comp.Model.Version)
		if _, err := os.Stat(componentPath); errors.Is(err, os.ErrNotExist) {
			err := os.Mkdir(componentPath, os.ModePerm)
			if err != nil {
				return err
			}
			log.Debug("created versioned directory ", comp.Model.Version)
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
