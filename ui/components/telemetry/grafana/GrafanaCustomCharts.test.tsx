import { describe, it } from 'vitest';

// GrafanaCustomCharts iterates over a board+panel config and mounts
// either a GrafanaCustomChart or a GrafanaCustomGaugeChart per panel.
// The leaf chart components are skipped (billboard.js + fetch
// dependency); this wrapper has no additional logic worth isolating.
describe.skip('GrafanaCustomCharts', () => {
  it('skipped - thin map over GrafanaCustomChart / GrafanaCustomGaugeChart', () => {});
});
