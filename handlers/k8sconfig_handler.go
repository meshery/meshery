//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	// for GKE kube API authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"

	"os"

	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	ErrInvalidKubeConfig     = fmt.Errorf("given file is not a valid Kubernetes config file, please try again")
	ErrInvalidK8sContextName = fmt.Errorf("given context name is not valid, please try again with a valid value")
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

	var kc *models.K8SConfig
	var gerr error

	if inClusterConfig != "" {
		kc, gerr = h.setupK8sConfig(true, nil, "")
	} else {
		k8sfile, _, err := req.FormFile("k8sfile")
		if err != nil {
			logrus.Errorf("error getting k8s file: %v", err)
			http.Error(w, "Unable to get Kubernetes config file", http.StatusBadRequest)
			return
		}
		defer func() {
			_ = k8sfile.Close()
		}()
		k8sConfigBytes, err := ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Errorf("error reading config: %v", err)
			http.Error(w, "Unable to read the Kubernetes config file, please try again", http.StatusBadRequest)
			return
		}

		kc, gerr = h.setupK8sConfig(false, k8sConfigBytes, req.FormValue("contextName"))
	}

	if gerr != nil {
		status := http.StatusInternalServerError
		if gerr == ErrInvalidK8sContextName || gerr == ErrInvalidKubeConfig {
			status = http.StatusBadRequest
		}

		logrus.Errorf("failed to get K8s Config: %s", gerr)
		http.Error(w, gerr.Error(), status)
		return
	}

	prefObj.K8SConfig = kc
	if err := provider.RecordPreferences(req, user.UserID, prefObj); err != nil {
		logrus.Errorf("unable to save session: %v", err)
		http.Error(w, "unable to save session", http.StatusInternalServerError)
		return
	}

	kc.Config = nil
	if err := json.NewEncoder(w).Encode(kc); err != nil {
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
	return h.setupK8sConfig(true, nil, "")
}

func (h *Handler) loadK8SConfigFromDisk() (*models.K8SConfig, error) {
	// try to load k8s config from local disk
	configFile := path.Join(h.config.KubeConfigFolder, "config") // is it ok to hardcode the name 'config'?
	if _, err := os.Stat(configFile); err != nil {
		err = errors.Wrapf(err, "unable to open file: %s", configFile)
		logrus.Error(err)
		return nil, err
	}
	k8sConfigBytes, err := utils.ReadFileSource(fmt.Sprintf("file://%s", configFile))
	if err != nil {
		err = errors.Wrapf(err, "error reading file: %s", configFile)
		logrus.Error(err)
		return nil, err
	}

	return h.setupK8sConfig(false, []byte(k8sConfigBytes), "")
}

// ATM used only in the SessionSyncHandler
func (h *Handler) checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req *http.Request, user *models.User, prefObj *models.Preference, provider models.Provider) error {
	if prefObj == nil {
		prefObj = &models.Preference{
			AnonymousUsageStats:  true,
			AnonymousPerfResults: true,
		}
	}
	if prefObj.K8SConfig == nil || h.kubeclient == nil {
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

	version, err := h.kubeclient.KubeClient.ServerVersion()
	if err != nil {
		err = errors.Wrap(err, "unable to ping Kubernetes")
		logrus.Error(err)
		http.Error(w, "unable to ping Kubernetes", http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{
		"server_version": version.String(),
	}); err != nil {
		err = errors.Wrap(err, "unable to marshal the payload")
		logrus.Error(err)
		http.Error(w, "unable to marshal the payload", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) setupK8sConfig(inClusterConfig bool, k8sConfigBytes []byte, contextName string) (*models.K8SConfig, error) {
	kc := &models.K8SConfig{
		InClusterConfig:   inClusterConfig,
		Config:            k8sConfigBytes,
		ContextName:       contextName,
		ClusterConfigured: inClusterConfig,
	}

	mclient, err := mesherykube.New(k8sConfigBytes)
	if err != nil {
		return nil, err
	}

	kc.Server = mclient.RestConfig.Host

	version, err := mclient.KubeClient.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("unable to ping the Kubernetes server")
	}
	kc.ServerVersion = version.String()

	//kc.Nodes, err = helpers.FetchKubernetesNodes(kc.Config, kc.ContextName)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to fetch nodes metadata from the Kubernetes server")
	//}
	h.kubeclient = mclient
	return kc, nil
}
