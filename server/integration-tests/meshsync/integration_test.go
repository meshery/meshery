package meshsync

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var runIntegrationTest bool

func init() {
	runIntegrationTest = os.Getenv("RUN_INTEGRATION_TESTS") == "true"
}

func TestMeshSyncIntegration(t *testing.T) {
	if !runIntegrationTest {
		t.Skip("skipping integration test")
	}

	t.Run(
		"test case will be here soon...",
		func(t *testing.T) {
			assert.Equal(t, 1, 1, "1 must be equal to 1")
		},
	)
}
