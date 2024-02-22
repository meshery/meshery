package utils

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	cuecsv "cuelang.org/go/pkg/encoding/csv"
	"github.com/gocarina/gocsv"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"google.golang.org/api/sheets/v4"
)

type SpreadsheetData struct {
	Model      *ModelCSV
	Components []v1alpha1.ComponentDefinition
}

var (
	compBatchSize  = 100
	modelBatchSize = 100
	// refers to all cells in the fourth row (explain how is table clulte in ghsset and henc this rang is valid and will not reuire udpate or when it hould require update)
	appendRange                = "A4"
	ComponentsSheetAppendRange = fmt.Sprintf("Components!%s", appendRange)
	ModelsSheetAppendRange     = fmt.Sprintf("Models!%s", appendRange)
)

// registrant:model:component:[true/false]
var modelToComponentsMap = make(map[string]map[string]map[string]bool)

func ProcessModelToComponentsMap(existingComponents map[string]map[string][]ComponentCSV) {
	modelToComponentsMap = make(map[string]map[string]map[string]bool, len(existingComponents))
	for registrant, models := range existingComponents {
		for model, comps := range models {
			if modelToComponentsMap[registrant] == nil {
				modelToComponentsMap[registrant] = make(map[string]map[string]bool)
			}
			fmt.Println("TEST ", model)
			for _, comp := range comps {
				if modelToComponentsMap[registrant][model] == nil {
					modelToComponentsMap[registrant][model] = make(map[string]bool)
				}
				modelToComponentsMap[registrant][model][comp.Component] = true
			}
		}
	}
}

func addEntriesInModelAndCompUpdateList(modelEntry *ModelCSV, compEntries []v1alpha1.ComponentDefinition, modelList []*ModelCSV, compList []*ComponentCSV) ([]*ModelCSV, []*ComponentCSV) {
	modelList = append(modelList, modelEntry)
	modelBatchSize--
	modelToComponentsMap[modelEntry.Registrant][modelEntry.Model] = make(map[string]bool)
	for _, comp := range compEntries {
		if modelToComponentsMap[modelEntry.Registrant][modelEntry.Model] == nil {
			modelToComponentsMap[modelEntry.Registrant][modelEntry.Model] = make(map[string]bool)
		}
		modelToComponentsMap[modelEntry.Registrant][modelEntry.Model][comp.Kind] = true
		compList = append(compList, ConvertCompDefToCompCSV(modelEntry, comp))
	}
	compBatchSize -= len(compEntries)

	return modelList, compList
}

// Verifies if the component entry already exist in the spreadsheet, otherwise updates the spreadshhet to include new component entry.
func VerifyandUpdateSpreadsheet(cred string, wg *sync.WaitGroup, srv *sheets.Service, spreadsheetUpdateChan chan SpreadsheetData, sheetId string) {
	defer wg.Done()

	entriesToBeAddedInCompSheet := []*ComponentCSV{}
	entriesToBeAddedInModelSheet := []*ModelCSV{}

	for data := range spreadsheetUpdateChan {
		for _, comp := range data.Components {
			existingModels, ok := modelToComponentsMap[data.Model.Registrant] // replace with registrantr
			if ok {
				fmt.Println("existing models : ", existingModels)
				existingComps, ok := existingModels[data.Model.Model]
				fmt.Println("existing comps: ", existingComps)
				if ok {
					entryExist := existingComps[comp.Kind]
					fmt.Println("test 57 : ", entryExist)
					if !entryExist {
						entriesToBeAddedInCompSheet = append(entriesToBeAddedInCompSheet, ConvertCompDefToCompCSV(data.Model, comp))
						compBatchSize--
						modelToComponentsMap[data.Model.Registrant][data.Model.Model][comp.Kind] = true
					}
				} else {
					entriesToBeAddedInModelSheet, entriesToBeAddedInCompSheet = addEntriesInModelAndCompUpdateList(data.Model, data.Components, entriesToBeAddedInModelSheet, entriesToBeAddedInCompSheet)
				}
			} else {
				fmt.Println("No registrant exist : ", data.Model.Registrant)
				entriesToBeAddedInModelSheet, entriesToBeAddedInCompSheet = addEntriesInModelAndCompUpdateList(data.Model, data.Components, entriesToBeAddedInModelSheet, entriesToBeAddedInCompSheet)
			}
		}

		if modelBatchSize <= 0 {
			// update model spreadsheet
			err := updateModelsSheet(srv, cred, sheetId, entriesToBeAddedInModelSheet)
			// Reset the list
			entriesToBeAddedInModelSheet = []*ModelCSV{}
			if err != nil {
				fmt.Println(err)
			}
		}

		if compBatchSize <= 0 {
			// update comp spreadsheet
			err := updateComponentsSheet(srv, cred, sheetId, entriesToBeAddedInCompSheet)
			// Reset the list
			entriesToBeAddedInCompSheet = []*ComponentCSV{}
			entriesToBeAddedInModelSheet = []*ModelCSV{}
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	fmt.Println("EXITING AFTER APPLYING PENDING UPDATES")
	if len(entriesToBeAddedInModelSheet) > 0 {
		err := updateModelsSheet(srv, cred, sheetId, entriesToBeAddedInModelSheet)
		if err != nil {
			fmt.Println(err)
		}
	}

	if len(entriesToBeAddedInModelSheet) > 0 {
		err := updateComponentsSheet(srv, cred, sheetId, entriesToBeAddedInCompSheet)
		if err != nil {
			fmt.Println(err)
		}
		return
	}
}

func updateModelsSheet(srv *sheets.Service, cred, sheetId string, values []*ModelCSV) error {
	marshalledValues, err := marshalStructToCSValues[ModelCSV](values)
	if err != nil {
		return err
	}
	err = appendSheet(srv, cred, sheetId, ModelsSheetAppendRange, marshalledValues)
	fmt.Println("appending models", err)
	return err
}

func updateComponentsSheet(srv *sheets.Service, cred, sheetId string, values []*ComponentCSV) error {
	marshalledValues, err := marshalStructToCSValues[ComponentCSV](values)
	if err != nil {
		return err
	}
	err = appendSheet(srv, cred, sheetId, ComponentsSheetAppendRange, marshalledValues)
	fmt.Println("appending compoennts", err, len(marshalledValues))
	return err
}

func appendSheet(srv *sheets.Service, cred, sheetId, appendRange string, values [][]interface{}) error {
	fmt.Println("LINE 97 : ")

	_, err := srv.Spreadsheets.Values.BatchUpdate(sheetId, &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
		Data: []*sheets.ValueRange{
			{
				MajorDimension: "ROWS",
				Range:          appendRange,
				Values:         values,
			},
		},
	}).Context(context.Background()).Do()
	if len(values) == 0 {
		return nil
	}

	if err != nil {
		return ErrAppendToSheet(err, sheetId)
	}
	return nil
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
