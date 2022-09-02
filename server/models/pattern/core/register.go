package core

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	cueJson "cuelang.org/go/encoding/json"
	"github.com/layer5io/meshery/server/internal/store"
	"github.com/layer5io/meshery/server/models/pattern/patterns/k8s"
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

// RegisterWorkload will register a workload definition into the database
func RegisterWorkload(data []byte) (err error) {
	var workload WorkloadCapability

	if err = json.Unmarshal(data, &workload); err != nil {
		return
	}

	// Store it in the global store
	key := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		workload.OAMDefinition.APIVersion,
		workload.OAMDefinition.Kind,
		workload.OAMDefinition.Name,
	)

	schema := map[string]interface{}{}
	_ = json.Unmarshal([]byte(workload.OAMRefSchema), &schema)
	if k8s.Format {
		k8s.Format.Prettify(schema)
	}
	temp, _ := json.Marshal(schema)
	workload.OAMRefSchema = string(temp)
	if workload.Metadata == nil {
		workload.Metadata = map[string]string{}
	}
	workload.Metadata["display.ui.meshery.io/name"], _ = schema["title"].(string)

	store.Set(key, &workload)

	return
}

// TraitCapability is the struct for capturing the workload definition
// of a particular type
type TraitCapability struct {
	OAMDefinition v1alpha1.TraitDefinition `json:"oam_definition,omitempty"`

	genericCapability
}

// RegisterTrait will register a TraitDefinition into the database
func RegisterTrait(data []byte) (err error) {
	var trait TraitCapability

	if err = json.Unmarshal(data, &trait); err != nil {
		return
	}

	// Store it in the global store
	key := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		trait.OAMDefinition.APIVersion,
		trait.OAMDefinition.Kind,
		trait.OAMDefinition.Name,
	)

	schema := map[string]interface{}{}
	_ = json.Unmarshal([]byte(trait.OAMRefSchema), &schema)
	if trait.Metadata == nil {
		trait.Metadata = map[string]string{}
	}
	trait.Metadata["display.ui.meshery.io/name"], _ = schema["title"].(string)

	store.Set(key, &trait)

	return
}

// ScopeCapability is the struct for capturing the ScopeDefinition
// of a particular type
type ScopeCapability struct {
	OAMDefinition v1alpha1.ScopeDefinition `json:"oam_definition,omitempty"`

	genericCapability
}

// RegisterScope will register a Scope definition into the database
func RegisterScope(data []byte) (err error) {
	var scope ScopeCapability

	if err = json.Unmarshal(data, &scope); err != nil {
		return
	}

	// Store it in the global store
	key := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		scope.OAMDefinition.APIVersion,
		scope.OAMDefinition.Kind,
		scope.OAMDefinition.Name,
	)

	schema := map[string]interface{}{}
	_ = json.Unmarshal([]byte(scope.OAMRefSchema), &schema)
	if scope.Metadata == nil {
		scope.Metadata = map[string]string{}
	}
	scope.Metadata["display.ui.meshery.io/name"], _ = schema["title"].(string)

	store.Set(key, &scope)

	return
}

// GetWorkloads return all of the workloads
func GetWorkloads() (caps []WorkloadCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/WorkloadDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(*WorkloadCapability)
		if ok {
			caps = append(caps, *casted)
		}
	}

	return
}

// GetTraits return all of the traits
func GetTraits() (traits []TraitCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/TraitDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(*TraitCapability)
		if ok {
			traits = append(traits, *casted)
		}
	}

	return
}

// GetScopes return all of the scopes
func GetScopes() (scopes []ScopeCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/ScopesDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(*ScopeCapability)
		if ok {
			scopes = append(scopes, *casted)
		}
	}

	return
}

// GetWorkload takes in a workload name and will return a SLICE of all of the workloads
// registered against the name
func GetWorkload(name string) (w []WorkloadCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/WorkloadDefinition/" + name

	res := store.GetAll(key)
	for _, wc := range res {
		casted, ok := wc.(*WorkloadCapability)
		if ok {
			w = append(w, *casted)
		}
	}

	return
}

// GetTrait takes in a trait name and will return a SLICE of all of the traits
// registered against the name
func GetTrait(name string) (t []TraitCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/TraitDefinition/" + name

	res := store.GetAll(key)
	for _, wc := range res {
		casted, ok := wc.(*TraitCapability)
		if ok {
			t = append(t, *casted)
		}
	}

	return
}

// GetScope takes in a scope name and will return a SLICE of all of the scopes
// registered against the name
func GetScope(name string) (s []ScopeCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/ScopeDefinition/" + name

	res := store.GetAll(key)
	for _, wc := range res {
		casted, ok := wc.(*ScopeCapability)
		if ok {
			s = append(s, *casted)
		}
	}

	return
}

// GetWorkloadsByK8sAPIVersionKind takes in kubernetes API version of a resource and its kind
// and returns all of the resources that matches that description
//
// Meshery core resources DO NOT have these attached to themselves, hence in order to get them just pass
// in empty string for both the parameters
func GetWorkloadsByK8sAPIVersionKind(apiVersion, kind string) (w []WorkloadCapability) {
	wcs := GetWorkloads()

	for _, wc := range wcs {
		aversion := wc.OAMDefinition.Spec.Metadata["k8sAPIVersion"]
		if (aversion == apiVersion || aversion == "/"+apiVersion) && wc.OAMDefinition.Spec.Metadata["k8sKind"] == kind {
			w = append(w, wc)
		}
	}

	return
}

// GetTraitsByK8sAPIVersionKind takes in kubernetes API version of a resource and its kind
// and returns all of the resources that matches that description
//
// Meshery core resources DO NOT have these attached to themselves, hence in order to get them just pass
// in empty string for both the parameters
func GetTraitsByK8sAPIVersionKind(apiVersion, kind string) (t []TraitCapability) {
	tcs := GetTraits()

	for _, tc := range tcs {
		aversion := tc.OAMDefinition.Spec.Metadata["k8sAPIVersion"]
		if (aversion == apiVersion || aversion == "/"+apiVersion) && tc.OAMDefinition.Spec.Metadata["k8sKind"] == kind {
			t = append(t, tc)
		}
	}

	return
}

// GetScopesByK8sAPIVersionKind takes in kubernetes API version of a resource and its kind
// and returns all of the resources that matches that description
//
// Meshery core resources DO NOT have these attached to themselves, hence in order to get them just pass
// in empty string for both the parameters
func GetScopesByK8sAPIVersionKind(apiVersion, kind string) (s []ScopeCapability) {
	scs := GetScopes()

	for _, sc := range scs {
		aversion := sc.OAMDefinition.Spec.Metadata["k8sAPIVersion"]
		if (aversion == apiVersion || aversion == "/"+apiVersion) && sc.OAMDefinition.Spec.Metadata["k8sKind"] == kind {
			s = append(s, sc)
		}
	}

	return
}

// GetWorkloadByID takes an id of a workload and returns the workload
func GetWorkloadByID(name, id string) *WorkloadCapability {
	res := GetWorkload(name)
	for _, f := range res {
		if f.GetID() == id {
			return &f
		}
	}

	return nil
}

// GetTraitByID takes an id of a trait and returns the trait
func GetTraitByID(name, id string) *TraitCapability {
	res := GetTrait(name)
	for _, f := range res {
		if f.GetID() == id {
			return &f
		}
	}

	return nil
}

// GetScopeByID takes an id of a scope and returns the scope
func GetScopeByID(name, id string) (w *ScopeCapability) {
	res := GetScope(name)
	for _, f := range res {
		if f.GetID() == id {
			return &f
		}
	}

	return nil
}

// RegisterMesheryOAMTraits will register local meshery traits with meshery server
func RegisterMesheryOAMTraits() error {
	// rootPath is the relative path to the traits
	// if the file is moved then this path MUST be changed
	// accordingly
	rootPath, _ := filepath.Abs("../oam/traits")

	return registerMesheryServerOAM(rootPath, []string{"meshmap"}, RegisterTrait)
}

// RegisterMesheryOAMWorkloads will register local meshery workloads with meshery server
func RegisterMesheryOAMWorkloads() error {
	// rootPath is the relative path to the workloads
	// if the file is moved then this path MUST be changed
	// accordingly
	rootPath, _ := filepath.Abs("../oam/workloads")

	return registerMesheryServerOAM(rootPath, []string{"application"}, RegisterWorkload)
}

// registerMesheryServerOAM will read the oam definition file and its corresponding schema file
// and then will call given regFn for registering the serialized data
//
// registerMesheryServerOAM expects that if a construct called "meshmap" is given then a definition
// file like "meshmap_definition.json" and "meshmap.meshery.layer5.io.schema.json" will be present
// in the given rootPath
func registerMesheryServerOAM(rootPath string, constructs []string, regFn func([]byte) error) error {
	var errs []string

	for _, construct := range constructs {
		path := filepath.Join(rootPath, construct)

		// Read the file definition file
		defpath := fmt.Sprintf("%s_definition.json", path)
		defFile, err := os.ReadFile(defpath)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		var tempDef map[string]interface{}
		err = json.Unmarshal(defFile, &tempDef)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		// Read the schema file
		schemapath := fmt.Sprintf("%s.meshery.layer5.io.schema.json", path)
		schemaFile, err := os.ReadFile(schemapath)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		// Create the pseudo structure for serializing data
		data := map[string]interface{}{
			"oam_ref_schema": string(schemaFile),
			"oam_definition": tempDef,
			"host":           "<none-local>",
			"metadata": map[string]string{
				"adapter.meshery.io/name": "core",
			},
		}

		// Serialize the data
		byt, err := json.Marshal(data)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}

		// Call the regFn which is supposed to process the generated bytes
		if err := regFn(byt); err != nil {
			errs = append(errs, err.Error())
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("%s", strings.Join(errs, "\n"))
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
			schema["$schema"] = "http://json-schema.org/draft-04/schema"

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
			b, err = referenceResolver.ResolveReferences(b, definitions)
			if err != nil {
				return
			}
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
