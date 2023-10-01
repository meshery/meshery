package planner

import "github.com/layer5io/meshkit/utils/patterns"

// Plan struct represents a node of an execution plan
type Plan struct {
	Data patterns.Pattern
	*Graph
}

// IsFeasible returns true if the plan execution is feasible
func (p *Plan) IsFeasible() bool {
	return !p.DetectCycle()
}

// Execute traverses the plan and calls the callback function
// on each of the node
func (p *Plan) Execute(cb func(string, patterns.Service) bool) error {
	parallelGraph := NewParallelProcessGraph(p.Graph)
	parallelGraph.Traverse(cb)
	return nil
}

// CreatePlan takes in the application components and creates a plan of execution for it
func CreatePlan(pattern patterns.Pattern, invert bool) (*Plan, error) {
	g := NewGraph()

	for name, svc := range pattern.Services {
		g.AddNode(name, *svc)
	}

	for name, svc := range pattern.Services {
		for _, dep := range svc.DependsOn {
			from := dep
			to := name

			if invert {
				from = name
				to = dep
			}

			g.AddEdge(from, to)
		}
	}

	return &Plan{pattern, g}, nil
}
