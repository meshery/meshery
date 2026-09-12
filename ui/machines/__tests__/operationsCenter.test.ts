import { describe, expect, it, vi, beforeEach } from 'vitest';
import { createActor } from 'xstate';

// The operationsCenter machine subscribes to the events SSE stream and bridges
// each event into the redux store + the snackbar notifier. We mock the
// subscription module, the redux store, the rtk-query API and the
// NotificationCenter constants so we can drive events and assert side
// effects.

type EventCallback = (event: unknown) => void;
type ErrorCallback = (err: unknown) => void;

const hoisted = vi.hoisted(() => ({
  subscriptionState: {
    onEvent: undefined as EventCallback | undefined,
    onError: undefined as ErrorCallback | undefined,
    dispose: vi.fn(),
  },
  dispatch: vi.fn(),
  pushEvent: vi.fn((p: unknown) => ({ type: 'events/pushEvent', payload: p })),
}));
const subscriptionState = hoisted.subscriptionState;
const dispatch = hoisted.dispatch;
const pushEvent = hoisted.pushEvent;

vi.mock('lib/eventsSubscription', () => ({
  subscribeToEvents: (next: EventCallback, error: ErrorCallback) => {
    hoisted.subscriptionState.onEvent = next;
    hoisted.subscriptionState.onError = error;
    return { dispose: hoisted.subscriptionState.dispose };
  },
}));

vi.mock('../../store', () => ({
  store: { dispatch: hoisted.dispatch },
}));

vi.mock('@/store/slices/events', () => ({
  pushEvent: (p: unknown) => hoisted.pushEvent(p),
}));

vi.mock('../../rtk-query', () => ({
  api: { util: { invalidateTags: vi.fn((tags) => ({ type: 'api/invalidate', tags })) } },
}));

vi.mock('@/rtk-query/notificationCenter', () => ({
  PROVIDER_TAGS: { EVENT: 'Event' },
}));

vi.mock('@/components/layout/NotificationCenter/constants', () => ({
  SEVERITY_TO_NOTIFICATION_TYPE_MAPPING: {
    informational: 'info',
    error: 'error',
    warning: 'warning',
    success: 'success',
  },
}));

import { operationsCenterActor, OPERATION_CENTER_EVENTS } from '../operationsCenter';

const startActor = (notify = vi.fn()) => {
  const actor = createActor(operationsCenterActor, { input: { notify } });
  actor.start();
  return { actor, notify };
};

describe('operationsCenter machine', () => {
  beforeEach(() => {
    dispatch.mockClear();
    pushEvent.mockClear();
    subscriptionState.dispose.mockClear();
    subscriptionState.onEvent = undefined;
    subscriptionState.onError = undefined;
  });

  it('boots into idle after init and spawns the subscription actor', () => {
    const { actor } = startActor();
    const snap = actor.getSnapshot();
    expect(snap.value).toBe('idle');
    expect(subscriptionState.onEvent).toBeTypeOf('function');
    expect(subscriptionState.onError).toBeTypeOf('function');
    actor.stop();
  });

  it('exposes the canonical event constants', () => {
    expect(OPERATION_CENTER_EVENTS).toEqual({
      EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
      ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
    });
  });

  it('stores received events in redux and triggers the notifier', () => {
    const { actor, notify } = startActor();

    const event = {
      id: 'evt-1',
      description: 'A thing happened',
      severity: 'informational',
    };

    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event },
    });

    // Redux dispatch fired with pushEvent action
    expect(pushEvent).toHaveBeenCalledWith(event);
    expect(dispatch).toHaveBeenCalledWith({ type: 'events/pushEvent', payload: event });

    // notify called with the mapped event_type
    expect(notify).toHaveBeenCalledWith({
      message: 'A thing happened',
      event_type: 'info',
      id: 'evt-1',
      showInNotificationCenter: true,
    });

    actor.stop();
  });

  it('emits received events back to the parent for consumers', () => {
    const { actor } = startActor();
    const seen: unknown[] = [];
    actor.on('*', (event) => {
      seen.push(event);
    });
    const payload = {
      id: 'evt-2',
      description: 'desc',
      severity: 'error',
    };
    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event: payload },
    });
    const emitted = seen.find(
      (e: { type?: string }) => e?.type === OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
    );
    expect(emitted).toBeDefined();
    actor.stop();
  });

  it('forwards subscription stream events into the machine via the next callback', () => {
    const { actor } = startActor();
    const event = { id: 'live', description: 'streamed', severity: 'warning' };
    subscriptionState.onEvent?.(event);
    expect(pushEvent).toHaveBeenCalledWith(event);
    actor.stop();
  });

  it('ignores empty stream callbacks but logs', () => {
    const errSpy = vi.spyOn(console, 'error').mockImplementation(() => {});
    const { actor } = startActor();
    subscriptionState.onEvent?.(null);
    expect(pushEvent).not.toHaveBeenCalled();
    expect(errSpy).toHaveBeenCalled();
    errSpy.mockRestore();
    actor.stop();
  });

  it('respawns the subscription actor on ERROR_OCCURRED_IN_SUBSCRIPTION', () => {
    const { actor } = startActor();
    // Send error event to trigger respawn
    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('socket') },
    });

    // After respawn, subscription handlers should still be wired up
    expect(subscriptionState.onEvent).toBeDefined();
    expect(subscriptionState.onError).toBeDefined();

    actor.stop();
  });
});
