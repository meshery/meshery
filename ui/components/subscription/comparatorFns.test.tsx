import { describe, expect, it, vi } from 'vitest';

vi.mock('@/graphql/subscriptions/MesheryControllersStatusSubscription', () => ({
  default: vi.fn(),
}));

import {
  isMeshSyncSubscriptionDataUpdated,
  isOperatorStateSubscriptionDataUpdated,
  isMesheryControllerStateSubscriptionDataUpdated,
} from './comparatorFns';

describe('isMeshSyncSubscriptionDataUpdated', () => {
  it('returns true when there is no prior state', () => {
    expect(isMeshSyncSubscriptionDataUpdated(undefined, { contextID: 'a' })).toBe(true);
    expect(isMeshSyncSubscriptionDataUpdated(null as any, { contextID: 'a' })).toBe(true);
  });

  it('returns true when contextID is new', () => {
    const prior = [{ contextID: 'a', value: 1 }];
    expect(isMeshSyncSubscriptionDataUpdated(prior, { contextID: 'b', value: 1 })).toBe(true);
  });

  it('returns false when same contextID and equal content', () => {
    const prior = [{ contextID: 'a', value: 1 }];
    expect(isMeshSyncSubscriptionDataUpdated(prior, { contextID: 'a', value: 1 })).toBe(false);
  });

  it('returns true when same contextID but content differs', () => {
    const prior = [{ contextID: 'a', value: 1 }];
    expect(isMeshSyncSubscriptionDataUpdated(prior, { contextID: 'a', value: 2 })).toBe(true);
  });
});

describe('isOperatorStateSubscriptionDataUpdated', () => {
  it('behaves identically to the MeshSync variant', () => {
    expect(isOperatorStateSubscriptionDataUpdated(undefined, { contextID: 'x' })).toBe(true);
    const prior = [{ contextID: 'x', value: 1 }];
    expect(isOperatorStateSubscriptionDataUpdated(prior, { contextID: 'x', value: 1 })).toBe(false);
  });
});

describe('isMesheryControllerStateSubscriptionDataUpdated', () => {
  it('returns true when there is no prior state', () => {
    expect(
      isMesheryControllerStateSubscriptionDataUpdated(undefined, [
        { contextId: '1', controller: 'meshsync' },
      ]),
    ).toBe(true);
  });

  it('returns true when prior content does not match newData object structure', () => {
    const prior = [{ contextId: '1', controller: 'meshsync', status: 'ok' }];
    const newData = [{ contextId: '1', controller: 'meshsync', status: 'down' }];
    expect(isMesheryControllerStateSubscriptionDataUpdated(prior, newData)).toBe(true);
  });

  it('returns false when prior and new data are identical', () => {
    const item = { contextId: '1', controller: 'meshsync' };
    expect(isMesheryControllerStateSubscriptionDataUpdated([item], [item])).toBe(false);
  });
});
