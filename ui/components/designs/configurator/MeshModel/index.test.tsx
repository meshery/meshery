import { describe, it } from 'vitest';

// MeshModel/index.tsx is the entry to the design configurator: it owns the
// Cytoscape canvas, a useDesignLifecycle hook, a Redux store provider,
// xstate actors, and tens of sistent UI primitives. End-to-end coverage is
// provided by playwright; the leaf pieces (NodeIcon, utils, etc.) are
// unit-tested individually in this directory.
describe.skip('MeshModel/index (DesignConfigurator)', () => {
  it('skipped - top-level container with cytoscape + xstate + redux integration', () => {});
});
