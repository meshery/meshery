// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/machines"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/meshmodel/core/policies"
	meshmodel "github.com/layer5io/meshkit/models/meshmodel/registry"
	"github.com/layer5io/meshkit/utils/events"
	"github.com/spf13/viper"
	"github.com/vmihailenco/taskq/v3"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config *models.HandlerConfig
	task   *taskq.Task
	MeshsyncChannel chan struct{}
	log             logger.Handler
	// to be removed
	brokerConn                              broker.Handler
	K8sCompRegHelper                        *models.ComponentsRegistrationHelper
	MesheryCtrlsHelper                      *models.MesheryControllersHelper
	Provider                                string // When set, all endpoints consider tokens / identities / capabilities valid from the single, designated provider.
	SystemID                                *uuid.UUID
	dbHandler                               *database.Handler
	registryManager                         *meshmodel.RegistryManager
	EventsBuffer                            *events.EventStreamer
	Rego                                    *policies.Rego
	ConnectionToStateMachineInstanceTracker *machines.ConnectionToStateMachineInstanceTracker
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
	rego *policies.Rego,
	connToInstanceTracker *machines.ConnectionToStateMachineInstanceTracker,
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
		Rego:                                    rego,
		SystemID:                                viper.Get("INSTANCE_ID").(*uuid.UUID),
		ConnectionToStateMachineInstanceTracker: connToInstanceTracker,
	}

	h.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}
