package adapter

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/meshery/meshery/mesheryctl/internal/cli/root/config"
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
	"github.com/meshery/meshkit/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWaitForValidateResponseWithMissingToken(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.SetupMeshkitLoggerTesting(t, false)

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	if err != nil {
		t.Fatal(err)
	}

	originalToken := utils.TokenFlag
	t.Cleanup(func() { utils.TokenFlag = originalToken })
	utils.TokenFlag = filepath.Join(t.TempDir(), "absent-token.json")

	_, err = waitForValidateResponse(mctlCfg, "Smi conformance test", "")

	assert.Error(t, err)
	assert.Equal(t, ErrCreatingValidateResponseRequestCode, errors.GetCode(err))
}

func TestParseAndValidateMesh(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.SetupMeshkitLoggerTesting(t, false)
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "sync.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/system/sync",
		httpmock.NewStringResponder(200, apiResponse))
	utils.TokenFlag = filepath.Join(fixturesDir, "token.golden")

	mctlCfg, err := config.GetMesheryCtl(viper.GetViper())
	assert.NoError(t, err)

	// Single word arg
	err = parseAndValidateMesh(mctlCfg, []string{"linkerd"})
	assert.NoError(t, err)
	assert.Equal(t, "LINKERD", meshName)

	// Multi-word positional args
	err = parseAndValidateMesh(mctlCfg, []string{"cilium", "service", "mesh"})
	assert.NoError(t, err)
	assert.Equal(t, "CILIUM_SERVICE_MESH", meshName)

	// Hyphenated arg
	err = parseAndValidateMesh(mctlCfg, []string{"cilium-service-mesh"})
	assert.NoError(t, err)
	assert.Equal(t, "CILIUM_SERVICE_MESH", meshName)

	// Unknown mesh returns ErrNoAdapters
	err = parseAndValidateMesh(mctlCfg, []string{"unknown-mesh"})
	assert.Error(t, err)
	assert.Equal(t, ErrNoAdaptersCode, errors.GetCode(err))
}

func TestDeployAndRemovePreRunE(t *testing.T) {
	utils.SetupContextEnv(t)
	utils.SetupMeshkitLoggerTesting(t, false)
	utils.StartMockery(t)
	defer utils.StopMockery(t)

	testContext := utils.NewTestHelper(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")

	apiResponse := utils.NewGoldenFile(t, "sync.golden", fixturesDir).Load()
	httpmock.RegisterResponder("GET", testContext.BaseURL+"/api/system/sync",
		httpmock.NewStringResponder(200, apiResponse))
	utils.TokenFlag = filepath.Join(fixturesDir, "token.golden")

	// Test deployCmd PreRunE parses positional args into meshName
	meshName = ""
	err := deployCmd.PreRunE(deployCmd, []string{"linkerd"})
	assert.NoError(t, err)
	assert.Equal(t, "LINKERD", meshName)

	// Test removeCmd PreRunE parses positional args into meshName
	meshName = ""
	err = removeCmd.PreRunE(removeCmd, []string{"istio"})
	assert.NoError(t, err)
	assert.Equal(t, "ISTIO", meshName)
}
