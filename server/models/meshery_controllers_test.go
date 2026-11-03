package models

import (
	"errors"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
)

func TestControllerEventActedUponPrefersConnectionID(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	connectionID := uuid.Must(uuid.NewV4())

	actedUpon := controllerEventActedUpon(userID, map[string]any{
		"connectionID": connectionID.String(),
	})

	if actedUpon != connectionID {
		t.Fatalf("expected actedUpon to use connectionID, got %s", actedUpon)
	}
}

func TestControllerEventActedUponFallsBackToUserID(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())

	actedUpon := controllerEventActedUpon(userID, map[string]any{
		"connectionID": "not-a-uuid",
	})

	if actedUpon != userID {
		t.Fatalf("expected actedUpon to fall back to userID, got %s", actedUpon)
	}
}

func TestControllerEventActedUponReadsConnectionId(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	connectionID := uuid.Must(uuid.NewV4())

	actedUpon := controllerEventActedUpon(userID, map[string]any{
		"connectionId": connectionID,
	})

	if actedUpon != connectionID {
		t.Fatalf("expected actedUpon to use connectionId, got %s", actedUpon)
	}
}

func TestControllerEventActedUponReturnsNilWithoutValidIDs(t *testing.T) {
	actedUpon := controllerEventActedUpon(uuid.Nil, map[string]any{
		"connectionID": "",
	})

	if actedUpon != uuid.Nil {
		t.Fatalf("expected actedUpon to be nil UUID, got %s", actedUpon)
	}
}

func TestShouldPersistControllerEvent(t *testing.T) {
	userID := uuid.Must(uuid.NewV4())
	resourceID := uuid.Must(uuid.NewV4())

	if !shouldPersistControllerEvent(userID, resourceID) {
		t.Fatal("expected controller event to be persisted when user and resource IDs are valid")
	}

	if shouldPersistControllerEvent(uuid.Nil, resourceID) {
		t.Fatal("expected controller event persistence to be skipped when user ID is nil")
	}

	if shouldPersistControllerEvent(userID, uuid.Nil) {
		t.Fatal("expected controller event persistence to be skipped when actedUpon is nil")
	}
}

func TestAddCtxControllerHandlersReturnsEarlyOnInvalidConfig(t *testing.T) {
	// Create an empty K8sContext which will result in an invalid kubeconfig
	// and eventually a failure in mesherykube.New()
	ctx := K8sContext{
		ID:   "test-context",
		Name: "test-cluster",
	}

	log, _ := logger.New("test", logger.Options{})

	mch := NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		nil,
		nil,
	)

	// A successful execution of this function without panicking indicates the fix is working.
	// We wrap in a defer-recover just to explicitly fail the test if a panic occurs.
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddCtxControllerHandlers panicked, indicating the nil pointer bug is still present: %v", r)
		}
	}()

	mch.AddCtxControllerHandlers(ctx)

	// Additionally, verify that ctxControllerHandlers is still empty/nil because it should have returned early
	if len(mch.ctxControllerHandlers) != 0 {
		t.Fatalf("expected ctxControllerHandlers to be empty, got %d", len(mch.ctxControllerHandlers))
	}
}

func recvControllerEvent(t *testing.T, ch <-chan interface{}) *events.Event {
	t.Helper()
	select {
	case raw := <-ch:
		event, ok := raw.(*events.Event)
		if !ok {
			t.Fatalf("expected *events.Event, got %T", raw)
		}
		return event
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for controller event")
	}
	return nil
}

func TestEmitEventActedUponAndSuppressesRepeat(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}
	systemID := uuid.Must(uuid.NewV4())
	connectionID := uuid.Must(uuid.NewV4())
	broadcaster := NewBroadcaster("test")
	ch, unsubscribe := broadcaster.Subscribe(systemID)
	defer unsubscribe()

	mch := &MesheryControllersHelper{
		log:              log,
		eventBroadcaster: broadcaster,
		systemID:         &systemID,
		connectionID:     connectionID,
	}

	mch.emitErrorEvent("Failed to deploy Meshery Operator", errors.New("boom"), map[string]any{
		"k8sContextID": "ctx",
	}, uuid.Nil)
	mch.emitErrorEvent("Failed to deploy Meshery Operator", errors.New("boom"), map[string]any{
		"k8sContextID": "ctx",
	}, uuid.Nil)

	got := recvControllerEvent(t, ch)
	if got.ActedUpon != connectionID {
		t.Fatalf("expected ActedUpon %s, got %s", connectionID, got.ActedUpon)
	}
	if got.Metadata["connectionId"] != connectionID.String() {
		t.Fatalf("expected metadata connectionId %s, got %#v", connectionID, got.Metadata["connectionId"])
	}
	select {
	case extra := <-ch:
		t.Fatalf("repeat controller error was emitted: %#v", extra)
	case <-time.After(200 * time.Millisecond):
	}

	mch.emitEvent("MeshSync connected", events.Informational, map[string]any{}, uuid.Nil)
	if success := recvControllerEvent(t, ch); success.Severity != events.Informational {
		t.Fatalf("expected the success event to still be emitted, got %q", success.Severity)
	}

	mch.emitErrorEvent("Failed to deploy Meshery Operator", errors.New("boom"), nil, uuid.Nil)
	again := recvControllerEvent(t, ch)
	if again.ActedUpon != connectionID {
		t.Fatalf("expected ActedUpon %s after success cleared suppression, got %s", connectionID, again.ActedUpon)
	}
}
