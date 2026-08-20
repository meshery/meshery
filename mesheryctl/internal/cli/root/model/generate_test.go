package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	mesheryctlflags "github.com/meshery/meshery/mesheryctl/internal/cli/pkg/flags"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// generateTestImportRequest represents only the wire-format fields
// inspected by the model generate request validators.
type generateTestImportRequest struct {
	Register   bool `json:"register"`
	ImportBody struct {
		Model struct {
			Model            string `json:"model"`
			ModelDisplayName string `json:"modelDisplayName"`
			Registrant       string `json:"registrant"`
			Category         string `json:"category"`
			SubCategory      string `json:"subCategory"`
		} `json:"model"`
	} `json:"importBody"`
}

func TestModelGenerate(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.StartMockery(t)

	testContext := utils.NewTestHelper(t)

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	testdataDir := filepath.Join(currDir, "testdata")
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiURL := "/api/registry/register"

	type tc struct {
		Name             string
		Args             []string
		URL              string
		Fixture          string
		ExpectedResponse string
		ExpectHelp       bool
		ExpectErr        bool
		RaisedError      error
		HttpCode         int
		ValidateRequest  func(req *http.Request, t *testing.T)
	}

	tests := []tc{
		{
			Name:             "model generate: no args prints help and errors",
			Args:             []string{"generate"},
			ExpectedResponse: "generate.no-args.output.golden",
			ExpectHelp:       true,
			ExpectErr:        true,
			RaisedError:      utils.ErrInvalidArgument(fmt.Errorf(errGenerateMissingArgsMsg, errGenerateUsageMsg)),
		},
		{
			Name:             "model generate: from CSV directory",
			Args:             []string{"generate", "--file", filepath.Join(fixturesDir, "templates", "template-csvs")},
			ExpectedResponse: "generate.dir.registered.output.golden",
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			HttpCode:         200,
		},
		{
			Name:             "model generate: from URL with template",
			Args:             []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "template.json"), "--skip-registration"},
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			ExpectedResponse: "generate.dir.skipped.output.golden",
			HttpCode:         200,
		},
		{
			Name:             "model generate: from URL with real nested model init template",
			Args:             []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "model_template.json"), "--skip-registration"},
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			ExpectedResponse: "generate.dir.skipped.output.golden",
			HttpCode:         200,
			ValidateRequest: func(req *http.Request, t *testing.T) {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				var importReq generateTestImportRequest
				err = json.Unmarshal(bodyBytes, &importReq)
				if err != nil {
					t.Fatalf("Failed to parse request body: %v", err)
				}

				if importReq.ImportBody.Model.Model != "untitled-model" {
					t.Errorf("Expected model 'untitled-model', got '%s'", importReq.ImportBody.Model.Model)
				}
				if importReq.ImportBody.Model.ModelDisplayName != "Untitled Model" {
					t.Errorf("Expected displayName 'Untitled Model', got '%s'", importReq.ImportBody.Model.ModelDisplayName)
				}
				if importReq.ImportBody.Model.Registrant != "artifacthub" {
					t.Errorf("Expected registrant 'artifacthub', got '%s'", importReq.ImportBody.Model.Registrant)
				}
				if string(importReq.ImportBody.Model.Category) != "Uncategorized" {
					t.Errorf("Expected category 'Uncategorized', got '%s'", importReq.ImportBody.Model.Category)
				}
				if string(importReq.ImportBody.Model.SubCategory) != "Uncategorized" {
					t.Errorf("Expected subCategory 'Uncategorized', got '%s'", importReq.ImportBody.Model.SubCategory)
				}
				if importReq.Register {
					t.Errorf("Expected Register to be false, got true")
				}
			},
		},
		{
			Name:             "model generate: from URL with default registration",
			Args:             []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "template.json")},
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			ExpectedResponse: "generate.dir.skipped.output.golden",
			HttpCode:         200,
			ValidateRequest: func(req *http.Request, t *testing.T) {
				bodyBytes, _ := io.ReadAll(req.Body)
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				var importReq generateTestImportRequest
				if err := json.Unmarshal(bodyBytes, &importReq); err != nil {
					t.Fatalf("Failed to parse request body: %v", err)
				}

				if importReq.ImportBody.Model.Model != "cert-manager" {
					t.Errorf("Expected model 'cert-manager', got '%s'", importReq.ImportBody.Model.Model)
				}
				if importReq.ImportBody.Model.Registrant != "github" {
					t.Errorf("Expected registrant 'github', got '%s'", importReq.ImportBody.Model.Registrant)
				}
				if string(importReq.ImportBody.Model.Category) != "Security" {
					t.Errorf("Expected category 'Security', got '%s'", importReq.ImportBody.Model.Category)
				}
				if string(importReq.ImportBody.Model.SubCategory) != "Certificates" {
					t.Errorf("Expected subCategory 'Certificates', got '%s'", importReq.ImportBody.Model.SubCategory)
				}
				if !importReq.Register {
					t.Errorf("Expected Register to be true, got %v", importReq.Register)
				}
			},
		},
		{
			Name:             "model generate: from URL with minimal legacy template",
			Args:             []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "template_minimal.json"), "--skip-registration"},
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			ExpectedResponse: "generate.dir.skipped.output.golden",
			HttpCode:         200,
			ValidateRequest: func(req *http.Request, t *testing.T) {
				bodyBytes, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("Failed to read request body: %v", err)
				}
				// Restore the body so other readers (if any) can read it
				req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				var importReq generateTestImportRequest
				err = json.Unmarshal(bodyBytes, &importReq)
				if err != nil {
					t.Fatalf("Failed to parse request body: %v", err)
				}

				if importReq.ImportBody.Model.Model != "minimal-model" {
					t.Errorf("Expected model 'minimal-model', got '%s'", importReq.ImportBody.Model.Model)
				}
				if importReq.ImportBody.Model.ModelDisplayName != "Minimal Model" {
					t.Errorf("Expected displayName 'Minimal Model', got '%s'", importReq.ImportBody.Model.ModelDisplayName)
				}
			},
		},
	}

	var resetFlags func(*cobra.Command, *testing.T)
	resetFlags = func(c *cobra.Command, t *testing.T) {
		c.Flags().VisitAll(func(f *pflag.Flag) {
			if err := f.Value.Set(f.DefValue); err != nil {
				t.Fatalf("failed to reset flag %q: %v", f.Name, err)
			}
		})
		for _, sub := range c.Commands() {
			resetFlags(sub, t)
		}
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			httpmock.Reset()
			resetFlags(ModelCmd, t)

			if tt.URL != "" {
				apiResponse := utils.NewGoldenFile(t, tt.Fixture, fixturesDir).LoadByte()

				httpmock.RegisterResponder("POST", testContext.BaseURL+tt.URL, func(req *http.Request) (*http.Response, error) {
					if tt.ValidateRequest != nil {
						tt.ValidateRequest(req, t)
					}
					return httpmock.NewBytesResponse(tt.HttpCode, apiResponse), nil
				})
			}

			utils.TokenFlag = utils.GetToken(t)

			golden := utils.NewGoldenFile(t, tt.ExpectedResponse, testdataDir)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			ModelCmd.SetOut(b)
			ModelCmd.SetArgs(tt.Args)
			mesheryctlflags.InitValidators(ModelCmd)
			err := ModelCmd.Execute()

			if tt.ExpectHelp || tt.ExpectErr {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
				t.Logf("[%s] stderr (error):\n%s", tt.Name, err.Error())
				utils.AssertMeshkitErrorsEqual(t, tt.RaisedError, err)
				return
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			actualResponse := utils.StripAnsiEscapeCodes(b.String())
			t.Logf("[%s] stdout:\n%s", tt.Name, actualResponse)

			expectedResponse := strings.TrimSpace(golden.Load())
			utils.Equals(t, expectedResponse, strings.TrimSpace(actualResponse))
		})
	}

	utils.StopMockery(t)
}
