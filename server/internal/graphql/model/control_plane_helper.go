package model

import (
	"context"
	"strings"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/encoding"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	corev1 "k8s.io/api/core/v1"
)

func GetControlPlaneState(ctx context.Context, selectors []MeshType, provider models.Provider, cid []string) ([]*ControlPlane, error) {
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
func haveCommonElements(a []string, b map[string]bool) bool {
	for _, ae := range a {
		if b[ae] {
			return true
		}
	}
	return false
}
