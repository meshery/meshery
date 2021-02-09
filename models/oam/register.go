package oam

import (
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/models/oam/core/v1alpha1"
	"github.com/layer5io/meshery/store"
	"github.com/sirupsen/logrus"
)

type genericCapability struct {
	// OAMRefSchema is the json schema for the workload
	OAMRefSchema string `json:"oam_ref_schema,omitempty"`

	// Host is the address of the service registering the capability
	//
	// this can be either in format grpc://hostname:port, http://hostname:port
	// https://hostname:port
	//
	// If the protocol of choice is grpc then it should conform the
	// meshops.proto and if protocol of choice is http then /oam/operation
	// endpoint should be available for "POST" method
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
