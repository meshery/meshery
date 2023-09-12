/*
Meshery Component Updater
Uses a spreadsheet of centralized information about MeshModel components and their metadata like color, icon, and so on. Script is used to update components metada (svgs, icons etc) for Meshery, Websites (Layer5.io, Meshery.io), and Remote Provider.

Usage: (order of flags matters)

    ./main [path-to-spreadsheet] [--system] [<system-name>] [relative path to docs in layer5 website] [relative path to docs in meshery website] [--only-published]

Examples:

	1. ./main https://docs.google.com/spreadsheets/d/e/2PACX-1vSgOXuiqbhUgtC9oNbJlz9PYpOEaFVoGNUFMIk4NZciFfQv1ewZg8ahdrWHKI79GkKK9TbmnZx8CqIe/pub\?gid\=0\&single\=true\&output\=csv --system docs layer5/src/collections/integrations meshery.io/integrations docs/ --published-only
	2. ./main https://docs.google.com/spreadsheets/d/e/2PACX-1vSgOXuiqbhUgtC9oNbJlz9PYpOEaFVoGNUFMIk4NZciFfQv1ewZg8ahdrWHKI79GkKK9TbmnZx8CqIe/pub\?gid\=0\&single\=true\&output\=csv --system remote-provider <remote-provider>/meshmodels/models <remote-provider>/ui/public/img/meshmodels
	3. ./main https://docs.google.com/spreadsheets/d/e/2PACX-1vSgOXuiqbhUgtC9oNbJlz9PYpOEaFVoGNUFMIk4NZciFfQv1ewZg8ahdrWHKI79GkKK9TbmnZx8CqIe/pub\?gid\=0\&single\=true\&output\=csv --system meshery ../../server/meshmodel

The flags are:

  --system
        defined type of system to update. Can be one of "meshery", "docs", or "remote-provider".

	--only-published
        Only handle components that have a value of "true" under the "Published?" column in spreadsheet.
*/

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/layer5io/component_scraper/pkg"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/manifests"
)

var (
	ColumnNamesToExtract        = []string{"modelDisplayName", "model", "category", "subCategory", "shape", "primaryColor", "secondaryColor", "logoURL", "svgColor", "svgWhite", "Publish?", "CRDs", "component", "svgComplete", "genealogy", "styleOverrides"}
	ColumnNamesToExtractForDocs = []string{"modelDisplayName", "Page Subtitle", "Docs URL", "category", "subCategory", "Feature 1", "Feature 2", "Feature 3", "howItWorks", "howItWorksDetails", "Publish?", "About Project", "Standard Blurb", "svgColor", "svgWhite", "Full Page", "model"}
	PrimaryColumnName           = "model"
	OutputPath                  = ""
	ExcludeDirs = []string{"relationships", "policies"}
)

var System string

const (
	SVG_WIDTH  = 20
	SVG_HEIGHT = 20
)

func main() {
	url := os.Args[1]
	if url == "" {
		log.Fatal("provide a valid spreadsheet URL")
		return
	}

	if len(os.Args) <= 2 {
		log.Fatal("system flag is missing")
		return
	}

	if os.Args[2] != "--system" {
		log.Fatal("invalid system flag provided")
		return
	}

	if len(os.Args) <= 3 {
		log.Fatal("system name is missing")
		return
	}

	if os.Args[3] == "" {
		log.Fatal("system name missing")
		return
	}
	System = os.Args[3]

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
	output, err := pkg.GetEntries(csvReader, ColumnNamesToExtractForDocs)
	if err != nil {
		log.Fatal(err)
		return
	}
	file.Close()
	os.Remove(file.Name())

	switch System {
	case pkg.Docs.String():
		docsUpdater(output)
	case pkg.Meshery.String():
		mesheryUpdater(output)
	case pkg.RemoteProvider.String():
		remoteProviderUpdater(output)
	default:
		log.Fatal("invalid system name")
		return
	}

}

// returns the index of column. Returns -1 if doesn't exist
func contains(key string, col []string) int {
	for i, n := range col {
		if n == key {
			return i
		}
	}
	return -1
}

// For Docs:: entries with empty Component field are preferred as they are considered general
// In other words, the absence of a component name indicates that a given row is a Model-level entry.
// And that for docs/websites updates, the values found in this row should be used to represent the
// integration overall (whether there is 1 or many 10s of components contained in the package / in the integration).
func cleanupDuplicatesAndPreferEmptyComponentField(out []map[string]string, groupBykey string) (out2 []map[string]string) {
	keyToEntry := make(map[string]map[string]string)
	for _, o := range out {
		gkey := o[groupBykey]
		//If the row with given gkey is encountered for the first time, or the given row already exists but with a non-empty component field then use the new entry.
		//This logic will prioritize empty component fields to not be overriden
		if keyToEntry[gkey] == nil || keyToEntry[gkey]["component"] != "" {
			keyToEntry[gkey] = o
		}

	}
	for _, entry := range keyToEntry {
		out2 = append(out2, entry)
	}
	return out2
}

func docsUpdater(output []map[string]string) {
	if len(os.Args) < 8 {
		log.Fatal("docsUpdater: invalid number of arguments; missing website and docs path")
		return
	}
	pathToIntegrationsLayer5 := os.Args[4]
	pathToIntegrationsMeshery := os.Args[5]
	pathToIntegrationsMesheryDocs := os.Args[6]	
	updateOnlyPublished := true
	if len(os.Args) > 7 {
		if os.Args[7] == "--published-only" {
			updateOnlyPublished = true
		}
	}
	output = cleanupDuplicatesAndPreferEmptyComponentField(output, "model")
	mesheryDocsJSON := "const data = ["
	for _, out := range output {
		var t pkg.TemplateAttributes
		publishValue, err := strconv.ParseBool(out["Publish?"])
		if err != nil {
			publishValue = false
		}
		if updateOnlyPublished && !publishValue {
			continue
		}
		for key, val := range out {
			switch key {
			case "modelDisplayName":
				t.Title = val
			case "model":
				t.ModelName = val
			case "Page Subtitle":
				t.Subtitle = val
			case "Docs URL":
				t.DocURL = val
			case "category":
				t.Category = val
			case "subCategory":
				t.Subcategory = val
			case "howItWorks":
				t.HowItWorks = val
			case "hotItWorksDetails":
				t.HowItWorksDetails = val
			case "Publish?":
				t.Published = val
			case "About Project":
				t.AboutProject = val
			case "Standard Blurb":
				t.StandardBlurb = val
			case "Full Page":
				t.FullPage = val
			case "svgColor":
				t.ColorSVG = val
			case "svgWhite":
				t.WhiteSVG = val
			}
		}
		t.FeatureList = "[" + strings.Join([]string{out["Feature 1"], out["Feature 2"], out["Feature 3"]}, ",") + "]"
		t.WorkingSlides = `[
				../_images/meshmap-visualizer.png,
				../_images/meshmap-designer.png]`

		md := t.CreateMarkDown()
		jsonItem := t.CreateJSONItem()
		mesheryDocsJSON += jsonItem + ","
		modelName := strings.TrimSpace(out["model"])
		pathToIntegrationsLayer5, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsLayer5, modelName))
		pathToIntegrationsMeshery, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsMeshery))
		pathToIntegrationsMesheryDocs, _ := filepath.Abs(filepath.Join("../../", pathToIntegrationsMesheryDocs, "assets/img/meshmodel/", modelName))	
		err = os.MkdirAll(pathToIntegrationsLayer5, 0777)
		if err != nil {
			panic(err)
		}
		_ = pkg.WriteToFile(filepath.Join(pathToIntegrationsLayer5, "index.mdx"), md)
		svgcolor := out["svgColor"]
		svgwhite := out["svgWhite"]

		// Write SVGs to Layer5 docs
		err = os.MkdirAll(filepath.Join(pathToIntegrationsLayer5, "icon", "color"), 0777)
		if err != nil {
			panic(err)
		}

		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsLayer5, "icon", "color", modelName+"-color.svg"), svgcolor) //CHANGE PATH
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(filepath.Join(pathToIntegrationsLayer5, "icon", "white"), 0777)
		if err != nil {
			panic(err)
		}
		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsLayer5, "icon", "white", modelName+"-white.svg"), svgwhite) //CHANGE PATH
		if err != nil {
			panic(err)
		}

		// Write SVGs to Meshery website
		err = os.MkdirAll(filepath.Join(pathToIntegrationsMeshery, "../", "images"), 0777)
		if err != nil {
			panic(err)
		}

		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsMeshery, "../", "images", modelName+"-color.svg"), svgcolor) //CHANGE PATH
		if err != nil {
			panic(err)
		}

		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsMeshery, "../", "images", modelName+"-white.svg"), svgwhite) //CHANGE PATH
		if err != nil {
			panic(err)
		}

		// Write SVGs to Meshery docs
		err = os.MkdirAll(filepath.Join(pathToIntegrationsMesheryDocs), 0777)
		if err != nil {
			panic(err)
		}

		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsMesheryDocs, modelName+"-color.svg"), svgcolor) //CHANGE PATH
		if err != nil {
			panic(err)
		}

		err = pkg.WriteSVG(filepath.Join(pathToIntegrationsMesheryDocs, modelName+"-white.svg"), svgwhite) //CHANGE PATH
		if err != nil {
			panic(err)
		}
	}

	mesheryDocsJSON = strings.TrimSuffix(mesheryDocsJSON, ",")
	mesheryDocsJSON += "]; export default data"
	if err := pkg.WriteToFile(filepath.Join("../../../", pathToIntegrationsMeshery, "data.js"), mesheryDocsJSON); err != nil {
		log.Fatal(err)
	}

	if err := pkg.WriteToFile(filepath.Join("../../", pathToIntegrationsMesheryDocs, "_data/integrations/", "data.js"), mesheryDocsJSON); err != nil {
		log.Fatal(err)
	}
}

func mesheryUpdater(output []map[string]string) {
	if len(os.Args) < 5 {
		log.Fatal("mesheryUpdater: invalid number of arguments; missing meshmodels path in meshery server")
		return
	}
	OutputPath = os.Args[4]
	if OutputPath == "" {
		OutputPath = "../../server/meshmodel" // default path for meshery server
	}
	publishedModels := make(map[string]bool)
	countWithoutCrds := 0
	_ = pkg.PopulateEntries(OutputPath, output, PrimaryColumnName, func(dirpath string, changeFields map[string]string) error {
		if changeFields["CRDs"] == "" {
			countWithoutCrds++
		}
		if changeFields["Publish?"] == "TRUE" { //For a component level field
			publishedModels[changeFields[PrimaryColumnName]] = true
		}
		return nil
	})
	err := pkg.PopulateEntries(OutputPath, output, PrimaryColumnName, func(dirpath string, changeFields map[string]string) error {
		entries, err := os.ReadDir(dirpath)
		if err != nil {
			return err
		}
		for _, versionentry := range entries {
			if versionentry.IsDir() {
				if contains(versionentry.Name(), ExcludeDirs) != -1 {
					continue
				}
				entries, err := os.ReadDir(filepath.Join(dirpath, versionentry.Name()))
				if err != nil {
					return err
				}
				for _, entry := range entries {
					if !publishedModels[changeFields[PrimaryColumnName]] && changeFields["component"] == "" { //Ignore Publish flag for component level rows
						fmt.Println("(Publish? is not set to TRUE) will not update", changeFields["modelDisplayName"])
						continue
					}
					name := strings.TrimSuffix(strings.TrimSpace(entry.Name()), ".json")
					if changeFields["component"] != "" && changeFields["component"] != name { //This is a component specific entry and only fill this when the filename matches component name
						continue
					}

					path, err := filepath.Abs(filepath.Join(dirpath, versionentry.Name(), entry.Name()))
					if err != nil {
						fmt.Printf("[Error] for %s: %s\n", name, err.Error())
						continue
					}
					byt, err := os.ReadFile(path)
					if err != nil {
						fmt.Printf("[Error] for %s: %s\n", name, err.Error())
						continue
					}
					var component v1alpha1.ComponentDefinition
					err = json.Unmarshal(byt, &component)
					if err != nil {
						fmt.Printf("[Error] for %s: %s\n", name, err.Error())
						continue
					}
					component.DisplayName = manifests.FormatToReadableString(component.Kind)
					if component.Metadata == nil {
						component.Metadata = make(map[string]interface{})
					}
					for key, value := range changeFields {
						if key == "category" {
							component.Model.Category = v1alpha1.Category{
								Name: value,
							}
						} else if key == "svgColor" || key == "svgWhite" {
							svg, err := pkg.UpdateSVGString(value, SVG_WIDTH, SVG_HEIGHT)
							if err != nil {
								fmt.Println("err for: ", component.Kind, err.Error())
							}
							if changeFields["component"] == "" { //If it is a model level entry then update model svgs
								if component.Model.Metadata == nil {
									component.Model.Metadata = make(map[string]interface{})
								}
								component.Model.Metadata[key] = svg
							}
							if changeFields["component"] != "" || component.Metadata[key] == nil { // If it is a component level SVG or component already doesn't have an SVG. Use this svg at component level.
								component.Metadata[key] = svg
							}
						} else if contains(key, ColumnNamesToExtract) != -1 {
							component.Metadata[key] = value
						}
					}
					if i := contains("modelDisplayName", ColumnNamesToExtract); i != -1 {
						component.Model.DisplayName = changeFields[ColumnNamesToExtract[i]]
					}
					isAnnotation, _ := component.Metadata["isAnnotation"].(bool)
					if isAnnotation {
						component.Metadata["isNamespaced"] = false
					}
					//Either component is set to published or the parent model is set to published
					if component.Metadata["Publish?"] == "TRUE" || publishedModels[component.Model.Name] { //Publish? is an invalid field for putting inside kubernetes annotations
						component.Metadata["published"] = true
					} else {
						component.Metadata["published"] = false
					}
					fmt.Println("updating for ", changeFields["modelDisplayName"], "--", component.Kind, "-- published=", component.Metadata["published"])
					delete(component.Metadata, "Publish?")
					delete(component.Metadata, "CRDs")
					delete(component.Metadata, "component")
					modelDisplayName := component.Metadata["modelDisplayName"].(string)
					component.Model.DisplayName = modelDisplayName
					byt, err = json.Marshal(component)
					if err != nil {
						fmt.Printf("[Error] for %s: %s\n", name, err.Error())
						continue
					}
					err = os.WriteFile(filepath.Join(dirpath, versionentry.Name(), entry.Name()), byt, 0777)
					if err != nil {
						fmt.Printf("[Error] for %s: %s\n", name, err.Error())
						continue
					}
				}
			}

		}
		return nil
	})
	fmt.Println("Total models without CRDs in spreadsheet are: ", countWithoutCrds)
	if err != nil {
		log.Fatal(err)
	}
}

func remoteProviderUpdater(output []map[string]string) {
	if len(os.Args) < 6 {
		log.Fatal("remoteProvider updater: invalid number of arguments; missing meshmodels path in remote provider")
		return
	}

	output = cleanupDuplicatesAndPreferEmptyComponentField(output, "model")
	pathForModals := os.Args[4]
	pathForIcons := os.Args[5]
	_path := ""
	for _, f := range strings.Split(pathForIcons, "/")[1:] {
		_path = filepath.Join(_path, f)
	}
	for _, out := range output {
		var m v1alpha1.Model
		var svgColor, svgWhite string
		publishValue, err := strconv.ParseBool(out["Publish?"])
		if err != nil {
			publishValue = false
		}
		if !publishValue {
			continue
		}
		modelName := strings.TrimSpace(out["model"])
		if m.Metadata == nil {
			m.Metadata = make(map[string]interface{})
		}
		for key, val := range out {
			switch key {
			case "modelDisplayName":
				m.DisplayName = val
			case "model":
				m.Name = val
			case "category":
				m.Category = v1alpha1.Category{
					Name: val,
				}
			// case "subCategory":
			// 	m.SubCategory = val
			case "svgColor":
				svgColor, err = pkg.UpdateSVGString(val, SVG_WIDTH, SVG_HEIGHT)
				if err != nil {
					fmt.Println("err for: ", modelName, err.Error())
				}
				m.Metadata["svgColor"] = fmt.Sprintf("%s/%s/color/%s-color.svg", _path, modelName, modelName)
			case "svgWhite":
				svgWhite, err = pkg.UpdateSVGString(val, SVG_WIDTH, SVG_HEIGHT)
				if err != nil {
					fmt.Println("err for: ", modelName, err.Error())
				}
				m.Metadata["svgWhite"] = fmt.Sprintf("%s/%s/white/%s-white.svg", _path, modelName, modelName)
			}
		}
		pathForModals, _ := filepath.Abs(filepath.Join("../../../", pathForModals, modelName))
		err = os.MkdirAll(pathForModals, 0777)
		if err != nil {
			fmt.Println("Error creating directory: ", err.Error())
		}
		pathForIcons, _ := filepath.Abs(filepath.Join("../../../", pathForIcons, modelName))
		err = os.MkdirAll(pathForIcons, 0777)
		if err != nil {
			fmt.Println("Error creating directory: ", err.Error())
		}
		byt, err := json.Marshal(m)
		if err != nil {
			fmt.Println("Error marshalling model: ", err.Error())
		}
		err = os.WriteFile(filepath.Join(pathForModals, "model.json"), byt, 0777)
		if err != nil {
			fmt.Println("Error writing model: ", err.Error())
			continue
		}
		err = os.MkdirAll(filepath.Join(pathForIcons, "color"), 0777)
		if err != nil {
			panic(err)
		}
		err = pkg.WriteSVG(filepath.Join(pathForIcons, "color", modelName+"-color.svg"), svgColor) //CHANGE PATH
		if err != nil {
			panic(err)
		}
		err = os.MkdirAll(filepath.Join(pathForIcons, "white"), 0777)
		if err != nil {
			panic(err)
		}
		err = pkg.WriteSVG(filepath.Join(pathForIcons, "white", modelName+"-white.svg"), svgWhite) //CHANGE PATH
		if err != nil {
			panic(err)
		}
	}
}
