package utils

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	cuecsv "cuelang.org/go/pkg/encoding/csv"
	"github.com/gocarina/gocsv"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetData struct {
	Model              *ModelCSV
	Components         []component.ComponentDefinition
	NewComponentCount  int
	ExistingComponents int
	Version            string
}

var (
	compBatchSize  = 100
	modelBatchSize = 100

	// refers to all cells in the fourth row (explain how is table clulte in ghsset and henc this rang is valid and will not reuire udpate or when it hould require update)
	ComponentsSheetAppendRange = "Components!A4"
	ModelsSheetAppendRange     = "Models!A4"
)

// registrant:model:component:[true/false]
// Tracks if component sheet requires update
var RegistrantToModelsToComponentsMap = make(map[string]map[string]map[string]bool)

// registrant:model:[true/false]
// Tracks if component sheet requires update
var RegistrantToModelsMap = make(map[string]map[string]bool)

func ProcessModelToComponentsMap(existingComponents map[string]map[string][]ComponentCSV) {

	RegistrantToModelsToComponentsMap = make(map[string]map[string]map[string]bool, len(existingComponents))
	for registrant, models := range existingComponents {
		for model, comps := range models {
			if RegistrantToModelsToComponentsMap[registrant] == nil {
				RegistrantToModelsToComponentsMap[registrant] = make(map[string]map[string]bool)
			}
			for _, comp := range comps {
				if RegistrantToModelsToComponentsMap[registrant][model] == nil {
					RegistrantToModelsToComponentsMap[registrant][model] = make(map[string]bool)
				}
				RegistrantToModelsToComponentsMap[registrant][model][comp.Component] = true
			}
		}
	}
}

func addEntriesInCompUpdateList(modelEntry *ModelCSV, compEntries []component.ComponentDefinition, compList []*ComponentCSV) []*ComponentCSV {
	registrant := modelEntry.Registrant
	model := modelEntry.Model

	if RegistrantToModelsToComponentsMap[registrant][model] == nil {
		RegistrantToModelsToComponentsMap[registrant][model] = make(map[string]bool)
	}

	for _, comp := range compEntries {
		if !RegistrantToModelsToComponentsMap[registrant][model][comp.Component.Kind] {
			RegistrantToModelsToComponentsMap[registrant][model][comp.Component.Kind] = true
			compList = append(compList, ConvertCompDefToCompCSV(modelEntry, comp))
			compBatchSize--
		}
	}

	return compList
}

func addEntriesInModelUpdateList(modelEntry *ModelCSV, modelList []*ModelCSV) []*ModelCSV {
	registrant := modelEntry.Registrant

	if RegistrantToModelsMap[registrant] == nil {
		RegistrantToModelsMap[registrant] = make(map[string]bool)
	}
	RegistrantToModelsMap[registrant][modelEntry.Model] = true
	modelBatchSize--

	return modelList
}

// Verifies if the component entry already exist in the spreadsheet, otherwise updates the spreadshhet to include new component entry.
func VerifyandUpdateSpreadsheet(cred string, wg *sync.WaitGroup, srv *sheets.Service, spreadsheetUpdateChan chan SpreadsheetData, sheetId, modelCSVFilePath, componentCSVFilePath string) {
	defer wg.Done()

	entriesToBeAddedInCompSheet := []*ComponentCSV{}
	entriesToBeAddedInModelSheet := []*ModelCSV{}

	for data := range spreadsheetUpdateChan {
		_, ok := RegistrantToModelsMap[data.Model.Registrant]
		if !ok {
			entriesToBeAddedInModelSheet = addEntriesInModelUpdateList(data.Model, entriesToBeAddedInModelSheet)
		}

		for _, comp := range data.Components {
			existingModels, ok := RegistrantToModelsToComponentsMap[data.Model.Registrant] // replace with registrantr
			if ok {
				existingComps, ok := existingModels[data.Model.Model]

				if ok {
					entryExist := existingComps[comp.Component.Kind]

					if !entryExist {
						entriesToBeAddedInCompSheet = append(entriesToBeAddedInCompSheet, ConvertCompDefToCompCSV(data.Model, comp))
						compBatchSize--
						RegistrantToModelsToComponentsMap[data.Model.Registrant][data.Model.Model][comp.Component.Kind] = true
					}
				} else {
					entriesToBeAddedInCompSheet = addEntriesInCompUpdateList(data.Model, data.Components, entriesToBeAddedInCompSheet)
				}
			} else {
				RegistrantToModelsToComponentsMap[data.Model.Registrant] = make(map[string]map[string]bool)
				entriesToBeAddedInCompSheet = addEntriesInCompUpdateList(data.Model, data.Components, entriesToBeAddedInCompSheet)
			}
		}

		if modelBatchSize <= 0 {
			err := updateModelsSheet(srv, cred, sheetId, entriesToBeAddedInModelSheet)
			// Reset the list
			entriesToBeAddedInModelSheet = []*ModelCSV{}
			if err != nil {
				Log.Error(err)
			}
		}

		if compBatchSize <= 0 {
			err := updateComponentsSheet(srv, cred, sheetId, entriesToBeAddedInCompSheet, componentCSVFilePath)
			// update comp spreadsheet

			// Reset the list
			entriesToBeAddedInCompSheet = []*ComponentCSV{}
			entriesToBeAddedInModelSheet = []*ModelCSV{}
			if err != nil {
				Log.Error(err)
			}
		}
	}

	if len(entriesToBeAddedInModelSheet) > 0 {
		err := updateModelsSheet(srv, cred, sheetId, entriesToBeAddedInModelSheet)
		if err != nil {
			Log.Error(err)
		}
	}

	if len(entriesToBeAddedInCompSheet) > 0 {
		err := updateComponentsSheet(srv, cred, sheetId, entriesToBeAddedInCompSheet, componentCSVFilePath)
		if err != nil {
			Log.Error(err)
		}
		return
	}
}

func (mrh *RelationshipCSVHelper) UpdateRelationshipSheet(srv *sheets.Service, cred, sheetId, csvPath string) error {
	if len(mrh.UpdatedRelationships) == 0 {
		return nil
	}

	// Prepare batch update data with logging
	Log.Infof("Preparing batch update for %d relationships", len(mrh.UpdatedRelationships))
	var dataToUpdate []*sheets.ValueRange

	for _, rel := range mrh.UpdatedRelationships {
		row := rel.RowIndex
		// Column O corresponds to column 15 in A1 notation (O is the 15th letter)
		rangeStr := fmt.Sprintf("Relationships!O%d", row)

		Log.Debugf("Adding update for relationship at row %d: %s", row, rel.Filename)
		valueRange := &sheets.ValueRange{
			Range:  rangeStr,
			Values: [][]interface{}{{rel.Filename}},
		}
		dataToUpdate = append(dataToUpdate, valueRange)
	}

	// Perform batch update
	Log.Info("Executing batch update to spreadsheet...")
	batchUpdateValuesRequest := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data:             dataToUpdate,
	}

	_, err := srv.Spreadsheets.Values.BatchUpdate(sheetId, batchUpdateValuesRequest).
		Context(context.Background()).
		Do()

	if err != nil {
		return ErrUpdateToSheet(fmt.Errorf("Failed to update relationships (sheet ID: %s): %v", sheetId, err), sheetId)
	}

	Log.Infof("Successfully updated %d relationships in spreadsheet", len(dataToUpdate))
	return nil
}

func updateModelsSheet(srv *sheets.Service, cred, sheetId string, values []*ModelCSV) error {
	marshalledValues, err := marshalStructToCSValues[ModelCSV](values)
	if err != nil {
		return err
	}
	Log.Info("Appending", len(marshalledValues), " components into models sheet") // appendSheet or appendToCSVFile
	err = appendSheet(srv, cred, sheetId, ModelsSheetAppendRange, marshalledValues)

	return err
}

func updateComponentsSheet(srv *sheets.Service, cred, sheetId string, values []*ComponentCSV, csvPath string) error {
	marshalledValues, err := marshalStructToCSValues[ComponentCSV](values)
	Log.Info("Appending ", len(marshalledValues), "in the components sheet")
	if err != nil {
		return err
	}
	err = appendSheet(srv, cred, sheetId, ComponentsSheetAppendRange, marshalledValues, csvPath)

	return err
}

func appendSheet(srv *sheets.Service, cred, sheetId, appendRange string, values [][]interface{}, csvPaths ...string) error {
	csvPath := csvPaths[0]
	if csvPath != "" {
		// Open the CSV file in append mode, create if it doesn't exist
		file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer file.Close()
		// Create a new CSV writer
		writer := csv.NewWriter(file)

		// Write each row to the CSV
		for _, row := range values {
			stringRow := make([]string, len(row))
			for i, val := range row {
				stringRow[i] = fmt.Sprintf("%v", val) // Convert each value to a string
			}

			err := writer.Write(stringRow)
			if err != nil {
				return err
			}
		}

		// Flush the writer to ensure all data is written to the file
		writer.Flush()

		// Check for any error during flush
		if err := writer.Error(); err != nil {
			return err
		}

		return nil
	} else {
		if len(values) == 0 {
			return nil
		}
		_, err := srv.Spreadsheets.Values.Append(sheetId, appendRange, &sheets.ValueRange{
			MajorDimension: "ROWS",
			Range:          appendRange,
			Values:         values,
		}).InsertDataOption("INSERT_ROWS").ValueInputOption("USER_ENTERED").Context(context.Background()).Do()

		if err != nil {
			return ErrAppendToSheet(err, sheetId)
		}
		return nil
	}
}

func marshalStructToCSValues[K any](data []*K) ([][]interface{}, error) {
	csvString, err := gocsv.MarshalString(data)
	if err != nil {
		return nil, ErrMarshalStructToCSV(err)
	}
	csvReader := bytes.NewBufferString(csvString)
	decodedCSV, err := cuecsv.Decode(csvReader)

	if err != nil {
		return nil, ErrMarshalStructToCSV(err)
	}

	results := make([][]interface{}, 0)

	// The ouput is [ [col-names...] [row1][row2]]
	// delete the first entry i.e. [col-names..], as it contains the column names and is not required as we are concerened only with rows
	if len(decodedCSV) > 0 {
		for idx, val := range decodedCSV {
			if idx == 0 {
				continue
			}
			result := make([]interface{}, 0, cap(val))
			for _, r := range val {
				result = append(result, r)
			}
			results = append(results, result)
		}
		return results, nil
	}

	return results, nil
}
func GetCsv(csvDirectory string) (string, string, string, error) {
	files, err := os.ReadDir(csvDirectory)
	if err != nil {
		return "", "", "", utils.ErrReadDir(err, csvDirectory)
	}
	var modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath string

	for _, file := range files {
		filePath := filepath.Join(csvDirectory, file.Name())
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".csv") {
			headers, secondRow, err := getCSVHeader(filePath)
			if Contains("modelDisplayName", headers) != -1 || Contains("modelDisplayName", secondRow) != -1 {
				modelCSVFilePath = filePath
			} else if Contains("component", headers) != -1 || Contains("component", secondRow) != -1 {
				componentCSVFilePath = filePath
			} else if Contains("kind", headers) != -1 || Contains("kind", secondRow) != -1 { // Check if the file matches the relationshipCSV structure
				relationshipCSVFilePath = filePath
			}
			if err != nil {
				return "", "", "", err
			}

		}
	}

	if modelCSVFilePath == "" || componentCSVFilePath == "" {
		return "", "", "", ErrCSVFileNotFound(csvDirectory)
	}
	return modelCSVFilePath, componentCSVFilePath, relationshipCSVFilePath, nil
}
func getCSVHeader(filePath string) (headers, secondRow []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		err = utils.ErrOpenFile(filePath)
		return headers, secondRow, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err = reader.Read()
	if err != nil {
		err = ErrReadCSVRow(err, "header")
		return headers, secondRow, err
	}

	secondRow, err = reader.Read()
	if err != nil {
		err = ErrReadCSVRow(err, "second row")
		return headers, secondRow, err
	}
	return headers, secondRow, nil
}
