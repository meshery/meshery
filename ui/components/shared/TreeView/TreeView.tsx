/**
 * Shared TreeView wrapper.
 *
 * This module is the single application-level boundary for `@mui/x-tree-view`.
 * All app code MUST consume tree-view primitives from `ui/components/shared/TreeView`
 * rather than importing `@mui/x-tree-view` directly. The ESLint rule
 * `no-restricted-imports` enforces this boundary; the only intentional escape
 * hatch is this file.
 *
 * Re-exports the primitives currently used by the Registry surfaces. Add more
 * primitives here (e.g. `RichTreeView`, `TreeView`) on an as-needed basis;
 * keep the surface intentionally small.
 */
/* eslint-disable no-restricted-imports */
export { SimpleTreeView, TreeItem, treeItemClasses } from '@mui/x-tree-view';
export type {
  SimpleTreeViewProps,
  SimpleTreeViewSlots,
  SimpleTreeViewSlotProps,
} from '@mui/x-tree-view/SimpleTreeView';
export type { TreeItemProps } from '@mui/x-tree-view/TreeItem';
