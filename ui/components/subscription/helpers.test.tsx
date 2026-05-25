import { describe, expect, it, vi } from 'vitest';

vi.mock('@/graphql/subscriptions/MesheryControllersStatusSubscription', () => ({
  default: vi.fn(),
}));

import { isControllerObjectEqual, fnMapping, MESHERY_CONTROLLER_SUBSCRIPTION } from './helpers';

describe('isControllerObjectEqual', () => {
  it('returns true only when both contextId AND controller match', () => {
    expect(
      isControllerObjectEqual(
        { contextId: 'a', controller: 'broker' },
        { contextId: 'a', controller: 'broker' },
      ),
    ).toBe(true);
  });

  it('returns false when contextId differs', () => {
    expect(
      isControllerObjectEqual(
        { contextId: 'a', controller: 'broker' },
        { contextId: 'b', controller: 'broker' },
      ),
    ).toBe(false);
  });

  it('returns false when controller differs', () => {
    expect(
      isControllerObjectEqual(
        { contextId: 'a', controller: 'broker' },
        { contextId: 'a', controller: 'meshsync' },
      ),
    ).toBe(false);
  });
});

describe('fnMapping', () => {
  it('exposes a MESHERY_CONTROLLER_SUBSCRIPTION entry with required handlers', () => {
    expect(MESHERY_CONTROLLER_SUBSCRIPTION).toBe('MESHERY_CONTROLLER_SUBSCRIPTION');
    const entry = fnMapping[MESHERY_CONTROLLER_SUBSCRIPTION];
    expect(entry).toBeDefined();
    expect(entry.eventName).toBe('subscribeMesheryControllersStatus');
    expect(typeof entry.subscriptionFn).toBe('function');
    expect(typeof entry.mergeFn).toBe('function');
    expect(typeof entry.comparatorFn).toBe('function');
  });
});
