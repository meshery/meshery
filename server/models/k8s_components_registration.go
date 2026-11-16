package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

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
	RegistrationPartial
)

// ErrPartialRegistration indicates that component registration completed partially with warnings.
type ErrPartialRegistration struct {
	Err error
}

func (e *ErrPartialRegistration) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return "partial registration completed with warnings"
}

func (e *ErrPartialRegistration) Unwrap() error {
	return e.Err
}

func NewErrPartialRegistration(err error) error {
	return &ErrPartialRegistration{Err: err}
}

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
	case RegistrationPartial:
		return "partial"
	default:
		return ""
	}
}

type ComponentsRegistrationHelper struct {
	// map that holds the registration status of each of the contexts in this runtime of the server
	// it should be private
	ctxRegStatusMap map[string]RegistrationStatus
	ctxRegTokenMap  map[string]uint64
	regTokenSeq     uint64
	log             logger.Handler
	mx              sync.RWMutex
}

// NewComponentsRegistrationHelper initializes and returns a new ComponentsRegistrationHelper.
func NewComponentsRegistrationHelper(logger logger.Handler) *ComponentsRegistrationHelper {
	return &ComponentsRegistrationHelper{
		ctxRegStatusMap: make(map[string]RegistrationStatus),
		ctxRegTokenMap:  make(map[string]uint64),
		log:             logger,
		mx:              sync.RWMutex{},
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

// GetContextRegistrationStatus returns the current registration status of the given context ID.
func (cg *ComponentsRegistrationHelper) GetContextRegistrationStatus(ctxID string) (RegistrationStatus, bool) {
	cg.mx.RLock()
	defer cg.mx.RUnlock()
	status, ok := cg.ctxRegStatusMap[ctxID]
	return status, ok
}

// RemoveContext removes a context from the registration map, allowing future registrations.
func (cg *ComponentsRegistrationHelper) RemoveContext(ctxID string) *ComponentsRegistrationHelper {
	cg.mx.Lock()
	defer cg.mx.Unlock()
	delete(cg.ctxRegStatusMap, ctxID)
	delete(cg.ctxRegTokenMap, ctxID)
	return cg
}

type K8sRegistrationFunction func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID core.Uuid, reg *meshmodel.RegistryManager, eb *Broadcast, log logger.Handler, ctxName string) error

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
		if ctx == nil {
			continue
		}
		ctxID := ctx.ID
		connectionID := uuid.FromStringOrNil(ctx.ConnectionID)
		ctxName := ctx.Name
		var mesheryInstanceID core.Uuid
		if ctx.MesheryInstanceID != nil {
			mesheryInstanceID = *ctx.MesheryInstanceID
		}

		cg.mx.Lock()
		// do not do anything about the contexts that are not present in the ctxRegStatusMap
		// only start registering components for contexts whose status is NotRegistered or RegistrationPartial
		status, ok := cg.ctxRegStatusMap[ctxID]
		if !ok || (status != NotRegistered && status != RegistrationPartial) {
			cg.mx.Unlock()
			continue
		}

		// update the status
		cg.ctxRegStatusMap[ctxID] = Registering
		cg.regTokenSeq++
		regToken := cg.regTokenSeq
		cg.ctxRegTokenMap[ctxID] = regToken
		cg.mx.Unlock()
		cg.log.Info("Registration of ", ctxName, " components started for contextID: ", ctxID)

		event := events.NewEvent().ActedUpon(connectionID).FromSystem(mesheryInstanceID).WithSeverity(events.Informational).WithCategory("connection").WithAction(Registering.String()).FromOwner(userUUID).WithDescription(fmt.Sprintf("Registration for Kubernetes context %s started", ctxName)).Build()
		// Even if the event was not persisted, the event is published to the user
		// and the operation continues.
		if err := PersistAndPublishSystemEvent(provider, eventsBrodcaster, userUUID, event); err != nil {
			cg.log.Warn(err)
		}

		go func(ctx *K8sContext, ctxID string, connectionID core.Uuid, ctxName string, mesheryInstanceID core.Uuid, regToken uint64) {
			var regErr error
			var isPartial bool
			var partialErr error
			defer func() {
				cg.mx.Lock()
				// Only update status and emit events if context was not removed
				// and this registration generation is still current.
				currentToken, exists := cg.ctxRegTokenMap[ctxID]
				if !exists || currentToken != regToken {
					cg.mx.Unlock()
					return
				}
				if regErr != nil {
					cg.ctxRegStatusMap[ctxID] = NotRegistered
				} else if isPartial {
					cg.ctxRegStatusMap[ctxID] = RegistrationPartial
				} else {
					cg.ctxRegStatusMap[ctxID] = RegistrationComplete
				}
				cg.mx.Unlock()

				if regErr != nil {
					cg.log.Errorf("failed to register components for context %s (ID: %s): %v", ctxName, ctxID, regErr)
					failEvent := events.NewEvent().
						ActedUpon(connectionID).
						FromSystem(mesheryInstanceID).
						WithSeverity(events.Error).
						WithCategory("connection").
						WithAction(Registering.String()).
						FromOwner(userUUID).
						WithDescription(fmt.Sprintf("Failed to register components for Kubernetes context %s: %s", ctxName, regErr.Error())).
						WithMetadata(map[string]interface{}{"error": regErr.Error()}).
						Build()
					if perr := PersistAndPublishSystemEvent(provider, eventsBrodcaster, userUUID, failEvent); perr != nil {
						cg.log.Warn(perr)
					}
				} else if isPartial {
					cg.log.Warnf("components partially registered with warnings for context %s (ID: %s): %v", ctxName, ctxID, partialErr)
				} else {
					cg.log.Info("components registered for context ", ctxName, " ID:", ctxID)
				}
			}()

			// start registration
			cfg, err := ctx.GenerateKubeConfig()
			if err != nil {
				regErr = err
				cg.log.Error(err)
				return
			}
			for _, f := range regFunc {
				err = f(&provider, context.Background(), cfg, ctxID, ctx.ConnectionID, userID, mesheryInstanceID, reg, eventsBrodcaster, cg.log, ctxName)
				if err != nil {
					var pErr *ErrPartialRegistration
					if errors.As(err, &pErr) {
						isPartial = true
						partialErr = err
						continue
					}
					regErr = ErrUnreachableKubeAPI(err, ctx.Server)
					cg.log.Error(regErr)
					return
				}
			}
		}(ctx, ctxID, connectionID, ctxName, mesheryInstanceID, regToken)
	}
}

// loadK8sMeshModelMetadata reads and unmarshals the Kubernetes mesh model template from the specified path,
// ensuring the file descriptor is closed across all return paths.
func loadK8sMeshModelMetadata(path string) (ModelTemplate, error) {
	byt, err := os.ReadFile(path)
	if err != nil {
		return ModelTemplate{}, err
	}

	componentStyles := ModelTemplate{}
	err = json.Unmarshal(byt, &componentStyles)
	if err != nil {
		return ModelTemplate{}, err
	}

	return componentStyles, nil
}

// Caches k8sMeshModel metadatas in memory to use at the time of dynamic k8s component generation
func init() {
	metadata, err := loadK8sMeshModelMetadata(k8sMeshModelPath)
	if err != nil {
		return
	}
	K8sMeshModelMetadata = metadata
}
