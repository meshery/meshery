package istio

import (
	"errors"

	"github.com/layer5io/meshery/graphql/model"
	"github.com/layer5io/meshkit/database"
	syncModel "github.com/layer5io/meshsync/pkg/model"
)

// IstioHandler will implement the resolvers for istio
type IstioHandler struct {
	DB *database.Handler
}

// MapEdges to get edges
func (i *IstioHandler) MapEdges(view *model.View) ([]*model.MapEdge, error) {
	switch *view {
	case model.ViewServicemeshview:
		return i.serviceMeshViewEdges()
	case model.ViewWorkloadview:
		return i.workloadViewEdges()
	}
	return nil, errors.New("View not defined")
}

// Type of edges to get
// 1. IngressGateway to LogicalGateway
// 2. LogicalGateway to VirtualService
// 3. VirtualService to node(4.)
// ... maybe others after deciding on the 4th point above
func (i *IstioHandler) serviceMeshViewEdges() ([]*model.MapEdge, error) {
	return nil, nil
}

// Type of edges to get
// 1. Service Endpoint to Destination Rule
// 2. DestinationRule to Pod
func (i *IstioHandler) workloadViewEdges() ([]*model.MapEdge, error) {
	return nil, nil
}

// MapNodes to get nodes
func (i *IstioHandler) MapNodes(view *model.View) ([]*model.MapNode, error) {
	switch *view {
	case model.ViewServicemeshview:
		return i.serviceMeshViewNodes()
	case model.ViewWorkloadview:
		return i.workloadViewNodes()
	}
	return nil, errors.New("View not defined")
}

// Type of nodes to get
// 1. Ingress Gateway
// 2. Logical Gateway
// 3. Virtual Service
// 4. Whole application as one node || divided by service endpoint
func (i *IstioHandler) serviceMeshViewNodes() ([]*model.MapNode, error) {
	var objs []syncModel.Object
	result := i.DB.Find(&objs)
	if result.Error != nil {
		return nil, result.Error
	}

	var typemetas []syncModel.ResourceTypeMeta
	result = i.DB.Find(&typemetas)
	if result.Error != nil {
		return nil, result.Error
	}

	var objectmetas []syncModel.ResourceObjectMeta
	result = i.DB.Find(&objectmetas)
	if result.Error != nil {
		return nil, result.Error
	}

	for i, obj := range objs {
		obj.TypeMeta = typemetas[i]
		obj.ObjectMeta = objectmetas[i]
		objs[i] = obj
	}

	var objects []*model.MapNode
	for _, obj := range objs {
		if obj.TypeMeta.Kind == "gateway" || obj.TypeMeta.Kind == "virtual-service" {
			objects = append(objects, &model.MapNode{
				ID: &obj.Index.ResourceID,
				Metadata: &model.MetaData{
					Name:      &obj.ObjectMeta.Name,
					Namespace: &obj.ObjectMeta.Name,
					NodeType:  &obj.TypeMeta.Kind,
				},
			})
		}
	}

	return objects, nil
}

// Type of nodes to get
// 1. Service Enpoint
// 2. Destination Rule
// 3. Pods
func (i *IstioHandler) workloadViewNodes() ([]*model.MapNode, error) {
	var objs []syncModel.Object
	result := i.DB.Find(&objs)
	if result.Error != nil {
		return nil, result.Error
	}

	var typemetas []syncModel.ResourceTypeMeta
	result = i.DB.Find(&typemetas)
	if result.Error != nil {
		return nil, result.Error
	}

	var objectmetas []syncModel.ResourceObjectMeta
	result = i.DB.Find(&objectmetas)
	if result.Error != nil {
		return nil, result.Error
	}

	var objects []*model.MapNode
	for _, obj := range objs {
		if obj.TypeMeta.Kind == "service" || obj.TypeMeta.Kind == "destination-rule" || obj.TypeMeta.Kind == "pod" {
			objects = append(objects, &model.MapNode{
				ID: &obj.Index.ResourceID,
				Metadata: &model.MetaData{
					Name:      &obj.ObjectMeta.Name,
					Namespace: &obj.ObjectMeta.Name,
					NodeType:  &obj.TypeMeta.Kind,
				},
			})
		}
	}
	return objects, nil
}
