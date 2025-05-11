package meshsync

import (
	"testing"

	"github.com/layer5io/meshkit/database"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
)

type testCaseBasedOnDatabaseContentStruct struct {
	name         string
	setupHooks   []func()
	cleanupHooks []func()
	run          func(database.Handler) func(*testing.T)
}

var testCaseBasedOnDatabaseContentData []testCaseBasedOnDatabaseContentStruct = []testCaseBasedOnDatabaseContentStruct{
	{
		name: "number of entries in kubernetes_resources must be greater than zero",
		run: func(handler database.Handler) func(*testing.T) {
			return func(t *testing.T) {
				handler.Lock()
				defer handler.Unlock()

				k8sResources := make([]*meshsyncmodel.KubernetesResource, 0, 256)

				// Query the database for the complete component definition
				dbresult := handler.Model(meshsyncmodel.KubernetesResource{}).Find(&k8sResources)

				if dbresult == nil {
					t.Fatal("db result is nil")
				}
				if dbresult.Error != nil {
					t.Fatalf("db result ended with an error %v", dbresult.Error)
				}

				assert.True(t, len(k8sResources) > 0, "table must contains resources")
				t.Logf("found %d k8s resources", len(k8sResources))

				for _, resource := range k8sResources {
					t.Logf(
						"found k8s resource with api = %s, kind = %s",
						resource.APIVersion,
						resource.Kind,
					)
				}
			}
		},
	},
}
