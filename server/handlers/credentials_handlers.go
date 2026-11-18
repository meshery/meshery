package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/errors"
)

func (h *Handler) SaveUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}
	token, _ := req.Context().Value(models.TokenCtxKey).(string)

	bd, err := io.ReadAll(req.Body)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	userUUID := user.ID
	credential := models.Credential{
		Secret: map[string]interface{}{},
	}

	err = json.Unmarshal(bd, &credential)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrDecoding(err, "credential"))
		}
		writeMeshkitError(w, ErrDecoding(err, "credential"), http.StatusBadRequest)
		return
	}

	// Bind credential ownership to the authenticated user AFTER unmarshal so a
	// client-supplied `userId` in the request body cannot redirect a credential
	// onto another user's account.
	credential.UserId = userUUID

	createdCredential, err := provider.SaveUserCredential(token, &credential)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrSaveUserCredential(err))
		}
		writeMeshkitError(w, ErrSaveUserCredential(err), http.StatusInternalServerError)
		return
	}

	if h.log != nil {
		h.log.Info("Credential saved.", createdCredential.Name, createdCredential.ID)
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetUserCredentialByID(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}
	raw := mux.Vars(req)["credentialID"]
	credentialID, err := uuid.FromString(raw)
	if err != nil || credentialID.IsNil() {
		invalidID := models.ErrInvalidUUID(fmt.Errorf("credentialID must be a valid UUID"))
		if h.log != nil {
			h.log.Error(invalidID)
		}
		writeMeshkitError(w, invalidID, http.StatusBadRequest)
		return
	}
	token, _ := req.Context().Value(models.TokenCtxKey).(string)
	credential, statusCode, err := provider.GetCredentialByID(token, credentialID)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrGetUserCredential(err))
		}
		writeMeshkitError(w, ErrGetUserCredential(err), statusCode)
		return
	}

	// Enforce user ownership: fail closed. Only the explicit owner can read a credential;
	// unowned credentials (credential.UserId == uuid.Nil) or requests without a resolved caller are rejected with 404.
	if user == nil || user.ID == uuid.Nil || credential.UserId == uuid.Nil || credential.UserId != user.ID {
		notOwned := models.ErrResultNotFound(fmt.Errorf("credential with id %s not found", credentialID))
		writeMeshkitError(w, notOwned, http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(credential); err != nil {
		if h.log != nil {
			h.log.Error(ErrEncodeUserCredential(err))
		}
		writeMeshkitError(w, ErrEncodeUserCredential(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetUserCredentials(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}
	q := req.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	order := q.Get("order")
	search := q.Get("search")
	pageSize, _ := strconv.Atoi(getPageSizeParam(q))

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

	if h.log != nil {
		h.log.Debug(fmt.Sprintf("page: %d, page size: %d, search: %s, order: %s", page+1, pageSize, search, order))
	}

	credentialsPage, err := provider.GetUserCredentials(req, user.ID.String(), page, pageSize, search, order)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrGetUserCredential(err))
		}
		writeMeshkitError(w, ErrGetUserCredential(err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(credentialsPage); err != nil {
		if h.log != nil {
			h.log.Error(ErrEncodeUserCredential(err))
		}
		writeMeshkitError(w, ErrEncodeUserCredential(err), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if user == nil {
		writeJSONError(w, "user unauthorized", http.StatusUnauthorized)
		return
	}
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}
	bd, err := io.ReadAll(req.Body)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrRequestBody(err))
		}
		writeMeshkitError(w, ErrRequestBody(err), http.StatusInternalServerError)
		return
	}

	userUUID := user.ID
	credential := &models.Credential{
		Secret: map[string]interface{}{},
	}
	err = json.Unmarshal(bd, credential)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrDecoding(err, "credential"))
		}
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
		// The local provider reports a missing or deleted credential with
		// ErrResultNotFound, the remote provider with its own 404.
		if errors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			writeMeshkitError(w, err, http.StatusNotFound)
			return
		}
		if h.log != nil {
			h.log.Error(ErrUpdateUserCredential(err))
		}
		writeMeshkitError(w, ErrUpdateUserCredential(err), http.StatusInternalServerError)
		return
	}

	if h.log != nil {
		h.log.Info("Credential updated.")
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, _ *models.Preference, user *models.User, provider models.Provider) {
	if provider == nil {
		if h.log != nil {
			h.log.Error(ErrNilProvider())
		}
		writeMeshkitError(w, ErrNilProvider(), http.StatusInternalServerError)
		return
	}
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
		invalidID := models.ErrInvalidUUID(fmt.Errorf("credentialId must be a valid UUID"))
		if h.log != nil {
			h.log.Error(invalidID)
		}
		writeMeshkitError(w, invalidID, http.StatusBadRequest)
		return
	}

	if user != nil {
		req = req.WithContext(context.WithValue(req.Context(), models.UserCtxKey, user))
	}

	_, err = provider.DeleteUserCredential(req, credentialID)
	if err != nil {
		if errors.GetCode(err) == models.ErrResultNotFoundCode || providerStatus(err) == http.StatusNotFound {
			writeMeshkitError(w, err, http.StatusNotFound)
			return
		}
		deleteErr := ErrDeleteUserCredential(err)
		if h.log != nil {
			h.log.Error(deleteErr)
		}
		writeMeshkitError(w, deleteErr, http.StatusInternalServerError)
		return
	}

	if h.log != nil {
		h.log.Info("Credential deleted.")
	}
	w.WriteHeader(http.StatusOK)
}
