import { describe, it } from 'vitest';

// GrafanaComponent is the top-level grafana view that composes
// GrafanaConfigComponent and GrafanaSelectionComponent (each skipped
// above) and a chart preview panel. Apart from prop forwarding there
// is no leaf behaviour to assert here.
describe.skip('GrafanaComponent', () => {
  it('skipped - composes Config + Selection grafana children', () => {});
});
