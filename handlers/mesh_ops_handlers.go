package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func (h *Handler) MeshOpsHandler(ctx context.Context, meshClient meshes.MeshClient) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		opName := req.PostFormValue("query")

		session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		if err != nil {
			logrus.Error("unable to get session data")
			http.Error(w, "unable to get user data", http.StatusUnauthorized)
			return
		}
		userName, _ := session.Values["user"].(string)

		err = meshClient.ApplyRule(ctx, opName, userName, "default")
		if err != nil {
			logrus.Error(err)
			http.Error(w, "there was an error creating the services", http.StatusInternalServerError)
			return
		}
	}
}
