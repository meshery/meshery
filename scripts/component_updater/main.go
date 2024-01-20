/*
Meshery Component Updater
Uses a spreadsheet of centralized information about MeshModel components and their metadata like color, icon, and so on. Script is used to update components metada (svgs, icons etc) for Meshery, Websites (Layer5.io, Meshery.io), and Remote Provider.

Secret - this script expects that an environment variable `CRED` is available
and it contains a token for Google Sheets API interactions.

Example:
	export CRED='{
		"type": "service_account",
		"project_id": "",
		"private_key_id": "",
		"private_key": "-----BEGIN PRIVATE KEY-----\nn-----END PRIVATE KEY-----\n",
		"client_email": "",
		"client_id": "",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url": "",
		"universe_domain": "googleapis.com"
		}'

Usage: (order of flags matters)

    ./main [path-to-spreadsheet] [--system] [<system-name>] [relative path to docs in layer5 website] [relative path to docs in meshery website] [--only-published]

Examples:

	1. CRED=$CRED go run main.go 1DZHnzxYWOlJ69Oguz4LkRVTFM79kC2tuvdwizOJmeMw --system docs layer5/src/collections/integrations meshery.io/integrations docs

	following are to be migrated from older logic
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
	ColumnNamesToExtract        = []string{"modelDisplayName", "model", "category", "subCategory", "shape", "primaryColor", "secondaryColor", "logoURL", "svgColor", "svgWhite", "isAnnotation", "isModelAnnotation", "PublishToRegistry", "CRDs", "component", "svgComplete", "capabilities", "genealogy", "styleOverrides", "capabilities"}
	ColumnNamesToExtractForDocs = []string{"modelDisplayName", "Page Subtitle", "Docs URL", "category", "subCategory", "Feature 1", "Feature 2", "Feature 3", "howItWorks", "howItWorksDetails", "Publish?", "About Project", "Standard Blurb", "svgColor", "svgWhite", "Full Page", "model"}
	PrimaryColumnName           = "model"
	OutputPath                  = ""
	ExcludeDirs                 = []string{"relationships", "policies"}
	GoogleSpreadSheetURL        = "https://docs.google.com/spreadsheets/d/"
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

	srv, err := pkg.NewSheetSRV()
	if err != nil {
		log.Fatal(err)
		return
	}
	resp, err := srv.Spreadsheets.Get(url).Fields().Do()
	if err != nil || resp.HTTPStatusCode != 200 {
		log.Fatal(err)
		return
	}

	modelCSVHelper := &pkg.ModelCSVHelper{}
	componentCSVHelper := &pkg.ComponentCSVHelper{}
	GoogleSpreadSheetURL += url

	for _, v := range resp.Sheets {
		switch v.Properties.Title {
		case "Models":
			modelCSVHelper, err = pkg.NewModelCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
			if err != nil {
				log.Fatal(err)
				return
			}
			modelCSVHelper.ParseModelsSheet()
		case "Components":
			componentCSVHelper, err = pkg.NewComponentCSVHelper(GoogleSpreadSheetURL, v.Properties.Title, v.Properties.SheetId)
			if err != nil {
				log.Fatal(err)
				return
			}
			componentCSVHelper.ParseComponentsSheet()
		}
	}
	// filep, err := pkg.DownloadCSV(url)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// file, err := os.Open(filep)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// csvReader := csv.NewReader(file)
	// output, err := pkg.GetEntries(csvReader, ColumnNamesToExtractForDocs)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }
	// file.Close()
	// os.Remove(file.Name())

	switch System {
	case pkg.Docs.String():
		docsUpdater(modelCSVHelper.Models, componentCSVHelper.Components)
	// commenting these for now
	// case pkg.Meshery.String():
	// 	mesheryUpdater(output)
	// case pkg.RemoteProvider.String():
	// 	remoteProviderUpdater(output)
	default:
		log.Fatal("invalid system name")
		return
	}

	err = modelCSVHelper.Cleanup()
	if err != nil {
		log.Fatal(err)
	}

	err = componentCSVHelper.Cleanup()
	if err != nil {
		log.Fatal(err)
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

func docsUpdater(models []pkg.ModelCSV, components map[string]map[string][]pkg.ComponentCSV) {
	if len(os.Args) < 7 {
		log.Fatal("docsUpdater: invalid number of arguments; missing website and docs path")
		return
	}
	pathToIntegrationsLayer5 := os.Args[4]
	pathToIntegrationsMeshery := os.Args[5]
	pathToIntegrationsMesheryDocs := os.Args[6]
	mesheryioDocsJSON := "const data = ["
	for _, model := range models {
		pathForLayer5ioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsLayer5))
		pathForMesheryioIntegrations, _ := filepath.Abs(filepath.Join("../../../", pathToIntegrationsMeshery))
		pathForMesheryDocsIntegrations, _ := filepath.Abs(filepath.Join("../../", pathToIntegrationsMesheryDocs))

		comps, ok := components[model.Registrant][model.Model]
		if !ok {
			fmt.Println("no components found for ", model.Model)
			comps = []pkg.ComponentCSV{}
		}

		err := pkg.GenerateLayer5Docs(model, comps, pathForLayer5ioIntegrations)
		if err != nil {
			fmt.Printf("Error generating layer5 docs for model %s: %v\n", model.Model, err.Error())
		}

		mesheryioDocsJSON, err = pkg.GenerateMesheryioDocs(model, pathForMesheryioIntegrations, mesheryioDocsJSON)
		if err != nil {
			fmt.Printf("Error generating mesheryio docs for model %s: %v\n", model.Model, err.Error())
		}

		err = pkg.GenerateMesheryDocs(model, comps, pathForMesheryDocsIntegrations)
		if err != nil {
			fmt.Printf("Error generating meshery docs for model %s: %v\n", model.Model, err.Error())
		}

	}

	mesheryioDocsJSON = strings.TrimSuffix(mesheryioDocsJSON, ",")
	mesheryioDocsJSON += "]; export default data"
	if err := pkg.WriteToFile(filepath.Join("../../../", pathToIntegrationsMeshery, "data.js"), mesheryioDocsJSON); err != nil {
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
		if changeFields["PublishToRegistry"] == "TRUE" { //For a component level field
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
						} else if key == "isModelAnnotation" && changeFields["component"] == "" {
							if component.Model.Metadata == nil {
								component.Model.Metadata = make(map[string]interface{})
							}
							if value == "TRUE" || value == "true" {
								component.Model.Metadata["isAnnotation"] = true
							} else {
								component.Model.Metadata["isAnnotation"] = false
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
					if component.Metadata["PublishToRegistry"] == "TRUE" || publishedModels[component.Model.Name] { //Publish? is an invalid field for putting inside kubernetes annotations
						component.Metadata["published"] = true
					} else {
						component.Metadata["published"] = false
					}
					if component.Metadata["isAnnotation"] == "TRUE" {
						component.Metadata["isAnnotation"] = true
					} else {
						component.Metadata["isAnnotation"] = false
					}

					fmt.Println("updating for ", changeFields["modelDisplayName"], "--", component.Kind, "-- published=", component.Metadata["published"])
					delete(component.Metadata, "Publish?")
					delete(component.Metadata, "PublishToRegistry")
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
