package design

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

// Test_designImport_FlagValidation covers the two validation paths
// `mesheryctl design import` actually has, and asserts the error each one
// returns is the expected one.
//
// It previously called importPattern directly. importPattern performs no flag
// validation at all - the source type is checked in RunE via
// retrieveProvidedSourceType, and the required -f flag in the command's Args -
// so both cases were passing on an incidental "no such file or directory" from
// os.ReadFile while the declared `want` was never compared to anything. Both
// assertions below fail if that error identity regresses.
func Test_designImport_FlagValidation(t *testing.T) {
	tests := []struct {
		name string
		run  func(t *testing.T) error
		want error
	}{
		{
			name: "given invalid source type when design import then error is thrown",
			run: func(_ *testing.T) error {
				_, err := retrieveProvidedSourceType("invalid source type", validDesignSourceTypes)
				return err
			},
			want: ErrInValidSource("invalid source type", validDesignSourceTypes),
		},
		{
			name: "given missing file flag when design import then error is thrown",
			run: func(t *testing.T) error {
				// `file` is a package-level flag target shared by every design
				// subcommand, so restore it rather than leaving it blank for
				// whatever test runs next.
				previous := file
				file = ""
				t.Cleanup(func() { file = previous })

				return importCmd.Args(importCmd, nil)
			},
			want: ErrDesignFileNotProvided(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The Args validator logs before returning; without this the
			// nil logger panics.
			_ = utils.SetupMeshkitLoggerTesting(t, false)

			err := tt.run(t)
			if err == nil {
				t.Fatalf("expected error %v, got nil", tt.want)
			}
			if err.Error() != tt.want.Error() {
				t.Fatalf("error mismatch\n got: %v\nwant: %v", err, tt.want)
			}
		})
	}
}

// Test_importPattern_SendsCamelCaseFileName is a regression test for the design
// import wire-contract fix. The /api/pattern/import File-import oneOf variant
// requires the camelCase `fileName` field; the legacy snake_case `file_name`
// left the variant unmatched and every file import failed with
// "Invalid design import request" (meshery-server-1422). This captures the exact
// JSON body mesheryctl puts on the wire and asserts the canonical field is sent.
func Test_importPattern_SendsCamelCaseFileName(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	var capturedBody map[string]interface{}
	httpmock.RegisterResponder(
		"POST",
		testContext.BaseURL+"/api/pattern/import",
		func(req *http.Request) (*http.Response, error) {
			raw, err := io.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.NoError(t, json.Unmarshal(raw, &capturedBody))
			t.Logf("Captured /api/pattern/import request body keys: %v", mapKeys(capturedBody))
			t.Logf("fileName field value: %q", capturedBody["fileName"])
			return httpmock.NewStringResponse(
				200,
				`[{"id":"3817ec9a-1d83-4f6f-9154-0fd4408ba9f0","name":"SampleApp"}]`,
			), nil
		},
	)

	DesignCmd.SetArgs([]string{"import", "-f", "fixtures/sampleDesign.golden", "-n", "SampleApp"})
	assert.NoError(t, DesignCmd.Execute())

	// The fix: canonical camelCase `fileName` present, legacy snake_case absent.
	assert.Contains(t, capturedBody, "fileName", "request body must carry camelCase fileName")
	assert.NotContains(t, capturedBody, "file_name", "request body must not carry legacy snake_case file_name")
	assert.Equal(t, "sampleDesign.golden", capturedBody["fileName"])
}

// Test_importPattern_SendsContractUrlPayload covers the URL oneOf arm, the one
// the camelCase regression test above cannot reach. Both arms are now built from
// the same generated MesheryPatternImportRequestBody, so the URL arm has to carry
// exactly the contract's `url`/`name` fields - and specifically NOT the legacy
// `save` key, which was never part of the request contract and which the server's
// DesignFileImportHandler never read.
func Test_importPattern_SendsContractUrlPayload(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	output := utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	var capturedBody map[string]interface{}
	httpmock.RegisterResponder(
		"POST",
		testContext.BaseURL+"/api/pattern/import",
		func(req *http.Request) (*http.Response, error) {
			raw, err := io.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.NoError(t, json.Unmarshal(raw, &capturedBody))
			t.Logf("Captured URL-import request body: %s", raw)
			return httpmock.NewStringResponse(
				200,
				`[{"id":"3817ec9a-1d83-4f6f-9154-0fd4408ba9f0","name":"RemoteApp"}]`,
			), nil
		},
	)

	DesignCmd.SetArgs([]string{
		"import", "-f", "https://example.com/design.yaml", "-n", "RemoteApp",
	})
	assert.NoError(t, DesignCmd.Execute())

	assert.Equal(t, "https://example.com/design.yaml", capturedBody["url"])
	assert.Equal(t, "RemoteApp", capturedBody["name"])
	assert.NotContains(t, capturedBody, "save", "`save` is not part of the import request contract")
	assert.ElementsMatch(t, []string{"url", "name"}, mapKeys(capturedBody))
	t.Logf("mesheryctl output: %s", output.String())
}

// Test_importPattern_EmptyResponseCollection pins the guard added when both oneOf
// arms were folded into postDesignImport: the server returns the saved design(s)
// as a collection, and both arms previously indexed response[0] unguarded, so an
// empty collection panicked with an index-out-of-range instead of surfacing an
// error the user can act on.
func Test_importPattern_EmptyResponseCollection(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer resetVariables()
	defer utils.ResetCommandFlags(DesignCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	httpmock.RegisterResponder(
		"POST",
		testContext.BaseURL+"/api/pattern/import",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, `[]`), nil
		},
	)

	DesignCmd.SetArgs([]string{"import", "-f", "fixtures/sampleDesign.golden", "-n", "SampleApp"})

	err := DesignCmd.Execute()
	assert.Error(t, err, "an empty design collection must surface as an error, not a panic")
	assert.Equal(t, ErrDesignInvalidApiResponse("design import returned no design").Error(), err.Error())
	t.Logf("structured error returned instead of panicking: %v", err)
}

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
