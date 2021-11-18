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

	"github.com/layer5io/meshery/helpers"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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

// swagger:route POST /api/system/kubernetes SystemAPI idPostK8SConfig
// Handle POST request for Kubernetes Config
//
// Used to add kubernetes config to System
// responses:
// 	200: k8sConfigRespWrapper

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
			logrus.Error(ErrFormFile(err))
			http.Error(w, ErrFormFile(err).Error(), http.StatusBadRequest)
			return
		}
		defer func() {
			_ = k8sfile.Close()
		}()
		k8sConfigBytes, err := ioutil.ReadAll(k8sfile)
		if err != nil {
			logrus.Error(ErrReadConfig(err))
			http.Error(w, ErrReadConfig(err).Error(), http.StatusBadRequest)
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
		logrus.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}

	kc.Config = nil
	if err := json.NewEncoder(w).Encode(kc); err != nil {
		logrus.Error(ErrMarshal(err, "kubeconfig"))
		http.Error(w, ErrMarshal(err, "kubeconfig").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route DELETE /api/system/kubernetes SystemAPI idDeleteK8SConfig
// Handle DELETE request for Kubernetes Config
//
// Used to delete kubernetes config to System
// responses:
// 	200:

func (h *Handler) deleteK8SConfig(user *models.User, prefObj *models.Preference, w http.ResponseWriter, req *http.Request, provider models.Provider) {
	prefObj.K8SConfig = nil
	ctxID := "0" //To be replaced with actual context ID after multi context support
	go core.DeleteK8sWorkloads(ctxID)
	err := provider.RecordPreferences(req, user.UserID, prefObj)
	if err != nil {
		logrus.Error(ErrRecordPreferences(err))
		http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write([]byte("{}"))
}

// swagger:route POST /api/system/kubernetes/contexts SystemAPI idPostK8SContexts
// Handle POST requests for Kubernetes Context list
//
// Returns the context list for a given k8s config
// responses:
// 	200: k8sContextsRespWrapper

// GetContextsFromK8SConfig returns the context list for a given k8s config
func (h *Handler) GetContextsFromK8SConfig(w http.ResponseWriter, req *http.Request) {
	_ = req.ParseMultipartForm(1 << 20)
	var k8sConfigBytes []byte

	k8sfile, _, err := req.FormFile("k8sfile")
	if err != nil {
		logrus.Error(ErrFormFile(err))
		http.Error(w, ErrFormFile(err).Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = k8sfile.Close()
	}()
	k8sConfigBytes, err = ioutil.ReadAll(k8sfile)
	if err != nil {
		logrus.Error(ErrReadConfig(err))
		http.Error(w, ErrReadConfig(err).Error(), http.StatusBadRequest)
		return
	}

	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		logrus.Error(ErrLoadConfig(err))
		http.Error(w, ErrLoadConfig(err).Error(), http.StatusBadRequest)
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
		logrus.Error(ErrMarshal(err, "kube-context"))
		http.Error(w, ErrMarshal(err, "kube-context").Error(), http.StatusInternalServerError)
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
		h.log.Error(ErrOpenFile(configFile))
		return nil, ErrOpenFile(configFile)
	}
	k8sConfig, err := utils.ReadFileSource(fmt.Sprintf("file://%s", configFile))
	if err != nil {
		h.log.Error(err)
		return nil, err
	}

	ncfg, err := helpers.FlattenMinifyKubeConfig([]byte(k8sConfig))
	if err != nil {
		h.log.Error(ErrLoadConfig(err))
		return nil, ErrLoadConfig(err)
	}

	ccfg, err := clientcmd.Load(ncfg)
	if err != nil {
		h.log.Error(ErrLoadConfig(err))
		return nil, ErrLoadConfig(err)
	}

	return h.setupK8sConfig(false, ncfg, ccfg.CurrentContext)
}

// ATM used only in the SessionSyncHandler
func (h *Handler) checkIfK8SConfigExistsOrElseLoadFromDiskOrK8S(req *http.Request, user *models.User, prefObj *models.Preference, provider models.Provider) error {
	if prefObj == nil {
		prefObj = &models.Preference{
			AnonymousUsageStats:  true,
			AnonymousPerfResults: true,
		}
	}
	if prefObj.K8SConfig == nil || h.config.KubeClient == nil {
		kc, err := h.loadK8SConfigFromDisk()
		if err != nil {
			kc, err = h.loadInClusterK8SConfig()
			if err != nil {
				return err
			}
		}
		prefObj.K8SConfig = kc
		if viper.GetBool("SKIP_COMP_GEN") {
			ctxID := "0" // To be replaced after multi-context support
			go func(l logger.Handler, config []byte, ctx string) {
				err := registerK8sComponents(h.log, prefObj.K8SConfig.Config, ctxID)
				if err != nil {
					logrus.Error(err)
				}
			}(h.log, prefObj.K8SConfig.Config, ctxID)
		}
		err = provider.RecordPreferences(req, user.UserID, prefObj)
		if err != nil {
			logrus.Error(ErrRecordPreferences(err))
			return ErrRecordPreferences(err)
		}
	}
	return nil
}

// swagger:route GET /api/system/kubernetes/ping SystemAPI idGetKubernetesPing
// Handle GET request for Kubernetes ping
//
// Fetches server version to simulate ping
// responses:
// 	200:

// KubernetesPingHandler - fetches server version to simulate ping
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	if prefObj.K8SConfig == nil {
		_, _ = w.Write([]byte("[]"))
		return
	}

	version, err := h.config.KubeClient.KubeClient.ServerVersion()
	if err != nil {
		logrus.Error(ErrKubeVersion(err))
		http.Error(w, ErrKubeVersion(err).Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(map[string]string{
		"server_version": version.String(),
	}); err != nil {
		err = errors.Wrap(err, "unable to marshal the payload")
		logrus.Error(ErrMarshal(err, "kube-server-version"))
		http.Error(w, ErrMarshal(err, "kube-server-version").Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) setupK8sConfig(inClusterConfig bool, k8sConfigBytes []byte, contextName string) (*models.K8SConfig, error) {
	kc := &models.K8SConfig{
		InClusterConfig: inClusterConfig,
		Config:          k8sConfigBytes,
		ContextName:     contextName,
	}

	mclient, err := mesherykube.New(k8sConfigBytes)
	if err != nil {
		kc.ClusterConfigured = false
		return nil, err
	}

	kc.Server = mclient.RestConfig.Host

	version, err := mclient.KubeClient.ServerVersion()
	if err != nil {
		kc.ClusterConfigured = false
		return nil, ErrKubeVersion(err)
	}
	kc.ServerVersion = version.String()

	ccfg, err := clientcmd.Load(k8sConfigBytes)
	if err != nil {
		h.log.Error(ErrLoadConfig(err))
		return nil, ErrLoadConfig(err)
	}

	kc.Contexts = []models.K8SContext{}
	for name, ctx := range ccfg.Contexts {
		kc.Contexts = append(kc.Contexts, models.K8SContext{
			ContextName:      name,
			ClusterName:      ctx.Cluster,
			IsCurrentContext: name == contextName,
		})
	}

	//kc.Nodes, err = helpers.FetchKubernetesNodes(kc.Config, kc.ContextName)
	//if err != nil {
	//	return nil, fmt.Errorf("unable to fetch nodes metadata from the Kubernetes server")
	//}
	*h.config.KubeClient = *mclient
	kc.ClusterConfigured = true
	return kc, nil
}
func registerK8sComponents(l logger.Handler, config []byte, ctx string) error {
	l.Info("Starting to register k8s native components")
	man, err := core.GetK8Components(config, ctx)
	if err != nil {
		return err
	}
	if man == nil {
		l.Error(errors.New("Could not get k8s components"))
		return err
	}
	for i, def := range man.Definitions {
		var ord core.WorkloadCapability
		ord.Metadata = make(map[string]string)
		ord.Metadata["io.meshery.ctxid"] = ctx
		ord.Metadata["adapter.meshery.io/name"] = "kubernetes"
		ord.Host = "<none-local>"
		ord.OAMRefSchema = man.Schemas[i]

		var definition v1alpha1.WorkloadDefinition
		err := json.Unmarshal([]byte(def), &definition)
		if err != nil {
			return err
		}
		ord.OAMDefinition = definition
		content, err := json.Marshal(ord)
		if err != nil {
			return err
		}
		err = core.RegisterWorkload(content)
		if err != nil {
			return err
		}
	}
	l.Info("Registration of k8s native components completed")
	return nil
}
