package planner

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/meshery/meshkit/logger"
)

// DefaultConcurrencyLimit is the default maximum number of components
// that can execute fn(...) concurrently in ParallelProcessGraph.Traverse.
const DefaultConcurrencyLimit = 10

// ParallelProcessGraph provides methods for parallel processing of the graph
type ParallelProcessGraph struct {
	ParallelProcessGraphNodeMap map[string]*ParallelProcessGraphNode

	// dependents maps a node to the nodes that declared a dependency on it,
	// resolved once while the graph is built so that traversal never has to
	// look a node up by name again.
	dependents map[string][]*ParallelProcessGraphNode
	wg         sync.WaitGroup

	// ConcurrencyLimit bounds the maximum number of components executing fn concurrently.
	// If <= 0, DefaultConcurrencyLimit is used.
	ConcurrencyLimit int

	// MaxJitter is an optional upper bound for random jitter added before executing fn.
	// Useful in high-scale deployments to desynchronize simultaneous requests.
	MaxJitter time.Duration
}

// ParallelProcessGraphNode is the node for ParallelProcessGraph
// it contains the constructs for channel based communication
type ParallelProcessGraphNode struct {
	Name string
	Data *Node

	DepUpdateCh chan struct{}
	// DepCancleCh carries the name of the dependency that did not succeed.
	DepCancleCh chan string

	DepCount int

	// sem bounds concurrent execution of fn across nodes in the graph.
	sem chan struct{}

	// maxJitter specifies maximum duration for random jitter before calling fn.
	maxJitter time.Duration
}

// NewParallelProcessGraph creates a parallel processing graph
// from a simple graph returns a pointer to ParallelProcessGraph
func NewParallelProcessGraph(g *Graph) (*ParallelProcessGraph, error) {
	// Create a new parallel process graph
	pg := &ParallelProcessGraph{
		ParallelProcessGraphNodeMap: make(map[string]*ParallelProcessGraphNode),
		dependents:                  make(map[string][]*ParallelProcessGraphNode),
		ConcurrencyLimit:            DefaultConcurrencyLimit,
	}

	// Copy the nodes data
	for node, data := range g.Nodes {
		pg.ParallelProcessGraphNodeMap[node] = &ParallelProcessGraphNode{
			Name:     node,
			Data:     data,
			DepCount: 0,
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

	// Sizing both channels to max(1, node.DepCount) ensures the buffer has sufficient
	// capacity for all incoming signals. Given the exactly-once reporting invariant
	// for each dependency, the buffer never saturates, allowing sends to proceed
	// without waiting or lock contention.
	for _, node := range pg.ParallelProcessGraphNodeMap {
		bufSize := node.DepCount
		if bufSize < 1 {
			bufSize = 1
		}
		node.DepUpdateCh = make(chan struct{}, bufSize)
		node.DepCancleCh = make(chan string, bufSize)
	}

	return pg, nil
}

// Traverse spins up the processes concurrently if it can, bounding concurrent
// visitor execution using a semaphore to avoid thundering herd on downstream APIs.
// If ctx is cancelled before or during traversal, Traverse aborts remaining work
// and returns ctx.Err().
func (g *ParallelProcessGraph) Traverse(ctx context.Context, fn VisitFn, withheld WithheldFn, log logger.Handler) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := ctx.Err(); err != nil {
		return err
	}

	limit := g.ConcurrencyLimit
	if limit <= 0 {
		limit = DefaultConcurrencyLimit
	}
	sem := make(chan struct{}, limit)

	// Spin up the processes concurrently
	for name, node := range g.ParallelProcessGraphNodeMap {
		node.sem = sem
		node.maxJitter = g.MaxJitter
		g.wg.Add(1)
		go node.Process(ctx, g.dependents[name], &g.wg, fn, withheld, log)
	}

	waitCh := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(waitCh)
	}()

	select {
	case <-ctx.Done():
		<-waitCh
		return ctx.Err()
	case <-waitCh:
		return ctx.Err()
	}
}

// Process starts an internal loop which listens for the signals on the channels
// and operate accordingly
func (v *ParallelProcessGraphNode) Process(ctx context.Context, deps []*ParallelProcessGraphNode, wg *sync.WaitGroup, fn VisitFn, withheld WithheldFn, log logger.Handler) {
	signaled := make(map[*ParallelProcessGraphNode]bool, len(deps))
	sendSignals := func(ok bool) {
		for _, dep := range deps {
			if !signaled[dep] {
				signaled[dep] = true
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Errorf("panic sending signal to dependent %s: %v", dep.Name, r)
						}
					}()
					sendSignalToDep(dep, ok, v.Name, log)
				}()
			}
		}
	}

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("panic in parallel process node %s: %v", v.Name, r)
			sendSignals(false)
		}
		wg.Done()
	}()

	log.Debug("started with:", v.Name, v.DepCount)

	depSuccessCount := 0
	depFailCount := 0
	failedDep := ""

	// Every dependency reports exactly once, either a success or a failure, so
	// the node is settled once as many signals have arrived as it has
	// dependencies.
	for depSuccessCount+depFailCount < v.DepCount {
		select {
		case <-ctx.Done():
			log.Debug("aborted due to context cancellation:", v.Name)
			sendSignals(false)
			return
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

		sendSignals(false)
		return
	}

	// Check context cancellation before initiating deployment work
	if err := ctx.Err(); err != nil {
		sendSignals(false)
		return
	}

	// The resources are deployed in correct order and dependsOn is respected but sometimes it has issues, for eg: CR depends on CRD, hence when deployment request to k8s to deploy CRD succeeds, we continue with the deployment of dependent CR.
	// The 200 response from k8s doesn’t guarantee that resource is available to use, it is just an indication that req is received and being worked on, therefore in ceratin cases, deployment failures are experienced and hence we need a mechanism to ensure that the dependent resource is actually deployed and ready to use before conitnuing.
	// Bounded concurrency & jitter: Apply optional jitter first to desynchronize
	// bursts, then acquire a token from the semaphore right before calling fn so
	// that slots are held only for actual deployment work.
	ok := func() bool {
		if v.maxJitter > 0 {
			jitter := time.Duration(rand.Int63n(int64(v.maxJitter)))
			timer := time.NewTimer(jitter)
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return false
			case <-timer.C:
			}
		}

		if v.sem != nil {
			select {
			case <-ctx.Done():
				return false
			case v.sem <- struct{}{}:
				defer func() { <-v.sem }()
			}
		}

		if err := ctx.Err(); err != nil {
			return false
		}

		return fn(ctx, v.Name, v.Data.Data)
	}()
	log.Debug("RESPONSE GOT FOR", v.Name, ":", ok)

	// Send the appropriate signal
	sendSignals(ok)
}

// sendSignalToDep sends signal either on the DepUpdateChannel or DepCancelChannel
// depending upon the "ok" parameter. "from" names the node the signal is sent on
// behalf of.
//
// Because DepUpdateCh and DepCancleCh are buffered to the node's DepCount and each
// dependency reports exactly once, the buffer never saturates, allowing standard channel
// sends to complete without waiting and without requiring mutex locking.
func sendSignalToDep(dep *ParallelProcessGraphNode, ok bool, from string, log logger.Handler) {
	if ok {
		log.Debug("Sending completed signal to:", dep.Name)
		dep.DepUpdateCh <- struct{}{}
		return
	}

	log.Debug("Sending abort signal to:", dep.Name)
	dep.DepCancleCh <- from
}
