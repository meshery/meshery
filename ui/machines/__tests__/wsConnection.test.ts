import { describe, expect, it, vi, beforeEach } from 'vitest';
import { createActor } from 'xstate';

// The wsConnection machine spawns an actor that subscribes to the
// connection-status event bus exposed by `lib/wsStatus`. We mock that
// module so the test can synthetically fire `connected | closed | error`
// events and assert the machine's resulting state transitions.

type Handler = (...args: unknown[]) => void;

const hoisted = vi.hoisted(() => {
  const handlers: Record<string, Handler[]> = {};
  const offFns: Record<string, ReturnType<typeof vi.fn>[]> = {};
  const on = (eventName: string, handler: Handler) => {
    handlers[eventName] = handlers[eventName] || [];
    handlers[eventName].push(handler);
    const off = vi.fn();
    offFns[eventName] = offFns[eventName] || [];
    offFns[eventName].push(off);
    return off;
  };
  return { handlers, offFns, on };
});
const handlers = hoisted.handlers;
const offFns = hoisted.offFns;

vi.mock('../../lib/wsStatus', () => ({
  on: hoisted.on,
  off: vi.fn(),
}));

import { wsConnectionMachine, WS_CONNECTION_EVENTS } from '../wsConnection';

const fire = (eventName: string, ...args: unknown[]) => {
  handlers[eventName]?.forEach((h) => h(...args));
};

const resetHandlers = () => {
  for (const k of Object.keys(handlers)) delete handlers[k];
  for (const k of Object.keys(offFns)) delete offFns[k];
};

describe('wsConnectionMachine', () => {
  beforeEach(() => {
    resetHandlers();
  });

  it('boots in the connecting state with a fresh context', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('connecting');
    expect(snap.context.retryCount).toBe(0);
    expect(snap.context.lastError).toBeNull();
    expect(snap.context.lastConnectedAt).toBeNull();
    actor.stop();
  });

  it('moves connecting -> connected on CONNECTED and resets retry context', () => {
    const before = Date.now();
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('connected');
    expect(snap.context.retryCount).toBe(0);
    expect(typeof snap.context.lastConnectedAt).toBe('number');
    expect(snap.context.lastConnectedAt as number).toBeGreaterThanOrEqual(before);
    actor.stop();
  });

  it('moves connecting -> disconnected on ERROR and increments retryCount', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.ERROR, data: { error: new Error('boom') } });
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('disconnected');
    expect(snap.context.retryCount).toBe(1);
    actor.stop();
  });

  it('moves connected -> reconnecting on DISCONNECTED and increments retry on entry', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    actor.send({ type: WS_CONNECTION_EVENTS.DISCONNECTED });
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('reconnecting');
    expect(snap.context.retryCount).toBe(1);
    actor.stop();
  });

  it('moves connected -> reconnecting on ERROR', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    actor.send({ type: WS_CONNECTION_EVENTS.ERROR });
    expect(actor.getSnapshot().value).toBe('reconnecting');
    actor.stop();
  });

  it('reconnecting -> connected on CONNECTED with reset retry', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    actor.send({ type: WS_CONNECTION_EVENTS.DISCONNECTED });
    expect(actor.getSnapshot().value).toBe('reconnecting');

    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('connected');
    expect(snap.context.retryCount).toBe(0);
    actor.stop();
  });

  it('reconnecting -> disconnected on ERROR', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    actor.send({ type: WS_CONNECTION_EVENTS.DISCONNECTED });
    actor.send({ type: WS_CONNECTION_EVENTS.ERROR });
    expect(actor.getSnapshot().value).toBe('disconnected');
    actor.stop();
  });

  it('disconnected -> connected on CONNECTED with reset retry', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.send({ type: WS_CONNECTION_EVENTS.ERROR });
    expect(actor.getSnapshot().value).toBe('disconnected');
    actor.send({ type: WS_CONNECTION_EVENTS.CONNECTED });
    expect(actor.getSnapshot().value).toBe('connected');
    expect(actor.getSnapshot().context.retryCount).toBe(0);
    actor.stop();
  });

  it('routes subscriptionClient events into the machine via the actor', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    expect(actor.getSnapshot().value).toBe('connecting');

    fire('connected');
    expect(actor.getSnapshot().value).toBe('connected');

    fire('closed', { code: 1006 });
    expect(actor.getSnapshot().value).toBe('reconnecting');

    fire('error', new Error('socket'));
    expect(actor.getSnapshot().value).toBe('disconnected');

    actor.stop();
  });

  it('unsubscribes from all three subscriptionClient events on stop', () => {
    const actor = createActor(wsConnectionMachine);
    actor.start();
    actor.stop();

    // Each topic should have its unsubscribe callback invoked
    expect(offFns.connected?.[0]).toHaveBeenCalled();
    expect(offFns.closed?.[0]).toHaveBeenCalled();
    expect(offFns.error?.[0]).toHaveBeenCalled();
  });

  it('exposes the canonical WS_CONNECTION_EVENTS constants', () => {
    expect(WS_CONNECTION_EVENTS).toEqual({
      CONNECTED: 'WS_CONNECTED',
      DISCONNECTED: 'WS_DISCONNECTED',
      ERROR: 'WS_ERROR',
      RECONNECTING: 'WS_RECONNECTING',
    });
  });
});
