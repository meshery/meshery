//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/broker"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils/broadcast"
	"github.com/vmihailenco/taskq/v3"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config          *models.HandlerConfig
	task            *taskq.Task
	meshsyncChannel chan struct{}
	log             logger.Handler
	brokerConn      broker.Handler
	broadcaster     broadcast.Broadcaster
}

// NewHandlerInstance returns a Handler instance
func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
	meshSyncCh chan struct{},
	logger logger.Handler,
	brokerConn broker.Handler,
	broadcaster broadcast.Broadcaster,
) models.HandlerInterface {
	h := &Handler{
		config:          handlerConfig,
		meshsyncChannel: meshSyncCh,
		log:             logger,
		brokerConn:      brokerConn,
		broadcaster:     broadcaster,
	}

	h.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}
