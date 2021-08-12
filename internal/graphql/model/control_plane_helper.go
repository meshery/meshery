package model

import (
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
			Preload("ObjectMeta", "namespace = ?", controlPlaneNamespace[MeshType(selector)]).
			Preload("ObjectMeta.Labels", "kind = ?", meshsyncmodel.KindLabel).
			Preload("ObjectMeta.Annotations", "kind = ?", meshsyncmodel.KindAnnotation).
			Preload("Spec").
			Preload("Status").
			Find(&object, "kind = ?", "Pod")
		if result.Error != nil {
			return nil, ErrQuery(result.Error)
		}

		members := make([]*ControlPlaneMember, 0)
		for _, obj := range object {
			if meshsyncmodel.IsObject(obj) {
				objspec := corev1.PodSpec{}
				objstatus := corev1.PodStatus{}
				err := utils.Unmarshal(obj.Spec.Attribute, &objspec)
				if err != nil {
					return nil, err
				}
				err = utils.Unmarshal(obj.Status.Attribute, &objstatus)
				if err != nil {
					return nil, err
				}
				proxyContainers := make([]*Container, 0)

				containers := objspec.Containers
				statuses := objstatus.ContainerStatuses
				if len(containers) == len(statuses) {
					for i := range containers {
						container := containers[i]
						status := statuses[i]
						var proxyStatus *ContainerStatus
						// Statuses
						if strings.Contains(status.Name, "proxy") || strings.Contains(status.Image, "proxy") {
							proxyStatus = &ContainerStatus{
								Name:  status.Name,
								State: status.State,
								// State: &ContainerStatusState{
								// 	Waiting: &ContainerStatusStateWaiting{
								// 		Reason:  &status.State.Waiting.Reason,
								// 		Message: &status.State.Waiting.Message,
								// 	},
								// 	Running: &ContainerStatusStateRunning{
								// 		StartedAt: &status.State.Running.StartedAt.Time,
								// 	},
								// 	Terminated: &ContainerStatusStateTerminated{
								// 		Reason:  &status.State.Terminated.Reason,
								// 		Message: &status.State.Terminated.Message,
								// 		// ExitCode:    &exitCode,
								// 		// Signal:      &signal,
								// 		StartedAt:   &status.State.Terminated.StartedAt.Time,
								// 		FinishedAt:  &status.State.Terminated.FinishedAt.Time,
								// 		ContainerID: &status.State.Terminated.ContainerID,
								// 	},
								// },
								Started: *status.Started,
								Ready:   status.Ready,
								// RestartCount: reflect.ValueOf(status.RestartCount).Addr().Int(),
							}
						}
						// Container
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
								Name:   container.Name,
								Image:  container.Image,
								Ports:  proxyPorts,
								Status: proxyStatus,
							})
						}
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
