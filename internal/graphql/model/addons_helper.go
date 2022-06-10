package model

import (
	"github.com/layer5io/meshery/models"
	"gorm.io/gorm"

	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func GetAddonsState(selectors []MeshType, provider models.Provider, cid []string) ([]*AddonList, error) {
	addonlist := make([]*AddonList, 0)
	objects := make([]meshsyncmodel.Object, 0)
	cidMap := make(map[string]bool)
	var foundall bool
	for _, c := range cid {
		if c == "all" {
			foundall = true
			break
		}
		cidMap[c] = true
	}
	for _, selector := range selectors {
		//subquery1 := r.DBHandler.Select("id").Where("kind = ? AND key = ? AND value = ?", meshsyncmodel.KindAnnotation, "meshery/component-type", "control-plane").Table("key_values")
		//subquery2 := r.DBHandler.Select("id").Where("id IN (?) AND kind = ? AND key = ? AND value IN (?)", subquery1, meshsyncmodel.KindAnnotation, "meshery/maintainer", selectors).Table("key_values")
		var result *gorm.DB
		result = provider.GetGenericPersister().
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
			if !foundall && !cidMap[obj.ClusterID] {
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
