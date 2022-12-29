package resolver

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/layer5io/meshkit/models/meshmodel"
)

func (r *Resolver) subscribeMeshModelSummary(ctx context.Context, provider models.Provider, selector model.MeshModelSummarySelector) (<-chan *model.MeshModelSummary, error) {
	ch := make(chan struct{}, 1)
	respChan := make(chan *model.MeshModelSummary)

	go func() {
		r.Log.Info("Initializing MeshModelSummary subscription")
		for {
			select {
			case <-ch: 
			  meshModelSummary, err := r.getMeshModelSummary(ctx, provider, selector)
				if err != nil {
					logrus.Error("Error while fetching MeshModelSummary: ", err)
					// r.Log.Error(err)
					break
				}
				respChan <- meshModelSummary
			case <-ctx.Done():
				r.Log.Info("Closing MeshModelSummary subscription")
				return
			}
		}
	}()

	return respChan, nil
}


func (r *Resolver) getMeshModelSummary(ctx context.Context, provider models.Provider, selector model.MeshModelSummarySelector) (*model.MeshModelSummary, error) {
	regManager, ok := ctx.Value(models.RegistryManagerKey).(*meshmodel.RegistryManager)
	summary := &model.MeshModelSummary{}
	if !ok {
		logrus.Error("Error while fetching RegistryManager")
		err1 := errors.New("Error while fetching RegistryManager")
		return nil, err1
	}
	switch selector.Type {
	case "Components" :
		components := getMeshModelComponents(regManager)
		logrus.Debug("components: ", components)
		summary.Components = components
		summary.Relationships = []*model.MeshModelRelationship{}
		break
	// case "relationships" :
	// 	relationships := getMeshModelRelationships(regManager)
	// 	summary.Relationships = relationships
	// 	summary.Components = []*model.MeshModelComponent{}
	// 	break
	}
	return summary, nil
}

// func getMeshModelRelationships(regManager *meshmodel.RegistryManager) []*model.MeshModelRelationship{
// 	res := regManager.GetEntities(&v1alpha1.RelationshipFilter{})
// 	relationships := make([]*model.MeshModelRelationship, 0)
// 	var relmap := make(map[string]*relationships)
// 	for _, r := range res {
// 		def, _ := r.(v1alpha1.RelationshipDefinition)
// 		if relmap[def.Model.Name] == nil {

// 		}
// 	}
// }


type typesResponseWithModelname struct {
	DisplayName string   `json:"display-name"`
	Versions    []string `json:"versions"`
}

func getMeshModelComponents(regManager *meshmodel.RegistryManager) []*model.MeshModelComponent {
	res := regManager.GetEntities(&v1alpha1.ComponentFilter{})
	logrus.Debug("res: ", res)
	components := make([]*model.MeshModelComponent, 0)
	var response = make(map[string]*typesResponseWithModelname)
	for _, r := range res {
		def, _ := r.(v1alpha1.ComponentDefinition)
		if response[def.Model.Name] == nil {
			response[def.Model.Name] = &typesResponseWithModelname{
				DisplayName: def.Model.DisplayName,
				Versions:    []string{def.Model.Version},
			}
		} else {
			response[def.Model.Name].Versions = append(response[def.Model.Name].Versions, def.Model.Version)
		}
	}
	for _, x := range response {
		x.Versions = filterUniqueElementsArray(x.Versions)
	}
	for _, x := range response {
		components = append(components, &model.MeshModelComponent{
			Name: x.DisplayName,
			Count:  len(x.Versions),
		})
	}
	return components
}

func filterUniqueElementsArray(s []string) []string {
	m := make(map[string]bool)
	for _, ele := range s {
		m[ele] = true
	}
	ans := make([]string, 0)
	for a := range m {
		ans = append(ans, a)
	}
	return ans
}
