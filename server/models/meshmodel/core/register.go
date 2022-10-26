package core

import (
	"encoding/json"
	"fmt"

	"github.com/layer5io/meshery/server/internal/store"
	"github.com/layer5io/meshkit/database"
	"github.com/layer5io/meshkit/models/meshmodel/core/v1alpha1"
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
	v1alpha1.Component `gorm:"embedded"`
	capability         `gorm:"embedded"`
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
		component.Component.GetMetadataValue("name").(string),
	)
	store.Set(key, &component)
	return
}

// GetComponents return all of the components
func GetComponents(dbHandler *database.Handler) (caps []ComponentCapability) {
	err := dbHandler.DB.Find(&caps).Error
	if err != nil {
		fmt.Println("error: ", err.Error())
	}
	// key := "/meshery/registry/definition/core.meshery.io/v1alpha1/ComponentDefinition"

	// res := store.PrefixMatch(key)
	// for _, cc := range res {
	// 	casted, ok := cc.(*ComponentCapability)
	// 	if ok {
	// 		caps = append(caps, *casted)
	// 	}
	// }

	return
}
