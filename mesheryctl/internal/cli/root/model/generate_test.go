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

func extractRequestFields(t *testing.T, bodyBytes []byte) (model, displayName, registrant, category, subCategory string, register bool) {
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		t.Fatalf("Failed to parse request body: %v", err)
	}

	if registerRaw, ok := payload["register"]; ok {
		if err := json.Unmarshal(registerRaw, &register); err != nil {
			t.Fatalf("Failed to parse register field: %v", err)
		}
	}

	importBodyRaw, ok := payload["importBody"]
	if !ok {
		t.Fatalf("missing required field: importBody")
	}

	var importBody map[string]json.RawMessage
	if err := json.Unmarshal(importBodyRaw, &importBody); err != nil {
		t.Fatalf("Failed to parse importBody: %v", err)
	}

	modelRaw, ok := importBody["model"]
	if !ok {
		t.Fatalf("missing required field: model")
	}

	var modelObj map[string]json.RawMessage
	if err := json.Unmarshal(modelRaw, &modelObj); err != nil {
		t.Fatalf("Failed to parse model: %v", err)
	}

	decodeString := func(field string, raw json.RawMessage, target *string) {
		if raw != nil {
			if err := json.Unmarshal(raw, target); err != nil {
				t.Fatalf("Failed to parse string field %s: %v", field, err)
			}
		}
	}

	decodeString("model", modelObj["model"], &model)
	decodeString("modelDisplayName", modelObj["modelDisplayName"], &displayName)
	decodeString("registrant", modelObj["registrant"], &registrant)
	decodeString("category", modelObj["category"], &category)
	decodeString("subCategory", modelObj["subCategory"], &subCategory)

	return
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

				model, displayName, registrant, category, subCategory, register := extractRequestFields(t, bodyBytes)

				if model != "untitled-model" {
					t.Errorf("Expected model 'untitled-model', got '%s'", model)
				}
				if displayName != "Untitled Model" {
					t.Errorf("Expected displayName 'Untitled Model', got '%s'", displayName)
				}
				if registrant != "artifacthub" {
					t.Errorf("Expected registrant 'artifacthub', got '%s'", registrant)
				}
				if category != "Uncategorized" {
					t.Errorf("Expected category 'Uncategorized', got '%s'", category)
				}
				if subCategory != "Uncategorized" {
					t.Errorf("Expected subCategory 'Uncategorized', got '%s'", subCategory)
				}
				if register {
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
				model, _, registrant, category, subCategory, register := extractRequestFields(t, bodyBytes)

				if model != "cert-manager" {
					t.Errorf("Expected model 'cert-manager', got '%s'", model)
				}
				if registrant != "github" {
					t.Errorf("Expected registrant 'github', got '%s'", registrant)
				}
				if category != "Security" {
					t.Errorf("Expected category 'Security', got '%s'", category)
				}
				if subCategory != "Certificates" {
					t.Errorf("Expected subCategory 'Certificates', got '%s'", subCategory)
				}
				if !register {
					t.Errorf("Expected Register to be true, got %v", register)
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

				model, displayName, _, _, _, register := extractRequestFields(t, bodyBytes)

				if model != "minimal-model" {
					t.Errorf("Expected model 'minimal-model', got '%s'", model)
				}
				if displayName != "Minimal Model" {
					t.Errorf("Expected displayName 'Minimal Model', got '%s'", displayName)
				}
				if register {
					t.Errorf("Expected Register to be false, got true")
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
