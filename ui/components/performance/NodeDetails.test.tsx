import { describe, it } from 'vitest';

// NodeDetails renders cluster-node detail panels backed by Cytoscape
// node selection events and a Prometheus boards config in redux. It
// mounts a custom Grafana chart per board panel and is exercised by
// the playwright performance suite when a connected k8s cluster is
// configured.
describe.skip('NodeDetails', () => {
  it('skipped - cytoscape event-driven panel coupled to redux + grafana charts', () => {});
});
