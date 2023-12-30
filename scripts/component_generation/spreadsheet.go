package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type spreadsheetUpdater struct {
	spreadsheetID               string
	sheetName                   string
	spreadsheetChan             chan componentWrapper
	availableModels             map[string][]interface{}
	availableComponentsPerModel map[string]map[string]bool
}

func NewSheetSRV() *sheets.Service {
	ctx := context.Background()
	byt, _ := base64.StdEncoding.DecodeString(os.Getenv("CRED"))
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(byt, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Error(ErrorFailedConnectingGCP(err))
		return nil
	}
	// create client with config and context
	client := config.Client(ctx)
	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Error(ErrorFailedConnectingGCP(err))
		return nil
	}
	return srv
}

func (updater *spreadsheetUpdater) update(srv *sheets.Service) {
	rangeString := updater.sheetName + APPENDRANGE
	// Get the value of the specified cell.
	resp, err := srv.Spreadsheets.Values.Get(updater.spreadsheetID, rangeString).Do()
	if err != nil {
		log.Error(ErrorFailedRetreivingSheet(err))
		return
	}
	values := make([][]interface{}, 0)
	for entry := range updater.spreadsheetChan {
		if len(entry.comps) == 0 {
			continue
		}
		for _, comp := range entry.comps {
			if updater.availableComponentsPerModel[entry.model][comp.Kind] {
				log.Debug(fmt.Sprintf("Skipping spreadsheet updation for %s %s", entry.model, comp.Kind))
				continue
			}
			var newValues []interface{}
			if updater.availableModels[entry.model] != nil {
				newValues = make([]interface{}, len(updater.availableModels[entry.model]))
				copy(newValues, updater.availableModels[entry.model])
			} else {
				newValues = make([]interface{}, len(resp.Values[0]))
				copy(newValues, resp.Values[0])
				newValues[NameToIndex["modelDisplayName"]] = entry.model
				newValues[NameToIndex["model"]] = entry.model
			}
			newValues[NameToIndex["component"]] = comp.Kind
			newValues[NameToIndex["hasSchema?"]] = comp.Schema != ""
			newValues[NameToIndex["link"]] = entry.helmURL
			if updater.availableComponentsPerModel[entry.model] == nil {
				updater.availableComponentsPerModel[entry.model] = make(map[string]bool)
			}
			updater.availableComponentsPerModel[entry.model][comp.Kind] = true

			values = append(values, newValues)
			if len(values) >= 100 {
				values = updater.flushBatch(srv, values)
			}
		}
		if updater.availableModels[entry.model] != nil {
			log.Debug(fmt.Sprintf("Skipping spreadsheet updation for %s %s", entry.model))
			continue
		}
		newValues := make([]interface{}, len(resp.Values[0]))
		copy(newValues, resp.Values[0])
		newValues[NameToIndex["modelDisplayName"]] = entry.model
		newValues[NameToIndex["model"]] = entry.model
		newValues[NameToIndex["CRDs"]] = len(entry.comps)
		newValues[NameToIndex["link"]] = entry.helmURL
		updater.availableModels[entry.model] = newValues

		values = append(values, newValues)
		if len(values) >= 100 {
			values = updater.flushBatch(srv, values)
		}
	}
}

func (updater *spreadsheetUpdater) flushBatch(srv *sheets.Service, values [][]interface{}) [][]interface{} {
	appendRange := updater.sheetName + "!A4:AV4"
	row := &sheets.ValueRange{
		Values: values,
	}
	response, err := srv.Spreadsheets.Values.Append(spreadsheetID, appendRange, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
	values = make([][]interface{}, 0)
	if err != nil || response.HTTPStatusCode != 200 {
		log.Error(ErrorFailedUpdatingSheet(err))
	}
	return make([][]interface{}, 0)
}
