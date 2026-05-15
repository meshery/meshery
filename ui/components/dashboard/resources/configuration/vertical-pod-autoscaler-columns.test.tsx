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

import { buildVerticalPodAutoscalerColumns } from './vertical-pod-autoscaler-columns';

describe('buildVerticalPodAutoscalerColumns', () => {
  it('returns the VerticalPodAutoscaler name and non-empty columns', () => {
    const cfg = buildVerticalPodAutoscalerColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'VerticalPodAutoscaler',
      ping: () => {},
    });
    expect(cfg.name).toBe('VerticalPodAutoscaler');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
