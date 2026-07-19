import { describe, expect, it, vi } from 'vitest';

vi.mock('@/assets/icons', () => ({
  AddCircleOutlined: () => null,
}));

import {
  AddIconStyled,
  BtnText,
  YamlDialogTitle,
  YamlDialogTitleText,
} from './MesheryPatterns.styled';

describe('designs/patterns/MesheryPatterns.styled', () => {
  it('exports the expected styled wrappers', () => {
    // ViewSwitchButton, CreateButton and SearchWrapper were removed here —
    // their layout responsibilities moved to Sistent's DataTableToolbar
    // (see meshery/meshery#20658).
    //
    // NOTE: these are MUI/emotion styled() components, which React wraps in
    // forwardRef internally — so typeof is 'object', not 'function'. We
    // just confirm the exports exist rather than pin down that detail.
    expect(AddIconStyled).toBeDefined();
    expect(BtnText).toBeDefined();
    expect(YamlDialogTitle).toBeDefined();
    expect(YamlDialogTitleText).toBeDefined();
  });
});
