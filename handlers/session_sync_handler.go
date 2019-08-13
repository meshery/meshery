package handlers

import (
	"net/http"

	"encoding/json"

	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) SessionSyncHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	var user *models.User
	user, _ = session.Values["user"].(*models.User)

	// h.config.SessionPersister.Lock(user.UserId)
	// defer h.config.SessionPersister.Unlock(user.UserId)

	sessObj, err := h.config.SessionPersister.Read(user.UserId)
	if err != nil {
		logrus.Errorf("error retrieving user config data: %v", err)
		http.Error(w, "unable to get user config data", http.StatusInternalServerError)
		return
	}

	err = h.checkIfK8SConfigExistsOrElseLoadFromDisk(user, sessObj)
	// if err != nil {
	// // We can ignore the errors here. They are logged in the other method
	// }

	// clearing out the config just for displaying purposes
	if sessObj.K8SConfig != nil && sessObj.K8SConfig.Config != nil && len(sessObj.K8SConfig.Config) > 0 {
		sessObj.K8SConfig.Config = nil
	}

	err = json.NewEncoder(w).Encode(sessObj)
	if err != nil {
		logrus.Errorf("error marshalling user config data: %v", err)
		http.Error(w, "unable to process the request", http.StatusInternalServerError)
		return
	}
}
