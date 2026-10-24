import { describe, expect, it, vi } from 'vitest';

vi.mock('@mui/x-tree-view', () => ({
  SimpleTreeView: () => null,
  TreeItem: () => null,
  treeItemClasses: { selected: 'mui-selected', focused: 'mui-focused' },
}));

vi.mock('@mui/x-tree-view/SimpleTreeView', () => ({
  // type-only file
}));

vi.mock('@mui/x-tree-view/TreeItem', () => ({}));

describe('TreeView shared wrapper', () => {
  it('re-exports SimpleTreeView and TreeItem from @mui/x-tree-view', async () => {
    const mod = await import('./TreeView');
    expect(mod.SimpleTreeView).toBeDefined();
    expect(mod.TreeItem).toBeDefined();
    expect(mod.treeItemClasses).toBeDefined();
    expect(mod.treeItemClasses.selected).toBe('mui-selected');
  });
});
