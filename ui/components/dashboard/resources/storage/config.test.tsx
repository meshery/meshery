import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  timeAgo: (val: string) => `AGO:${val}`,
}));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (id: string) => ({
    name: `ctx-${id}`,
    server: 's',
    connectionId: 'c',
  }),
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

import { StorageTableConfig } from './config';

describe('StorageTableConfig', () => {
  const config = StorageTableConfig(vi.fn(), [], {}, {}, 'PersistentVolume');

  it('exposes three storage workload configs', () => {
    expect(Object.keys(config)).toEqual([
      'PersistentVolume',
      'PersistentVolumeClaim',
      'StorageClass',
    ]);
  });

  it('declares the PersistentVolume name and columns', () => {
    expect(config.PersistentVolume.name).toBe('PersistentVolume');
    expect(config.PersistentVolume.columns.length).toBeGreaterThan(0);
  });

  it('declares the PersistentVolumeClaim name and columns', () => {
    expect(config.PersistentVolumeClaim.name).toBe('PersistentVolumeClaim');
    expect(config.PersistentVolumeClaim.columns.length).toBeGreaterThan(0);
  });

  it('declares the StorageClass name and columns', () => {
    expect(config.StorageClass.name).toBe('StorageClass');
    expect(config.StorageClass.columns.length).toBeGreaterThan(0);
  });

  it('formats PV storage class via spec attribute parsing', () => {
    const pvCols = config.PersistentVolume.columns;
    const scCol = pvCols.find((c: any) => c.label === 'Storage Class');
    const { container } = require('@testing-library/react').render(
      scCol.options.customBodyRender(JSON.stringify({ storageClassName: 'gp2' })),
    );
    expect(container.textContent).toBe('gp2');
  });
});
