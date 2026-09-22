export {};

declare global {
  // Provided by Go's wasm_exec.js runtime.
  // Keep as `any` so we don't couple UI build to a specific Go version.
  const Go: any;

  // Exposed by server/policies/wasm/main.go.
  function initEngine(relsJson: string): string;
  function evaluateDesign(designJson: string): string;
}
