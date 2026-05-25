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
vi.mock('../../../../utils/utils', () => ({
  ResizableCell: ({ value }: any) => value,
}));

import { buildPodColumns } from './pod-columns';

describe('buildPodColumns', () => {
  it('returns the Pod name and non-empty columns', () => {
    const cfg = buildPodColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Pod',
      ping: () => {},
    });
    expect(cfg.name).toBe('Pod');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });
});
