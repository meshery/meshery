package core

import (
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"

	"github.com/layer5io/meshery/server/internal/store"
)

type capability struct {
	ID string `json:"id,omitempty"`
	// Host is the address of the service registering the capability
	Host string `json:"host,omitempty"`
}

// SetID sets the ID of the capability
func (cap *capability) SetID(id string) {
	cap.ID = id
}

// GetID returns the ID of the capability
func (cap *capability) GetID() string {
	return cap.ID
}

type ComponentCapability struct {
	v1alpha1.Component
	capability
}

// RegisterComponent will register a component definition into the database
func RegisterComponent(data []byte, host string) (err error) {
	var component ComponentCapability

	if err = json.Unmarshal(data, &component); err != nil {
		return
	}
	component.Host = host
	// Store it in the global store
	key := fmt.Sprintf(
		"/meshery/registry/definition/%s/%s/%s",
		component.APIVersion,
		component.Kind,
		component.Metadata["name"].(string),
	)
	store.Set(key, &component)
	return
}

// GetComponents return all of the components
func GetComponents() (caps []ComponentCapability) {
	key := "/meshery/registry/definition/core.meshery.io/v1alpha1/ComponentDefinition"

	res := store.PrefixMatch(key)
	for _, cc := range res {
		casted, ok := cc.(*ComponentCapability)
		if ok {
			caps = append(caps, *casted)
		}
	}

	return
}
