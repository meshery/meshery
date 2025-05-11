package meshsync

import (
	"fmt"
	"testing"

	"github.com/layer5io/meshkit/database"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	"github.com/stretchr/testify/assert"
)

const CUSTOM_K8S_NAMESPACE = "calm-koala"

// as defined in infrastructure/test-deployment.yaml
const CUSTOM_APP_NAME = "nginx-deployment"
const CUSTOM_APP_REPLICAS_NUM = 3

type testCaseBasedOnDatabaseContentStruct struct {
	name         string
	setupHooks   []func()
	cleanupHooks []func()
	run          func(database.Handler) func(*testing.T)
}

// TODO add test case which checks connection for test cluster
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
	{
		name: "kubernetes_resources must contain one meshery.io/* Broker entity",
		run: func(handler database.Handler) func(*testing.T) {
			return func(t *testing.T) {
				handler.Lock()
				defer handler.Unlock()

				k8sResources := make([]*meshsyncmodel.KubernetesResource, 0, 2)

				// Query the database for the complete component definition
				dbresult := handler.
					Model(meshsyncmodel.KubernetesResource{}).
					Where("api_version LIKE ?", "meshery.io/%").
					Where("kind = ?", "Broker").
					Find(&k8sResources)

				if dbresult == nil {
					t.Fatal("db result is nil")
				}
				if dbresult.Error != nil {
					t.Fatalf("db result ended with an error %v", dbresult.Error)
				}

				assert.Equal(t, 1, len(k8sResources), "must contains exactly one Broker entity")
				assert.Equal(t, "Broker", k8sResources[0].Kind)
			}
		},
	},
	{
		name: "kubernetes_resources must contain one meshery.io/* MeshSync entity",
		run: func(handler database.Handler) func(*testing.T) {
			return func(t *testing.T) {
				handler.Lock()
				defer handler.Unlock()

				k8sResources := make([]*meshsyncmodel.KubernetesResource, 0, 2)

				// Query the database for the complete component definition
				dbresult := handler.
					Model(meshsyncmodel.KubernetesResource{}).
					Where("api_version LIKE ?", "meshery.io/%").
					Where("kind = ?", "MeshSync").
					Find(&k8sResources)

				if dbresult == nil {
					t.Fatal("db result is nil")
				}
				if dbresult.Error != nil {
					t.Fatalf("db result ended with an error %v", dbresult.Error)
				}

				assert.Equal(t, 1, len(k8sResources), "must contains exactly one Broker entity")
				assert.Equal(t, "MeshSync", k8sResources[0].Kind)
			}
		},
	},
	{
		// deployment should be only one, as it has the defined name
		// replica set and pods, could be more, in case it were some failures
		name: "custom namespace must contain custom app resources",
		run: func(handler database.Handler) func(*testing.T) {
			return func(t0 *testing.T) {
				t0.Run("1 deployment", func(t *testing.T) {
					handler.Lock()
					defer handler.Unlock()

					k8sResources := make([]*meshsyncmodel.KubernetesResourceObjectMeta, 0, 8)

					dbresult := handler.
						Model(&meshsyncmodel.KubernetesResourceObjectMeta{}).
						Joins("JOIN kubernetes_resources ON kubernetes_resource_object_meta.id = kubernetes_resources.id").
						Where("namespace = ?", CUSTOM_K8S_NAMESPACE).
						Where("kubernetes_resources.kind = ?", "Deployment").
						Find(&k8sResources)

					if dbresult == nil {
						t.Fatal("db result is nil")
					}
					if dbresult.Error != nil {
						t.Fatalf("db result ended with an error %v", dbresult.Error)
					}

					assert.Equal(t, 1, len(k8sResources), "must contains one deployment")
					if len(k8sResources) > 0 {
						assert.Equal(t, CUSTOM_APP_NAME, k8sResources[0].Name, "deployment must have correct name")
					}
				})
				t0.Run("at least 1 replica set", func(t *testing.T) {
					handler.Lock()
					defer handler.Unlock()

					k8sResources := make([]*meshsyncmodel.KubernetesResourceObjectMeta, 0, 8)

					dbresult := handler.
						Model(&meshsyncmodel.KubernetesResourceObjectMeta{}).
						Joins("JOIN kubernetes_resources ON kubernetes_resource_object_meta.id = kubernetes_resources.id").
						Where("namespace = ?", CUSTOM_K8S_NAMESPACE).
						Where("kubernetes_resources.kind = ?", "ReplicaSet").
						Find(&k8sResources)

					if dbresult == nil {
						t.Fatal("db result is nil")
					}
					if dbresult.Error != nil {
						t.Fatalf("db result ended with an error %v", dbresult.Error)
					}

					assert.GreaterOrEqual(t, len(k8sResources), 1, "must contains at least one replica set")
					for _, rs := range k8sResources {
						assert.Contains(t, rs.Name, CUSTOM_APP_NAME, "replica set must have correct name")
					}
				})
				t0.Run(fmt.Sprintf("at least %d pods", CUSTOM_APP_REPLICAS_NUM), func(t *testing.T) {})
			}
		},
	},
}
