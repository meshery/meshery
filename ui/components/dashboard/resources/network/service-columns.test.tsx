import { describe, expect, it, vi } from 'vitest';
import { render, screen } from '@testing-library/react';

vi.mock('../../../../utils/k8s-utils', () => ({
  getStatus: (v: any) => `STATUS:${v}`,
  timeAgo: (v: string) => `AGO:${v}`,
}));
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

import { buildServiceColumns } from './service-columns';

describe('buildServiceColumns', () => {
  it('returns the Service name and a non-empty columns array', () => {
    const cfg = buildServiceColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Service',
      ping: () => {},
    });
    expect(cfg.name).toBe('Service');
    expect(cfg.columns.length).toBeGreaterThan(0);
  });

  const renderExternalIP = (attribute: unknown) => {
    const cfg = buildServiceColumns({
      switchView: () => {},
      meshSyncResources: [],
      k8sConfig: {},
      connectionMetadataState: {},
      workloadType: 'Service',
      ping: () => {},
    });
    const column = cfg.columns.find((c) => c.label === 'External IP');
    return render(<>{column?.options.customBodyRender(JSON.stringify(attribute))}</>);
  };

  it('shows load balancer ingress IPs and hostnames from status.loadBalancer', () => {
    const { container } = renderExternalIP({
      loadBalancer: { ingress: [{ ip: '203.0.113.10' }, { hostname: 'abc.elb.amazonaws.com' }] },
    });
    expect(container).toHaveTextContent('203.0.113.10,abc.elb.amazonaws.com');
  });

  it('renders nothing when the Service has no load balancer ingress', () => {
    const { container } = renderExternalIP({ loadBalancer: {} });
    expect(container).toBeEmptyDOMElement();
  });
});
