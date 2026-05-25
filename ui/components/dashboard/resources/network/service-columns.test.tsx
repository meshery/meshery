import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  getStatus: (v: any) => `STATUS:${v}`,
  timeAgo: (v: string) => `AGO:${v}`,
}));
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

import { buildServiceColumns } from './service-columns';

describe('buildServiceColumns', () => {
  it('returns the Service name and a non-empty columns array', () => {
    const cfg = buildServiceColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Service',
      ping: () => {},
    });
    expect(cfg.name).toBe('Service');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
