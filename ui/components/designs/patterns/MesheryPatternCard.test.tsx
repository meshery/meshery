import { describe, it } from 'vitest';

// MesheryPatternCard renders a flippable card with a CodeMirror editor, a
// YAMLDialog (fullscreen), useTheme, VisibilityChipMenu, FlipCard, ActionButton,
// and 8+ TooltipButtons. Mocking the entire surface would more than triple the
// size of the component; the per-leaf pieces it depends on are individually
// tested (FlipCard, ActionButton, ActionPopover). Visual + interactive coverage
// of the card is delegated to e2e tests.
//
// The card's owner resolution - preferring the design's embedded `user` profile
// and skipping the by-id lookup when it is present, plus suppressing the
// Meshery Cloud avatar link on the built-in provider - is not part of that gap:
// it lives in the shared useResourceOwner hook and is covered directly by
// ui/utils/hooks/useResourceOwner.test.ts.
describe.skip('MesheryPatternCard', () => {
  it('skipped - too many leaf components and a redux/rtk-query owner', () => {});
});
