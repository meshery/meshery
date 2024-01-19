package pkg

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type SystemType int

const (
	Meshery SystemType = iota
	Docs
	RemoteProvider
	rowIndex               = 1
	shouldRegisterColIndex = -1
)

func (dt SystemType) String() string {
	switch dt {
	case Meshery:
		return "meshery"

	case Docs:
		return "docs"

	case RemoteProvider:
		return "remote-provider"
	}
	return ""
}

func GetIndexForRegisterCol(cols []string, shouldRegister string) int {
	if shouldRegisterColIndex != -1 {
		return shouldRegisterColIndex
	}

	for index, col := range cols {
		if col == shouldRegister {
			return index
		}
	}
	return shouldRegisterColIndex
}

func FormatName(input string) string {
	formattedName := strings.ReplaceAll(input, " ", "-")
	formattedName = strings.ToLower(formattedName)
	return formattedName
}

func CreateFiles(path, filetype, name, content string) {
	formattedName := FormatName(name)

	fullPath := path + "/" + formattedName + filetype

	file, err := os.Create(fullPath)
	fmt.Println("Creating file:", fullPath)
	if err != nil {
		log.Println("Error creating filetype file:", err)
		return
	}
	defer file.Close()

	_, err = io.WriteString(file, content)
	if err != nil {
		log.Println("Error writing to filetype file:", err)
	}
}

func GenerateLayer5Docs(model ModelCSV, components []ComponentCSV, path string) error {
	modelName := FormatName(model.Model)
	// create dir for model
	modelDir := filepath.Join(path, modelName)
	err := os.MkdirAll(modelDir, 0777)
	if err != nil {
		return err
	}

	// create dir for color model icons
	iconsDir := filepath.Join(modelDir, "icons", "color")
	err = os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	err = WriteSVG(filepath.Join(iconsDir, modelName+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}

	// create dir for white model icons
	iconsDir = filepath.Join(modelDir, "icons", "white")
	err = os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	err = WriteSVG(filepath.Join(iconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// generate components metadata and create svg files
	compIconsSubDir := filepath.Join("icons", "components")
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForLayer5io(components, modelDir, compIconsSubDir)
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForLayer5(componentMetadata)
	err = WriteToFile(filepath.Join(modelDir, "index.mdx"), md)
	if err != nil {
		return err
	}

	return nil
}

func GenerateMesheryioDocs(model ModelCSV, path, mesheryioDocsJSON string) (string, error) {

	formattedName := FormatName(model.Model)

	iconDir := filepath.Join("../images", "integrations", formattedName)

	// generate data.js file
	jsonItem := model.CreateJSONItem(iconDir)
	mesheryioDocsJSON += jsonItem + ","

	// create color dir for icons
	colorIconsDir := filepath.Join(path, iconDir, "icons", "color")
	// create svg dir
	err := os.MkdirAll(colorIconsDir, 0777)
	if err != nil {
		return "", err
	}

	// write color svg
	err = WriteSVG(filepath.Join(colorIconsDir, formattedName+"-color.svg"), model.SVGColor)
	if err != nil {
		return "", err
	}

	// create white dir for icons
	whiteIconsDir := filepath.Join(path, iconDir, "icons", "white")
	// create svg dir
	err = os.MkdirAll(whiteIconsDir, 0777)
	if err != nil {
		return "", err
	}

	// write white svg
	err = WriteSVG(filepath.Join(whiteIconsDir, formattedName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return "", err
	}

	return mesheryioDocsJSON, nil
}

func GenerateMesheryDocs(model ModelCSV, components []ComponentCSV, path string) error {

	modelName := FormatName(model.Model)

	// dir for markdown
	mdDir := filepath.Join(path, "pages", "integrations")
	err := os.MkdirAll(mdDir, 0777)
	if err != nil {
		return err
	}

	// dir for icons
	_iconsSubDir := filepath.Join("assets", "img", "integrations", modelName)
	iconsDir := filepath.Join(path, _iconsSubDir)
	err = os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	// for color icons
	colorIconsDir := filepath.Join(iconsDir, "icons", "color")
	err = os.MkdirAll(colorIconsDir, 0777)
	if err != nil {
		return err
	}

	err = WriteSVG(filepath.Join(colorIconsDir, modelName+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}

	// for white icons
	whiteIconsDir := filepath.Join(iconsDir, "icons", "white")
	err = os.MkdirAll(whiteIconsDir, 0777)
	if err != nil {
		return err
	}

	err = WriteSVG(filepath.Join(whiteIconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// generate components metadata and create svg files
	compIconsSubDir := filepath.Join(_iconsSubDir, "components")
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForMesheryDocs(components, path, compIconsSubDir)
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForMesheryDocs(componentMetadata)
	file, err := os.Create(filepath.Join(mdDir, modelName+".md"))
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, md)
	if err != nil {
		return err
	}

	return nil
}

func NewSheetSRV() (*sheets.Service, error) {
	ctx := context.Background()
	byt, _ := base64.StdEncoding.DecodeString(os.Getenv("CRED"))
	// authenticate and get configuration
	config, err := google.JWTConfigFromJSON(byt, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, err
	}
	// create client with config and context
	client := config.Client(ctx)
	// create new service using client
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}
	return srv, nil
}
