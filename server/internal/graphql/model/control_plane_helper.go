package model

import (
	"context"
	"strings"

	"github.com/meshery/meshery/server/models"
	"github.com/meshery/meshkit/encoding"

	meshsyncmodel "github.com/meshery/meshsync/pkg/model"

	corev1 "k8s.io/api/core/v1"
)

func GetControlPlaneState(ctx context.Context, selectors []MeshType, provider models.Provider, cid []string) ([]*ControlPlane, error) {
	if len(cid) > 0 {
		batchMap, err := BatchGetControlPlaneState(ctx, selectors, provider, cid)
		if err != nil {
			return nil, err
		}
		controlplanelist := make([]*ControlPlane, 0)
		for _, selector := range selectors {
			selectorName := strings.ToLower(selector.String())
			var members []*ControlPlaneMember
			for _, c := range cid {
				for _, cp := range batchMap[c] {
					if cp.Name == selectorName {
						members = append(members, cp.Members...)
					}
				}
			}
			controlplanelist = append(controlplanelist, &ControlPlane{
				Name:    selectorName,
				Members: members,
			})
		}
		return controlplanelist, nil
	}

	object := []meshsyncmodel.KubernetesResource{}
	controlplanelist := make([]*ControlPlane, 0)
	cidMap := make(map[string]bool)

	for _, c := range cid {
		cidMap[c] = true
	}

	for _, selector := range selectors {
		result := provider.GetGenericPersister().Model(&meshsyncmodel.KubernetesResource{}).
			Preload("KubernetesResourceMeta", "namespace IN ?", controlPlaneNamespace[MeshType(selector)]).
			Preload("KubernetesResourceMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("KubernetesResourceMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status").
			Find(&object, "kind = ?", "Pod")
		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}
		members := make([]*ControlPlaneMember, 0)
		for _, obj := range object {
			if !cidMap[obj.ClusterID] {
				continue
			}
			if meshsyncmodel.IsObject(obj) { //As a fallback extract objectmeta manually, if possible
				objspec := corev1.PodSpec{}
				err := encoding.Unmarshal([]byte(obj.Spec.Attribute), &objspec)
				if err != nil {
					return nil, err
				}
				var imageOrgs = make(map[string]bool)
				for _, c := range objspec.Containers {
					if len(strings.Split(c.Image, "/")) > 1 {
						imageOrgs[strings.Split(c.Image, "/")[1]] = true // Extracting image org from <domainname>/<imageorg>/<imagename>
					}
				}
				version := "unknown"
				//If image orgs are not passed on in from controlPlaneImageOrgs variable, then skip this filtering (for backward compatibility)
				if len(controlPlaneImageOrgs[MeshType(selector)]) != 0 && !haveCommonElements(controlPlaneImageOrgs[MeshType(selector)], imageOrgs) {
					continue
				}

				if len(strings.Split(objspec.Containers[0].Image, ":")) > 1 {
					version = strings.Split(objspec.Containers[0].Image, ":")[1]
				}

				members = append(members, &ControlPlaneMember{
					Name:      obj.KubernetesResourceMeta.Name,
					Component: strings.Split(obj.KubernetesResourceMeta.GenerateName, "-")[0],
					Version:   strings.Split(version, "@")[0],
					Namespace: obj.KubernetesResourceMeta.Namespace,
				})
			}
		}
		controlplanelist = append(controlplanelist, &ControlPlane{
			Name:    strings.ToLower(selector.String()),
			Members: members,
		})
	}

	return controlplanelist, nil
}

// BatchGetControlPlaneState fetches control plane data for multiple clusters in a single batched query
// grouping results by clusterID.
func BatchGetControlPlaneState(ctx context.Context, selectors []MeshType, provider models.Provider, clusterIDs []string) (map[string][]*ControlPlane, error) {
	resultMap := make(map[string][]*ControlPlane)
	for _, cid := range clusterIDs {
		resultMap[cid] = make([]*ControlPlane, 0)
	}

	if provider == nil || provider.GetGenericPersister() == nil {
		return resultMap, nil
	}

	for _, selector := range selectors {
		object := []meshsyncmodel.KubernetesResource{}
		query := provider.GetGenericPersister().Model(&meshsyncmodel.KubernetesResource{}).
			Preload("KubernetesResourceMeta", "namespace IN ?", controlPlaneNamespace[MeshType(selector)]).
			Preload("KubernetesResourceMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("KubernetesResourceMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status")

		if len(clusterIDs) > 0 {
			query = query.Where("cluster_id IN (?)", clusterIDs)
		}

		result := query.Find(&object, "kind = ?", "Pod")
		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}

		membersByCluster := make(map[string][]*ControlPlaneMember)
		for _, cid := range clusterIDs {
			membersByCluster[cid] = make([]*ControlPlaneMember, 0)
		}

		for _, obj := range object {
			if meshsyncmodel.IsObject(obj) {
				objspec := corev1.PodSpec{}
				err := encoding.Unmarshal([]byte(obj.Spec.Attribute), &objspec)
				if err != nil {
					return nil, err
				}
				var imageOrgs = make(map[string]bool)
				for _, c := range objspec.Containers {
					if len(strings.Split(c.Image, "/")) > 1 {
						imageOrgs[strings.Split(c.Image, "/")[1]] = true
					}
				}
				version := "unknown"
				if len(controlPlaneImageOrgs[MeshType(selector)]) != 0 && !haveCommonElements(controlPlaneImageOrgs[MeshType(selector)], imageOrgs) {
					continue
				}

				if len(objspec.Containers) > 0 && len(strings.Split(objspec.Containers[0].Image, ":")) > 1 {
					version = strings.Split(objspec.Containers[0].Image, ":")[1]
				}

				membersByCluster[obj.ClusterID] = append(membersByCluster[obj.ClusterID], &ControlPlaneMember{
					Name:      obj.KubernetesResourceMeta.Name,
					Component: strings.Split(obj.KubernetesResourceMeta.GenerateName, "-")[0],
					Version:   strings.Split(version, "@")[0],
					Namespace: obj.KubernetesResourceMeta.Namespace,
				})
			}
		}

		for _, cid := range clusterIDs {
			resultMap[cid] = append(resultMap[cid], &ControlPlane{
				Name:    strings.ToLower(selector.String()),
				Members: membersByCluster[cid],
			})
		}
	}

	return resultMap, nil
}
func haveCommonElements(a []string, b map[string]bool) bool {
	for _, ae := range a {
		if b[ae] {
			return true
		}
	}
	return false
}
