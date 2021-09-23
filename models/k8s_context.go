package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/helpers/utils"
	"github.com/layer5io/meshery/internal/sql"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sContext struct {
	ID                string     `json:"id,omitempty" yaml:"id,omitempty"`
	Name              string     `json:"name,omitempty" yaml:"name,omitempty"`
	Auth              sql.Map    `json:"auth,omitempty" yaml:"auth,omitempty"`
	Cluster           sql.Map    `json:"cluster,omitempty" yaml:"cluster,omitempty"`
	Server            string     `json:"server,omitempty" yaml:"server,omitempty"`
	Owner             *uuid.UUID `json:"owner,omitempty" gorm:"-" yaml:"owner,omitempty"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty" gorm:"-" yaml:"created_by,omitempty"`
	IsCurrentContext  bool       `json:"is_current_context,omitempty" yaml:"is_current_context,omitempty"`
	MesheryInstanceID *uuid.UUID `json:"meshery_instance_id,omitempty" yaml:"meshery_instance_id,omitempty"`

	UpdatedAt *time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

type InternalKubeConfig struct {
	ApiVersion     string                   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
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
		kcfg.CurrentContext == name,
		instanceID,
	)
}

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
		kcs = append(kcs, kcfg.K8sContext(name, instanceID))
	}

	return kcs
}

func NewK8sContext(
	contextName string,
	clusters map[string]interface{},
	users map[string]interface{},
	server string,
	isCurrentContext bool,
	instanceID *uuid.UUID,
) K8sContext {
	ctx := K8sContext{
		Name:              contextName,
		Cluster:           clusters,
		Auth:              users,
		Server:            server,
		MesheryInstanceID: instanceID,
		IsCurrentContext:  isCurrentContext,
	}

	ID, err := K8sContextGenerateID(ctx)
	if err != nil {
		return ctx
	}

	ctx.ID = ID

	return ctx
}

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

	_, err = h.KubeClient.ServerVersion()
	return err
}
