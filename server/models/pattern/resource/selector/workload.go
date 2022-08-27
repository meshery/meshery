package selector

import (
	"github.com/layer5io/meshery/server/internal/store"
	"github.com/layer5io/meshery/server/models/pattern/core"
)

func (s *Selector) Workload(name string) (core.WorkloadCapability, bool) {
	data := store.GetAll(generateWorkloadKey(name))
	workloads := convertValueInterfaceSliceToWorkloadSlice(data)

	filteredWorkloads, typ := filterWorkloadByType(workloads)

	if typ == CoreResource {
		return s.selectCoreWorkload(filteredWorkloads)
	}

	if typ == K8sResource {
		return s.selectK8sWorkload(filteredWorkloads)
	}

	return s.selectMeshWorkload(filteredWorkloads)
}

// selectCoreWorkload selects a core workload - first workload from the list is selected
// at the moment
func (s *Selector) selectCoreWorkload(ws []core.WorkloadCapability) (core.WorkloadCapability, bool) {
	if len(ws) == 0 {
		return core.WorkloadCapability{}, false
	}

	return ws[0], true
}

// selectK8sWorkload selects a k8s workload - resource with latest apiVersion will be selected
//
// TODO: Change implementation to look into the cluster to identify appropriate
func (s *Selector) selectK8sWorkload(ws []core.WorkloadCapability) (core.WorkloadCapability, bool) {
	if len(ws) == 0 {
		return core.WorkloadCapability{}, false
	}

	var selected *core.WorkloadCapability
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

// selectMeshWorkload will look through all the candidate workloads and will return most suitable workload
//
// The selection of workload is based on the following criterion
//   - If a workload definition is of a service mesh version that is NOT in the cluster then reject that candidate
//   - If
func (s *Selector) selectMeshWorkload(ws []core.WorkloadCapability) (core.WorkloadCapability, bool) {
	if len(ws) == 0 {
		return core.WorkloadCapability{}, false
	}

	meshName, meshVersion := s.helpers.GetServiceMesh()

	// If we failed to get service mesh version or we didn't find any service mesh
	// running in the cluster then proceed to randomly select a workload
	if meshName == "" || meshVersion == "" {
		return ws[0], true
	}

	var selected *core.WorkloadCapability
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

	return core.WorkloadCapability{}, false
}

// filterWorkloadByType filters workloads based on their type and returns a slice
// of filtered workloads along with the type
//
// Following is the priority of the resource types:
// pattern.meshery.io/core > pattern.meshery.io/mesh/workload > pattern.meshery.io/k8s
func filterWorkloadByType(ws []core.WorkloadCapability) ([]core.WorkloadCapability, string) {
	cor := []core.WorkloadCapability{}
	mesh := []core.WorkloadCapability{}
	k8s := []core.WorkloadCapability{}

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

func getResourceType(metadata map[string]string) string {
	typ, ok := metadata["@type"]
	if !ok {
		// Legacy resource => For now mark it as core
		return CoreResource
	}

	return typ
}

func convertValueInterfaceSliceToWorkloadSlice(data []store.Value) []core.WorkloadCapability {
	res := []core.WorkloadCapability{}

	for _, el := range data {
		elc, ok := el.(*core.WorkloadCapability)
		if ok {
			res = append(res, *elc)
		}
	}

	return res
}
