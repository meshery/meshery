package utils

import (
	"os"
	"testing"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1beta1"
	"github.com/stretchr/testify/assert"
)

func TestCreateComponentDefinition(t *testing.T) {
	c := ComponentCSV{
		Registrant: "meshery",
		Component:  "TestComponent",
	}

	tests := []struct {
		isModelPublished bool
		defVersion       string
	}{
		{true, "v1.0.0"},
		{false, "v2.0.0"},
	}

	for _, tt := range tests {
		t.Run(tt.defVersion, func(t *testing.T) {
			compDef, err := c.CreateComponentDefinition(tt.isModelPublished, tt.defVersion)
			assert.NoError(t, err)
			assert.Equal(t, tt.defVersion, compDef.Version)
			assert.Equal(t, tt.isModelPublished, compDef.Metadata["published"])
			assert.Equal(t, "TestComponent", compDef.DisplayName)
		})
	}
}

func TestCreateComponentsMetadataAndCreateSVGsForMDXStyle(t *testing.T) {
	tests := []struct {
		name        string
		model       ModelCSV
		components  []ComponentCSV
		path        string
		svgDir      string
		expected    string
		expectError bool
	}{
		{
			name: "Basic test with one component",
			model: ModelCSV{
				SVGColor: "modelColorSVG",
				SVGWhite: "modelWhiteSVG",
			},
			components: []ComponentCSV{
				{
					Component:   "test-component",
					Description: "Test Description",
					SVGColor:    "componentColorSVG",
					SVGWhite:    "componentWhiteSVG",
				},
			},
			path:   "/tmp/test",
			svgDir: "svgs",
			expected: `[
{
"name": "test-component",
"colorIcon": "svgs/test-component/icons/color/test-component-color.svg",
"whiteIcon": "svgs/test-component/icons/white/test-component-white.svg",
"description": "Test Description",
}]`,
			expectError: false,
		},
		{
			name: "Error creating directories",
			model: ModelCSV{
				SVGColor: "modelColorSVG",
				SVGWhite: "modelWhiteSVG",
			},
			components: []ComponentCSV{
				{
					Component:   "test-component",
					Description: "Test Description",
					SVGColor:    "componentColorSVG",
					SVGWhite:    "componentWhiteSVG",
				},
			},
			path:        "/invalid/path",
			svgDir:      "svgs",
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test")
			assert.NoError(t, err)
			defer os.RemoveAll(tempDir)

			// Update path to use temporary directory
			if tt.path == "/tmp/test" {
				tt.path = tempDir
			}

			result, err := CreateComponentsMetadataAndCreateSVGsForMDXStyle(tt.model, tt.components, tt.path, tt.svgDir)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCreateComponentsMetadataAndCreateSVGsForMDStyle(t *testing.T) {
	tests := []struct {
		name       string
		model      ModelCSV
		components []ComponentCSV
		path       string
		svgDir     string
		expected   string
	}{
		{
			name: "Single component",
			model: ModelCSV{
				Registrant: "meshery",
				Model:      "TestModel",
				SVGColor:   "modelColorSVG",
				SVGWhite:   "modelWhiteSVG",
			},
			components: []ComponentCSV{
				{
					Registrant:  "meshery",
					Model:       "TestModel",
					Component:   "test-component",
					Description: "Test Description",
					SVGColor:    "componentColorSVG",
					SVGWhite:    "componentWhiteSVG",
				},
			},
			path:   "/tmp/testpath",
			svgDir: "svgs",
			expected: `
- name: test-component
  colorIcon: svgs/test-component/icons/color/test-component-color.svg
  whiteIcon: svgs/test-component/icons/white/test-component-white.svg
  description: Test Description`,
		},
		{
			name: "Component with missing SVGs",
			model: ModelCSV{
				Registrant: "meshery",
				Model:      "TestModel",
				SVGColor:   "modelColorSVG",
				SVGWhite:   "modelWhiteSVG",
			},
			components: []ComponentCSV{
				{
					Registrant:  "meshery",
					Model:       "TestModel",
					Component:   "test-component",
					Description: "Test Description",
					SVGColor:    "",
					SVGWhite:    "",
				},
			},
			path:   "/tmp/testpath",
			svgDir: "svgs",
			expected: `
- name: test-component
  colorIcon: svgs/test-component/icons/color/test-component-color.svg
  whiteIcon: svgs/test-component/icons/white/test-component-white.svg
  description: Test Description`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := os.MkdirAll(tt.path, 0777)
			assert.NoError(t, err)
			defer os.RemoveAll(tt.path)

			result, err := CreateComponentsMetadataAndCreateSVGsForMDStyle(tt.model, tt.components, tt.path, tt.svgDir)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)

		})
	}
}

func TestConvertCompDefToCompCSV(t *testing.T) {
	tests := []struct {
		name     string
		modelCSV ModelCSV
		compDef  v1beta1.ComponentDefinition
		expected ComponentCSV
	}{
		{
			name: "Basic conversion",
			modelCSV: ModelCSV{
				Registrant:       "meshery",
				Model:            "TestModel",
				ModelDisplayName: "Test Model Display Name",
				Category:         "TestCategory",
				SubCategory:      "TestSubCategory",
			},
			compDef: v1beta1.ComponentDefinition{
				Metadata: map[string]interface{}{
					"component": "TestComponent",
				},
				Component: v1beta1.ComponentEntity{
					TypeMeta: v1beta1.TypeMeta{
						Kind: "TestComponent",
					},
				},
			},
			expected: ComponentCSV{
				Registrant:       "meshery",
				Model:            "TestModel",
				Component:        "TestComponent",
				ModelDisplayName: "Test Model Display Name",
				Category:         "TestCategory",
				SubCategory:      "TestSubCategory",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertCompDefToCompCSV(&tt.modelCSV, tt.compDef)
			assert.Equal(t, tt.expected.Registrant, result.Registrant)
			assert.Equal(t, tt.expected.Model, result.Model)
			assert.Equal(t, tt.expected.Component, result.Component)
			assert.Equal(t, tt.expected.ModelDisplayName, result.ModelDisplayName)
			assert.Equal(t, tt.expected.Category, result.Category)
			assert.Equal(t, tt.expected.SubCategory, result.SubCategory)
		})
	}
}

func TestGetSVGForComponent(t *testing.T) {
	tests := []struct {
		name          string
		model         ModelCSV
		component     ComponentCSV
		expectedColor string
		expectedWhite string
	}{
		{
			name: "Component with both SVGs",
			model: ModelCSV{
				SVGColor: "modelColorSVG",
				SVGWhite: "modelWhiteSVG",
			},
			component: ComponentCSV{
				SVGColor: "componentColorSVG",
				SVGWhite: "componentWhiteSVG",
			},
			expectedColor: "componentColorSVG",
			expectedWhite: "componentWhiteSVG",
		},
		{
			name: "Component with only color SVG",
			model: ModelCSV{
				SVGColor: "modelColorSVG",
				SVGWhite: "modelWhiteSVG",
			},
			component: ComponentCSV{
				SVGColor: "componentColorSVG",
				SVGWhite: "",
			},
			expectedColor: "componentColorSVG",
			expectedWhite: "modelWhiteSVG",
		},
		{
			name: "Component with only white SVG",
			model: ModelCSV{
				SVGColor: "modelColorSVG",
				SVGWhite: "modelWhiteSVG",
			},
			component: ComponentCSV{
				SVGColor: "",
				SVGWhite: "componentWhiteSVG",
			},
			expectedColor: "modelColorSVG",
			expectedWhite: "componentWhiteSVG",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			colorSVG, whiteSVG := getSVGForComponent(tt.model, tt.component)
			assert.Equal(t, tt.expectedColor, colorSVG)
			assert.Equal(t, tt.expectedWhite, whiteSVG)
		})
	}
}
