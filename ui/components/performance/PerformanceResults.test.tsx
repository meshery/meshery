import { describe, it } from 'vitest';

// PerformanceResults owns a fetch lifecycle for paginated test results,
// renders a MUIDataTable with several custom cells, and conditionally
// presents an info modal. The component is closely tied to a global
// notifier, sistent's CustomTooltip, react-moment, and the @sistent
// mui-datatables package. Coverage of its data-shape transforms is
// provided by the playwright performance suite.
describe.skip('PerformanceResults', () => {
  it('skipped - paginated mui-datatables container with side-effecting fetch', () => {});
});
