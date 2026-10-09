package planner

import (
	"fmt"
	"sync"

	"github.com/meshery/meshkit/logger"
)

// ParallelProcessGraph provides methods for parallel processing of the graph
type ParallelProcessGraph struct {
	ParallelProcessGraphNodeMap map[string]*ParallelProcessGraphNode

	// dependents maps a node to the nodes that declared a dependency on it,
	// resolved once while the graph is built so that traversal never has to
	// look a node up by name again.
	dependents map[string][]*ParallelProcessGraphNode
	wg         sync.WaitGroup
}

// ParallelProcessGraphNode is the node for ParallelProcessGraph
// it contains the constructs for channel based communication
type ParallelProcessGraphNode struct {
	Name string
	Data *Node

	DepUpdateCh chan struct{}
	// DepCancleCh carries the name of the dependency that did not succeed.
	DepCancleCh chan string
	depLock     sync.Mutex

	DepCount int
}

// NewParallelProcessGraph creates a parallel processing graph
// from a simple graph returns a pointer to ParallelProcessGraph
func NewParallelProcessGraph(g *Graph) (*ParallelProcessGraph, error) {
	// Create a new parallel process graph
	pg := &ParallelProcessGraph{
		ParallelProcessGraphNodeMap: make(map[string]*ParallelProcessGraphNode),
		dependents:                  make(map[string][]*ParallelProcessGraphNode),
	}

	// Copy the nodes data
	for node, data := range g.Nodes {
		pg.ParallelProcessGraphNodeMap[node] = &ParallelProcessGraphNode{
			Name:        node,
			Data:        data,
			DepUpdateCh: make(chan struct{}, 1),
			DepCancleCh: make(chan string, 1),
			DepCount:    0,
		}
	}

	// Copy the edges data
	for node, adjacentNodes := range g.Edges {
		if _, known := pg.ParallelProcessGraphNodeMap[node]; !known {
			return nil, ErrUnknownPlanNode(node, "the execution plan's dependency list")
		}

		for _, aNode := range adjacentNodes {
			dependent, known := pg.ParallelProcessGraphNodeMap[aNode]
			if !known {
				return nil, ErrUnknownPlanNode(aNode, fmt.Sprintf("the dependency %s -> %s", node, aNode))
			}

			dependent.DepCount++
			pg.dependents[node] = append(pg.dependents[node], dependent)
		}
	}

	return pg, nil
}

// Traverse spins up the processes concurrently if it can
func (g *ParallelProcessGraph) Traverse(fn VisitFn, withheld WithheldFn, log logger.Handler) {
	// Spin up the processes concurrently
	for name, node := range g.ParallelProcessGraphNodeMap {
		g.wg.Add(1)
		go node.Process(g.dependents[name], &g.wg, fn, withheld, log)
	}

	g.wg.Wait()
}

// Process starts an internal loop which listens for the signals on the channels
// and operate accordingly
func (v *ParallelProcessGraphNode) Process(deps []*ParallelProcessGraphNode, wg *sync.WaitGroup, fn VisitFn, withheld WithheldFn, log logger.Handler) {
	defer wg.Done()

	log.Debug("started with:", v.Name, v.DepCount)

	depSuccessCount := 0
	depFailCount := 0
	failedDep := ""

	// Every dependency reports exactly once, either a success or a failure, so
	// the node is settled once as many signals have arrived as it has
	// dependencies.
	for depSuccessCount+depFailCount < v.DepCount {
		select {
		case name := <-v.DepCancleCh:
			depFailCount++

			if failedDep == "" {
				failedDep = name
			}

			log.Debug("RECEIVED DEP FAILURE:", v.Name)
		case <-v.DepUpdateCh:
			depSuccessCount++
		}
	}

	// A node whose dependency did not succeed is never visited, and propagates
	// the abort to the nodes that depend on it in turn.
	if depFailCount > 0 {
		if withheld != nil {
			withheld(v.Name, v.Data.Data, failedDep)
		}

		sendSignalToDeps(deps, false, v.Name, log)

		return
	}

	// The resources are deployed in correct order and dependsOn is respected but sometimes it has issues, for eg: CR depends on CRD, hence when deployment request to k8s to deploy CRD succeeds, we continue with the deployment of dependent CR.
	// The 200 response from k8s doesn’t guarantee that resource is available to use, it is just an indication that req is received and being worked on, therefore in ceratin cases, deployment failures are experienced and hence we need a mechanism to ensure that the dependent resource is actually deployed and ready to use before conitnuing.
	ok := fn(v.Name, v.Data.Data)
	log.Debug("RESPONSE GOT FOR", v.Name, ":", ok)

	// Send the appropriate signal
	sendSignalToDeps(deps, ok, v.Name, log)
}

// sendSignalToDeps sends signal on the channel
func sendSignalToDeps(deps []*ParallelProcessGraphNode, ok bool, from string, log logger.Handler) {
	for _, dep := range deps {
		sendSignalToDep(dep, ok, from, log)
	}
}

// sendSignalToDep sends signal either on the DepUpdateChannel or DepCancelChannel
// depending upon the "ok" parameter. "from" names the node the signal is sent on
// behalf of.
//
// sendSignalToDep will first acquire a lock on the channel operation and then will
// perform the acttion, hence this method SHOULD be used for sending the messages
// on the channels to avoid races
func sendSignalToDep(dep *ParallelProcessGraphNode, ok bool, from string, log logger.Handler) {
	// Lock the channel
	dep.depLock.Lock()
	defer dep.depLock.Unlock()

	if ok {
		log.Debug("Sending completed signal to:", dep.Name)
		dep.DepUpdateCh <- struct{}{}
		return
	}

	log.Debug("Sending abort signal to:", dep.Name)
	dep.DepCancleCh <- from
}
