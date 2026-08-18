package planner

import (
	"fmt"

	"github.com/meshery/meshkit/errors"
)

const (
	ErrInvalidDependencyCode    = "meshery-server-1475"
	ErrUnresolvedDependencyCode = "meshery-server-1476"
	ErrAmbiguousDependencyCode  = "meshery-server-1477"
	ErrUnknownPlanNodeCode      = "meshery-server-1478"
	ErrCyclicDependencyCode     = "meshery-server-1479"
)

// ErrCyclicDependency reports a design whose components depend on each other in
// a loop, which leaves no order in which they can be deployed.
func ErrCyclicDependency(design string) error {
	return errors.New(
		ErrCyclicDependencyCode,
		errors.Alert,
		[]string{fmt.Sprintf("The components of design %q depend on each other in a loop", design)},
		[]string{"Meshery deploys a component only once every component it depends on has been deployed. A loop of dependencies leaves no component that can be deployed first, so there is no order in which the design can be deployed."},
		[]string{"Two or more components of the design declare a dependency on each other, either directly or through a chain of other components"},
		[]string{fmt.Sprintf("Open design %q in Meshery and remove one of the dependencies that closes the loop, then deploy again", design)},
	)
}

// ErrInvalidDependency reports a "dependsOn" declaration that is not a list of
// component names. reason describes the shape that was found instead.
func ErrInvalidDependency(design, component, reason string) error {
	return errors.New(
		ErrInvalidDependencyCode,
		errors.Alert,
		[]string{fmt.Sprintf("Component %q of design %q declares its dependencies in a form Meshery cannot read", component, design)},
		[]string{fmt.Sprintf("\"dependsOn\" has to be a list of component names: %s", reason)},
		[]string{"The design was written by a tool that spells \"dependsOn\" differently, or was hand-edited outside of Meshery"},
		[]string{fmt.Sprintf("Edit design %q so that \"dependsOn\" on component %q is a list of the names of other components in the same design", design, component)},
	)
}

// ErrUnresolvedDependency reports a "dependsOn" entry that names no component of
// the design.
func ErrUnresolvedDependency(design, component, dependency string) error {
	return errors.New(
		ErrUnresolvedDependencyCode,
		errors.Alert,
		[]string{fmt.Sprintf("Component %q of design %q depends on %q, which is not a component of that design", component, design, dependency)},
		[]string{fmt.Sprintf("Dependencies name other components of the same design by their name. No component of design %q is named %q.", design, dependency)},
		[]string{"The component that was depended upon has been renamed or removed while the dependency on it was left in place"},
		[]string{fmt.Sprintf("Open design %q in Meshery and either restore a component named %q or remove that dependency from component %q", design, dependency, component)},
	)
}

// ErrAmbiguousDependency reports a "dependsOn" entry that names more than one
// component of the design.
func ErrAmbiguousDependency(design, component, dependency string, matches int) error {
	return errors.New(
		ErrAmbiguousDependencyCode,
		errors.Alert,
		[]string{fmt.Sprintf("Component %q of design %q depends on %q, and %d components of that design carry that name", component, design, dependency, matches)},
		[]string{fmt.Sprintf("Dependencies name other components of the same design by their name, so a name shared by %d components does not identify which one to wait for.", matches)},
		[]string{"Component names are not required to be unique within a design, and more than one component was given the same name"},
		[]string{fmt.Sprintf("Open design %q in Meshery and rename the components named %q so that the one component %q depends on is uniquely named", design, dependency, component)},
	)
}

// ErrUnknownPlanNode reports an execution plan that refers to a component which
// is not part of that plan.
func ErrUnknownPlanNode(missing, referencedBy string) error {
	return errors.New(
		ErrUnknownPlanNodeCode,
		errors.Alert,
		[]string{fmt.Sprintf("The execution plan for this deployment refers to a component (%s) that is not part of it", missing)},
		[]string{fmt.Sprintf("%s refers to %s, but no component with that identifier is being deployed.", referencedBy, missing)},
		[]string{"The design is internally inconsistent: it declares a dependency on a component that it does not contain"},
		[]string{"Open the design in Meshery, remove or repoint the dangling dependency, and deploy again"},
	)
}
