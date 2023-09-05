package system

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	utils.SetupContextEnv(t)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Not able to get current working directory")
	}
	currDir := filepath.Dir(filename)
	fixturesDir := filepath.Join(currDir, "fixtures")
	tests := []CmdTestInput{
		{
			Name:             "Login with Meshery",
			Args:             []string{"login", "-p", "Meshery"},
			ExpectedResponse: "successfully authenticated." + "\n",
			Token:            filepath.Join(fixturesDir, "token.golden"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			utils.TokenFlag = tt.Token
			SetupFunc()
			SystemCmd.SetArgs(tt.Args)
			err = SystemCmd.Execute()
			if err != nil {
				t.Error(err)
			}
			actualResponse := b.String()
			expectedResponse := tt.ExpectedResponse
			assert.Equal(t, expectedResponse, actualResponse)
			BreakupFunc()
		})
	}
}
