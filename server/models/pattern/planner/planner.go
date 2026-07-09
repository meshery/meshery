package planner

import (
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta2/component"
	pattern "github.com/meshery/schemas/models/v1beta3/design"
	"github.com/pkg/errors"
)

// Plan struct represents a node of an execution plan
type Plan struct {
	Data pattern.PatternFile
	*Graph
}

// IsFeasible returns true if the plan execution is feasible
func (p *Plan) IsFeasible() bool {
	return !p.DetectCycle()
}

// Execute traverses the plan and calls the callback function
// on each of the node
func (p *Plan) Execute(cb func(string, component.ComponentDefinition) bool, log logger.Handler) error {
	parallelGraph := NewParallelProcessGraph(p.Graph)
	parallelGraph.Traverse(cb, log)
	return nil
}

// CreatePlan takes in the application components and creates a plan of execution for it
func CreatePlan(pattern pattern.PatternFile, invert bool) (*Plan, error) {
	g := NewGraph()

	for _, component := range pattern.Components {
		g.AddNode(component.ID.String(), *component)
	}

	// "dependsOn" entries reference sibling components by display name while
	// graph nodes are keyed by component ID; resolve the names so that edges
	// land on real nodes.
	idByDisplayName := make(map[string]string, len(pattern.Components))
	for _, component := range pattern.Components {
		if _, ok := idByDisplayName[string(component.DisplayName)]; !ok {
			idByDisplayName[string(component.DisplayName)] = component.ID.String()
		}
	}

	for _, component := range pattern.Components {
		_dependsOn, ok := component.Metadata.AdditionalProperties["dependsOn"]
		if !ok {
			continue
		}

		dependsOn, err := utils.Cast[[]string](_dependsOn)
		if err != nil {
			err = errors.Wrapf(err, "Failed to cast 'dependsOn' to []string for component %s", component.DisplayName)
			return nil, err
		}
		for _, dep := range dependsOn {
			depID, ok := idByDisplayName[dep]
			if !ok {
				if _, isID := g.Nodes[dep]; isID {
					// The entry already carries a component ID.
					depID = dep
				} else {
					return nil, errors.Errorf("component %s dependsOn %q, which does not match any component in the design", component.DisplayName, dep)
				}
			}

			from := depID
			to := component.ID.String()

			if invert {
				from, to = to, from
			}

			g.AddEdge(from, to)
		}
	}

	return &Plan{pattern, g}, nil
}
