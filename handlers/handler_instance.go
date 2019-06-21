package handlers

import (
	"github.com/layer5io/meshery/models"
	"github.com/vmihailenco/taskq"
)
type Handler struct {
	config *models.HandlerConfig
	task	*taskq.Task
}

func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
) models.HandlerInterface {
	h := &Handler{
		config: handlerConfig,
	}

	h.task = handlerConfig.Queue.NewTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	
	return h
}
