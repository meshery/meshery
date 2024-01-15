package pkg

import (
	"fmt"
	"strings"

	"github.com/layer5io/meshkit/utils/csv"
)

var (
	rowIndex = 1
	shouldRegisterMod        = "Publish?"
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
	Feature1 string `json:"feature1"`
	Feature2 string `json:"feature2"`
	Feature3 string `json:"feature3"`
	HowItWorks string `json:"howItWorks"`
	HowItWorksDetails string `json:"howItWorksDetails"`
	StandardBlurb string `json:"standardBlurb"`
	Screenshots string `json:"screenshots"`
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
		index := GetIndexForRegisterCol(columns, shouldRegisterMod)
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
			// fmt.Println("data: ", data)
		case err := <-errorChan:
			fmt.Println(err)

		case <-csvReader.Context.Done():
			return
		}
	}
}


func (m ModelCSV) CreateMarkDown() string {
	markdown := m.FullPage
	markdown = strings.ReplaceAll(markdown, "\r", "\n")
	return markdown
}

// Creates JSON formatted meshmodel attribute item for Meshery docs
func (m ModelCSV) CreateJSONItem() string {
	json := "{"
	json += fmt.Sprintf("\"name\":\"%s\"", m.Model)
	// If SVGs exist, then add the paths to json
	if m.SVGColor != "" {
		json += fmt.Sprintf(",\"color\":\"../images/integration/%s-color.svg\"", m.Model)
	}

	if m.SVGWhite != "" {
		json += fmt.Sprintf(",\"white\":\"../images/integration/%s-white.svg\"", m.Model)
	}

	json += fmt.Sprintf(",\"permalink\":\"https://docs.meshery.io/extensibility/integrations/%s\"", FormatName(m.ModelDisplayName))

	json += "}"
	return json
}