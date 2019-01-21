package istio

import (
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type VirtualService struct {
	meta_v1.TypeMeta `json:",inline" yaml:",inline"`
	ObjectMeta       meta_v1.ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec             Spec               `json:"spec" yaml:"spec"`
}

type Spec struct {
	Hosts []string `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	Http  []HTTP   `json:"http,omitempty" yaml:"http,omitempty"`
}

type HTTP struct {
	Route []Route                           `json:"route,omitempty"`
	Match []Match                           `json:"match,omitempty"`
	Fault map[string]map[string]interface{} `json:"fault,omitempty" yaml:"fault,omitempty"`
}

type Route struct {
	Destinations map[string]string `json:"destination,omitempty" yaml:"destination,omitempty"`
	Weight       int               `json:"weight,omitempty" yaml:"weight,omitempty"`
}

type Match struct {
	Headers map[string]map[string]interface{} `json:"headers,omitempty" yaml:"headers,omitempty"`
}

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

func (in *VirtualService) DeepCopyInto(out *VirtualService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
}

func (in *VirtualService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}

func (in *VirtualService) DeepCopy() *VirtualService {
	if in == nil {
		return nil
	}
	out := new(VirtualService)
	in.DeepCopyInto(out)
	return out
}

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

func (in *VirtualServiceList) DeepCopy() *VirtualServiceList {
	if in == nil {
		return nil
	}
	out := new(VirtualServiceList)
	in.DeepCopyInto(out)
	return out
}

func (in *VirtualServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}

	return nil
}
