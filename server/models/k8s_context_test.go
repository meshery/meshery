package models

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
)

// TestFlushMeshSyncDataNilKubernetesServerIDNoPanic verifies that FlushMeshSyncData
// does not panic when K8sContext entries have a nil KubernetesServerID. Before the
// fix, the refcount loop called KubernetesServerID.String() unconditionally,
// triggering a nil dereference whenever a context had no associated server ID.
func TestFlushMeshSyncDataNilKubernetesServerIDNoPanic(t *testing.T) {
	// All contexts have nil KubernetesServerID — the scenario that previously caused
	// a panic. With the nil guard in place, refCount stays 0 and the flush block is
	// skipped entirely, so no provider, broadcast, or logger calls are made.
	k8sctxs := []*K8sContext{
		{ID: "ctx-1", Name: "cluster-a", KubernetesServerID: nil},
		{ID: "ctx-2", Name: "cluster-b", KubernetesServerID: nil},
	}
	ctx := context.WithValue(context.Background(), AllKubeClusterKey, k8sctxs)
	k8sCtx := K8sContext{ID: "ctx-1", Name: "cluster-a", Server: "https://k8s.example.com"}

	// Should complete without panic. If the guard is missing this will panic with a
	// nil dereference on KubernetesServerID.String() inside the refcount loop.
	FlushMeshSyncData(ctx, k8sCtx, nil, nil, "00000000-0000-0000-0000-000000000000", nil, nil)
}

// TestFlushMeshSyncDataMixedNilAndPopulatedServerIDs verifies that the refcount loop
// correctly skips nil entries while still counting non-nil ones, and that no panic
// occurs when both nil and populated KubernetesServerID values are present.
func TestFlushMeshSyncDataMixedNilAndPopulatedServerIDs(t *testing.T) {
	serverID, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate UUID: %v", err)
	}

	k8sctxs := []*K8sContext{
		{ID: "ctx-1", Name: "cluster-a", KubernetesServerID: nil},
		{ID: "ctx-2", Name: "cluster-b", KubernetesServerID: &serverID},
		{ID: "ctx-3", Name: "cluster-c", KubernetesServerID: nil},
	}
	ctx := context.WithValue(context.Background(), AllKubeClusterKey, k8sctxs)
	// ctx-1 has nil ServerID — refCount for its sid ("") will be 0, flush skipped.
	k8sCtx := K8sContext{ID: "ctx-1", Name: "cluster-a", Server: "https://k8s.example.com"}

	FlushMeshSyncData(ctx, k8sCtx, nil, nil, "00000000-0000-0000-0000-000000000000", nil, nil)
}
