import { describe, expect, it, vi } from 'vitest';

vi.mock('@/graphql/subscriptions/MesheryControllersStatusSubscription', () => ({
  default: vi.fn(),
}));

import {
  mergeMeshSyncSubscription,
  mergeOperatorStateSubscription,
  mergeMesheryController,
} from './mergeFns';

describe('mergeMeshSyncSubscription', () => {
  it('returns a fresh single-element list when no prior state exists', () => {
    const newData = { contextID: 'ctx-1' };
    expect(mergeMeshSyncSubscription(undefined, newData)).toEqual([newData]);
    expect(mergeMeshSyncSubscription(null as any, newData)).toEqual([newData]);
  });

  it('replaces any existing entry with the same contextID', () => {
    const prior = [
      { contextID: 'a', value: 1 },
      { contextID: 'b', value: 2 },
    ];
    const result = mergeMeshSyncSubscription(prior, { contextID: 'a', value: 99 });
    expect(result).toEqual([
      { contextID: 'b', value: 2 },
      { contextID: 'a', value: 99 },
    ]);
  });

  it('appends new entries when their contextID is new', () => {
    const prior = [{ contextID: 'a', value: 1 }];
    const result = mergeMeshSyncSubscription(prior, { contextID: 'b', value: 2 });
    expect(result).toEqual([
      { contextID: 'a', value: 1 },
      { contextID: 'b', value: 2 },
    ]);
  });
});

describe('mergeOperatorStateSubscription', () => {
  it('is functionally identical to mergeMeshSyncSubscription', () => {
    const prior = [{ contextID: 'a', value: 1 }];
    expect(mergeOperatorStateSubscription(prior, { contextID: 'a', value: 5 })).toEqual([
      { contextID: 'a', value: 5 },
    ]);
  });
});

describe('mergeMesheryController', () => {
  it('returns newData as-is when there is no prior state', () => {
    const newData = [{ contextId: '1', controller: 'meshsync' }];
    expect(mergeMesheryController(undefined, newData)).toBe(newData);
  });

  it('keeps prior controllers that do NOT match anything in newData, then appends newData', () => {
    const prior = [
      { contextId: 'a', controller: 'meshsync' },
      { contextId: 'a', controller: 'broker' },
    ];
    const newData = [{ contextId: 'a', controller: 'meshsync', status: 'running' }];

    const result = mergeMesheryController(prior, newData);
    expect(result).toEqual([
      { contextId: 'a', controller: 'broker' },
      { contextId: 'a', controller: 'meshsync', status: 'running' },
    ]);
  });
});
