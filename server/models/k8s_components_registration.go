package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofrs/uuid"

	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/events"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/spf13/viper"
)

const k8sMeshModelPath = "../meshmodel/kubernetes/model_template.json"

type RegistrationStatus int

const (
	RegistrationComplete RegistrationStatus = iota
	NotRegistered
	Registering
)

var K8sMeshModelMetadata = make(map[string]interface{})

// INstead define a set of actions
func (rs RegistrationStatus) String() string {
	switch rs {
	case RegistrationComplete:
		return "register"
	case NotRegistered:
		return "not_registered"
	case Registering:
		return "registering"
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
}

func NewComponentsRegistrationHelper(logger logger.Handler) *ComponentsRegistrationHelper {
	return &ComponentsRegistrationHelper{
		ctxRegStatusMap: make(map[string]RegistrationStatus),
		log:             logger,
		mx:              sync.RWMutex{},
	}
}

// update the map with the given list of contexts
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

type K8sRegistrationFunction func(provider *Provider, ctxt context.Context, config []byte, ctxID string, connectionID string, userID string, MesheryInstanceID uuid.UUID, reg *meshmodel.RegistryManager, eb *Broadcast, ctxName string) error

type InitFunc func(ctx context.Context, machineCtx interface{}, log logger.Handler) (interface{}, *events.Event, error)

// start registration of components for the contexts
func (cg *ComponentsRegistrationHelper) RegisterComponents(ctxs []*K8sContext, regFunc []K8sRegistrationFunction, reg *meshmodel.RegistryManager, eventsBrodcaster *Broadcast, provider Provider, userID string, skip bool) {
	/* If flag "SKIP_COMP_GEN" is set but the registration is invoked in form of API request explicitly,
	then flag should not be respected and to control this behaviour skip is introduced.
	In case of API requests "skip" is set to false, otherise true and behaviour is controlled by "SKIP_COMP_GEN".
	*/
	if viper.GetBool("SKIP_COMP_GEN") && skip {
		return
	}

	userUUID, _ := uuid.FromString(userID)

	for _, ctx := range ctxs {
		ctxID := ctx.ID
		connectionID, _ := uuid.FromString(ctx.ConnectionID)
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

		event := events.NewEvent().ActedUpon(connectionID).FromSystem(*ctx.MesheryInstanceID).WithSeverity(events.Informational).WithCategory("connection").WithAction(Registering.String()).FromUser(userUUID).WithDescription(fmt.Sprintf("Registration for Kubernetes context %s started", ctxName)).Build()
		err := provider.PersistEvent(event)
		if err != nil {
			// Even if event was not persisted continue with the operation and publish the event to user.
			cg.log.Warn(err)
		}
		eventsBrodcaster.Publish(userUUID, event)

		go func(ctx *K8sContext) {
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
			for _, f := range regFunc {
				err = f(&provider, context.Background(), cfg, ctxID, ctx.ConnectionID, userID, *ctx.MesheryInstanceID, reg, eventsBrodcaster, ctxName)
				if err != nil {
					cg.log.Error(err)
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
	m := make(map[string]interface{})
	err = json.Unmarshal(byt, &m)
	if err != nil {
		return
	}
	K8sMeshModelMetadata = m
}
