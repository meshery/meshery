package pkg

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/utils/csv"
)

var (
	rowIndex = 1
	shouldRegister         = "Publish?"
	shouldRegisterColIndex = -1
	sheetID = "1551563103" // update
)

type ModelCSV struct {
	ModelDisplayName string `json:"modelDisplayName"`
	Model string `json:"model"`
	Registrant string `json:"registrant"`
	Category string `json:"category"`
	SubCategory string `json:"subCategory"`
	Source string `json:"source"`
	Website string `json:"website"`
	Docs string `json:"docs"`
	Shape string `json:"shape"`
	PrimaryColor string `json:"primaryColor"`
	SecondaryColor string `json:"secondaryColor"`
	StyleOverrides string `json:"styleOverrides"`
	Styles string `json:"styles"`
	ShapePolygonPoints string `json:"shapePolygonPoints"`
	DefaultData string `json:"defaultData"`
	Capabilities string `json:"capabilities"`
	LogoURL string `json:"logoURL"`
	SVGColor string `json:"svgColor"`
	SVGWhite string `json:"svgWhite"`
	SVGComplete string `json:"svgComplete"`
	PublishToRegistry string `json:"publishToRegistry"`
	ModelAnnotation string `json:"modelAnnotation"`
	AboutProject string `json:"aboutProject"`
	PageSubtTitle string `json:"pageSubtitle"`
	FeatureList string `json:"featureList"`
	HowItWorks string `json:"howItWorks"`
	HowItWorksDetails string `json:"howItWorksDetails"`
	StandardBlurb string `json:"standardBlurb"`
	WorkingSlides string `json:"workingSlides"`
	Published string `json:"published"`
	IntegrationIcon string `json:"integrationIcon"`
	DarkModeIntegrationIcon string `json:"darkModeIntegrationIcon"`
	FullPage string `json:"fullPage"`
}

type ModelCSVHelper struct {
	SheetID string
	CSVPath string
	Models []ModelCSV
}

func NewModelCSVHelper(sheetURL string) (*ModelCSVHelper, error) {
	sheetURL = sheetURL + "&gid=" + sheetID
	csvPath, err := DownloadCSV(sheetURL)
	if err != nil {
		return nil, err
	}

	return &ModelCSVHelper{
		SheetID: sheetID,
		CSVPath: csvPath,
		Models: []ModelCSV{},
	}, nil
}

func (mch *ModelCSVHelper) ParseModelsSheet(){
	ch := make(chan ModelCSV, 1)
	errorChan := make(chan error, 1)
	csvReader, err := csv.NewCSVParser[ModelCSV](mch.CSVPath, rowIndex, nil, func(columns []string, currentRow []string) bool {
		index := GetIndexForRegisterCol(columns)
		if index != -1 && index < len(currentRow) {
			shouldRegister := currentRow[index]
			return strings.ToLower(shouldRegister) == "true"
		}
		return false
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
			mch.Models = append(mch.Models, data)
		case err := <-errorChan:
			fmt.Println(err)

		case <-csvReader.Context.Done():
			return
		}
	}
}