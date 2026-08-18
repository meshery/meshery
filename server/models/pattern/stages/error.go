package stages

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrResolveReferenceCode       = "meshery-server-1361"
	ErrYAMLUnmarshalCode          = "meshery-server-1436"
	ErrDependencyNotSatisfiedCode = "meshery-server-1480"
)

// ErrDependencyNotSatisfied reports a component that was not deployed because a
// component it declared a dependency on did not deploy successfully.
func ErrDependencyNotSatisfied(component, dependency string) error {
	return errors.New(
		ErrDependencyNotSatisfiedCode,
		errors.Alert,
		[]string{fmt.Sprintf("Component %q was not deployed", component)},
		[]string{fmt.Sprintf("%q declares a dependency on %s, which did not deploy successfully. Meshery deploys a component only once every component it depends on is in place, so %q was withheld.", component, dependency, component)},
		[]string{"The component that was depended upon failed to deploy"},
		[]string{fmt.Sprintf("Resolve the failure reported for %s and deploy the design again", dependency)},
	)
}

func ErrResolveReference(err error) error {
	return errors.New(ErrResolveReferenceCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

func ErrYAMLUnmarshal(err error) error {
	return errors.New(ErrYAMLUnmarshalCode, errors.Alert, []string{"failed to parse YAML configuration"}, []string{err.Error()}, []string{"malformed YAML in component configuration"}, []string{"verify the YAML syntax in your design configuration"})
}
