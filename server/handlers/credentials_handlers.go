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

func providerUserID(user *models.User) string {
	if user == nil || user.ID == uuid.Nil {
		return ""
	}
	return user.ID.String()
}

func sessionRequest(req *http.Request, prefObj *models.Preference, user *models.User) *http.Request {
	if req == nil {
		return req
	}
	if prefObj != nil {
		req = req.WithContext(context.WithValue(req.Context(), models.PerfObjCtxKey, prefObj))
	}
	if user != nil {
		req = req.WithContext(context.WithValue(req.Context(), models.UserCtxKey, user))
	}
	return req
}

// writeCredentialProviderError writes a provider failure from a credential
// endpoint. 409 is the organization_selection_required contract body, not
// MeshKit. 403 and 404 stay distinct MeshKit errors. Missing local
// credentials (ErrResultNotFound) stay 404. Other failures are wrapped and
// written with the provider status, falling back to fallback (or 500).
func (h *Handler) writeCredentialProviderError(w http.ResponseWriter, err error, wrap func(error) error, fallback int) {
	if err == nil {
		return
	}
	if sel, ok := models.AsOrganizationSelectionRequired(err); ok {
		if h.log != nil {
			h.log.Error(err)
		}
		writeJSONMessage(w, sel.Body(), http.StatusConflict)
		return
	}
	code := errors.GetCode(err)
	status := providerStatusOrInternal(err)
	if status == http.StatusInternalServerError && fallback != 0 {
		status = fallback
	}
	// Organization-scope and local "credential not found" errors keep their
	// MeshKit body and real status. Other failures are wrapped in the
	// handler-specific error, as they were before this mapping existed.
	if code == models.ErrResultNotFoundCode || code == models.ErrNoQualifyingOrganizationCode || code == models.ErrOrganizationMembershipCode || status == http.StatusNotFound {
		if h.log != nil {
			h.log.Error(err)
		}
		if code == models.ErrResultNotFoundCode || status == http.StatusNotFound {
			status = http.StatusNotFound
		}
		writeMeshkitError(w, err, status)
		return
	}
	out := err
	if wrap != nil {
		out = wrap(err)
	}
	if h.log != nil {
		h.log.Error(out)
	}
	writeMeshkitError(w, out, status)
}

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

	if err := models.RejectUnsafeKubeCredential(credential.Secret); err != nil {
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}

	// SaveUserCredential takes a token, not the request, so the saved
	// selection is read from the local persister via credential.UserId.
	createdCredential, err := provider.SaveUserCredential(token, &credential)
	if err != nil {
		h.writeCredentialProviderError(w, err, ErrSaveUserCredential, 0)
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
	credential, statusCode, err := provider.GetCredentialByID(token, providerUserID(user), credentialID)
	if err != nil {
		h.writeCredentialProviderError(w, err, ErrGetUserCredential, statusCode)
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

func (h *Handler) GetUserCredentials(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
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
	// The generated client sends the canonical camelCase pageSize; keep
	// page_size as a legacy fallback so older callers keep working.
	pageSizeParam := q.Get("pageSize")
	if pageSizeParam == "" {
		pageSizeParam = q.Get("page_size")
	}
	pageSize, _ := strconv.Atoi(pageSizeParam)

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

	// The remote provider reads the saved selection off the session preference
	// (and auto-selects a sole organization when none is saved). Put it on the
	// request so GetUserCredentials can see the same Preference the middleware
	// already loaded, including a selection the handler just mutated.
	if prefObj != nil {
		req = req.WithContext(context.WithValue(req.Context(), models.PerfObjCtxKey, prefObj))
	}

	credentialsPage, err := provider.GetUserCredentials(req, user.ID.String(), page, pageSize, search, order)
	if err != nil {
		h.writeCredentialProviderError(w, err, ErrGetUserCredential, 0)
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

func (h *Handler) UpdateUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
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

	if err := models.RejectUnsafeKubeCredential(credential.Secret); err != nil {
		writeMeshkitError(w, err, http.StatusBadRequest)
		return
	}

	req = sessionRequest(req, prefObj, user)
	_, err = provider.UpdateUserCredential(req, credential)
	if err != nil {
		h.writeCredentialProviderError(w, err, ErrUpdateUserCredential, 0)
		return
	}

	if h.log != nil {
		h.log.Info("Credential updated.")
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteUserCredential(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
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

	req = sessionRequest(req, prefObj, user)

	_, err = provider.DeleteUserCredential(req, credentialID)
	if err != nil {
		h.writeCredentialProviderError(w, err, ErrDeleteUserCredential, 0)
		return
	}

	if h.log != nil {
		h.log.Info("Credential deleted.")
	}
	w.WriteHeader(http.StatusOK)
}
