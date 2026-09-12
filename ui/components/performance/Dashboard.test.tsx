import { describe, it } from 'vitest';

// Dashboard is the performance landing page; it wires multiple Redux
// selectors, several MUI Tabs, lazily loaded MesheryMetrics and
// PerformanceProfiles children, plus a graphql control plane fetch.
// Top-level rendering depends on a configured store and an active
// kubernetes context which are not available in a unit context.
describe.skip('Performance Dashboard', () => {
  it('skipped - container page owning multiple tabs with redux + graphql', () => {});
});
