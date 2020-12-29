package meshsync

type Schema struct {
	Preference  Preference
	Cluster     []ClusterView
	ServiceMesh []ServiceMeshView
}

type ClusterView struct {
	ID          ID
	Name        string
	Namespaces  Count
	Nodes       Count
	CPU         Count
	RAM         Count
	Utilization ProgressMeter
}

type ServiceMeshView struct {
	ID        ID
	Name      string
	Istios    []IstioInstance
	StartedAt TimeStamp
}

type Preference struct {
	// All user settings go in here
}
