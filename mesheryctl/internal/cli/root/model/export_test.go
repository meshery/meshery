package model

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

func TestExportModel(t *testing.T) {
	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	apiURL := "/api/meshmodels/export"

	// Create a custom directory for testing custom output location
	customDir := filepath.Join(currDir, "testCustomDir")
	os.MkdirAll(customDir, 0755)

	// Setup cleanup to remove exported files generated during testing
	t.Cleanup(func() {
		filesToClean := []string{
			"amd-gpu.tar",
			"amd-gpu.tar.gz",
			"amd-gpu.oci",
			"amd-gpu.yaml",
			"amd-gpu.json",
			"non-existent-model.tar",
			"non-existent-model.tar.gz",
		}

		for _, file := range filesToClean {
			os.Remove(file)
		}

		// Clean up custom directory
		os.RemoveAll(customDir)
	})

	// test scenarios for exporting models
	tests := []utils.MesheryListCommamdTest{
		{
			Name:             "Export model without arguments",
			Args:             []string{"export"},
			URL:              "",
			Fixture:          "empty.golden",
			ExpectedResponse: "export.model.no-args.output.golden",
			ExpectError:      true,
		},
		{
			Name:             "Export model with default parameters",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.default.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with YAML format",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.yaml.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with JSON format",
			Args:             []string{"export", "amd-gpu", "-t", "json", "-o", "oci"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=json&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.json.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with OCI output type",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.oci.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with TAR output type",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "tar"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=tar&components=true&relationships=true&pagesize=all", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.tar.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with version",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci", "--version", "v1.0.0"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=true&relationships=true&pagesize=all&version=v1.0.0", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.version.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model without components and relationships",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci", "--discard-components", "--discard-relationships"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=false&relationships=false&pagesize=all&version=v1.0.0", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.no-components-relationships.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with JSON format and TAR type",
			Args:             []string{"export", "amd-gpu", "-t", "json", "-o", "tar"},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=json&file_type=tar&components=false&relationships=false&pagesize=all&version=v1.0.0", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.json-tar.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export model with custom location",
			Args:             []string{"export", "amd-gpu", "-t", "yaml", "-o", "oci", "-l", customDir},
			URL:              fmt.Sprintf("%s?name=amd-gpu&output_format=yaml&file_type=oci&components=false&relationships=false&pagesize=all&version=v1.0.0", apiURL),
			Fixture:          "export.model.api.response.golden",
			ExpectedResponse: "export.model.custom-location.output.golden",
			ExpectError:      false,
		},
		{
			Name:             "Export non-existent model",
			Args:             []string{"export", "non-existent-model", "-t", "yaml", "-o", "oci"},
			URL:              fmt.Sprintf("%s?name=non-existent-model&output_format=yaml&file_type=oci&components=false&relationships=false&pagesize=all&version=v1.0.0", apiURL),
			Fixture:          "export.model.not-found.api.response.golden",
			ExpectedResponse: "export.model.not-found.output.golden",
			ExpectError:      false,
		},
	}

	utils.InvokeMesheryctlTestListCommand(t, update, ModelCmd, tests, currDir, "model")
}
