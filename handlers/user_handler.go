package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// UserHandler returns info about the logged in user
func (h *Handler) UserHandler(w http.ResponseWriter, req *http.Request, _ *sessions.Session, _ *models.Session, user *models.User) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		logrus.Errorf("error getting user data: %v", err)
		http.Error(w, "unable to get session", http.StatusInternalServerError)
		return
	}
}
