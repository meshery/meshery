//go:build js && wasm

// Wasm entry point for the native Go relationship policy engine.
// Builds with: GOOS=js GOARCH=wasm go build -o policy_engine.wasm .
package main

import (
	"encoding/json"
	"syscall/js"

	gopolicies "github.com/meshery/meshery/server/policies"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// noopLogger satisfies policies.Logger with no-ops.
type noopLogger struct{}

func (noopLogger) Info(...any)          {}
func (noopLogger) Warnf(string, ...any) {}

// Cached after initEngine; evaluateDesign reuses both per call so the
// registry only crosses the JS/wasm boundary once per session.
var (
	cachedEngine *gopolicies.GoEngine
	cachedRels   []*relationship.RelationshipDefinition
)

// initEngine is exposed to JS as globalThis.initEngine(relsJSON).
// Parses + caches the registered relationships and constructs a single
// GoEngine that subsequent evaluateDesign calls reuse. Returns "" on success
// or {"error": "..."} on failure.
func initEngine(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return errResult("initEngine(relsJSON): expected 1 argument")
	}
	var rels []*relationship.RelationshipDefinition
	if err := json.Unmarshal([]byte(args[0].String()), &rels); err != nil {
		return errResult("decode relationships: " + err.Error())
	}
	cachedRels = rels
	cachedEngine = gopolicies.NewGoEngine(noopLogger{})
	return ""
}

// evaluateDesign is exposed to JS as globalThis.evaluateDesign(designJSON).
// Returns a JSON string of pattern.EvaluationResponse, or {"error": "..."} on failure.
//
// Caller must invoke initEngine(relsJSON) once first; the registered
// relationships are reused across calls.
func evaluateDesign(_ js.Value, args []js.Value) any {
	if cachedEngine == nil {
		return errResult("engine not initialized; call initEngine(relsJSON) first")
	}
	if len(args) != 1 {
		return errResult("evaluateDesign(designJSON): expected 1 argument")
	}

	var design pattern.PatternFile
	if err := json.Unmarshal([]byte(args[0].String()), &design); err != nil {
		return errResult("decode design: " + err.Error())
	}

	resp, err := cachedEngine.EvaluateDesign(design, cachedRels)
	if err != nil {
		return errResult("evaluate: " + err.Error())
	}

	out, err := json.Marshal(resp)
	if err != nil {
		return errResult("encode response: " + err.Error())
	}
	return string(out)
}

func errResult(msg string) string {
	b, _ := json.Marshal(map[string]string{"error": msg})
	return string(b)
}

func main() {
	js.Global().Set("initEngine", js.FuncOf(initEngine))
	js.Global().Set("evaluateDesign", js.FuncOf(evaluateDesign))
	select {} // block forever; required so JS callbacks can keep firing
}
