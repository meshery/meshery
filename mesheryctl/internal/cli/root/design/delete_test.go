package design

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const nonExistentID = "a12b3c4d-5e6f-4890-abcd-ef1234567890"

func TestDeleteCmd(t *testing.T) {
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
			Name:             "given valid file when design delete then design is deleted",
			Args:             []string{"delete", "--yes", "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid file as positional argument when design delete then design is deleted",
			Args:             []string{"delete", "--yes", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given valid relative file as positional argument when design delete then design is deleted",
			Args:             []string{"delete", "--yes", filepath.Join("fixtures", "sampleDesign.golden")},
			ExpectedResponse: "delete.output.golden",
			URLs: []utils.MockURL{
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError: false,
		},
		{
			Name:             "given invalid file path when design delete then error is thrown",
			Args:             []string{"delete", "-f", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "given mistyped file path positionally when server unreachable then file read error is thrown",
			Args:             []string{"delete", invalidFilePath},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrFileRead(fmt.Errorf(errInvalidPathMsg, invalidFilePath)),
		},
		{
			Name:             "given nonexistent design ID when design delete then error is thrown",
			Args:             []string{"delete", nonExistentID},
			ExpectedResponse: "",
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "delete.idList.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/" + nonExistentID,
					Response:     "delete.error.response.golden",
					ResponseCode: 404,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  ErrDesignNotFound(nonExistentID),
		},
		{
			Name:             "given no arguments when design delete then error is thrown",
			Args:             []string{"delete"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("please provide a design name, ID, or file path")),
		},
		{
			Name:             "given empty positional argument when design delete then error is thrown",
			Args:             []string{"delete", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty")),
		},
		{
			Name:             "given whitespace positional argument when design delete then error is thrown",
			Args:             []string{"delete", "   "},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("target cannot be empty")),
		},
		{
			Name:             "given multiple positional arguments when design delete then error is thrown",
			Args:             []string{"delete", "a.yaml", "b.yaml"},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("expected at most one positional argument, got 2")),
		},
		{
			Name:             "given both positional argument and -f flag when design delete then error is thrown",
			Args:             []string{"delete", filepath.Join(fixturesDir, "sampleDesign.golden"), "-f", filepath.Join(fixturesDir, "sampleDesign.golden")},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("cannot specify both a positional argument and the --file flag")),
		},
		{
			Name:             "given empty file flag when design delete then error is thrown",
			Args:             []string{"delete", "-f", ""},
			ExpectedResponse: "",
			URLs:             []utils.MockURL{},
			ExpectError:      true,
			IsOutputGolden:   false,
			ExpectedError:    utils.ErrInvalidArgument(fmt.Errorf("file path cannot be empty")),
		},
		{
			Name:             "given ID prefix shared by two designs when design delete then ambiguity error is thrown and nothing is deleted",
			Args:             []string{"delete", "a12"},
			ExpectedResponse: "",
			// DELETE responders are registered so that any delete request would
			// succeed and fail this test, which expects the ambiguity error.
			URLs: []utils.MockURL{
				{
					Method:       "GET",
					URL:          testContext.BaseURL + "/api/pattern?populate=pattern_file&pagesize=10000",
					Response:     "delete.ambiguousIdList.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/a12b3c4d-5e6f-4890-abcd-ef1234567890",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/a12f9e8d-7c6b-4a59-8d4c-3b2a19087654",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
				{
					Method:       "DELETE",
					URL:          testContext.BaseURL + "/api/pattern/deploy?contexts=all",
					Response:     "delete.response.golden",
					ResponseCode: 200,
				},
			},
			ExpectError:    true,
			IsOutputGolden: false,
			ExpectedError:  utils.ErrAmbiguousID("a12", []string{"a12b3c4d-5e6f-4890-abcd-ef1234567890", "a12f9e8d-7c6b-4a59-8d4c-3b2a19087654"}),
		},
	}

	// Run tests
	utils.RunMesheryctlMultiURLTests(t, update, DesignCmd, tests, currDir, "design", resetVariables)
}

func TestDeleteFileConfirmation(t *testing.T) {
	target := filepath.Join("fixtures", "sampleDesign.golden")
	t.Run("decline", func(t *testing.T) {
		output, prompt, gets, _, deletes, err := runRemovalPromptTest(t, []string{"delete", target}, "n\n")
		require.NoError(t, err)
		assert.Equal(t, 1, gets)
		assert.Zero(t, deletes)
		assert.Contains(t, prompt, "Design 'SampleApp' will be undeployed from 2 connected Kubernetes cluster(s): east, west.")
		assert.Contains(t, output, "Delete cancelled")
	})
	t.Run("skip", func(t *testing.T) {
		_, prompt, gets, _, deletes, err := runRemovalPromptTest(t, []string{"delete", "-f", target, "--yes"}, "")
		require.NoError(t, err)
		assert.Zero(t, gets)
		assert.Equal(t, 1, deletes)
		assert.Empty(t, prompt)
	})
}
