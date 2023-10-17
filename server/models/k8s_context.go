package models

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/gofrs/uuid"
	"github.com/layer5io/meshery/server/helpers/utils"
	"github.com/layer5io/meshery/server/internal/sql"
	"github.com/layer5io/meshkit/models/events"
	"github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
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
	DeploymentType     string     `json:"deployment_type,omitempty" yaml:"deployment_type,omitempty" default:"out_cluster"`
	Version            string     `json:"version,omitempty" yaml:"version,omitempty"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	CreatedAt          *time.Time `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	ConnectionID       string     `json:"connection_id,omitempty" yaml:"connection_id,omitempty"`
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

func (kcfg InternalKubeConfig) K8sContext(name string, instanceID *uuid.UUID) (K8sContext, string) {
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
	ctx, _ := NewK8sContext(contextName, clusters, users, server, instanceID)

	// Perform Ping test on the cluster
	if err := ctx.PingTest(); err != nil {
		return nil, err
	}

	// Get a kubernetes handler
	handler, err := ctx.GenerateKubeHandler()
	if err != nil {
		return nil, err
	}

	err = ctx.AssignVersion(handler)
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
func K8sContextsFromKubeconfig(provider Provider, userID string, eventChan *Broadcast, kubeconfig []byte, instanceID *uuid.UUID) []*K8sContext {
	kcs := []*K8sContext{}
	parsed, err := clientcmd.Load(kubeconfig)
	if err != nil {
		return kcs
	}

	userUUID := uuid.FromStringOrNil(userID)

	kcfg := InternalKubeConfig{}
	if err := yaml.Unmarshal(kubeconfig, &kcfg); err != nil {
		return kcs
	}

	for name := range parsed.Contexts {
		var msg string
		kc, _ := kcfg.K8sContext(name, instanceID)
		eventBuilder := events.NewEvent().ActedUpon(uuid.FromStringOrNil(kc.ConnectionID)).WithCategory("connection").WithAction("register").FromSystem(*instanceID).FromUser(userUUID)

		handler, err := kc.GenerateKubeHandler()
		if err != nil {
			msg = fmt.Sprintf("error generating kubernetes handler, skipping context %s: %v", err, kc.Name)

			event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Error connecting with kubernetes context at %s, skipping %s", kc.Server, kc.Name)).WithMetadata(map[string]interface{}{
				"error": err,
			}).Build()

			_ = provider.PersistEvent(event)
			eventChan.Publish(userUUID, event)

			logrus.Warnf(msg)
			continue
		}

		// Perform Ping test on the cluster
		if err := kc.PingTest(); err != nil {
			msg = fmt.Sprintf("unable to ping kubernetes context at %s, skipping context %s %v \n", kc.Server, kc.Name, err)
			event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Unable to ping kubernetes context at %s, skipping %s", kc.Server, kc.Name)).WithMetadata(map[string]interface{}{
				"error": err,
			}).Build()

			_ = provider.PersistEvent(event)
			eventChan.Publish(userUUID, event)

			logrus.Warn(msg)
			continue
		}

		if err := kc.AssignServerID(handler); err != nil {
			msg = fmt.Sprintf("could not retrieve kubernetes cluster ID, skipping context %s: %v", kc.Name, err)

			event := eventBuilder.WithSeverity(events.Error).WithDescription(fmt.Sprintf("Could not assign server id, skipping context %s", kc.Name)).WithMetadata(map[string]interface{}{
				"error": err,
			}).Build()

			_ = provider.PersistEvent(event)
			eventChan.Publish(userUUID, event)
			logrus.Warn(msg)
			continue
		}

		err = kc.AssignVersion(handler)
		if err != nil {
			msg = fmt.Sprintf("could not retrieve kubernetes version for context %s: %v ", kc.Name, err)
			event := eventBuilder.WithSeverity(events.Warning).WithDescription(fmt.Sprintf("Could not retrieve Kubernetes version for %s", kc.Name)).WithMetadata(map[string]interface{}{
				"error": err,
			}).Build()

			_ = provider.PersistEvent(event)
			eventChan.Publish(userUUID, event)

			logrus.Warnf(msg)
			kcs = append(kcs, &kc)
			continue
		}

		kcs = append(kcs, &kc)
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
) (K8sContext, string) {
	ctx := K8sContext{
		Name:              contextName,
		Cluster:           cluster,
		Auth:              user,
		Server:            server,
		MesheryInstanceID: instanceID,
	}

	ID, err := K8sContextGenerateID(ctx)
	if err != nil {
		return ctx, ""
	}

	ctx.ID = ID
	msg := fmt.Sprintf("Generated context: %s\n", ctx.Name)

	logrus.Infof(msg)

	return ctx, msg
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

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return []byte{}, ErrMarshal(err, "kube config")
	}
	return data, nil
}

func (kc *K8sContext) GenerateKubeHandler() (*kubernetes.Client, error) {
	cfg, err := kc.GenerateKubeConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.New(cfg)
}

func (kc *K8sContext) AssignVersion(handler *kubernetes.Client) error {
	res, err := handler.KubeClient.DiscoveryClient.ServerVersion()
	if err != nil {
		return ErrUnreachableKubeAPI(err, kc.Server)
	}

	kc.Version = res.GitVersion
	return nil
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
		return ErrUnreachableKubeAPI(res.Error(), kc.Server)
	}

	return nil
}

// AssignServerID will attempt to assign kubernetes
// server ID to the kubernetes context
func (kc *K8sContext) AssignServerID(handler *kubernetes.Client) error {
	// Get Kubernetes API server ID by querying the "kube-system" namespace uuid
	ksns, err := handler.KubeClient.CoreV1().Namespaces().Get(context.TODO(), "kube-system", v1.GetOptions{})
	if err != nil {
		return ErrUnreachableKubeAPI(err, kc.Server)
	}
	uid := ksns.ObjectMeta.GetUID()
	ksUUID := uuid.FromStringOrNil(string(uid))

	kc.KubernetesServerID = &ksUUID

	return nil
}

// FlushMeshSyncData will flush the meshsync data for the passed kubernetes contextID
func FlushMeshSyncData(ctx context.Context, k8sContext K8sContext, provider Provider, eventsChan *Broadcast, userID string, mesheryInstanceID *uuid.UUID) {
	ctxID := k8sContext.ID
	ctxUUID, _ := uuid.FromString(ctxID)
	userUUID, _ := uuid.FromString(userID)
	// Gets all the available kubernetes contexts

	ctxName := k8sContext.Name
	serverURL := k8sContext.Server
	k8sctxs, ok := ctx.Value(AllKubeClusterKey).([]K8sContext)
	if !ok || len(k8sctxs) == 0 {
		event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription("No Kubernetes context specified, please choose a context from context switcher").FromUser(userUUID).Build()
		err := provider.PersistEvent(event)
		if err != nil {
			logrus.Error(err)
		}

		eventsChan.Publish(userUUID, event)
		return
	}
	var sid string
	var refCount int
	// Gets the serverID for the passed contextID
	for _, k8ctx := range k8sctxs {
		if k8ctx.ID == ctxID && k8ctx.KubernetesServerID != nil {
			sid = k8ctx.KubernetesServerID.String()
			break
		}
	}
	// Counts the reference of the serverID
	// As multiple context can have same serverID
	for _, k8ctx := range k8sctxs {
		if k8ctx.KubernetesServerID.String() == sid {
			refCount++
		}
	}
	// If the reference count is 1 then only flush the meshsync data
	// because this means its the last contextID referring to that Kubernetes Server
	if refCount == 1 {
		if provider.GetGenericPersister() == nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(errors.New("meshery Database handler is not accessible to perform operations"), ctxName, serverURL),
			}).Build()
			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}
			eventsChan.Publish(userUUID, event)
			return
		}

		err := provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("objects").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesKeyValue{}).Error
		if err != nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(err, ctxName, serverURL),
			}).Build()
			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}

			eventsChan.Publish(userUUID, event)
			return
		}

		err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("objects").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceSpec{}).Error
		if err != nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(err, ctxName, serverURL),
			}).Build()
			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}

			eventsChan.Publish(userUUID, event)
			return
		}

		err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("objects").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceStatus{}).Error
		if err != nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(err, ctxName, serverURL),
			}).Build()

			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}

			eventsChan.Publish(userUUID, event)
			return
		}

		err = provider.GetGenericPersister().Where("id IN (?)", provider.GetGenericPersister().Table("objects").Select("id").Where("cluster_id=?", sid)).Delete(&meshsyncmodel.KubernetesResourceObjectMeta{}).Error
		if err != nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(err, ctxName, serverURL),
			}).Build()
			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}

			eventsChan.Publish(userUUID, event)
			return
		}

		err = provider.GetGenericPersister().Where("cluster_id = ?", sid).Delete(&meshsyncmodel.KubernetesResource{}).Error
		if err != nil {

			event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Error).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("Error flushing MeshSync data for %s", ctxName)).FromUser(userUUID).WithMetadata(map[string]interface{}{
				"error": ErrFlushMeshSyncData(err, ctxName, serverURL),
			}).Build()
			err := provider.PersistEvent(event)
			if err != nil {
				logrus.Error(err)
			}

			eventsChan.Publish(userUUID, event)
			return
		}

		event := events.NewEvent().ActedUpon(ctxUUID).FromSystem(*mesheryInstanceID).WithSeverity(events.Informational).WithCategory("meshsync").WithAction("flush").WithDescription(fmt.Sprintf("MeshSync data flushed successfully for context %s", ctxName)).FromUser(userUUID).Build()
		// Also add context name, as id is not helpful
		err = provider.PersistEvent(event)
		if err != nil {
			logrus.Error(err)
		}
		eventsChan.Publish(userUUID, event)
	}
}
