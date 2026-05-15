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

import { buildDaemonSetColumns } from './daemonset-columns';

describe('buildDaemonSetColumns', () => {
  it('returns the DaemonSet name and non-empty columns', () => {
    const cfg = buildDaemonSetColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'DaemonSet',
      ping: () => {},
    });
    expect(cfg.name).toBe('DaemonSet');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
