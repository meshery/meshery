package meshsync

import (
	"os"
	"testing"

	"github.com/layer5io/meshkit/database"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
)

var runIntegrationTest bool
var pathToSQLFile string

func init() {
	runIntegrationTest = os.Getenv("RUN_INTEGRATION_TESTS") == "true"
	pathToSQLFile = os.Getenv("PATH_TO_SQL_FILE")
}

func TestMeshSyncIntegration(t *testing.T) {
	if !runIntegrationTest {
		t.Skip("skipping integration test")
	}

	t.Run(
		"test case will be here soon...",
		func(t *testing.T) {
			dboptions := database.Options{
				Engine: database.SQLITE,
				// TODO
				Filename: pathToSQLFile,
			}
			handler, errDatabaseNew := database.New(dboptions)
			if errDatabaseNew != nil {
				assert.NoError(t, errDatabaseNew, "expect no error on  database connection")
				t.Fatal(errDatabaseNew)
			}
			assert.NotNil(t, handler)

			handler.Lock()
			defer handler.Unlock()

			k8sResources := make([]*meshsyncmodel.KubernetesResource, 0, 8)

			// Query the database for the complete component definition
			dbresult := handler.Model(meshsyncmodel.KubernetesResource{}).Find(&k8sResources)

			if dbresult == nil {
				t.Fatal("db result is nil")
			}
			if dbresult.Error != nil {
				t.Fatalf("db result ended with an error %v", dbresult.Error)
			}

			// TODO add more meaningfull check
			assert.True(t, len(k8sResources) > 0, "table must contains resources")
			t.Logf("found %d k8s resources", len(k8sResources))

			for _, resource := range k8sResources {
				t.Logf(
					"found k8s resource with api = %s, kind = %s",
					resource.APIVersion,
					resource.Kind,
				)
			}
		},
	)
}
