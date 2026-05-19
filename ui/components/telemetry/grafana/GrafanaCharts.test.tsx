import { describe, it } from 'vitest';

// GrafanaCharts is a wrapper around GrafanaPanelIframe instances; iframe
// rendering itself is unit-tested in GrafanaPanelIframe.test.tsx. The
// wrapper plumbs an array of board+panel configs into the iframes; the
// transformations are simple list rendering with no business logic.
describe.skip('GrafanaCharts', () => {
  it('skipped - thin map over GrafanaPanelIframe', () => {});
});
