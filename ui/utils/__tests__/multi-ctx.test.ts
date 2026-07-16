import { beforeEach, describe, expect, it } from 'vitest';
import {
  ctxUrl,
  loadSelectedK8sContexts,
  persistSelectedK8sContexts,
  getK8sClusterIdsFromCtxId,
  getFirstCtxIdFromSelectedCtxIds,
  getK8sConfigIdsFromK8sConfig,
  getK8sClusterNamesFromCtxId,
  getK8sContextFromClusterId,
  getClusterNameFromClusterId,
  getClusterNameFromConnectionId,
  getConnectionIdFromClusterId,
  getClusterNameFromCtxId,
  getConnectionIDsFromContextIds,
  getControllerPollConnectionIDsFromContextIds,
} from '../multi-ctx';

const sampleConfig = [
  { id: 'ctx-1', name: 'minikube', kubernetesServerId: 'srv-1', connectionId: 'conn-1' },
  { id: 'ctx-2', name: 'gke-prod', kubernetesServerId: 'srv-2', connectionId: 'conn-2' },
  { id: 'ctx-3', name: 'eks-stage', kubernetesServerId: 'srv-3', connectionId: 'conn-3' },
];

describe('ctxUrl', () => {
  it('returns the URL unchanged when no contexts are supplied', () => {
    expect(ctxUrl('/api/foo', [])).toBe('/api/foo');
    expect(ctxUrl('/api/foo', null)).toBe('/api/foo');
    expect(ctxUrl('/api/foo', undefined)).toBe('/api/foo');
  });

  it('appends each context as its own query parameter joined by &', () => {
    expect(ctxUrl('/api/foo', ['a', 'b'])).toBe('/api/foo?contexts=a&contexts=b');
  });

  it('handles a single context', () => {
    expect(ctxUrl('/api/foo', ['only'])).toBe('/api/foo?contexts=only');
  });
});

describe('getK8sClusterIdsFromCtxId', () => {
  it('returns an empty array for empty or missing inputs', () => {
    expect(getK8sClusterIdsFromCtxId([], sampleConfig)).toEqual([]);
    expect(getK8sClusterIdsFromCtxId(null, sampleConfig)).toEqual([]);
    expect(getK8sClusterIdsFromCtxId(['ctx-1'], null)).toEqual([]);
  });

  it('returns every cluster id when "all" is selected', () => {
    expect(getK8sClusterIdsFromCtxId(['all'], sampleConfig)).toEqual(['srv-1', 'srv-2', 'srv-3']);
  });

  it('returns the cluster ids for the selected context ids', () => {
    expect(getK8sClusterIdsFromCtxId(['ctx-1', 'ctx-3'], sampleConfig)).toEqual(['srv-1', 'srv-3']);
  });

  it('skips contexts not found in the config', () => {
    expect(getK8sClusterIdsFromCtxId(['ctx-1', 'missing'], sampleConfig)).toEqual(['srv-1']);
  });
});

describe('getFirstCtxIdFromSelectedCtxIds', () => {
  it('returns an empty string when no contexts are selected', () => {
    expect(getFirstCtxIdFromSelectedCtxIds([], sampleConfig)).toBe('');
    expect(getFirstCtxIdFromSelectedCtxIds(undefined, sampleConfig)).toBe('');
  });

  it('returns the first config id when "all" is selected', () => {
    expect(getFirstCtxIdFromSelectedCtxIds(['all'], sampleConfig)).toBe('ctx-1');
  });

  it('returns the first selected context otherwise', () => {
    expect(getFirstCtxIdFromSelectedCtxIds(['ctx-2', 'ctx-3'], sampleConfig)).toBe('ctx-2');
  });
});

describe('getK8sConfigIdsFromK8sConfig', () => {
  it('returns all config ids', () => {
    expect(getK8sConfigIdsFromK8sConfig(sampleConfig)).toEqual(['ctx-1', 'ctx-2', 'ctx-3']);
  });

  it('returns an empty array for empty/missing config', () => {
    expect(getK8sConfigIdsFromK8sConfig([])).toEqual([]);
    expect(getK8sConfigIdsFromK8sConfig(null)).toEqual([]);
    expect(getK8sConfigIdsFromK8sConfig(undefined)).toEqual([]);
  });
});

describe('getK8sClusterNamesFromCtxId', () => {
  it('returns an empty array when no contexts are selected', () => {
    expect(getK8sClusterNamesFromCtxId([], sampleConfig)).toEqual([]);
  });

  it('returns ["all"] when "all" is included in selection', () => {
    expect(getK8sClusterNamesFromCtxId(['all'], sampleConfig)).toEqual(['all']);
  });

  it('returns the names of the matched contexts', () => {
    expect(getK8sClusterNamesFromCtxId(['ctx-1', 'ctx-3'], sampleConfig)).toEqual([
      'minikube',
      'eks-stage',
    ]);
  });

  it('skips unmatched context ids', () => {
    expect(getK8sClusterNamesFromCtxId(['ctx-1', 'bogus'], sampleConfig)).toEqual(['minikube']);
  });
});

describe('getK8sContextFromClusterId', () => {
  it('returns the matching cluster config object', () => {
    expect(getK8sContextFromClusterId('srv-2', sampleConfig)).toBe(sampleConfig[1]);
  });

  it('returns an empty object when no match is found', () => {
    expect(getK8sContextFromClusterId('does-not-exist', sampleConfig)).toEqual({});
  });
});

describe('getClusterNameFromClusterId', () => {
  it('returns the cluster name for a matching server id', () => {
    expect(getClusterNameFromClusterId('srv-1', sampleConfig)).toBe('minikube');
  });

  it('returns an empty string when no match is found', () => {
    expect(getClusterNameFromClusterId('missing', sampleConfig)).toBe('');
  });
});

describe('getClusterNameFromConnectionId', () => {
  it('returns the cluster name for a matching connection id', () => {
    expect(getClusterNameFromConnectionId('conn-2', sampleConfig)).toBe('gke-prod');
  });

  it('returns an empty string when no match is found', () => {
    expect(getClusterNameFromConnectionId('nope', sampleConfig)).toBe('');
  });
});

describe('getConnectionIdFromClusterId', () => {
  it('returns the connection id for a matching server id', () => {
    expect(getConnectionIdFromClusterId('srv-3', sampleConfig)).toBe('conn-3');
  });

  it('returns an empty string when no match is found', () => {
    expect(getConnectionIdFromClusterId('missing', sampleConfig)).toBe('');
  });
});

describe('getClusterNameFromCtxId', () => {
  it('returns the cluster name for a matching context id', () => {
    expect(getClusterNameFromCtxId('ctx-1', sampleConfig)).toBe('minikube');
  });

  it('returns an empty string when no match is found', () => {
    expect(getClusterNameFromCtxId('missing', sampleConfig)).toBe('');
  });
});

describe('getConnectionIDsFromContextIds', () => {
  it('returns connection ids whose config id is in the context list', () => {
    expect(getConnectionIDsFromContextIds(['ctx-1', 'ctx-3'], sampleConfig)).toEqual([
      'conn-1',
      'conn-3',
    ]);
  });

  it('returns an empty array when no contexts are matched', () => {
    expect(getConnectionIDsFromContextIds(['missing'], sampleConfig)).toEqual([]);
  });

  it('returns an empty array when contexts list is empty', () => {
    expect(getConnectionIDsFromContextIds([], sampleConfig)).toEqual([]);
  });
});

describe('getControllerPollConnectionIDsFromContextIds', () => {
  const pollConfig = [
    // operator + connected → eligible
    {
      id: 'ctx-1',
      connectionId: 'conn-1',
      meshsync_deployment_mode: 'operator',
      connectionStatus: 'connected',
    },
    // operator + connected (camelCase mode) → eligible
    {
      id: 'ctx-2',
      connectionId: 'conn-2',
      meshsyncDeploymentMode: 'operator',
      connectionStatus: 'connected',
    },
    // operator but not connected → excluded
    {
      id: 'ctx-3',
      connectionId: 'conn-3',
      meshsync_deployment_mode: 'operator',
      connectionStatus: 'registered',
    },
    // embedded but connected → excluded
    {
      id: 'ctx-4',
      connectionId: 'conn-4',
      meshsync_deployment_mode: 'embedded',
      connectionStatus: 'connected',
    },
    // operator + connected, but no status field → excluded (missing = not connected)
    { id: 'ctx-5', connectionId: 'conn-5', meshsync_deployment_mode: 'operator' },
  ];

  it('returns only connections that are operator-mode AND connected', () => {
    expect(
      getControllerPollConnectionIDsFromContextIds(
        ['ctx-1', 'ctx-2', 'ctx-3', 'ctx-4', 'ctx-5'],
        pollConfig,
      ),
    ).toEqual(['conn-1', 'conn-2']);
  });

  it('excludes not-connected, embedded, and status-less connections', () => {
    expect(
      getControllerPollConnectionIDsFromContextIds(['ctx-3', 'ctx-4', 'ctx-5'], pollConfig),
    ).toEqual([]);
  });

  it('respects the context filter', () => {
    expect(getControllerPollConnectionIDsFromContextIds(['ctx-1'], pollConfig)).toEqual(['conn-1']);
  });
});

describe('selected k8s contexts session persistence', () => {
  beforeEach(() => {
    window.sessionStorage.clear();
  });

  it('round-trips a selection through sessionStorage', () => {
    persistSelectedK8sContexts(['ctx-1', 'ctx-2']);
    expect(loadSelectedK8sContexts()).toEqual(['ctx-1', 'ctx-2']);
  });

  it('round-trips the sentinel and empty selections', () => {
    persistSelectedK8sContexts(['all']);
    expect(loadSelectedK8sContexts()).toEqual(['all']);

    persistSelectedK8sContexts([]);
    expect(loadSelectedK8sContexts()).toEqual([]);
  });

  it('returns null when nothing is persisted', () => {
    expect(loadSelectedK8sContexts()).toBeNull();
  });

  it('returns null for corrupt or non-array persisted values', () => {
    window.sessionStorage.setItem('selectedK8sContexts', 'not-json');
    expect(loadSelectedK8sContexts()).toBeNull();

    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify({ ctx: true }));
    expect(loadSelectedK8sContexts()).toBeNull();

    window.sessionStorage.setItem('selectedK8sContexts', JSON.stringify([1, 2]));
    expect(loadSelectedK8sContexts()).toBeNull();
  });

  it('ignores non-array input on persist', () => {
    persistSelectedK8sContexts('all' as unknown as string[]);
    expect(loadSelectedK8sContexts()).toBeNull();
  });
});
