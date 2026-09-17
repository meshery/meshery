package adapter

import (
	"path/filepath"
	"testing"

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

	_, err = waitForValidateResponse(mctlCfg, "Smi conformance test")

	assert.Error(t, err)
	assert.Equal(t, ErrCreatingValidateResponseRequestCode, errors.GetCode(err))
}
