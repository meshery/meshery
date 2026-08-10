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
vi.mock('../../../../utils/Enum', () => ({
  CoreConnectionKinds: {
    meshery: 'meshery',
    kubernetes: 'kubernetes',
    prometheus: 'prometheus',
    grafana: 'grafana',
    github: 'github',
  },
}));
vi.mock('@/components/data-formatter', () => ({ FormatId: () => null }));
vi.mock('../../view', () => ({ Title: () => null }));
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { buildNetworkPolicyColumns } from './network-policy-columns';

describe('buildNetworkPolicyColumns', () => {
  it('returns a name and non-empty columns array', () => {
    const cfg = buildNetworkPolicyColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'NetworkPolicy',
      ping: () => {},
    });
    expect(typeof cfg.name).toBe('string');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
