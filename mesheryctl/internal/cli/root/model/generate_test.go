package model

import (
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
	"github.com/stretchr/testify/assert"
)

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

	apiURL := "/api/meshmodels/register"

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
		AssertRequest    func(t *testing.T, body []byte)
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
			Args:             []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "template.json"), "--skip-registration=true"},
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			ExpectedResponse: "generate.dir.skipped.output.golden",
			HttpCode:         200,
		},
		{
			Name:             "model generate: from CSV directory with selected model",
			Args:             []string{"generate", "--file", filepath.Join(fixturesDir, "templates", "template-csvs"), "--model", "couchbase"},
			ExpectedResponse: "generate.dir.skip-register.output.golden",
			URL:              apiURL,
			Fixture:          "generate.api.ok.response.golden",
			HttpCode:         200,
			AssertRequest: func(t *testing.T, body []byte) {
				t.Helper()
				var payload struct {
					UploadType string `json:"uploadType"`
					ImportBody struct {
						Model struct {
							Model string `json:"model"`
						} `json:"model"`
					} `json:"importBody"`
				}
				err := json.Unmarshal(body, &payload)
				if err != nil {
					t.Fatalf("failed to unmarshal request body: %v", err)
				}
				assert.Equal(t, "csv", payload.UploadType)
				assert.Equal(t, "couchbase", payload.ImportBody.Model.Model)
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
					if tt.AssertRequest != nil {
						reqBody, err := io.ReadAll(req.Body)
						if err != nil {
							t.Fatalf("failed to read request body: %v", err)
						}
						tt.AssertRequest(t, reqBody)
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
