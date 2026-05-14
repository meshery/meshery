import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';

vi.mock('react-relay', () => ({
  requestSubscription: vi.fn(),
}));

vi.mock('../relayEnvironment', () => ({
  createRelayEnvironment: vi.fn(() => ({ __mockRelayEnv: true })),
}));

import { createSubscription } from '../subscriptionHelper';
import { requestSubscription } from 'react-relay';
import { createRelayEnvironment } from '../relayEnvironment';

describe('createSubscription', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('passes subscription, variables, onNext, and onError to requestSubscription', () => {
    const subscription = { name: 'TestSub' };
    const variables = { id: '1' };
    const onNext = vi.fn();
    const onError = vi.fn();

    createSubscription({
      subscription,
      variables,
      onNext,
      onError,
      subscriptionName: 'MySub',
    });

    expect(requestSubscription).toHaveBeenCalledTimes(1);
    const [env, config] = (requestSubscription as ReturnType<typeof vi.fn>).mock.calls[0];
    expect(env).toEqual({ __mockRelayEnv: true });
    expect(config.subscription).toBe(subscription);
    expect(config.variables).toBe(variables);
    expect(config.onNext).toBe(onNext);
    expect(config.onError).toBe(onError);
  });

  it('defaults variables to {} when not provided', () => {
    createSubscription({
      subscription: { name: 'NoVars' },
      onNext: vi.fn(),
    });
    const [, config] = (requestSubscription as ReturnType<typeof vi.fn>).mock.calls[0];
    expect(config.variables).toEqual({});
  });

  it('uses default onError which logs errors via console.error', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    createSubscription({
      subscription: { name: 'TestSub' },
      onNext: vi.fn(),
      subscriptionName: 'MySub',
    });

    const [, config] = (requestSubscription as ReturnType<typeof vi.fn>).mock.calls[0];
    const err = new Error('boom');
    config.onError(err);

    expect(consoleSpy).toHaveBeenCalledWith(
      expect.stringContaining('[GraphQL Subscription: MySub] Error:'),
      err,
    );
  });

  it('uses default subscription name "Unknown" when not supplied', () => {
    const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    createSubscription({
      subscription: { name: 'TestSub' },
      onNext: vi.fn(),
    });
    const [, config] = (requestSubscription as ReturnType<typeof vi.fn>).mock.calls[0];
    config.onError(new Error('x'));
    expect(consoleSpy.mock.calls[0][0]).toContain('Unknown');
  });

  it('prefers user-provided onError over the default', () => {
    const customOnError = vi.fn();
    createSubscription({
      subscription: { name: 'TestSub' },
      onNext: vi.fn(),
      onError: customOnError,
    });
    const [, config] = (requestSubscription as ReturnType<typeof vi.fn>).mock.calls[0];
    expect(config.onError).toBe(customOnError);
  });

  it('calls createRelayEnvironment with an empty config', () => {
    createSubscription({
      subscription: { name: 'TestSub' },
      onNext: vi.fn(),
    });
    expect(createRelayEnvironment).toHaveBeenCalledWith({});
  });

  it('returns whatever requestSubscription returns', () => {
    const disposable = { dispose: vi.fn() };
    (requestSubscription as ReturnType<typeof vi.fn>).mockReturnValueOnce(disposable);
    const result = createSubscription({
      subscription: { name: 'TestSub' },
      onNext: vi.fn(),
    });
    expect(result).toBe(disposable);
  });
});
