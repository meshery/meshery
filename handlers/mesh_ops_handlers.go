package handlers

import (
	"net/http"

	"github.com/layer5io/meshery/meshes"
	"github.com/sirupsen/logrus"
)

/*
func (h *Handler) K8SConfigHandler(w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1 << 20)

	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

	inClusterConfig := req.FormValue("inClusterConfig")
	logrus.Debugf("inClusterConfig: %s", inClusterConfig)

	var k8sConfigBytes []byte
	var contextName string

	if inClusterConfig == "" {
		// k8sfile, contextName
		k8sfile, _, err := req.FormFile("k8sfile")
		if err != nil {
			logrus.Errorf("error getting k8s file: %v", err)
			// http.Error(w, "error getting k8s file", http.StatusUnauthorized)
			// session.AddFlash("Unable to get kubernetes config file")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Unable to get kubernetes config file", http.StatusBadRequest)
			return
		}
		defer k8sfile.Close()
		k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Errorf("error reading config: %v", err)
			// http.Error(w, "unable to read config", http.StatusBadRequest)
			// session.AddFlash("Unable to read the kubernetes config file, please try again")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Unable to read the kubernetes config file, please try again", http.StatusBadRequest)
			return
		}

		contextName = req.FormValue("contextName")

		ccfg, err := clientcmd.Load(k8sConfigBytes)
		if err != nil {
			logrus.Errorf("error parsing k8s config: %v", err)
			// http.Error(w, "k8s config file not valid", http.StatusBadRequest)
			// session.AddFlash("Given file is not a valid kubernetes config file, please try again")
			// session.Save(req, w)
			// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
			http.Error(w, "Given file is not a valid kubernetes config file, please try again", http.StatusBadRequest)
			return
		}
		logrus.Debugf("current context: %s, contexts from config file: %v, clusters: %v", ccfg.CurrentContext, ccfg.Contexts, ccfg.Clusters)
		if contextName != "" {
			k8sCtx, ok := ccfg.Contexts[contextName]
			if !ok || k8sCtx == nil {
				logrus.Errorf("error specified context not found")
				// http.Error(w, "context not valid", http.StatusBadRequest)
				// session.AddFlash("Given context name is not valid, please try again with a valid value")
				// session.Save(req, w)
				// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
				http.Error(w, "Given context name is not valid, please try again with a valid value", http.StatusBadRequest)
				return
			}
			// all good, now set the current context to use
			ccfg.CurrentContext = contextName
		}

		// session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		// if err != nil {
		// 	logrus.Errorf("error getting session: %v", err)
		// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
		// 	return
		// }
		session.Values["k8sContext"] = contextName
		session.Values["k8sConfig"] = k8sConfigBytes
	}
	session.Values["k8sInCluster"] = inClusterConfig

	meshLocationURL := req.FormValue("meshLocationURL")
	logrus.Debugf("meshLocationURL: %s", meshLocationURL)
	session.Values["meshLocationURL"] = meshLocationURL

	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	mClient, err := meshes.CreateClient(req.Context(), k8sConfigBytes, contextName, meshLocationURL)
	if err != nil {
		logrus.Errorf("error creating a mesh client: %v", err)
		// http.Error(w, "unable to create an istio client", http.StatusInternalServerError)
		// session.AddFlash("Unable to connect to the mesh using the given kubernetes config file, please try again")
		// session.Save(req, w)
		// http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		http.Error(w, "Unable to connect to the mesh using the given kubernetes config file, please try again", http.StatusInternalServerError)
		return
	}
	defer mClient.Close()
	respOps, err := mClient.MClient.SupportedOperations(req.Context(), &meshes.SupportedOperationsRequest{})

	// meshClient, err := istio.CreateIstioClientWithK8SConfig(ctx, k8sConfigBytes, contextName)
	// if err != nil {
	// 	logrus.Errorf("error creating an istio client: %v", err)
	// 	// http.Error(w, "unable to create an istio client", http.StatusInternalServerError)
	// 	session.AddFlash("Unable to connect to Istio using the given kubernetes config file, please try again")
	// 	session.Save(req, w)
	// 	http.Redirect(w, req, "/play/dashboard", http.StatusFound)
	// 	return
	// }

	// ops, err := meshClient.Operations(ctx)
	if err != nil {
		logrus.Errorf("error getting operations for the mesh: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}

	meshNameOps, err := mClient.MClient.MeshName(req.Context(), &meshes.MeshNameRequest{})
	if err != nil {
		logrus.Errorf("error getting mesh name: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}

	result := map[string]interface{}{
		"Ops":  respOps.Ops,
		"Name": meshNameOps.GetName(),
		// "User": user,
	}

	// err = dashTempl.Execute(w, result)
	// if err != nil {
	// 	logrus.Errorf("error rendering the template for the dashboard: %v", err)
	// 	http.Error(w, "unable to render the page", http.StatusInternalServerError)
	// 	return
	// }
	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}
*/
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
