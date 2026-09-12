package planner

import (
	"fmt"

	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
)

// Plan struct represents a node of an execution plan
type Plan struct {
	Data pattern.PatternFile
	*Graph
}

// IsFeasible returns true if the plan execution is feasible
func (p *Plan) IsFeasible() (bool, error) {
	hasCycle, err := p.DetectCycle()
	if err != nil {
		return false, err
	}

	return !hasCycle, nil
}

// Execute traverses the plan and calls the callback function
// on each of the node.
//
// A component that is never dispatched, because a component it declared a
// dependency on did not succeed, is reported to withheld instead.
func (p *Plan) Execute(cb VisitFn, withheld WithheldFn, log logger.Handler) error {
	parallelGraph, err := NewParallelProcessGraph(p.Graph)
	if err != nil {
		return err
	}

	parallelGraph.Traverse(cb, withheld, log)

	return nil
}

// DeclaredDependencies reads the names a component declares a dependency on.
//
// A design travels between Meshery Server and its clients as JSON, so the
// entries arrive as []string while the design is held in memory and as
// []interface{} once the design has been round-tripped over the wire.
//
// Every reader of "dependsOn" resolves it through this one helper, so that a
// declaration Meshery cannot read is reported the same way wherever it is met
// rather than being read as one shape in one place and another elsewhere.
func DeclaredDependencies(design string, c *component.ComponentDefinition) ([]string, error) {
	raw, ok := c.Metadata.AdditionalProperties["dependsOn"]
	if !ok || raw == nil {
		return nil, nil
	}

	switch dependsOn := raw.(type) {
	case []string:
		return dependsOn, nil
	case []interface{}:
		names := make([]string, 0, len(dependsOn))
		for _, dep := range dependsOn {
			name, isName := dep.(string)
			if !isName {
				return nil, ErrInvalidDependency(design, c.DisplayName, fmt.Sprintf("one of its entries is a %T", dep))
			}

			names = append(names, name)
		}

		return names, nil
	default:
		return nil, ErrInvalidDependency(design, c.DisplayName, fmt.Sprintf("found a %T", raw))
	}
}

// CreatePlan takes in the application components and creates a plan of execution for it.
//
// A component's "dependsOn" entries name other components of the same design by
// their name, while the plan is keyed by component id, so each declared
// dependency is resolved to the id of the component it names.
//
// annotations carries the non-semantic components that were set aside before
// planning. They are never deployed, so a dependency naming one is always
// satisfied and contributes no edge.
func CreatePlan(design pattern.PatternFile, annotations []*component.ComponentDefinition, invert bool) (*Plan, error) {
	g := NewGraph()

	idsByName := map[string][]string{}
	for _, component := range design.Components {
		g.AddNode(component.ID.String(), *component)
		idsByName[component.DisplayName] = append(idsByName[component.DisplayName], component.ID.String())
	}

	annotationNames := map[string]struct{}{}
	for _, annotation := range annotations {
		annotationNames[annotation.DisplayName] = struct{}{}
	}

	for _, component := range design.Components {
		dependsOn, err := DeclaredDependencies(design.Name, component)
		if err != nil {
			return nil, err
		}

		for _, dep := range dependsOn {
			ids := idsByName[dep]

			switch len(ids) {
			case 1:
			case 0:
				if _, isAnnotation := annotationNames[dep]; isAnnotation {
					continue
				}

				return nil, ErrUnresolvedDependency(design.Name, component.DisplayName, dep)
			default:
				return nil, ErrAmbiguousDependency(design.Name, component.DisplayName, dep, len(ids))
			}

			from := ids[0]
			to := component.ID.String()

			if invert {
				from, to = to, from
			}

			g.AddEdge(from, to)
		}
	}

	return &Plan{design, g}, nil
}
