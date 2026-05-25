import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({ timeAgo: (v: string) => `AGO:${v}` }));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (id: string) => ({ name: `ctx-${id}` }),
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

import { buildSecretColumns } from './secret-columns';

describe('buildSecretColumns', () => {
  it('returns the documented Secret columns', () => {
    const cfg = buildSecretColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Secret',
      ping: () => {},
    });

    expect(cfg.name).toBe('Secret');
    expect(cfg.columns.map((c: any) => c.name)).toEqual([
      'id',
      'metadata.name',
      'apiVersion',
      'type',
      'metadata.namespace',
      'cluster_id',
      'metadata.creationTimestamp',
    ]);
  });
});
