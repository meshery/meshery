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

func Test_importPattern_DisplayErrorsMissingFlags(t *testing.T) {
	type args struct {
		sourceType string
		file       string
		patternURL string
		save       bool
	}

	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name:    "given invalid source type when design import then error is thrown",
			args:    args{"invalid source type", "file.yaml", "", false},
			want:    ErrInValidSource("invalid source type", validDesignSourceTypes),
			wantErr: true,
		},
		{
			name:    "given missing file flag when design import then error is thrown",
			args:    args{"helm", "", "", false},
			want:    ErrDesignFileNotProvided(),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := importPattern(tt.args.sourceType, tt.args.file, tt.args.patternURL, tt.args.save)
			if (err != nil) != tt.wantErr {
				t.Errorf("importPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
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

func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
