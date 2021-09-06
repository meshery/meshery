package models

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/internal/sql"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"gopkg.in/yaml.v2"
)

type K8sContext struct {
	ID                string     `json:"id,omitempty"`
	Name              string     `json:"name,omitempty"`
	Auth              sql.Map    `json:"auth,omitempty"`
	Cluster           sql.Map    `json:"cluster,omitempty"`
	Server            string     `json:"server,omitempty"`
	Owner             *uuid.UUID `json:"owner,omitempty" gorm:"-"`
	CreatedBy         *uuid.UUID `json:"created_by,omitempty" gorm:"-"`
	IsCurrentContext  bool       `json:"is_current_context,omitempty"`
	MesheryInstanceID *uuid.UUID `json:"meshery_instance_id,omitempty"`

	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func K8sContextsFromKubeconfig(kubeconfig []byte, instanceID *uuid.UUID) []K8sContext {
	kcs := []K8sContext{}

	parsed, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return kcs
	}

	for name := range parsed.Contexts {
		kcs = append(kcs, NewK8sContext(
			name,
			parsed.Clusters,
			parsed.AuthInfos,
			parsed.CurrentContext == name,
			instanceID,
		))
	}

	return kcs
}

func NewK8sContext(
	contextName string,
	clusters map[string]*api.Cluster,
	users map[string]*api.AuthInfo,
	isCurrentContext bool,
	instanceID *uuid.UUID,
) K8sContext {
	ctx := K8sContext{
		Name: contextName,
		Cluster: map[string]interface{}{
			"cluster": clusters[contextName],
			"name":    contextName,
		},
		Auth: map[string]interface{}{
			"user": users[contextName],
			"name": contextName,
		},
		Server:            clusters[contextName].Server,
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
	config := map[string]interface{}{
		"apiVersion": "v1",
		"clusters": []map[string]interface{}{
			kc.Cluster,
		},
		"contexts": []map[string]interface{}{
			{
				"context": map[string]interface{}{
					"cluster": kc.Name,
					"user":    kc.Name,
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

	return yaml.Marshal(config)
}

// ToMapStringInterface takes in any data and converts it into map[string]interface{}
// if the data's underlying type has `MarshalJSON` and `UnmarshalJSON` implemented
func ToMapStringInterface(data interface{}) map[string]interface{} {
	mp := map[string]interface{}{}

	byt, err := json.Marshal(data)
	if err != nil {
		return mp
	}

	_ = json.Unmarshal(byt, &mp)

	return mp
}
