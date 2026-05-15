import React from 'react';
import { render } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';

vi.mock('../../../../utils/k8s-utils', () => ({
  getResourceStr: (val: string) => `mem-${val}`,
  resourceParsers: { memory: (v: string) => `parsed:${v}` },
  timeAgo: (val: string) => `AGO:${val}`,
  getStatus: (val: string) => `STATUS:${val}`,
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
vi.mock('../config', () => ({ SINGLE_VIEW: 'single' }));

import { NodeTableConfig } from './config';

describe('NodeTableConfig', () => {
  const config = NodeTableConfig(vi.fn(), [{}], {}, {}, 'Node');

  it('describes the Node columns with the documented column names', () => {
    expect(config.name).toBe('Node');
    expect(config.columns.map((c: any) => c.name)).toEqual([
      'id',
      'metadata.name',
      'apiVersion',
      'status.attribute',
      'status.attribute',
      'cluster_id',
      'status.attribute',
      'status.attribute',
      'metadata.creationTimestamp',
      'status.attribute',
    ]);
  });

  it('parses CPU/Memory from status.attribute', () => {
    const cpuCol = config.columns[3];
    const { container } = render(
      cpuCol.options.customBodyRender(JSON.stringify({ capacity: { cpu: '4', memory: '8Gi' } })),
    );
    expect(container.textContent).toBe('4');

    const memCol = config.columns[4];
    const { container: memContainer } = render(
      memCol.options.customBodyRender(JSON.stringify({ capacity: { memory: '8Gi' } })),
    );
    expect(memContainer.textContent).toBe('mem-parsed:8Gi');
  });

  it('renders Internal and External IP from status attributes', () => {
    const internalIpCol = config.columns[6];
    const { container: intCtr } = render(
      internalIpCol.options.customBodyRender(
        JSON.stringify({ addresses: [{ type: 'InternalIP', address: '10.0.0.1' }] }),
      ),
    );
    expect(intCtr.textContent).toBe('10.0.0.1');

    const externalIpCol = config.columns[7];
    const { container: extCtr } = render(
      externalIpCol.options.customBodyRender(
        JSON.stringify({ addresses: [{ type: 'ExternalIP', address: '34.1.2.3' }] }),
      ),
    );
    expect(extCtr.textContent).toBe('34.1.2.3');
  });

  it('renders a dash for missing External IP', () => {
    const externalIpCol = config.columns[7];
    const { container } = render(
      externalIpCol.options.customBodyRender(JSON.stringify({ addresses: [] })),
    );
    expect(container.textContent).toContain('-');
  });

  it('renders the conditions cell using getStatus output when present', () => {
    const conditionsCol = config.columns[9];
    const { container } = render(conditionsCol.options.customBodyRender('Ready'));
    expect(container.textContent).toBe('STATUS:Ready');
  });
});
