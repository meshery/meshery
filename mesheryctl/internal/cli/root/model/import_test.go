package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasCSVs(t *testing.T) {
	tests := []struct {
		name           string
		dirPath        string
		expectedResult bool
	}{
		{
			name:           "directory with CSVs",
			dirPath:        "./fixtures/with_csvs",
			expectedResult: true,
		},
		{
			name:           "directory without CSVs",
			dirPath:        "./fixtures/without_csvs",
			expectedResult: false,
		},
		{
			name:           "non-existent directory",
			dirPath:        "./fixtures/invalid_path",
			expectedResult: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := hasCSVs(tc.dirPath)
			assert.Equal(t, tc.expectedResult, res)
		})
	}
}

// TestScaffoldToModel_MapsInitScaffoldToURLTemplate tests Path 1: canonical ModelDefinition format (from "model init").
func TestScaffoldToModel_MapsInitScaffoldToURLTemplate(t *testing.T) {
	input := []byte(`{
		"name": "test-model",
		"displayName": "Test Model",
		"category": {"name": "Uncategorized"},
		"subCategory": "Utilities",
		"registrant": {"kind": "github", "name": ""},
		"metadata": {
			"primaryColor": "#00b39f",
			"secondaryColor": "#00D3A9",
			"shape": "circle",
			"svgColor": "<svg-color/>",
			"svgWhite": "<svg-white/>",
			"isAnnotation": true
		}
	}`)

	got, err := scaffoldToModel(input)
	require.NoError(t, err)

	assert.Equal(t, "test-model", got.Model)
	assert.Equal(t, "Test Model", got.ModelDisplayName)
	assert.Equal(t, "github", got.Registrant)
	assert.Equal(t, "#00b39f", got.PrimaryColor)
	assert.Equal(t, "#00D3A9", got.SecondaryColor)
	assert.Equal(t, "<svg-color/>", got.SvgColor)
	assert.Equal(t, "<svg-white/>", got.SvgWhite)
	assert.Equal(t, true, got.IsAnnotation)
	assert.Equal(t, "circle", got.Shape, "shape from metadata must not be silently dropped")
}

// TestScaffoldToModel_PreservesExistingURLTemplateKeys tests Path 2: flat template with "modelDisplayName" key.
func TestScaffoldToModel_PreservesExistingURLTemplateKeys(t *testing.T) {
	input := []byte(`{
		"model": "cert-manager",
		"modelDisplayName": "Cert Manager",
		"registrant": "github",
		"category": "Security",
		"subCategory": "Certificates",
		"primaryColor": "#0B5FFF",
		"secondaryColor": "#121212",
		"isAnnotation": false,
		"publishToRegistry": false
	}`)

	got, err := scaffoldToModel(input)
	require.NoError(t, err)

	assert.Equal(t, "cert-manager", got.Model)
	assert.Equal(t, "Cert Manager", got.ModelDisplayName)
	assert.Equal(t, "github", got.Registrant)
	assert.Equal(t, "#0B5FFF", got.PrimaryColor)
	assert.Equal(t, "#121212", got.SecondaryColor)
	assert.Equal(t, false, got.IsAnnotation)
	assert.Equal(t, false, got.PublishToRegistry)
}

// TestScaffoldToModel_ResolvesDisplayNameAlias tests Path 2 alias: flat template using "displayName" maps to ModelDisplayName.
func TestScaffoldToModel_ResolvesDisplayNameAlias(t *testing.T) {
	input := []byte(`{
		"registrant": "github",
		"model": "cert-manager",
		"displayName": "Cert Manager",
		"category": "Security",
		"subCategory": "Certificates",
		"primaryColor": "#0B5FFF",
		"secondaryColor": "#121212",
		"isAnnotation": false,
		"publishToRegistry": false
	}`)

	got, err := scaffoldToModel(input)
	require.NoError(t, err)

	assert.Equal(t, "cert-manager", got.Model)
	assert.Equal(t, "github", got.Registrant)
	assert.Equal(t, "Cert Manager", got.ModelDisplayName,
		"displayName from flat template must not be silently dropped (json tag alias)")
}

// TestScaffoldToModel_ParsesFlatFixtureTemplate tests Path 2: flat fixture template parses correctly.
func TestScaffoldToModel_ParsesFlatFixtureTemplate(t *testing.T) {
	input := []byte(`{
		"registrant": "github",
		"model": "cert-manager",
		"modelDisplayName": "Cert Manager",
		"category": "Security",
		"primaryColor": "#0B5FFF",
		"isAnnotation": false,
		"publishToRegistry": false
	}`)

	got, err := scaffoldToModel(input)
	require.NoError(t, err)

	assert.Equal(t, "cert-manager", got.Model)
	assert.Equal(t, "github", got.Registrant)
	assert.Equal(t, "#0B5FFF", got.PrimaryColor)
	assert.Equal(t, false, got.IsAnnotation)
}
