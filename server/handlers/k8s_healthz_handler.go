package handlers

import "net/http"

func (h *Handler) K8sHealthzHandler(w http.ResponseWriter, r *http.Request) {

	if _, err := w.Write([]byte("ok")); err != nil {
		http.Error(w, "K8s Health Probe Error", http.StatusInternalServerError)
	}
}
