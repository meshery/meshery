package resolver

import (
	"context"
	"strings"

	"github.com/layer5io/meshery/internal/graphql/model"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
	corev1 "k8s.io/api/core/v1"
)

func (r *Resolver) getControlPlanes(ctx context.Context, filter *model.ControlPlaneFilter) ([]*model.ControlPlane, error) {
	objects := make([]meshsyncmodel.Object, 0)
	controlplanelist := make([]*model.ControlPlane, 0)
	selectors := make([]model.MeshType, 0)
	if filter.Type == nil {
		for _, mesh := range model.AllMeshType {
			selectors = append(selectors, mesh)
		}
	} else {
		selectors = append(selectors, *filter.Type)
	}

	for _, selector := range selectors {
		//subquery := r.DBHandler.Select("id").Where("kind = ? AND key = ? AND value IN (?)", meshsyncmodel.KindAnnotation, "meshery/maintainer", selector).Table("key_values")
		result := r.DBHandler.
			Preload("ObjectMeta", "namespace = ?", controlPlaneNamespace[model.MeshType(selector)]).
			Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status").
			Find(&objects, "kind = ?", "Pod")
		if result.Error != nil {
			r.Log.Error(ErrQuery(result.Error))
			return nil, ErrQuery(result.Error)
		}
		members := make([]*model.ControlPlaneMember, 0)
		for _, obj := range objects {
			if meshsyncmodel.IsObject(obj) {
				objspec := corev1.PodSpec{}
				err := utils.Unmarshal(obj.Spec.Attribute, &objspec)
				if err != nil {
					r.Log.Error(err)
					return nil, err
				}

				version := "unknown"
				if len(strings.Split(objspec.Containers[0].Image, ":")) > 0 {
					version = strings.Split(objspec.Containers[0].Image, ":")[1]
				}
				members = append(members, &model.ControlPlaneMember{
					Name:      obj.ObjectMeta.Name,
					Component: strings.Split(obj.ObjectMeta.GenerateName, "-")[0],
					Version:   version,
					Namespace: obj.ObjectMeta.Namespace,
				})
			}
		}
		controlplanelist = append(controlplanelist, &model.ControlPlane{
			Name:    strings.ToLower(selector.String()),
			Members: members,
		})
	}
	return controlplanelist, nil
}
func (r *Resolver) listenToControlPlaneState(ctx context.Context, filter *model.ControlPlaneFilter) (<-chan []*model.ControlPlane, error) {
	if r.controlPlaneChannel == nil {
		r.controlPlaneChannel = make(chan []*model.ControlPlane)
	}
	go func() {
		r.Log.Info("ControlPlane subscription started")
		err := r.connectToBroker(context.TODO())
		if err != nil && err != ErrNoMeshSync {
			r.Log.Error(err)
			return
		}
		select {
		case <-r.MeshSyncChannel:
			status, err := r.getControlPlanes(ctx, filter)
			if err != nil {
				r.Log.Error(ErrControlPlaneSubscription(err))
				return
			}
			r.controlPlaneChannel <- status
		}
	}()
	return r.controlPlaneChannel, nil
}
