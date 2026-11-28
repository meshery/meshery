package design

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

var invalidDesignSourceType = "invalid-source"

func TestDeployCmd(t *testing.T) {
	// create a test helper
	testContext := utils.NewTestHelper(t)

	// get current directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	// test scenrios for fetching data
	tests := []utils.MesheryMultiURLCommamdTest{
		{
			Name:             "given valid file and source type when design deploy then design is deployed",
			Args:             []string{"deploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-s", "Kubernetes Manifest"},
			ExpectedResponse: "deploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/types",
					Response:     "view.designTypes.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/Kubernetes%20Manifest",
					Response:     "deploy.applicationSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "deploy.designdeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file and --skip-save when design deploy then design is deployed without being saved",
			Args:             []string{"deploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "--skip-save", "-s", "Kubernetes Manifest"},
			ExpectedResponse: "deploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/types",
					Response:     "view.designTypes.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "deploy.designdeploy.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given invalid source type when design deploy then throw error",
			Args:             []string{"deploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-s", invalidDesignSourceType},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/types",
					Response:     "view.designTypes.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  utils.ErrFlagsInvalid(fmt.Errorf("Invalid value for --source-type '%s': valid values are %s", invalidDesignSourceType, strings.Join(GetDesignSourceTypes(), ", "))),
		},
		{
			Name:             "given valid file and shorthand source type when design deploy then design is deployed",
			Args:             []string{"deploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-s", "manifest"},
			ExpectedResponse: "deploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "deploy.designdeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file and canonical schema source type when design deploy then design is deployed",
			Args:             []string{"deploy", "-f", filepath.Join(fixturesDir, "sampleDesign.golden"), "-s", "k8s-manifest"},
			ExpectedResponse: "deploy.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/import",
					Response:     "apply.designSave.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "POST",
					URL:          testContext.BaseURL + "/api/pattern/deploy",
					Response:     "deploy.designdeploy.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given non-existent design name when design deploy then throw error",
			Args:             []string{"deploy", "nonexistent-design"},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern/types",
					Response:     "view.designTypes.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&search=nonexistent-design",
					Response:     "pattern.empty.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound("nonexistent-design"),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

func Test_designDeploy_FlagValidation(t *testing.T) {
	tests := []struct {
		name       string
		run        func(t *testing.T) error
		want       error
		wantSubstr string
	}{
		{
			name: "given no file and no design name when design deploy then ErrDeployDesign is returned",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = ""
				t.Cleanup(func() { designDeployFlags.File = previous })

				return deployDesignCmd.Args(deployDesignCmd, nil)
			},
			want: ErrDeployDesign(),
		},
		{
			name: "given whitespace-only file flag when design deploy then ErrDeployDesign is returned",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = "   "
				t.Cleanup(func() { designDeployFlags.File = previous })

				return deployDesignCmd.Args(deployDesignCmd, nil)
			},
			want: ErrDeployDesign(),
		},
		{
			name: "given whitespace-only design name when design deploy then ErrInvalidArgument is returned",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = ""
				t.Cleanup(func() { designDeployFlags.File = previous })

				return deployDesignCmd.Args(deployDesignCmd, []string{"   "})
			},
			wantSubstr: "design name cannot be empty",
		},
		{
			name: "given both design name and --file flag when design deploy then ErrInvalidArgument is returned",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = "sampleDesign.yaml"
				t.Cleanup(func() { designDeployFlags.File = previous })

				return deployDesignCmd.Args(deployDesignCmd, []string{"my-design"})
			},
			wantSubstr: "cannot specify both a design name and the --file flag",
		},
		{
			name: "given valid file with surrounding whitespace when design deploy then whitespace is trimmed and succeeds",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = "  sampleDesign.yaml  "
				t.Cleanup(func() { designDeployFlags.File = previous })

				err := deployDesignCmd.Args(deployDesignCmd, nil)
				if err != nil {
					return err
				}
				if designDeployFlags.File != "sampleDesign.yaml" {
					return fmt.Errorf("expected trimmed file 'sampleDesign.yaml', got %q", designDeployFlags.File)
				}
				return nil
			},
			want: nil,
		},
		{
			name: "given valid design name when design deploy then succeeds",
			run: func(t *testing.T) error {
				previous := designDeployFlags.File
				designDeployFlags.File = ""
				t.Cleanup(func() { designDeployFlags.File = previous })

				return deployDesignCmd.Args(deployDesignCmd, []string{"my-design"})
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			err := tt.run(t)
			if tt.want == nil && tt.wantSubstr == "" {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tt.wantSubstr != "" {
				assert.Contains(t, err.Error(), tt.wantSubstr)
				return
			}
			if err.Error() != tt.want.Error() {
				t.Fatalf("error mismatch\n got: %v\nwant: %v", err, tt.want)
			}
		})
	}
}
