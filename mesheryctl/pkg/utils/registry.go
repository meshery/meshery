package utils

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshkit/utils"
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

func GenerateMDXStyleDocs(model ModelCSV, components []ComponentCSV, modelPath, imgPath string) error {
	modelName := utils.FormatName(model.Model)
	// create dir for model
	modelDir, _ := filepath.Abs(filepath.Join("../", modelPath, modelName))
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

	err = utils.WriteToFile(filepath.Join(iconsDir, modelName+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}

	// create dir for white model icons
	iconsDir = filepath.Join(modelDir, "icons", "white")
	err = os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	err = utils.WriteToFile(filepath.Join(iconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// generate components metadata and create svg files
	compIconsSubDir := filepath.Join("icons", "components")
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForMDXStyle(model, components, modelDir, compIconsSubDir)
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForMDXStyle(componentMetadata)
	err = utils.WriteToFile(filepath.Join(modelDir, "index.mdx"), md)
	if err != nil {
		return err
	}

	return nil
}

func GenerateJSStyleDocs(model ModelCSV, docsJSON, imgPath string) (string, error) {
	modelName := utils.FormatName(model.Model)

	iconDir := filepath.Join(filepath.Join(strings.Split(imgPath, "/")[1:]...), modelName) // "../images", "integrations"

	// generate data.js file
	jsonItem := model.CreateJSONItem(iconDir)
	docsJSON += jsonItem + ","

	// create color dir for icons
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, modelName))
	colorIconsDir := filepath.Join(imgsOutputPath, "icons", "color")
	// create svg dir
	err := os.MkdirAll(colorIconsDir, 0777)
	if err != nil {
		return "", err
	}

	// write color svg
	err = utils.WriteToFile(filepath.Join(colorIconsDir, modelName+"-color.svg"), model.SVGColor)
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
	err = utils.WriteToFile(filepath.Join(whiteIconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return "", err
	}

	return docsJSON, nil
}

func GenerateMDStyleDocs(model ModelCSV, components []ComponentCSV, modelPath, imgPath string) error {

	modelName := utils.FormatName(model.Model)

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

	err = utils.WriteToFile(filepath.Join(colorIconsDir, modelName+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}

	// for white icons
	whiteIconsDir := filepath.Join(iconsDir, "icons", "white")
	err = os.MkdirAll(whiteIconsDir, 0777)
	if err != nil {
		return err
	}

	err = utils.WriteToFile(filepath.Join(whiteIconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// generate components metadata and create svg files
	_iconsSubDir := filepath.Join(filepath.Join(strings.Split(imgPath, "/")[1:]...), modelName, "components") // "assets", "img", "integrations"
	_imgOutputPath := filepath.Join(imgsOutputPath, "components")
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGsForMDStyle(model, components, _imgOutputPath, _iconsSubDir)
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

func GenerateIcons(model ModelCSV, components []ComponentCSV, imgPath string) error {
	modelName := utils.FormatName(model.Model)

	// Dir for icons
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, modelName))
	iconsDir := filepath.Join(imgsOutputPath)
	err := os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	// For color icons
	colorIconsDir := filepath.Join(iconsDir, "icons", "color")
	err = os.MkdirAll(colorIconsDir, 0777)
	if err != nil {
		return err
	}

	err = utils.WriteToFile(filepath.Join(colorIconsDir, modelName+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}

	// For white icons
	whiteIconsDir := filepath.Join(iconsDir, "icons", "white")
	err = os.MkdirAll(whiteIconsDir, 0777)
	if err != nil {
		return err
	}

	err = utils.WriteToFile(filepath.Join(whiteIconsDir, modelName+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// Generate components metadata and create SVG files
	_iconsSubDir := filepath.Join(filepath.Join(strings.Split(imgPath, "/")[1:]...), modelName, "components")
	_imgOutputPath := filepath.Join(imgsOutputPath, "components")
	_, err = CreateComponentsMetadataAndCreateSVGsForMDStyle(model, components, _imgOutputPath, _iconsSubDir)
	if err != nil {
		return err
	}

	return nil
}
