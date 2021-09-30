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
	// OAMRefSchema is the json schema for the workload
	OAMRefSchema string `json:"oam_ref_schema,omitempty"`

	// Host is the address of the grpc service registering the capability
	Host string `json:"host,omitempty"`

	Restricted bool `json:"restricted,omitempty"`

	Metadata map[string]string `json:"metadata,omitempty"`
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
	store.Set(key, workload)

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
	store.Set(key, trait)

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
	store.Set(key, scope)

	return
}

// GetWorkloads return all of the workloads
func GetWorkloads() (caps []WorkloadCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/WorkloadDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(WorkloadCapability)
		if ok {
			caps = append(caps, casted)
		}
	}

	return
}

// GetTraits return all of the traits
func GetTraits() (traits []TraitCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/TraitDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(TraitCapability)
		if ok {
			traits = append(traits, casted)
		}
	}

	return
}

// GetScopes return all of the scopes
func GetScopes() (scopes []ScopeCapability) {
	key := "/meshery/registry/definition/core.oam.dev/v1alpha1/ScopesDefinition"

	res := store.PrefixMatch(key)
	for _, wc := range res {
		casted, ok := wc.(ScopeCapability)
		if ok {
			scopes = append(scopes, casted)
		}
	}

	return
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
			RootFilter:    []string{"$.definitions", "--resolve", "$"},
			VersionFilter: []string{"$[0]"},
			GroupFilter:   []string{"$[0]"},
			ItrFilter:     "$..[\"x-kubernetes-group-version-kind\"][?(@.kind",
			ItrSpecFilter: "$[0][?(@[\"x-kubernetes-group-version-kind\"][0][\"kind\"]",
			GField:        "group",
			VField:        "version",
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
		value, ok := val.(WorkloadCapability)
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
