package handlers

import (
	"testing"

	patterncore "github.com/meshery/meshery/server/models/pattern/core"
)

type fakePatternExecutor struct {
	called   bool
	opts     *patterncore.ProcessPatternOptions
	response map[string]interface{}
	err      error
}

func (f *fakePatternExecutor) Execute(opts *patterncore.ProcessPatternOptions) (map[string]interface{}, error) {
	f.called = true
	f.opts = opts

	return f.response, f.err
}

func TestExecutePatternUsesConfiguredExecutor(t *testing.T) {
	opts := &patterncore.ProcessPatternOptions{}

	executor := &fakePatternExecutor{
		response: map[string]interface{}{
			"executor": "fake",
		},
	}

	response, err := executePattern(executor, opts)
	if err != nil {
		t.Fatalf("executePattern() returned unexpected error: %v", err)
	}

	if !executor.called {
		t.Fatal("expected configured pattern executor to be called")
	}

	if executor.opts != opts {
		t.Fatal("expected ProcessPatternOptions to be passed to configured executor")
	}

	if got := response["executor"]; got != "fake" {
		t.Fatalf("expected response from configured executor, got %v", got)
	}
}
