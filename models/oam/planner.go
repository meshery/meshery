package oam

// Plan struct represents a node of an execution plan
type Plan struct {
	Data Pattern
	*Graph
}

// IsFeasible returns true if the plan execution is feasible
func (p *Plan) IsFeasible() bool {
	return !p.DetectCycle()
}

// Execute traverses the plan and calls the callback function
// on each of the node
func (p *Plan) Execute(cb func(string, Service) bool) error {
	parallelGraph := NewParallelProcessGraph(p.Graph)
	parallelGraph.Traverse(cb)
	return nil
}

// CreatePlan takes in the application components and creates a plan of execution for it
func CreatePlan(pattern Pattern, policies [][2]string) (*Plan, error) {
	g := NewGraph()

	for name, svc := range pattern.Services {
		g.AddNode(name, *svc)
	}

	for name, svc := range pattern.Services {
		for _, deps := range svc.DependsOn {
			g.AddEdge(deps, name)
		}

		// Inject "policies"
		for _, policy := range policies {
			g.AddEdge(policy[0], policy[1])
		}
	}

	return &Plan{pattern, g}, nil
}
