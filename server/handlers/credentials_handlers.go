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

	// The canonical query param is camelCase `credentialId` (schemas
	// deleteUserCredential, and the identifier-naming contract). `credential_id`
	// is the legacy spelling this handler used to read exclusively, kept as a
	// fallback so older clients keep working.
	raw := q.Get("credentialId")
	if raw == "" {
		raw = q.Get("credential_id")
	}

	// FromStringOrNil turns a missing or malformed id into the zero UUID, which
	// would reach the provider as a real delete for the nil id and report
	// success having removed nothing. Reject it here instead.
	credentialID, err := uuid.FromString(raw)
	if err != nil || credentialID.IsNil() {
		// Built once and reused: constructing the error separately for the log
		// and the response duplicates the message, which is how the two drift
		// apart when only one is later edited.
		invalidID := ErrInvalidRequestObject("credentialId must be a valid UUID")
		h.log.Error(invalidID)
		writeMeshkitError(w, invalidID, http.StatusBadRequest)
		return
	}

	_, err = provider.DeleteUserCredential(req, credentialID)
	if err != nil {
		deleteErr := ErrDeleteUserCredential(err)
		h.log.Error(deleteErr)
		writeMeshkitError(w, deleteErr, http.StatusInternalServerError)
		return
	}

	h.log.Info("Credential deleted.")
	w.WriteHeader(http.StatusOK)
}
