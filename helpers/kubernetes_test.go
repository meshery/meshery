package helpers

import (
	"testing"

	"github.com/layer5io/meshery/mesheryctl/pkg/utils"
)

func TestFetchKubernetesVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping tests")
	}

	utils.SetupContextEnv(t)
	utils.StartMockery(t)
	testByte := []byte{0, 0}
	var testContext = "nil"

	_, err := FetchKubernetesVersion(testByte, testContext)
	if err != nil {
		t.Error("FetchKubernetesVersion() failed")
	}
}
