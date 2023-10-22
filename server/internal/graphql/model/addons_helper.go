package model

import (
	"context"

	"github.com/layer5io/meshery/server/models"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func GetAddonsState(ctx context.Context, selectors []MeshType, provider models.Provider, cid []string) ([]*AddonList, error) {
	addonlist := make([]*AddonList, 0)
	objects := make([]meshsyncmodel.KubernetesResource, 0)
	cidMap := make(map[string]bool)

	for _, c := range cid {
		cidMap[c] = true
	}

	for _, selector := range selectors {
		result := provider.GetGenericPersister().
			Where("cluster_id IN ?", cid).
			Preload("KubernetesResourceMeta", "namespace = ?", controlPlaneNamespace[selector]).
			Preload("KubernetesResourceMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("KubernetesResourceMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
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
			if meshsyncmodel.IsObject(obj) && len(addonPortSelector[obj.KubernetesResourceMeta.Name]) > 0 {
				addonlist = append(addonlist, &AddonList{
					Name:  obj.KubernetesResourceMeta.Name,
					Owner: selector.String(),
				})
			}
		}
	}

	return addonlist, nil
}
