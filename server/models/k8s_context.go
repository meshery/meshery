package models

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/sql"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sContext struct {
	ID                 string     `json:"id,omitempty" yaml:"id,omitempty"`
	Name               string     `json:"name,omitempty" yaml:"name,omitempty"`
	Auth               sql.Map    `json:"auth,omitempty" yaml:"auth,omitempty"`
	Cluster            sql.Map    `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	Server             string     `json:"server,omitempty" yaml:"server,omitempty"`
	Owner              *uuid.UUID `json:"owner,omitempty" gorm:"-" yaml:"owner,omitempty"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty" gorm:"-" yaml:"created_by,omitempty"`
	MesheryInstanceID  *uuid.UUID `json:"meshery_instance_id,omitempty" yaml:"meshery_instance_id,omitempty"`
	KubernetesServerID *uuid.UUID `json:"kubernetes_server_id,omitempty" yaml:"kubernetes_server_id,omitempty"`

	UpdatedAt *time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type InternalKubeConfig struct {
	APIVersion     string                   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind           string                   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Clusters       []map[string]interface{} `json:"clusters,omitempty" yaml:"clusters,omitempty"`
	Contexts       []map[string]interface{} `json:"contexts,omitempty" yaml:"contexts,omitempty"`
	CurrentContext string                   `json:"current-context,omitempty" yaml:"current-context,omitempty"`
	Preferences    map[string]interface{}   `json:"preferences,omitempty" yaml:"preferences,omitempty"`
	Users          []map[string]interface{} `json:"users,omitempty" yaml:"users,omitempty"`
}

func (kcfg InternalKubeConfig) K8sContext(name string, instanceID *uuid.UUID) K8sContext {
	cluster := map[string]interface{}{}
	user := map[string]interface{}{}
	context := map[string]interface{}{}

	// Find context data
	for _, ctx := range kcfg.Contexts {
		ctx = utils.RecursiveCastMapStringInterfaceToMapStringInterface(ctx)
		if ctx["name"] == name {
			context = ctx
			break
		}
	}

	ctxInfo, _ := context["context"].(map[string]interface{})

	// Find cluster data associated with the context
	clusterName := ctxInfo["cluster"]
	for _, cl := range kcfg.Clusters {
		cl = utils.RecursiveCastMapStringInterfaceToMapStringInterface(cl)
		if cl["name"] == clusterName {
			cluster = cl
			break
		}
	}

	clusterInfo, _ := cluster["cluster"].(map[string]interface{})
	server, _ := clusterInfo["server"].(string)

	// Find Auth data associated with the context
	userName := ctxInfo["user"]
	for _, u := range kcfg.Users {
		u = utils.RecursiveCastMapStringInterfaceToMapStringInterface(u)
		if u["name"] == userName {
			user = u
			break
		}
	}

	return NewK8sContext(
		name,
		cluster,
		user,
		server,
		instanceID,
	)
}

func NewK8sContextWithServerID(
	contextName string,
	clusters map[string]interface{},
	users map[string]interface{},
	server string,
	instanceID *uuid.UUID,
) (*K8sContext, error) {
	ctx := NewK8sContext(contextName, clusters, users, server, instanceID)

	// Perform Ping test on the cluster
	if err := ctx.PingTest(); err != nil {
		return nil, err
	}

	// Get a kubernetes handler
	handler, err := ctx.GenerateKubeHandler()
	if err != nil {
		return nil, err
	}

	// Get Kubernetes API server ID by querying the "kube-system" namespace uuid
	ksns, err := handler.KubeClient.CoreV1().Namespaces().Get(context.TODO(), "kube-system", v1.GetOptions{})
	if err != nil {
		return nil, err
	}
	uid := ksns.ObjectMeta.GetUID()
	ksUUID := uuid.FromStringOrNil(string(uid))

	ctx.KubernetesServerID = &ksUUID

	return &ctx, nil
}

// K8sContextsFromKubeconfig takes in a kubeconfig and meshery instance ID and generates
// kubernetes contexts from it
func K8sContextsFromKubeconfig(kubeconfig []byte, instanceID *uuid.UUID) []K8sContext {
	kcs := []K8sContext{}

	parsed, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return kcs
	}

	kcfg := InternalKubeConfig{}
	if err := yaml.Unmarshal(kubeconfig, &kcfg); err != nil {
		return kcs
	}

	for name := range parsed.Contexts {
		kc := kcfg.K8sContext(name, instanceID)
		if err := kc.AssignServerID(); err != nil {
			logrus.Warn("Skipping context: Reason => ", err)

			continue
		}

		kcs = append(kcs, kc)
	}

	return kcs
}

func NewK8sContextFromInClusterConfig(contextName string, instanceID *uuid.UUID) (*K8sContext, error) {
	const (
		tokenFile  = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		rootCAFile = "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	)
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	if len(host) == 0 || len(port) == 0 {
		return nil, ErrMesheryNotInCluster
	}

	token, err := os.ReadFile(tokenFile)
	if err != nil {
		return nil, err
	}

	server := "https://" + net.JoinHostPort(host, port)

	caData, err := os.ReadFile(rootCAFile)
	if err != nil {
		return nil, err
	}

	return NewK8sContextWithServerID(
		contextName,
		map[string]interface{}{
			"cluster": map[string]interface{}{
				"certificate-authority-data": base64.StdEncoding.EncodeToString(caData),
				"server":                     server,
			},
			"name": contextName,
		},
		map[string]interface{}{
			"user": map[string]interface{}{
				"token": string(token),
			},
			"name": contextName,
		},
		server,
		instanceID,
	)
}

// NewK8sContext takes in name of the context, cluster info of the contexts,
// auth info, server address and meshery instance ID and will return a K8sContext from it
//
// This function does NOT assigns kubernetes server ID to the context, either the ID
// can be assigned manually by invoking `AssignServerID` method or instead use
// `NewK8sContextWithServerID` to create a context
func NewK8sContext(
	contextName string,
	cluster map[string]interface{},
	user map[string]interface{},
	server string,
	instanceID *uuid.UUID,
) K8sContext {
	ctx := K8sContext{
		Name:              contextName,
		Cluster:           cluster,
		Auth:              user,
		Server:            server,
		MesheryInstanceID: instanceID,
	}

	ID, err := K8sContextGenerateID(ctx)
	if err != nil {
		return ctx
	}

	ctx.ID = ID

	logrus.Infof("Generated context: %s\n", ctx.Name)

	return ctx
}

// K8sContextGenerateID takes in a kubernetes context and generates an ID for it
//
// If the context remains the same, it is guaranteed that the ID will be same
func K8sContextGenerateID(kc K8sContext) (string, error) {
	data := map[string]interface{}{
		"cluster": kc.Cluster,
		"auth":    kc.Auth,
		"meshery": kc.MesheryInstanceID.String(),
		"name":    kc.Name,
	}

	byt, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	hash := md5.Sum(byt)

	return hex.EncodeToString(hash[:]), nil
}

// GenerateKubeConfig will generate a kubeconfig from the context object
// and will set the "current-context" to the current context's name
func (kc K8sContext) GenerateKubeConfig() ([]byte, error) {
	cfg := map[string]interface{}{
		"apiVersion": "v1",
		"clusters": []map[string]interface{}{
			kc.Cluster,
		},
		"contexts": []map[string]interface{}{
			{
				"context": map[string]interface{}{
					"cluster": kc.Cluster["name"],
					"user":    kc.Auth["name"],
				},
				"name": kc.Name,
			},
		},
		"current-context": kc.Name,
		"kind":            "Config",
		"users": []map[string]interface{}{
			kc.Auth,
		},
	}

	return yaml.Marshal(cfg)
}

func (kc K8sContext) GenerateKubeHandler() (*kubernetes.Client, error) {
	cfg, err := kc.GenerateKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.New(cfg)
}

// PingTest uses the k8scontext to to "ping" the kubernetes cluster
// if the return value is nil then the succeeds or else it has failed
func (kc K8sContext) PingTest() error {
	h, err := kc.GenerateKubeHandler()
	if err != nil {
		return err
	}

	res := h.KubeClient.DiscoveryClient.RESTClient().Get().RequestURI("/livez").Timeout(1 * time.Second).Do(context.TODO())
	if res.Error() != nil {
		return res.Error()
	}

	return nil
}

// AssignServerID will attempt to assign kubernetes
// server ID to the kubernetes context
func (kc *K8sContext) AssignServerID() error {
	// Perform Ping test on the cluster
	if err := kc.PingTest(); err != nil {
		return err
	}

	// Get a kubernetes handler
	handler, err := kc.GenerateKubeHandler()
	if err != nil {
		return err
	}

	// Get Kubernetes API server ID by querying the "kube-system" namespace uuid
	ksns, err := handler.KubeClient.CoreV1().Namespaces().Get(context.TODO(), "kube-system", v1.GetOptions{})
	if err != nil {
		return err
	}
	uid := ksns.ObjectMeta.GetUID()
	ksUUID := uuid.FromStringOrNil(string(uid))

	kc.KubernetesServerID = &ksUUID

	return nil
}
