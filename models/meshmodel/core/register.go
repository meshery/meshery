package core

import (
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/internal/store"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
)

type RegistrantMetadata struct {
	// Host is the address of the grpc service of the registrant
	Host string `json:"host,omitempty"`
}

type CapabilityDefinition struct {
	v1alpha1.TypeMeta `json:",inline"`
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
}

type Capability struct {
	ID                 string `json:"id,omitempty"`
	Restricted         bool   `json:"restricted,omitempty"`
	RegistrantMetadata `json:"registrant_metadata,omitempty"`
	// A Capability is generic and need not be adhering to any sort of spec
	CapabilityDefinition interface{} `json:"capability_definition,omitempty"`
}

// SetID sets the ID of the capability
func (cap *Capability) SetID(id string) {
	cap.ID = id
}

// GetID returns the ID of the capability
func (cap *Capability) GetID() string {
	return cap.ID
}

func RegisterCapability(data []byte) (err error) {
	var cap Capability
	if err = json.Unmarshal(data, &cap); err != nil {
		return
	}

	capDef := cap.CapabilityDefinition.(map[string]interface{})
	// Store it in the global store
	key := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		capDef["apiVersion"],
		capDef["kind"],
		capDef["metadata"].(map[string]interface{})["name"],
	)

	store.Set(key, &cap)
	return nil
}

// GetComponents return all of the components
func GetComponents() (caps []v1alpha1.Component) {
	regex := `meshery\/registry\/definition\/.+?(?=\/).+?(?=\/)\/ComponentDefinition`
	res := store.RegexMatch(regex)
	for _, wc := range res {
		casted, ok := (wc.(*Capability))
		def := (casted.CapabilityDefinition).(v1alpha1.Component)
		if ok {
			caps = append(caps, def)
		}
	}

	return
}

// GetComponent takes in a component name and will return a SLICE of all of the components
// registered against the name
func GetComponent(name string) (comps []v1alpha1.Component) {
	regex := `meshery\/registry\/definition\/.+?(?=\/).+?(?=\/)\/ComponentDefinition/` + name
	res := store.RegexMatch(regex)
	for _, wc := range res {
		casted, ok := (wc.(*Capability))
		def := (casted.CapabilityDefinition).(v1alpha1.Component)
		if ok {
			comps = append(comps, def)
		}
	}
	return
}
