package system

import (
	"github.com/meshery/meshery/mesheryctl/pkg/utils"
)

var providerFlag string

// validateComponents checks if the specified components are valid according to utils.Services
func validateComponents(components []string) error {
	for _, component := range components {
		if utils.Services[component].Image == "" {
			return ErrInvalidComponent(component)
		}
	}
	return nil
}
