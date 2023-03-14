package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetID = "1BexOMvNdJ8zz9Ux15_w8NqMljhsF_gHC7UJ5lpPeTOw"

func main() {
	ctx := context.Background()

	// get bytes from base64 encoded google service accounts key
	f, _ := os.Open("./credentials.json")
	byt, _ := io.ReadAll(f)
	credBytes, err := base64.StdEncoding.DecodeString(string(byt))
	if err != nil {
		fmt.Println("ERR1", err)
		return
	}

	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(credBytes, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		fmt.Println("ERR2", err)
		return
	}

	// create client with config and context
	client := config.Client(ctx)

	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		fmt.Println("ERR3", err)
		return
	}
	rangeToUpdate := "Sheet1!A2" // The cell to update.
	valueInputOption := "USER_ENTERED"
	values := [][]interface{}{{"New Value"}}
	data := &sheets.ValueRange{
		Values: values,
	}
	_, err = srv.Spreadsheets.Values.Update(spreadsheetID, rangeToUpdate, data).ValueInputOption(valueInputOption).Do()
	if err != nil {
		log.Fatalf("Unable to update data: %v", err)
	}
}
