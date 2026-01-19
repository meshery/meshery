package model

import (
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiURL := "/api/meshmodels/register"

	type tc struct {
		Name                 string
		Args                 []string
		URL                  string
		Fixture              string
		ExpectedContains     []string
		ExpectedErrContains  []string
		ExpectHelp           bool
		ExpectErr            bool
		HttpCode             int
	}

	tests := []tc{
		{
			Name:                "model generate: no args prints help and errors",
			Args:                []string{"generate"},
			ExpectedErrContains: []string{"isn't specified", "mesheryctl model generate"},
			ExpectHelp:          true,
			ExpectErr:           true,
		},
		{
			Name: "model generate: from CSV directory",
			Args: []string{"generate", "--file", filepath.Join(fixturesDir, "templates", "template-csvs")},
			ExpectedContains: []string{
				"Generating model from CSV files",
				"Imported model couchbase",
				"Model can be accessed from .meshery/models",
				"Logs for the csv generation can be accessed .meshery/logs/registry",
			},
			URL:      apiURL,
			Fixture:  "generate.api.ok.response.golden",
			HttpCode: 200,
		},
		{
			Name: "model generate: from URL with template",
			Args: []string{"generate", "--file", "https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.crds.yaml", "--template", filepath.Join(fixturesDir, "templates", "template.json"), "--register=true"},
			URL:      apiURL,
			Fixture:  "generate.api.ok.response.golden",
			HttpCode: 200,
			ExpectedContains: []string{
				"Generating model from URL",
				"Imported model couchbase",
				"Model can be accessed from .meshery/models",
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
				apiResponse := utils.ReadTestFixtureBytes(t, fixturesDir, tt.Fixture)

				httpmock.RegisterResponder("POST", testContext.BaseURL+tt.URL, func(req *http.Request) (*http.Response, error) {

					return httpmock.NewBytesResponse(tt.HttpCode, apiResponse), nil
				})
			}

			utils.TokenFlag = utils.GetToken(t)

			b := utils.SetupMeshkitLoggerTesting(t, false)
			ModelCmd.SetOut(b)
			ModelCmd.SetArgs(tt.Args)
			err := ModelCmd.Execute()

			if tt.ExpectHelp || tt.ExpectErr {
				if err == nil {
					t.Fatal("expected an error, but got nil")
				}
				t.Logf("[%s] stderr (error):\n%s", tt.Name, err.Error())
				for _, expected := range tt.ExpectedErrContains {
					if !strings.Contains(err.Error(), expected) {
						t.Fatalf("expected error to contain %q, got %q", expected, err.Error())
					}
				}
				return
			}

			if err != nil {
				t.Fatalf("did not expect an error, but got: %v", err)
			}

			actualResponse := utils.StripAnsiEscapeCodes(b.String())
			t.Logf("[%s] stdout:\n%s", tt.Name, actualResponse)
			for _, expected := range tt.ExpectedContains {
				if !strings.Contains(actualResponse, expected) {
					t.Fatalf("expected output to contain %q, got %q", expected, actualResponse)
				}
			}
		})
	}

	utils.StopMockery(t)
}
