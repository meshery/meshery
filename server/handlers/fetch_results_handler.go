// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	meshkiterrors "github.com/meshery/meshkit/errors"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
)

// FetchResultsHandler fetchs pages of results from Remote Provider and presents it to the UI
func (h *Handler) FetchResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	profileID := mux.Vars(req)["id"]
	if profileID == "" {
		profileID = mux.Vars(req)["performanceProfileId"]
	}

	err := req.ParseForm()
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrParseForm(err))
		}
		writeMeshkitError(w, ErrParseForm(err), http.StatusBadRequest)
		return
	}
	q := req.Form

	if p == nil {
		writeJSONError(w, "provider not initialized", http.StatusInternalServerError)
		return
	}

	tokenString, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		fetchTokenErr := ErrFetchToken(fmt.Errorf("token not found in request context"))
		if h.log != nil {
			h.log.Error(fetchTokenErr)
		}
		writeMeshkitError(w, fetchTokenErr, http.StatusInternalServerError)
		return
	}

	bdr, err := p.FetchResults(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), profileID)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrFetchResults(err))
		}
		writeMeshkitError(w, ErrFetchResults(err), providerStatusOrInternal(err))
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bdr)
}

// FetchAllResultsHandler fetchs pages of results from Remote Provider and presents it to the UI
func (h *Handler) FetchAllResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	err := req.ParseForm()
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrParseForm(err))
		}
		writeMeshkitError(w, ErrParseForm(err), http.StatusBadRequest)
		return
	}
	q := req.Form

	if p == nil {
		writeJSONError(w, "provider not initialized", http.StatusInternalServerError)
		return
	}

	tokenString, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		fetchTokenErr := ErrFetchToken(fmt.Errorf("token not found in request context"))
		if h.log != nil {
			h.log.Error(fetchTokenErr)
		}
		writeMeshkitError(w, fetchTokenErr, http.StatusInternalServerError)
		return
	}

	bdr, err := p.FetchAllResults(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("from"), q.Get("to"))
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrFetchResults(err))
		}
		writeMeshkitError(w, ErrFetchResults(err), providerStatusOrInternal(err))
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bdr)
}

// GetResultHandler gets an individual result from provider
func (h *Handler) GetResultHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	id := mux.Vars(req)["id"]
	if id == "" {
		if h.log != nil {
			h.log.Error(ErrQueryGet("id"))
		}
		writeMeshkitError(w, ErrMissingResultID(), http.StatusBadRequest)
		return
	}
	key, err := uuid.FromString(id)
	if err != nil || key == uuid.Nil {
		if h.log != nil {
			h.log.Error(ErrQueryGet("key"))
		}
		writeMeshkitError(w, models.ErrInvalidUUID(fmt.Errorf("invalid result id: %q", id)), http.StatusBadRequest)
		return
	}

	if p == nil {
		writeJSONError(w, "provider not initialized", http.StatusInternalServerError)
		return
	}

	tokenString, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		fetchTokenErr := ErrFetchToken(fmt.Errorf("token not found in request context"))
		if h.log != nil {
			h.log.Error(fetchTokenErr)
		}
		writeMeshkitError(w, fetchTokenErr, http.StatusInternalServerError)
		return
	}

	bdr, err := p.GetResult(tokenString, key)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrGetResult(err))
		}
		status := providerStatusOrInternal(err)
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode {
			status = http.StatusNotFound
		}
		writeMeshkitError(w, ErrGetResult(err), status)
		return
	}
	if bdr == nil {
		notFoundErr := ErrGetResult(fmt.Errorf("result not found for id: %s", key))
		if h.log != nil {
			h.log.Error(notFoundErr)
		}
		writeMeshkitError(w, notFoundErr, http.StatusNotFound)
		return
	}

	sp, err := bdr.ConvertToSpec(h.log)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrConvertToSpec(err))
		}
		writeMeshkitError(w, ErrConvertToSpec(err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/yaml")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="result_%s.yaml"`, bdr.ID))
	b, err := yaml.Marshal(sp)
	if err != nil {
		if h.log != nil {
			h.log.Error(models.ErrMarshal(err, "test result"))
		}
		// yaml.Marshal failed before we wrote the body; send JSON error with
		// the marshal failure context. Override content-type back to JSON.
		w.Header().Set("content-type", "application/json")
		writeMeshkitError(w, models.ErrMarshal(err, "test result"), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) FetchSmiResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	w.Header().Set("content-type", "application/json")
	err := req.ParseForm()
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrParseForm(err))
		}
		writeMeshkitError(w, ErrParseForm(err), http.StatusBadRequest)
		return
	}
	q := req.Form

	if p == nil {
		writeJSONError(w, "provider not initialized", http.StatusInternalServerError)
		return
	}

	bdr, err := p.FetchSmiResults(req, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"))
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrFetchSMIResults(err))
		}
		writeMeshkitError(w, ErrFetchSMIResults(err), providerStatusOrInternal(err))
		return
	}
	_, _ = w.Write(bdr)
}

// FetchSingleSmiResultHandler gets the result of single smi conformance test
func (h *Handler) FetchSingleSmiResultHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	w.Header().Set("content-type", "application/json")
	err := req.ParseForm()
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrParseForm(err))
		}
		writeMeshkitError(w, ErrParseForm(err), http.StatusBadRequest)
		return
	}
	q := req.Form

	if p == nil {
		writeJSONError(w, "provider not initialized", http.StatusInternalServerError)
		return
	}

	id := mux.Vars(req)["id"]
	key, err := uuid.FromString(id)
	if err != nil || key == uuid.Nil {
		if h.log != nil {
			h.log.Error(ErrQueryGet("key"))
		}
		writeMeshkitError(w, models.ErrInvalidUUID(fmt.Errorf("invalid result id: %q", id)), http.StatusBadRequest)
		return
	}

	bdr, err := p.FetchSmiResult(req, q.Get("page"), q.Get("pageSize"), q.Get("search"), q.Get("order"), key)
	if err != nil {
		if h.log != nil {
			h.log.Error(ErrFetchSMIResults(err))
		}
		status := providerStatusOrInternal(err)
		if errors.Is(err, gorm.ErrRecordNotFound) || meshkiterrors.GetCode(err) == models.ErrResultNotFoundCode {
			status = http.StatusNotFound
		}
		writeMeshkitError(w, ErrFetchSMIResults(err), status)
		return
	}
	if len(bdr) == 0 {
		notFoundErr := ErrFetchSMIResults(fmt.Errorf("SMI result not found for id: %s", key))
		if h.log != nil {
			h.log.Error(notFoundErr)
		}
		writeMeshkitError(w, notFoundErr, http.StatusNotFound)
		return
	}
	_, _ = w.Write(bdr)
}
