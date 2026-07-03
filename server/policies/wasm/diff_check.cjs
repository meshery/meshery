// Boots the wasm engine in Node and runs one evaluation.
//
// Reads a JSON envelope `{design, rels}` from stdin, calls
// `initEngine(relsJSON)` once, then `evaluateDesign(designJSON)`, and prints
// the resulting `EvaluationResponse` JSON to stdout. Errors go to stderr with
// non-zero exit.
//
// Used by server/policies/wasm_diff_test.go to diff wasm output against the
// in-process Go engine. Standalone usage:
//   echo '{"design":{...},"rels":[...]}' | node diff_check.cjs

const fs = require('fs');
const path = require('path');

require(path.join(__dirname, 'wasm_exec.js'));

const go = new Go();
const wasmBuf = fs.readFileSync(path.join(__dirname, 'policy_engine.wasm'));

WebAssembly.instantiate(wasmBuf, go.importObject)
  .then(({ instance }) => {
    go.run(instance);
    const env = JSON.parse(fs.readFileSync(0, 'utf8'));

    const initRes = globalThis.initEngine(JSON.stringify(env.rels));
    if (initRes) {
      const parsed = JSON.parse(initRes);
      if (parsed && parsed.error) {
        console.error(parsed.error);
        process.exit(1);
      }
    }

    const out = globalThis.evaluateDesign(JSON.stringify(env.design));
    const parsed = JSON.parse(out);
    if (parsed.error) {
      console.error(parsed.error);
      process.exit(1);
    }
    // Wait for the write to flush before exiting; on a piped stdout, calling
    // process.exit(0) immediately can truncate large payloads at the OS pipe
    // buffer size (commonly 64 KB on macOS/Linux).
    process.stdout.write(out, () => process.exit(0));
  })
  .catch((e) => {
    console.error('wasm boot error:', e);
    process.exit(2);
  });
