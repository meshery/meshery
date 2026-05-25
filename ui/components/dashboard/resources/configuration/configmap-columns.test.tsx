import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  timeAgo: (val: string) => `AGO:${val}`,
}));
vi.mock('../../../../utils/multi-ctx', () => ({
  getK8sContextFromClusterId: (id: string) => ({ name: `ctx-${id}` }),
}));
vi.mock('../../../connections/ConnectionChip', () => ({
  TooltipWrappedConnectionChip: ({ title }: any) => <div data-testid="chip">{title}</div>,
}));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: () => null,
  SortableTableCell: () => null,
}));
vi.mock('../../../../utils/Enum', () => ({
  CONNECTION_KINDS: { KUBERNETES: 'kubernetes' },
}));
vi.mock('@/components/data-formatter', () => ({
  FormatId: ({ id }: any) => <span>{id}</span>,
}));
vi.mock('../../view', () => ({
  Title: ({ value }: any) => <a>{value}</a>,
}));
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { buildConfigMapColumns } from './configmap-columns';

describe('buildConfigMapColumns', () => {
  const switchView = vi.fn();
  const ping = vi.fn();
  const args = {
    switchView,
    meshSyncResources: [{ metadata: { name: 'cm-1' } }],
    k8sConfig: {},
    connectionMetadataState: { kubernetes: { icon: '/k.svg' } },
    workloadType: 'ConfigMap',
    ping,
  };
  const config = buildConfigMapColumns(args);

  it('returns the canonical ConfigMap column set', () => {
    expect(config.name).toBe('ConfigMap');
    expect(config.columns.map((c: any) => c.name)).toEqual([
      'id',
      'metadata.name',
      'apiVersion',
      'data',
      'metadata.namespace',
      'cluster_id',
      'metadata.creationTimestamp',
    ]);
  });

  it('parses the data column into a comma separated list of keys', () => {
    const dataCol = config.columns.find((c: any) => c.name === 'data');
    const value = JSON.stringify({ alpha: '1', beta: '2' });
    const { container } = render(dataCol.options.customBodyRender(value));
    expect(container.textContent).toBe('alpha, beta');
  });

  it('renders a dash when the data column is empty', () => {
    const dataCol = config.columns.find((c: any) => c.name === 'data');
    const { container } = render(dataCol.options.customBodyRender(undefined));
    expect(container.textContent).toBe('-');
  });

  it('renders the cluster chip via the multi-ctx helper', () => {
    const clusterCol = config.columns.find((c: any) => c.name === 'cluster_id');
    const { getByTestId } = render(clusterCol.options.customBodyRender('cluster-1'));
    expect(getByTestId('chip')).toHaveTextContent('ctx-cluster-1');
  });

  it('renders age via timeAgo', () => {
    const ageCol = config.columns.find((c: any) => c.name === 'metadata.creationTimestamp');
    const { container } = render(ageCol.options.customBodyRender('2024-01-01'));
    expect(container.textContent).toBe('AGO:2024-01-01');
  });
});
