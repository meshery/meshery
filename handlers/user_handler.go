package handlers

import (
	"context"
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) UserHandler(ctx context.Context) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		if err != nil {
			logrus.Errorf("error getting session: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}

		var user *models.User
		user, _ = session.Values["user"].(*models.User)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			logrus.Errorf("error getting user data: %v", err)
			http.Error(w, "unable to get session", http.StatusUnauthorized)
			return
		}
		// json.Marshal(user)
	}
}
