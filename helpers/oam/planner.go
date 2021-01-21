package oam

import (
	"fmt"
	"sync"
	"time"

	"github.com/layer5io/meshery/helpers/oam/core/v1alpha1"
)

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
func (p *Plan) Execute(cb func(v1alpha1.Component) bool) error {
	comps, err := p.Data.ToApplicationComponents()
	if err != nil {
		return err
	}

	p.Traverse(func(node string, svc Service) bool {
		time.Sleep(1 * time.Second)

		for _, comp := range comps {
			if comp.Name == node {
				return cb(comp)
			}
		}

		return true
	})

	return nil
}

// CreatePlan takes in the application components and creates a plan of execution for it
func CreatePlan(pattern Pattern, policies [][2]string) (*Plan, error) {
	g := NewGraph()

	for name, svc := range pattern.Services {
		g.AddNode(name, svc)

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

// Graph represents the graph data structure
type Graph struct {
	Nodes map[string]*Node
	Edges map[string][]string

	sync.RWMutex
}

// Node is a graph node
type Node struct {
	Data Service
}

// NewGraph creates a new instance of the graph and returns a pointer to it
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}
}

// VisitFn is the function defintion for the visitor function
type VisitFn func(name string, node Service) bool

// AddNode adds a node to the graph
func (g *Graph) AddNode(name string, data Service) *Graph {
	g.Lock()
	defer g.Unlock()

	g.Nodes[name] = &Node{
		Data: data,
	}

	return g
}

// AddEdge adds edge from source to destination nodes
//
// These edges are unidirectional, hence for bidirectional
// edges this method needs to be called twice
func (g *Graph) AddEdge(src, dest string) *Graph {
	g.Lock()
	defer g.Unlock()

	if g.Edges[src] == nil {
		g.Edges[src] = make([]string, 0)
	}

	// Add iff it doesn't already exists
	exists := false
	for _, edge := range g.Edges[src] {
		if dest == edge {
			exists = true
			break
		}
	}

	if !exists {
		g.Edges[src] = append(g.Edges[src], dest)
	}

	return g
}

// DetectCycle will return true if there is a cycle
// in the graph
func (g *Graph) DetectCycle() bool {
	return !g.topologicalSort(func(_ string, _ Service) bool { return true })
}

// Traverse traverses the graph in topological sorted order
// and executes the visit function on each visited node
func (g *Graph) Traverse(fn VisitFn) {
	g.topologicalSort(fn)
}

func (g *Graph) topologicalSort(fn VisitFn) bool {
	g.RLock()
	defer g.RUnlock()

	indegree := make(map[string]int, g.Order())

	for node := range g.Nodes {
		indegree[node] = 0
	}

	g.Visit(func(name string, _ Service) bool {
		indegree[name]++
		return true
	})
	for k, v := range indegree {
		fmt.Println(k, ": ", v)
	}

	// Hold all the vertices with 0 degree
	var queue []string
	for v, degree := range indegree {
		if degree == 0 {
			queue = append(queue, v)
		}
	}

	vertexCount := 0
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]

		if ok := fn(v, g.Nodes[v].Data); !ok {
			return false
		}

		vertexCount++

		for _, node := range g.Edges[v] {
			indegree[node]--

			if indegree[node] == 0 {
				queue = append(queue, node)
			}
		}
	}

	return vertexCount == g.Order()
}

// Order returns the count for number of edges
func (g *Graph) Order() int {
	return len(g.Nodes)
}

// Visit visits each node in the graph but does not keep
// track of the pre-visited nodes
func (g *Graph) Visit(fn VisitFn) {
	for node := range g.Nodes {
		for _, edgeNode := range g.Edges[node] {
			fn(edgeNode, g.Nodes[edgeNode].Data)
		}
	}
}
