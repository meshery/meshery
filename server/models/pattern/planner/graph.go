package planner

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/meshery/schemas/models/v1beta2/component"
)

// Graph represents the graph data structure
type Graph struct {
	Nodes map[string]*Node
	Edges map[string][]string

	sync.RWMutex
}

// Node is a graph node
type Node struct {
	Data component.ComponentDefinition
}

// NewGraph creates a new instance of the graph and returns a pointer to it
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}
}

// VisitFn is the function definition for the visitor function
type VisitFn func(name string, node component.ComponentDefinition) bool

// WithheldFn is the function definition invoked for a node that is never
// visited because a node it depends on did not succeed. failedDependency names
// that node.
type WithheldFn func(name string, node component.ComponentDefinition, failedDependency string)

// AddNode adds a node to the graph
func (g *Graph) AddNode(name string, data component.ComponentDefinition) *Graph {
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
func (g *Graph) DetectCycle() (bool, error) {
	sorted, err := g.topologicalSort(func(_ string, _ component.ComponentDefinition) bool { return true })
	if err != nil {
		return false, err
	}

	return !sorted, nil
}

// Traverse traverses the graph in topological sorted order
// and executes the visit function on each visited node
func (g *Graph) Traverse(fn VisitFn) error {
	_, err := g.topologicalSort(fn)
	return err
}

func (g *Graph) topologicalSort(fn VisitFn) (bool, error) {
	g.RLock()
	defer g.RUnlock()

	indegree := map[string]int{}
	for node := range g.Nodes {
		indegree[node] = 0
	}

	if err := g.Visit(func(name string, _ component.ComponentDefinition) bool {
		indegree[name]++
		return true
	}); err != nil {
		return false, err
	}

	// Hold all the vertices with 0 degree
	var queue Queue
	for v, degree := range indegree {
		if degree == 0 {
			queue.Enqueue(v)
		}
	}

	var vertexCount int64
	for queue.Length() > 0 {
		v := queue.Dequeue()

		node, known := g.Nodes[v]
		if !known {
			return false, ErrUnknownPlanNode(v, "the execution plan")
		}

		if ok := fn(v, node.Data); !ok {
			return false, nil
		}

		atomic.AddInt64(&vertexCount, 1)

		for _, node := range g.Edges[v] {
			indegree[node]--

			if indegree[node] == 0 {
				queue.Enqueue(node)
			}
		}
	}

	return vertexCount == int64(g.Order()), nil
}

// Order returns the count for number of edges
func (g *Graph) Order() int {
	return len(g.Nodes)
}

// Visit visits each node in the graph but does not keep
// track of the pre-visited nodes
func (g *Graph) Visit(fn VisitFn) error {
	for node := range g.Nodes {
		for _, edgeNode := range g.Edges[node] {
			target, known := g.Nodes[edgeNode]
			if !known {
				return ErrUnknownPlanNode(edgeNode, fmt.Sprintf("the dependency %s -> %s", node, edgeNode))
			}

			fn(edgeNode, target.Data)
		}
	}

	return nil
}
