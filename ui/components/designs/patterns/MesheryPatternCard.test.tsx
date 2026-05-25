import { describe, it } from 'vitest';

// MesheryPatternCard renders a flippable card with a CodeMirror editor, a
// YAMLDialog (fullscreen), a useGetUserByIdQuery RTK hook, useTheme,
// VisibilityChipMenu, FlipCard, ActionButton, and 8+ TooltipButtons. Mocking
// the entire surface would more than triple the size of the component; the
// per-leaf pieces it depends on are individually tested (FlipCard,
// ActionButton, ActionPopover). Visual + interactive coverage of the card
// is delegated to e2e tests.
describe.skip('MesheryPatternCard', () => {
  it('skipped - too many leaf components and a redux/rtk-query owner', () => {});
});
