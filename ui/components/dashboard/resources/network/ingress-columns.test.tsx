import { describe, expect, it, vi } from 'vitest';
import { render } from '@testing-library/react';

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

const getColumn = (label: string) => {
  const column = buildIngressColumns({
    switchView: () => {},
    meshSyncResources: [],
    k8sConfig: {},
    connectionMetadataState: {},
    workloadType: 'Ingress',
    ping: () => {},
  }).columns.find((col: any) => col.label === label) as any;
  expect(column).toBeDefined();
  return column;
};

const renderCell = (label: string, value: unknown) => {
  const { container } = render(
    <>{getColumn(label).options.customBodyRender(JSON.stringify(value), {})}</>,
  );
  return container.textContent;
};

describe('Ingress LoadBalancer column', () => {
  it('shows the hostname when the load balancer has no IP', () => {
    expect(
      renderCell('LoadBalancer IP', {
        loadBalancer: { ingress: [{ hostname: 'lb.example.com' }] },
      }),
    ).toBe('lb.example.com');
  });

  it('shows every address when there are several', () => {
    expect(
      renderCell('LoadBalancer IP', {
        loadBalancer: { ingress: [{ ip: '203.0.113.10' }, { hostname: 'lb.example.com' }] },
      }),
    ).toBe('203.0.113.10, lb.example.com');
  });

  it('leaves out an entry that has no address', () => {
    expect(
      renderCell('LoadBalancer IP', { loadBalancer: { ingress: [{}, { ip: '203.0.113.10' }] } }),
    ).toBe('203.0.113.10');
  });

  it.each([null, {}, { loadBalancer: {} }, { loadBalancer: { ingress: [{}] } }])(
    'is empty when there is no address: %j',
    (status) => {
      expect(renderCell('LoadBalancer IP', status)).toBe('');
    },
  );
});

describe('Ingress Rules column', () => {
  it('shows * for a rule without a host', () => {
    expect(renderCell('Rules', { rules: [{ http: { paths: [] } }] })).toBe('*');
  });

  it('shows the host of every rule', () => {
    expect(renderCell('Rules', { rules: [{ host: 'a.example.com' }, { http: {} }] })).toBe(
      'a.example.com, *',
    );
  });

  it('is empty when there are no rules', () => {
    expect(renderCell('Rules', { defaultBackend: {} })).toBe('');
  });
});
