package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetID = "1BexOMvNdJ8zz9Ux15_w8NqMljhsF_gHC7UJ5lpPeTOw"
const sheetID = 0

func Spreadsheet(srv *sheets.Service, sheetName string, spreadsheet chan struct {
	comps []v1alpha1.ComponentDefinition
	model string
}, am map[string]bool, acpm map[string]map[string]bool) {
	start := time.Now()
	rangeString := sheetName + "!A4:AB4"
	for entry := range spreadsheet {
		if len(entry.comps) == 0 {
			continue
		}
		values := [][]interface{}{}
		for _, comp := range entry.comps {
			if acpm[entry.model][comp.Kind] {
				fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model, comp.Kind)
				continue
			}
			// Get the value of the specified cell.
			resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
			if err != nil {
				fmt.Println("Unable to retrieve data from sheet: ", err)
				return
			}
			newValues := resp.Values[0]
			newValues[6] = comp.Kind
			newValues[1] = entry.model
			newValues[0] = entry.model
			values = append(values, newValues)
			time.Sleep(500 * time.Millisecond)
		}
		if am[entry.model] {
			fmt.Println("[Debug][Spreadsheet] Skipping spreadsheet updation for ", entry.model)
			continue
		}
		// Get the value of the specified cell.
		resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeString).Do()
		if err != nil {
			fmt.Println("Unable to retrieve data from sheet: ", err)
			return
		}
		newValues := resp.Values[0]
		newValues[1] = entry.model
		newValues[0] = entry.model
		newValues[4] = len(entry.comps)
		values = append(values, newValues)
		row := &sheets.ValueRange{
			Values: values,
		}
		// // srv.Spreadsheets.Values.Get(spreadsheetID, sheetName)
		// // srv.Spreadsheets.Values.Update()
		response2, err := srv.Spreadsheets.Values.Append(spreadsheetID, sheetName, row).ValueInputOption("USER_ENTERED").InsertDataOption("INSERT_ROWS").Context(context.Background()).Do()
		if err != nil || response2.HTTPStatusCode != 200 {
			fmt.Println(err)
			continue
		}
		time.Sleep(500 * time.Millisecond)
	}
	elapsed := time.Now().Sub(start)
	fmt.Printf("Time taken by spreadsheet updater (including the time it required to generate components): %f", elapsed.Minutes())
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
