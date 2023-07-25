package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueJson "cuelang.org/go/encoding/json"
	"github.com/layer5io/meshery/server/internal/store"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type genericCapability struct {
	ID string `json:"id,omitempty"`

	// OAMRefSchema is the json schema for the workload
	OAMRefSchema string `json:"oam_ref_schema,omitempty"`

	// Host is the address of the grpc service registering the capability
	Host string `json:"host,omitempty"`

	Restricted bool `json:"restricted,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`
}

// TraitCapability is the struct for capturing the workload definition
// of a particular type
type TraitCapability struct {
	OAMDefinition v1alpha1.TraitDefinition `json:"oam_definition,omitempty"`

	genericCapability
}

// SetID sets the ID of the capability
func (cap *genericCapability) SetID(id string) {
	cap.ID = id
}

// GetID returns the ID of the capability
func (cap *genericCapability) GetID() string {
	return cap.ID
}

// WorkloadCapability is the struct for capturing the workload definition
// of a particular type
type WorkloadCapability struct {
	OAMDefinition v1alpha1.WorkloadDefinition `json:"oam_definition,omitempty"`

	genericCapability
}

// present in metadata."adapter.meshery.io/name". Ex- core,kubernetes,istio,linkerd,etc
type ComponentTypes struct {
	Names                     map[string]bool
	LatestVersionForComponent map[string]string
	mx                        sync.Mutex
}

func (c *ComponentTypes) Set(name string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.Names[name] = true
}
func (c *ComponentTypes) SetLatestVersion(typ string, ver string) {
	c.mx.Lock()
	defer c.mx.Unlock()
	c.LatestVersionForComponent[typ] = ver
}
func (c *ComponentTypes) Get() (names []string) {
	for n := range c.Names {
		names = append(names, n)
	}
	return
}

// ComponentTypesSingleton is initialized per meshery instance and acts as a helper middleware between client facing API and capability registry.
// Examples of names stored in this struct are: core,kubernetes,istio,linkerd
var ComponentTypesSingleton = ComponentTypes{
	Names:                     make(map[string]bool),
	LatestVersionForComponent: make(map[string]string),
}

const customResourceKey = "isCustomResource"

type crd struct {
	Items []crdhelper `json:"items"`
}
type crdhelper struct {
	Metadata map[string]interface{} `json:"metadata"`
}

// GetK8Components returns all the generated definitions and schemas for available api resources
func GetK8Components(ctxt context.Context, config []byte) (*manifests.Component, error) {
	cli, err := kubernetes.New(config)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	req := cli.KubeClient.RESTClient().Get().RequestURI("/openapi/v2")
	k8version, err := cli.KubeClient.ServerVersion()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	var customResources = make(map[string]bool)
	crdresult, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/apiextensions.k8s.io/v1/customresourcedefinitions").Do(context.Background()).Raw()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}

	var xcrd crd
	err = json.Unmarshal(crdresult, &xcrd)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	for _, item := range xcrd.Items {
		customResources[item.Metadata["name"].(string)] = true
	}
	res := req.Do(context.Background())
	content, err := res.Raw()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	apiResources, err := getAPIRes(cli)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}

	var arrAPIResources []string
	for res := range apiResources {
		arrAPIResources = append(arrAPIResources, res)
	}
	groups, err := getGroupsFromResource(cli) //change this
	if err != nil {
		return nil, err
	}
	manifest := string(content)
	man, err := manifests.GenerateComponents(ctxt, manifest, manifests.K8s, manifests.Config{
		Name: "Kubernetes",
		CrdFilter: manifests.NewCueCrdFilter(manifests.ExtractorPaths{
			NamePath:    `"x-kubernetes-group-version-kind"[0].kind`,
			IdPath:      `"x-kubernetes-group-version-kind"[0].kind`,
			VersionPath: `"x-kubernetes-group-version-kind"[0].version`,
			GroupPath:   `"x-kubernetes-group-version-kind"[0].group`,
			SpecPath:    "",
		}, true),
		ExtractCrds: func(manifest string) []string {
			crds := make([]string, 0)
			cuectx := cuecontext.New()
			cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
			parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
			definitions := parsedManifest.LookupPath(cue.ParsePath("definitions"))
			if err != nil {
				fmt.Printf("%v", err)
				return nil
			}
			for _, resource := range arrAPIResources {
				resource = strings.ToLower(resource)
				fields, err := definitions.Fields()
				if err != nil {
					fmt.Printf("%v\n", err)
					continue
				}
				for fields.Next() {
					fieldVal := fields.Value()
					kindCue := fieldVal.LookupPath(cue.ParsePath(`"x-kubernetes-group-version-kind"[0].kind`))
					if kindCue.Err() != nil {
						continue
					}
					kind, err := kindCue.String()
					kind = strings.ToLower(kind)
					if err != nil {
						fmt.Printf("%v", err)
						continue
					}
					if kind == resource {
						crd, err := fieldVal.MarshalJSON()
						if err != nil {
							fmt.Printf("%v", err)
							continue
						}
						crds = append(crds, string(crd))
					}
				}
			}
			return crds
		},
		K8sVersion: k8version.String(),
		ModifyDefSchema: func(s1, s2 *string) { //s1 is the definition and s2 is the schema
			var schema map[string]interface{}
			err = json.Unmarshal([]byte(*s2), &schema)
			if err != nil {
				return
			}
			prop, ok := schema["properties"].(map[string]interface{})
			if !ok {
				return
			}
			var def v1alpha1.WorkloadDefinition
			err := json.Unmarshal([]byte(*s1), &def)
			if err != nil {
				return
			}
			//Add additional info in metadata like whether the resource is namespaced or not. Or whether it is a custom resource.
			kind := strings.TrimSuffix(def.Spec.Metadata["k8sKind"], ".K8s")
			if apiResources[kind].Namespaced {
				def.Spec.Metadata["namespaced"] = "true"
			} else {
				def.Spec.Metadata["namespaced"] = "false"
			}
			def.Spec.Metadata[customResourceKey] = "false" //default
			for cr := range customResources {
				if groups[kind][cr] {
					def.Spec.Metadata[customResourceKey] = "true"
					break
				}
			}
			b, err := json.Marshal(def)
			if err != nil {
				return
			}
			*s1 = string(b)
			// The schema generated has few fields that are not required and can break things, so they are removed here
			delete(prop, "apiVersion")
			delete(prop, "metadata")
			delete(prop, "kind")
			delete(prop, "status")
			schema["properties"] = prop
			schema["$schema"] = "http://json-schema.org/draft-07/schema"

			b, err = json.Marshal(schema)
			if err != nil {
				return
			}

			cuectx := cuecontext.New()
			cueParsedManExpr, err := cueJson.Extract("", []byte(manifest))
			if err != nil {
				return
			}
			parsedManifest := cuectx.BuildExpr(cueParsedManExpr)
			definitions := parsedManifest.LookupPath(cue.ParsePath("definitions"))

			referenceResolver := manifests.ResolveOpenApiRefs{}
			cache := make(map[string][]byte)
			newResRef, err := referenceResolver.ResolveReferences(b, definitions, cache)
			if err != nil {
				return
			}
			b, _ = json.Marshal(newResRef)
			resolved := make(map[string]interface{})
			err = json.Unmarshal(b, &resolved)
			if err != nil {
				return
			}

			b, err = json.Marshal(resolved)
			if err != nil {
				return
			}
			*s2 = string(b)
		},
	})
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	return man, nil
}

// DeleteK8sWorkloads deletes the registered in memory k8s workloads for a given k8s contextID.
func DeleteK8sWorkloads(ctx string) {
	//Iterate through entire store
	vals := store.PrefixMatch("")
	for _, val := range vals {
		value, ok := val.(*WorkloadCapability)
		if !ok {
			continue
		}
		var workload = value
		//delete only the ones with given context in metadata
		if workload.OAMDefinition.Spec.Metadata["@type"] == "pattern.meshery.io/k8s" && workload.Metadata["io.meshery.ctxid"] == ctx {
			key := fmt.Sprintf(
				"/meshery/registry/definition/%s/%s/%s",
				workload.OAMDefinition.APIVersion,
				workload.OAMDefinition.Kind,
				workload.OAMDefinition.Name,
			)
			store.Delete(key, value)
		}
	}
}

// TODO: To be moved in meshkit
// getAPIRes gets all the available api resources from kube-api server. It is equivalent to the output of `kubectl api-resources`
// Returns a map of api resources with key as api-resource kind and value as api-resource object
func getAPIRes(cli *kubernetes.Client) (map[string]v1.APIResource, error) {
	var apiRes = make(map[string]v1.APIResource)
	lists, err := cli.KubeClient.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		for _, name := range list.APIResources {
			apiRes[name.Kind] = name
		}
	}
	return apiRes, nil
}

// TODO: To be moved in meshkit
// Return a set of resource.groups(with resourcename appended) for a given kind //Example: {"CSIDriver":{"csidrivers.storage.k8s.io"}}
func getGroupsFromResource(cli *kubernetes.Client) (gr map[string]map[string]bool, err error) {
	gr = make(map[string]map[string]bool)

	var gl v1.APIGroupList
	gs, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis").Do(context.Background()).Raw()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(gs, &gl)
	if err != nil {
		return nil, err
	}

	for _, g := range gl.Groups {
		groupName := g.Name
		var apig v1.APIGroup
		apigbytes, err := cli.KubeClient.RESTClient().Get().RequestURI("/apis/" + groupName).Do(context.Background()).Raw()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(apigbytes, &apig)
		if err != nil {
			return nil, err
		}
		apiRes, err := cli.KubeClient.DiscoveryClient.ServerResourcesForGroupVersion(apig.PreferredVersion.GroupVersion)
		if err != nil {
			return nil, err
		}
		for _, res := range apiRes.APIResources {
			if gr[res.Kind] == nil {
				gr[res.Kind] = make(map[string]bool)
			}
			gr[res.Kind][res.Name+"."+groupName] = true
		}
	}
	return
}
