// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/events"
	"github.com/vmihailenco/taskq/v3"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config *models.HandlerConfig
	task   *taskq.Task
	// to be removed
	meshsyncChannel chan struct{}
	log             logger.Handler
	// to be removed
	brokerConn         broker.Handler
	K8sCompRegHelper   *models.ComponentsRegistrationHelper
	MesheryCtrlsHelper *models.MesheryControllersHelper
	dbHandler          *database.Handler
	EventsBuffer       *events.EventStreamer
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
) models.HandlerInterface {
	h := &Handler{
		config:             handlerConfig,
		meshsyncChannel:    meshSyncCh,
		log:                logger,
		brokerConn:         brokerConn,
		K8sCompRegHelper:   compRegHelper,
		MesheryCtrlsHelper: mctrlHelper,
		dbHandler:          dbHandler,
		EventsBuffer:       eb,
	}

	h.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}
