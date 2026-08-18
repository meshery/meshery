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

// Each of these errors leads every string array with a static literal, because
// the error reference published from this file is extracted statically and a
// formatted string is not carried into it. The formatted detail follows as a
// further element, where it reaches the user at runtime.

// ErrCyclicDependency reports a design whose components depend on each other in
// a loop, which leaves no order in which they can be acted on.
func ErrCyclicDependency(design string, isDelete bool) error {
	rule := "Meshery deploys a component only once every component it depends on has been deployed, so a loop of dependencies leaves no component that can be deployed first."
	verb := "deploy"

	if isDelete {
		rule = "Meshery removes a component only once everything that depends on it has been removed, so a loop of dependencies leaves no component that can be removed first."
		verb = "undeploy"
	}

	return errors.New(
		ErrCyclicDependencyCode,
		errors.Alert,
		[]string{
			"The components of a design depend on each other in a loop",
			fmt.Sprintf("The components of design %q depend on each other in a loop", design),
		},
		[]string{
			"A loop of dependencies leaves no order in which the components of the design can be acted on",
			rule,
		},
		[]string{"Two or more components of the design declare a dependency on each other, either directly or through a chain of other components"},
		[]string{
			"Open the design in Meshery and remove one of the dependencies that closes the loop, then try again",
			fmt.Sprintf("Open design %q in Meshery and remove one of the dependencies that closes the loop, then %s again", design, verb),
		},
	)
}

// ErrInvalidDependency reports a "dependsOn" declaration that is not a list of
// component names. reason describes the shape that was found instead.
func ErrInvalidDependency(design, component, reason string) error {
	return errors.New(
		ErrInvalidDependencyCode,
		errors.Alert,
		[]string{
			"A component declares its dependencies in a form Meshery cannot read",
			fmt.Sprintf("Component %q of design %q declares its dependencies in a form Meshery cannot read", component, design),
		},
		[]string{
			"\"dependsOn\" has to be a list of the names of other components in the same design",
			fmt.Sprintf("\"dependsOn\" has to be a list of component names: %s", reason),
		},
		[]string{"The design was written by a tool that spells \"dependsOn\" differently, or was hand-edited outside of Meshery"},
		[]string{
			"Edit the design so that \"dependsOn\" on that component is a list of the names of other components in the same design",
			fmt.Sprintf("Edit design %q so that \"dependsOn\" on component %q is a list of the names of other components in the same design", design, component),
		},
	)
}

// ErrUnresolvedDependency reports a "dependsOn" entry that names no component of
// the design.
func ErrUnresolvedDependency(design, component, dependency string) error {
	return errors.New(
		ErrUnresolvedDependencyCode,
		errors.Alert,
		[]string{
			"A component depends on something that is not a component of the design",
			fmt.Sprintf("Component %q of design %q depends on %q, which is not a component of that design", component, design, dependency),
		},
		[]string{
			"Dependencies name other components of the same design by their name, and no component of the design carries the name that was declared",
			fmt.Sprintf("Dependencies name other components of the same design by their name. No component of design %q is named %q.", design, dependency),
		},
		[]string{"The component that was depended upon has been renamed or removed while the dependency on it was left in place"},
		[]string{
			"Open the design in Meshery and either restore a component with the name that was depended upon, or remove that dependency",
			fmt.Sprintf("Open design %q in Meshery and either restore a component named %q or remove that dependency from component %q", design, dependency, component),
		},
	)
}

// ErrAmbiguousDependency reports a "dependsOn" entry that names more than one
// component of the design.
func ErrAmbiguousDependency(design, component, dependency string, matches int) error {
	return errors.New(
		ErrAmbiguousDependencyCode,
		errors.Alert,
		[]string{
			"A component depends on a name that more than one component of the design carries",
			fmt.Sprintf("Component %q of design %q depends on %q, and %d components of that design carry that name", component, design, dependency, matches),
		},
		[]string{
			"Dependencies name other components of the same design by their name, so a name shared by several components does not identify which one to wait for",
			fmt.Sprintf("Dependencies name other components of the same design by their name, so a name shared by %d components does not identify which one to wait for.", matches),
		},
		[]string{"Component names are not required to be unique within a design, and more than one component was given the same name"},
		[]string{
			"Open the design in Meshery and rename those components so that the one being depended upon is uniquely named",
			fmt.Sprintf("Open design %q in Meshery and rename the components named %q so that the one component %q depends on is uniquely named", design, dependency, component),
		},
	)
}

// ErrUnknownPlanNode reports an execution plan that refers to a component which
// is not part of that plan.
func ErrUnknownPlanNode(missing, referencedBy string) error {
	return errors.New(
		ErrUnknownPlanNodeCode,
		errors.Alert,
		[]string{
			"The execution plan refers to a component that is not part of it",
			fmt.Sprintf("The execution plan for this deployment refers to a component (%s) that is not part of it", missing),
		},
		[]string{
			"A dependency in the plan points at a component that is not being acted on",
			fmt.Sprintf("%s refers to %s, but no component with that identifier is being deployed.", referencedBy, missing),
		},
		[]string{"The design is internally inconsistent: it declares a dependency on a component that it does not contain"},
		[]string{"Open the design in Meshery, remove or repoint the dangling dependency, and try again"},
	)
}
