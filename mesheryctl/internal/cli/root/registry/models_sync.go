// # Copyright Meshery Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package registry

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	meshkitRegistryUtils "github.com/meshery/meshkit/registry"
	"github.com/meshery/schemas/models/v1beta1/category"
	"github.com/meshery/schemas/models/v1beta1/subcategory"
)

// modelDefJSON represents the structure of a committed model.json file.
type modelDefJSON struct {
	ID            string `json:"id"`
	SchemaVersion string `json:"schemaVersion"`
	Version       string `json:"version"`
	Name          string `json:"name"`
	DisplayName   string `json:"displayName"`
	Status        string `json:"status"`
	Description   string `json:"description,omitempty"`
	Registrant    struct {
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Type    string `json:"type"`
		SubType string `json:"subType"`
	} `json:"registrant"`
	Category struct {
		Name string `json:"name"`
	} `json:"category"`
	SubCategory string                 `json:"subCategory"`
	Metadata    map[string]interface{} `json:"metadata"`
	Model       struct {
		Version string `json:"version"`
	} `json:"model"`
}

// componentDefJSON represents the structure of a committed component JSON file.
type componentDefJSON struct {
	ID            string `json:"id"`
	SchemaVersion string `json:"schemaVersion"`
	Version       string `json:"version"`
	DisplayName   string `json:"displayName"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	Model         struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Version     string `json:"version"`
		Registrant  struct {
			Name string `json:"name"`
			Kind string `json:"kind"`
		} `json:"registrant"`
		Metadata map[string]interface{} `json:"metadata"`
		Model    struct {
			Version string `json:"version"`
		} `json:"model"`
	} `json:"model"`
	Styles       map[string]interface{} `json:"styles"`
	Capabilities interface{}            `json:"capabilities"`
	Metadata     map[string]interface{} `json:"metadata"`
	Component    struct {
		Version string `json:"version"`
		Kind    string `json:"kind"`
		Schema  string `json:"schema"`
	} `json:"component"`
}

func getStringFromMap(m map[string]interface{}, keys ...string) string {
	if m == nil {
		return ""
	}
	for _, k := range keys {
		if val, ok := m[k]; ok && val != nil {
			if s, ok := val.(string); ok && s != "" {
				return s
			}
		}
	}
	return ""
}

func stringifyJSONField(val interface{}) string {
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return ""
		}
		s := string(b)
		if s == "null" || s == "{}" || s == "[]" {
			return ""
		}
		return s
	}
}

func boolToCSVString(val interface{}) string {
	if val == nil {
		return "FALSE"
	}
	switch v := val.(type) {
	case bool:
		if v {
			return "TRUE"
		}
		return "FALSE"
	case string:
		if strings.EqualFold(v, "true") {
			return "TRUE"
		}
		return "FALSE"
	default:
		return "FALSE"
	}
}

func logWarnf(format string, args ...interface{}) {
	if utils.Log != nil {
		utils.Log.Warnf(format, args...)
	}
}

func logInfof(format string, args ...interface{}) {
	if utils.Log != nil {
		utils.Log.Infof(format, args...)
	}
}

func logErrorf(format string, args ...interface{}) {
	if utils.Log != nil {
		utils.Log.Errorf(format, args...)
	}
}

// extractSourceURLFromComponentsDir checks the sibling components directory for source_uri.
func extractSourceURLFromComponentsDir(compDir string) string {
	entries, err := os.ReadDir(compDir)
	if err != nil {
		return ""
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		compPath := filepath.Join(compDir, entry.Name())
		raw, err := os.ReadFile(compPath)
		if err != nil {
			continue
		}
		var compBrief struct {
			Metadata struct {
				SourceURI string `json:"source_uri"`
			} `json:"metadata"`
			Model struct {
				Metadata struct {
					SourceURI string `json:"source_uri"`
				} `json:"metadata"`
			} `json:"model"`
		}
		if err := json.Unmarshal(raw, &compBrief); err == nil {
			if compBrief.Metadata.SourceURI != "" {
				return compBrief.Metadata.SourceURI
			}
			if compBrief.Model.Metadata.SourceURI != "" {
				return compBrief.Model.Metadata.SourceURI
			}
		}
	}
	return ""
}

// ParseModelJSONFile parses a single model.json file into a meshkitRegistryUtils.ModelCSV.
func ParseModelJSONFile(filePath string) (*meshkitRegistryUtils.ModelCSV, error) {
	rawBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var def modelDefJSON
	if err := json.Unmarshal(rawBytes, &def); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	modelName := def.Name
	if modelName == "" {
		slashPath := filepath.ToSlash(filePath)
		parts := strings.Split(slashPath, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == "model.json" && i >= 3 {
				modelName = parts[i-3]
				break
			}
		}
	}

	displayName := def.DisplayName
	if displayName == "" {
		displayName = modelName
	}

	registrant := def.Registrant.Name
	if registrant == "" {
		registrant = def.Registrant.Kind
	}
	if registrant == "" {
		registrant = "Github"
	}

	cat := def.Category.Name
	subCat := def.SubCategory
	if subCat == "" && def.Metadata != nil {
		subCat = getStringFromMap(def.Metadata, "subCategory")
	}

	desc := def.Description
	if desc == "" && def.Metadata != nil {
		desc = getStringFromMap(def.Metadata, "description")
	}

	var shape, primaryColor, secondaryColor, styleOverrides, shapePolygonPoints, logoURL, svgColor, svgWhite, svgComplete string
	var stylesStr, defaultDataStr, capabilitiesStr, sourceURL string
	isAnnotationStr := "FALSE"

	if def.Metadata != nil {
		shape = getStringFromMap(def.Metadata, "shape")
		primaryColor = getStringFromMap(def.Metadata, "primaryColor")
		secondaryColor = getStringFromMap(def.Metadata, "secondaryColor")
		styleOverrides = getStringFromMap(def.Metadata, "styleOverrides")
		shapePolygonPoints = getStringFromMap(def.Metadata, "shapePolygonPoints", "shape-polygon-points")
		logoURL = getStringFromMap(def.Metadata, "logoURL")
		svgColor = getStringFromMap(def.Metadata, "svgColor")
		svgWhite = getStringFromMap(def.Metadata, "svgWhite")
		svgComplete = getStringFromMap(def.Metadata, "svgComplete")
		sourceURL = getStringFromMap(def.Metadata, "source_uri", "sourceURL")

		stylesStr = stringifyJSONField(def.Metadata["styles"])
		defaultDataStr = stringifyJSONField(def.Metadata["defaultData"])
		capabilitiesStr = stringifyJSONField(def.Metadata["capabilities"])
		isAnnotationStr = boolToCSVString(def.Metadata["isAnnotation"])
	}

	// If sourceURL is still empty, look into sibling components directory
	if sourceURL == "" {
		compDir := filepath.Join(filepath.Dir(filePath), "components")
		sourceURL = extractSourceURLFromComponentsDir(compDir)
	}

	publishToRegistry := "FALSE"
	if strings.EqualFold(def.Status, "enabled") {
		publishToRegistry = "TRUE"
	} else if def.Metadata != nil {
		if pub, ok := def.Metadata["published"].(bool); ok && pub {
			publishToRegistry = "TRUE"
		}
	}

	return &meshkitRegistryUtils.ModelCSV{
		Registrant:         registrant,
		ModelDisplayName:   displayName,
		Model:              modelName,
		Category:           category.CategoryDefinitionName(cat),
		SubCategory:        subcategory.SubCategoryDefinition(subCat),
		Description:        desc,
		SourceURL:          sourceURL,
		Website:            "",
		Docs:               "",
		Shape:              shape,
		PrimaryColor:       primaryColor,
		SecondaryColor:     secondaryColor,
		StyleOverrides:     styleOverrides,
		Styles:             stylesStr,
		ShapePolygonPoints: shapePolygonPoints,
		DefaultData:        defaultDataStr,
		Capabilities:       capabilitiesStr,
		LogoURL:            logoURL,
		SVGColor:           svgColor,
		SVGWhite:           svgWhite,
		SVGComplete:        svgComplete,
		IsAnnotation:       isAnnotationStr,
		PublishToRegistry:  publishToRegistry,
		Group:              "",
	}, nil
}

// ParseComponentJSONFile parses a single component JSON file into a meshkitRegistryUtils.ComponentCSV.
func ParseComponentJSONFile(filePath string) (*meshkitRegistryUtils.ComponentCSV, error) {
	rawBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var def componentDefJSON
	if err := json.Unmarshal(rawBytes, &def); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	modelName := def.Model.Name
	if modelName == "" {
		slashPath := filepath.ToSlash(filePath)
		parts := strings.Split(slashPath, "/")
		for i := len(parts) - 1; i >= 0; i-- {
			if parts[i] == "components" && i >= 2 {
				modelName = parts[i-2]
				break
			}
		}
	}

	compName := def.Component.Kind
	if compName == "" {
		compName = def.DisplayName
	}
	if compName == "" {
		base := filepath.Base(filePath)
		compName = strings.TrimSuffix(base, filepath.Ext(base))
	}

	registrant := def.Model.Registrant.Name
	if registrant == "" {
		registrant = def.Model.Registrant.Kind
	}
	if registrant == "" {
		registrant = "Github"
	}

	desc := def.Description
	if desc == "" && def.Metadata != nil {
		desc = getStringFromMap(def.Metadata, "description")
	}

	shape := getStringFromMap(def.Styles, "shape")
	if shape == "" {
		shape = getStringFromMap(def.Metadata, "shape")
	}
	if shape == "" {
		shape = getStringFromMap(def.Model.Metadata, "shape")
	}

	primaryColor := getStringFromMap(def.Styles, "primaryColor")
	if primaryColor == "" {
		primaryColor = getStringFromMap(def.Metadata, "primaryColor")
	}
	if primaryColor == "" {
		primaryColor = getStringFromMap(def.Model.Metadata, "primaryColor")
	}

	secondaryColor := getStringFromMap(def.Styles, "secondaryColor")
	if secondaryColor == "" {
		secondaryColor = getStringFromMap(def.Metadata, "secondaryColor")
	}
	if secondaryColor == "" {
		secondaryColor = getStringFromMap(def.Model.Metadata, "secondaryColor")
	}

	svgColor := getStringFromMap(def.Styles, "svgColor")
	if svgColor == "" {
		svgColor = getStringFromMap(def.Metadata, "svgColor")
	}
	if svgColor == "" {
		svgColor = getStringFromMap(def.Model.Metadata, "svgColor")
	}

	svgWhite := getStringFromMap(def.Styles, "svgWhite")
	if svgWhite == "" {
		svgWhite = getStringFromMap(def.Metadata, "svgWhite")
	}
	if svgWhite == "" {
		svgWhite = getStringFromMap(def.Model.Metadata, "svgWhite")
	}

	svgComplete := getStringFromMap(def.Styles, "svgComplete")
	if svgComplete == "" {
		svgComplete = getStringFromMap(def.Metadata, "svgComplete")
	}
	if svgComplete == "" {
		svgComplete = getStringFromMap(def.Model.Metadata, "svgComplete")
	}

	styleOverrides := getStringFromMap(def.Styles, "styleOverrides")
	if styleOverrides == "" {
		styleOverrides = getStringFromMap(def.Metadata, "styleOverrides")
	}
	if styleOverrides == "" {
		styleOverrides = getStringFromMap(def.Model.Metadata, "styleOverrides")
	}

	shapePolygonPoints := getStringFromMap(def.Styles, "shapePolygonPoints", "shape-polygon-points")
	if shapePolygonPoints == "" {
		shapePolygonPoints = getStringFromMap(def.Metadata, "shapePolygonPoints", "shape-polygon-points")
	}

	var defaultDataStr, logoURL, genealogy string
	isAnnotationStr := "FALSE"
	if def.Metadata != nil {
		defaultDataStr = stringifyJSONField(def.Metadata["defaultData"])
		logoURL = getStringFromMap(def.Metadata, "logoURL")
		genealogy = getStringFromMap(def.Metadata, "genealogy")
		isAnnotationStr = boolToCSVString(def.Metadata["isAnnotation"])
	}

	stylesStr := stringifyJSONField(def.Styles)
	capabilitiesStr := stringifyJSONField(def.Capabilities)

	version := def.Component.Version
	if version == "" {
		version = def.Model.Model.Version
	}
	if version == "" {
		version = def.Model.Version
	}

	status := def.Status
	if status == "" {
		status = "enabled"
	}

	return &meshkitRegistryUtils.ComponentCSV{
		Registrant:         registrant,
		Model:              modelName,
		Component:          compName,
		Description:        desc,
		Shape:              shape,
		PrimaryColor:       primaryColor,
		SecondaryColor:     secondaryColor,
		SVGColor:           svgColor,
		SVGWhite:           svgWhite,
		SVGComplete:        svgComplete,
		Schema:             def.Component.Schema,
		Docs:               "",
		StyleOverrides:     styleOverrides,
		Styles:             stylesStr,
		ShapePolygonPoints: shapePolygonPoints,
		DefaultData:        defaultDataStr,
		Capabilities:       capabilitiesStr,
		LogoURL:            logoURL,
		Genealogy:          genealogy,
		IsAnnotation:       isAnnotationStr,
		Version:            version,
		ModelDisplayName:   def.Model.DisplayName,
		Status:             status,
	}, nil
}

// ScanCommittedModels scans the models directory and parses all model.json files.
func ScanCommittedModels(modelsPath string, targetModel string) ([]meshkitRegistryUtils.ModelCSV, error) {
	var rows []meshkitRegistryUtils.ModelCSV
	seen := make(map[string]bool)

	err := filepath.WalkDir(modelsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Base(path) != "model.json" {
			return nil
		}

		relPath, _ := filepath.Rel(modelsPath, path)
		pathParts := strings.Split(filepath.ToSlash(relPath), "/")

		if targetModel != "" && len(pathParts) > 0 && pathParts[0] != targetModel {
			return nil
		}

		row, err := ParseModelJSONFile(path)
		if err != nil {
			logWarnf("Failed to parse model file %s: %v", path, err)
			return nil
		}

		modelVersion := ""
		if len(pathParts) >= 2 {
			modelVersion = pathParts[1]
		}
		key := fmt.Sprintf("%s:%s", row.Model, modelVersion)
		if seen[key] {
			return nil
		}
		seen[key] = true

		rows = append(rows, *row)
		return nil
	})

	return rows, err
}

// ScanCommittedComponents scans the models directory and parses all component JSON files.
func ScanCommittedComponents(modelsPath string, targetModel string) ([]meshkitRegistryUtils.ComponentCSV, error) {
	var rows []meshkitRegistryUtils.ComponentCSV
	seen := make(map[string]bool)

	err := filepath.WalkDir(modelsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slashPath := filepath.ToSlash(path)
		if d.IsDir() || !strings.HasSuffix(path, ".json") || !strings.Contains(slashPath, "/components/") {
			return nil
		}

		relPath, _ := filepath.Rel(modelsPath, path)
		pathParts := strings.Split(filepath.ToSlash(relPath), "/")

		if targetModel != "" && len(pathParts) > 0 && pathParts[0] != targetModel {
			return nil
		}

		row, err := ParseComponentJSONFile(path)
		if err != nil {
			logWarnf("Failed to parse component file %s: %v", path, err)
			return nil
		}

		modelVersion := ""
		if len(pathParts) >= 2 {
			modelVersion = pathParts[1]
		}
		if row.Version == "" {
			row.Version = modelVersion
		}

		key := fmt.Sprintf("%s:%s:%s", row.Model, modelVersion, row.Component)
		if seen[key] {
			return nil
		}
		seen[key] = true

		rows = append(rows, *row)
		return nil
	})

	return rows, err
}

// ExportModelsToCSV exports model rows to a CSV file matching the Models spreadsheet schema.
func ExportModelsToCSV(rows []meshkitRegistryUtils.ModelCSV, outputPath string) (err error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %w", outputPath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if err := gocsv.MarshalFile(&rows, file); err != nil {
		return fmt.Errorf("failed to write models CSV %s: %w", outputPath, err)
	}

	return nil
}

// ExportComponentsToCSV exports component rows to a CSV file matching the Components spreadsheet schema.
func ExportComponentsToCSV(rows []meshkitRegistryUtils.ComponentCSV, outputPath string) (err error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %w", outputPath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if err := gocsv.MarshalFile(&rows, file); err != nil {
		return fmt.Errorf("failed to write components CSV %s: %w", outputPath, err)
	}

	return nil
}

// InvokeModelsExport scans committed models and exports them to a CSV file.
func InvokeModelsExport(modelsPath string, targetModel string, exportCSVPath string) error {
	logInfof("Scanning committed model definitions in %s...", modelsPath)
	rows, err := ScanCommittedModels(modelsPath, targetModel)
	if err != nil {
		return err
	}
	logInfof("Found %d committed model definitions in %s", len(rows), modelsPath)

	if exportCSVPath != "" {
		if err := ExportModelsToCSV(rows, exportCSVPath); err != nil {
			logErrorf("Failed to export models CSV: %v", err)
			return fmt.Errorf("failed to export models CSV to %s: %w", exportCSVPath, err)
		}
		logInfof("Successfully exported %d model rows to %s", len(rows), exportCSVPath)
	}

	return nil
}

// InvokeComponentsExport scans committed components and exports them to a CSV file.
func InvokeComponentsExport(modelsPath string, targetModel string, exportCSVPath string) error {
	logInfof("Scanning committed component definitions in %s...", modelsPath)
	rows, err := ScanCommittedComponents(modelsPath, targetModel)
	if err != nil {
		return err
	}
	logInfof("Found %d committed component definitions in %s", len(rows), modelsPath)

	if exportCSVPath != "" {
		if err := ExportComponentsToCSV(rows, exportCSVPath); err != nil {
			logErrorf("Failed to export components CSV: %v", err)
			return fmt.Errorf("failed to export components CSV to %s: %w", exportCSVPath, err)
		}
		logInfof("Successfully exported %d component rows to %s", len(rows), exportCSVPath)
	}

	return nil
}
