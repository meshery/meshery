package meshsync

import (
	"os"
	"testing"

	"github.com/layer5io/meshkit/database"
	"github.com/stretchr/testify/assert"
)

var runIntegrationTest bool
var pathToSQLFile string

func init() {
	runIntegrationTest = os.Getenv("RUN_INTEGRATION_TESTS") == "true"
	pathToSQLFile = os.Getenv("PATH_TO_SQL_FILE")
}

func TestBasedOnDatabaseContentMeshSyncIntegration(t *testing.T) {
	if !runIntegrationTest {
		t.Skip("skipping integration test")
	}

	dboptions := database.Options{
		Engine:   database.SQLITE,
		Filename: pathToSQLFile,
	}
	handler, errDatabaseNew := database.New(dboptions)
	if errDatabaseNew != nil {
		assert.NoError(t, errDatabaseNew, "expect no error on  database connection")
		t.Fatal(errDatabaseNew)
	}
	// NOTE : didn't find how to close the handler

	for _, tc := range testCaseBasedOnDatabaseContentData {
		for _, cleanupHook := range tc.cleanupHooks {
			defer cleanupHook()
		}

		for _, setupHook := range tc.setupHooks {
			setupHook()
		}

		t.Run(
			tc.name,
			tc.run(handler),
		)
	}
}
