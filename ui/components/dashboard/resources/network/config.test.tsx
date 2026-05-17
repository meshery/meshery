import { describe, expect, it, vi } from 'vitest';

vi.mock('@/utils/hooks/useKubernetesHook', () => ({
  default: () => () => undefined,
}));

vi.mock('./service-columns', () => ({
  buildServiceColumns: () => ({ name: 'Service', columns: [] }),
}));
vi.mock('./endpoints-columns', () => ({
  buildEndpointsColumns: () => ({ name: 'Endpoints', columns: [] }),
}));
vi.mock('./endpoint-slice-columns', () => ({
  buildEndpointSliceColumns: () => ({ name: 'EndpointSlice', columns: [] }),
}));
vi.mock('./ingress-columns', () => ({
  buildIngressColumns: () => ({ name: 'Ingress', columns: [] }),
}));
vi.mock('./ingress-class-columns', () => ({
  buildIngressClassColumns: () => ({ name: 'IngressClass', columns: [] }),
}));
vi.mock('./network-policy-columns', () => ({
  buildNetworkPolicyColumns: () => ({ name: 'NetworkPolicy', columns: [] }),
}));

import { NetWorkTableConfig } from './config';

describe('NetWorkTableConfig', () => {
  it('exposes the documented network workload configs', () => {
    const result = NetWorkTableConfig(vi.fn(), [], {}, {}, 'Service');
    expect(Object.keys(result)).toEqual([
      'Service',
      'Endpoints',
      'EndpointSlice',
      'Ingress',
      'IngressClass',
      'NetworkPolicy',
    ]);
    expect(result.Service.name).toBe('Service');
  });
});
