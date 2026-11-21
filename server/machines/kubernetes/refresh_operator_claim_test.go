package kubernetes

import (
	"reflect"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/internal/sql"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// refreshTestContext builds a context whose kubeconfig GenerateClientSetAction
// can turn into a clientset without contacting the cluster.
func refreshTestContext(connectionID, server, token string, serverID core.Uuid) models.K8sContext {
	return models.K8sContext{
		ID:                 "ctx-" + connectionID,
		Name:               "ctx-" + connectionID,
		ConnectionID:       connectionID,
		Server:             server,
		KubernetesServerID: &serverID,
		Cluster: sql.Map{
			"name":    "cluster",
			"cluster": map[string]interface{}{"server": server},
		},
		Auth: sql.Map{
			"name": "user",
			"user": map[string]interface{}{"token": token},
		},
	}
}

// claimAsOperatorConsumer records k8sctx's connection as a consumer of its
// cluster's operator the way the connect setup pass does.
func claimAsOperatorConsumer(t *testing.T, log logger.Handler, ot *models.OperatorTracker, k8sctx models.K8sContext) {
	t.Helper()
	models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, nil, nil).
		SetMeshsyncDeploymentMode(connections.MeshsyncDeploymentModeOperator).
		DeployUndeployedOperators(ot, k8sctx)
}

func TestRefreshK8sContext_ReleasesOperatorClaimOnPreviousCluster(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to build test logger: %v", err)
	}
	oldClusterID := core.Uuid(uuid.Must(uuid.NewV4()))
	newClusterID := core.Uuid(uuid.Must(uuid.NewV4()))
	const oldURL = "https://old.example.invalid:6443"

	tests := []struct {
		name       string
		next       func(prev models.K8sContext) models.K8sContext
		wantOnPrev []string // connection IDs still consuming the previous cluster's operator
	}{
		{
			name: "moves to a different cluster",
			next: func(prev models.K8sContext) models.K8sContext {
				return refreshTestContext(prev.ConnectionID, "https://new.example.invalid:6443", "rotated", newClusterID)
			},
			wantOnPrev: []string{"other"},
		},
		{
			name: "rotates credentials on the same cluster",
			next: func(prev models.K8sContext) models.K8sContext {
				return refreshTestContext(prev.ConnectionID, oldURL, "rotated", oldClusterID)
			},
			wantOnPrev: []string{"other", "refreshed"},
		},
		{
			name: "reaches the same cluster through another address",
			next: func(prev models.K8sContext) models.K8sContext {
				return refreshTestContext(prev.ConnectionID, "https://node.example.invalid:6443", "rotated", oldClusterID)
			},
			wantOnPrev: []string{"other", "refreshed"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := models.NewOperatorTracker(false)
			prev := refreshTestContext("refreshed", oldURL, "original", oldClusterID)
			other := refreshTestContext("other", "https://lb.example.invalid:6443", "other", oldClusterID)
			claimAsOperatorConsumer(t, log, ot, prev)
			claimAsOperatorConsumer(t, log, ot, other)

			machinectx := &MachineCtx{K8sContext: prev, OperatorTracker: ot, log: log}
			changed, err := machinectx.RefreshK8sContext(tt.next(prev), events.NewEvent(), log)
			if err != nil {
				t.Fatalf("RefreshK8sContext: %v", err)
			}
			if !changed {
				t.Fatal("RefreshK8sContext reported no change for new connection material")
			}

			if got := ot.OperatorConsumers(other); !reflect.DeepEqual(got, tt.wantOnPrev) {
				t.Fatalf("consumers of the previous cluster after the refresh = %v, want %v", got, tt.wantOnPrev)
			}
		})
	}
}
