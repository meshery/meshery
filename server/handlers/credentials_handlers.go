package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/models"
)

func (h *Handler) SaveUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
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

	err = provider.SaveCredential(&credential)
	if err != nil {
		h.log.Error(fmt.Errorf("error saving user credentials: %v", err))
		http.Error(w, "unable to save user credentials", http.StatusInternalServerError)
		return
	}

	h.log.Info("credential saved successfully")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUserCredentials(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	if pageSize > 25 {
		pageSize = 25
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	if page < 0 {
		page = 0
	}
	if order == "" {
		order = "created_at desc"
	}

	h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s", page+1, pageSize, search, order))

	credentialsPage, err := provider.GetCredentials(user.ID, page, pageSize, order, search)
	if err != nil {
		h.log.Error(fmt.Errorf("error getting user credentials: %v", err))
		http.Error(w, "unable to get user credentials", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(credentialsPage); err != nil {
		h.log.Error(fmt.Errorf("error encoding user credentials: %v", err))
		http.Error(w, "unable to encode user credentials", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(fmt.Errorf("error reading request body: %v", err))
		http.Error(w, "unable to read credential data", http.StatusInternalServerError)
		return
	}

	credential := &models.Credential{}
	err = json.Unmarshal(bd, credential)
	if err != nil {
		h.log.Error(fmt.Errorf("error unmarshal request body: %v", err))
		http.Error(w, "unable to parse credential data", http.StatusInternalServerError)
		return
	}

	result := provider.GetGenericPersister().Where("user_id = ? AND id = ? AND deleted_at is NULL", user.UserID, credential.ID).First(&models.Credential{})
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error getting user credential: %v", result.Error))
		http.Error(w, "unable to get user credential", http.StatusInternalServerError)
		return
	}

	result = provider.GetGenericPersister().Save(credential)
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error updating user credential: %v", result.Error))
		http.Error(w, "unable to update user credential", http.StatusInternalServerError)
		return
	}

	h.log.Info("credential updated successfully")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	q := req.URL.Query()

	credentialID := q.Get("credential_id")

	result := provider.GetGenericPersister().Model(&models.Credential{}).Where("user_id = ? AND id = ? AND deleted_at is NULL", user.UserID, credentialID).Update("deleted_at", time.Now())
	if result.Error != nil {
		h.log.Error(fmt.Errorf("error deleting user credential: %v", result.Error))
		http.Error(w, "unable to delete user credential", http.StatusInternalServerError)
		return
	}

	h.log.Info("credential deleted successfully")
	w.WriteHeader(http.StatusOK)
}
