package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) DashboardHandler(ctx context.Context, meshClient meshes.MeshClient, w http.ResponseWriter, req *http.Request) {
	var user *models.User
	if h.config.ByPassAuth {
		userName := req.FormValue("user_name")
		if userName == "" {
			userName = "Test User"
		}
		user = h.setupSession(userName, req, w)
		if user == nil {
			return
		}
	} else {
		session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		if err != nil {
			logrus.Errorf("error getting session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}
		user, _ = session.Values["user"].(*models.User)
	}

	ops, err := meshClient.Operations(ctx)
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		http.Error(w, "unable to render the page", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"Ops":  ops,
		"Name": meshClient.MeshName(),
		"Url":  h.config.LoadTestURL,
		"User": user,
	}

	err = dashTempl.Execute(w, result)
	if err != nil {
		logrus.Errorf("error rendering the template for the dashboard: %v", err)
		http.Error(w, "unable to render the page", http.StatusInternalServerError)
		return
	}
}
