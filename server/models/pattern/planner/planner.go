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
	idsByDisplayName := make(map[string][]string, len(pattern.Components))
	for _, component := range pattern.Components {
		name := string(component.DisplayName)
		idsByDisplayName[name] = append(idsByDisplayName[name], component.ID.String())
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
			var depID string
			switch ids := idsByDisplayName[dep]; len(ids) {
			case 1:
				depID = ids[0]
			case 0:
				if _, isID := g.Nodes[dep]; !isID {
					return nil, errors.Errorf("component %s dependsOn %q, which does not match any component in the design", component.DisplayName, dep)
				}
				// The entry already carries a component ID.
				depID = dep
			default:
				return nil, errors.Errorf("component %s dependsOn %q, which matches %d components in the design; give them distinct display names or reference the component ID", component.DisplayName, dep, len(ids))
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
