package model

import (
	"context"

	"github.com/layer5io/meshery/server/models"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func GetAddonsState(ctx context.Context, selectors []MeshType, provider models.Provider, cid []string) ([]*AddonList, error) {
	addonlist := make([]*AddonList, 0)
	objects := make([]meshsyncmodel.Object, 0)
	cidMap := make(map[string]bool)
	if len(cid) == 1 && cid[0] == "all" {
		k8sctxs, ok := ctx.Value(models.AllKubeClusterKey).([]models.K8sContext)
		if !ok || len(k8sctxs) == 0 {
			return nil, ErrMesheryClientNil(nil)
		}
		for _, k8ctx := range k8sctxs {
			if k8ctx.KubernetesServerID != nil {
				cidMap[k8ctx.KubernetesServerID.String()] = true
			}
		}
	} else {
		for _, c := range cid {
			cidMap[c] = true
		}
	}
	for _, selector := range selectors {
		//subquery1 := r.DBHandler.Select("id").Where("kind = ? AND key = ? AND value = ?", meshsyncmodel.KindAnnotation, "meshery/component-type", "control-plane").Table("key_values")
		//subquery2 := r.DBHandler.Select("id").Where("id IN (?) AND kind = ? AND key = ? AND value IN (?)", subquery1, meshsyncmodel.KindAnnotation, "meshery/maintainer", selectors).Table("key_values")
		result := provider.GetGenericPersister().
			Where("cluster_id IN ?", cid).
			Preload("ObjectMeta", "namespace = ?", controlPlaneNamespace[selector]).
			Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status").
			Find(&objects, "kind = ?", "Service")

		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}

		for _, obj := range objects {
			if !cidMap[obj.ClusterID] {
				continue
			}
			if meshsyncmodel.IsObject(obj) && len(addonPortSelector[obj.ObjectMeta.Name]) > 0 {
				addonlist = append(addonlist, &AddonList{
					Name:  obj.ObjectMeta.Name,
					Owner: selector.String(),
				})
			}
		}
	}

	return addonlist, nil
}
