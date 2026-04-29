package model

import (
	"encoding/json"
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

func TestFlattenModelScaffold_MapsInitScaffoldToURLTemplate(t *testing.T) {
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

	out, err := flattenModelScaffold(input)
	require.NoError(t, err)

	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &got))

	assert.Equal(t, "test-model", got["model"])
	assert.Equal(t, "Test Model", got["modelDisplayName"])
	assert.Equal(t, "Uncategorized", got["category"])
	assert.Equal(t, "Utilities", got["subCategory"])
	assert.Equal(t, "github", got["registrant"])
	assert.Equal(t, "#00b39f", got["primaryColor"])
	assert.Equal(t, "#00D3A9", got["secondaryColor"])
	assert.Equal(t, "circle", got["shape"])
	assert.Equal(t, "<svg-color/>", got["svgColor"])
	assert.Equal(t, "<svg-white/>", got["svgWhite"])
	assert.Equal(t, true, got["isAnnotation"])
}

func TestFlattenModelScaffold_PreservesExistingURLTemplateKeys(t *testing.T) {
	input := []byte(`{
		"model": "cert-manager",
		"displayName": "Cert Manager",
		"registrant": "github",
		"category": "Security",
		"subCategory": "Certificates",
		"primaryColor": "#0B5FFF",
		"secondaryColor": "#121212",
		"isAnnotation": false,
		"publishToRegistry": false
	}`)

	out, err := flattenModelScaffold(input)
	require.NoError(t, err)

	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &got))

	assert.Equal(t, "cert-manager", got["model"])
	assert.Equal(t, "Cert Manager", got["modelDisplayName"])
	assert.Equal(t, "github", got["registrant"])
	assert.Equal(t, "Security", got["category"])
	assert.Equal(t, "Certificates", got["subCategory"])
	assert.Equal(t, "#0B5FFF", got["primaryColor"])
	assert.Equal(t, "#121212", got["secondaryColor"])
	assert.Equal(t, false, got["isAnnotation"])
	assert.Equal(t, false, got["publishToRegistry"])
}

// TestFlattenModelScaffold_PreservesLogoFromFixtureTemplate verifies that the logo
// field present in the real fixture template is not silently dropped — addresses
// miacycle's review comment on PR #18906.
func TestFlattenModelScaffold_PreservesLogoFromFixtureTemplate(t *testing.T) {
	// Mirrors the exact shape of mesheryctl/internal/cli/root/model/fixtures/templates/template.json
	input := []byte(`{
		"registrant": "github",
		"model": "cert-manager",
		"displayName": "Cert Manager",
		"category": "Security",
		"subCategory": "Certificates",
		"primaryColor": "#0B5FFF",
		"secondaryColor": "#121212",
		"logo": "https://raw.githubusercontent.com/cert-manager/cert-manager/master/logo/logo.png",
		"isAnnotation": false,
		"publishToRegistry": false
	}`)

	out, err := flattenModelScaffold(input)
	require.NoError(t, err)

	var got map[string]interface{}
	require.NoError(t, json.Unmarshal(out, &got))

	assert.Equal(t, "cert-manager", got["model"])
	assert.Equal(t, "github", got["registrant"])
	assert.Equal(t, "https://raw.githubusercontent.com/cert-manager/cert-manager/master/logo/logo.png", got["logo"],
		"logo field must be preserved and not silently dropped (miacycle review #18906)")
}
