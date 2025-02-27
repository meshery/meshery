package planner

import (
	"github.com/layer5io/meshkit/logger"
	"github.com/layer5io/meshkit/utils"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
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
		g.AddNode(component.Id.String(), *component)
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
			from := dep
			to := component.DisplayName

			if invert {
				from = component.DisplayName
				to = dep
			}

			g.AddEdge(from, to)
		}
	}

	return &Plan{pattern, g}, nil
}
