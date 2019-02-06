package handlers

import (
	"context"
	"net/http"

	"github.com/layer5io/meshery/meshes/istio"
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
		// if !ok && contextNameI == nil {
		// 	logrus.Error("unable to get session data")
		// 	http.Error(w, "unable to get user data", http.StatusUnauthorized)
		// 	return
		// }
		logrus.Debugf("session values: %v", session.Values)
		k8sConfigBytesI, ok := session.Values["k8sConfig"]
		if !ok || k8sConfigBytesI == nil {
			logrus.Error("no valid kubernetes config found")
			http.Error(w, `No valid kubernetes config found. Please go <a href="/play/dashboard">here</a> to upload your config file and try again.`, http.StatusBadRequest)
			return
		}
		k8sConfigBytes, _ := k8sConfigBytesI.([]byte)

		meshClient, err := istio.CreateIstioClientWithK8SConfig(ctx, k8sConfigBytes, contextName)
		if err != nil {
			logrus.Fatalf("Error creating an istio client: %v", err)
			logrus.Errorf("error creating an istio client: %v", err)
			http.Error(w, "unable to create an istio client", http.StatusBadRequest)
			return
		}

		err = meshClient.ApplyRule(ctx, opName, userName, namespace)
		if err != nil {
			logrus.Error(err)
			http.Error(w, "there was an error creating the services", http.StatusInternalServerError)
			return
		}
	}
}
