//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"github.com/layer5io/meshery/models"
	"github.com/vmihailenco/taskq/v3"
)

// Handler type is the bucket for configs and http handlers
type Handler struct {
	config *models.HandlerConfig
	task   *taskq.Task
}

// NewHandlerInstance returns a Handler instance
func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
) models.HandlerInterface {
	h := &Handler{
		config: handlerConfig,
	}

	h.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}
