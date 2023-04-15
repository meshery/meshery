package resolver

import (
	"context"
	"errors"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	"github.com/layer5io/meshkit/models/meshmodel"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (r *Resolver) subscribeMeshModelSummary(ctx context.Context, provider models.Provider, selector model.MeshModelSummarySelector) (<-chan *model.MeshModelSummary, error) {
	ch := make(chan struct{}, 1)
	ch <- struct{}{}
	respChan := make(chan *model.MeshModelSummary)

	r.Config.MeshModelSummaryChannel.Subscribe(ch)
	go func() {
		r.Log.Info("Initializing MeshModelSummary subscription")
		for {
			select {
			case <-ch:
				meshModelSummary, err := r.getMeshModelSummary(ctx, provider, selector)
				if err != nil {
					logrus.Error(ErrMeshModelSummarySubscription(err))
					break
				}
				respChan <- meshModelSummary
			case <-ctx.Done():
				close(respChan)
				close(ch)
				r.Log.Info("Closing MeshModelSummary subscription")
				return
			}
		}
	}()

	return respChan, nil
}

func (r *Resolver) getMeshModelSummary(ctx context.Context, _ models.Provider, selector model.MeshModelSummarySelector) (*model.MeshModelSummary, error) {
	regManager, ok := ctx.Value(models.RegistryManagerKey).(*meshmodel.RegistryManager)
	summary := &model.MeshModelSummary{}
	if !ok {
		err := errors.New("unable to get registry manager from context")
		return nil, ErrGettingRegistryManager(err)
	}
	switch selector.Type {
	case "components":
		components := getMeshModelComponents(regManager)
		summary.Components = components
		summary.Relationships = []*model.MeshModelRelationship{}
	case "relationships":
		relationships := getMeshModelRelationships(regManager)
		summary.Relationships = relationships
		summary.Components = []*model.MeshModelComponent{}
	}
	return summary, nil
}

type MeshModelComponentResponse struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type MeshModelRelationshipResponse struct {
	Name    string   `json:"name"`
	Subtype []string `json:"subType"`
}

func getMeshModelRelationships(regManager *meshmodel.RegistryManager) []*model.MeshModelRelationship {
	res := regManager.GetEntities(&v1alpha1.RelationshipFilter{})
	relationships := make([]*model.MeshModelRelationship, 0)
	var relmap = make(map[string]*MeshModelRelationshipResponse)
	for _, r := range res {
		def, _ := r.(v1alpha1.RelationshipDefinition)
		if relmap[def.Kind] == nil {
			relmap[def.Kind] = &MeshModelRelationshipResponse{
				Name:    def.Kind,
				Subtype: []string{def.SubType},
			}
		} else {
			relmap[def.Kind].Subtype = append(relmap[def.Kind].Subtype, def.SubType)
		}
	}
	for _, x := range relmap {
		x.Subtype = filterUniqueElementsArray(x.Subtype)
	}
	for _, x := range relmap {
		relationships = append(relationships, &model.MeshModelRelationship{
			Name:  x.Name,
			Count: len(x.Subtype),
		})
	}
	return relationships
}

func getMeshModelComponents(regManager *meshmodel.RegistryManager) []*model.MeshModelComponent {
	res := regManager.GetEntities(&v1alpha1.ComponentFilter{})
	components := make([]*model.MeshModelComponent, 0)
	var response = make(map[string]*MeshModelComponentResponse)
	for _, r := range res {
		def, _ := r.(v1alpha1.ComponentDefinition)
		if response[def.Model.Name] == nil {
			response[def.Model.Name] = &MeshModelComponentResponse{
				Name:     def.Model.Name,
				Versions: []string{def.Model.Version},
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
			Name:  x.Name,
			Count: len(x.Versions),
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
