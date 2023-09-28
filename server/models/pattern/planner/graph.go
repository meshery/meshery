package planner

import (
	"sync"
	"sync/atomic"

	"github.com/layer5io/meshkit/utils/patterns"
)

// Graph represents the graph data structure
type Graph struct {
	Nodes map[string]*Node
	Edges map[string][]string

	sync.RWMutex
}

// Node is a graph node
type Node struct {
	Data patterns.Service
}

// NewGraph creates a new instance of the graph and returns a pointer to it
func NewGraph() *Graph {
	return &Graph{
		Nodes: make(map[string]*Node),
		Edges: make(map[string][]string),
	}
}

// VisitFn is the function definition for the visitor function
type VisitFn func(name string, node patterns.Service) bool

// AddNode adds a node to the graph
func (g *Graph) AddNode(name string, data patterns.Service) *Graph {
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
	return !g.topologicalSort(func(_ string, _ patterns.Service) bool { return true })
}

// Traverse traverses the graph in topological sorted order
// and executes the visit function on each visited node
func (g *Graph) Traverse(fn VisitFn) {
	g.topologicalSort(fn)
}

func (g *Graph) topologicalSort(fn VisitFn) bool {
	g.RLock()
	defer g.RUnlock()

	indegree := map[string]int{}
	for node := range g.Nodes {
		indegree[node] = 0
	}

	g.Visit(func(name string, _ patterns.Service) bool {
		indegree[name]++
		return true
	})

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

		if ok := fn(v, g.Nodes[v].Data); !ok {
			return false
		}

		atomic.AddInt64(&vertexCount, 1)

		for _, node := range g.Edges[v] {
			indegree[node]--

			if indegree[node] == 0 {
				queue.Enqueue(node)
			}
		}
	}

	return vertexCount == int64(g.Order())
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
