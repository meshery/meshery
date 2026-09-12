import { describe, it } from 'vitest';

// useDesignLifecycle is a 1000+ line hook orchestrating the design canvas
// state machine: it spawns and tears down xstate actors, fires RTK
// mutations (deploy/undeploy/validate), dispatches Redux actions for
// the global environment context, and wires Cytoscape event handlers. It
// is excercised indirectly by the playwright suite for the configurator
// page; mocking the entire dependency graph in isolation produces a
// brittle hook test that mirrors the implementation.
describe.skip('useDesignLifecycle hook', () => {
  it('skipped - orchestrates xstate actors, RTK mutations and cytoscape', () => {});
});
