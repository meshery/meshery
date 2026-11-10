package policies

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/meshery/meshkit/logger"
	"github.com/meshery/schemas/models/v1beta1/component"
	"github.com/meshery/schemas/models/v1beta1/pattern"
	"github.com/meshery/schemas/models/v1beta2/relationship"
)

// repoRoot walks up from the test working directory to the directory containing the Makefile.
func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for range 8 {
		if _, err := os.Stat(filepath.Join(dir, "Makefile")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not locate repo root from " + dir)
	return ""
}

var (
	wasmEnvOnce  sync.Once
	wasmRoot     string
	wasmEnvReady bool
	wasmSkipMsg  string
)

// wasmDiffEnv resolves whether the wasm artifact + node are both present.
// Resolved once per test process and reused; absent env causes silent skip.
func wasmDiffEnv(t *testing.T) (root string, ok bool, skip string) {
	t.Helper()
	wasmEnvOnce.Do(func() {
		wasmRoot = repoRoot(t)
		if _, err := os.Stat(filepath.Join(wasmRoot, "server/policies/wasm/policy_engine.wasm")); err != nil {
			wasmSkipMsg = "policy_engine.wasm not built; run `make wasm-engine` first"
			return
		}
		if _, err := exec.LookPath("node"); err != nil {
			wasmSkipMsg = "node not installed; skipping wasm diff"
			return
		}
		wasmEnvReady = true
	})
	return wasmRoot, wasmEnvReady, wasmSkipMsg
}

// assertWasmMatches feeds (design, rels) through the wasm engine via Node and
// asserts byte-equal JSON against respNative. Silently no-ops when the wasm
// artifact or node are not available, so existing test runs are unaffected
// unless the caller has explicitly built the wasm binary.
func assertWasmMatches(
	t *testing.T,
	design pattern.PatternFile,
	rels []*relationship.RelationshipDefinition,
	respNative pattern.EvaluationResponse,
) {
	t.Helper()
	root, ok, _ := wasmDiffEnv(t)
	if !ok {
		return
	}

	envelope, err := json.Marshal(struct {
		Design any `json:"design"`
		Rels   any `json:"rels"`
	}{design, rels})
	if err != nil {
		t.Fatalf("wasm diff: encode envelope: %v", err)
	}

	script := filepath.Join(root, "server/policies/wasm/diff_check.cjs")
	cmd := exec.Command("node", script)
	cmd.Stdin = bytes.NewReader(envelope)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("wasm diff: node run: %v\nstderr: %s", err, stderr.String())
	}

	// Re-encode native through json.Marshal so both sides go through the same
	// canonicalization (nil-vs-empty slice, key ordering) before comparison.
	nativeJSON, err := json.Marshal(respNative)
	if err != nil {
		t.Fatalf("wasm diff: encode native: %v", err)
	}
	if !bytes.Equal(nativeJSON, bytes.TrimRight(stdout.Bytes(), "\n")) {
		t.Fatalf("native vs wasm mismatch\nnative=%s\nwasm  =%s", nativeJSON, stdout.String())
	}
}

// TestWasmEngineMatchesNativeOnFixture runs the production-grade fixture
// shared with the rego-eval Make target through both engines and asserts the
// wasm output matches native byte-for-byte. Doubles as a determinism check:
// any pointer-address leak in canonicalSeed reappears here as a mismatch.
//
// This bypasses assertWasmMatches because we need both engines to start from
// the *same raw fixture bytes*. The default helper marshals the caller's typed
// Go struct into the envelope, which would round-trip through Go's encoder
// before the wasm sees it; the fixture contains shapes where that round-trip
// is not byte-stable, so we feed the original JSON through directly.
func TestWasmEngineMatchesNativeOnFixture(t *testing.T) {
	root, ok, msg := wasmDiffEnv(t)
	if !ok {
		t.Skip(msg)
	}

	designBytes, err := os.ReadFile(filepath.Join(root, "policies/test/design_all_relationships.yaml"))
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	relsBytes, err := os.ReadFile(filepath.Join(root, "policies/test/all_relationships.json"))
	if err != nil {
		t.Skipf("fixture missing: %v", err)
	}

	var design pattern.PatternFile
	if err := json.Unmarshal(designBytes, &design); err != nil {
		t.Fatalf("decode design fixture: %v", err)
	}
	var rels []*relationship.RelationshipDefinition
	if err := json.Unmarshal(relsBytes, &rels); err != nil {
		t.Fatalf("decode rels fixture: %v", err)
	}

	log, err := logger.New("wasm-diff", logger.Options{Format: logger.SyslogLogFormat})
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	respNative, err := NewGoEngine(log).EvaluateDesign(design, rels)
	if err != nil {
		t.Fatalf("native eval: %v", err)
	}

	// Build envelope from raw fixture bytes (no typed-struct re-marshal).
	var envelope bytes.Buffer
	envelope.WriteString(`{"design":`)
	envelope.Write(designBytes)
	envelope.WriteString(`,"rels":`)
	envelope.Write(relsBytes)
	envelope.WriteByte('}')

	script := filepath.Join(root, "server/policies/wasm/diff_check.cjs")
	cmd := exec.Command("node", script)
	cmd.Stdin = &envelope
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("wasm diff: node run: %v\nstderr: %s", err, stderr.String())
	}

	nativeJSON, err := json.Marshal(respNative)
	if err != nil {
		t.Fatalf("encode native: %v", err)
	}
	wasmJSON := bytes.TrimRight(stdout.Bytes(), "\n")
	if !bytes.Equal(nativeJSON, wasmJSON) {
		t.Fatalf("native vs wasm mismatch (sizes: native=%d wasm=%d)", len(nativeJSON), len(wasmJSON))
	}
}

// TestWasmEngineMatchesNative is the dedicated diff smoke test using the
// orphaned-edge fixture from TestValidateRelationshipsInDesign. The other
// EvaluateDesign-touching tests opt into the same diff via assertWasmMatches.
func TestWasmEngineMatchesNative(t *testing.T) {
	if _, ok, msg := wasmDiffEnv(t); !ok {
		t.Skip(msg)
	}

	compA := &component.ComponentDefinition{}
	compA.ID, _ = uuid.FromString("00000000-0000-0000-0000-0000000000aa")
	compMissing := &component.ComponentDefinition{}
	compMissing.ID, _ = uuid.FromString("00000000-0000-0000-0000-00000000dead")

	relStatus := relationship.RelationshipDefinitionStatus("approved")
	rel := &relationship.RelationshipDefinition{
		Kind:             relationship.RelationshipDefinitionKind("edge"),
		RelationshipType: "non-binding",
		Status:           &relStatus,
		Selectors: &relationship.SelectorSet{
			relationship.SelectorSetItem{
				Allow: relationship.Selector{
					From: []relationship.SelectorItem{{ID: &compA.ID}},
					To:   []relationship.SelectorItem{{ID: &compMissing.ID}},
				},
			},
		},
	}
	rel.ID, _ = uuid.FromString("00000000-0000-0000-0000-000000000001")

	design := *makePatternFile([]*component.ComponentDefinition{compA}, []*relationship.RelationshipDefinition{rel})
	rels := []*relationship.RelationshipDefinition{rel}

	log, err := logger.New("wasm-diff", logger.Options{Format: logger.SyslogLogFormat})
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	respNative, err := NewGoEngine(log).EvaluateDesign(design, rels)
	if err != nil {
		t.Fatalf("native eval: %v", err)
	}

	assertWasmMatches(t, design, rels, respNative)
}
