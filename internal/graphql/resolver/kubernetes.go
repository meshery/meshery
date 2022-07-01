package resolver

import (
	"context"
	"database/sql"
	"strings"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshery/models/pattern/core"
	meshkitKube "github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/kubernetes/describe"
)

func (r *Resolver) getAvailableNamespaces(ctx context.Context, provider models.Provider, k8sClusterIDs []string) ([]*model.NameSpace, error) {
	var cids []string
	if len(k8sClusterIDs) != 0 {
		cids = k8sClusterIDs
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
	namespaces := make([]string, 0)
	var rows *sql.Rows
	var err error
	if len(cids) == 1 && cids[0] == "all" {
		rows, err = provider.GetGenericPersister().Raw("SELECT DISTINCT rom.name as name FROM objects o LEFT JOIN resource_object_meta rom ON o.id = rom.id WHERE o.kind = 'Namespace'").Rows()
	} else {
		rows, err = provider.GetGenericPersister().Raw("SELECT DISTINCT rom.name as name FROM objects o LEFT JOIN resource_object_meta rom ON o.id = rom.id WHERE o.kind = 'Namespace' AND o.cluster_id IN ?", cids).Rows()
	}

	if err != nil {
		r.Log.Error(ErrGettingNamespace(err))
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			r.Log.Error(ErrGettingNamespace(err))
			return nil, err
		}

		namespaces = append(namespaces, name)
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
