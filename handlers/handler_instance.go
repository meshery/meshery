package handlers

import (
	"github.com/layer5io/meshery/models"
	"github.com/vmihailenco/taskq"
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

	h.task = handlerConfig.Queue.NewTask(&taskq.TaskOptions{
		Name:    "submitMetrics",
		Handler: h.CollectStaticMetrics,
	})

	return h
}

// GetProviderType - retrieves the provider type
func (h *Handler) GetProviderType() models.ProviderType {
	return h.config.Provider.GetProviderType()
}
