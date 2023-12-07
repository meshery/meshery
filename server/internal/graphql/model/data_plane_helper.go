package model

import (
	"context"
	"reflect"
	"strings"

	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/utils"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"

	corev1 "k8s.io/api/core/v1"
)

func GetDataPlaneState(ctx context.Context, selectors []MeshType, provider models.Provider, cid []string) ([]*DataPlane, error) {
	object := []meshsyncmodel.KubernetesResource{}
	dataPlaneList := make([]*DataPlane, 0)
	cidMap := make(map[string]bool)

	for _, c := range cid {
		cidMap[c] = true
	}

	for _, selector := range selectors {
		result := provider.GetGenericPersister().Model(&meshsyncmodel.KubernetesResource{}).
			Preload("KubernetesResourceMeta", "namespace IN ?", controlPlaneNamespace[MeshType(selector)]).
			Preload("Status").
			Preload("Spec"). // get only resources specs that has proxy string inside its attributes
			Where("EXISTS(SELECT 1 FROM kubernetes_resource_specs rsp WHERE rsp.attribute LIKE ? AND rsp.id = kubernetes_resources.id)", `%proxy%`).
			// get only resources statuses that has proxy string inside its attributes
			Where("EXISTS(SELECT 1 FROM kubernetes_resource_statuses rst WHERE rst.attribute LIKE ? AND rst.id = kubernetes_resources.id)", `%proxy%`).
			Find(&object, "kind = ?", "Pod")

		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}
		proxies := make([]*Container, 0)
		for _, obj := range object {
			if !cidMap[obj.ClusterID] {
				continue
			}
			if meshsyncmodel.IsObject(obj) {
				objspec := corev1.PodSpec{}
				objstatus := corev1.PodStatus{}

				// unmarshal resource_specs
				err := utils.Unmarshal(obj.Spec.Attribute, &objspec)
				if err != nil {
					return nil, err
				}

				// unmarshal resource_statuses
				err = utils.Unmarshal(obj.Status.Attribute, &objstatus)
				if err != nil {
					return nil, err
				}

				// make initial list of proxy containers
				containers := objspec.Containers
				statuses := objstatus.ContainerStatuses
				// check if the length of containers and its statuses is the same
				if len(containers) == len(statuses) {
					for i := range statuses {
						var container corev1.Container
						status := statuses[i]
						// filter based on Name
						// since apparently the GORM query gets all the child images on the pods
						// and not all child images are proxies
						if strings.Contains(status.Name, "proxy") || strings.Contains(status.Image, "proxy") {
							// Name of the container specified as a DNS_LABEL.
							// Each container in a pod must have a unique name (DNS_LABEL).
							for i := range containers {
								if containers[i].Name == status.Name {
									// Found!
									container = containers[i]
									break
								}
							}
							proxyStatus := &ContainerStatus{
								ContainerStatusName: status.Name,
								State:               status.State,
								LastState:           status.LastTerminationState,
								Started:             *status.Started,
								Ready:               status.Ready,
								RestartCount:        status.RestartCount,
								ImageID:             status.ImageID,
								Image:               status.Image,
								ContainerID:         status.ContainerID,
							}
							proxyPorts := make([]*ContainerPort, 0)
							for _, port := range container.Ports {
								proxyPorts = append(proxyPorts, &ContainerPort{
									Name:          &port.Name,
									ContainerPort: int(port.ContainerPort),
									Protocol:      reflect.ValueOf(port.Protocol).String(),
								})
							}
							proxies = append(proxies, &Container{
								ContainerName:          container.Name,
								Image:                  container.Image,
								Ports:                  proxyPorts,
								Status:                 proxyStatus,
								ControlPlaneMemberName: obj.KubernetesResourceMeta.Name,
								Resources:              container.Resources,
							})
						}
					}
				}
			}
			dataPlaneList = append(dataPlaneList, &DataPlane{
				Name:    strings.ToLower(selector.String()),
				Proxies: proxies,
			})
		}
	}
	return dataPlaneList, nil
}
