package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
)

func (h *Handler) SaveUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	userUUID := user.ID
	credential := models.Credential{
		Secret: map[string]interface{}{},
	}

	err = json.Unmarshal(bd, &credential)
	if err != nil {
		h.log.Error(ErrDecoding(err, "credential"))
		writeMeshkitError(w, ErrDecoding(err, "credential"), http.StatusBadRequest)
		return
	}

	// Bind credential ownership to the authenticated user AFTER unmarshal so a
	// client-supplied `userId` in the request body cannot redirect a credential
	// onto another user's account.
	credential.UserId = userUUID

	createdCredential, err := provider.SaveUserCredential(token, &credential)
	if err != nil {
		h.log.Error(ErrSaveUserCredential(err))
		writeMeshkitError(w, ErrSaveUserCredential(err), http.StatusInternalServerError)
		return
	}

	h.log.Info("Credential saved.", createdCredential.Name, createdCredential.ID)
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUserCredentialByID(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	credentialID := uuid.FromStringOrNil(mux.Vars(req)["credentialID"])
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	credential, statusCode, err := provider.GetCredentialByID(token, credentialID)
	if err != nil {
		h.log.Error(ErrGetUserCredential(err))
		writeMeshkitError(w, ErrGetUserCredential(err), statusCode)
		return
	}

	if err := json.NewEncoder(w).Encode(credential); err != nil {
		h.log.Error(ErrEncodeUserCredential(err))
		writeMeshkitError(w, ErrEncodeUserCredential(err), http.StatusInternalServerError)
		return
	}
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

	credentialsPage, err := provider.GetUserCredentials(req, user.ID.String(), page, pageSize, search, order)
	if err != nil {
		h.log.Error(ErrGetUserCredential(err))
		writeMeshkitError(w, ErrGetUserCredential(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(credentialsPage); err != nil {
		h.log.Error(ErrEncodeUserCredential(err))
		writeMeshkitError(w, ErrEncodeUserCredential(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		h.log.Error(ErrRequestBody(err))
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	userUUID := user.ID
	credential := &models.Credential{
		Secret: map[string]interface{}{},
	}
	err = json.Unmarshal(bd, credential)
	if err != nil {
		h.log.Error(ErrDecoding(err, "credential"))
		writeMeshkitError(w, ErrDecoding(err, "credential"), http.StatusBadRequest)
		return
	}

	// Bind credential ownership to the authenticated user AFTER unmarshal so a
	// client-supplied `userId` in the request body cannot hijack another user's
	// credential (the provider layer's authorization check should rely on this
	// field to confirm the caller owns the credential being updated).
	credential.UserId = userUUID

	_, err = provider.UpdateUserCredential(req, credential)
	if err != nil {
		h.log.Error(ErrUpdateUserCredential(err))
		writeMeshkitError(w, ErrUpdateUserCredential(err), http.StatusInternalServerError)
		return
	}

	h.log.Info("Credential updated.")
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, provider models.Provider) {
	q := req.URL.Query()

	credentialID := uuid.FromStringOrNil(q.Get("credential_id"))
	_, err := provider.DeleteUserCredential(req, credentialID)
	if err != nil {
		h.log.Error(ErrDeleteUserCredential(err))
		writeMeshkitError(w, ErrDeleteUserCredential(err), http.StatusInternalServerError)
		return
	}

	h.log.Info("Credential deleted.")
	w.WriteHeader(http.StatusOK)
}
