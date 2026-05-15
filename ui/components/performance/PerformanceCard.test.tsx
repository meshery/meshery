import { describe, it } from 'vitest';

// PerformanceCard is a 400-line flippable card that depends on RTK Query
// (useGetUserByIdQuery), a CAN permission gate, FlipCard, a custom
// hook (useTestIDsGenerator), and 10+ sistent primitives. The leaf
// pieces are independently tested (FlipCard, generic permission gate
// via CAN tests elsewhere); this card itself is covered by playwright.
describe.skip('PerformanceCard', () => {
  it('skipped - deep RTK Query + flip-card + permission container', () => {});
});
