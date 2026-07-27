/**
 * Shared TreeView wrapper.
 *
 * This module is the single application-level boundary for `@mui/x-tree-view`.
 * All app code MUST consume tree-view primitives from `ui/components/shared/TreeView`
 * rather than importing `@mui/x-tree-view` directly. The ESLint rule
 * `no-restricted-imports` enforces this boundary; the only intentional escape
 * hatch is this file.
 */
/* eslint-disable no-restricted-imports */
import React from 'react';
import { SimpleTreeView, TreeItem as MuiTreeItem, treeItemClasses } from '@mui/x-tree-view';
import { Collapse } from '@sistent/sistent';

const GroupTransitionComponent = (Collapse as any)?.default || Collapse;

export const TreeItem = React.forwardRef((props: any, ref: any) => {
  return (
    <MuiTreeItem
      ref={ref}
      {...props}
      slots={{ groupTransition: GroupTransitionComponent, ...props?.slots }}
    />
  );
});
TreeItem.displayName = 'TreeItem';

export { SimpleTreeView, treeItemClasses };
export type {
  SimpleTreeViewProps,
  SimpleTreeViewSlots,
  SimpleTreeViewSlotProps,
} from '@mui/x-tree-view/SimpleTreeView';
export type { TreeItemProps } from '@mui/x-tree-view/TreeItem';
