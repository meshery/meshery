// Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/meshery/meshery/server/models"
	"gopkg.in/yaml.v2"
)

// FetchResultsHandler fetchs pages of results from Remote Provider and presents it to the UI
func (h *Handler) FetchResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	profileID := mux.Vars(req)["id"]

	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrParseForm(err))
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form

	tokenString := req.Context().Value(models.TokenCtxKey).(string)

	bdr, err := p.FetchResults(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), profileID)
	if err != nil {
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bdr)
}

// FetchAllResultsHandler fetchs pages of results from Remote Provider and presents it to the UI
func (h *Handler) FetchAllResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrParseForm(err))
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form

	tokenString := req.Context().Value(models.TokenCtxKey).(string)

	bdr, err := p.FetchAllResults(tokenString, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"), q.Get("from"), q.Get("to"))
	if err != nil {
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bdr)
}

// GetResultHandler gets an individual result from provider
func (h *Handler) GetResultHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	// if req.Method != http.MethodGet {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// TODO: may be force login if token not found?????
	id := mux.Vars(req)["id"]
	if id == "" {
		h.log.Error(ErrQueryGet("id"))
		http.Error(w, "please provide a result id", http.StatusBadRequest)
		return
	}
	key := uuid.FromStringOrNil(id)
	if key == uuid.Nil {
		h.log.Error(ErrQueryGet("key"))
		http.Error(w, "please provide a valid result id", http.StatusBadRequest)
		return
	}

	tokenString := req.Context().Value(models.TokenCtxKey).(string)

	bdr, err := p.GetResult(tokenString, key)
	if err != nil {
		h.log.Error(ErrGetResult(err))
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	sp, err := bdr.ConvertToSpec(h.log)
	if err != nil {
		h.log.Error(ErrConvertToSpec(err))
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/yaml")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="result_%s.yaml"`, bdr.ID))
	b, err := yaml.Marshal(sp)
	if err != nil {
		h.log.Error(models.ErrMarshal(err, "test result"))
		http.Error(w, "error while getting test result", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(b)
}

func (h *Handler) FetchSmiResultsHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	w.Header().Set("content-type", "application/json")
	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrParseForm(err))
		http.Error(w, ErrParseForm(err).Error(), http.StatusForbidden)
	}
	q := req.Form

	bdr, err := p.FetchSmiResults(req, q.Get("page"), q.Get("pagesize"), q.Get("search"), q.Get("order"))
	if err != nil {
		h.log.Error(ErrFetchSMIResults(err))
		http.Error(w, ErrFetchSMIResults(err).Error(), http.StatusInternalServerError)
	}
	_, _ = w.Write(bdr)
}

// FetchSingleSmiResultHandler gets the result of single smi conformance test
func (h *Handler) FetchSingleSmiResultHandler(w http.ResponseWriter, req *http.Request, _ *models.Preference, _ *models.User, p models.Provider) {
	w.Header().Set("content-type", "application/json")
	err := req.ParseForm()
	if err != nil {
		h.log.Error(ErrParseForm(err))
		http.Error(w, ErrParseForm(err).Error(), http.StatusForbidden)
	}
	q := req.Form
	id := mux.Vars(req)["id"]
	key := uuid.FromStringOrNil(id)
	if key == uuid.Nil {
		h.log.Error(ErrQueryGet("key"))
		http.Error(w, "please provide a valid result id", http.StatusBadRequest)
		return
	}
	bdr, err := p.FetchSmiResult(req, q.Get("page"), q.Get("pageSize"), q.Get("search"), q.Get("order"), key)
	if err != nil {
		h.log.Error(ErrFetchSMIResults(err))
		http.Error(w, ErrFetchSMIResults(err).Error(), http.StatusInternalServerError)
	}
	_, _ = w.Write(bdr)
}
