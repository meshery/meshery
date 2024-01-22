package pkg

import (
	"fmt"
	"os"
	"strconv"

	"github.com/layer5io/meshkit/utils/csv"
)

type ComponentCSV struct {
	Registrant         string `json:"registrant"`
	Model              string `json:"model"`
	Component          string `json:"component"`
	Shape              string `json:"shape"`
	PrimaryColor       string `json:"primaryColor"`
	SecondaryColor     string `json:"secondaryColor"`
	SVGColor           string `json:"svgColor"`
	SVGWhite           string `json:"svgWhite"`
	SVGComplete        string `json:"svgComplete"`
	HasSchema          string `json:"hasSchema"`
	Description        string `json:"description"`
	Docs               string `json:"docs"`
	StyleOverrides     string `json:"styleOverrides"`
	Styles             string `json:"styles"`
	ShapePolygonPoints string `json:"shapePolygonPoints"`
	DefaultData        string `json:"defaultData"`
	Capabilities       string `json:"capabilities"`
	LogoURL            string `json:"logoURL"`
	Genealogy          string `json:"genealogy"`
	IsAnnotation       string `json:"isAnnotation"`
}
type ComponentCSVHelper struct {
	SpreadsheetID  int64
	SpreadsheetURL string
	Title          string
	CSVPath        string
	Components     map[string]map[string][]ComponentCSV
}

func NewComponentCSVHelper(sheetURL, spreadsheetName string, spreadsheetID int64) (*ComponentCSVHelper, error) {
	sheetURL = sheetURL + "/pub?output=csv" + "&gid=" + strconv.FormatInt(spreadsheetID, 10)
	fmt.Println("Downloading CSV from:", sheetURL)
	csvPath, err := DownloadCSV(sheetURL)
	if err != nil {
		return nil, err
	}

	return &ComponentCSVHelper{
		SpreadsheetID:  spreadsheetID,
		SpreadsheetURL: sheetURL,
		Title:          spreadsheetName,
		CSVPath:        csvPath,
		Components:     make(map[string]map[string][]ComponentCSV),
	}, nil
}

func (mch *ComponentCSVHelper) ParseComponentsSheet() {
	ch := make(chan ComponentCSV, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[ComponentCSV](mch.CSVPath, rowIndex, nil, func(_ []string, _ []string) bool {
		return true
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		err := csvReader.Parse(ch, errorChan)
		if err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {

		case data := <-ch:
			if mch.Components[data.Registrant] == nil {
				mch.Components[data.Registrant] = make(map[string][]ComponentCSV, 0)
			}
			if mch.Components[data.Registrant][data.Model] == nil {
				mch.Components[data.Registrant][data.Model] = make([]ComponentCSV, 0)
			}
			mch.Components[data.Registrant][data.Model] = append(mch.Components[data.Registrant][data.Model], data)
			fmt.Printf("Reading Component: %s for Modal: %s from Registrant: %s\n", data.Component, data.Model, data.Registrant)
		case err := <-errorChan:
			fmt.Println(err)

		case <-csvReader.Context.Done():
			return
		}
	}
}

func (m ComponentCSVHelper) Cleanup() error {
	// remove csv file
	err := os.Remove(m.CSVPath)
	if err != nil {
		return err
	}
	return nil
}
