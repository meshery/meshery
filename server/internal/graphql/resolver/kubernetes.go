package resolver

import (
	"context"
	"database/sql"
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

func (r *Resolver) subscribeClusterResources(ctx context.Context, provider models.Provider, k8scontextIDs []string, namespace string) (<-chan *model.ClusterResources, error) {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	respChan := make(chan *model.ClusterResources)

	r.Config.DashboardK8sResourcesChan.SubscribeDashbordK8Resources(ch)

	go func() {
		// Remove this listener when the subscription ends so PublishDashboardK8sResources
		// stops tracking it (the channel is never closed, so removal is enough).
		defer r.Config.DashboardK8sResourcesChan.UnsubscribeDashboardK8sResources(ch)
		r.Log.Info("Initializing Cluster Resources subscription")
		for {
			select {
			case <-ch:
				clusterResources, err := r.getClusterResources(ctx, provider, k8scontextIDs, namespace)
				if err != nil {
					r.Log.Error(ErrClusterResourcesSubscription(err))
					break
				}
				respChan <- clusterResources
			case <-ctx.Done():
				r.Log.Info("Cluster Resources subscription stopped")
				return
			}
		}
	}()

	return respChan, nil
}

func (r *Resolver) getClusterResources(ctx context.Context, provider models.Provider, k8scontextIDs []string, namespace string) (*model.ClusterResources, error) {
	var cids []string
	query := `
		SELECT count(kind) as count, kind FROM kubernetes_resources kr LEFT JOIN kubernetes_resource_object_meta rom on kr.id = rom.id 
			WHERE kr.kind <> 'Namespace' AND rom.namespace = '' AND kr.cluster_id IN (?) GROUP BY kind
				UNION 
		SELECT count(kind) as count, kind FROM kubernetes_resources kr LEFT JOIN kubernetes_resource_object_meta rom on kr.id = rom.id 
			WHERE rom.namespace IN (?) AND kr.cluster_id IN (?) GROUP BY kind 
				UNION			
		SELECT count(kind) as count, kind FROM kubernetes_resources kr 
			WHERE kr.kind = 'Namespace' AND kr.cluster_id IN (?) GROUP BY kind`

	var rows *sql.Rows
	var err error
	if len(k8scontextIDs) == 0 {
		return nil, ErrEmptyCurrentK8sContext
	}

	cids = k8scontextIDs

	rows, err = provider.GetGenericPersister().Raw(query, cids, namespace, cids, cids).Rows()

	if err != nil {
		r.Log.Error(ErrGettingClusterResources(err))
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			r.Log.Error(closeErr)
		}
	}()

	resources := make([]*model.Resource, 0)
	for rows.Next() {
		var resource model.Resource
		err := rows.Scan(&resource.Count, &resource.Kind)
		if err != nil {
			r.Log.Error(ErrGettingClusterResources(err))
			return nil, err
		}
		resources = append(resources, &resource)
	}

	return &model.ClusterResources{
		Resources: resources,
	}, nil
}
