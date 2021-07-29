package pattern

import "sync"

// ChainStageFunction is the type for function that will be invoked on each stage of the chain
type ChainStageFunction func(plan *Plan, err error, next ChainStageFunction) (*Plan, error)

// ChainStages type represents a slice of ChainStageFunction
type ChainStages []ChainStageFunction

// Chain allows to add any number of stages to be added to itself
// allowing "chaining" all of those functions
type Chain struct {
	stages ChainStages

	mu *sync.Mutex
}

// CreateChain returns a pointer to the chain object
func CreateChain() *Chain {
	return &Chain{
		stages: make(ChainStages, 0),
	}
}

// Add adds a function to the chain and returns a pointer to the Chain object
func (ch *Chain) Add(fn ChainStageFunction) *Chain {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	ch.stages = append(ch.stages, fn)

	return ch
}

// Consume takes in a plan and starts the chain of the functions
// when the the functions returns all the stages of the chain would
// be cleaned up
//
// Returns a pointer to the Chain object
func (ch *Chain) Consume(plan *Plan) *Chain {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	var err error
	for len(ch.stages) > 0 {
		var next ChainStageFunction
		if len(ch.stages) > 1 {
			next = ch.stages[1]
		}

		plan, err = ch.stages[0](plan, err, next)

		ch.stages = ch.stages[1:]
	}

	return ch
}
