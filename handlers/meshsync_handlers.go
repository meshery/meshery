package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/models"
)

// MeshSyncHandler - handles that parses meshsync response
func (h *Handler) MeshSyncHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	data := []byte("mock data")
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
