/**
 * Barrel of typed SVG icons for Meshery UI.
 *
 * This file exposes the canonical, typed SVG replacements for repeated
 * `@mui/icons-material` imports, as part of Phase 2 of the UI restructure
 * (parent #18657, sub-issue #18730).
 *
 * Conventions:
 *  - Each icon is a typed React functional component accepting `IconProps`
 *    (standard SVG props plus optional `width`/`height`/`fill` overrides).
 *  - `fill` defaults to `currentColor` so icons inherit color from CSS.
 *  - Each module also provides a `default` export for ergonomic
 *    `import EditIcon from '...'` consumption.
 *
 * Downstream sub-issues (#18733-#18740) will swap the corresponding
 * `@mui/icons-material/<Name>` imports across the codebase to use these
 * barrel exports (or Sistent equivalents where Sistent ships the glyph).
 * This issue is purely additive — no call sites are modified here.
 */

export type { IconProps } from './types';

// Typed SVG replacements added in #18730 for icons used >=3 times in `ui/`.
export { ArrowDropDownIcon } from './ArrowDropDownIcon';
export { ChevronLeftIcon } from './ChevronLeftIcon';
export { ChevronRightIcon } from './ChevronRightIcon';
export { EditIcon } from './EditIcon';
export { FullscreenExitIcon } from './FullscreenExitIcon';
export { FullscreenIcon } from './FullscreenIcon';
export { GetAppIcon } from './GetAppIcon';
export { LockIcon } from './LockIcon';
export { SaveIcon } from './SaveIcon';
export { SettingsIcon } from './SettingsIcon';

// Pre-existing typed icons in `ui/assets/icons/` whose glyphs also appear
// >=3 times via `@mui/icons-material`. Re-exported here so that downstream
// migrations can import the entire canonical set from a single path.
//
// Note: these are intentionally re-exported as `default` since the source
// files only provide default exports today. Renaming their source exports
// is out of scope for this purely-additive issue.
export { default as DeleteIcon } from './DeleteIcon';
export { default as ExpandMoreIcon } from './ExpandMoreIcon';
export { default as InfoOutlinedIcon } from './InfoOutlined';
