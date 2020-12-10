package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/models"
)

func (h *Handler) SmiTestHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// To be implemented
}
