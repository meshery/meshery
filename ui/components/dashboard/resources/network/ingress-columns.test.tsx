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

import { buildIngressColumns } from './ingress-columns';

describe('buildIngressColumns', () => {
  it('returns the Ingress name and non-empty columns', () => {
    const cfg = buildIngressColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Ingress',
      ping: () => {},
    });
    expect(cfg.name).toBe('Ingress');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
