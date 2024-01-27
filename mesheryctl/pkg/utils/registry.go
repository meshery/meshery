package utils

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
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

func GenerateMDXStyleDocs(model ModelCSV, components []ComponentCSV, modelPath, imgPath string) error {
	modelName := FormatName(model.Model)
	// create dir for model
	modelsOutputPath, _ := filepath.Abs(filepath.Join("../", modelPath))
	modelDir := filepath.Join(modelsOutputPath, modelName)
	err := os.MkdirAll(modelDir, 0777)
	if err != nil {
		return err
	}

	// create img for model
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath))
	imgDir := filepath.Join(imgsOutputPath, modelName)
	err = os.MkdirAll(imgDir, 0777)
	if err != nil {
		return err
	}

	// create dir for color model icons
	iconsDir := filepath.Join(imgDir, "icons", "color")
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
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForMDXStyle(components, modelDir, compIconsSubDir)
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForMDXStyle(componentMetadata)
	err = WriteToFile(filepath.Join(modelDir, "index.mdx"), md)
	if err != nil {
		return err
	}

	return nil
}

func GenerateJSStyleDocs(model ModelCSV, docsJSON, imgPath string) (string, error) {
	formattedName := FormatName(model.Model)

	iconDir := filepath.Join(filepath.Join(strings.Split(imgPath, "/")[1:]...), formattedName) // "../images", "integrations"

	// generate data.js file
	jsonItem := model.CreateJSONItem(iconDir)
	docsJSON += jsonItem + ","

	// create color dir for icons
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath))
	colorIconsDir := filepath.Join(imgsOutputPath, "icons", "color")
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
	whiteIconsDir := filepath.Join(imgsOutputPath, "icons", "white")
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

	return docsJSON, nil
}

func GenerateMDStyleDocs(model ModelCSV, components []ComponentCSV, modelPath, imgPath string) error {

	modelName := FormatName(model.Model)

	// dir for markdown
	modelsOutputPath, _ := filepath.Abs(filepath.Join("../", modelPath))
	mdDir := filepath.Join(modelsOutputPath) // path, "pages", "integrations"
	err := os.MkdirAll(mdDir, 0777)
	if err != nil {
		return err
	}

	// dir for icons
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, modelName))
	iconsDir := filepath.Join(imgsOutputPath)
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
	_iconsSubDir := filepath.Join(filepath.Join(strings.Split(imgPath, "/")[1:]...), modelName, "components") // "assets", "img", "integrations"
	_imgOutputPath := filepath.Join(imgsOutputPath, "components")
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForMDStyle(components, _imgOutputPath, _iconsSubDir)
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForMDStyle(componentMetadata)
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

func NewSheetSRV(cred string) (*sheets.Service, error) {
	ctx := context.Background()
	byt, _ := base64.StdEncoding.DecodeString(cred)
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

func DownloadCSV(url string) (string, error) {
	file, err := os.CreateTemp("./", "*.csv")
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	byt, _ := io.ReadAll(resp.Body)
	file.WriteString(string(byt))
	byt, _ = io.ReadAll(file)
	path, _ := filepath.Abs(file.Name())
	return path, nil
}

func WriteSVG(path string, svg string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = file.WriteString(svg)
	if err != nil {
		return err
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		return err
	}
	return nil
}

func WriteToFile(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}

	_, err = file.WriteString(content)
	if err != nil {
		panic(err)
	}
	// Close the file to save the changes.
	err = file.Close()
	if err != nil {
		panic(err)
	}
	return nil
}
