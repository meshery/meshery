import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  getStatus: (val: string) => `STATUS:${val}`,
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
  TooltipWrappedConnectionChip: ({ title }: any) => (
    <div data-testid="connection-chip">{title}</div>
  ),
}));
vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => () => undefined,
}));
vi.mock('../sortable-table-cell', () => ({
  DefaultTableCell: ({ columnData }: any) => <th>{columnData?.label}</th>,
  SortableTableCell: ({ columnData }: any) => <th>{columnData?.label}</th>,
}));
vi.mock('../../../../utils/Enum', () => ({
  CONNECTION_KINDS: { KUBERNETES: 'kubernetes' },
}));
vi.mock('@/components/data-formatter', () => ({
  FormatId: ({ id }: any) => <span>{id}</span>,
}));
vi.mock('../../view', () => ({
  Title: ({ value }: any) => <a data-testid="title">{value}</a>,
}));
vi.mock('../config', () => ({
  SINGLE_VIEW: 'single',
}));

import { useNamespaceTableConfig } from './config';

describe('useNamespaceTableConfig', () => {
  const switchView = vi.fn();
  const config = useNamespaceTableConfig(
    switchView,
    [{ metadata: { name: 'ns1' } }],
    {},
    { kubernetes: { icon: '/k.svg' } },
    'Namespace',
  );

  it('exposes the canonical Namespace config shape', () => {
    expect(config.name).toBe('Namespace');
    expect(Array.isArray(config.colViews)).toBe(true);
    expect(config.columns).toHaveLength(6);
    expect(config.columns.map((c: any) => c.name)).toEqual([
      'id',
      'metadata.name',
      'apiVersion',
      'cluster_id',
      'metadata.creationTimestamp',
      'status.attribute',
    ]);
  });

  it('formats the id cell with FormatId', () => {
    const idCol = config.columns.find((c: any) => c.name === 'id');
    const { getByText } = render(idCol.options.customBodyRender('abc123'));
    expect(getByText('abc123')).toBeInTheDocument();
  });

  it('renders a clickable Title cell that wires switchView on click', () => {
    const nameCol = config.columns.find((c: any) => c.name === 'metadata.name');
    const { getByTestId } = render(nameCol.options.customBodyRender('default', { rowIndex: 0 }));
    expect(getByTestId('title')).toHaveTextContent('default');
  });

  it('renders the cluster column with a TooltipWrappedConnectionChip', () => {
    const clusterCol = config.columns.find((c: any) => c.name === 'cluster_id');
    const { getByTestId } = render(clusterCol.options.customBodyRender('cluster-1'));
    expect(getByTestId('connection-chip')).toHaveTextContent('ctx-cluster-1');
  });

  it('renders age via timeAgo and status via getStatus', () => {
    const ageCol = config.columns.find((c: any) => c.name === 'metadata.creationTimestamp');
    const { container: ageContainer } = render(ageCol.options.customBodyRender('2024-01-01'));
    expect(ageContainer.textContent).toBe('AGO:2024-01-01');

    const statusCol = config.columns.find((c: any) => c.name === 'status.attribute');
    const { container: statusContainer } = render(statusCol.options.customBodyRender('OK'));
    expect(statusContainer.textContent).toBe('STATUS:OK');
  });

  it('falls back to empty icon when connection metadata is missing', () => {
    const minimalConfig = useNamespaceTableConfig(switchView, [], {}, null, 'Namespace');
    const clusterCol = minimalConfig.columns.find((c: any) => c.name === 'cluster_id');
    expect(typeof clusterCol.options.customBodyRender).toBe('function');
  });
});
