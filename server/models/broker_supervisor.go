package models

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/schemas/models/core"
)

// An operator-mode broker connection is built once, from the broker endpoint
// and the NATS token (secret meshery-nats-auth) as they were at that moment.
// The NATS client reconnects forever to the URLs it was given, and treats an
// authorization rejection as terminal - so when the broker is recreated (the
// operator reinstalled, the namespace rebuilt) with a new endpoint or a new
// token, the connection never recovers, and MeshSync data stops until the
// connection's state machine re-initializes, which in practice meant a server
// restart.
//
// The supervisor closes that gap. While a connection is up it does nothing;
// once it has been down for brokerStaleAfter it re-resolves the endpoint and
// token, and rebuilds the data handler in place when either changed. A broker
// that is merely down - same endpoint, same token - is left to the client's own
// reconnect loop. Re-resolving only after a sustained outage was chosen over
// watching the Broker resource and the auth Secret: it needs no watch RBAC and
// holds no long-lived watch stream per remote cluster, and it costs API calls
// only while the connection is already broken.

const (
	defaultBrokerRecheckInterval = 15 * time.Second
	defaultBrokerStaleAfter      = 45 * time.Second
	// defaultBrokerConnectWait bounds how long a rebuild waits for its
	// replacement handler to report connected before giving up on it and
	// keeping the old handler. It stays well under staleAfter so the
	// replacement's own supervisor cannot start a rebuild of its own while
	// the wait is running.
	defaultBrokerConnectWait = 10 * time.Second
	brokerReachTimeout       = 5 * time.Second
)

// waitBrokerConnected reports whether handler reports connected within wait.
// It polls without holding meshsyncMutex: the NATS client connects
// asynchronously, so attachment - or a nil return from
// MeshsyncDataHandler.Run - says nothing about connectivity.
func waitBrokerConnected(handler *MeshsyncDataHandler, wait time.Duration) bool {
	if handler.IsConnected() {
		return true
	}
	if wait <= 0 {
		return false
	}
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.NewTimer(wait)
	defer deadline.Stop()
	for {
		select {
		case <-deadline.C:
			return handler.IsConnected()
		case <-ticker.C:
			if handler.IsConnected() {
				return true
			}
		}
	}
}

// brokerAnswers reports whether endpoint accepts a TCP connection.
func brokerAnswers(endpoint string) bool {
	conn, err := net.DialTimeout("tcp", endpoint, brokerReachTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// brokerTarget is what one broker connection was built from. It also owns the
// lifetime of that connection's supervisor.
type brokerTarget struct {
	controller controllers.IMesheryController
	endpoint   string
	token      string

	// reattach rebuilds the data handler with the inputs the replaced one was
	// built with. Called with meshsyncMutex held.
	reattach func()
	// Event metadata for a rebuild.
	k8scontext K8sContext
	userID     core.Uuid

	stop     chan struct{}
	stopOnce sync.Once
}

func newBrokerTarget(controller controllers.IMesheryController, endpoint, token string) *brokerTarget {
	return &brokerTarget{
		controller: controller,
		endpoint:   endpoint,
		token:      token,
		stop:       make(chan struct{}),
	}
}

// halt ends the supervisor that watches this target. Idempotent.
func (t *brokerTarget) halt() {
	t.stopOnce.Do(func() { close(t.stop) })
}

// replacedBy reports whether a freshly resolved endpoint and token describe a
// different broker than the one this connection was built for. An empty token
// means the operator has not provisioned the new broker yet, so there is
// nothing to rebuild against. An empty endpoint means none could be resolved,
// which is not evidence that it changed - unless a managed port-forward is
// carrying the connection: then the forward's local address is the transport
// the rebuild connects through, so a new token alone is a different broker.
// That mirrors the attach path, which only gives up when there is neither a
// published endpoint nor a forward. The resolved endpoint is the first
// published candidate that answers, falling back to the internal one when none
// does, so a down broker resolves differently from the one the connection was
// built for without having moved: a different endpoint counts only if it
// answers.
func (t *brokerTarget) replacedBy(endpoint, token string, answers func(string) bool, forwardUp bool) bool {
	if token == "" {
		return false
	}
	if token != t.token {
		return endpoint != "" || forwardUp
	}
	return endpoint != "" && endpoint != t.endpoint && answers(endpoint)
}

func (mch *MesheryControllersHelper) brokerTimings() (interval, staleAfter time.Duration) {
	interval, staleAfter = mch.brokerRecheckInterval, mch.brokerStaleAfter
	if interval <= 0 {
		interval = defaultBrokerRecheckInterval
	}
	if staleAfter <= 0 {
		staleAfter = defaultBrokerStaleAfter
	}
	return interval, staleAfter
}

// brokerConnectWaitTimeout is how long a rebuild waits for its replacement
// handler to report connected before giving up on it.
func (mch *MesheryControllersHelper) brokerConnectWaitTimeout() time.Duration {
	if mch.brokerConnectWait <= 0 {
		return defaultBrokerConnectWait
	}
	return mch.brokerConnectWait
}

// superviseBroker records target as the attached connection's and starts its
// supervisor. The caller holds meshsyncMutex.
func (mch *MesheryControllersHelper) superviseBroker(target *brokerTarget, k8scontext K8sContext, userID core.Uuid, reattach func()) {
	mch.stopBrokerSupervisorLocked()
	target.reattach = reattach
	target.k8scontext = k8scontext
	target.userID = userID
	mch.brokerTarget = target

	interval, _ := mch.brokerTimings()
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		var disconnectedSince time.Time
		for {
			select {
			case <-target.stop:
				return
			case now := <-ticker.C:
				if mch.recheckBroker(target, now, &disconnectedSince) {
					return
				}
			}
		}
	}()
}

// stopBrokerSupervisorLocked ends the attached connection's supervisor, if
// any. The caller holds meshsyncMutex. It does not wait for the supervisor to
// exit - the supervisor takes meshsyncMutex itself - but a stopped supervisor
// never acts: every action is re-checked under the mutex against brokerTarget.
func (mch *MesheryControllersHelper) stopBrokerSupervisorLocked() {
	if mch.brokerTarget != nil {
		mch.brokerTarget.halt()
		mch.brokerTarget = nil
	}
}

// recheckBroker is one supervisor tick. disconnectedSince carries how long the
// connection has been down across ticks. It reports true once this supervisor
// is finished: its connection was replaced, removed, or just rebuilt (which
// starts a supervisor of its own). A rebuild that attaches nothing - or whose
// replacement never connects - leaves the old handler and this supervisor in
// place, so the next tick tries again.
func (mch *MesheryControllersHelper) recheckBroker(target *brokerTarget, now time.Time, disconnectedSince *time.Time) (finished bool) {
	mch.meshsyncMutex.Lock()
	current, handler := mch.brokerTarget, mch.ctxMeshsyncDataHandler
	mch.meshsyncMutex.Unlock()
	if current != target || handler == nil {
		return true
	}
	if handler.IsConnected() {
		*disconnectedSince = time.Time{}
		return false
	}
	if disconnectedSince.IsZero() {
		*disconnectedSince = now
		return false
	}
	if _, staleAfter := mch.brokerTimings(); now.Sub(*disconnectedSince) < staleAfter {
		return false
	}

	// Resolved without the mutex: these are Kubernetes API calls, and status
	// readers take the same mutex.
	endpoint, _ := target.controller.GetPublicEndpoint()
	token, err := brokerAuthToken(target.controller)
	if err != nil {
		mch.log.Warn(err)
	}
	// RemoveMeshSyncDataHandler clears the managed port-forward under
	// meshsyncMutex, so read it under the same mutex. GetBrokerPortForwardAddr
	// itself takes no lock, because callers may already hold the mutex.
	mch.meshsyncMutex.Lock()
	hasForward := mch.GetBrokerPortForwardAddr() != ""
	mch.meshsyncMutex.Unlock()
	if !target.replacedBy(endpoint, token, mch.brokerReachable, hasForward) {
		return false
	}

	mch.meshsyncMutex.Lock()
	if mch.brokerTarget != target || mch.ctxMeshsyncDataHandler == nil {
		mch.meshsyncMutex.Unlock()
		return true
	}
	// The connection may have recovered while the endpoint and token were
	// being re-resolved; a live handler must not be stopped.
	if mch.ctxMeshsyncDataHandler.IsConnected() {
		*disconnectedSince = time.Time{}
		mch.meshsyncMutex.Unlock()
		return false
	}
	k8scontext := target.k8scontext
	mch.log.Info(fmt.Sprintf("Meshery Broker for Kubernetes context (%s) was recreated with a new endpoint or token; rebuilding the MeshSync data connection", k8scontext.ID))

	// The managed port-forward is kept: it re-resolves the NATS pod by label,
	// so it already follows a recreated broker.
	stale := mch.ctxMeshsyncDataHandler
	mch.ctxMeshsyncDataHandler = nil
	mch.brokerTarget = nil
	target.reattach()
	if mch.ctxMeshsyncDataHandler == nil {
		mch.ctxMeshsyncDataHandler = stale
		mch.brokerTarget = target
		mch.meshsyncMutex.Unlock()
		return false
	}
	replacement := mch.ctxMeshsyncDataHandler
	replacementTarget := mch.brokerTarget

	// The replacement attaches before it connects: Run returns once the
	// subscriptions are set up, without waiting for the broker connection, and
	// a token the broker rejects is terminal for the NATS client, so the
	// replacement may never connect. Keep the old handler until the
	// replacement reports connected, waiting without holding meshsyncMutex:
	// status readers and the replacement's own supervisor take it.
	wait := mch.brokerConnectWaitTimeout()
	mch.meshsyncMutex.Unlock()
	connected := waitBrokerConnected(replacement, wait)
	mch.meshsyncMutex.Lock()
	defer mch.meshsyncMutex.Unlock()

	if mch.ctxMeshsyncDataHandler != replacement {
		// Someone else removed or replaced the handler while the rebuild
		// waited: the old handler has no owner left, so stop it rather than
		// leave it reconnecting in the background.
		stale.Stop()
		return true
	}
	if !connected {
		// The replacement never connected - the broker rejected its token, or
		// its endpoint never came up. Discard it and keep the old handler and
		// its supervisor, so the next tick retries the rebuild instead of the
		// connection staying down until a restart or another token change.
		mch.log.Warnf("Rebuilt MeshSync data connection for Kubernetes context (%s) did not connect within %s; keeping the previous connection", k8scontext.ID, wait)
		replacement.Stop()
		if replacementTarget != nil {
			replacementTarget.halt()
		}
		mch.ctxMeshsyncDataHandler = stale
		mch.brokerTarget = target
		return false
	}
	// Only a rebuild that actually replaced the handler announces itself:
	// announcing before the reattach repeats the event on every retry tick
	// while the rebuild keeps failing, and resetting the flag then would
	// re-emit the connected event for the stale handler.
	mch.emitEvent("Meshery Broker changed; reconnecting MeshSync", events.Informational, map[string]any{
		"k8sContextID":   k8scontext.ID,
		"k8sContextName": k8scontext.Name,
		"connectionID":   k8scontext.ConnectionID,
		"brokerEndpoint": endpoint,
	}, target.userID)
	mch.meshsyncConnectedEventEmitted.Store(false)
	// The replacement is connected, so emit its connected event now instead of
	// waiting for a later reconcile call to notice it.
	mch.emitMeshsyncConnectedEventLocked(k8scontext, target.userID)
	target.halt()
	stale.Stop()
	return true
}
