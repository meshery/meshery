import { describe, expect, it, vi } from 'vitest';

vi.mock('@mui/x-tree-view', () => ({
  SimpleTreeView: () => null,
  TreeItem: () => null,
  treeItemClasses: { selected: 'mui-selected' },
}));

describe('TreeView index re-export', () => {
  it('forwards SimpleTreeView and TreeItem from ./TreeView', async () => {
    const mod = await import('./index');
    expect(mod.SimpleTreeView).toBeDefined();
    expect(mod.TreeItem).toBeDefined();
    expect(mod.treeItemClasses).toBeDefined();
  });
});
