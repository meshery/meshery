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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"google.golang.org/api/sheets/v4"
)

// RelationshipDefModel represents the schema of a committed relationship JSON file in the models directory.
type RelationshipDefModel struct {
	ID              string `json:"id"`
	EvaluationQuery string `json:"evaluationQuery"`
	Kind            string `json:"kind"`
	Type            string `json:"type,omitempty"`
	SubType         string `json:"subType,omitempty"`
	SchemaVersion   string `json:"schemaVersion"`
	Metadata        struct {
		Description  string                 `json:"description"`
		IsAnnotation bool                   `json:"isAnnotation"`
		Styles       map[string]interface{} `json:"styles"`
	} `json:"metadata"`
	Model struct {
		Name        string `json:"name"`
		DisplayName string `json:"displayName"`
		Version     string `json:"version"`
		Model       struct {
			Version string `json:"version"`
		} `json:"model"`
	} `json:"model"`
	Selectors []struct {
		Allow struct {
			From []map[string]interface{} `json:"from"`
			To   []map[string]interface{} `json:"to"`
		} `json:"allow"`
		Deny struct {
			From []map[string]interface{} `json:"from"`
			To   []map[string]interface{} `json:"to"`
		} `json:"deny"`
	} `json:"selectors"`
}

// RelationshipRow represents a row formatted for the Integrations Spreadsheet "Relationships" tab.
type RelationshipRow struct {
	Model                string
	Version              string
	Kind                 string
	Type                 string
	SubType              string
	SchemaVersion        string
	MetadataDescription  string
	Docs                 string
	MetadataStyles       string
	EvalPolicy           string
	SelectorsDenyFrom    string
	SelectorsDenyTo      string
	SelectorsAllowFrom   string
	SelectorsAllowTo     string
	Selector             string
	Filename             string
	CompleteDefinition   string
	MetadataIsAnnotation string
	VisualizationExample string
	FromKind             string
	ToKind               string
}

// ToCSVRow converts RelationshipRow to a slice of strings matching the Relationships spreadsheet schema.
func (r *RelationshipRow) ToCSVRow() []string {
	return []string{
		r.Model, r.Version, r.Kind, r.Type, r.SubType,
		r.MetadataDescription, r.Docs, r.MetadataStyles, r.EvalPolicy,
		r.SelectorsDenyFrom, r.SelectorsDenyTo, r.SelectorsAllowFrom, r.SelectorsAllowTo,
		r.Selector, r.Filename, r.CompleteDefinition, r.MetadataIsAnnotation, r.VisualizationExample,
		"", r.FromKind, r.Type, r.ToKind, "",
	}
}

// ToSheetRow converts RelationshipRow to a slice of interface values for the Google Sheets API.
func (r *RelationshipRow) ToSheetRow() []interface{} {
	csvRow := r.ToCSVRow()
	row := make([]interface{}, len(csvRow))
	for i, v := range csvRow {
		row[i] = v
	}
	return row
}

// ParseRelationshipJSONFile parses a single relationship JSON file into a RelationshipRow.
func ParseRelationshipJSONFile(filePath string) (*RelationshipRow, error) {
	rawBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var def RelationshipDefModel
	if err := json.Unmarshal(rawBytes, &def); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	fileName := filepath.Base(filePath)
	nameWithoutExt := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	tokens := strings.Split(nameWithoutExt, "-")

	kind := def.Kind
	if kind == "" && len(tokens) > 0 {
		kind = tokens[0]
	}

	relType := def.Type
	if relType == "" && len(tokens) >= 2 {
		if tokens[1] == "non" && len(tokens) >= 3 && tokens[2] == "binding" {
			relType = "non-binding"
		} else {
			relType = tokens[1]
		}
	}

	subType := def.SubType
	if subType == "" {
		for _, candidate := range []string{"reference", "inventory", "firewall", "mount", "network", "permission"} {
			if strings.Contains(nameWithoutExt, candidate) {
				subType = candidate
				break
			}
		}
		if subType == "" && len(tokens) >= 3 {
			subType = tokens[2]
		}
	}

	modelVersion := def.Model.Model.Version
	if modelVersion == "" {
		modelVersion = def.Model.Version
	}
	if modelVersion == "" {
		modelVersion = "*"
	}

	var fromKind, toKind, allowFromStr, allowToStr, denyFromStr, denyToStr, selectorsStr string
	if len(def.Selectors) > 0 {
		first := def.Selectors[0]
		if len(first.Allow.From) > 0 {
			if k, ok := first.Allow.From[0]["kind"].(string); ok {
				fromKind = k
			}
			b, _ := json.Marshal(first.Allow.From)
			allowFromStr = string(b)
		}
		if len(first.Allow.To) > 0 {
			if k, ok := first.Allow.To[0]["kind"].(string); ok {
				toKind = k
			}
			b, _ := json.Marshal(first.Allow.To)
			allowToStr = string(b)
		}
		if len(first.Deny.From) > 0 {
			b, _ := json.Marshal(first.Deny.From)
			denyFromStr = string(b)
		}
		if len(first.Deny.To) > 0 {
			b, _ := json.Marshal(first.Deny.To)
			denyToStr = string(b)
		}
		allB, _ := json.Marshal(def.Selectors)
		selectorsStr = string(allB)
	}

	stylesStr := ""
	if len(def.Metadata.Styles) > 0 {
		b, _ := json.Marshal(def.Metadata.Styles)
		stylesStr = string(b)
	}

	isAnnotationStr := "FALSE"
	if def.Metadata.IsAnnotation {
		isAnnotationStr = "TRUE"
	}

	var compactBuf bytes.Buffer
	completeDefStr := string(rawBytes)
	if err := json.Compact(&compactBuf, rawBytes); err == nil {
		completeDefStr = compactBuf.String()
	}

	return &RelationshipRow{
		Model:                def.Model.Name,
		Version:              modelVersion,
		Kind:                 kind,
		Type:                 relType,
		SubType:              subType,
		SchemaVersion:        def.SchemaVersion,
		MetadataDescription:  def.Metadata.Description,
		MetadataStyles:       stylesStr,
		EvalPolicy:           def.EvaluationQuery,
		SelectorsDenyFrom:    denyFromStr,
		SelectorsDenyTo:      denyToStr,
		SelectorsAllowFrom:   allowFromStr,
		SelectorsAllowTo:     allowToStr,
		Selector:             selectorsStr,
		Filename:             fileName,
		CompleteDefinition:   completeDefStr,
		MetadataIsAnnotation: isAnnotationStr,
		FromKind:             fromKind,
		ToKind:               toKind,
	}, nil
}

// ScanCommittedRelationships scans the models directory and parses all relationship JSON definitions.
func ScanCommittedRelationships(modelsPath string, targetModel string) ([]RelationshipRow, error) {
	var rows []RelationshipRow

	err := filepath.WalkDir(modelsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		slashPath := filepath.ToSlash(path)
		if d.IsDir() || !strings.HasSuffix(path, ".json") || !strings.Contains(slashPath, "/relationships/") {
			return nil
		}

		if targetModel != "" {
			relPath, _ := filepath.Rel(modelsPath, path)
			pathParts := strings.Split(filepath.ToSlash(relPath), "/")
			if len(pathParts) > 0 && pathParts[0] != targetModel {
				return nil
			}
		}

		row, err := ParseRelationshipJSONFile(path)
		if err != nil {
			utils.Log.Warnf("Failed to parse relationship file %s: %v", path, err)
			return nil
		}

		if row.Model == "" {
			relPath, _ := filepath.Rel(modelsPath, path)
			pathParts := strings.Split(filepath.ToSlash(relPath), "/")
			if len(pathParts) > 0 {
				row.Model = pathParts[0]
			}
		}
		rows = append(rows, *row)
		return nil
	})

	return rows, err
}

// ExportRelationshipsToCSV exports relationship rows to a CSV matching the Relationships spreadsheet template.
func ExportRelationshipsToCSV(rows []RelationshipRow, outputPath string) (err error) {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file %s: %w", outputPath, err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	writer := csv.NewWriter(file)

	header1 := []string{
		"Meshery Relationship Definitions",
		"If the specified version * it acts a wildcard that writes the relationships to all the version present for that particular model",
		"", "", "", "", "", "", "", "", "", "", "",
		"Generated from the combination of column J-M. \nJ - Deny from\nK - Deny to\nL - Allow from\nM - Allow to",
		"Generated by Golang.\nFilename: kind-type-subtype-randomstring.json",
		"", "", "", "",
		"These columns track components and whether they have been accounted for within an existing relationship or not.",
		"", "", "",
	}
	header2 := []string{
		"Model", "Version", "kind", "type", "subType",
		"metadata.description", "docs", "metadata.styles", "evalPolicy",
		"selectors.deny.from", "selectors.deny.to", "selectors.allow.from", "selectors.allow.to",
		"selector", "filename", "Complete Definition", "metadata.isAnnotation", "Visualization Example",
		"", "Component", "Relationship Type", "Relationship With", "Constraints",
	}

	if err := writer.Write(header1); err != nil {
		return err
	}
	if err := writer.Write(header2); err != nil {
		return err
	}

	for _, r := range rows {
		if err := writer.Write(r.ToCSVRow()); err != nil {
			return err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to write CSV %s: %w", outputPath, err)
	}

	return nil
}

// DeduplicateRelationshipRows filters out relationship rows that already exist in the spreadsheet values.
func DeduplicateRelationshipRows(rows []RelationshipRow, existingValues [][]interface{}) [][]interface{} {
	existingRows := make(map[string]bool)
	for idx, sRow := range existingValues {
		if idx < 2 || len(sRow) == 0 {
			continue
		}
		modelName := fmt.Sprintf("%v", sRow[0])
		filename := ""
		if len(sRow) > 14 {
			filename = fmt.Sprintf("%v", sRow[14])
		}
		existingRows[fmt.Sprintf("%s:%s", modelName, filename)] = true
	}

	var newSheetValues [][]interface{}
	for _, r := range rows {
		key := fmt.Sprintf("%s:%s", r.Model, r.Filename)
		if !existingRows[key] {
			newSheetValues = append(newSheetValues, r.ToSheetRow())
			existingRows[key] = true
		}
	}
	return newSheetValues
}

// InvokeRelationshipSync scans committed relationships and syncs them to Google Sheets and/or exports to CSV.
func InvokeRelationshipSync(srv *sheets.Service, spreadsheetID string, modelsPath string, targetModel string, exportCSVPath string) error {
	utils.Log.Info("Scanning committed relationship definitions in ", modelsPath, "...")
	rows, err := ScanCommittedRelationships(modelsPath, targetModel)
	if err != nil {
		return err
	}
	utils.Log.Info(fmt.Sprintf("Found %d committed relationship definitions in %s", len(rows), modelsPath))

	if exportCSVPath != "" {
		if err := ExportRelationshipsToCSV(rows, exportCSVPath); err != nil {
			utils.Log.Errorf("Failed to export relationships CSV: %v", err)
			return fmt.Errorf("failed to export relationships CSV to %s: %w", exportCSVPath, err)
		}
		utils.Log.Info(fmt.Sprintf("Successfully exported %d relationship rows to %s", len(rows), exportCSVPath))
	}

	if srv == nil || spreadsheetID == "" {
		utils.Log.Info("Spreadsheet credentials not provided; skipping Google Sheet remote sync.")
		return nil
	}

	utils.Log.Info("Syncing relationships to Google Spreadsheet...")
	sheetResp, err := srv.Spreadsheets.Values.Get(spreadsheetID, "Relationships").Do()
	if err != nil {
		return fmt.Errorf("failed to fetch existing relationships from sheet: %w", err)
	}

	newSheetValues := DeduplicateRelationshipRows(rows, sheetResp.Values)

	if len(newSheetValues) == 0 {
		utils.Log.Info("All committed relationships are already synced with the spreadsheet. No new rows to add.")
		return nil
	}

	appendCall := srv.Spreadsheets.Values.Append(spreadsheetID, "Relationships!A:W", &sheets.ValueRange{Values: newSheetValues})
	appendCall.ValueInputOption("RAW")
	appendResp, err := appendCall.Do()
	if err != nil {
		return fmt.Errorf("failed to append relationship rows to Google Spreadsheet: %w", err)
	}

	utils.Log.Info(fmt.Sprintf("Successfully appended %d new relationship definitions to Google Spreadsheet (Updates: %d rows affected)", len(newSheetValues), appendResp.Updates.UpdatedRows))
	return nil
}
