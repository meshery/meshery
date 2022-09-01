package resolver

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/layer5io/meshery/server/handlers"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/models/pattern/core"
	meshkitKube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/kubernetes/describe"
)

func (r *Resolver) getAvailableNamespaces(ctx context.Context, provider models.Provider, k8sClusterIDs []string) ([]*model.NameSpace, error) {
	var cids []string
	if len(k8sClusterIDs) != 0 {
		k8sCtxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(k8sCtxs) == 0 {
			return nil, ErrMesheryClient(nil)
		}
		if len(k8sClusterIDs) == 1 && k8sClusterIDs[0] == "all" {
			for _, k8sContext := range k8sCtxs {
				if k8sContext.KubernetesServerID != nil {
					clusterID := k8sContext.KubernetesServerID.String()
					cids = append(cids, clusterID)
				}
			}
		} else {
			cids = k8sClusterIDs
		}
	} else { //This is a fallback
		k8sctxs, ok := ctx.Value(models.KubeClustersKey).([]models.K8sContext)
		if !ok || len(k8sctxs) == 0 {
			r.Log.Error(ErrEmptyCurrentK8sContext)
			return nil, ErrEmptyCurrentK8sContext
		}
		for _, context := range k8sctxs {
			if context.KubernetesServerID == nil {
				r.Log.Error(ErrEmptyCurrentK8sContext)
				return nil, ErrEmptyCurrentK8sContext
			}
			cids = append(cids, context.KubernetesServerID.String())
		}
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

// getWorkloads return workloads
func (r *Resolver) getWorkloads(ctx context.Context, name, id *string, trim *bool) (res []*model.OAMCapability, err error) {
	if name != nil && *name != "" && id != nil && *id != "" {
		workload := core.GetWorkloadByID(*name, *id)
		// If trim is set to true then remove the schema from the response
		if trim != nil && *trim {
			workload.OAMRefSchema = ""
		}
		res = append(res, &model.OAMCapability{
			OamDefinition: workload.OAMDefinition,
			OamRefSchema:  &workload.OAMRefSchema,
			Host:          &workload.Host,
			ID:            &workload.ID,
			Restricted:    &workload.Restricted,
			Metadata:      workload.Metadata,
		})
	} else if name != nil && *name != "" {
		workloads := core.GetWorkload(*name)
		for _, w := range workloads {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				w.OAMRefSchema = ""
			}
			res = append(res, &model.OAMCapability{
				OamDefinition: w.OAMDefinition,
				OamRefSchema:  &w.OAMRefSchema,
				Host:          &w.Host,
				ID:            &w.ID,
				Restricted:    &w.Restricted,
				Metadata:      w.Metadata,
			})
		}
	} else {
		// return all the workloads
		workloads := core.GetWorkloads()
		for _, w := range workloads {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				w.OAMRefSchema = ""
			}
			res = append(res, &model.OAMCapability{
				OamDefinition: w.OAMDefinition,
				OamRefSchema:  &w.OAMRefSchema,
				Host:          &w.Host,
				ID:            &w.ID,
				Restricted:    &w.Restricted,
				Metadata:      w.Metadata,
			})
		}
	}
	return
}

// getTraits return traits
func (r *Resolver) getTraits(ctx context.Context, name, id *string, trim *bool) (res []*model.OAMCapability, err error) {
	if name != nil && *name != "" && id != nil && *id != "" {
		trait := core.GetTraitByID(*name, *id)
		// If trim is set to true then remove the schema from the response
		if trim != nil && *trim {
			trait.OAMRefSchema = ""
		}
		res = append(res, &model.OAMCapability{
			OamDefinition: trait.OAMDefinition,
			OamRefSchema:  &trait.OAMRefSchema,
			Host:          &trait.Host,
			ID:            &trait.ID,
			Restricted:    &trait.Restricted,
			Metadata:      trait.Metadata,
		})
	} else if name != nil && *name != "" {
		traits := core.GetWorkload(*name)
		for _, t := range traits {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				t.OAMRefSchema = ""
			}
			res = append(res, &model.OAMCapability{
				OamDefinition: t.OAMDefinition,
				OamRefSchema:  &t.OAMRefSchema,
				Host:          &t.Host,
				ID:            &t.ID,
				Restricted:    &t.Restricted,
				Metadata:      t.Metadata,
			})
		}
	} else {
		// return all the traits
		traits := core.GetWorkloads()
		for _, t := range traits {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				t.OAMRefSchema = ""
			}
			res = append(res, &model.OAMCapability{
				OamDefinition: t.OAMDefinition,
				OamRefSchema:  &t.OAMRefSchema,
				Host:          &t.Host,
				ID:            &t.ID,
				Restricted:    &t.Restricted,
				Metadata:      t.Metadata,
			})
		}
	}

	return
}

// getScopes return scopes
func (r *Resolver) getScopes(ctx context.Context, name, id *string, trim *bool) (res []*model.OAMCapability, err error) {
	if name != nil && *name != "" && id != nil && *id != "" {
		scope := core.GetScopeByID(*name, *id)
		// If trim is set to true then remove the schema from the response
		if trim != nil && *trim {
			scope.OAMRefSchema = ""
		}
		res = append(res, &model.OAMCapability{
			OamDefinition: scope.OAMDefinition,
			OamRefSchema:  &scope.OAMRefSchema,
			Host:          &scope.Host,
			ID:            &scope.ID,
			Restricted:    &scope.Restricted,
			Metadata:      scope.Metadata,
		})
	} else if name != nil && *name != "" {
		scopes := core.GetWorkload(*name)
		for _, s := range scopes {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				s.OAMRefSchema = ""
			}

			res = append(res, &model.OAMCapability{
				OamDefinition: s.OAMDefinition,
				OamRefSchema:  &s.OAMRefSchema,
				Host:          &s.Host,
				ID:            &s.ID,
				Restricted:    &s.Restricted,
				Metadata:      s.Metadata,
			})
		}
	} else {
		// return all the scopes
		scopes := core.GetWorkloads()
		for _, s := range scopes {
			// If trim is set to true then remove the schema from the response
			if trim != nil && *trim {
				s.OAMRefSchema = ""
			}
			res = append(res, &model.OAMCapability{
				OamDefinition: s.OAMDefinition,
				OamRefSchema:  &s.OAMRefSchema,
				Host:          &s.Host,
				ID:            &s.ID,
				Restricted:    &s.Restricted,
				Metadata:      s.Metadata,
			})
		}
	}
	return
}

func (r *Resolver) getKubectlDescribe(ctx context.Context, name string, kind string, namespace string) (*model.KctlDescribeDetails, error) {
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
		r.Log.Error(ErrMesheryClient(err))
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
	respChan := make(chan *model.ClusterResources)

	r.Config.DashboardK8sResourcesChan.SubscribeDashbordK8Resources(ch)

	go func() {
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
		SELECT count(kind) as count, kind FROM objects o LEFT JOIN resource_object_meta rom on o.id = rom.id 
			WHERE o.kind <> 'Namespace' AND rom.namespace = '' AND o.cluster_id IN (?) GROUP BY kind
				UNION 
		SELECT count(kind) as count, kind FROM objects o LEFT JOIN resource_object_meta rom on o.id = rom.id 
			WHERE rom.namespace IN (?) AND o.cluster_id IN (?) GROUP BY kind 
				UNION			
		SELECT count(kind) as count, kind FROM objects o 
			WHERE o.kind = 'Namespace' AND o.cluster_id IN (?) GROUP BY kind`

	var rows *sql.Rows
	var err error
	k8sCtxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
	if !ok || len(k8sCtxs) == 0 {
		return nil, ErrMesheryClient(nil)
	}

	if len(k8scontextIDs) == 1 && k8scontextIDs[0] == "all" {
		for _, k8sContext := range k8sCtxs {
			if k8sContext.KubernetesServerID != nil {
				clusterID := k8sContext.KubernetesServerID.String()
				cids = append(cids, clusterID)
			}
		}
	} else {
		cids = k8scontextIDs
	}

	rows, err = provider.GetGenericPersister().Raw(query, cids, namespace, cids, cids).Rows()

	if err != nil {
		r.Log.Error(ErrGettingClusterResources(err))
	}

	defer rows.Close()

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

func (r *Resolver) subscribeK8sContexts(ctx context.Context, provider models.Provider, selector model.PageFilter) (<-chan *model.K8sContextsPage, error) {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	contextsChan := make(chan *model.K8sContextsPage)

	r.Config.K8scontextChannel.SubscribeContext(ch)
	r.Log.Info("K8s context subscription started")

	go func() {
		for {
			select {
			case <-ch:
				contexts, err := r.getK8sContexts(ctx, provider, selector)
				if err != nil {
					r.Log.Error(ErrK8sContextSubscription(err))
					break
				}
				contextsChan <- contexts

			case <-ctx.Done():
				r.Log.Info("K8s context subscription stopped")
				return
			}
		}
	}()
	return contextsChan, nil
}

func (r *Resolver) getK8sContexts(ctx context.Context, provider models.Provider, selector model.PageFilter) (*model.K8sContextsPage, error) {
	tokenString := ctx.Value(models.TokenCtxKey).(string)
	resp, err := provider.GetK8sContexts(tokenString, selector.Page, selector.PageSize, *selector.Search, *selector.Order)
	if err != nil {
		return nil, err
	}
	var k8sContext model.K8sContextsPage
	err = json.Unmarshal(resp, &k8sContext)
	if err != nil {
		obj := "k8s context"
		return nil, handlers.ErrEncoding(err, obj)
	}
	return &k8sContext, nil
}
