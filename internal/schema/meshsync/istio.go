package meshsync

import (
	istiobetav1 "istio.io/api/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
)

type Service struct {
	Name string
	Type ServiceType
	Pods []corev1.Pod
}

type Policy struct {
	Authorizations         []istiobetav1.AuthorizationPolicy
	PeerAuthentications    []istiobetav1.PeerAuthentication
	RequestAuthentications []istiobetav1.RequestAuthentication
}

type Workload struct {
	Name     string
	Labels   map[string]string
	Service  Service
	Sidecar  istiobetav1.Sidecar
	Policies []Policy
}

type Rule struct {
	Host          string
	CreatedAt     TimeStamp
	TrafficPolicy istiobetav1.TrafficPolicy
	Subset        []istiobetav1.Subset
}

type WorkloadGroup struct {
	Workloads []Workload
	Rules     []istiobetav1.DestinationRule
}

type Namespace struct {
	Name            string
	MonitoredBy     bool
	Gateways        []istiobetav1.Gateway
	VirtualServices []istiobetav1.VirtualService
	WorkloadGroups  []WorkloadGroup
}

type DataPlane struct {
	Namespaces []Namespace
}

type ControlPlane struct {
	Istiod          corev1.Service
	IngressGateways []corev1.Service
	EgressGateways  []corev1.Service
	Addons          []corev1.Service
}

type IstioInstance struct {
	ID           ID
	ControlPlane ControlPlane
	Namespaces   []Namespace
}
