import { HighlightStyle, syntaxHighlighting } from '@codemirror/language';
import { tags } from '@lezer/highlight';
import { materialDarkStyle } from '@uiw/codemirror-theme-material';

// Same green Material already uses for quoted `string` tokens (`theme={material}`).
export const MATERIAL_STRING = materialDarkStyle.find(
  ({ tag }) => Array.isArray(tag) && tag.includes(tags.string),
)?.color;

/**
 * @lezer/yaml marks quoted scalars as `string` but plain (unquoted) scalars
 * as `content`. Themes only color `string`, which is why `name: foo` looks
 * unstyled next to `name: "foo"`. Color those content tokens as strings too.
 */
export const yamlPlainScalarHighlight = syntaxHighlighting(
  HighlightStyle.define([{ tag: tags.content, color: MATERIAL_STRING }]),
);
