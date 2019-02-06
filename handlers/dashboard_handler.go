package handlers

import (
	"context"
	"io/ioutil"
	"net/http"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/layer5io/meshery/meshes/istio"
	"github.com/layer5io/meshery/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) DashboardHandler(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1 << 20)

	session, err := h.config.SessionStore.Get(req, h.config.SessionName)
	if err != nil {
		logrus.Errorf("error getting session: %v", err)
		http.Error(w, "unable to get session", http.StatusUnauthorized)
		return
	}

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
		// session, err := h.config.SessionStore.Get(req, h.config.SessionName)
		// if err != nil {
		// 	logrus.Errorf("error getting session: %v", err)
		// 	http.Error(w, "unable to get session", http.StatusUnauthorized)
		// 	return
		// }
		user, _ = session.Values["user"].(*models.User)
	}

	// k8sfile, contextName
	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Errorf("error getting k8s file: %v", err)
		// http.Error(w, "error getting k8s file", http.StatusUnauthorized)
		session.AddFlash("Unable to get kubernetes config file")
		session.Save(req, w)
		http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		return
	}
	defer k8sfile.Close()
	k8sConfigBytes, err := ioutil.ReadAll(k8sfile)
	if err != nil {
		logrus.Errorf("error reading config: %v", err)
		// http.Error(w, "unable to read config", http.StatusBadRequest)
		session.AddFlash("Unable to read the kubernetes config file, please try again")
		session.Save(req, w)
		http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		return
	}

	contextName := req.FormValue("contextName")

	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		logrus.Errorf("error parsing k8s config: %v", err)
		// http.Error(w, "k8s config file not valid", http.StatusBadRequest)
		session.AddFlash("Given file is not a valid kubernetes config file, please try again")
		session.Save(req, w)
		http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		return
	}
	logrus.Debugf("current context: %s, contexts from config file: %v, clusters: %v", ccfg.CurrentContext, ccfg.Contexts, ccfg.Clusters)
	if contextName != "" {
		k8sCtx, ok := ccfg.Contexts[contextName]
		if !ok || k8sCtx == nil {
			logrus.Errorf("error specified context not found")
			// http.Error(w, "context not valid", http.StatusBadRequest)
			session.AddFlash("Given context name is not valid, please try again with a valid value")
			session.Save(req, w)
			http.Redirect(w, req, "/play/dashboard", http.StatusFound)
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
	err = session.Save(req, w)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}
	meshClient, err := istio.CreateIstioClientWithK8SConfig(ctx, k8sConfigBytes, contextName)
	if err != nil {
		logrus.Errorf("error creating an istio client: %v", err)
		// http.Error(w, "unable to create an istio client", http.StatusInternalServerError)
		session.AddFlash("Unable to connect to Istio using the given kubernetes config file, please try again")
		session.Save(req, w)
		http.Redirect(w, req, "/play/dashboard", http.StatusFound)
		return
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
		"User": user,
	}

	err = dashTempl.Execute(w, result)
	if err != nil {
		logrus.Errorf("error rendering the template for the dashboard: %v", err)
		http.Error(w, "unable to render the page", http.StatusInternalServerError)
		return
	}
}
