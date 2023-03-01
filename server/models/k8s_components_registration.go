package models

import (
	"context"
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/meshes"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/utils/events"
	"github.com/spf13/viper"
)

type RegistrationStatus int

const (
	RegistrationComplete RegistrationStatus = iota
	NotRegistered
	Registering
)

type ComponentsRegistrationHelper struct {
	// map that holds the registration status of each of the contexts in this runtime of the server
	// it should be private
	ctxRegStatusMap map[string]RegistrationStatus
	log             logger.Handler
}

func NewComponentsRegistrationHelper(logger logger.Handler) *ComponentsRegistrationHelper {
	return &ComponentsRegistrationHelper{
		ctxRegStatusMap: make(map[string]RegistrationStatus),
		log:             logger,
	}
}

// update the map with the given list of contexts
func (cg *ComponentsRegistrationHelper) UpdateContexts(ctxs []*K8sContext) *ComponentsRegistrationHelper {
	for _, ctx := range ctxs {
		ctxID := ctx.ID
		if _, ok := cg.ctxRegStatusMap[ctxID]; !ok {
			cg.ctxRegStatusMap[ctxID] = NotRegistered
		}
	}
	return cg
}

type K8sRegistrationFunction func(ctxt context.Context, config []byte, ctxID string, reg *meshmodel.RegistryManager, es *events.EventStreamer, ctxName string) error

// start registration of components for the contexts
func (cg *ComponentsRegistrationHelper) RegisterComponents(ctxs []*K8sContext, regFunc []K8sRegistrationFunction, eb *events.EventStreamer, reg *meshmodel.RegistryManager) {
	if viper.GetBool("SKIP_COMP_GEN") {
		return
	}

	for _, ctx := range ctxs {
		ctxID := ctx.ID

		// do not do anything about the contexts that are not present in the ctxRegStatusMap
		// only start registering components for contexts whose status is NotRegistered
		status, ok := cg.ctxRegStatusMap[ctxID]
		if !ok || status != NotRegistered {
			continue
		}

		ctxName := ctx.Name
		id, _ := uuid.NewV4()

		// update the status
		cg.ctxRegStatusMap[ctxID] = Registering
		cg.log.Info("Registration of ", ctxName, " components started for contextID: ", ctxID)
		req := meshes.EventsResponse{
			Component:     "core",
			ComponentName: "Kubernetes",
			EventType:     meshes.EventType_INFO,
			Summary:       fmt.Sprintf("Registration for Kubernetes context \"%s\" started", ctxName),
			Details:       fmt.Sprintf("Registration for Kubernetes context \"%s\" started with context ID %s", ctxName, ctxID),
			OperationId:   id.String(),
		}
		eb.Publish(&req)

		var mu sync.Mutex // declare a mutex for synchronizing access to the map

		go func(ctx *K8sContext) {
			// set the status to RegistrationComplete
			defer func() {
				mu.Lock()
				cg.ctxRegStatusMap[ctxID] = RegistrationComplete
				mu.Unlock()

				cg.log.Info(ctxName, " components for contextID:", ctxID, " registered")
			}()

			// start registration
			cfg, err := ctx.GenerateKubeConfig()
			if err != nil {
				cg.log.Error(err)
				return
			}
			for _, f := range regFunc {
				err = f(context.Background(), cfg, ctxID, reg, eb, ctxName)
				if err != nil {
					cg.log.Error(err)
					return
				}
			}
		}(ctx)
	}
}
