package istio

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// VirtualService - represents VirtualService in Istio
type VirtualService struct {
	meta_v1.TypeMeta `json:",inline" yaml:",inline"`
	ObjectMeta       meta_v1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec             Spec               `json:"spec" yaml:"spec"`
}

// Spec - Spec block in VirtualService definitions
type Spec struct {
	Hosts []string `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Http  []HTTP   `json:"http,omitempty" yaml:"http,omitempty"`
}

// HTTP - represents HTTP block in VirtualService definitions
type HTTP struct {
	Route []Route                           `json:"route,omitempty"`
	Match []Match                           `json:"match,omitempty"`
	Fault map[string]map[string]interface{} `json:"fault,omitempty" yaml:"fault,omitempty"`
}

// Route - represents Route block in VirtualService definitions
type Route struct {
	Destinations map[string]string `json:"destination,omitempty" yaml:"destination,omitempty"`
	Weight       int               `json:"weight,omitempty" yaml:"weight,omitempty"`
}

// Match - represents Match block in VirtualService definitions
type Match struct {
	Headers map[string]map[string]interface{} `json:"headers,omitempty" yaml:"headers,omitempty"`
}

// VirtualServiceList - represents a collection of VirtualServices
type VirtualServiceList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []VirtualService `json:"items"`
}

const (
	virtualServices = "virtualservices"
	virtualservice  = "VirtualService"
)

var (
	istioConfigGroupVersion = schema.GroupVersion{
		Group:   "config.istio.io",
		Version: "v1alpha2",
	}
	istioNetworkingGroupVersion = schema.GroupVersion{
		Group:   "networking.istio.io",
		Version: "v1alpha3",
	}

	istioKnownTypes = map[string]struct {
		object       runtime.Object
		collection   runtime.Object
		groupVersion *schema.GroupVersion
	}{
		virtualServices: {
			object: &VirtualService{
				TypeMeta: meta_v1.TypeMeta{
					Kind:       virtualservice,
					APIVersion: istioNetworkingGroupVersion.Group + "/" + istioNetworkingGroupVersion.Version,
				},
			},
			collection:   &VirtualServiceList{},
			groupVersion: &istioNetworkingGroupVersion,
		},
	}
)

// DeepCopyInto - method to deep copy one VS into another
func (in *VirtualService) DeepCopyInto(out *VirtualService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

// DeepCopyObject - helper for Deep Copy
func (in *VirtualService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}

// DeepCopy - method to deep copy one VS to create a new VS
func (in *VirtualService) DeepCopy() *VirtualService {
	if in == nil {
		return nil
	}
	out := new(VirtualService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto - method to deep copy one VS list into another
func (in *VirtualServiceList) DeepCopyInto(out *VirtualServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy - method to deep copy one VS list to create a new VS list
func (in *VirtualServiceList) DeepCopy() *VirtualServiceList {
	if in == nil {
		return nil
	}
	out := new(VirtualServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject - helper for Deep Copy
func (in *VirtualServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}
