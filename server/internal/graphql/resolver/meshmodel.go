package resolver

import (
	"context"
	"errors"

	"github.com/layer5io/meshery/server/internal/graphql/model"
	"github.com/layer5io/meshery/server/models"
	registry "github.com/layer5io/meshkit/models/meshmodel/registry"
	regv1alpha3 "github.com/layer5io/meshkit/models/meshmodel/registry/v1alpha3"
	regv1beta1 "github.com/layer5io/meshkit/models/meshmodel/registry/v1beta1"
	"github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/component"
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
					r.Log.Error(ErrMeshModelSummarySubscription(err))
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
	regManager, ok := ctx.Value(models.RegistryManagerKey).(*registry.RegistryManager)
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

func getMeshModelRelationships(regManager *registry.RegistryManager) []*model.MeshModelRelationship {
	res, _, _, _ := regManager.GetEntities(&regv1alpha3.RelationshipFilter{})
	relationships := make([]*model.MeshModelRelationship, 0)
	var relmap = make(map[string]*MeshModelRelationshipResponse)
	for _, r := range res {
		def, _ := r.(*relationship.RelationshipDefinition)
		if relmap[string(def.Kind)] == nil {
			relmap[string(def.Kind)] = &MeshModelRelationshipResponse{
				Name:    string(def.Kind),
				Subtype: []string{def.SubType},
			}
		} else {
			relmap[string(def.Kind)].Subtype = append(relmap[string(def.Kind)].Subtype, def.SubType)
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

func getMeshModelComponents(regManager *registry.RegistryManager) []*model.MeshModelComponent {
	res, _, _, _ := regManager.GetEntities(&regv1beta1.ComponentFilter{})
	components := make([]*model.MeshModelComponent, 0)
	var response = make(map[string]*MeshModelComponentResponse)
	for _, r := range res {
		def, _ := r.(*component.ComponentDefinition)
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
