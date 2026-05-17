import { describe, expect, it, vi, beforeEach } from 'vitest';

const subscriptionFn = vi.fn();
const comparatorFn = vi.fn();
const mergeFn = vi.fn();
const dispose = vi.fn();

vi.mock('./helpers', () => ({
  fnMapping: {
    MESHERY_CONTROLLER_SUBSCRIPTION: {
      eventName: 'subscribeMesheryControllersStatus',
      subscriptionFn: (...args: unknown[]) => {
        subscriptionFn(...args);
        return { dispose };
      },
      mergeFn: (...args: unknown[]) => mergeFn(...args),
      comparatorFn: (...args: unknown[]) => comparatorFn(...args),
    },
  },
}));

import { GQLSubscription } from './subscriptionhandler';

describe('GQLSubscription', () => {
  beforeEach(() => {
    subscriptionFn.mockReset();
    comparatorFn.mockReset();
    mergeFn.mockReset();
    dispose.mockReset();
  });

  it('throws when initialised without a type', () => {
    const sub = new GQLSubscription({
      type: undefined as any,
      connectionIDs: ['c-1'],
      callbackFunction: () => {},
    });
    expect(() => sub.initSubscription()).toThrow(/Subscription Type is Empty/);
  });

  it('initSubscription registers a subscription with the provided callback and connection IDs', () => {
    const callback = vi.fn();
    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: ['c-1', 'c-2'],
      callbackFunction: callback,
    });

    sub.initSubscription();

    expect(subscriptionFn).toHaveBeenCalledTimes(1);
    expect(subscriptionFn.mock.calls[0][1]).toEqual(['c-1', 'c-2']);
  });

  it('defaults connectionIDs to an empty array when none is passed', () => {
    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: undefined as any,
      callbackFunction: () => {},
    });
    expect(sub.connectionIDs).toEqual([]);
  });

  it('doCallback triggers callback when comparatorFn returns true', () => {
    const callback = vi.fn();
    comparatorFn.mockReturnValue(true);
    mergeFn.mockReturnValue([{ contextId: 'a', controller: 'broker' }]);

    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: [],
      callbackFunction: callback,
    });

    sub.doCallback({
      subscribeMesheryControllersStatus: [{ contextId: 'a', controller: 'broker' }],
    });

    expect(comparatorFn).toHaveBeenCalled();
    expect(mergeFn).toHaveBeenCalled();
    expect(callback).toHaveBeenCalledWith([{ contextId: 'a', controller: 'broker' }]);
    expect(sub.state).toEqual([{ contextId: 'a', controller: 'broker' }]);
  });

  it('doCallback does NOT trigger callback when comparatorFn returns false', () => {
    const callback = vi.fn();
    comparatorFn.mockReturnValue(false);

    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: [],
      callbackFunction: callback,
    });
    sub.doCallback({ subscribeMesheryControllersStatus: [] });

    expect(callback).not.toHaveBeenCalled();
    expect(mergeFn).not.toHaveBeenCalled();
  });

  it('flushSubscription disposes of the active subscription', () => {
    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: [],
      callbackFunction: vi.fn(),
    });
    sub.initSubscription();
    sub.flushSubscription();
    expect(dispose).toHaveBeenCalledTimes(1);
  });

  it('updateSubscription disposes prior subscription and re-creates with new IDs', () => {
    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: ['old'],
      callbackFunction: vi.fn(),
    });
    sub.initSubscription();
    expect(subscriptionFn).toHaveBeenCalledTimes(1);

    sub.updateSubscription(['new-1', 'new-2']);

    expect(dispose).toHaveBeenCalled();
    expect(subscriptionFn).toHaveBeenCalledTimes(2);
    expect(subscriptionFn.mock.calls[1][1]).toEqual(['new-1', 'new-2']);
  });

  it('setState updates the cached state', () => {
    const sub = new GQLSubscription({
      type: 'MESHERY_CONTROLLER_SUBSCRIPTION' as any,
      connectionIDs: [],
      callbackFunction: vi.fn(),
    });
    sub.setState({ foo: 'bar' });
    expect(sub.state).toEqual({ foo: 'bar' });
  });
});
