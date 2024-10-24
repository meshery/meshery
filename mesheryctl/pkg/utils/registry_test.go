package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	tests := []struct {
		input    SystemType
		expected string
	}{
		{Meshery, "meshery"},
		{Docs, "docs"},
		{RemoteProvider, "remote-provider"},
		{SystemType(999), ""}, // Test for an undefined SystemType
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetIndexForRegisterCol(t *testing.T) {
	tests := []struct {
		cols           []string
		shouldRegister string
		expected       int
	}{
		{[]string{"name", "age", "shouldRegister"}, "shouldRegister", 2},
		{[]string{"name", "shouldRegister", "age"}, "shouldRegister", 1},
		{[]string{"shouldRegister", "name", "age"}, "shouldRegister", 0},
		{[]string{"name", "age", "address"}, "shouldRegister", -1},
		{[]string{}, "shouldRegister", -1},
	}

	for _, tt := range tests {
		t.Run(tt.shouldRegister, func(t *testing.T) {
			result := GetIndexForRegisterCol(tt.cols, tt.shouldRegister)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGenerateMDXStyleDocs(t *testing.T) {
	model := ModelCSV{
		Model:    "TestModel",
		SVGColor: "<svg>color</svg>",
		SVGWhite: "<svg>white</svg>",
	}
	components := []ComponentCSV{
		{Component: "Component1"},
		{Component: "Component2"},
	}

	modelPath := "modelPath"
	imgPath := "imgPath"

	err := GenerateMDXStyleDocs(model, components, modelPath, imgPath)
	assert.NoError(t, err)

	// Cleanup the directories
	modelDir, _ := filepath.Abs(filepath.Join("../", modelPath, "testmodel"))
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath))
	imgDir := filepath.Join(imgsOutputPath, "testmodel")
	defer os.RemoveAll(modelDir)
	defer os.RemoveAll(imgDir)

}
func TestGenerateJSStyleDocs(t *testing.T) {
	model := ModelCSV{
		Model:    "TestModel",
		SVGColor: "<svg>color</svg>",
		SVGWhite: "<svg>white</svg>",
	}
	docsJSON := ""
	imgPath := "assets/img/integrations"
	expected := `{"name":"TestModel","color":"img/integrations/testmodel/icons/color/testmodel-color.svg","white":"img/integrations/testmodel/icons/white/testmodel-white.svg","permalink":""},`

	result, err := GenerateJSStyleDocs(model, docsJSON, imgPath)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Cleanup the directories
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, "testmodel"))
	colorIconsDir := filepath.Join(imgsOutputPath, "icons", "color")
	whiteIconsDir := filepath.Join(imgsOutputPath, "icons", "white")
	defer os.RemoveAll(colorIconsDir)
	defer os.RemoveAll(whiteIconsDir)

}

func TestGenerateMDStyleDocs(t *testing.T) {
	model := ModelCSV{
		Model:    "TestModel",
		SVGColor: "<svg>color</svg>",
		SVGWhite: "<svg>white</svg>",
	}
	components := []ComponentCSV{
		{Component: "Component1"},
		{Component: "Component2"},
	}
	modelPath := "modelPath"
	imgPath := "imgPath"
	modelName := "testmodel"

	err := GenerateMDStyleDocs(model, components, modelPath, imgPath)
	assert.NoError(t, err)

	// Cleanup the directories
	modelsOutputPath, _ := filepath.Abs(filepath.Join("../", modelPath))
	mdDir := filepath.Join(modelsOutputPath)
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, modelName))
	iconsDir := filepath.Join(imgsOutputPath)
	defer os.RemoveAll(mdDir)
	defer os.RemoveAll(iconsDir)

}

func TestGenerateIcons(t *testing.T) {
	model := ModelCSV{
		Model:    "TestModel",
		SVGColor: "<svg>color</svg>",
		SVGWhite: "<svg>white</svg>",
	}
	components := []ComponentCSV{
		{Component: "Component1"},
		{Component: "Component2"},
	}
	imgPath := "assets/img/integrations"
	modelName := "testmodel"

	err := GenerateIcons(model, components, imgPath)
	assert.NoError(t, err)

	// Cleanup the directories
	imgsOutputPath, _ := filepath.Abs(filepath.Join("../", imgPath, modelName))
	iconsDir := filepath.Join(imgsOutputPath)

	defer os.RemoveAll(iconsDir)

}
