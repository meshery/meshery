package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/layer5io/meshery/models"
)

func (h *Handler) GetAllContexts(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	q := req.URL.Query()

	vals, err := provider.GetK8sContexts(token, q.Get("page"), q.Get("pageSize"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(w, "failed to get contexts", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(vals); err != nil {
		http.Error(w, "failed to encode contexts", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetContext(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.URL.Query().Get("current") != "" {
	// 	context, ok := req.Context().Value(models.KubeContextKey).(*models.K8sContext)
	// 	if !ok || context == nil {
	// 		http.Error(w, "failed to get context", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	if err := json.NewEncoder(w).Encode(context); err != nil {
	// 		http.Error(w, "failed to encode context", http.StatusInternalServerError)
	// 		return
	// 	}

	// 	return
	// }

	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	val, err := provider.GetK8sContext(token, mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "failed to get context", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(val); err != nil {
		http.Error(w, "failed to encode context", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) DeleteContext(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		http.Error(w, "failed to get token", http.StatusInternalServerError)
		return
	}

	_, err := provider.DeleteK8sContext(token, mux.Vars(req)["id"])
	if err != nil {
		http.Error(w, "failed to delete context", http.StatusInternalServerError)
		return
	}
}

// func (h *Handler) GetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	token, ok := req.Context().Value(models.TokenCtxKey).(string)
// 	if !ok {
// 		http.Error(w, "failed to get token", http.StatusInternalServerError)
// 		return
// 	}

// 	val, err := h.GetCurrentContext(token, provider)
// 	if err != nil {
// 		http.Error(w, "failed to get current context", http.StatusInternalServerError)
// 		return
// 	}

// 	if err := json.NewEncoder(w).Encode(val); err != nil {
// 		http.Error(w, "failed to encode context", http.StatusInternalServerError)
// 		return
// 	}
// }

// func (h *Handler) SetCurrentContextHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
// 	token, ok := req.Context().Value(models.TokenCtxKey).(string)
// 	if !ok {
// 		http.Error(w, "failed to get token", http.StatusInternalServerError)
// 		return
// 	}

// 	_, err := provider.SetCurrentContext(token, mux.Vars(req)["id"])
// 	if err != nil {
// 		http.Error(w, "failed to set current context", http.StatusInternalServerError)
// 		return
// 	}
// }
