import { describe, it } from 'vitest';

// PerformanceProfiles is the table+grid view manager for performance
// profiles. It is wired to RTK Query mutations / queries, a redux slice
// for the selected context, a notification hook, a router, and several
// modals. The grid sub-view is unit-tested in
// PerformanceProfileGrid.test.tsx; this top-level page is covered by
// the playwright performance suite.
describe.skip('PerformanceProfiles', () => {
  it('skipped - container backed by RTK Query, redux and several modals', () => {});
});
