package filter

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/require"
)

func TestImportCmdLocalRelativeFile(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(FilterCmd, t)

	utils.TokenFlag = utils.GetToken(t)
	previousDir, err := os.Getwd()
	require.NoError(t, err)
	tempDir := t.TempDir()
	require.NoError(t, os.Chdir(tempDir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(previousDir)) })

	const fileName = "my-filter.wasm"
	fileContent := []byte("wasm filter bytes")
	require.NoError(t, os.WriteFile(fileName, fileContent, 0600))

	var request map[string]interface{}
	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/filter", func(req *http.Request) (*http.Response, error) {
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payload, &request); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(200, "{}"), nil
	})

	FilterCmd.SetArgs([]string{"import", fileName})
	require.NoError(t, FilterCmd.Execute())

	require.NotContains(t, request, "url")
	filterData, ok := request["filterData"].(map[string]interface{})
	require.True(t, ok)
	require.Equal(t, fileName, filterData["name"])
	require.Equal(t, base64.StdEncoding.EncodeToString(fileContent), filterData["filterFile"])
}

func TestImportCmdRemoteURL(t *testing.T) {
	testContext := utils.InitTestEnvironment(t)
	_ = utils.SetupMeshkitLoggerTesting(t, false)
	defer utils.StopMockery(t)
	defer utils.ResetCommandFlags(FilterCmd, t)

	utils.TokenFlag = utils.GetToken(t)

	var request map[string]interface{}
	httpmock.RegisterResponder("POST", testContext.BaseURL+"/api/filter", func(req *http.Request) (*http.Response, error) {
		payload, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(payload, &request); err != nil {
			return nil, err
		}
		return httpmock.NewStringResponse(200, "{}"), nil
	})

	const remoteURL = "https://example.com/my-filter.wasm"
	FilterCmd.SetArgs([]string{"import", remoteURL})
	require.NoError(t, FilterCmd.Execute())

	require.Equal(t, remoteURL, request["url"])
}

