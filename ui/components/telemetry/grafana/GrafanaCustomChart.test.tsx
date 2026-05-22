import { describe, it } from 'vitest';

// GrafanaCustomChart renders a billboard.js chart that polls a grafana
// or prometheus datasource. The full surface requires mocking
// billboard.js + the time-series fetch hooks + window.ResizeObserver.
// Visual coverage is provided by playwright.
describe.skip('GrafanaCustomChart', () => {
  it('skipped - billboard.js polling chart with fetch lifecycle', () => {});
});
