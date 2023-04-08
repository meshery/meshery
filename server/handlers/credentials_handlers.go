package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
)

func (h *Handler) CreateUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(fmt.Errorf("error reading request body: %v", err))
		http.Error(w, "unable to read result data", http.StatusInternalServerError)
		return
	}

	userUUID := uuid.FromStringOrNil(user.ID)
	credential := models.Credential{
		UserID: &userUUID,
		Secret: map[string]interface{}{},
	}

	err = json.Unmarshal(bd, &credential)
	if err != nil {
		h.log.Error(fmt.Errorf("error unmarshal request body: %v", err))
		http.Error(w, "unable to parse credential data", http.StatusInternalServerError)
		return
	}

	result := provider.GetGenericPersister().Table("credentials").Create(&credential) // CredentialsDAO.SaveUserCredentials(&credential)
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error saving user credentials: %v", result.Error))
		http.Error(w, "unable to save user credentials", http.StatusInternalServerError)
		return
	}

	h.log.Info("credential saved successfully")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) ReadUserCredentials(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
}
