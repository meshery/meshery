package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/meshery/schemas/models/core"

	"github.com/gofrs/uuid"
	"github.com/meshery/schemas/models/v1beta1/capability"
	"github.com/meshery/schemas/models/v1beta3/component"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/models/events"
	meshmodel "github.com/meshery/meshkit/models/meshmodel/registry"
	"github.com/spf13/viper"
)

const k8sMeshModelPath = "../../models/kubernetes/model_template.json"

type RegistrationStatus int

const (
	RegistrationComplete RegistrationStatus = iota
	NotRegistered
	Registering
)

type ModelTemplate struct {
	component.Styles
	ModelName    string                   `json:"modelName"`
	Capabilities *[]capability.Capability `json:"capabilities"`
}

var K8sMeshModelMetadata = ModelTemplate{}

// String converts a RegistrationStatus into its string representation
func (rs RegistrationStatus) String() string {
	switch rs {
	case RegistrationComplete:
		return "register"
	case NotRegistered:
		return "not_registered"
	case Registering:
		return "register"
	default:
		return ""
	}
}

type ComponentsRegistrationHelper struct {
	// map that holds the registration status of each of the contexts in this runtime of the server
	// it should be private
	ctxRegStatusMap map[string]RegistrationStatus
	log             logger.Handler
	mx              sync.RWMutex
	// regSem bounds how many contexts are registered concurrently across all
	// RegisterComponents calls on this helper. It is shared (not per-call) because
	// the helper itself is a single long-lived instance reused across concurrent
	// requests; a per-call semaphore would let N concurrent requests each get their
	// own budget, multiplying the effective concurrency by N.
	regSem chan struct{}
}

func NewComponentsRegistrationHelper(logger logger.Handler) *ComponentsRegistrationHelper {
	return &ComponentsRegistrationHelper{
		ctxRegStatusMap: make(map[string]RegistrationStatus),
		log:             logger,
		mx:              sync.RWMutex{},
		regSem:          make(chan struct{}, maxConcurrentK8sRegistrations),
	}
}

// UpdateContexts updates the map with the given list of contexts
func (cg *ComponentsRegistrationHelper) UpdateContexts(ctxs []*K8sContext) *ComponentsRegistrationHelper {
	for _, ctx := range ctxs {
		ctxID := ctx.ID
		cg.mx.Lock()
		if _, ok := cg.ctxRegStatusMap[ctxID]; !ok {
			cg.ctxRegStatusMap[ctxID] = NotRegistered
		}
		cg.mx.Unlock()
	}
	return cg
}

type K8sRegistrationFunction func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID core.Uuid, reg *meshmodel.RegistryManager, eb *Broadcast, log logger.Handler, ctxName string) error

// maxConcurrentK8sRegistrations bounds how many contexts RegisterComponents will
// register concurrently, independent of how many contexts a single request
// supplies. Without this cap, a kubeconfig with a large number of contexts (which
// need not point at real, reachable clusters) drives a goroutine, file descriptor,
// and outbound-dial spike proportional to attacker-controlled input.
const maxConcurrentK8sRegistrations = 10

// k8sRegistrationTimeout bounds how long a single context's registration may run.
// Without a deadline, a slow or unreachable API server can hold a semaphore slot
// (see regSem) indefinitely, starving every other context waiting to register.
const k8sRegistrationTimeout = 30 * time.Second

// RegisterComponents starts registration of components for the contexts
func (cg *ComponentsRegistrationHelper) RegisterComponents(ctxs []*K8sContext, regFunc []K8sRegistrationFunction, reg *meshmodel.RegistryManager, eventsBrodcaster *Broadcast, provider Provider, userID string, skip bool) {
	/* If flag "SKIP_COMP_GEN" is set but the registration is invoked in form of API request explicitly,
	then flag should not be respected and to control this behaviour skip is introduced.
	In case of API requests "skip" is set to false, otherise true and behaviour is controlled by "SKIP_COMP_GEN".
	*/
	if viper.GetBool("SKIP_COMP_GEN") && skip {
		return
	}

	userUUID, err := uuid.FromString(userID)
	if err != nil {
		// Every event below is attributed FromOwner(userUUID) and broadcast to the
		// user key; a nil owner would mis-attribute the whole registration. Bail.
		cg.log.Error(ErrInvalidUUID(fmt.Errorf("invalid user id %q: %w", userID, err)))
		return
	}

	for _, ctx := range ctxs {
		ctxID := ctx.ID
		connectionID := uuid.FromStringOrNil(ctx.ConnectionID)
		ctxName := ctx.Name

		cg.mx.Lock()
		// do not do anything about the contexts that are not present in the ctxRegStatusMap
		// only start registering components for contexts whose status is NotRegistered
		status, ok := cg.ctxRegStatusMap[ctxID]
		if !ok || status != NotRegistered {
			cg.mx.Unlock()
			continue
		}

		// update the status
		cg.ctxRegStatusMap[ctxID] = Registering
		cg.mx.Unlock()
		cg.log.Info("Registration of ", ctxName, " components started for contextID: ", ctxID)

		event := events.NewEvent().ActedUpon(connectionID).FromSystem(*ctx.MesheryInstanceID).WithSeverity(events.Informational).WithCategory("connection").WithAction(Registering.String()).FromOwner(userUUID).WithDescription(fmt.Sprintf("Registration for Kubernetes context %s started", ctxName)).Build()
		err := provider.PersistSystemEvent(*event)
		if err != nil {
			// Even if event was not persisted continue with the operation and publish the event to user.
			cg.log.Warn(err)
		}
		eventsBrodcaster.Publish(userUUID, event)

		// Acquire the shared semaphore inside the goroutine, not on this loop, so a
		// full semaphore blocks only queued registrations and never the caller (the
		// HTTP handler still returns its 202 Accepted immediately).
		go func(ctx *K8sContext) {
			cg.regSem <- struct{}{}
			defer func() { <-cg.regSem }()

			// set the status to RegistrationComplete
			defer func() {
				cg.mx.Lock()
				cg.ctxRegStatusMap[ctxID] = RegistrationComplete
				cg.mx.Unlock()

				cg.log.Info("components registered for context ", ctxName, " ID:", ctxID)
			}()

			// start registration
			cfg, err := ctx.GenerateKubeConfig()
			if err != nil {
				cg.log.Error(err)
				return
			}

			// Bound how long this context can hold its semaphore slot; see
			// k8sRegistrationTimeout.
			regCtx, cancel := context.WithTimeout(context.Background(), k8sRegistrationTimeout)
			defer cancel()

			for _, f := range regFunc {
				err = f(&provider, regCtx, cfg, ctxID, ctx.ConnectionID, userID, *ctx.MesheryInstanceID, reg, eventsBrodcaster, cg.log, ctxName)
				if err != nil {
					cg.log.Error(ErrUnreachableKubeAPI(err, ctx.Server))
					return
				}
			}
		}(ctx)
	}
}

// Caches k8sMeshModel metadatas in memory to use at the time of dynamic k8s component generation
func init() {
	f, err := os.Open(filepath.Join(k8sMeshModelPath))
	if err != nil {
		return
	}
	byt, err := io.ReadAll(f)
	if err != nil {
		return
	}
	componentStyles := ModelTemplate{}
	err = json.Unmarshal(byt, &componentStyles)
	if err != nil {
		return
	}
	K8sMeshModelMetadata = componentStyles
}
