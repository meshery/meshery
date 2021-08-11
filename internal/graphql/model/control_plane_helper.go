package model

import (
	// "encoding/json"
	// "fmt"
	"reflect"
	"strings"

	"github.com/layer5io/meshery/models"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	corev1 "k8s.io/api/core/v1"
)

func GetControlPlaneState(selectors []MeshType, provider models.Provider) ([]*ControlPlane, error) {
	object := []meshsyncmodel.Object{}
	controlplanelist := make([]*ControlPlane, 0)

	for _, selector := range selectors {
		result := provider.GetGenericPersister().Model(&meshsyncmodel.Object{}).
			Preload("ObjectMeta", "namespace = ?", controlPlaneNamespace[MeshType(selector)]). // milih istio/kuma atau service mesh lainnya
			Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status").
			Find(&object, "kind = ?", "Pod")
		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}

		// result2 := provider.GetGenericPersister().Model(&meshsyncmodel.Object{}).
		// Preload("ObjectMeta", "namespace = ?", controlPlaneNamespace[MeshType(selector)]). // milih istio/kuma atau service mesh lainnya
		// Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
		// Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
		// Preload("Spec").
		// Preload("Status").
		// Find(&object, "kind = ?", "Pod")

		members := make([]*ControlPlaneMember, 0)
		for _, obj := range object {
			if meshsyncmodel.IsObject(obj) {
				objspec := corev1.PodSpec{}
				err := utils.Unmarshal(obj.Spec.Attribute, &objspec)
				if err != nil {
					return nil, err
				}

				// fmt.Println("OBJSPEC!")
				// fmt.Println("..")
				// fmt.Println("..")
				// fmt.Println("..")

				proxyContainers := make([]*Container, 0)

				for _, container := range objspec.Containers {
					// fmt.Println("CONTAINER!")
					// s, _ := json.MarshalIndent(container, "", "\t")
					// fmt.Println(string(s))
					// fmt.Println("..")
					// fmt.Println("..")
					if strings.Contains(container.Name, "proxy") || strings.Contains(container.Image, "proxy") {

						proxyPorts := make([]*ContainerPort, 0)
						for _, port := range container.Ports {
							proxyPorts = append(proxyPorts, &ContainerPort{
								Name:          &port.Name,
								ContainerPort: int(port.ContainerPort),
								Protocol:      reflect.ValueOf(port.Protocol).String(),
							})
						}

						proxyContainers = append(proxyContainers, &Container{
							Name:  container.Name,
							Image: container.Image,
							Ports: proxyPorts,
						})
					}
				}

				version := "unknown"
				if len(strings.Split(objspec.Containers[0].Image, ":")) > 0 {
					version = strings.Split(objspec.Containers[0].Image, ":")[1]
				}
				members = append(members, &ControlPlaneMember{
					Name:       obj.ObjectMeta.Name,
					Component:  strings.Split(obj.ObjectMeta.GenerateName, "-")[0],
					Version:    version,
					Namespace:  obj.ObjectMeta.Namespace,
					DataPlanes: proxyContainers,
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
