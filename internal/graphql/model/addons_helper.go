package model

import (
	"github.com/layer5io/meshery/models"

	mesherykube "github.com/layer5io/meshkit/utils/kubernetes"
	meshsyncmodel "github.com/layer5io/meshsync/pkg/model"
)

func GetAddonsState(selectors []MeshType, kubeClient *mesherykube.Client, provider models.Provider, cid []string) ([]*AddonList, error) {
	addonlist := make([]*AddonList, 0)
	objects := make([]meshsyncmodel.Object, 0)

	for _, selector := range selectors {
		//subquery1 := r.DBHandler.Select("id").Where("kind = ? AND key = ? AND value = ?", meshsyncmodel.KindAnnotation, "meshery/component-type", "control-plane").Table("key_values")
		//subquery2 := r.DBHandler.Select("id").Where("id IN (?) AND kind = ? AND key = ? AND value IN (?)", subquery1, meshsyncmodel.KindAnnotation, "meshery/maintainer", selectors).Table("key_values")
		result := provider.GetGenericPersister().
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
			if meshsyncmodel.IsObject(obj) && len(addonPortSelector[obj.ObjectMeta.Name]) > 0 {
				// objstatus := corev1.ServiceStatus{}
				// err := utils.Unmarshal(obj.Status.Attribute, &objstatus)
				// if err != nil && len(obj.Status.Attribute) > 0 {
				// 	return nil, err
				// }

				// objspec := corev1.ServiceSpec{}
				// err = utils.Unmarshal(obj.Spec.Attribute, &objspec)
				// if err != nil && len(obj.Spec.Attribute) > 0 {
				// 	return nil, err
				// }

				// endpoint, err := mesherykube.GetEndpoint(context.TODO(),
				// 	&mesherykube.ServiceOptions{
				// 		APIServerURL: kubeClient.RestConfig.Host,
				// 		PortSelector: addonPortSelector[obj.ObjectMeta.Name],
				// 	},
				// 	&corev1.Service{
				// 		Spec:   objspec,
				// 		Status: objstatus,
				// 	})
				// if err != nil {
				// 	return nil, err
				// }

				// if endpoint.External == nil {
				// 	endpoint.External = endpoint.Internal
				// } else {
				// 	if !utils.TcpCheck(&utils.HostPort{
				// 		Address: endpoint.External.Address,
				// 		Port:    endpoint.External.Port,
				// 	}, nil) {
				// 		if !utils.TcpCheck(&utils.HostPort{
				// 			Address: "host.docker.internal",
				// 			Port:    endpoint.External.Port,
				// 		}, nil) {
				// 			u, _ := url.Parse(kubeClient.RestConfig.Host)
				// 			if utils.TcpCheck(&utils.HostPort{
				// 				Address: u.Hostname(),
				// 				Port:    endpoint.External.Port,
				// 			}, nil) {
				// 				u, _ := url.Parse(kubeClient.RestConfig.Host)
				// 				endpoint.External.Address = u.Hostname()
				// 			}
				// 		} else {
				// 			endpoint.External.Address = "host.docker.internal"
				// 		}
				// 	}
				// }

				addonlist = append(addonlist, &AddonList{
					Name:  obj.ObjectMeta.Name,
					Owner: selector.String(),
					// Endpoint: fmt.Sprintf("%s:%d", endpoint.External.Address, endpoint.External.Port),
				})
			}
		}
	}

	return addonlist, nil
}
