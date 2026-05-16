//go:build js && wasm

// Command relationship-engine builds the Meshery relationship policy
// engine to WebAssembly so the in-browser policy-engine consumer can
// evaluate designs entirely client-side, with output that is at parity
// with Meshery Server's HTTP relationship-evaluation endpoint.
//
// Meshery is the source of truth for the policy engine. This entrypoint
// is a thin syscall/js boundary around the SAME shared,
// registry-INDEPENDENT orchestration that the HTTP handler uses
// (server/policies.(*GoEngine).EvaluateDesignOrchestrated). It performs
// NO registry hydration — there is no component registry in the browser;
// the in-browser client compensates for component hydration on its side.
//
// Contract (attached to globalThis.__mesheryRelationshipEngine):
//
//	contract: 1                              (integer, bump on breaking change)
//	version:  <meshery build version / SHA>  (string, diagnostics only)
//	init(relationshipsJSON string) string     -> "{}" | {"error":"..."}
//	evaluate(designJSON, optionsJSON) string  -> EvaluationResponse JSON | {"error":"..."}
//	selfTest() string                         -> {"ok":true,...} | {"error":"..."}
//
// All exported functions are synchronous string->string, JSON in/out,
// and never panic across the JS boundary (every call is wrapped in a
// recover() that converts a panic into {"error":"..."}).
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"syscall/js"

	"github.com/meshery/meshery/server/policies"
	"github.com/meshery/meshkit/logger"
	schemav1alpha3 "github.com/meshery/schemas/models/v1alpha3"
	relationshipv1alpha3 "github.com/meshery/schemas/models/v1alpha3/relationship"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	relationshipv1beta2 "github.com/meshery/schemas/models/v1beta2/relationship"

	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

// contractVersion is the integer contract version the in-browser client
// pins against. Bump ONLY on a breaking change to the exported surface.
const contractVersion = 1

// relationshipSchemaVersionV1alpha3 is the wire schemaVersion value the
// registry uses for v1alpha3 relationship definitions. It is pinned to
// the schemas source of truth so it cannot drift. Everything else
// (notably relationships.meshery.io/v1beta2) is decoded as v1beta2,
// which is the policy engine's native relationship type.
const relationshipSchemaVersionV1alpha3 = schemav1alpha3.RelationshipSchemaVersion

// buildVersion is injected at build time via
// -ldflags "-X main.buildVersion=<meshery git sha>". It is diagnostic
// only and never affects evaluation output.
var buildVersion = "dev"

// engineState holds the initialized engine and its registered
// relationships. Guarded by mu because, although the Go wasm runtime is
// single-threaded today, the contract permits init/evaluate to be
// called in any order and we must never read a half-written state.
type engineState struct {
	mu     sync.RWMutex
	engine *policies.GoEngine
	rels   []*relationshipv1beta2.RelationshipDefinition
	ready  bool
}

var state engineState

// noopLogger is a wasm-safe logger.Handler that discards all output.
// The policy engine only ever calls Info on its logger; the remaining
// interface methods exist solely to satisfy logger.Handler and are
// inert. Using a no-op handler (instead of logger.New, which writes to
// os.Stdout/os.Stderr) keeps the browser console clean and avoids any
// runtime surprises from the syslog/terminal formatters.
type noopLogger struct{}

func (noopLogger) Info(...interface{})                  {}
func (noopLogger) Infof(string, ...interface{})         {}
func (noopLogger) Debug(...interface{})                 {}
func (noopLogger) Debugf(string, ...interface{})        {}
func (noopLogger) Warn(error)                           {}
func (noopLogger) Warnf(string, ...interface{})         {}
func (noopLogger) Error(error)                          {}
func (noopLogger) Errorf(string, ...interface{})        {}
func (noopLogger) Fatal(error)                          {}
func (noopLogger) Fatalf(string, ...interface{})        {}
func (noopLogger) SetLevel(logrus.Level)                {}
func (noopLogger) GetLevel() logrus.Level               { return logrus.InfoLevel }
func (noopLogger) UpdateLogOutput(io.Writer)            {}
func (noopLogger) UpdateErrorLogOutput(io.Writer)       {}
func (noopLogger) ControllerLogger() logr.Logger        { return logr.Discard() }
func (noopLogger) DatabaseLogger() gormlogger.Interface { return gormlogger.Default }

// compile-time assertion that noopLogger satisfies the interface the
// engine constructor expects.
var _ logger.Handler = noopLogger{}

// errorJSON renders {"error":"..."} with the error text JSON-escaped.
func errorJSON(format string, args ...interface{}) string {
	b, _ := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: fmt.Sprintf(format, args...)})
	return string(b)
}

// recoverToError converts a recovered panic into the standard error
// envelope so a panic can never cross the JS boundary.
func recoverToError(out *string) {
	if r := recover(); r != nil {
		*out = errorJSON("relationship engine panic: %v", r)
	}
}

// initEngine builds (or rebuilds) the engine from the relationships JSON
// array exactly as returned by
// GET /api/meshmodels/relationships?page=0&pagesize=all (the
// `.relationships` array — mixed v1alpha3/v1beta2). Re-calling init
// replaces the registered relationships.
func initEngine(relationshipsJSON string) (out string) {
	defer recoverToError(&out)

	var raw []json.RawMessage
	if err := json.Unmarshal([]byte(relationshipsJSON), &raw); err != nil {
		return errorJSON("failed to parse relationships JSON: %v", err)
	}

	// Discriminate each element by schemaVersion and decode into the
	// matching typed RelationshipDefinition, then hand the heterogeneous
	// slice to policies.ConvertRelationships exactly as the HTTP handler
	// does with the registry's typed entities.
	relInterfaces := make([]interface{}, 0, len(raw))
	for i, rm := range raw {
		var probe struct {
			SchemaVersion string `json:"schemaVersion"`
		}
		if err := json.Unmarshal(rm, &probe); err != nil {
			return errorJSON("failed to inspect relationship[%d]: %v", i, err)
		}

		if probe.SchemaVersion == relationshipSchemaVersionV1alpha3 {
			var rel relationshipv1alpha3.RelationshipDefinition
			if err := json.Unmarshal(rm, &rel); err != nil {
				return errorJSON("failed to decode v1alpha3 relationship[%d]: %v", i, err)
			}
			relInterfaces = append(relInterfaces, &rel)
			continue
		}

		// Default path: v1beta2 (engine-native) and anything not
		// explicitly v1alpha3.
		var rel relationshipv1beta2.RelationshipDefinition
		if err := json.Unmarshal(rm, &rel); err != nil {
			return errorJSON("failed to decode v1beta2 relationship[%d]: %v", i, err)
		}
		relInterfaces = append(relInterfaces, &rel)
	}

	converted := policies.ConvertRelationships(relInterfaces)
	eng := policies.NewGoEngine(noopLogger{})

	state.mu.Lock()
	state.engine = eng
	state.rels = converted
	state.ready = true
	state.mu.Unlock()

	return "{}"
}

// evaluate runs the shared registry-independent orchestration against
// the supplied design. designJSON is a v1beta1 pattern.PatternFile;
// optionsJSON is the optional {"enableTrace":bool,"returnDiffOnly":bool}.
func evaluate(designJSON, optionsJSON string) (out string) {
	defer recoverToError(&out)

	state.mu.RLock()
	eng := state.engine
	rels := state.rels
	ready := state.ready
	state.mu.RUnlock()

	if !ready || eng == nil {
		return errorJSON("engine not initialized")
	}

	var design pattern.PatternFile
	if err := json.Unmarshal([]byte(designJSON), &design); err != nil {
		return errorJSON("failed to parse design JSON: %v", err)
	}

	opts := policies.EvaluationOptions{}
	if optionsJSON != "" {
		var parsed pattern.Options
		if err := json.Unmarshal([]byte(optionsJSON), &parsed); err != nil {
			return errorJSON("failed to parse options JSON: %v", err)
		}
		if parsed.EnableTrace != nil {
			opts.EnableTrace = *parsed.EnableTrace
		}
		if parsed.ReturnDiffOnly != nil {
			opts.ReturnDiffOnly = *parsed.ReturnDiffOnly
		}
	}

	resp, err := eng.EvaluateDesignOrchestrated(design, rels, policies.MAX_RE_EVALUATION_DEPTH, opts)
	if err != nil {
		return errorJSON("evaluation failed: %v", err)
	}

	b, err := json.Marshal(resp)
	if err != nil {
		return errorJSON("failed to encode evaluation response: %v", err)
	}
	return string(b)
}

// selfTest initializes the engine with a minimal hardcoded relationship
// set and evaluates a minimal hardcoded design, exercising the real
// engine path end to end. It returns {"ok":true,...} on success.
func selfTest() (out string) {
	defer recoverToError(&out)

	// A single, well-formed v1beta2 hierarchical relationship. It does
	// not need to match the test design — selfTest only proves the
	// init -> evaluate path runs without error through the real engine.
	const minimalRelationships = `[
	  {
	    "schemaVersion": "relationships.meshery.io/v1beta2",
	    "version": "v1.0.0",
	    "kind": "hierarchical",
	    "type": "parent",
	    "subType": "inventory",
	    "model": {"name": "kubernetes", "model": {"version": "*"}},
	    "selectors": []
	  }
	]`

	if res := initEngine(minimalRelationships); res != "{}" {
		return errorJSON("selfTest init failed: %s", res)
	}

	const minimalDesign = `{
	  "id": "00000000-0000-0000-0000-000000000000",
	  "name": "wasm-selftest",
	  "schemaVersion": "designs.meshery.io/v1beta1",
	  "version": "0.0.1",
	  "components": [
	    {
	      "id": "11111111-1111-1111-1111-111111111111",
	      "component": {"kind": "Namespace", "version": "v1"},
	      "model": {"name": "kubernetes", "model": {"version": "*"}},
	      "configuration": {}
	    },
	    {
	      "id": "22222222-2222-2222-2222-222222222222",
	      "component": {"kind": "Pod", "version": "v1"},
	      "model": {"name": "kubernetes", "model": {"version": "*"}},
	      "configuration": {}
	    }
	  ],
	  "relationships": []
	}`

	res := evaluate(minimalDesign, `{"enableTrace":true}`)
	var probe struct {
		Error string `json:"error"`
	}
	if err := json.Unmarshal([]byte(res), &probe); err == nil && probe.Error != "" {
		return errorJSON("selfTest evaluate failed: %s", probe.Error)
	}

	b, _ := json.Marshal(struct {
		OK       bool   `json:"ok"`
		Version  string `json:"version"`
		Contract int    `json:"contract"`
	}{OK: true, Version: buildVersion, Contract: contractVersion})
	return string(b)
}

// stringFuncOf adapts a Go (args...) -> string function into a
// js.Func, coercing every argument to a string and never letting a
// panic escape into JS.
func stringFuncOf(fn func(args []string) string) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) (result interface{}) {
		defer func() {
			if r := recover(); r != nil {
				result = errorJSON("relationship engine panic: %v", r)
			}
		}()
		strArgs := make([]string, len(args))
		for i, a := range args {
			strArgs[i] = a.String()
		}
		return fn(strArgs)
	})
}

func main() {
	api := map[string]interface{}{
		"contract": contractVersion,
		"version":  buildVersion,
		"init": stringFuncOf(func(args []string) string {
			if len(args) < 1 {
				return errorJSON("init requires the relationships JSON argument")
			}
			return initEngine(args[0])
		}),
		"evaluate": stringFuncOf(func(args []string) string {
			if len(args) < 1 {
				return errorJSON("evaluate requires the design JSON argument")
			}
			optionsJSON := ""
			if len(args) >= 2 {
				optionsJSON = args[1]
			}
			return evaluate(args[0], optionsJSON)
		}),
		"selfTest": stringFuncOf(func(_ []string) string {
			return selfTest()
		}),
	}

	js.Global().Set("__mesheryRelationshipEngine", js.ValueOf(api))

	// Keep the wasm instance alive; the in-browser host invokes the
	// exported functions synchronously over the lifetime of the page.
	select {}
}
