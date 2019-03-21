package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func (h *Handler) MeshOpsHandler(w http.ResponseWriter, req *http.Request) {
	opName := req.PostFormValue("query")
	customBody := req.PostFormValue("customBody")
	namespace := req.PostFormValue("namespace")
	delete := req.PostFormValue("deleteOp")
	if namespace == "" {
		namespace = "default"
	}

	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Error("unable to get session data")
		http.Error(w, "unable to get user data", http.StatusUnauthorized)
		return
	}
	userNameI, ok := session.Values["user"]
	if !ok && userNameI == nil {
		logrus.Error("unable to get session data")
		http.Error(w, "unable to get user data", http.StatusUnauthorized)
		return
	}
	userName, _ := userNameI.(string)

	contextName := ""
	contextNameI, ok := session.Values["k8sContext"]
	if ok && contextNameI != nil {
		contextName, _ = contextNameI.(string)
	}

	inClusterConfig := ""
	inClusterConfigI, ok := session.Values["k8sInCluster"]
	if ok && contextNameI != nil {
		inClusterConfig, _ = inClusterConfigI.(string)
	}

	// logrus.Debugf("session values: %v", session.Values)
	k8sConfigBytesI, ok := session.Values["k8sConfig"]
	if inClusterConfig == "" && (!ok || k8sConfigBytesI == nil) {
		logrus.Error("no valid kubernetes config found")
		http.Error(w, `No valid kubernetes config found.`, http.StatusBadRequest)
		return
	}
	k8sConfigBytes, _ := k8sConfigBytesI.([]byte)

	meshLocationURL := ""
	meshLocationURLI, ok := session.Values["meshLocationURL"]
	if ok && meshLocationURLI != nil {
		meshLocationURL, _ = meshLocationURLI.(string)
	} else {
		logrus.Error("no valid url for mesh adapter found")
		http.Error(w, `No valid url for mesh adapter found.`, http.StatusBadRequest)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	if err != nil {
		logrus.Errorf("error creating a mesh client: %v", err)
		http.Error(w, "Unable to create a mesh client", http.StatusBadRequest)
		return
	}
	defer mClient.Close()

	_, err = mClient.MClient.ApplyOperation(req.Context(), &meshes.ApplyRuleRequest{
		OpName:     opName,
		Username:   userName,
		Namespace:  namespace,
		CustomBody: customBody,
		DeleteOp:   (delete != ""),
	})
	if err != nil {
		logrus.Error(err)
		http.Error(w, "There was an error applying the change", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{}"))
}
