package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

// FetchResultsHandler fetchs pages of results from SaaS and presents it to the UI
func (h *Handler) FetchResultsHandler(w http.ResponseWriter, req *http.Request, session *sessions.Session, _ *models.Preference, user *models.User, p models.Provider) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// TODO: may be force login if token not found?????

	err := req.ParseForm()
	if err != nil {
		logrus.Errorf("Error: unable to parse form: %v", err)
		http.Error(w, "unable to process the received data", http.StatusForbidden)
		return
	}
	q := req.Form

	bdr, err := p.FetchResults(req, q.Get("page"), q.Get("pageSize"), q.Get("search"), q.Get("order"))
	if err != nil {
		http.Error(w, "error while getting load test results", http.StatusInternalServerError)
		return
	}
	w.Header().Set("content-type", "application/json")
	_, _ = w.Write(bdr)
}
