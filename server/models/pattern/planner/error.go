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

// Each error below divides its text between the two audiences that read it,
// because one string array cannot serve both. MeshKit renders an array by
// joining its elements, so two elements saying the same thing read as one
// sentence stated twice; the tool that publishes the error reference, in turn,
// reads only string literals and skips anything formatted at runtime.
//
// So the short description and the remediation are single literals, generic
// enough to stand alone in the published reference, and the long description is
// the formatted sentence naming the design and components involved. The long
// description is therefore blank in the published reference by design - it has
// nothing to say until there is a real design to name.

// ErrCyclicDependency reports a design whose components depend on each other in
// a loop, which leaves no order in which they can be acted on.
func ErrCyclicDependency(design string, isDelete bool) error {
	rule := "Meshery deploys a component only once every component it depends on has been deployed, so a loop of dependencies leaves no component that can be deployed first."
	if isDelete {
		rule = "Meshery removes a component only once everything that depends on it has been removed, so a loop of dependencies leaves no component that can be removed first."
	}

	return errors.New(
		ErrCyclicDependencyCode,
		errors.Alert,
		[]string{"The components of a design depend on each other in a loop"},
		[]string{fmt.Sprintf("The components of design %q depend on each other in a loop. %s", design, rule)},
		[]string{"Two or more components of the design declare a dependency on each other, either directly or through a chain of other components"},
		[]string{"Open the design in Meshery and remove one of the dependencies that closes the loop, then try again"},
	)
}

// ErrInvalidDependency reports a "dependsOn" declaration that is not a list of
// component names. reason describes the shape that was found instead.
func ErrInvalidDependency(design, component, reason string) error {
	return errors.New(
		ErrInvalidDependencyCode,
		errors.Alert,
		[]string{"A component declares its dependencies in a form Meshery cannot read"},
		[]string{fmt.Sprintf("Component %q of design %q declares its dependencies in a form Meshery cannot read: %q has to be a list of component names, but %s", component, design, "dependsOn", reason)},
		[]string{"The design was written by a tool that spells the dependsOn entry differently, or was hand-edited outside of Meshery"},
		[]string{"Edit the design so that the dependsOn entry on that component is a list of the names of other components in the same design"},
	)
}

// ErrUnresolvedDependency reports a "dependsOn" entry that names no component of
// the design.
func ErrUnresolvedDependency(design, component, dependency string) error {
	return errors.New(
		ErrUnresolvedDependencyCode,
		errors.Alert,
		[]string{"A component depends on something that is not a component of the design"},
		[]string{fmt.Sprintf("Component %q of design %q depends on %q, and no component of that design carries that name", component, design, dependency)},
		[]string{"The component that was depended upon has been renamed or removed while the dependency on it was left in place"},
		[]string{"Open the design in Meshery and either restore a component with the name that was depended upon, or remove that dependency"},
	)
}

// ErrAmbiguousDependency reports a "dependsOn" entry that names more than one
// component of the design.
func ErrAmbiguousDependency(design, component, dependency string, matches int) error {
	return errors.New(
		ErrAmbiguousDependencyCode,
		errors.Alert,
		[]string{"A component depends on a name that more than one component of the design carries"},
		[]string{fmt.Sprintf("Component %q of design %q depends on %q, and %d components of that design carry that name, so it does not identify which one to wait for", component, design, dependency, matches)},
		[]string{"Component names are not required to be unique within a design, and more than one component was given the same name"},
		[]string{"Open the design in Meshery and rename those components so that the one being depended upon is uniquely named"},
	)
}

// ErrUnknownPlanNode reports an execution plan that refers to a component which
// is not part of that plan.
func ErrUnknownPlanNode(missing, referencedBy string) error {
	return errors.New(
		ErrUnknownPlanNodeCode,
		errors.Alert,
		[]string{"The execution plan refers to a component that is not part of it"},
		[]string{fmt.Sprintf("%s refers to %s, but no component with that identifier is being acted on", referencedBy, missing)},
		[]string{"The design is internally inconsistent: it declares a dependency on a component that it does not contain"},
		[]string{"Open the design in Meshery, remove or repoint the dangling dependency, and try again"},
	)
}
