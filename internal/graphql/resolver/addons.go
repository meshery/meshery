package resolver

import (
	"context"
	"fmt"
	"strings"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/utils"
	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

var (
	portMapper = map[string]int32{
		"prometheus": 9090,
		"grafana":    3000,
	}
)

func (r *Resolver) changeAddonStatus(ctx context.Context) (model.Status, error) {
	return model.StatusProcessing, nil
}

func (r *Resolver) getAvailableAddons(ctx context.Context, selector *model.MeshType) ([]*model.AddonList, error) {
	addonlist := make([]*model.AddonList, 0)
	objects := make([]meshsyncmodel.Object, 0)

	selectors := make([]string, 0)
	if selector == nil || *selector == model.MeshTypeAll {
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

	subquery1 := r.DBHandler.Select("id").Where("key = ? AND value IN (?)", "meshery/maintainer", selectors).Table("key_values")
	subquery2 := r.DBHandler.Select("id").Where("id IN (?) AND key = ? AND value = ?", subquery1, "meshery/component-type", "control-plane").Table("key_values")
	result := r.DBHandler.
		Preload("ObjectMeta", "name IN ?", names).
		Preload("ObjectMeta.Labels").
		Preload("ObjectMeta.Annotations", "id IN (?)", subquery2).
		Preload("Spec").
		Preload("Status").
		Find(&objects, "kind = ?", "Service")
	if result.Error != nil {
		r.Log.Error(result.Error)
		return nil, result.Error
	}

	for _, obj := range objects {
		if meshsyncmodel.IsObject(obj) {
			objstatus := corev1.ServiceStatus{}
			err := utils.Unmarshal(obj.Status.Attribute, &objstatus)
			if err != nil && len(obj.Status.Attribute) > 0 {
				r.Log.Error(err)
				return nil, err
			}

			objspec := corev1.ServiceSpec{}
			err = utils.Unmarshal(obj.Spec.Attribute, &objspec)
			if err != nil && len(obj.Spec.Attribute) > 0 {
				r.Log.Error(err)
				return nil, err
			}

			port := portMapper[obj.ObjectMeta.Name]
			endpoint, err := mesherykube.GetEndpoint(context.TODO(),
				&mesherykube.ServiceOptions{
					APIServerURL: r.KubeClient.RestConfig.Host,
				},
				&corev1.Service{
					Spec:   objspec,
					Status: objstatus,
				})
			if err != nil {
				r.Log.Error(err)
				return nil, err
			}

			addonlist = append(addonlist, &model.AddonList{
				Type: strings.ToLower(selector.String()),
				Config: &model.AddonConfig{
					ServiceName: obj.ObjectMeta.Name,
					Endpoint:    fmt.Sprintf("%s:%d", endpoint.External.Address, port),
				},
			})
		}
	}

	return addonlist, nil
}

func (r *Resolver) listenToAddonState(ctx context.Context, selector *model.MeshType) (<-chan []*model.AddonList, error) {
	if r.addonChannel == nil {
		r.addonChannel = make(chan []*model.AddonList, 0)
	}

	go func() {
		r.Log.Info("Addons subscription started")
		err := r.connectToBroker(context.TODO())
		if err != nil {
			if err == ErrNoMeshSync {
				r.Log.Warn(err)
			} else {
				r.Log.Error(err)
				return
			}
		}

		select {
		case <-r.MeshSyncChannel:
			status, err := r.getAvailableAddons(ctx, selector)
			if err != nil {
				r.Log.Error(ErrAddonSubscription(err))
				return
			}
			r.addonChannel <- status
		}
	}()

	return r.addonChannel, nil
}
