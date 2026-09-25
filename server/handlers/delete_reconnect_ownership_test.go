package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/machines/kubernetes"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

type deleteBlockingProvider struct {
	lifecycleTestMockProvider
	entered   chan struct{}
	release   chan struct{}
	enterOnce sync.Once
}

func (p *deleteBlockingProvider) UpdateConnectionById(token string, conn *connections.ConnectionPayload, connID string) (*connections.Connection, error) {
	if conn != nil && conn.Status == connections.DELETED {
		p.enterOnce.Do(func() { close(p.entered) })
		<-p.release
	}
	return p.lifecycleTestMockProvider.UpdateConnectionById(token, conn, connID)
}

// TestDeleteCleanupKeepsMachineAdoptedByConcurrentReconnect pins the ownership
// hand-off between a delete and a reconnect that arrives while the delete's
// cleanup is still running. The provider blocks the Delete transition inside
// its status update while it holds the machine lock; the reconnect parks on
// that lock and commits as soon as the Delete releases it. The reconnect owns
// the machine from that point, so it must still be tracked once the delete's
// cleanup finishes. Measured before the fix: 6 evictions in 20 runs.
func TestDeleteCleanupKeepsMachineAdoptedByConcurrentReconnect(t *testing.T) {
	connID := uuid.Must(uuid.NewV4())
	sysID := core.Uuid(uuid.Must(uuid.NewV4()))
	log := newTestLogger(t)

	provider := &deleteBlockingProvider{
		lifecycleTestMockProvider: lifecycleTestMockProvider{
			k8sContext: models.K8sContext{ID: connID.String(), Name: "test-k8s", ConnectionID: connID.String()},
		},
		entered: make(chan struct{}),
		release: make(chan struct{}),
	}

	tracker := &machines.ConnectionToStateMachineInstanceTracker{ConnectToInstanceMap: make(map[core.Uuid]*machines.StateMachine)}
	h := &Handler{
		ConnectionToStateMachineInstanceTracker: tracker,
		config:                                  &models.HandlerConfig{EventBroadcaster: models.NewBroadcaster("test")},
		log:                                     log,
		SystemID:                                &sysID,
		MesheryCtrlsHelper:                      models.NewMesheryControllersHelper(log, controllers.OperatorDeploymentConfig{}, nil, nil, provider, &sysID),
	}

	sm, err := kubernetes.New(connID.String(), core.Uuid(connID), log)
	if err != nil {
		t.Fatal(err)
	}
	sm.Provider = provider
	// Keep the PRODUCTION transition table (DELETED must accept Connect); swap
	// only the actions so no cluster work runs.
	del := sm.States[machines.DELETED]
	sm.States[machines.DELETED] = *del.RegisterAction(&successfulLifecycleAction{})
	con := sm.States[machines.CONNECTED]
	sm.States[machines.CONNECTED] = *con.RegisterAction(&successfulLifecycleAction{})
	sm.CurrentState = machines.CONNECTED

	user := &models.User{ID: uuid.Must(uuid.NewV4())}
	ctx := context.WithValue(context.Background(), models.UserCtxKey, user)
	ctx = context.WithValue(ctx, models.SystemIDKey, &sysID)
	ctx = context.WithValue(ctx, models.TokenCtxKey, "test-token")
	machineCtx := &kubernetes.MachineCtx{K8sContext: provider.k8sContext, ActionMutex: &sync.Mutex{}}
	if _, err := sm.Start(ctx, machineCtx, log, func(c context.Context, mc interface{}, l logger.Handler) (interface{}, *events.Event, error) {
		return mc, nil, nil
	}); err != nil {
		t.Fatal(err)
	}
	tracker.Add(core.Uuid(connID), sm)

	// 1. DELETE, blocked inside the transition while holding the machine lock.
	reqDel := httptest.NewRequest(http.MethodDelete, "/api/system/kubernetes/contexts/"+connID.String(), nil)
	reqDel = mux.SetURLVars(reqDel, map[string]string{"id": connID.String()})
	cleanupDone := make(chan struct{})
	dctx := context.WithValue(reqDel.Context(), models.TokenCtxKey, "test-token")
	dctx = context.WithValue(dctx, models.UserCtxKey, user)
	dctx = context.WithValue(dctx, models.SystemIDKey, &sysID)
	dctx = context.WithValue(dctx, trackerCleanupDoneKey, cleanupDone)
	reqDel = reqDel.WithContext(dctx)
	recDel := httptest.NewRecorder()
	delDone := make(chan struct{})
	go func() { h.DeleteContext(recDel, reqDel, nil, user, provider); close(delDone) }()

	select {
	case <-provider.entered:
	case <-time.After(10 * time.Second):
		t.Fatal("delete transition never reached the provider status update")
	}

	// 2. RECONNECT parks on the machine lock while the Delete still holds it.
	reqRec := httptest.NewRequest(http.MethodPut, "/api/integrations/connections/"+connID.String(), strings.NewReader(`{"status":"connected"}`))
	reqRec = mux.SetURLVars(reqRec, map[string]string{"connectionId": connID.String()})
	rctx := context.WithValue(reqRec.Context(), models.TokenCtxKey, "test-token")
	rctx = context.WithValue(rctx, models.UserCtxKey, user)
	rctx = context.WithValue(rctx, models.SystemIDKey, &sysID)
	reqRec = reqRec.WithContext(rctx)
	recRec := httptest.NewRecorder()
	recDone := make(chan struct{})
	go func() { h.UpdateConnectionById(recRec, reqRec, nil, user, provider); close(recDone) }()

	time.Sleep(150 * time.Millisecond) // let the reconnect reach sm.mx.Lock()
	close(provider.release)

	<-delDone
	<-recDone
	select {
	case <-cleanupDone:
	case <-time.After(10 * time.Second):
		t.Fatal("delete cleanup goroutine never finished")
	}

	state := sm.GetCurrentState()
	_, tracked := tracker.Get(core.Uuid(connID))
	t.Logf("delete=%d reconnect=%d finalState=%s tracked=%v", recDel.Code, recRec.Code, state, tracked)

	if recRec.Code != http.StatusOK || state != machines.CONNECTED {
		t.Fatalf("setup did not produce the race: reconnect status=%d state=%s", recRec.Code, state)
	}
	if !tracked {
		t.Fatalf("delete cleanup left a live CONNECTED connection untracked: the reconnect adopted this StateMachine, so the next request will build a second one for the same connection")
	}
}
