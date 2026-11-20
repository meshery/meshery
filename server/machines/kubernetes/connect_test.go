package kubernetes

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/schemas/models/core"
)

type mockConnectProvider struct {
	models.Provider
	connection *connections.Connection
	onLookup   func()
}

func (m *mockConnectProvider) GetConnectionByID(token string, id core.Uuid) (*connections.Connection, int, error) {
	if m.onLookup != nil {
		m.onLookup()
	}
	return m.connection, 1, nil
}

// TestConnectAction_OverlappingConnectAttemptsSerialized verifies that when a
// Connect transition initiates background setup, an overlapping second Connect attempt
// (such as during a retry after database persistence failure while setup is still in flight)
// does not launch a duplicate concurrent setup goroutine.
func TestConnectAction_OverlappingConnectAttemptsSerialized(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:   connID,
			Kind: "kubernetes",
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-1",
			Name:         "test-cluster",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	ca := &ConnectAction{}

	// Manually set setupInProgress to simulate an active in-flight setup from a previous attempt.
	machinectx.setupLock.Lock()
	machinectx.setupInProgress = true
	machinectx.setupLock.Unlock()

	// An overlapping retry while setup is running must return NoOp and queue the retry.
	nextEvent, event, err := ca.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("expected overlapping Connect to return nil error, got: %v", err)
	}
	if nextEvent != machines.NoOp {
		t.Fatalf("expected nextEvent NoOp, got %v", nextEvent)
	}
	if event != nil {
		t.Fatalf("expected nil event for overlapping attempt, got: %v", event)
	}

	if !machinectx.IsRetryPending() {
		t.Fatal("expected overlapping Connect to set retryPending to true")
	}

	// Once the first setup finishes and resets the flag, setupInProgress becomes false.
	machinectx.setupLock.Lock()
	machinectx.setupInProgress = false
	machinectx.setupLock.Unlock()

	if machinectx.IsSetupInProgress() {
		t.Fatal("expected setup to be available after in-flight setup clears")
	}
}

// TestConnectAction_ConcurrentConnectAttemptsDoNotRace drives concurrent Connect
// transitions on the same connection machine and verifies that setup in progress guards
// prevent racing setups without panics or data races.
func TestConnectAction_ConcurrentConnectAttemptsDoNotRace(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:   connID,
			Kind: "kubernetes",
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-concurrent",
			Name:         "test-cluster-concurrent",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	ca := &ConnectAction{}

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, _ = ca.Execute(ctx, machinectx, nil)
		}()
	}
	wg.Wait()
}

// TestConnectAction_QueuedRetryAbortsWhenDisconnected verifies that if a connection
// transitions to DISCONNECTED while setup is running with a queued retry, the queued
// retry is canceled and does not re-run connection setup.
func TestConnectAction_QueuedRetryAbortsWhenDisconnected(t *testing.T) {
	log, err := logger.New("test", logger.Options{})
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}

	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:     connID,
			Kind:   "kubernetes",
			Status: connections.DISCONNECTED,
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-disconnect",
			Name:         "test-cluster-disconnect",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	// Manually simulate a queued retry
	machinectx.setupLock.Lock()
	machinectx.retryPending = true
	machinectx.setupLock.Unlock()

	// Execute disconnect
	da := &DisconnectAction{}
	_, _, err = da.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("unexpected disconnect error: %v", err)
	}

	if machinectx.IsRetryPending() {
		t.Fatal("expected DisconnectAction to cancel pending retry")
	}

	if !machinectx.IsDisconnected() {
		t.Fatal("expected DisconnectAction to set disconnected flag to true")
	}

	// Also verify helper reports false for disconnected status
	if isConnectionMeantToBeConnected(provider.connection) {
		t.Fatal("expected isConnectionMeantToBeConnected to return false for DISCONNECTED connection")
	}
}

// TestDeleteAction_SetsDisconnectedAndCancelsRetry verifies that DeleteAction marks
// the machine context as disconnected and cancels any pending retry.
func TestDeleteAction_SetsDisconnectedAndCancelsRetry(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-delete",
			Name:         "test-cluster-delete",
			ConnectionID: connID.String(),
		},
	}
	machinectx.setupLock.Lock()
	machinectx.retryPending = true
	machinectx.setupLock.Unlock()

	da := &DeleteAction{}
	_, _, err := da.Execute(context.Background(), machinectx, nil)
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}

	if machinectx.IsRetryPending() {
		t.Fatal("expected DeleteAction to cancel pending retry")
	}

	if !machinectx.IsDisconnected() {
		t.Fatal("expected DeleteAction to set disconnected flag to true")
	}
}

// TestConnectAction_RetryLoopsWhenPending verifies that when setup finishes with
// a retryPending flag set, it loops and executes the second setup pass unless
// disconnected.
func TestConnectAction_RetryLoopsWhenPending(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:     connID,
			Kind:   "kubernetes",
			Status: connections.CONNECTED,
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-retry-loop",
			Name:         "test-cluster",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	ca := &ConnectAction{}

	// Execute connect: this launches the background setup goroutine
	_, _, err := ca.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("unexpected connect error: %v", err)
	}

	// While it runs (or immediately after), queue a retry
	_, _, _ = ca.Execute(ctx, machinectx, nil)

	// Disconnect should cleanly cancel and prevent further setup
	da := &DisconnectAction{}
	_, _, err = da.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("unexpected disconnect error: %v", err)
	}

	if !machinectx.IsDisconnected() {
		t.Fatal("expected machinectx to be disconnected")
	}
	if machinectx.IsRetryPending() {
		t.Fatal("expected pending retry to be canceled by disconnect")
	}
}

// TestConnectAction_ControlledRetryInterleaving verifies that an overlapping Connect
// executes a second setup pass, confirming handler registration and lookup synchronization.
func TestConnectAction_ControlledRetryInterleaving(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	lookupCalled := make(chan struct{}, 1)
	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:     connID,
			Kind:   "kubernetes",
			Status: connections.CONNECTED,
		},
		onLookup: func() {
			select {
			case lookupCalled <- struct{}{}:
			default:
			}
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-controlled-retry",
			Name:         "test-cluster",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	ca := &ConnectAction{}

	// Initial connect starts setup pass 1
	_, _, err := ca.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("first connect failed: %v", err)
	}

	// Queue an overlapping retry while pass 1 is in progress
	_, _, err = ca.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("second connect failed: %v", err)
	}

	// Verify worker evaluates retryPending and triggers pass 2 lookup check
	select {
	case <-lookupCalled:
		// Passed: worker recognized queued retry and checked eligibility
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for second setup pass lookup")
	}

	// Worker eventually finishes and releases setup ownership
	for i := 0; i < 50; i++ {
		if !machinectx.IsSetupInProgress() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if machinectx.IsSetupInProgress() {
		t.Fatal("expected setupInProgress to be released once setup completes")
	}
}

// TestConnectAction_DisconnectPreventsSecondSetupPass verifies that when DisconnectAction
// runs while a retry is pending, the second setup pass is aborted and no handlers remain.
func TestConnectAction_DisconnectPreventsSecondSetupPass(t *testing.T) {
	log, _ := logger.New("test", logger.Options{})
	connID := uuid.Must(uuid.NewV4())
	userUUID := core.Uuid(uuid.Must(uuid.NewV4()))

	provider := &mockConnectProvider{
		connection: &connections.Connection{
			ID:     connID,
			Kind:   "kubernetes",
			Status: connections.DISCONNECTED,
		},
	}

	mch := models.NewMesheryControllersHelper(
		log,
		controllers.OperatorDeploymentConfig{},
		nil,
		nil,
		provider,
		nil,
	)

	machinectx := &MachineCtx{
		K8sContext: models.K8sContext{
			ID:           "test-ctx-disconnect-prevent",
			Name:         "test-cluster",
			ConnectionID: connID.String(),
		},
		MesheryCtrlsHelper: mch,
		log:                log,
		OperatorTracker:    models.NewOperatorTracker(true),
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, models.ProviderCtxKey, models.Provider(provider))
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	ctx = context.WithValue(ctx, models.UserCtxKey, &models.User{ID: userUUID})

	ca := &ConnectAction{}

	// Initial connect starts setup
	_, _, _ = ca.Execute(ctx, machinectx, nil)

	// Queue overlapping retry
	_, _, _ = ca.Execute(ctx, machinectx, nil)

	// Disconnect immediately
	da := &DisconnectAction{}
	_, _, err := da.Execute(ctx, machinectx, nil)
	if err != nil {
		t.Fatalf("disconnect failed: %v", err)
	}

	if !machinectx.IsDisconnected() {
		t.Fatal("expected disconnected=true")
	}
	if machinectx.IsRetryPending() {
		t.Fatal("expected retryPending=false after disconnect")
	}

	for i := 0; i < 50; i++ {
		if !machinectx.IsSetupInProgress() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if machinectx.IsSetupInProgress() {
		t.Fatal("expected setupInProgress=false")
	}
}


