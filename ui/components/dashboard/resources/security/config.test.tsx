import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  timeAgo: (val: string) => `AGO:${val}`,
}));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (id: string) => ({ name: `ctx-${id}` }),
}));
vi.mock('../../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title }: any) => <div>{title}</div>,
}));
vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => () => undefined,
}));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('../../../../utils/Enum', () => ({
  CONNECTION_KINDS: { KUBERNETES: 'kubernetes' },
}));
vi.mock('@/components/data-formatter', () => ({
  FormatId: ({ id }: any) => id,
}));
vi.mock('../../view', () => ({
  Title: ({ value }: any) => value,
}));
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { SecurityTypesConfig } from './config';

describe('SecurityTypesConfig', () => {
  const config = SecurityTypesConfig(vi.fn(), [], {}, {}, 'ServiceAccount');

  it('exposes the documented security resource kinds', () => {
    expect(Object.keys(config)).toEqual([
      'ServiceAccount',
      'ClusterRole',
      'Role',
      'ClusterRoleBinding',
      'RoleBinding',
    ]);
  });

  it('each kind has a name and columns array', () => {
    for (const [key, cfg] of Object.entries(config)) {
      expect(cfg.name).toBe(key);
      expect(Array.isArray(cfg.columns)).toBe(true);
      expect(cfg.columns.length).toBeGreaterThan(0);
    }
  });

  it('exposes the canonical ClusterRole columns including namespace', () => {
    const names = config.ClusterRole.columns.map((c: any) => c.name);
    expect(names).toEqual([
      'id',
      'metadata.name',
      'apiVersion',
      'metadata.namespace',
      'cluster_id',
      'metadata.creationTimestamp',
    ]);
  });
});
