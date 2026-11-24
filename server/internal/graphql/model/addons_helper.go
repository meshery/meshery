package model

import (
	"context"

	"github.com/meshery/meshery/server/models"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"
)

// GetAddonsState fetches addons across the specified clusters and mesh types.
func GetAddonsState(ctx context.Context, selectors []MeshType, provider models.GenericPersisterGetter, cid []string) ([]*AddonList, error) {
	if len(cid) > 0 {
		seen := make(map[string]struct{}, len(cid))
		uniq := make([]string, 0, len(cid))
		for _, c := range cid {
			if _, ok := seen[c]; !ok {
				seen[c] = struct{}{}
				uniq = append(uniq, c)
			}
		}
		cid = uniq

		batchMap, err := BatchGetAddonsState(ctx, selectors, provider, cid)
		if err != nil {
			return nil, err
		}
		addonlist := make([]*AddonList, 0)
		for _, selector := range selectors {
			selectorName := selector.String()
			for _, c := range cid {
				for _, addon := range batchMap[c] {
					if addon.Owner == selectorName {
						addonlist = append(addonlist, addon)
					}
				}
			}
		}
		return addonlist, nil
	}

	addonlist := make([]*AddonList, 0)
	objects := make([]meshsyncmodel.KubernetesResource, 0)
	cidMap := make(map[string]bool)

	for _, c := range cid {
		cidMap[c] = true
	}

	if provider == nil || provider.GetGenericPersister() == nil {
		return addonlist, nil
	}

	for _, selector := range selectors {
		result := provider.GetGenericPersister().
			WithContext(ctx).
			Where("cluster_id IN ?", cid).
			Preload("KubernetesResourceMeta", "namespace IN ?", controlPlaneNamespace[selector]).
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

// BatchGetAddonsState fetches addon data for multiple clusters in a single batched query
// grouping results by clusterID.
func BatchGetAddonsState(ctx context.Context, selectors []MeshType, provider models.GenericPersisterGetter, clusterIDs []string) (map[string][]*AddonList, error) {
	if len(clusterIDs) > 0 {
		seen := make(map[string]struct{}, len(clusterIDs))
		uniq := make([]string, 0, len(clusterIDs))
		for _, c := range clusterIDs {
			if _, ok := seen[c]; !ok {
				seen[c] = struct{}{}
				uniq = append(uniq, c)
			}
		}
		clusterIDs = uniq
	}

	resultMap := make(map[string][]*AddonList)
	for _, cid := range clusterIDs {
		resultMap[cid] = make([]*AddonList, 0)
	}

	if provider == nil || provider.GetGenericPersister() == nil {
		return resultMap, nil
	}

	cidMap := make(map[string]bool, len(clusterIDs))
	for _, cid := range clusterIDs {
		cidMap[cid] = true
	}

	for _, selector := range selectors {
		objects := make([]meshsyncmodel.KubernetesResource, 0)
		query := provider.GetGenericPersister().Model(&meshsyncmodel.KubernetesResource{}).
			WithContext(ctx).
			Preload("KubernetesResourceMeta", "namespace IN ?", controlPlaneNamespace[selector]).
			Preload("KubernetesResourceMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("KubernetesResourceMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status")

		if len(clusterIDs) > 0 {
			query = query.Where("cluster_id IN (?)", clusterIDs)
		}

		result := query.Find(&objects, "kind = ?", "Service")
		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}

		for _, obj := range objects {
			if len(clusterIDs) > 0 && !cidMap[obj.ClusterID] {
				continue
			}
			if meshsyncmodel.IsObject(obj) && len(addonPortSelector[obj.KubernetesResourceMeta.Name]) > 0 {
				addon := &AddonList{
					Name:  obj.KubernetesResourceMeta.Name,
					Owner: selector.String(),
				}
				resultMap[obj.ClusterID] = append(resultMap[obj.ClusterID], addon)
			}
		}
	}

	return resultMap, nil
}
