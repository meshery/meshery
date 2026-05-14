import { describe, expect, it, vi } from 'vitest';

vi.mock('./resources-table.tsx', () => ({
  default: function ResourcesTableMock() {
    return null;
  },
  ACTION_TYPES: {
    FETCH_MESHSYNC_RESOURCES: {
      name: 'FETCH_MESHSYNC_RESOURCES',
      error_msg: 'Failed to fetch meshsync resources',
    },
  },
}));

describe('resources-table.ts barrel', () => {
  it('re-exports the default ResourcesTable component and ACTION_TYPES', async () => {
    const mod = await import('./resources-table');
    expect(mod.default).toBeDefined();
    expect(mod.ACTION_TYPES).toBeDefined();
    expect(mod.ACTION_TYPES.FETCH_MESHSYNC_RESOURCES.name).toBe('FETCH_MESHSYNC_RESOURCES');
  });
});
