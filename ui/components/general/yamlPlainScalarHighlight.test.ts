import { describe, expect, it, vi } from 'vitest';
import { tags } from '@lezer/highlight';

vi.mock('@codemirror/language', () => ({
  HighlightStyle: {
    define: (rules: unknown) => rules,
  },
  syntaxHighlighting: (style: unknown) => ({ marker: 'yaml-scalar', style }),
}));

import { MATERIAL_STRING, yamlPlainScalarHighlight } from './yamlPlainScalarHighlight';

describe('yamlPlainScalarHighlight', () => {
  it('applies the Material string color to unquoted YAML content tokens only', () => {
    const extension = yamlPlainScalarHighlight as unknown as {
      marker: string;
      style: { tag: unknown; color: string }[];
    };
    expect(MATERIAL_STRING).toBe('#99d066');
    expect(extension.marker).toBe('yaml-scalar');
    expect(extension.style).toEqual([{ tag: tags.content, color: MATERIAL_STRING }]);
  });
});
