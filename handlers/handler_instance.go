package handlers

import "github.com/layer5io/meshery/models"

type Handler struct {
	config *models.HandlerConfig
}

func NewHandlerInstance(
	handlerConfig *models.HandlerConfig,
) models.HandlerInterface {
	h := &Handler{
		config: handlerConfig,
	}
	return h
}
