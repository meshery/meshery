package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

func (h *Handler) MeshOpsHandler(ctx context.Context) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		opName := req.PostFormValue("query")
		namespace := req.PostFormValue("namespace")
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
		// if !ok && contextNameI == nil {
		// 	logrus.Error("unable to get session data")
		// 	http.Error(w, "unable to get user data", http.StatusUnauthorized)
		// 	return
		// }
		logrus.Debugf("session values: %v", session.Values)
		k8sConfigBytesI, ok := session.Values["k8sConfig"]
		if inClusterConfig == "" && (!ok || k8sConfigBytesI == nil) {
			logrus.Error("no valid kubernetes config found")
			http.Error(w, `No valid kubernetes config found. Please go <a href="/play/dashboard">here</a> to upload your config file and try again.`, http.StatusBadRequest)
			return
		}
		k8sConfigBytes, _ := k8sConfigBytesI.([]byte)

		meshLocationURL := ""
		meshLocationURLI, ok := session.Values["meshLocationURL"]
		if ok && meshLocationURLI != nil {
			meshLocationURL, _ = meshLocationURLI.(string)
		} else {
			logrus.Error("no valid url for mesh adapter found")
			http.Error(w, `No valid url for mesh adapter found. Please go <a href="/play/dashboard">here</a> to the mesh adapter url and try again.`, http.StatusBadRequest)
			return
		}

		mClient, err := meshes.CreateClient(ctx, k8sConfigBytes, contextName, meshLocationURL)
		if err != nil {
			logrus.Errorf("error creating a mesh client: %v", err)
			http.Error(w, "unable to create a mesh client", http.StatusBadRequest)
			return
		}
		defer mClient.Close()

		_, err = mClient.MClient.ApplyOperation(ctx, &meshes.ApplyRuleRequest{
			OpName:    opName,
			Username:  userName,
			Namespace: namespace,
		}) //ApplyRule(ctx, opName, userName, namespace)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "there was an error creating the services", http.StatusInternalServerError)
			return
		}
	}
}
