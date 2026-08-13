package handlers

import patterncore "github.com/meshery/meshery/server/models/pattern/core"

// patternExecutor provides an abstraction over pattern execution.
//
// The native executor preserves the current pattern-engine behavior while
// allowing alternative orchestration backends to be introduced without
// coupling HTTP handlers directly to a concrete implementation.
type patternExecutor interface {
	Execute(opts *patterncore.ProcessPatternOptions) (map[string]interface{}, error)
}

type nativePatternExecutor struct{}

func (nativePatternExecutor) Execute(opts *patterncore.ProcessPatternOptions) (map[string]interface{}, error) {
	return _processPattern(opts)
}

func executePattern(executor patternExecutor, opts *patterncore.ProcessPatternOptions) (map[string]interface{}, error) {
	if executor == nil {
		executor = nativePatternExecutor{}
	}

	return executor.Execute(opts)
}

var _ patternExecutor = nativePatternExecutor{}
