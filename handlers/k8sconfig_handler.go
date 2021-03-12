//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"os"

	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// K8SConfigHandler is used for persisting kubernetes config and context info
func (h *Handler) K8SConfigHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	// if req.Method != http.MethodPost && req.Method != http.MethodDelete {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }

	if req.Method == http.MethodPost {
		h.addK8SConfig(user, prefObj, w, req, provider)
		return
	}
	if req.Method == http.MethodDelete {
		h.deleteK8SConfig(user, prefObj, w, req, provider)
		return
	}
}

func (h *Handler) addK8SConfig(user *models.User, prefObj *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
	_ = req.ParseMultipartForm(1 << 20)

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
			http.Error(w, "Unable to get Kubernetes config file", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = k8sfile.Close()
		}()
		k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Errorf("error reading config: %v", err)
			http.Error(w, "Unable to read the Kubernetes config file, please try again", http.StatusBadRequest)
			return
		}

		contextName = req.FormValue("contextName")

		ccfg, err := clientcmd.Load(k8sConfigBytes)
		if err != nil {
			logrus.Errorf("error parsing k8s config: %v", err)
			http.Error(w, "Given file is not a valid Kubernetes config file, please try again", http.StatusBadRequest)
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
	prefObj.K8SConfig = kc

	var err error
	prefObj.K8SConfig.ServerVersion, err = helpers.FetchKubernetesVersion(kc.Config, kc.ContextName)
	if err != nil {
		http.Error(w, "unable to ping the Kubernetes server", http.StatusInternalServerError)
		return
	}

	prefObj.K8SConfig.Nodes, err = helpers.FetchKubernetesNodes(kc.Config, kc.ContextName)
	if err != nil {
		http.Error(w, "unable to fetch nodes metadata from the Kubernetes server", http.StatusInternalServerError)
		return
	}

	if err = provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	kc.Config = nil

	if err = json.NewEncoder(w).Encode(kc); err != nil {
		logrus.Errorf("error marshaling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) deleteK8SConfig(user *models.User, prefObj *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
	prefObj.K8SConfig = nil
	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// GetContextsFromK8SConfig returns the context list for a given k8s config
func (h *Handler) GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request) {
	// if req.Method != http.MethodPost {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	_ = req.ParseMultipartForm(1 << 20)
	var k8sConfigBytes []byte

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Errorf("error getting k8s file: %v", err)
		http.Error(w, "Unable to get Kubernetes config file", http.StatusBadRequest)
		return
	}
	defer func() {
		_ = k8sfile.Close()
	}()
	k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
	if err != nil {
		logrus.Errorf("error reading config: %v", err)
		http.Error(w, "Unable to read the Kubernetes config file, please try again", http.StatusBadRequest)
		return
	}

	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		logrus.Errorf("error parsing k8s config: %v", err)
		http.Error(w, "Given file is not a valid Kubernetes config file, please try again", http.StatusBadRequest)
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
		logrus.Errorf("error marshaling data: %v", err)
		http.Error(w, "unable to retrieve the requested data", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) loadInClusterK8SConfig() (*models.K8SConfig, error) {
	// try to load k8s config from incluster config
	_, err := rest.InClusterConfig()
	if err != nil {
		err = errors.Wrap(err, "error parsing incluster k8s config")
		logrus.Error(err)
		return nil, err
	}
	return &models.K8SConfig{
		InClusterConfig: true,
		// ContextName:       ccfg.CurrentContext,
		ClusterConfigured: true,
	}, nil
}

func (h *Handler) loadK8SConfigFromDisk() (*models.K8SConfig, error) {
	// // try to load k8s config from local disk
	configFile := path.Join(h.config.KubeConfigFolder, "config") // is it ok to hardcode the name 'config'?
	if _, err := os.Stat(configFile); err != nil {
		// if os.IsNotExist(err) {
		// 	logrus.Warnf("%s location does not exist", configFile)
		// 	return nil
		// } else {
		err = errors.Wrapf(err, "unable to open file: %s", configFile)
		logrus.Error(err)
		return nil, err
		// }
	}
	k8sConfigBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		err = errors.Wrapf(err, "error reading file: %s", configFile)
		logrus.Error(err)
		return nil, err
	}
	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		err = errors.Wrap(err, "error parsing k8s config")
		logrus.Error(err)
		return nil, err
	}
	return &models.K8SConfig{
		InClusterConfig:   false,
		Config:            k8sConfigBytes,
		ContextName:       ccfg.CurrentContext,
		ClusterConfigured: true,
	}, nil
}

// ATM used only in the SessionSyncHandler
func (h *Handler) checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req *http.Request, user *models.User, prefObj *models.Preference, provider models.Provider) error {
	if prefObj == nil {
		prefObj = &models.Preference{
			AnonymousUsageStats:  true,
			AnonymousPerfResults: true,
		}
	}
	if prefObj.K8SConfig == nil {
		kc, err := h.loadK8SConfigFromDisk()
		if err != nil {
			kc, err = h.loadInClusterK8SConfig()
			if err != nil {
				return err
			}
		}
		prefObj.K8SConfig = kc
		err = provider.RecordPreferences(req, user.UserID, prefObj)
		if err != nil {
			err = errors.Wrapf(err, "unable to persist k8s config")
			logrus.Error(err)
			return err
		}
	}
	return nil
}

// KubernetesPingHandler - fetches server version to simulate ping
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if prefObj.K8SConfig == nil {
		_, _ = w.Write([]byte("[]"))
		return
	}

	version, err := helpers.FetchKubernetesVersion(prefObj.K8SConfig.Config, prefObj.K8SConfig.ContextName)
	if err != nil {
		err = errors.Wrap(err, "unable to ping Kubernetes")
		logrus.Error(err)
		http.Error(w, "unable to ping Kubernetes", http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(map[string]string{
		"server_version": version,
	}); err != nil {
		err = errors.Wrap(err, "unable to marshal the payload")
		logrus.Error(err)
		http.Error(w, "unable to marshal the payload", http.StatusInternalServerError)
		return
	}
}
