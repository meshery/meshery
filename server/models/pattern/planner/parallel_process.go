package planner

import (
	"sync"

	"github.com/layer5io/meshkit/logger"
)

// ParallelProcessGraph provides methods for parallel processing of the graph
type ParallelProcessGraph struct {
	Edges map[string][]string

	ParallelProcessGraphNodeMap map[string]*ParallelProcessGraphNode
	wg                          sync.WaitGroup
}

// ParallelProcessGraphNode is the node for ParallelProcessGraph
// it contains the constructs for channel based communication
type ParallelProcessGraphNode struct {
	Name string
	Data *Node

	DepUpdateCh chan struct{}
	DepCancleCh chan struct{}
	depLock     sync.Mutex

	DepCount int
}

// NewParallelProcessGraph creates a parallel processing graph
// from a simple graph returns a pointer to ParallelProcessGraph
func NewParallelProcessGraph(g *Graph) *ParallelProcessGraph {
	// Create a new parallel process graph
	pg := &ParallelProcessGraph{
		Edges:                       make(map[string][]string),
		ParallelProcessGraphNodeMap: make(map[string]*ParallelProcessGraphNode),
	}

	// Copy the nodes data
	for node, data := range g.Nodes {
		pg.ParallelProcessGraphNodeMap[node] = &ParallelProcessGraphNode{
			Name:        node,
			Data:        data,
			DepUpdateCh: make(chan struct{}, 1),
			DepCancleCh: make(chan struct{}, 1),
			DepCount:    0,
		}
	}

	// Copy the edges data
	for node, adjacentNodes := range g.Edges {
		pg.Edges[node] = adjacentNodes

		for _, aNode := range adjacentNodes {
			pg.ParallelProcessGraphNodeMap[aNode].DepCount++
		}
	}

	return pg
}

// Traverse spins up the processes concurrently if it can
func (g *ParallelProcessGraph) Traverse(fn VisitFn, log logger.Handler) {
	// Create map for node deps
	depsMap := make(map[string][]*ParallelProcessGraphNode)
	for node, edges := range g.Edges {
		for _, edge := range edges {
			depsMap[node] = append(depsMap[node], g.ParallelProcessGraphNodeMap[edge])
		}
	}

	// Spin up the processes concurrently
	for name, node := range g.ParallelProcessGraphNodeMap {
		g.wg.Add(1)
		go node.Process(depsMap[name], &g.wg, fn, log)
	}

	g.wg.Wait()
}

// Process starts an internal loop which listens for the signals on the channels
// and operate accordingly
func (v *ParallelProcessGraphNode) Process(deps []*ParallelProcessGraphNode, wg *sync.WaitGroup, fn VisitFn, log logger.Handler) {
	defer wg.Done()

	log.Debug("started with:", v.Name, v.DepCount)

	if v.DepCount == 0 {
		ok := fn(v.Name, v.Data.Data)
		log.Debug("RESPONSE GOT FOR", v.Name, ":", ok)

		// Send the appropriate signal
		sendSignalToDeps(deps, ok, log)

		return
	}

	depSuccessCount := 0
	depFailCount := 0

	for {
		select {
		case <-v.DepCancleCh:
			// Increment a dep failure
			depFailCount++

			log.Debug("REVCEIVED DEP FAILURE:", v.Name)

			// Now that we know that some of the deps have completed
			// successfully while some have failed we can proceed
			// to shut down this node
			if depFailCount+depSuccessCount == v.DepCount {
				// Send the appropriate signal
				// Propagate the abort
				sendSignalToDeps(deps, false, log)
				return
			}
		case <-v.DepUpdateCh:
			// Increment a dep success
			depSuccessCount++

			// If all of the deps were successful then
			// proceed to execute current node's function
			// and upon completion send the appropriate signal
			// to the dependent nodes

			// The resources are deployed in correct order and dependsOn is respected but sometimes it has issues, for eg: CR depends on CRD, hence when deployment request to k8s to deploy CRD succeeds, we continue with the deployment of dependent CR.
			// The 200 response from k8s doesn’t guarantee that resource is available to use, it is just an indication that req is received and being worked on, therefore in ceratin cases, deployment failures are experienced and hence we need a mechanism to ensure that the dependent resource is actually deployed and ready to use before conitnuing.
			if depSuccessCount == v.DepCount {
				log.Debug("Now deps 0 hence:", v.Name)
				ok := fn(v.Name, v.Data.Data)

				// Send the appropriate signal
				sendSignalToDeps(deps, ok, log)
				return
			}
		}
	}
}

// sendSignalToDeps sends signal on the channel
func sendSignalToDeps(deps []*ParallelProcessGraphNode, ok bool, log logger.Handler) {
	for _, dep := range deps {
		sendSignalToDep(dep, ok, log)
	}
}

// sendSignalToDep sends signal either on the DepUpdateChannel or DepCancelChannel
// depending upon the "ok" parameter.
//
// sendSignalToDep will first acquire a lock on the channel operation and then will
// perform the acttion, hence this method SHOULD be used for sending the messages
// on the channels to avoid races
func sendSignalToDep(dep *ParallelProcessGraphNode, ok bool, log logger.Handler) {
	// Lock the channel
	dep.depLock.Lock()
	defer dep.depLock.Unlock()

	if ok {
		log.Debug("Sending completed signal to:", dep.Name)
		dep.DepUpdateCh <- struct{}{}
		return
	}

	log.Debug("Sending abort signal to:", dep.Name)
	dep.DepCancleCh <- struct{}{}
}
