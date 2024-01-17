package pkg

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SystemType int

const (
	Meshery SystemType = iota
	Docs
	RemoteProvider
	rowIndex = 1
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
	// create dir for model
	modelDir := filepath.Join(path, model.Model)
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

	err = WriteSVG(filepath.Join(iconsDir, model.Model+"-color.svg"), model.SVGColor)
	if err != nil {
		return err
	}


	// create dir for white model icons
	iconsDir = filepath.Join(modelDir, "icons", "white")
	err = os.MkdirAll(iconsDir, 0777)
	if err != nil {
		return err
	}

	err = WriteSVG(filepath.Join(iconsDir, model.Model+"-white.svg"), model.SVGWhite)
	if err != nil {
		return err
	}

	// generate components metadata and create svg files

	componentMetadata, err := CreateComponentsMetadataAndCreateSVGs(components, modelDir, "components") 
	if err != nil {
		return err
	}


	// generate markdown file
	md := model.CreateMarkDownForLayer5(componentMetadata)
		_ = WriteToFile(filepath.Join(modelDir, "index.mdx"), md)

	return nil
}

func GenerateMesheryioDocs(model ModelCSV, path, mesheryioDocsJSON string) (string, error) {
	// generate data.js file
	jsonItem := model.CreateJSONItem()
	mesheryioDocsJSON += jsonItem + ","


	// create svg dir
	err := os.MkdirAll(filepath.Join(path, "../", "images"), 0777)
		if err != nil {
			return "", err
		}

		// write color svg
		err = WriteSVG(filepath.Join(path, "../", "images", model.Model+"-color.svg"), model.SVGColor)
		if err != nil {
			return "", err
		}

		// write white svg
		err = WriteSVG(filepath.Join(path, "../", "images", model.Model+"-white.svg"), model.SVGWhite)
		if err != nil {
				return "", err
		}


	return mesheryioDocsJSON, nil
}

func GenerateMesheryDocs(model ModelCSV, components []ComponentCSV, path string) error {

	modelName := FormatName(model.Model)
	
	// dir for markdown
	mdDir := filepath.Join(path, "pages", "integrations")

	// dir for icons
	iconsDir := filepath.Join(mdDir, "assets", "img", "integrations", modelName)
	err := os.MkdirAll(iconsDir, 0777)
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
	componentMetadata, err := CreateComponentsMetadataAndCreateSVGs(components, iconsDir, "components")
	if err != nil {
		return err
	}

	// generate markdown file
	md := model.CreateMarkDownForMesheryDocs(componentMetadata)
	_ = WriteToFile(filepath.Join(mdDir, modelName+".mdx"), md)


	return nil
}
