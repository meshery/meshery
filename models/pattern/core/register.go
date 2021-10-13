package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshkit/models/oam/core/v1alpha1"
	"github.com/layer5io/meshkit/utils/kubernetes"
	"github.com/layer5io/meshkit/utils/manifests"
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

	return registerMesheryServerOAM(rootPath, []string{"application", "service.k8s"}, RegisterWorkload)
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
		defFile, err := ioutil.ReadFile(defpath)
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
		schemaFile, err := ioutil.ReadFile(schemapath)
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

// GetK8Components returns all the generated definitions and schemas for available api resources
func GetK8Components(config []byte, ctx string) (*manifests.Component, error) {
	cli, err := kubernetes.New(config)
	if err != nil {
		return nil, ErrGetK8sComponents(err)
	}
	req := cli.KubeClient.RESTClient().Get().RequestURI("/openapi/v2")
	k8version, err := cli.KubeClient.ServerVersion()
	if err != nil {
		return nil, ErrGetK8sComponents(err)
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
	manifest := string(content)
	man, err := manifests.GenerateComponents(manifest, manifests.K8s, manifests.Config{
		Name: "Kubernetes",
		Filter: manifests.CrdFilter{
			IsJson:        true,
			OnlyRes:       apiResources, //When crd or api-resource names are directly given, we dont need NameFilter
			RootFilter:    []string{"$.definitions"},
			VersionFilter: []string{"$[0]"},
			GroupFilter:   []string{"$[0]"},
			ItrFilter:     []string{"$..[\"x-kubernetes-group-version-kind\"][?(@.kind"},
			ItrSpecFilter: []string{"$[0][?(@[\"x-kubernetes-group-version-kind\"][0][\"kind\"]"},
			ResolveFilter: []string{"--resolve", "$"},
			GField:        "group",
			VField:        "version",
		},
		K8sVersion: k8version.String(),
		ModifyDefSchema: func(s1, s2 *string) {
			var schema map[string]interface{}
			err := json.Unmarshal([]byte(*s2), &schema)
			if err != nil {
				return
			}
			prop, ok := schema["properties"].(map[string]interface{})
			if !ok {
				return
			}
			// The schema generated has few fields that are not required and can break things, so they are removed here
			delete(prop, "apiVersion")
			delete(prop, "metadata")
			delete(prop, "kind")
			delete(prop, "status")
			schema["properties"] = prop
			schema["$schema"] = "http://json-schema.org/draft-04/schema"
			b, err := json.Marshal(schema)
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

// getAPIRes gets all the available api resources from kube-api server. It is equivalent to the output of `kubectl api-resources`
func getAPIRes(cli *kubernetes.Client) ([]string, error) {
	var apiRes []string
	lists, err := cli.KubeClient.DiscoveryClient.ServerPreferredResources()
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		for _, name := range list.APIResources {
			apiRes = append(apiRes, name.Kind)
		}
	}
	return apiRes, nil
}
