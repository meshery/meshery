package handlers

import "net/http"

func (h *Handler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
