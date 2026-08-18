package stages

import (
	"fmt"
	"strings"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrResolveReferenceCode       = "meshery-server-1361"
	ErrYAMLUnmarshalCode          = "meshery-server-1436"
	ErrDependencyNotSatisfiedCode = "meshery-server-1480"
)

// ErrDependencyNotSatisfied reports a component that was not acted on because a
// component it declared a dependency on did not succeed.
func ErrDependencyNotSatisfied(component, dependency string, isDelete bool) error {
	verb := "deploy"
	explanation := fmt.Sprintf("%q declares a dependency on %s, which did not deploy successfully. Meshery deploys a component only once every component it depends on is in place, so %q was withheld.", component, dependency, component)
	cause := "The component that was depended upon failed to deploy"

	if isDelete {
		verb = "undeploy"
		explanation = fmt.Sprintf("%s depends on %q and did not undeploy successfully. Meshery removes a component only once everything that depends on it has been removed, so %q was withheld.", dependency, component, component)
		cause = "Something that depends on this component could not be removed"
	}

	return errors.New(
		ErrDependencyNotSatisfiedCode,
		errors.Alert,
		[]string{
			"A component was withheld because a component it depends on did not succeed",
			fmt.Sprintf("Component %q was not %sed", component, verb),
		},
		[]string{
			"Meshery acts on a component only once the components it depends on have been dealt with, so this component was withheld",
			explanation,
		},
		[]string{
			"A component this one depends on did not succeed",
			cause,
		},
		[]string{
			"Resolve the failure reported for the dependency, then try the design again",
			fmt.Sprintf("Resolve the failure reported for %s and %s the design again", dependency, verb),
		},
	)
}

func ErrResolveReference(err error) error {
	return errors.New(ErrResolveReferenceCode, errors.Alert, []string{}, []string{err.Error()}, []string{}, []string{})
}

// resolveReferenceError reports a failure met while filling a design in.
//
// An error that already carries a MeshKit code is returned unchanged: it
// describes the failure more precisely than "failed to resolve reference" can,
// and re-wrapping would drop that code before it reached the API boundary.
func resolveReferenceError(err error) error {
	if errors.GetCode(err) != strings.Join(errors.NoneString, "") {
		return err
	}

	return ErrResolveReference(err)
}

func ErrYAMLUnmarshal(err error) error {
	return errors.New(ErrYAMLUnmarshalCode, errors.Alert, []string{"failed to parse YAML configuration"}, []string{err.Error()}, []string{"malformed YAML in component configuration"}, []string{"verify the YAML syntax in your design configuration"})
}
