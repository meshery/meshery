package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

func (r *Resolver) changeAddonStatus(ctx context.Context) (*model.Status, error) {
	return nil, nil
}

func (r *Resolver) getAvailableAddons(ctx context.Context, selector *model.MeshType) ([]*model.AddonList, error) {
	addonlist := make([]*model.AddonList, 0)
	objects := make([]meshsyncmodel.Object, 0)
	status := model.StatusEnabled

	selectors := make([]string, 0)
	if selector == nil {
		for _, obj := range model.AllAddonSelector {
			selectors = append(selectors, strings.ToLower(obj.String()))
		}
	} else {
		selectors = append(selectors, strings.ToLower(selector.String()))
	}

	names := make([]string, 0)
	for _, addon := range model.AllAddonSelector {
		names = append(names, strings.ToLower(addon.String()))
	}

	subquery1 := r.DBHandler.Select("id").Where("key = ? AND value IN ?", "meshery/maintainer", selectors).Table("key_values")
	subquery2 := r.DBHandler.Select("id").Where("id IN (?) AND key = ? AND value = ?", subquery1, "meshery/component-type", "control-plane").Table("key_values")
	result := r.DBHandler.
		Preload("TypeMeta", "kind = ?", "Service").
		Preload("ObjectMeta", "name IN ?", names).
		Preload("ObjectMeta.Labels").
		Preload("ObjectMeta.Annotations").
		Preload("Spec").
		Preload("Status").
		Find(&objects, "id IN (?)", subquery2)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, obj := range objects {
		if obj.TypeMeta != nil && obj.ObjectMeta != nil && obj.Spec != nil && obj.Status != nil {

			objstatus := corev1.ServiceStatus{}
			endpoint := ""
			err := utils.Unmarshal(obj.Status.Attribute, &objstatus)
			if err != nil {
				return nil, err
			}

			if &objstatus.LoadBalancer != nil && len(objstatus.LoadBalancer.Ingress) > 0 {
				if objstatus.LoadBalancer.Ingress[0].IP != "" {
					endpoint = objstatus.LoadBalancer.Ingress[0].IP
				} else {
					endpoint = objstatus.LoadBalancer.Ingress[0].Hostname
				}
			}

			addonlist = append(addonlist, &model.AddonList{
				Type:   strings.ToLower(selector.String()),
				Status: &status,
				Config: &model.AddonConfig{
					ServiceName: obj.ObjectMeta.Name,
					Endpoint:    endpoint,
				},
			})
			fmt.Println(obj.Status.Attribute)
		}
	}

	return addonlist, nil
}

func (r *Resolver) listenToAddonEvents(ctx context.Context) (<-chan []*model.AddonList, error) {
	r.addonChannel = make(chan []*model.AddonList, 0)

	go func() {
		select {
		case <-r.meshsyncChannel:
			status, err := r.getAvailableAddons(ctx, nil)
			if err != nil {
				return
			}
			r.addonChannel <- status
		}
	}()

	return r.addonChannel, nil
}
