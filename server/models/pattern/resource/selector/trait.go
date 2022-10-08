package selector

import (
	"github.com/layer5io/meshery/server/internal/store"
	"github.com/layer5io/meshery/server/models/pattern/core"
)

func (s *Selector) Trait(name string) (core.TraitCapability, bool) {
	data := store.GetAll(generateTraitKey(name))
	traits := convertValueInterfaceSliceToTraitSlice(data)

	filteredTraits, typ := filterTraitByType(traits)

	if typ == CoreResource {
		return s.selectCoreTrait(filteredTraits)
	}

	if typ == K8sResource {
		return s.selectK8sTrait(filteredTraits)
	}

	return s.selectMeshTrait(filteredTraits)
}

// selectCoreTrait selects a core trait - first trait from the list is selected
// at the moment
func (s *Selector) selectCoreTrait(ts []core.TraitCapability) (core.TraitCapability, bool) {
	if len(ts) == 0 {
		return core.TraitCapability{}, false
	}

	return ts[0], true
}

// selectK8sTrait selects a k8s trait - resource with latest apiVersion will be selected
//
// TODO: Change implementation to look into the cluster to identify appropriate
func (s *Selector) selectK8sTrait(ws []core.TraitCapability) (core.TraitCapability, bool) {
	if len(ws) == 0 {
		return core.TraitCapability{}, false
	}

	var selected *core.TraitCapability
	for _, w := range ws {
		if selected == nil {
			selected = &w
			continue
		}

		mversion := selected.OAMDefinition.Spec.Metadata["k8sAPIVersion"]
		version := w.OAMDefinition.Spec.Metadata["k8sAPIVersion"]

		if version > mversion {
			selected = &w
		}
	}

	return *selected, true
}

// selectMeshTrait will look through all the candidate traits and will return most suitable trait
//
// The selection of trait is based on the following criterion
//   - If a trait definition is of a service mesh version that is NOT in the cluster then reject that candidate
//   - If
func (s *Selector) selectMeshTrait(ws []core.TraitCapability) (core.TraitCapability, bool) {
	if len(ws) == 0 {
		return core.TraitCapability{}, false
	}

	meshName, meshVersion := s.helpers.GetServiceMesh()

	// If we failed to get service mesh version or we didn't find any service mesh
	// running in the cluster then proceed to randomly select a workload
	if meshName == "" || meshVersion == "" {
		return ws[0], true
	}

	var selected *core.TraitCapability
	for _, w := range ws {
		version := w.OAMDefinition.Spec.Metadata["meshVersion"]
		name := w.OAMDefinition.Spec.Metadata["meshName"]

		if version == meshVersion && name == meshName { // Potential candidate found
			if selected == nil {
				selected = &w
				continue
			}

			sAPIVersion := selected.OAMDefinition.Spec.Metadata["k8sAPIVersion"]
			apiVersion := w.OAMDefinition.Spec.Metadata["k8sAPIVersion"]

			if apiVersion > sAPIVersion {
				selected = &w
			}
		}
	}

	if selected != nil {
		return *selected, true
	}

	return core.TraitCapability{}, false
}

// filterTraitByType filters traits based on their type and returns a slice
// of filtered traits along with the type
//
// Following is the priority of the resource types:
// pattern.meshery.io/core > pattern.meshery.io/mesh/workload > pattern.meshery.io/k8s
func filterTraitByType(ws []core.TraitCapability) ([]core.TraitCapability, string) {
	cor := []core.TraitCapability{}
	mesh := []core.TraitCapability{}
	k8s := []core.TraitCapability{}

	for _, w := range ws {
		typ := getResourceType(w.OAMDefinition.Spec.Metadata)

		switch typ {
		case CoreResource:
			cor = append(cor, w)
		case MeshResource:
			mesh = append(mesh, w)
		case K8sResource:
			k8s = append(k8s, w)
		}
	}

	if len(cor) > 0 {
		return cor, CoreResource
	}

	if len(mesh) > 0 {
		return mesh, MeshResource
	}

	return k8s, K8sResource
}

func convertValueInterfaceSliceToTraitSlice(data []store.Value) []core.TraitCapability {
	res := []core.TraitCapability{}

	for _, el := range data {
		elc, ok := el.(*core.TraitCapability)
		if ok {
			res = append(res, *elc)
		}
	}

	return res
}
