import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({ timeAgo: (v: string) => `AGO:${v}` }));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: () => ({ name: 'ctx' }),
}));
vi.mock('../../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: () => null,
}));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('../../../../utils/Enum', () => ({ CONNECTION_KINDS: { KUBERNETES: 'kubernetes' } }));
vi.mock('@/components/data-formatter', () => ({ FormatId: () => null }));
vi.mock('../../view', () => ({ Title: () => null }));
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { buildReplicationControllerColumns } from './replication-controller-columns';

describe('buildReplicationControllerColumns', () => {
  it('returns the ReplicationController name and non-empty columns', () => {
    const cfg = buildReplicationControllerColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'ReplicationController',
      ping: () => {},
    });
    expect(cfg.name).toBe('ReplicationController');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
