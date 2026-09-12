package models

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/meshery/schemas/models/core"

	"maps"

	"github.com/gofrs/uuid"
	mesheryutils "github.com/meshery/meshery/server/helpers/utils"
	"github.com/meshery/meshery/server/models/connections"
	"github.com/meshery/meshkit/broker"
	channelBroker "github.com/meshery/meshkit/broker/channel"
	"github.com/meshery/meshkit/broker/nats"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/controllers"
	"github.com/meshery/meshkit/models/events"
	"github.com/meshery/meshkit/utils"
	mesherykube "github.com/meshery/meshkit/utils/kubernetes"
	libmeshsync "github.com/meshery/meshsync/pkg/lib/meshsync"
	controllersconfig "github.com/meshery/schemas/models/v1alpha1/controllers_config"
	"github.com/spf13/viper"
)

const (
	ChartRepo                     = "https://meshery.github.io/meshery.io/charts"
	MesheryServerBrokerConnection = "meshery-server"
)

type MesheryControllerStatusAndVersion struct {
	Status  controllers.MesheryControllerStatus
	Version string
}

type MesheryController int

const (
	MesheryBroker MesheryController = iota
	Meshsync
	MesheryOperator
)

// MesheryControllers is every controller Meshery manages on a Kubernetes
// connection, in enum order. It is the authoritative set: AddCtxControllerHandlers
// attaches a handler for each of these, so a caller that must account for every
// controller - whether or not a handler could be attached - iterates this rather
// than the handler map, which is missing entries precisely when something went
// wrong.
var MesheryControllers = []MesheryController{MesheryBroker, Meshsync, MesheryOperator}

type MesheryControllersHelper struct {
	// contextID is the Kubernetes context this helper serves.
	//
	// KNOWN BUG - nothing ever assigns it, so it is permanently "". The one
	// reader is UpdateOperatorsStatusMap, whose ot.IsUndeployed(mch.contextID)
	// therefore always probes the empty-string key and never matches a real
	// context: an operator explicitly undeployed for a context is not reported
	// as Undeployed there, and its live status is read from the handler instead.
	// Assigning it would silently change that operator-undeploy behavior, so it
	// is deliberately left alone here and tracked as a follow-up. Code that
	// needs a context identifier takes one from its caller (the K8sContext is
	// available at every call site) rather than reading this field.
	contextID string
	//  controller handlers for a particular context
	// this will be used as the source of truth
	ctxControllerHandlers map[MesheryController]controllers.IMesheryController

	// operator status for a particular context
	ctxOperatorStatus controllers.MesheryControllerStatus

	// meshsync data handler for a particular context
	ctxMeshsyncDataHandler *MeshsyncDataHandler

	// meshsyncConnectedEventEmitted tracks whether we already published the
	// "MeshSync connected in <mode> mode" event for the current data-handler
	// session. AddMeshsyncDataHandlers can be invoked concurrently while the
	// handler stays attached; atomic avoids a data race and re-broadcast spam.
	meshsyncConnectedEventEmitted atomic.Bool

	// brokerPortForward is the self-healing port-forward to the NATS pod used when
	// Meshery runs out-of-cluster and the broker is only reachable on a ClusterIP.
	// nil when running in-cluster or when managed forwarding is disabled.
	brokerPortForward *mesherykube.PortForwarder

	log          logger.Handler
	oprDepConfig controllers.OperatorDeploymentConfig
	dbHandler    *database.Handler

	meshsyncDeploymentMode connections.MeshsyncDeploymentMode

	// controllersConfig is the resolved (merged, explicitly-set) Meshery
	// Operator / MeshSync / Broker configuration for the context this helper
	// serves: per-connection override merged over the server-wide defaults.
	// Set alongside the deployment mode when a connection connects or its
	// configuration changes; consumed by the embedded meshsync run options.
	controllersConfig *controllersconfig.MesheryControllersConfig

	// attachedOperatorChartVersion is the Meshery Operator Helm chart version the
	// currently attached operator controller handler was constructed with. The
	// meshkit handler captures its deployment config at construction, so this is
	// the only record of what a later `operator.version` change has to differ
	// from before a redeploy is warranted. "" until handlers are attached.
	attachedOperatorChartVersion string

	// chartVersions lists the versions a Helm repository publishes for a chart.
	// Injected rather than called directly so tests resolve against a fixed
	// catalogue instead of the network. Never nil once constructed.
	chartVersions func(repo, chart string) ([]string, error)

	// reattachControllerHandlers re-runs handler attachment ahead of an install,
	// which is what re-resolves the chart version. It is AddCtxControllerHandlers
	// and is injected for the same reason chartVersions is: the state it produces
	// is otherwise unreachable from a test.
	//
	// Specifically, AddCtxControllerHandlers clears operatorChartError exactly
	// where it resolves a version, so it can never return having both succeeded
	// and left a refusal standing. That is the state operatorInstallTarget guards
	// the install sites against, and without this seam the guard is unreachable -
	// so nothing would notice an install site that stopped consulting it, which is
	// how this gap survived once already. Never nil once constructed.
	reattachControllerHandlers func(K8sContext)

	// event broadcasting dependencies
	eventBroadcaster *Broadcast
	provider         Provider
	systemID         *core.Uuid

	// lastOperatorError records the most recent failure to set up the operator for
	// this context — an unreadable kubeconfig, a failed Kubernetes client, or a
	// failed Deploy — so the connection-diagnostics API can surface *why* the
	// operator (and therefore MeshSync/broker) isn't running, not just that it
	// isn't. nil once setup succeeds. Guarded by opErrMu because it is written
	// from the connect goroutine and read from the diagnostics HTTP handler.
	opErrMu           sync.RWMutex
	lastOperatorError error

	// operatorChartError records that no publishable Meshery Operator chart
	// version could be resolved for this context. It is narrower than
	// lastOperatorError, which also carries deploy and client failures: this one
	// answers exactly "may Meshery install the operator right now?", and it is
	// what the install call sites refuse on. Observation is unaffected — meshkit's
	// operator handler reads the deployment config only in Deploy and Undeploy,
	// never in GetStatus or GetVersion — so the handler stays attached and keeps
	// reporting a healthy operator's status and image tag. Guarded by opErrMu.
	operatorChartError error
}

// setOperatorError records (or clears, when err is nil) the latest operator
// setup failure for this context. See lastOperatorError.
func (mch *MesheryControllersHelper) setOperatorError(err error) {
	mch.opErrMu.Lock()
	mch.lastOperatorError = err
	mch.opErrMu.Unlock()
}

// setOperatorErrorIfUnset records err only when no operator failure is already
// recorded for this context.
//
// It exists for the one diagnostic that is strictly less informative than
// anything it could replace: "no operator handler is attached" is the
// *consequence* of an unreadable kubeconfig or a failed Kubernetes client, both
// of which AddCtxControllerHandlers has already recorded by name. Letting the
// consequence overwrite the cause meant tearing down such a connection replaced
// the actionable diagnostic with a generic one.
func (mch *MesheryControllersHelper) setOperatorErrorIfUnset(err error) {
	mch.opErrMu.Lock()
	if mch.lastOperatorError == nil {
		mch.lastOperatorError = err
	}
	mch.opErrMu.Unlock()
}

// GetOperatorError returns the most recent operator setup failure for this
// context, or nil if operator setup last succeeded.
func (mch *MesheryControllersHelper) GetOperatorError() error {
	mch.opErrMu.RLock()
	defer mch.opErrMu.RUnlock()
	return mch.lastOperatorError
}

// setOperatorChartError records (or clears, when err is nil) the failure to
// resolve a publishable operator chart version. See operatorChartError.
func (mch *MesheryControllersHelper) setOperatorChartError(err error) {
	mch.opErrMu.Lock()
	mch.operatorChartError = err
	mch.opErrMu.Unlock()
}

// GetOperatorChartError returns why Meshery may not install the Meshery
// Operator on this context, or nil when a chart version was resolved.
func (mch *MesheryControllersHelper) GetOperatorChartError() error {
	mch.opErrMu.RLock()
	defer mch.opErrMu.RUnlock()
	return mch.operatorChartError
}

func (mch *MesheryControllersHelper) GetControllerHandlersForEachContext() map[MesheryController]controllers.IMesheryController {
	return mch.ctxControllerHandlers
}

func (mch *MesheryControllersHelper) GetMeshSyncDataHandlersForEachContext() *MeshsyncDataHandler {
	return mch.ctxMeshsyncDataHandler
}

func (mch *MesheryControllersHelper) GetOperatorsStatusMap() controllers.MesheryControllerStatus {
	return mch.ctxOperatorStatus
}

func NewMesheryControllersHelper(
	log logger.Handler,
	operatorDepConfig controllers.OperatorDeploymentConfig,
	dbHandler *database.Handler,
	eventBroadcaster *Broadcast,
	provider Provider,
	systemID *core.Uuid,
) *MesheryControllersHelper {
	mch := &MesheryControllersHelper{
		ctxControllerHandlers: make(map[MesheryController]controllers.IMesheryController),
		log:                   log,
		oprDepConfig:          operatorDepConfig,
		ctxOperatorStatus:     controllers.Unknown,
		// The nil check is performed for the ctxMeshsyncDataHandler and if it is nil, then a new dataHandler for the context is assigned.
		// The presence of a handler for a context in a map indicate that the meshsync data for that context is properly being handled.
		// Resetting this value results in again subscribing to the Broker.
		ctxMeshsyncDataHandler: nil,
		dbHandler:              dbHandler,
		chartVersions:          mesheryutils.PublishedChartVersions,
		meshsyncDeploymentMode: connections.MeshsyncDeploymentModeOperator,
		eventBroadcaster:       eventBroadcaster,
		provider:               provider,
		systemID:               systemID,
	}
	mch.reattachControllerHandlers = func(ctx K8sContext) { mch.AddCtxControllerHandlers(ctx) }
	return mch
}

// reattach re-runs handler attachment for ctx. See reattachControllerHandlers.
func (mch *MesheryControllersHelper) reattach(ctx K8sContext) {
	if mch.reattachControllerHandlers == nil {
		mch.AddCtxControllerHandlers(ctx)
		return
	}
	mch.reattachControllerHandlers(ctx)
}

func (mch *MesheryControllersHelper) SetMeshsyncDeploymentMode(value connections.MeshsyncDeploymentMode) *MesheryControllersHelper {
	mch.meshsyncDeploymentMode = value
	return mch
}

func (mch *MesheryControllersHelper) GetMeshsyncDeploymentMode() connections.MeshsyncDeploymentMode {
	return mch.meshsyncDeploymentMode
}

// GetBrokerPortForwardAddr returns the local address of the managed broker
// port-forward when one is active (out-of-cluster Meshery), else "". Used to
// surface how Meshery reaches the broker.
func (mch *MesheryControllersHelper) GetBrokerPortForwardAddr() string {
	if mch.brokerPortForward == nil {
		return ""
	}
	return mch.brokerPortForward.LocalAddr()
}

// SetControllersConfig stashes the resolved controllers configuration for the
// context this helper serves. Chainable, mirroring SetMeshsyncDeploymentMode.
//
// Call it BEFORE AddCtxControllerHandlers: the operator controller handler is
// constructed with the Helm chart version this document resolves to
// (`operator.version`), and the meshkit handler captures that at construction.
// On an already-attached helper, ReconcileOperatorChartVersion is what turns a
// changed version into a redeploy.
func (mch *MesheryControllersHelper) SetControllersConfig(value *controllersconfig.MesheryControllersConfig) *MesheryControllersHelper {
	mch.controllersConfig = value
	return mch
}

// ResolveControllersConfigForConnection resolves the layered controllers
// configuration for a connection's metadata against the server-wide defaults
// persisted in this server's database. It returns the merged
// (explicitly-set) document and the fully-resolved effective document.
//
// A malformed per-connection override invalidates only the override layer:
// it is logged and treated as absent so the Settings defaults still apply
// (the per-connection GET endpoint surfaces the parse error to the user).
// A non-nil error therefore always means the defaults store itself failed,
// in which case no resolution is returned.
func (mch *MesheryControllersHelper) ResolveControllersConfigForConnection(metadata core.Map) (merged, effective *controllersconfig.MesheryControllersConfig, err error) {
	override, overrideErr := connections.ControllersConfigFromMetadata(metadata)
	if overrideErr != nil {
		if mch.log != nil {
			mch.log.Error(overrideErr)
		}
		override = nil
	}
	serverDefaults, err := GetControllersConfigDefaults(mch.dbHandler)
	if err != nil {
		return nil, nil, err
	}
	merged, effective = connections.ResolveControllersConfig(override, serverDefaults)
	return merged, effective, nil
}

// AddMeshsyncDataHandlers initializes Meshsync data handler for the contexts for whom it has not been
// initialized yet. Apart from updating the map, it also runs the handler after
// updating the map. The presence of a handler for a context in a map indicate that
// the meshsync data for that context is properly being handled
func (mch *MesheryControllersHelper) AddMeshsyncDataHandlers(ctx context.Context, k8scontext K8sContext, userID, mesheryInstanceID core.Uuid, provider Provider) *MesheryControllersHelper {
	// only checking those contexts whose MesheryControllers are active
	// go func(mch *MesheryControllersHelper) {

	ctxID := k8scontext.ID
	if mch.ctxMeshsyncDataHandler == nil {
		var brokerHandler broker.Handler
		var stopFunc func()

		switch mch.meshsyncDeploymentMode {
		case connections.MeshsyncDeploymentModeOperator:
			brokerHandler = mch.meshsyncDataHandlersNatsBroker(k8scontext, userID)
		case connections.MeshsyncDeploymentModeEmbedded:
			brokerHandler = channelBroker.NewChannelBrokerHandler()
			// use a standalone context here context.Background(), as
			// meshsync run must be stopped only when meshsync data handler is deregistered
			// and ctx which is passed from above, could be closed earlier
			stop, err := mch.meshsyncDataHandlersStartLibMeshsyncRun(context.Background(), brokerHandler, k8scontext, userID)
			if err != nil {
				mch.log.Error(err)
				mch.emitErrorEvent("Failed to start MeshSync library run", err, map[string]any{
					"k8sContextID":           ctxID,
					"k8sContextName":         k8scontext.Name,
					"connectionID":           k8scontext.ConnectionID,
					"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
				}, userID)
				return mch
			}
			stopFunc = stop
		default:
			mch.log.Warnf(
				"MesheryControllersHelper unsupported meshsyncDeploymentMode %s",
				mch.meshsyncDeploymentMode,
			)
			mch.emitWarningEvent("Unsupported MeshSync deployment mode", nil, map[string]any{
				"k8sContextID":           ctxID,
				"k8sContextName":         k8scontext.Name,
				"connectionID":           k8scontext.ConnectionID,
				"meshsyncDeploymentMode": string(mch.meshsyncDeploymentMode),
			}, userID)
			return mch
		}

		if brokerHandler == nil {
			mch.log.Warnf("MesheryControllersHelper::AddMeshsyncDataHandlers brokerHandler is nil")
			mch.emitWarningEvent("MeshSync data handler broker is nil", nil, map[string]any{
				"k8sContextID":   ctxID,
				"k8sContextName": k8scontext.Name,
				"connectionID":   k8scontext.ConnectionID,
			}, userID)
			return mch
		}
		token, _ := ctx.Value(TokenCtxKey).(string)
		msDataHandler := NewMeshsyncDataHandler(brokerHandler, *mch.dbHandler, mch.log, provider, userID, uuid.FromStringOrNil(k8scontext.ConnectionID), mesheryInstanceID, token, stopFunc)
		err := msDataHandler.Run()
		if err != nil {
			mch.log.Warn(err)
			mch.log.Info(fmt.Sprintf("Unable to connect MeshSync for Kubernetes context (%s) due to: %s", ctxID, err.Error()))
			mch.emitErrorEvent("Unable to connect MeshSync", err, map[string]any{
				"k8sContextID":   ctxID,
				"k8sContextName": k8scontext.Name,
				"connectionID":   k8scontext.ConnectionID,
			}, userID)
			return mch
		}
		mch.ctxMeshsyncDataHandler = msDataHandler
		mch.log.Info(fmt.Sprintf("MeshSync connected for Kubernetes context (%s)", ctxID))
	}

	// }(mch)

	// Emit the success event only when the data path is actually up. With the
	// self-healing broker connection, the handler may be attached but still
	// connecting in the background; in that case meshsyncDataHandlersNatsBroker has
	// already surfaced the "unreachable, retrying" warning + remediation, and a
	// later AddMeshsyncDataHandlers call (once IsConnected) emits the event once.
	// Deduplicate so repeated reconcile calls do not spam the same snackbar.
	// CompareAndSwap so concurrent AddMeshsyncDataHandlers callers emit once.
	if mch.ctxMeshsyncDataHandler != nil &&
		mch.ctxMeshsyncDataHandler.IsConnected() &&
		mch.meshsyncConnectedEventEmitted.CompareAndSwap(false, true) {
		description := "MeshSync connected"
		if mch.meshsyncDeploymentMode != "" {
			description = fmt.Sprintf("MeshSync connected in %s mode", string(mch.meshsyncDeploymentMode))
		}
		mch.emitEvent(
			description,
			events.Informational,
			map[string]any{
				"k8sContextID":           k8scontext.ID,
				"k8sContextName":         k8scontext.Name,
				"connectionID":           k8scontext.ConnectionID,
				"meshsyncDeploymentMode": string(mch.meshsyncDeploymentMode),
			},
			userID,
		)
	}

	return mch
}

// BrokerUnreachableLongDescription explains why Meshery cannot reach the broker.
// It is shared by the "broker unreachable" event and the controller-diagnostics
// API so both surface identical wording.
const BrokerUnreachableLongDescription = "Meshery could not obtain a reachable endpoint for the Meshery Broker (NATS) in this cluster, so it cannot receive MeshSync data. This typically happens when the broker is only exposed on a cluster-internal (ClusterIP) address while Meshery runs outside the cluster."

// BrokerUnreachableRemediation is surfaced to the user (as an event's
// SuggestedRemediation and via the controller-diagnostics API) when Meshery
// cannot obtain a reachable Meshery Broker (NATS) endpoint. The most common
// cause is Meshery running outside the cluster while the broker is only exposed
// on a cluster-internal (ClusterIP) address — newer operators deploy NATS via
// the upstream Helm chart as ClusterIP-only, so the Broker publishes no external
// endpoint. Each entry renders as its own step in the UI.
var BrokerUnreachableRemediation = []string{
	"Port-forward the broker so Meshery can reach it, e.g.: kubectl port-forward -n meshery svc/meshery-nats 4222:4222",
	"Or expose the broker via a NodePort / LoadBalancer Service so it has an external endpoint.",
	"Or run Meshery inside the same cluster, so it can reach the broker's internal (ClusterIP) address.",
}

// BrokerTokenUnavailableLongDescription explains why Meshery has deferred the
// broker connection: the operator provisions NATS with token auth (secret
// meshery-nats-auth) and Meshery cannot authenticate until that token exists.
// This is normally a brief startup race that resolves on its own once the
// operator finishes provisioning the broker.
const BrokerTokenUnavailableLongDescription = "Meshery could not read the Meshery Broker (NATS) authentication token yet, so it deferred connecting to avoid an authorization failure. The Meshery Operator provisions the broker's auth token (secret meshery-nats-auth) as part of bringing NATS up; this is usually a brief startup race that clears on its own once provisioning completes, and Meshery retries automatically on the next connect cycle."

// BrokerTokenUnavailableRemediation is surfaced when the broker's auth token is
// not yet available. It self-heals in the common case; the steps below help only
// if the token never appears (a stuck or misconfigured operator).
var BrokerTokenUnavailableRemediation = []string{
	"No action is usually needed — Meshery retries automatically once the operator provisions the broker.",
	"If this persists, check the Meshery Operator is healthy: kubectl get pods -n meshery",
	"Confirm the broker auth secret exists: kubectl get secret meshery-nats-auth -n meshery",
}

func (mch *MesheryControllersHelper) meshsyncDataHandlersNatsBroker(
	k8scontext K8sContext,
	userID core.Uuid,
) broker.Handler {
	ctxID := k8scontext.ID
	controllerHandlers := mch.ctxControllerHandlers

	// The broker controller handler is only populated by AddCtxControllerHandlers,
	// which bails early (leaving ctxControllerHandlers empty) when it can't read
	// the kubeconfig or build a Kubernetes client for this context. Guard against a
	// nil handler here so that failure surfaces as an actionable diagnostic instead
	// of a nil-pointer panic when we go to read the broker's endpoint below.
	brokerController := controllerHandlers[MesheryBroker]
	if brokerController == nil {
		opErr := mch.GetOperatorError()
		mch.log.Warnf("Meshery Broker controller unavailable for Kubernetes context (%v); operator setup likely failed", ctxID)
		mch.emitWarningEvent("Meshery Broker controller unavailable", opErr, map[string]any{
			"k8sContextID":   ctxID,
			"k8sContextName": k8scontext.Name,
			"connectionID":   k8scontext.ConnectionID,
		}, userID)
		return nil
	}

	// Out-of-cluster Meshery can't reach a ClusterIP-only broker directly. Start a
	// self-healing managed port-forward to the NATS pod first — it selects the pod
	// by label, so it does not depend on the broker having published an endpoint
	// yet and can come up (and keep retrying) even when GetPublicEndpoint is still
	// empty right after deploy. In-cluster Meshery reaches the ClusterIP directly,
	// so this is a no-op there.
	forwardAddr := mch.ensureBrokerPortForward(brokerController)

	// brokerStatus := brokerController.GetStatus()
	// do something if broker is being deployed , maybe try again after sometime
	brokerEndpoint, err := brokerController.GetPublicEndpoint()
	if brokerEndpoint == "" && forwardAddr == "" {
		// Nothing to connect through: no published endpoint and no managed
		// port-forward to fall back on. Only now do we give up.
		if err != nil {
			mch.log.Warn(err)
		}
		mch.log.Info(
			fmt.Sprintf("Meshery Broker unreachable for Kubernetes context (%v)", ctxID),
		)
		// Surface a single, actionable event. The err (which carries the broker's
		// published internal/external endpoints) is attached via emitWarningEvent,
		// and the LongDescription/SuggestedRemediation keys are rendered by the
		// notification center so the user knows what to do next.
		mch.emitWarningEvent("Meshery Broker unreachable", err, map[string]any{
			"k8sContextID":         ctxID,
			"k8sContextName":       k8scontext.Name,
			"connectionID":         k8scontext.ConnectionID,
			"LongDescription":      BrokerUnreachableLongDescription,
			"SuggestedRemediation": BrokerUnreachableRemediation,
		}, userID)
		return nil
	}
	// The operator (>= 1.0.2) provisions NATS with token auth: without the token
	// the connection is rejected with an authorization violation even when the
	// endpoint is reachable, and the NATS client treats that rejection as terminal
	// (RetryOnFailedConnect only retries unreachability, never auth), so a handler
	// built with the wrong/empty token is permanently poisoned. The token lives in
	// the meshery-nats-auth secret, which the operator may create moments *after*
	// the broker pod becomes selectable by label — so proceeding here as soon as we
	// have a port-forward races that secret and can capture an empty token.
	//
	// Read the token and, if it isn't available yet, treat the broker as not-ready
	// and bail: leaving ctxMeshsyncDataHandler nil means the next connect cycle
	// retries, by which time the secret exists and we connect cleanly. We no longer
	// support tokenless brokers, so an empty token is always "not ready yet", never
	// a legitimately unauthenticated broker.
	var brokerToken string
	tokenProvider, ok := brokerController.(interface{ GetToken() (string, error) })
	if ok {
		token, terr := tokenProvider.GetToken()
		if terr != nil {
			mch.log.Warn(terr)
		}
		brokerToken = token
	}
	if brokerToken == "" {
		mch.log.Info(
			fmt.Sprintf("Meshery Broker auth token not available yet for Kubernetes context (%v); the operator is still provisioning it, retrying on the next connect cycle", ctxID),
		)
		mch.emitWarningEvent("Meshery Broker not ready", nil, map[string]any{
			"k8sContextID":         ctxID,
			"k8sContextName":       k8scontext.Name,
			"connectionID":         k8scontext.ConnectionID,
			"LongDescription":      BrokerTokenUnavailableLongDescription,
			"SuggestedRemediation": BrokerTokenUnavailableRemediation,
		}, userID)
		return nil
	}

	// Self-healing connection: hand NATS the (managed forward, if any) plus the
	// resolved endpoint and localhost / host.docker.internal on the same port,
	// retry on failed connect, and reconnect forever. NATS connects as soon as any
	// candidate becomes reachable (e.g. once the forward is up) and reconnects if
	// the broker or the tunnel drops — no restart or manual reconnect needed. When
	// no endpoint is published yet, the managed forward alone carries the
	// connection until the broker publishes one.
	var urls []string
	if forwardAddr != "" {
		urls = append(urls, forwardAddr)
	}
	if brokerEndpoint != "" {
		urls = append(urls, brokerConnectURLs(brokerEndpoint)...)
	}
	brokerHandler, err := nats.New(nats.Options{
		URLS:                 urls,
		ConnectionName:       MesheryServerBrokerConnection,
		Token:                brokerToken,
		ReconnectWait:        2 * time.Second,
		MaxReconnect:         -1, // reconnect indefinitely
		RetryOnFailedConnect: true,
	})
	if err != nil {
		// With RetryOnFailedConnect, nats.New only errors on misconfiguration, not
		// unreachability, so this is rare.
		mch.log.Warn(err)
		mch.log.Info(fmt.Sprintf("MeshSync not configured for Kubernetes context (%v) due to '%v'", ctxID, err.Error()))
		mch.emitWarningEvent("Failed to configure Meshery Broker connection", err, map[string]any{
			"k8sContextID":         ctxID,
			"k8sContextName":       k8scontext.Name,
			"connectionID":         k8scontext.ConnectionID,
			"brokerEndpoint":       brokerEndpoint,
			"LongDescription":      BrokerUnreachableLongDescription,
			"SuggestedRemediation": BrokerUnreachableRemediation,
		}, userID)
		return nil
	}

	if brokerHandler.IsConnected() {
		mch.log.Info(fmt.Sprintf("Connected to Meshery Broker (%v) for Kubernetes context (%v)", brokerEndpoint, ctxID))
	} else {
		// Not reachable yet. The handler keeps retrying in the background and will
		// connect (and MeshSync data will flow) as soon as the broker becomes
		// reachable — e.g. once a port-forward is up — with no restart. Surface the
		// remediation; the status/diagnostics flip to Connected on their own.
		mch.log.Info(fmt.Sprintf("Meshery Broker not reachable yet for Kubernetes context (%v); retrying in the background", ctxID))
		mch.emitWarningEvent("Meshery Broker unreachable", nil, map[string]any{
			"k8sContextID":         ctxID,
			"k8sContextName":       k8scontext.Name,
			"connectionID":         k8scontext.ConnectionID,
			"brokerEndpoint":       brokerEndpoint,
			"LongDescription":      BrokerUnreachableLongDescription,
			"SuggestedRemediation": BrokerUnreachableRemediation,
		}, userID)
	}
	return brokerHandler
}

// brokerConnectURLs returns the NATS server pool for the resolved broker
// endpoint. It always includes localhost and host.docker.internal on the same
// port so a port-forward that comes up later (or a Docker Desktop host mapping)
// is picked up by NATS's reconnect logic without re-resolving — the client
// simply skips unreachable entries. This is what makes the connection
// self-healing across a coming-and-going port-forward.
func brokerConnectURLs(endpoint string) []string {
	urls := []string{endpoint}
	_, port, err := net.SplitHostPort(endpoint)
	if err != nil || port == "" {
		return urls
	}
	for _, extra := range []string{
		net.JoinHostPort("127.0.0.1", port),
		net.JoinHostPort("host.docker.internal", port),
	} {
		if extra != endpoint {
			urls = append(urls, extra)
		}
	}
	return urls
}

// managedBrokerPortForwardEnabled reports whether Meshery should manage a
// port-forward to the broker. On by default; set
// MESHERY_MANAGED_BROKER_PORTFORWARD=false to disable (e.g. for a deterministic
// "unreachable" state in tests, or when a manual/other path is preferred).
func managedBrokerPortForwardEnabled() bool {
	return !strings.EqualFold(os.Getenv("MESHERY_MANAGED_BROKER_PORTFORWARD"), "false")
}

// runningInCluster reports whether Meshery itself runs inside a Kubernetes
// cluster, where the broker's ClusterIP is directly reachable and no port-forward
// is needed.
func runningInCluster() bool {
	return os.Getenv("KUBERNETES_SERVICE_HOST") != ""
}

// ensureBrokerPortForward starts (once) a self-healing managed port-forward to the
// broker's NATS pod when Meshery runs out-of-cluster, and returns its stable local
// address (or "" when not applicable). The forwarder is stored on the helper and
// torn down in RemoveMeshSyncDataHandler.
func (mch *MesheryControllersHelper) ensureBrokerPortForward(broker controllers.IMesheryController) string {
	if !managedBrokerPortForwardEnabled() || runningInCluster() || broker == nil {
		return ""
	}
	if mch.brokerPortForward != nil {
		return mch.brokerPortForward.LocalAddr()
	}
	provider, ok := broker.(interface {
		GetPortForwarder(logger.Handler) (*mesherykube.PortForwarder, error)
	})
	if !ok {
		return ""
	}
	pf, err := provider.GetPortForwarder(mch.log)
	if err != nil || pf == nil {
		if err != nil {
			mch.log.Warn(err)
		}
		return ""
	}
	pf.Start()
	mch.brokerPortForward = pf
	mch.log.Info(fmt.Sprintf("Started managed broker port-forward at %s", pf.LocalAddr()))
	return pf.LocalAddr()
}

// meshsyncDataHandlersStartLibMeshsyncRun starts the libmeshsync run for the given context.
// returns stop function to stop goroutine
func (mch *MesheryControllersHelper) meshsyncDataHandlersStartLibMeshsyncRun(
	ctx context.Context,
	brokerHandler broker.Handler,
	k8sContext K8sContext,
	userID core.Uuid,
) (func(), error) {
	kubeConfig, err := k8sContext.GenerateKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("MesheryControllersHelper::meshsyncDataHandlersStartLibMeshsyncRun error generating kubeconfig from context: %v", err)
	}

	cancelCtx, stopFunc := context.WithCancel(ctx)

	runOptions := []libmeshsync.OptionsSetter{
		libmeshsync.WithOutputMode("broker"),
		libmeshsync.WithBrokerHandler(brokerHandler),
		libmeshsync.WithKubeConfig(kubeConfig),
		libmeshsync.WithContext(cancelCtx),
	}
	// Apply the resolved controllers configuration to the in-process run.
	// Output filters are the embedded-mode knobs libmeshsync exposes today;
	// env-driven knobs (secret redaction, broker content dedup, debug
	// logging) follow the Meshery Server process environment in embedded
	// mode, and the watch-list is read from the MeshSync CR when the target
	// cluster has one. The run restarts whenever the configuration changes.
	if cfg := mch.controllersConfig; cfg != nil && cfg.Meshsync != nil {
		if len(cfg.Meshsync.OutputNamespaces) > 0 {
			runOptions = append(runOptions, libmeshsync.WithOnlyK8sNamespaces(cfg.Meshsync.OutputNamespaces...))
		}
		if len(cfg.Meshsync.OutputResources) > 0 {
			runOptions = append(runOptions, libmeshsync.WithOnlyK8sResources(cfg.Meshsync.OutputResources))
		}
	}

	go func() {
		if err := libmeshsync.Run(
			mch.log,
			runOptions...,
		); err != nil {
			meshsyncErr := fmt.Errorf("MesheryControllersHelper::meshsyncDataHandlersStartLibMeshsyncRun error running meshsync lib: %v", err)
			mch.log.Error(meshsyncErr)
			mch.emitErrorEvent("Error running MeshSync library", meshsyncErr, map[string]any{
				"k8sContextID":           k8sContext.ID,
				"k8sContextName":         k8sContext.Name,
				"connectionID":           k8sContext.ConnectionID,
				"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
			}, userID)
		}
	}()

	return stopFunc, nil
}

func (mch *MesheryControllersHelper) RemoveMeshSyncDataHandler(ctx context.Context, contextID string) {
	if mch.ctxMeshsyncDataHandler != nil {
		mch.log.Infof("MesheryControllersHelper::RemoveMeshSyncDataHandler for contextID = %s", contextID)
		mch.ctxMeshsyncDataHandler.Stop()
		mch.ctxMeshsyncDataHandler = nil
	}
	// Allow a fresh "MeshSync connected" event when a new handler is attached.
	mch.meshsyncConnectedEventEmitted.Store(false)
	// Tear down the managed broker port-forward alongside the data handler.
	if mch.brokerPortForward != nil {
		mch.brokerPortForward.Stop()
		mch.brokerPortForward = nil
	}
}

func (mch *MesheryControllersHelper) ResyncMeshsync(ctx context.Context) error {
	if mch.ctxMeshsyncDataHandler != nil {
		return mch.ctxMeshsyncDataHandler.Resync()
	}
	return nil
}

// AddCtxControllerHandlers attaches a MesheryController for each context if
// 1. the config is valid
// 2. if it is not already attached
func (mch *MesheryControllersHelper) AddCtxControllerHandlers(ctx K8sContext) *MesheryControllersHelper {
	// go func(mch *MesheryControllersHelper) {

	// resetting this value as a specific controller handler instance does not have any significance opposed to
	// a MeshsyncDataHandler instance where it signifies whether or not a listener is attached

	// Fresh setup attempt: clear any error from a previous attempt. It is re-set
	// below (and in DeployUndeployedOperators) if this attempt fails, so the
	// diagnostics API always reflects the latest operator setup outcome.
	//
	// operatorChartError is deliberately NOT cleared here. It is the refusal that
	// keeps an unresolved chart version away from Helm, and the two paths below
	// return without reaching resolution while leaving the previously attached
	// operator handler - which still carries that unresolved version - in place.
	// Clearing it up here therefore produced the one state the install guard
	// cannot see: a stale handler with no refusal against it. It is cleared on
	// the success path instead, so it lifts only when a resolution actually
	// succeeded.
	mch.setOperatorError(nil)

	cfg, err := ctx.GenerateKubeConfig()
	if err != nil {
		mch.setOperatorError(err)
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to generate kubeconfig", err, map[string]any{
			"k8sContextID":   ctx.ID,
			"k8sContextName": ctx.Name,
			"connectionID":   ctx.ConnectionID,
		}, uuid.Nil)
		return mch
	}

	client, err := mesherykube.New(cfg)
	// means that the config is invalid
	if err != nil {
		mch.setOperatorError(err)
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to create Kubernetes client", err, map[string]any{
			"k8sContextID":   ctx.ID,
			"k8sContextName": ctx.Name,
			"connectionID":   ctx.ConnectionID,
		}, uuid.Nil)
		return mch
	}

	// Broker and MeshSync handlers only observe what is already in the cluster,
	// so they need no chart version and are always attached.
	ctxHandlers := map[MesheryController]controllers.IMesheryController{
		MesheryBroker: controllers.NewMesheryBrokerHandler(client),
		Meshsync:      controllers.NewMeshsyncHandler(client),
	}

	// The operator handler is different: it captures its chart version at
	// construction and installs from it. Pin that version first - attaching a
	// handler carrying a version the repository does not publish only defers the
	// failure to the Helm install, where it surfaces as an opaque
	// chart-not-found error rather than as this connection's operator error.
	depConfig, substitution, err := mch.pinnedOperatorDeploymentConfig()
	if err != nil {
		// A chart version is what it takes to *install* the operator, not to
		// *observe* one: meshkit's handler consults the deployment config in
		// Deploy and Undeploy only, so GetStatus and GetVersion work without it.
		// The handler is therefore attached with the unresolved config and
		// installation alone is withheld, recorded on operatorChartError for the
		// install call sites to refuse on. Withholding the handler instead cost
		// the operator card its status and image tag - the very value the
		// troubleshooting guide tells users to read - on a cluster whose operator
		// was installed and healthy all along.
		mch.setOperatorChartError(err)
		mch.setOperatorError(err)
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to resolve the Meshery Operator chart version", err, map[string]any{
			"k8sContextID":                  ctx.ID,
			"k8sContextName":                ctx.Name,
			"connectionID":                  ctx.ConnectionID,
			"requestedOperatorChartVersion": mch.operatorDeploymentConfig().MesheryReleaseVersion,
		}, uuid.Nil)
		ctxHandlers[MesheryOperator] = controllers.NewMesheryOperatorHandler(client, mch.operatorDeploymentConfig())
		mch.ctxControllerHandlers = ctxHandlers
		// The attached handler carries no installable version, so the next
		// reconcile must not mistake a stale value for "already at the desired
		// version".
		mch.attachedOperatorChartVersion = ""
		return mch
	}
	if substitution != "" {
		mch.log.Warn(ErrOperatorChartSubstituted(substitution))
		mch.emitWarningEvent("Meshery Operator chart version adjusted", ErrOperatorChartSubstituted(substitution), map[string]any{
			"k8sContextID":         ctx.ID,
			"k8sContextName":       ctx.Name,
			"connectionID":         ctx.ConnectionID,
			"operatorChartVersion": depConfig.MesheryReleaseVersion,
		}, uuid.Nil)
	}

	ctxHandlers[MesheryOperator] = controllers.NewMesheryOperatorHandler(client, depConfig)
	mch.ctxControllerHandlers = ctxHandlers
	mch.attachedOperatorChartVersion = depConfig.MesheryReleaseVersion
	mch.setOperatorChartError(nil)

	// }(mch)
	return mch
}

// operatorDeploymentConfig resolves the Meshery Operator deployment
// configuration to attach a controller handler with: the server-wide
// configuration assembled once at boot (NewOperatorDeploymentConfig), with the
// Helm chart version replaced by this connection's resolved `operator.version`
// when a layer sets one.
//
// The chart version this returns is the one that was *asked for*; it is not yet
// known to exist. pinnedOperatorDeploymentConfig is what turns it into a
// version the chart repository actually publishes, and every path that reaches
// a cluster goes through that instead.
//
// Leaving `operator.version` unset at every layer yields the boot-time chart
// version verbatim. The operator handler captures this configuration at
// construction, so the resolved controllers configuration must be stashed
// (SetControllersConfig) before AddCtxControllerHandlers runs - see
// ReconcileOperatorChartVersion for the path that re-attaches when the version
// changes on a live connection.
func (mch *MesheryControllersHelper) operatorDeploymentConfig() controllers.OperatorDeploymentConfig {
	depConfig := mch.oprDepConfig
	if version := connections.OperatorChartVersionFromControllersConfig(mch.controllersConfig); version != "" {
		depConfig.MesheryReleaseVersion = version
	}
	return depConfig
}

// pinnedOperatorDeploymentConfig is operatorDeploymentConfig with the chart
// version pinned to one the chart repository actually publishes.
//
// This is the only deployment config that may reach Helm. The version
// operatorDeploymentConfig returns is either an explicit `operator.version` or
// the Meshery Server release stamped in at boot, and neither is guaranteed to
// name a published chart: chart publishing trails server releases, so a current
// server routinely asks for a chart that does not exist yet, and a server old
// enough to predate MinimumOperatorChartVersion asks for one that exists but
// cannot run. ResolveOperatorChartVersion settles both, failing loudly when the
// version was explicitly requested and substituting - with a reason - when it
// was merely derived.
//
// The returned reason is empty unless a substitution happened; callers log it
// and emit it so no correction is silent.
func (mch *MesheryControllersHelper) pinnedOperatorDeploymentConfig() (depConfig controllers.OperatorDeploymentConfig, reason string, err error) {
	depConfig = mch.operatorDeploymentConfig()

	source := OperatorChartVersionDerived
	if connections.OperatorChartVersionFromControllersConfig(mch.controllersConfig) != "" {
		source = OperatorChartVersionRequested
	}

	lister := mch.chartVersions
	if lister == nil {
		lister = mesheryutils.PublishedChartVersions
	}
	published, err := lister(depConfig.HelmChartRepo, OperatorChartName)
	if err != nil {
		return depConfig, "", err
	}

	version, reason, err := ResolveOperatorChartVersion(depConfig.HelmChartRepo, published, depConfig.MesheryReleaseVersion, source)
	if err != nil {
		return depConfig, "", err
	}
	depConfig.MesheryReleaseVersion = version
	return depConfig, reason, nil
}

// ReconcileOperatorChartVersion brings the operator running on this connection's
// cluster into line with the chart version its resolved controllers
// configuration asks for.
//
// The operator handler captures its chart version at construction, so a changed
// `operator.version` means re-attaching the handlers and re-running the Helm
// release. That release is applied with UpgradeIfInstalled, so an operator
// already present is upgraded in place rather than reinstalled.
//
// It reports the chart version now in force and whether a redeploy was
// performed. It no-ops - reporting redeployed=false - when:
//   - the connection is not in operator deployment mode (nothing is installed
//     into the cluster at all, so the field is inert and the UI says so),
//   - the operator is disabled server-wide or was explicitly undeployed for
//     this context, or
//   - the resolved chart version already matches the attached handler's.
func (mch *MesheryControllersHelper) ReconcileOperatorChartVersion(k8sctx K8sContext, ot *OperatorTracker) (chartVersion string, redeployed bool, err error) {
	requested := mch.operatorDeploymentConfig().MesheryReleaseVersion
	if mch.meshsyncDeploymentMode != connections.MeshsyncDeploymentModeOperator {
		return requested, false, nil
	}
	if ot == nil || ot.DisableOperator || ot.IsUndeployed(k8sctx.ID) {
		return requested, false, nil
	}

	// Compare pinned against pinned. attachedOperatorChartVersion records what
	// the handler was built with, which is always a published version, so
	// comparing the raw request against it would redeploy on every call
	// whenever the request is being substituted (an unpublished server release
	// or a below-floor one) - an endless upgrade loop against the cluster.
	pinned, _, err := mch.pinnedOperatorDeploymentConfig()
	if err != nil {
		mch.setOperatorError(err)
		return requested, false, ErrReconcileOperatorChartVersion(err)
	}
	desired := pinned.MesheryReleaseVersion
	if desired == mch.attachedOperatorChartVersion {
		return desired, false, nil
	}
	// Any substitution is logged and emitted by AddCtxControllerHandlers below,
	// which resolves against the same catalogue; reporting it here too would
	// duplicate every warning.

	// Re-attaching records the desired version as attached, which the guard
	// above then treats as "already reconciled". If anything below fails, that
	// leaves a failed upgrade permanently ineligible for retry - the next call
	// would short-circuit and never try again. Restore the previous value on
	// every failure path so a transient Helm error is retried rather than
	// silently becoming the resting state.
	//
	// Restoring is conditional on AddCtxControllerHandlers having left an
	// attached version behind: when it could not resolve an installable one it
	// cleared the value *for the same reason*, and putting the stale one back
	// would withhold operator lifecycle until the connection is rebuilt.
	previousChartVersion := mch.attachedOperatorChartVersion
	mch.reattach(k8sctx)
	if setupErr := mch.GetOperatorError(); setupErr != nil {
		if mch.attachedOperatorChartVersion != "" {
			mch.attachedOperatorChartVersion = previousChartVersion
		}
		return desired, false, ErrReconcileOperatorChartVersion(setupErr)
	}
	// Through the same guard as every other install, rather than reading the
	// handler out of the map directly. Reaching Deploy here was gated only by
	// GetOperatorError being nil after the re-attach above, which holds solely
	// because AddCtxControllerHandlers writes both errors together - a coupling
	// nothing enforces, and one that five other setOperatorError writers could
	// break into an unresolved version reaching Helm.
	operatorHandler, _, targetErr := mch.operatorInstallTarget(k8sctx.ID)
	if targetErr != nil {
		if mch.attachedOperatorChartVersion != "" {
			mch.attachedOperatorChartVersion = previousChartVersion
		}
		return desired, false, ErrReconcileOperatorChartVersion(targetErr)
	}
	// Deploy(false) is a no-op against an operator this handler undeployed and a
	// Helm upgrade against one that is installed, so an operator the user turned
	// off is not resurrected by a chart-version change.
	if deployErr := operatorHandler.Deploy(false); deployErr != nil {
		mch.attachedOperatorChartVersion = previousChartVersion
		mch.setOperatorError(deployErr)
		return desired, false, ErrReconcileOperatorChartVersion(deployErr)
	}
	return desired, true, nil
}

func (mch *MesheryControllersHelper) RemoveCtxControllerHandler(ctx context.Context, contextID string) {
	mch.ctxControllerHandlers = nil
}

// UpdateOperatorsStatusMap updates the status of MesheryOperator in all the contexts
// for whom MesheryControllers are attached
// should be called after AddCtxControllerHandlers
func (mch *MesheryControllersHelper) UpdateOperatorsStatusMap(ot *OperatorTracker) *MesheryControllersHelper {
	// go func(mch *MesheryControllersHelper) {
	if mch.meshsyncDeploymentMode != connections.MeshsyncDeploymentModeOperator {
		return mch
	}

	if ot.IsUndeployed(mch.contextID) {
		// this code is probably never reached as mch.contextID is never set
		mch.ctxOperatorStatus = controllers.Undeployed
	} else {
		if mch.ctxControllerHandlers != nil {
			operatorHandler, ok := mch.ctxControllerHandlers[MesheryOperator]
			if ok {
				mch.ctxOperatorStatus = operatorHandler.GetStatus()
			}
		}
	}

	// }(mch)
	return mch
}

type OperatorTracker struct {
	ctxIDtoDeploymentStatus map[string]bool
	mx                      sync.Mutex
	DisableOperator         bool
}

func NewOperatorTracker(disabled bool) *OperatorTracker {
	return &OperatorTracker{
		ctxIDtoDeploymentStatus: make(map[string]bool),
		mx:                      sync.Mutex{},
		DisableOperator:         disabled,
	}
}

func (ot *OperatorTracker) Undeployed(ctxID string, undeployed bool) {
	if ot.DisableOperator { //no-op when operator is disabled
		return
	}
	if ot.ctxIDtoDeploymentStatus == nil {
		ot.ctxIDtoDeploymentStatus = make(map[string]bool)
	}
	ot.ctxIDtoDeploymentStatus[ctxID] = undeployed
}

func (ot *OperatorTracker) IsUndeployed(ctxID string) bool {
	if ot.DisableOperator { //Return true everytime so that operators stay in undeployed state across all contexts
		return true
	}
	if ot.ctxIDtoDeploymentStatus == nil {
		ot.ctxIDtoDeploymentStatus = make(map[string]bool)
		return false
	}
	return ot.ctxIDtoDeploymentStatus[ctxID]
}

// attachedOperatorHandler returns the Meshery Operator controller handler for
// this context, or nil when none is attached.
//
// A nil handler is never a no-op a caller may pass over in silence: it means
// AddCtxControllerHandlers could not read the kubeconfig or build a Kubernetes
// client, so nothing was done to the cluster.
func (mch *MesheryControllersHelper) attachedOperatorHandler() controllers.IMesheryController {
	if mch.ctxControllerHandlers == nil {
		return nil
	}
	return mch.ctxControllerHandlers[MesheryOperator]
}

// operatorInstallTarget returns the handler to install the Meshery Operator
// through, or a structured error naming why installation may not proceed.
// attached reports whether a handler was present at all, which is the one
// condition callers report differently (see setOperatorErrorIfUnset).
//
// This is the single guard every install path goes through. There are exactly
// three, and all three call this: the FSM's DeployUndeployedOperators, the
// user-initiated SetOperatorDeployment, and ReconcileOperatorChartVersion. They
// are the only callers of Deploy on an operator handler in this repository, so
// none of them can drift into handing Helm an unpublishable chart version while
// the others refuse it.
//
// Installing is the one operator action that needs a chart version. The handler
// is attached for observation even when none could be resolved, so refusing
// here keeps the structured resolution failure in front of the user instead of
// the opaque chart-not-found error Helm would raise from an unpublished version.
func (mch *MesheryControllersHelper) operatorInstallTarget(contextID string) (handler controllers.IMesheryController, attached bool, err error) {
	operatorHandler := mch.attachedOperatorHandler()
	if operatorHandler == nil {
		return nil, false, ErrOperatorHandlerNotAttached(contextID)
	}
	if chartErr := mch.GetOperatorChartError(); chartErr != nil {
		return nil, true, chartErr
	}
	return operatorHandler, true, nil
}

// reportMissingOperatorHandler records and surfaces the absence of an operator
// handler for an action that needed one, so a lifecycle request that did
// nothing does not read to the user as one that succeeded.
func (mch *MesheryControllersHelper) reportMissingOperatorHandler(action, contextID string) {
	err := ErrOperatorHandlerNotAttached(contextID)
	mch.setOperatorErrorIfUnset(err)
	mch.log.Error(err)
	mch.emitErrorEvent(action, err, map[string]any{
		"k8sContextID":           contextID,
		"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
		"operatorStatus":         mch.ctxOperatorStatus,
	}, uuid.Nil)
}

// operatorStatusObserved reports whether Meshery has ever actually read this
// context's operator status, as opposed to still carrying the Unknown that
// NewMesheryControllersHelper seeds. It distinguishes "we saw an operator and
// can no longer act on it" from "we never reached this cluster at all".
func (mch *MesheryControllersHelper) operatorStatusObserved() bool {
	return mch.ctxOperatorStatus != controllers.Unknown
}

// SetOperatorDeployment applies a user-initiated Meshery Operator lifecycle
// request for k8sctx: deploy=true installs or upgrades it, deploy=false removes
// it. It is the counterpart to the FSM's DeployUndeployedOperators and
// UndeployDeployedOperators for callers that must be told the outcome - the
// GraphQL changeOperatorStatus mutation - and shares their guard rather than
// restating it.
//
// A latched chart-resolution failure is retried here rather than being final.
// operatorChartError is written only when handlers are attached, so a connect
// that landed during a chart-repository outage would otherwise refuse every
// later install for the life of that connection, no matter how long ago the
// repository came back - and the user is told, correctly, to confirm the
// repository is reachable and retry. Re-attaching is what makes that
// instruction true: clearing the latch alone would not, because the handler the
// failure path attached still carries the raw unresolved chart version, so the
// version and the latch have to be corrected together. AddCtxControllerHandlers
// does exactly that, and only a resolution that actually succeeds clears the
// latch; one that fails again replaces it with the fresh error, which the guard
// below then refuses on instead of falling through to a Deploy.
//
// The error is returned rather than emitted: the caller owns how a request it
// initiated is surfaced. A failure is still recorded for the diagnostics API.
func (mch *MesheryControllersHelper) SetOperatorDeployment(k8sctx K8sContext, deploy bool) error {
	if !deploy {
		operatorHandler := mch.attachedOperatorHandler()
		if operatorHandler == nil {
			err := ErrOperatorHandlerNotAttached(k8sctx.ID)
			mch.setOperatorErrorIfUnset(err)
			return err
		}
		if err := operatorHandler.Undeploy(); err != nil {
			mch.setOperatorError(err)
			return err
		}
		return nil
	}

	if mch.GetOperatorChartError() != nil {
		mch.reattach(k8sctx)
	}

	operatorHandler, attached, err := mch.operatorInstallTarget(k8sctx.ID)
	if err != nil {
		if attached {
			mch.setOperatorError(err)
		} else {
			mch.setOperatorErrorIfUnset(err)
		}
		return err
	}
	if err := operatorHandler.Deploy(true); err != nil {
		mch.setOperatorError(err)
		return err
	}
	return nil
}

// DeployUndeployedOperators looks at the status of Meshery Operator for each cluster and takes necessary action.
// it will deploy the operator only when it is in NotDeployed state
func (mch *MesheryControllersHelper) DeployUndeployedOperators(ot *OperatorTracker, contextID string) *MesheryControllersHelper {
	if ot.DisableOperator { //Return true everytime so that operators stay in undeployed state across all contexts
		return mch
	}
	if mch.meshsyncDeploymentMode != connections.MeshsyncDeploymentModeOperator {
		return mch
	}
	if mch.ctxOperatorStatus != controllers.NotDeployed {
		return mch
	}

	operatorHandler, attached, err := mch.operatorInstallTarget(contextID)
	if err != nil {
		if !attached {
			mch.reportMissingOperatorHandler("Failed to deploy Meshery Operator", contextID)
			return mch
		}
		// The chart-resolution failure is the same one AddCtxControllerHandlers
		// already emitted for this connection, so it is recorded and logged
		// rather than emitted a second time on every reconcile.
		mch.setOperatorError(err)
		mch.log.Error(err)
		return mch
	}

	if err := operatorHandler.Deploy(false); err != nil {
		mch.setOperatorError(err)
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to deploy Meshery Operator", err, map[string]any{
			"k8sContextID":           contextID,
			"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
			"operatorStatus":         mch.ctxOperatorStatus,
		}, uuid.Nil)
	}

	return mch
}

func (mch *MesheryControllersHelper) UndeployDeployedOperators(ot *OperatorTracker, contextID string) *MesheryControllersHelper {
	if mch.ctxOperatorStatus == controllers.Undeployed {
		return mch
	}

	// Unlike deploy, this is not gated on GetOperatorChartError. Removal is the
	// direction to attempt rather than refuse: refusing would leave the operator
	// running on a cluster the user asked to have it taken off.
	//
	// Attempted is all it is, though. meshkit's ApplyHelmChart downloads and
	// loads the chart archive before dispatching UNINSTALL just as it does for
	// INSTALL, from the same repository whose index could not be read - so in the
	// very case that motivates not refusing, Undeploy fails at chart download,
	// the operator stays on the cluster, and the user gets a Helm error. A
	// removal path that does not need the archive would have to come from
	// meshkit; until then this surfaces the failure rather than hiding it.
	operatorHandler := mch.attachedOperatorHandler()
	if operatorHandler == nil {
		// A connection whose cluster was never reached has no operator to
		// remove and never had one, so reporting a failed removal on its
		// teardown would be a second alert for a non-event - and would bury the
		// kubeconfig or Kubernetes-client diagnostic that explains it. Report
		// only when an operator was actually observed here at some point.
		if mch.operatorStatusObserved() {
			mch.reportMissingOperatorHandler("Failed to undeploy Meshery Operator", contextID)
		}
		return mch
	}

	if err := operatorHandler.Undeploy(); err != nil {
		mch.setOperatorError(err)
		mch.log.Error(err)
		mch.emitErrorEvent("Failed to undeploy Meshery Operator", err, map[string]any{
			"k8sContextID":           contextID,
			"meshsyncDeploymentMode": mch.meshsyncDeploymentMode,
			"operatorStatus":         mch.ctxOperatorStatus,
		}, uuid.Nil)
	}

	return mch
}

func NewOperatorDeploymentConfig(adapterTracker AdaptersTrackerInterface) controllers.OperatorDeploymentConfig {
	// The chart version is the Meshery Server release this binary was stamped
	// with, which is a *request*, not a promise: it is pinned to a version the
	// chart repository actually publishes by pinnedOperatorDeploymentConfig
	// before it ever reaches Helm.
	//
	// An unstamped build (a source run, or an edge image) has no release to
	// name, so it is left empty deliberately rather than resolved against the
	// GitHub releases API. The newest *GitHub release* is not the newest
	// *published chart* - charts are republished at server releases and trail
	// them - so asking GitHub yields a version that frequently does not exist
	// in the chart repository at all, on top of spending an unauthenticated,
	// rate-limited API call at boot. An empty version resolves to the newest
	// published chart, which is what was actually wanted.
	mesheryReleaseVersion := viper.GetString("BUILD")
	if !isPinnedChartVersion(mesheryReleaseVersion) {
		mesheryReleaseVersion = ""
	}

	return controllers.OperatorDeploymentConfig{
		MesheryReleaseVersion: mesheryReleaseVersion,
		GetHelmOverrides: func(delete bool) map[string]interface{} {
			return setOverrideValues(delete, adapterTracker)
		},
		HelmChartRepo: ChartRepo,
	}
}

// CheckLatestVersion takes in the current server version compares it with the target
// and returns the (isOutdated, latestVersion, error)
func CheckLatestVersion(serverVersion string) (*bool, string, error) {
	// Inform user of the latest release version
	versions, err := utils.GetLatestReleaseTagsSorted("meshery", "meshery")
	if err != nil {
		return nil, "", ErrCreateOperatorDeploymentConfig(err)
	}
	// Index only after the error check: a failed fetch returns a nil slice, and
	// indexing it panicked - taking down whichever request or boot step called
	// in - every time GitHub was unreachable or rate-limited.
	if len(versions) == 0 {
		return nil, "", ErrCreateOperatorDeploymentConfig(ErrNoMesheryReleasesFound())
	}
	latestVersion := versions[len(versions)-1]
	isOutdated := false
	// Compare current running Meshery server version to the latest available Meshery release on GitHub.
	if latestVersion != serverVersion {
		isOutdated = true
		return &isOutdated, latestVersion, nil
	}

	return &isOutdated, latestVersion, nil
}

// setOverrideValues detects the currently insalled adapters and sets appropriate
// overrides so as to not uninstall them. It also sets override values for
// operator so that it can be enabled or disabled depending on the need
func setOverrideValues(delete bool, adapterTracker AdaptersTrackerInterface) map[string]interface{} {
	installedAdapters := make([]string, 0)
	adapters := adapterTracker.GetAdapters(context.TODO())

	for _, adapter := range adapters {
		if adapter.Name != "" {
			installedAdapters = append(installedAdapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{}{
		"fullnameOverride": "meshery-operator",
		"meshery": map[string]interface{}{
			"enabled": false,
		},
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-operator": map[string]interface{}{
			"enabled": true,
		},
	}

	for _, adapter := range installedAdapters {
		if _, ok := overrideValues[adapter]; ok {
			overrideValues[adapter] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	if delete {
		overrideValues["meshery-operator"] = map[string]interface{}{
			"enabled": false,
		}
	}

	return overrideValues
}

// SetOverrideValuesForMesheryDeploy detects the currently installed adapters and sets appropriate
// overrides so as to not uninstall them, then sets the override for the given adapter based on install.
func SetOverrideValuesForMesheryDeploy(adapters []Adapter, adapter Adapter, install bool) map[string]interface{} {
	installedAdapters := make([]string, 0)
	for _, adapter := range adapters {
		if adapter.Name != "" {
			installedAdapters = append(installedAdapters, strings.Split(adapter.Location, ":")[0])
		}
	}

	overrideValues := map[string]interface{}{
		"meshery-istio": map[string]interface{}{
			"enabled": false,
		},
		"meshery-cilium": map[string]interface{}{
			"enabled": false,
		},
		"meshery-linkerd": map[string]interface{}{
			"enabled": false,
		},
		"meshery-consul": map[string]interface{}{
			"enabled": false,
		},
		"meshery-kuma": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nsm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-nginx-sm": map[string]interface{}{
			"enabled": false,
		},
		"meshery-traefik-mesh": map[string]interface{}{
			"enabled": false,
		},
		"meshery-app-mesh": map[string]interface{}{
			"enabled": false,
		},
	}

	for _, adapter := range installedAdapters {
		if _, ok := overrideValues[adapter]; ok {
			overrideValues[adapter] = map[string]interface{}{
				"enabled": true,
			}
		}
	}

	// based on deploy/undeploy action change the status of adapter override
	if _, ok := overrideValues[strings.Split(adapter.Location, ":")[0]]; ok {
		overrideValues[strings.Split(adapter.Location, ":")[0]] = map[string]interface{}{
			"enabled": install,
		}
	}

	return overrideValues
}

// General helper method to emit events for system-level operations
func (mch *MesheryControllersHelper) emitEvent(description string, severity events.EventSeverity, metadata map[string]any, userID core.Uuid) {
	if mch.eventBroadcaster != nil && mch.systemID != nil {
		actedUpon := controllerEventActedUpon(userID, metadata)

		// Own the event by the initiating user when there is one; otherwise fall
		// back to the system. Controller events are frequently emitted from
		// background FSM reconciles (server-restart reconnect, MeshSync mode
		// reconcile) with no user in scope. Owning those by the system keeps them
		// both persisted and retrievable, since GetEvents matches
		// owner IN (userID, systemID) — previously a userless event set no owner
		// and was dropped by the persistence guard, so it was never saved.
		owner := userID
		if owner == uuid.Nil {
			owner = *mch.systemID
		}

		eventBuilder := events.NewEvent().
			FromSystem(*mch.systemID).
			FromOwner(owner).
			WithCategory("connection").
			WithAction("update").
			WithSeverity(severity).
			WithDescription(description).
			WithMetadata(metadata)

		if actedUpon != uuid.Nil {
			eventBuilder = eventBuilder.ActedUpon(actedUpon)
		}

		event := eventBuilder.Build()

		if mch.provider != nil {
			if shouldPersistControllerEvent(owner, actedUpon) {
				if err := mch.provider.PersistSystemEvent(*event); err != nil {
					mch.log.Error(fmt.Errorf("failed to persist event: %w", err))
				}
			} else {
				mch.log.Debug("skipping persistence for controller event without a resource association")
			}
		}
		go mch.eventBroadcaster.Publish(owner, event)
	}
}

func controllerEventActedUpon(userID core.Uuid, metadata map[string]any) core.Uuid {
	if metadata != nil {
		switch connectionID := metadata["connectionID"].(type) {
		case string:
			if parsedID := uuid.FromStringOrNil(connectionID); parsedID != uuid.Nil {
				return parsedID
			}
		case core.Uuid:
			if connectionID != uuid.Nil {
				return connectionID
			}
		case *core.Uuid:
			if connectionID != nil && *connectionID != uuid.Nil {
				return *connectionID
			}
		}
	}

	if userID != uuid.Nil {
		return userID
	}

	return uuid.Nil
}

// shouldPersistControllerEvent reports whether a controller event is worth
// persisting: it must have an owner (the initiating user, or the system for
// background reconciles) and be tied to a specific resource (the connection).
// This keeps truly context-free events out of the DB while still persisting the
// background/system-emitted controller events that the FSM raises.
func shouldPersistControllerEvent(owner, actedUpon core.Uuid) bool {
	return owner != uuid.Nil && actedUpon != uuid.Nil
}

// Common helper for both error and warning events with error information
func (mch *MesheryControllersHelper) emitEventWithError(description string, severity events.EventSeverity, err error, metadata map[string]any, userID core.Uuid) {
	eventMetadata := make(map[string]any)

	if metadata != nil {
		maps.Copy(eventMetadata, metadata)
	}

	if err != nil {
		eventMetadata["error"] = err.Error()
	}

	mch.emitEvent(description, severity, eventMetadata, userID)
}

// Helper method to emit error events
func (mch *MesheryControllersHelper) emitErrorEvent(description string, err error, metadata map[string]any, userID core.Uuid) {
	mch.emitEventWithError(description, events.Error, err, metadata, userID)
}

// Helper method to emit warning events
func (mch *MesheryControllersHelper) emitWarningEvent(description string, err error, metadata map[string]any, userID core.Uuid) {
	mch.emitEventWithError(description, events.Warning, err, metadata, userID)
}
