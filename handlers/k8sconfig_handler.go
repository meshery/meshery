//Package handlers :  collection of handlers (aka "HTTP middleware")
package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	// for GKE kube API authentication
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/models"
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
	token, ok := req.Context().Value(models.TokenCtxKey).(string)
	if !ok {
		logrus.Error(fmt.Errorf("failed to retrieve user token")) // TODO: Replace with meshkit errors
		http.Error(w, fmt.Errorf("failed to retrieve user token").Error(), http.StatusInternalServerError)
		return
	}

	_ = req.ParseMultipartForm(1 << 20)

	inClusterConfig := req.FormValue("inClusterConfig")
	logrus.Debugf("inClusterConfig: %s", inClusterConfig)

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

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(fmt.Errorf("failed to read INSTANCE_ID")) // TODO: Replace with meshkit errors
		http.Error(w, fmt.Errorf("failed to read INSTANCE_ID").Error(), http.StatusInternalServerError)
		return
	}

	contexts := models.K8sContextsFromKubeconfig(k8sConfigBytes, mid)

	for _, ctx := range contexts {
		_, err := provider.SaveK8sContext(token, ctx) // Ignore errors
		if err != nil {
			logrus.Error("failed to persist context")
		}
	}

	if err := json.NewEncoder(w).Encode(contexts); err != nil {
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
	// prefObj.K8SConfig = nil
	// err := provider.RecordPreferences(req, user.UserID, prefObj)
	// if err != nil {
	// 	logrus.Error(ErrRecordPreferences(err))
	// 	http.Error(w, ErrRecordPreferences(err).Error(), http.StatusInternalServerError)
	// 	return
	// }

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

	// Get meshery instance ID
	mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
	if !ok {
		logrus.Error(fmt.Errorf("failed to read INSTANCE_ID")) // TODO: Replace with meshkit errors
		http.Error(w, fmt.Errorf("failed to read INSTANCE_ID").Error(), http.StatusInternalServerError)
		return
	}

	contexts := models.K8sContextsFromKubeconfig(k8sConfigBytes, mid)

	err = json.NewEncoder(w).Encode(contexts)
	if err != nil {
		logrus.Error(ErrMarshal(err, "kube-context"))
		http.Error(w, ErrMarshal(err, "kube-context").Error(), http.StatusInternalServerError)
		return
	}
}

// swagger:route GET /api/system/kubernetes/ping SystemAPI idGetKubernetesPing
// Handle GET request for Kubernetes ping
//
// Fetches server version to simulate ping
// responses:
// 	200:

// KubernetesPingHandler - fetches server version to simulate ping
func (h *Handler) KubernetesPingHandler(w http.ResponseWriter, req *http.Request, prefObj *models.Preference, user *models.User, provider models.Provider) {
	kubeclient, ok := req.Context().Value(models.KubeHanderKey).(*mesherykube.Client)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "failed to get kube client for the user")
		return
	}

	version, err := kubeclient.KubeClient.ServerVersion()
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

func (h *Handler) GetCurrentContext(token string, prov models.Provider) (*models.K8sContext, error) {
	// Try to get current context
	cc, err := prov.GetCurrentContext(token)
	if err != nil {
		// No current context implies that this is the first time meshery is loading up

		// Get meshery instance ID
		mid, ok := viper.Get("INSTANCE_ID").(*uuid.UUID)
		if !ok {
			return nil, fmt.Errorf("not instance ID found") // TODO: Replace with meshkit error
		}

		// Attempt to get kubeconfig from the filesystem
		data, err := utils.ReadFileSource(fmt.Sprintf("file://%s", filepath.Join(h.config.KubeConfigFolder, "config")))
		if err != nil {
			// Could be an in-cluster deployment
			cfg, err := rest.InClusterConfig()
			if err == rest.ErrNotInCluster {
				return nil, err // TODO: Replace with meshkit error
			}

			ctxName := "in-cluster"

			cc = models.NewK8sContext(
				ctxName,
				map[string]*api.Cluster{
					ctxName: {
						CertificateAuthorityData: cfg.CAData,
						Server:                   cfg.Host,
					},
				},
				map[string]*api.AuthInfo{
					ctxName: {
						ClientCertificateData: cfg.CertData,
						ClientKeyData:         cfg.KeyData,
					},
				},
				true,
				mid,
			)
		} else {
			ctxs := models.K8sContextsFromKubeconfig([]byte(data), mid)

			// Persist the generated contexts
			for _, ctx := range ctxs {
				_, err := prov.SaveK8sContext(token, ctx)
				if err != nil {
					// TODO: log the error here
					continue
				}

				if ctx.IsCurrentContext {
					cc = ctx
				}
			}
		}

		// Set cc as the current context
		// fail silently - even if it fails we can go through the entire process in
		// in the next request
		// failing to persist should not block the flow
		prov.SetCurrentContext(token, cc.ID)
	}

	// TODO: Handle edge case when there is no valid context
	return &cc, nil
}
