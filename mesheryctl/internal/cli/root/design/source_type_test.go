package design

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/schemas/models/core"
	"github.com/stretchr/testify/assert"
)

func Test_ValidateAndNormalizeSourceType(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    core.IaCFileTypes
		expectError bool
	}{
		// Canonical schema values
		{"canonical k8s-manifest", "k8s-manifest", core.K8sManifest, false},
		{"canonical helm-chart", "helm-chart", core.HelmChart, false},
		{"canonical docker-compose", "docker-compose", core.DockerCompose, false},
		{"canonical k8s-kustomize", "k8s-kustomize", core.K8sKustomize, false},
		{"canonical meshery-design", "meshery-design", core.MesheryDesign, false},

		// Shorthand aliases
		{"shorthand manifest", "manifest", core.K8sManifest, false},
		{"shorthand helm", "helm", core.HelmChart, false},
		{"shorthand compose", "compose", core.DockerCompose, false},
		{"shorthand kustomize", "kustomize", core.K8sKustomize, false},
		{"shorthand design", "design", core.MesheryDesign, false},
		{"shorthand docker", "docker", core.DockerCompose, false},
		{"shorthand k8s", "k8s", core.K8sManifest, false},

		// Human-readable titles & case-insensitivity
		{"title Kubernetes Manifest", "Kubernetes Manifest", core.K8sManifest, false},
		{"lowercase kubernetes manifest", "kubernetes manifest", core.K8sManifest, false},
		{"uppercase HELM", "HELM", core.HelmChart, false},
		{"title Docker Compose", "Docker Compose", core.DockerCompose, false},
		{"title Helm Chart", "Helm Chart", core.HelmChart, false},
		{"title K8s Kustomize", "K8s Kustomize", core.K8sKustomize, false},
		{"title Meshery Design", "Meshery Design", core.MesheryDesign, false},
		{"with whitespace", "  manifest  ", core.K8sManifest, false},

		// Invalid values
		{"invalid type random", "invalid-type", "", true},
		{"invalid type empty", "", "", true},
		{"invalid type typo", "manifst", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateAndNormalizeSourceType(tt.input)
			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, core.IaCFileTypes(""), got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func Test_retrieveProvidedSourceType_ShorthandsAndCanonical(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
	}{
		{"shorthand manifest", "manifest", string(core.K8sManifest), false},
		{"shorthand compose", "compose", string(core.DockerCompose), false},
		{"shorthand helm", "helm", string(core.HelmChart), false},
		{"shorthand kustomize", "kustomize", string(core.K8sKustomize), false},
		{"shorthand design", "design", string(core.MesheryDesign), false},
		{"canonical k8s-manifest", "k8s-manifest", string(core.K8sManifest), false},
		{"title Kubernetes Manifest", "Kubernetes Manifest", string(core.K8sManifest), false},
		{"invalid source type", "invalid source type", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := retrieveProvidedSourceType(tt.input, GetDesignSourceTypes())
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func Test_SourceTypeFlagHelp(t *testing.T) {
	expectedDesc := "Type of source file (ex. manifest / compose / helm / kustomize / design)"

	deployFlag := deployDesignCmd.Flags().Lookup("source-type")
	assert.NotNil(t, deployFlag)
	assert.Equal(t, expectedDesc, deployFlag.Usage)

	importFlag := importCmd.Flags().Lookup("source-type")
	assert.NotNil(t, importFlag)
	assert.Equal(t, expectedDesc, importFlag.Usage)
}

func Test_importPattern_WithSourceTypeFlag(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/pattern/import", func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(200, `[{"id":"3817ec9a-1d83-4f6f-9154-0fd4408ba9f0","name":"SampleApp"}]`), nil
	})

	// Test with shorthand -s manifest
	DesignCmd.SetArgs([]string{"import", "-f", filepath.Join("fixtures", "sampleDesign.golden"), "-s", "manifest"})
	assert.NoError(t, DesignCmd.Execute())
	assert.Equal(t, string(core.K8sManifest), sourceType)

	// Test with canonical -s k8s-manifest
	DesignCmd.SetArgs([]string{"import", "-f", filepath.Join("fixtures", "sampleDesign.golden"), "-s", "k8s-manifest"})
	assert.NoError(t, DesignCmd.Execute())
	assert.Equal(t, string(core.K8sManifest), sourceType)

	// Test with legacy -s "Kubernetes Manifest"
	DesignCmd.SetArgs([]string{"import", "-f", filepath.Join("fixtures", "sampleDesign.golden"), "-s", "Kubernetes Manifest"})
	assert.NoError(t, DesignCmd.Execute())
	assert.Equal(t, string(core.K8sManifest), sourceType)

	// Test with invalid source type
	DesignCmd.SetArgs([]string{"import", "-f", filepath.Join("fixtures", "sampleDesign.golden"), "-s", "invalid-source"})
	assert.Error(t, DesignCmd.Execute())
}

func Test_deployDesign_NormalizesSourceTypeInPlace(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)
	utils.TokenFlag = utils.GetToken(t)

	patternContent := "name: SampleApp\nservices:\n  myapp:\n    type: Application\n"
	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/pattern/import", func(req *http.Request) (*http.Response, error) {
		resp := fmt.Sprintf(`[{"id":"3817ec9a-1d83-4f6f-9154-0fd4408ba9f0","name":"SampleApp","patternFile":%q,"pattern_file":%q}]`, patternContent, patternContent)
		return httpmock.NewStringResponse(200, resp), nil
	})
	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/pattern/deploy", func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		var payload models.MesheryPatternFileDeployPayload
		err = json.Unmarshal(body, &payload)
		assert.NoError(t, err)
		assert.Contains(t, payload.PatternFile, "SampleApp")
		return httpmock.NewStringResponse(200, "design deployed"), nil
	})

	DesignCmd.SetArgs([]string{"deploy", "-f", filepath.Join("fixtures", "sampleDesign.golden"), "-s", "manifest"})
	assert.NoError(t, DesignCmd.Execute())
	assert.Equal(t, string(core.K8sManifest), designDeployFlags.SourceType)
}
