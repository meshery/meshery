//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
	"github.com/vmihailenco/taskq/v3"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config          *models.HandlerConfig
	task            *taskq.Task
	meshsyncChannel chan struct{}
	log             logger.Handler
	brokerConn      broker.Handler
}

// NewHandlerInstance returns a Handler instance
func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
	meshSyncCh chan struct{},
	logger logger.Handler,
	brokerConn broker.Handler,
) models.HandlerInterface {
	h := &Handler{
		config:          handlerConfig,
		meshsyncChannel: meshSyncCh,
		log:             logger,
		brokerConn:      brokerConn,
	}

	h.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}
