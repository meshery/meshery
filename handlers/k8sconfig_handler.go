package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"k8s.io/client-go/tools/clientcmd"

	"os"

	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (h *Handler) K8SConfigHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost && req.Method != http.MethodDelete {
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
		logrus.Warn("unable to read session from the session persister, starting with a new one")
	}

	if sessObj == nil {
		sessObj = &models.Session{}
	}
	if req.Method == http.MethodPost {
		h.addK8SConfig(user, sessObj, w, req)
		return
	}
	if req.Method == http.MethodDelete {
		h.deleteK8SConfig(user, sessObj, w, req)
		return
	}

}

func (h *Handler) addK8SConfig(user *models.User, sessObj *models.Session, w http.ResponseWriter, req *http.Request) {
	req.ParseMultipartForm(1 << 20)

	inClusterConfig := req.FormValue("inClusterConfig")
	logrus.Debugf("inClusterConfig: %s", inClusterConfig)

	var k8sConfigBytes []byte
	var contextName string
	kc := &models.K8SConfig{
		InClusterConfig: (inClusterConfig != ""),
	}

	if inClusterConfig == "" {
		k8sfile, _, err := req.FormFile("k8sfile")
		if err != nil {
			logrus.Errorf("error getting k8s file: %v", err)
			http.Error(w, "Unable to get kubernetes config file", http.StatusBadRequest)
			return
		}
		defer k8sfile.Close()
		k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Errorf("error reading config: %v", err)
			http.Error(w, "Unable to read the kubernetes config file, please try again", http.StatusBadRequest)
			return
		}

		contextName = req.FormValue("contextName")

		ccfg, err := clientcmd.Load(k8sConfigBytes)
		if err != nil {
			logrus.Errorf("error parsing k8s config: %v", err)
			http.Error(w, "Given file is not a valid kubernetes config file, please try again", http.StatusBadRequest)
			return
		}
		logrus.Debugf("current context: %s, contexts from config file: %v, clusters: %v", ccfg.CurrentContext, ccfg.Contexts, ccfg.Clusters)
		if contextName != "" {
			k8sCtx, ok := ccfg.Contexts[contextName]
			if !ok || k8sCtx == nil {
				logrus.Errorf("error specified context not found")
				http.Error(w, "Given context name is not valid, please try again with a valid value", http.StatusBadRequest)
				return
			}
			ccfg.CurrentContext = contextName
		}

		kc.Config = k8sConfigBytes
		kc.ContextName = ccfg.CurrentContext

		k8sContext, ok := ccfg.Contexts[ccfg.CurrentContext]
		if ok {
			k8sServer, ok := ccfg.Clusters[k8sContext.Cluster]
			if ok {
				kc.Server = k8sServer.Server
			}
		}
	}
	kc.ClusterConfigured = true
	sessObj.K8SConfig = kc
	err := h.config.SessionPersister.Write(user.UserId, sessObj)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	kc.Config = nil

	err = json.NewEncoder(w).Encode(kc)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteK8SConfig(user *models.User, sessObj *models.Session, w http.ResponseWriter, req *http.Request) {
	sessObj.K8SConfig = nil
	err := h.config.SessionPersister.Write(user.UserId, sessObj)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("{}"))
}

func (h *Handler) GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	req.ParseMultipartForm(1 << 20)
	var k8sConfigBytes []byte

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Errorf("error getting k8s file: %v", err)
		http.Error(w, "Unable to get kubernetes config file", http.StatusBadRequest)
		return
	}
	defer k8sfile.Close()
	k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
	if err != nil {
		logrus.Errorf("error reading config: %v", err)
		http.Error(w, "Unable to read the kubernetes config file, please try again", http.StatusBadRequest)
		return
	}

	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		logrus.Errorf("error parsing k8s config: %v", err)
		http.Error(w, "Given file is not a valid kubernetes config file, please try again", http.StatusBadRequest)
		return
	}

	contexts := []*models.K8SContext{}
	for contextName, contextVal := range ccfg.Contexts {
		ct := &models.K8SContext{
			ContextName:      contextName,
			ClusterName:      contextVal.Cluster,
			IsCurrentContext: (contextName == ccfg.CurrentContext),
		}
		contexts = append(contexts, ct)
	}

	err = json.NewEncoder(w).Encode(contexts)
	if err != nil {
		logrus.Errorf("error marshalling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

// ATM used only in the SessionSyncHandler
func (h *Handler) checkIfK8SConfigExistsOrElseLoadFromDisk(user *models.User, sessObj *models.Session) error {
	if sessObj == nil {
		sessObj = &models.Session{}
	}
	if sessObj.K8SConfig == nil {
		// try to load k8s config from local disk
		_, err := os.Stat(h.config.KubeConfigFolder)
		if err != nil {
			if os.IsNotExist(err) {
				logrus.Warnf("%s location does not exist", h.config.KubeConfigFolder)
				return nil
			} else {
				err = errors.Wrapf(err, "unable to open folder: %s", h.config.KubeConfigFolder)
				logrus.Error(err)
				return err
			}
		}
		configFile := path.Join(h.config.KubeConfigFolder, "config")
		_, err = os.Stat(configFile)
		if err != nil {
			if os.IsNotExist(err) {
				logrus.Warnf("%s location does not exist", configFile)
				return nil
			} else {
				err = errors.Wrapf(err, "unable to open file: %s", configFile)
				logrus.Error(err)
				return err
			}
		}
		k8sConfigBytes, err := ioutil.ReadFile(configFile) // is it ok to hardcode the name 'config'?
		if err != nil {
			err = errors.Wrapf(err, "error reading file: %s", configFile)
			logrus.Error(err)
			return err
		}
		ccfg, err := clientcmd.Load(k8sConfigBytes)
		if err != nil {
			err = errors.Wrap(err, "error parsing k8s config")
			logrus.Error(err)
			return err
		}
		kc := &models.K8SConfig{
			InClusterConfig:   false,
			Config:            k8sConfigBytes,
			ContextName:       ccfg.CurrentContext,
			ClusterConfigured: true,
		}
		sessObj.K8SConfig = kc
		err = h.config.SessionPersister.Write(user.UserId, sessObj)
		if err != nil {
			err = errors.Wrapf(err, "unable to persist k8s config")
			logrus.Error(err)
			return err
		}
	}
	return nil
}
