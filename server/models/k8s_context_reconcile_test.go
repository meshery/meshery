package models

import (
	stderrors "errors"
	"net/http"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/schemas/models/core"
)

// reconcileFakeProvider is a partial Provider: it embeds the interface so the
// type satisfies Provider, and implements only the two methods
// ReconcileK8sContextServerID calls. Any other method would panic on the nil
// embedded interface, keeping the test honest about what the function touches.
type reconcileFakeProvider struct {
	Provider

	conn      *connections.Connection
	getErr    error
	updateErr error

	getCalls    int
	updateCalls int
	lastPayload *connections.ConnectionPayload
}

func (f *reconcileFakeProvider) GetConnectionByID(_ string, _ core.Uuid) (*connections.Connection, int, error) {
	f.getCalls++
	if f.getErr != nil {
		return nil, http.StatusInternalServerError, f.getErr
	}
	if f.conn == nil {
		return nil, http.StatusNotFound, nil
	}
	return f.conn, http.StatusOK, nil
}

func (f *reconcileFakeProvider) UpdateConnectionById(_ string, conn *connections.ConnectionPayload, _ string) (*connections.Connection, error) {
	f.updateCalls++
	f.lastPayload = conn
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return &connections.Connection{ID: conn.ID, Kind: conn.Kind, Status: conn.Status, Metadata: conn.MetaData}, nil
}

func newReconcileServerID(t *testing.T) (*core.Uuid, string) {
	t.Helper()
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatalf("failed to generate server UUID: %v", err)
	}
	u := core.Uuid(id)
	return &u, u.String()
}

func reconcileK8sCtx(connID string, serverID *core.Uuid) K8sContext {
	return K8sContext{
		Name:               "ctx",
		ConnectionID:       connID,
		KubernetesServerID: serverID,
	}
}

// A connection first registered while its cluster was unreachable persists an
// empty kubernetesServerId. Once the cluster is reachable and discovery resolves
// the real ID, the reconcile must back-fill it so the dashboard's cluster_id
// filter starts matching the streamed MeshSync rows.
func TestReconcileK8sContextServerID_BackfillsEmpty(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, serverIDStr := newReconcileServerID(t)

	f := &reconcileFakeProvider{
		conn: &connections.Connection{
			ID:       uuid.FromStringOrNil(connID),
			Kind:     "kubernetes",
			Status:   connections.CONNECTED,
			Metadata: core.Map{"kubernetesServerId": "", "name": "ctx"},
		},
	}

	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err != nil {
		t.Fatalf("reconcile returned error: %v", err)
	}
	if f.updateCalls != 1 {
		t.Fatalf("expected exactly 1 update call, got %d", f.updateCalls)
	}
	if got, _ := f.lastPayload.MetaData["kubernetesServerId"].(string); got != serverIDStr {
		t.Fatalf("expected metadata kubernetesServerId %q, got %q", serverIDStr, got)
	}
	// Unrelated metadata must be preserved through the update.
	if name, _ := f.lastPayload.MetaData["name"].(string); name != "ctx" {
		t.Fatalf("expected unrelated metadata preserved, got name %q", name)
	}
	// The reconcile must not change the connection's status.
	if f.lastPayload.Status != connections.CONNECTED {
		t.Fatalf("expected status preserved as CONNECTED, got %q", f.lastPayload.Status)
	}
}

// The remote provider persists a nil UUID (00000000-...) rather than an empty
// string for an unassigned server ID; the reconcile must treat it as stale too.
func TestReconcileK8sContextServerID_CorrectsNilUUID(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, serverIDStr := newReconcileServerID(t)

	f := &reconcileFakeProvider{
		conn: &connections.Connection{
			ID:       uuid.FromStringOrNil(connID),
			Kind:     "kubernetes",
			Status:   connections.CONNECTED,
			Metadata: core.Map{"kubernetesServerId": uuid.Nil.String()},
		},
	}

	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err != nil {
		t.Fatalf("reconcile returned error: %v", err)
	}
	if f.updateCalls != 1 {
		t.Fatalf("expected exactly 1 update call, got %d", f.updateCalls)
	}
	if got, _ := f.lastPayload.MetaData["kubernetesServerId"].(string); got != serverIDStr {
		t.Fatalf("expected metadata kubernetesServerId %q, got %q", serverIDStr, got)
	}
}

// A connection whose persisted server ID already matches must incur no write, so
// the per-request discovery re-drive costs nothing on the steady state.
func TestReconcileK8sContextServerID_NoopWhenAlreadyCorrect(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, serverIDStr := newReconcileServerID(t)

	f := &reconcileFakeProvider{
		conn: &connections.Connection{
			ID:       uuid.FromStringOrNil(connID),
			Kind:     "kubernetes",
			Status:   connections.CONNECTED,
			Metadata: core.Map{"kubernetesServerId": serverIDStr},
		},
	}

	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err != nil {
		t.Fatalf("reconcile returned error: %v", err)
	}
	if f.updateCalls != 0 {
		t.Fatalf("expected no update call when already correct, got %d", f.updateCalls)
	}
}

func TestReconcileK8sContextServerID_HandlesNilMetadata(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, serverIDStr := newReconcileServerID(t)

	f := &reconcileFakeProvider{
		conn: &connections.Connection{
			ID:       uuid.FromStringOrNil(connID),
			Kind:     "kubernetes",
			Status:   connections.CONNECTED,
			Metadata: nil,
		},
	}

	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err != nil {
		t.Fatalf("reconcile returned error: %v", err)
	}
	if f.updateCalls != 1 {
		t.Fatalf("expected exactly 1 update call, got %d", f.updateCalls)
	}
	if got, _ := f.lastPayload.MetaData["kubernetesServerId"].(string); got != serverIDStr {
		t.Fatalf("expected metadata kubernetesServerId %q, got %q", serverIDStr, got)
	}
}

// Guard the early returns: a missing/nil server id or an invalid connection id
// must not touch the provider at all (nothing authoritative to sync).
func TestReconcileK8sContextServerID_NoopOnMissingInputs(t *testing.T) {
	serverID, _ := newReconcileServerID(t)
	nilID := core.Uuid(uuid.Nil)

	cases := []struct {
		name string
		ctx  K8sContext
	}{
		{"nil server id", reconcileK8sCtx(uuid.Must(uuid.NewV4()).String(), nil)},
		{"nil-uuid server id", reconcileK8sCtx(uuid.Must(uuid.NewV4()).String(), &nilID)},
		{"empty connection id", reconcileK8sCtx("", serverID)},
		{"invalid connection id", reconcileK8sCtx("not-a-uuid", serverID)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &reconcileFakeProvider{}
			if err := ReconcileK8sContextServerID(f, "token", tc.ctx); err != nil {
				t.Fatalf("reconcile returned error: %v", err)
			}
			if f.getCalls != 0 || f.updateCalls != 0 {
				t.Fatalf("expected no provider calls, got get=%d update=%d", f.getCalls, f.updateCalls)
			}
		})
	}
}

func TestReconcileK8sContextServerID_PropagatesGetError(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, _ := newReconcileServerID(t)

	f := &reconcileFakeProvider{getErr: stderrors.New("boom")}
	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err == nil {
		t.Fatal("expected an error when GetConnectionByID fails, got nil")
	}
	if f.updateCalls != 0 {
		t.Fatalf("expected no update call when the read fails, got %d", f.updateCalls)
	}
}

func TestReconcileK8sContextServerID_PropagatesUpdateError(t *testing.T) {
	connID := uuid.Must(uuid.NewV4()).String()
	serverID, _ := newReconcileServerID(t)

	f := &reconcileFakeProvider{
		conn: &connections.Connection{
			ID:       uuid.FromStringOrNil(connID),
			Kind:     "kubernetes",
			Status:   connections.CONNECTED,
			Metadata: core.Map{"kubernetesServerId": ""},
		},
		updateErr: stderrors.New("boom"),
	}
	if err := ReconcileK8sContextServerID(f, "token", reconcileK8sCtx(connID, serverID)); err == nil {
		t.Fatal("expected an error when UpdateConnectionById fails, got nil")
	}
	if f.updateCalls != 1 {
		t.Fatalf("expected the update to be attempted once, got %d", f.updateCalls)
	}
}
