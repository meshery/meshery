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

// evaluateDesign is exposed to JS as globalThis.evaluateDesign(designJSON, relsJSON).
// Returns a JSON string of pattern.EvaluationResponse, or {"error": "..."} on failure.
//
// The caller is responsible for hydrating the design and providing the registered
// relationships; this entry point intentionally does no registry I/O.
func evaluateDesign(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		return errResult("evaluateDesign(designJSON, relsJSON): expected 2 arguments")
	}

	var design pattern.PatternFile
	if err := json.Unmarshal([]byte(args[0].String()), &design); err != nil {
		return errResult("decode design: " + err.Error())
	}

	var relsRaw []any
	if err := json.Unmarshal([]byte(args[1].String()), &relsRaw); err != nil {
		return errResult("decode relationships: " + err.Error())
	}
	rels := gopolicies.ConvertRelationships(relsRaw)

	engine := gopolicies.NewGoEngine(noopLogger{})
	resp, err := engine.EvaluateDesign(design, rels)
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

// Keep the relationship schema package in the import closure so
// ConvertRelationships's input type is not elided.
var _ = relationship.RelationshipDefinition{}

func main() {
	js.Global().Set("evaluateDesign", js.FuncOf(evaluateDesign))
	select {} // block forever; required so JS callbacks can keep firing
}
