package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/server/models"
)

func (h *Handler) CreateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) ReadUserCredentials(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}
