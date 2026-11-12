package stages

import (
	"fmt"

	"github.com/meshery/meshery/server/models/pattern/planner"
	"github.com/meshery/meshkit/errors"
)

const (
	ErrResolveReferenceCode       = "meshery-server-1361"
	ErrYAMLUnmarshalCode          = "meshery-server-1436"
	ErrDependencyNotSatisfiedCode = "meshery-server-1480"
)

// ErrDependencyNotSatisfied reports a component that was not acted on because a
// component it declared a dependency on did not succeed.
//
// The text is split between its two audiences the same way the planner's
// dependency errors are: single literals for the short description and the
// remediation, which is what the published error reference can carry, and the
// formatted sentence for the long description, which names the components
// involved and so has nothing to publish until there is a real design.
func ErrDependencyNotSatisfied(component, dependency string, isDelete bool) error {
	explanation := fmt.Sprintf("Component %q was not deployed. It declares a dependency on %s, which did not deploy successfully, and Meshery deploys a component only once every component it depends on is in place.", component, dependency)
	if isDelete {
		explanation = fmt.Sprintf("Component %q was not undeployed. %s depends on it and did not undeploy successfully, and Meshery removes a component only once everything that depends on it has been removed.", component, dependency)
	}

	return errors.New(
		ErrDependencyNotSatisfiedCode,
		errors.Alert,
		[]string{"A component was withheld because a component it depends on did not succeed"},
		[]string{explanation},
		[]string{"A component that this one's order depends on did not succeed"},
		[]string{"Resolve the failure reported for that component, then try the design again"},
	)
}

func ErrResolveReference(err error) error {
	return errors.New(ErrResolveReferenceCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

// dependencyErrorCodes are the errors that describe a design's declared
// dependencies rather than a reference that could not be resolved.
var dependencyErrorCodes = map[string]struct{}{
	planner.ErrInvalidDependencyCode:    {},
	planner.ErrUnresolvedDependencyCode: {},
	planner.ErrAmbiguousDependencyCode:  {},
	planner.ErrUnknownPlanNodeCode:      {},
	planner.ErrCyclicDependencyCode:     {},
}

// resolveReferenceError reports a failure met while filling a design in.
//
// A dependency error is returned unchanged, because it describes the failure
// far better than "failed to resolve reference" does and re-wrapping would drop
// that code on the way out of the stage. Every other failure keeps the resolve
// reference code it has always carried - several of them wrap a MeshKit type
// cast error, and passing those through would silently re-code them too.
func resolveReferenceError(err error) error {
	if _, isDependency := dependencyErrorCodes[errors.GetCode(err)]; isDependency {
		return err
	}

	return ErrResolveReference(err)
}

func ErrYAMLUnmarshal(err error) error {
	return errors.New(ErrYAMLUnmarshalCode, errors.Alert, []string{"failed to parse YAML configuration"}, []string{err.Error()}, []string{"malformed YAML in component configuration"}, []string{"verify the YAML syntax in your design configuration"})
}
