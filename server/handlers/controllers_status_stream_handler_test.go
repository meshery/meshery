package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
)

func newSSETestHandler(t *testing.T) *Handler {
	t.Helper()
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("logger init: %v", err)
	}
	return &Handler{log: log}
}

// nonFlusherRecorder is an http.ResponseWriter that does NOT implement
// http.Flusher, used to exercise the streaming-not-supported guard.
type nonFlusherRecorder struct {
	header http.Header
	status int
	buf    []byte
}

func (n *nonFlusherRecorder) Header() http.Header {
	if n.header == nil {
		n.header = http.Header{}
	}
	return n.header
}
func (n *nonFlusherRecorder) Write(p []byte) (int, error) { n.buf = append(n.buf, p...); return len(p), nil }
func (n *nonFlusherRecorder) WriteHeader(s int)           { n.status = s }

func TestControllersStatusStreamHandler_MissingConnectionID(t *testing.T) {
	h := newSSETestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/controllers/status/stream", nil)
	w := httptest.NewRecorder()

	h.ControllersStatusStreamHandler(w, req, nil, nil, nil)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status: got %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestControllersStatusStreamHandler_NonFlusherWriter(t *testing.T) {
	h := newSSETestHandler(t)
	req := httptest.NewRequest(http.MethodGet, "/api/system/kubernetes/controllers/status/stream?connectionID=abc", nil)
	w := &nonFlusherRecorder{}

	h.ControllersStatusStreamHandler(w, req, nil, nil, nil)

	if w.status != http.StatusInternalServerError {
		t.Fatalf("status: got %d, want %d", w.status, http.StatusInternalServerError)
	}
}

// TestDiffControllerObservations_InitialSeed exercises the !diffOnly path —
// every observation must be emitted regardless of whether snapshot already
// holds an entry for it. Mirrors the legacy resolver's first-frame behavior
// that the SSE handler preserves on the initial seed.
func TestDiffControllerObservations_InitialSeed(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{}
	obs := []controllerObservation{
		{
			connectionID: "conn-a",
			controller:   models.MesheryBroker,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
		},
		{
			connectionID: "conn-a",
			controller:   models.Meshsync,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
		},
	}

	items := diffControllerObservations(obs, snapshot, false)

	if got, want := len(items), 2; got != want {
		t.Fatalf("items emitted: got %d, want %d", got, want)
	}
	if _, ok := snapshot["conn-a"]; !ok {
		t.Fatalf("snapshot was not seeded for conn-a")
	}
	if got, want := len(snapshot["conn-a"]), 2; got != want {
		t.Fatalf("snapshot entries for conn-a: got %d, want %d", got, want)
	}
}

// TestDiffControllerObservations_UnchangedTickEmitsNothing is the core dedup
// test — when the same status+version is observed on a subsequent tick, no
// item should be emitted. This is the behavior that replaces the deleted
// client-side `_.isEqual` comparator (ui/components/subscription/comparatorFns).
func TestDiffControllerObservations_UnchangedTickEmitsNothing(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{
		"conn-a": {
			models.MesheryBroker: {Status: controllers.Deployed, Version: "v1"},
		},
	}
	obs := []controllerObservation{{
		connectionID: "conn-a",
		controller:   models.MesheryBroker,
		statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
	}}

	items := diffControllerObservations(obs, snapshot, true)

	if len(items) != 0 {
		t.Fatalf("unchanged observation emitted %d items, want 0", len(items))
	}
}

// TestDiffControllerObservations_StatusChangeEmits verifies that a changed
// status is detected and emitted. The snapshot must also be updated to the
// new value so a subsequent identical tick is then deduped.
func TestDiffControllerObservations_StatusChangeEmits(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{
		"conn-a": {
			models.MesheryBroker: {Status: controllers.Deployed, Version: "v1"},
		},
	}
	obs := []controllerObservation{{
		connectionID: "conn-a",
		controller:   models.MesheryBroker,
		statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.NotDeployed, Version: "v1"},
	}}

	items := diffControllerObservations(obs, snapshot, true)

	if len(items) != 1 {
		t.Fatalf("status change emitted %d items, want 1", len(items))
	}
	if got, want := items[0].Status, model.GetInternalControllerStatus(controllers.NotDeployed); got != want {
		t.Fatalf("emitted status: got %v, want %v", got, want)
	}
	if got, want := snapshot["conn-a"][models.MesheryBroker].Status, controllers.NotDeployed; got != want {
		t.Fatalf("snapshot status not updated: got %v, want %v", got, want)
	}
}

// TestDiffControllerObservations_VersionChangeEmits is the bug guard for the
// "or-chain on fields" alternative — using reflect.DeepEqual catches Version
// changes even when Status is unchanged. If the dedup ever regresses to
// status-only comparison this test fails.
func TestDiffControllerObservations_VersionChangeEmits(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{
		"conn-a": {
			models.MesheryBroker: {Status: controllers.Deployed, Version: "v1"},
		},
	}
	obs := []controllerObservation{{
		connectionID: "conn-a",
		controller:   models.MesheryBroker,
		statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v2"},
	}}

	items := diffControllerObservations(obs, snapshot, true)

	if len(items) != 1 {
		t.Fatalf("version change emitted %d items, want 1", len(items))
	}
	if got, want := items[0].Version, "v2"; got != want {
		t.Fatalf("emitted version: got %q, want %q", got, want)
	}
}

// TestDiffControllerObservations_NewConnectionEmits ensures a never-before-
// seen (connectionID, controller) tuple is emitted on the diffOnly path, not
// just on the initial seed. This matters when the UI rebinds to an additional
// connection mid-stream — that connection's first tick should always emit.
func TestDiffControllerObservations_NewConnectionEmits(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{
		"conn-a": {
			models.MesheryBroker: {Status: controllers.Deployed, Version: "v1"},
		},
	}
	obs := []controllerObservation{
		{
			connectionID: "conn-a",
			controller:   models.MesheryBroker,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
		},
		{
			connectionID: "conn-b",
			controller:   models.MesheryBroker,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
		},
	}

	items := diffControllerObservations(obs, snapshot, true)

	if len(items) != 1 {
		t.Fatalf("expected only the new conn-b item, got %d", len(items))
	}
	if got, want := items[0].ConnectionID, "conn-b"; got != want {
		t.Fatalf("emitted connectionID: got %q, want %q", got, want)
	}
	if _, ok := snapshot["conn-b"]; !ok {
		t.Fatalf("snapshot was not extended for conn-b")
	}
}

// TestDiffControllerObservations_MultipleControllersIndependent verifies that
// the dedup is per-(connection, controller) and not per-connection — a
// changed status on one controller emits without dragging in the other
// controllers' unchanged statuses.
func TestDiffControllerObservations_MultipleControllersIndependent(t *testing.T) {
	snapshot := map[string]map[models.MesheryController]models.MesheryControllerStatusAndVersion{
		"conn-a": {
			models.MesheryBroker: {Status: controllers.Deployed, Version: "v1"},
			models.Meshsync:      {Status: controllers.Deployed, Version: "v1"},
		},
	}
	obs := []controllerObservation{
		{
			connectionID: "conn-a",
			controller:   models.MesheryBroker,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.NotDeployed, Version: "v1"},
		},
		{
			connectionID: "conn-a",
			controller:   models.Meshsync,
			statusAndVer: models.MesheryControllerStatusAndVersion{Status: controllers.Deployed, Version: "v1"},
		},
	}

	items := diffControllerObservations(obs, snapshot, true)

	if len(items) != 1 {
		t.Fatalf("expected only the changed broker item, got %d", len(items))
	}
	if got, want := items[0].Controller, model.GetInternalController(models.MesheryBroker); got != want {
		t.Fatalf("emitted controller: got %v, want %v", got, want)
	}
}
