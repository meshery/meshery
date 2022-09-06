package stages

import (
	"sync"

	"github.com/layer5io/meshery/server/models/pattern/core"
)

// Data is the struct that will be passed on each stage
type Data struct {
	Pattern                        *core.Pattern
	PatternSvcWorkloadCapabilities map[string]core.WorkloadCapability
	PatternSvcTraitCapabilities    map[string][]core.TraitCapability

	// Other is for passing metadata across different stages
	Lock  sync.Mutex
	Other map[string]interface{}
}

// ChainStageFunction is the type for function that will be invoked on each stage of the chain
type ChainStageFunction func(data *Data, err error, next ChainStageNextFunction)

type ChainStageNextFunction func(data *Data, err error)

// ChainStages type represents a slice of ChainStageFunction
type ChainStages []ChainStageNextFunction

// Chain allows to add any number of stages to be added to itself
// allowing "chaining" all of those functions
type Chain struct {
	stages ChainStages
	nexts  ChainStages

	mu sync.Mutex
}

// CreateChain returns a pointer to the chain object
func CreateChain() *Chain {
	return &Chain{
		stages: make(ChainStages, 0),
		nexts:  make(ChainStages, 0),
	}
}

// Add adds a function to the chain and returns a pointer to the Chain object
func (ch *Chain) Add(fn ChainStageFunction) *Chain {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	// Add the next function for "fn"
	ch.nexts = append(ch.nexts, nil)

	nextIdxStageFn := len(ch.nexts) - 1

	// Create the stage function
	stageFn := func(data *Data, err error) {
		fn(data, err, ch.nexts[nextIdxStageFn])
	}

	// Modify next function of previous stage to point
	// to the newly added "fn"
	if nextIdxStageFn > 0 {
		ch.nexts[nextIdxStageFn-1] = stageFn
	}

	// Add the stageFn to stages
	ch.stages = append(ch.stages, stageFn)

	return ch
}

// Process takes in a plan and starts the chain of the functions
//
// Returns a pointer to the Chain object
func (ch *Chain) Process(data *Data) *Chain {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	if len(ch.stages) > 0 {
		ch.stages[0](data, nil)
	}

	return ch
}

// Clear clears the chain and returns a pointer to the chain object
func (ch *Chain) Clear() *Chain {
	ch.stages = []ChainStageNextFunction{}
	ch.nexts = []ChainStageNextFunction{}

	return ch
}
