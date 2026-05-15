import { describe, it } from 'vitest';

// GrafanaMetricsCompare is a compare-view container that pulls multiple
// time-ranges of telemetry data via GrafanaCustomChart and overlays
// them. Like the other charts it depends on billboard.js + ResizeObserver
// + a fetch polling lifecycle that is impractical to mock at unit
// granularity.
describe.skip('GrafanaMetricsCompare', () => {
  it('skipped - compare view of GrafanaCustomChart with timer polling', () => {});
});
