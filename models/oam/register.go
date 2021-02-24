package oam

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshery/models/oam/core/v1alpha1"
	"github.com/sirupsen/logrus"
)

type genericCapability struct {
	// OAMRefSchema is the json schema for the workload
	OAMRefSchema string `json:"oam_ref_schema,omitempty"`

	// Host is the address of the grpc service registering the capability
	Host string `json:"host,omitempty"`
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

	v, _ := store.Get(key)
	logrus.Debugf("Stored workload with key %s: %+v", key, v)

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

	v, _ := store.Get(key)
	logrus.Debugf("Stored trait: %+v", v)

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

	v, _ := store.Get(key)
	logrus.Debugf("Stored scope: %+v", v)

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

	return registerMesheryServerOAM(rootPath, []string{}, RegisterWorkload)
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
