// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/meshery/meshery/server/machines"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshery/server/models/connections"
	gopolicies "github.com/meshery/meshery/server/policies"
	"github.com/meshery/meshkit/broker"
	"github.com/meshery/meshkit/database"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/meshmodel/core/policies"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/meshery/meshkit/utils/events"
	"github.com/meshery/schemas/models/core"
	"github.com/spf13/viper"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config          *models.HandlerConfig
	MeshsyncChannel chan struct{}
	log             logger.Handler
	// to be removed
	brokerConn         broker.Handler
	K8sCompRegHelper   *models.ComponentsRegistrationHelper
	MesheryCtrlsHelper *models.MesheryControllersHelper
	Provider           string // When set, all endpoints consider tokens / identities / capabilities valid from the single, designated provider.
	SystemID           *core.Uuid
	dbHandler          *database.Handler
	registryManager    *meshmodel.RegistryManager
	EventsBuffer       *events.EventStreamer
	// rego holds the OPA/Rego relationship-policy engine, published once by
	// the boot seeding goroutine (server/cmd/main.go) after policies load.
	// Publication must go through SetRego: the goroutine's publish and the
	// request goroutines' reads are unordered, so a plain field write races
	// every h.Rego dereference, and an evaluation would run against a
	// zero-value Rego before seeding finished. Nil until published -
	// regoEngine reports not-ready instead of evaluating a zero-value
	// engine.
	rego                                    atomic.Pointer[policies.Rego]
	GoEngine                                *gopolicies.GoEngine
	ConnectionToStateMachineInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
	MeshsyncDefaultDeploymentMode           connections.MeshsyncDeploymentMode
	// notFoundProbeThrottle gates background Discovery probes of NOTFOUND
	// Kubernetes connections; see notFoundProbeThrottle in middlewares.go.
	notFoundProbeThrottle *notFoundProbeThrottle
	evalTracker           *evaluationTracker
	adapterOperations     *adapterOperationOwners
	// policyEvalTimeoutWarnOnce limits the invalid POLICY_EVAL_TIMEOUT
	// warning to one per Handler; see policyEvalTimeout.
	policyEvalTimeoutWarnOnce sync.Once
}

// NewHandlerInstance returns a Handler instance
func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
	meshSyncCh chan struct{},
	logger logger.Handler,
	brokerConn broker.Handler,
	compRegHelper *models.ComponentsRegistrationHelper,
	mctrlHelper *models.MesheryControllersHelper,
	dbHandler *database.Handler,
	eb *events.EventStreamer,
	regManager *meshmodel.RegistryManager,
	provider string,
	connToInstanceTracker *machines.ConnectionToStateMachineInstanceTracker,
	meshsyncDefaultDeploymentMode connections.MeshsyncDeploymentMode,
) models.HandlerInterface {

	h := &Handler{
		config:                                  handlerConfig,
		MeshsyncChannel:                         meshSyncCh,
		log:                                     logger,
		brokerConn:                              brokerConn,
		K8sCompRegHelper:                        compRegHelper,
		MesheryCtrlsHelper:                      mctrlHelper,
		dbHandler:                               dbHandler,
		EventsBuffer:                            eb,
		registryManager:                         regManager,
		Provider:                                provider,
		GoEngine:                                gopolicies.NewGoEngine(logger),
		adapterOperations:                       newAdapterOperationOwners(),
		SystemID:                                viper.Get("INSTANCE_ID").(*core.Uuid),
		ConnectionToStateMachineInstanceTracker: connToInstanceTracker,
		MeshsyncDefaultDeploymentMode:           meshsyncDefaultDeploymentMode,
		notFoundProbeThrottle:                   newNotFoundProbeThrottle(time.Now, notFoundProbeInterval),
		evalTracker:                             newEvaluationTracker(),
	}

	return h
}

// SetRego publishes the OPA/Rego relationship-policy engine to every future
// evaluation. Called whenever relationships are re-registered at runtime
// (component_generator_helper.go), which swaps the engine atomically so an
// in-flight evaluation keeps the instance it loaded. Boot-time publication
// must use PublishBootRego so it cannot overwrite a runtime reload.
func (h *Handler) SetRego(r *policies.Rego) {
	h.rego.Store(r)
}

// PublishBootRego publishes the boot seeding engine only while no engine is
// published. A runtime reload that already published keeps its engine.
func (h *Handler) PublishBootRego(r *policies.Rego) {
	h.rego.CompareAndSwap(nil, r)
}

// regoEngine returns the published OPA/Rego instance, or ErrPolicyEngineNotReady
// while boot seeding has not published one yet. Evaluated designs must never
// run against a zero-value Rego.
func (h *Handler) regoEngine() (*policies.Rego, error) {
	if r := h.rego.Load(); r != nil {
		return r, nil
	}
	return nil, ErrPolicyEngineNotReady()
}
