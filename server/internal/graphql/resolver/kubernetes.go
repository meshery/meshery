package resolver

import (
	"context"
	"strings"

	"github.com/meshery/meshery/server/internal/graphql/model"
	"github.com/meshery/meshery/server/models"
	meshkitKube "github.com/meshery/meshkit/utils/kubernetes"
	"github.com/meshery/meshkit/utils/kubernetes/describe"
)

func (r *Resolver) getAvailableNamespaces(ctx context.Context, provider models.Provider, k8sClusterIDs []string) ([]*model.NameSpace, error) {
	var cids []string
	if len(k8sClusterIDs) != 0 {
		cids = k8sClusterIDs
	} else {
		r.Log.Error(ErrEmptyCurrentK8sContext)
		return nil, ErrEmptyCurrentK8sContext
	}
	// resourceobjects := make([]meshsyncmodel.ResourceObjectMeta, 0)
	namespaces, err := model.SelectivelyFetchNamespaces(cids, provider)
	if err != nil {
		r.Log.Error(ErrGettingNamespace(err))
		return nil, err
	}

	modelnamespaces := make([]*model.NameSpace, 0)

	for _, ns := range namespaces {
		modelnamespaces = append(modelnamespaces, &model.NameSpace{
			Namespace: ns,
		})
	}
	return modelnamespaces, nil
}

func (r *Resolver) getKubectlDescribe(_ context.Context, name, kind, namespace string) (*model.KctlDescribeDetails, error) {
	var ResourceMap = map[string]describe.DescribeType{
		"pod":                       describe.Pod,
		"deployment":                describe.Deployment,
		"job":                       describe.Job,
		"cronjob":                   describe.CronJob,
		"statefulset":               describe.StatefulSet,
		"daemonset":                 describe.DaemonSet,
		"replicaset":                describe.ReplicaSet,
		"secret":                    describe.Secret,
		"service":                   describe.Service,
		"serviceaccount":            describe.ServiceAccount,
		"node":                      describe.Node,
		"limitrange":                describe.LimitRange,
		"resourcequota":             describe.ResourceQuota,
		"persistentvolume":          describe.PersistentVolume,
		"persistentvolumeclaim":     describe.PersistentVolumeClaim,
		"namespace":                 describe.Namespace,
		"endpoints":                 describe.Endpoints,
		"configmap":                 describe.ConfigMap,
		"priorityclass":             describe.PriorityClass,
		"ingress":                   describe.Ingress,
		"role":                      describe.Role,
		"clusterrole":               describe.ClusterRole,
		"rolebinding":               describe.RoleBinding,
		"clusterrolebinding":        describe.ClusterRoleBinding,
		"networkpolicy":             describe.NetworkPolicy,
		"replicationcontroller":     describe.ReplicationController,
		"certificatesigningrequest": describe.CertificateSigningRequest,
		"endpointslice":             describe.EndpointSlice,
	}

	options := describe.DescriberOptions{
		Name:      name,
		Namespace: namespace,
		Type:      ResourceMap[strings.ToLower(kind)],
	}

	client, err := meshkitKube.New([]byte(""))
	if err != nil {
		r.Log.Error(model.ErrMesheryClient(err))
		return nil, err
	}

	details, err := describe.Describe(client, options)
	if err != nil {
		r.Log.Error(ErrKubectlDescribe(err))
		return nil, err
	}

	return &model.KctlDescribeDetails{
		Describe: &details,
	}, nil
}
