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
  loadEvents: vi.fn((fetch: unknown, page: number, filters: unknown) => ({
    type: 'events/loadEvents',
    fetch,
    page,
    filters,
  })),
  initiate: vi.fn(),
  currentFilters: { status: 'unread' } as Record<string, unknown>,
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
  store: {
    dispatch: hoisted.dispatch,
    getState: () => ({ events: { current_view: { filters: hoisted.currentFilters } } }),
  },
}));

vi.mock('@/store/slices/events', () => ({
  pushEvent: (p: unknown) => hoisted.pushEvent(p),
  loadEvents: (fetch: unknown, page: number, filters: unknown) =>
    hoisted.loadEvents(fetch, page, filters),
}));

vi.mock('../../rtk-query', () => ({
  api: { util: { invalidateTags: vi.fn((tags) => ({ type: 'api/invalidate', tags })) } },
}));

vi.mock('@/rtk-query/notificationCenter', () => ({
  PROVIDER_TAGS: { EVENT: 'event' }, // matches rtk-query/notificationCenter
  notificationCenterApi: {
    endpoints: { getEvents: { initiate: (...args: unknown[]) => hoisted.initiate(...args) } },
  },
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
    dispatch.mockReset();
    hoisted.loadEvents.mockClear();
    hoisted.initiate.mockReset();
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

  it('refetches persisted events after re-subscribing, so events missed while down appear', () => {
    const { actor } = startActor();
    dispatch.mockClear();

    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('stream closed') },
    });

    expect(subscriptionState.dispose).toHaveBeenCalled();
    expect(dispatch).toHaveBeenCalledWith({ type: 'api/invalidate', tags: ['event'] });

    actor.stop();
  });

  it('reloads the first page of the current notification view after re-subscribing', () => {
    const { actor } = startActor();
    hoisted.loadEvents.mockClear();

    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('stream closed') },
    });

    expect(hoisted.loadEvents).toHaveBeenCalledWith(expect.any(Function), 0, {
      status: 'unread',
    });
    expect(dispatch).toHaveBeenCalledWith(
      expect.objectContaining({ type: 'events/loadEvents', page: 0 }),
    );
    actor.stop();
  });

  it('fetches the reloaded page fresh and releases the query subscription', async () => {
    const { actor } = startActor();
    hoisted.loadEvents.mockClear();
    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('stream closed') },
    });
    const fetchEvents = hoisted.loadEvents.mock.calls[0][0] as (args: unknown) => Promise<unknown>;

    const unsubscribe = vi.fn();
    const request = Object.assign(Promise.resolve({ data: { events: [] } }), { unsubscribe });
    hoisted.initiate.mockReturnValue({ type: 'api/initiate' });
    dispatch.mockReturnValueOnce(request);

    await expect(fetchEvents({ page: 0, filters: {} })).resolves.toEqual({ data: { events: [] } });
    expect(hoisted.initiate).toHaveBeenCalledWith({ page: 0, filters: {} }, { forceRefetch: true });
    expect(unsubscribe).toHaveBeenCalled();
    actor.stop();
  });

  const failSubscription = (actor) =>
    actor.send({
      type: OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION,
      data: { error: new Error('stream closed') },
    });

  it('rejects a failed reload so loadEvents keeps the current list', async () => {
    const { actor } = startActor();
    failSubscription(actor);
    const fetchEvents = hoisted.loadEvents.mock.calls[0][0] as (args: unknown) => Promise<unknown>;

    const unsubscribe = vi.fn();
    const failure = { status: 500, data: 'boom' };
    dispatch.mockReturnValueOnce(
      Object.assign(Promise.resolve({ error: failure }), { unsubscribe }),
    );

    await expect(fetchEvents({ page: 0, filters: {} })).rejects.toBe(failure);
    expect(unsubscribe).toHaveBeenCalled();
    actor.stop();
  });

  it('re-adds events that arrive while the reload is in flight', async () => {
    let finishReload: () => void = () => {};
    dispatch.mockImplementation((action) =>
      action?.type === 'events/loadEvents'
        ? new Promise<void>((resolve) => {
            finishReload = resolve;
          })
        : action,
    );
    const { actor } = startActor();
    failSubscription(actor);

    const live = { id: 'during-reload', description: 'd', severity: 'informational' };
    actor.send({ type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER, data: { event: live } });
    expect(pushEvent).toHaveBeenCalledTimes(1);

    finishReload();
    await vi.waitFor(() => expect(pushEvent).toHaveBeenCalledTimes(2));
    expect(pushEvent).toHaveBeenLastCalledWith(live);
    actor.stop();
  });

  it('does not re-add events that arrive after the reload finished', async () => {
    const { actor } = startActor();
    failSubscription(actor);
    await Promise.resolve();

    const later = { id: 'after-reload', description: 'd', severity: 'informational' };
    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event: later },
    });
    await Promise.resolve();

    expect(pushEvent).toHaveBeenCalledTimes(1);
    actor.stop();
  });

  it('does not refetch events for ordinary received events', () => {
    const { actor } = startActor();
    dispatch.mockClear();

    actor.send({
      type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
      data: { event: { id: 'e', description: 'd', severity: 'informational' } },
    });

    expect(dispatch).not.toHaveBeenCalledWith(expect.objectContaining({ type: 'api/invalidate' }));
    expect(hoisted.loadEvents).not.toHaveBeenCalled();
    actor.stop();
  });
});
