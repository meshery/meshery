import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

const useGetMeshSyncResourceKindsQuery = vi.fn();

vi.mock('../../../../utils/k8s-utils', () => ({ timeAgo: (v: string) => `AGO:${v}` }));
vi.mock('@/utils/multi-ctx', () => ({
  getK8sClusterIdsFromCtxId: () => [],
  getK8sContextFromClusterId: () => ({ name: 'ctx' }),
}));
vi.mock('../../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: () => null,
}));
vi.mock('@/utils/hooks/useKubernetesHook', () => ({ default: () => () => undefined }));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('../../../../utils/Enum', () => ({ CONNECTION_KINDS: { KUBERNETES: 'kubernetes' } }));
vi.mock('@/components/data-formatter', () => ({ FormatId: () => null }));
vi.mock('../../view', () => ({ Title: () => null }));
vi.mock('@/rtk-query/meshsync', () => ({
  useGetMeshSyncResourceKindsQuery: (...args: any[]) => useGetMeshSyncResourceKindsQuery(...args),
}));
vi.mock('../config', () => ({
  SINGLE_VIEW: 'single',
  getAllCustomResourceDefinitionsKinds: (kinds: any[]) => kinds ?? [],
}));

import { useCustomResourceConfig } from './config';

const callConfig = (kinds: any[]) => {
  useGetMeshSyncResourceKindsQuery.mockReturnValue({ data: { kinds } });
  let captured: any = null;
  const Probe = () => {
    captured = useCustomResourceConfig(vi.fn(), [], {}, {}, 'CRDS', []);
    return null;
  };
  render(<Probe />);
  return captured;
};

describe('useCustomResourceConfig', () => {
  it('returns an empty config when no kinds are present', () => {
    const result = callConfig([]);
    expect(result).toEqual({});
  });

  it('builds a config map keyed by custom resource Kind', () => {
    const result = callConfig([{ Kind: 'MyCRD', Model: 'my-model' }]);
    expect(result['MyCRD']).toBeDefined();
    expect(result['MyCRD'].name).toBe('MyCRD');
    expect(result['MyCRD'].model).toBe('my-model');
    expect(result['MyCRD'].columns.length).toBeGreaterThan(0);
  });
});
