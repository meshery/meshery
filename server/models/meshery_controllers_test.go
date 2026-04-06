package models

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
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
