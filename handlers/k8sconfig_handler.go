package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/sirupsen/logrus"
)

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

	result := map[string]interface{}{
		"inClusterConfig": (inClusterConfig != ""),
	}

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

		result["context"] = ccfg.CurrentContext
		k8sServer, ok := ccfg.Clusters[ccfg.CurrentContext]
		if ok {
			result["server"] = k8sServer.Server
		}
	}
	session.Values["k8sInCluster"] = inClusterConfig

	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(result)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}
