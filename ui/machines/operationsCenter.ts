import { SEVERITY_TO_NOTIFICATION_TYPE_MAPPING } from '@/components/layout/NotificationCenter/constants';
import { subscribeToEvents } from 'lib/eventsSubscription';
import { assign, emit, fromCallback, setup, spawnChild, stopChild } from 'xstate';
import { store } from '../store';
import { loadEvents, pushEvent } from '@/store/slices/events';
import { api as mesheryApi } from '../rtk-query';
import { notificationCenterApi, PROVIDER_TAGS } from '@/rtk-query/notificationCenter';
export const OPERATION_CENTER_EVENTS = {
  EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
  ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
  SUBSCRIPTION_OPENED: 'SUBSCRIPTION_OPENED',
};
const events = {
  eventReceivedFromServer: (event) => ({
    type: OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER,
    data: {
      event,
    },
  }),
  errorOccurredInSubscription: (error) => ({
    type: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
    data: {
      error,
    },
  }),
  subscriptionOpened: () => ({ type: OPERATION_CENTER_EVENTS.SUBSCRIPTION_OPENED }),
};

const subscriptionActor = fromCallback(({ sendBack }) => {
  const subscription = subscribeToEvents(
    (event) => {
      try {
        if (!event) {
          console.error('Invalid event received', event);
          return;
        }

        // The SSE stream delivers the raw meshkit Event JSON (camelCase), the
        // same shape the REST /api/system/events list endpoint returns, so it
        // is consumed as-is by the downstream UI.
        sendBack(events.eventReceivedFromServer(event));
      } catch (error) {
        console.error('[operationsCenter] An error occurred in processing event', error);
      }
    },
    (error) => {
      console.error('[operationsCenter] An error occurred in subscription to events', error);
      sendBack(events.errorOccurredInSubscription(error));
    },
    () => sendBack(events.subscriptionOpened()),
  );

  return () => {
    subscription.dispose();
  };
});

// Events received from the stream while a list reload is in flight, one buffer
// per reload. The reloaded page replaces the list, so these are re-added after.
const reloadBuffers = new Set<unknown[]>();

export const operationsCenterActor = setup({
  actions: {
    spawnSubscriptionActor: spawnChild(subscriptionActor, {
      id: 'subscriptionActor',
    }),
    storeInRedux: ({ event }) => {
      store.dispatch(pushEvent(event.data.event));
      reloadBuffers.forEach((buffer) => buffer.push(event.data.event));
    },
    invalidateRtk: () => {
      store.dispatch(mesheryApi.util.invalidateTags([PROVIDER_TAGS.EVENT]));
    },
    // The notification list renders the events slice, which only loadEvents
    // fills (RTK cache invalidation doesn't reach it). Reload the first page of
    // the current view, the same way the notification center does.
    reloadEvents: () => {
      const fetchEvents = async (args) => {
        const request = store.dispatch(
          notificationCenterApi.endpoints.getEvents.initiate(args, { forceRefetch: true }),
        );
        try {
          const result = await request;
          // A failed request resolves with `error` instead of throwing. Throw so
          // loadEvents keeps the current list rather than treating it as empty.
          if (result.error) throw result.error;
          return result;
        } finally {
          request.unsubscribe();
        }
      };
      const received: unknown[] = [];
      reloadBuffers.add(received);
      Promise.resolve(
        store.dispatch(loadEvents(fetchEvents, 0, store.getState().events.current_view.filters)),
      ).finally(() => {
        reloadBuffers.delete(received);
        // Re-add events the stream delivered while the page was loading; the
        // slice is keyed by id, so ones already in the page are skipped.
        received.forEach((event) => store.dispatch(pushEvent(event)));
      });
    },
    notifyUI: ({ context, event }) => {
      const validatedEvent = event.data.event;
      context.notify({
        message: validatedEvent.description,
        event_type: SEVERITY_TO_NOTIFICATION_TYPE_MAPPING[validatedEvent.severity],
        id: validatedEvent.id,
        showInNotificationCenter: true,
      });
    },
    emitBack: emit(({ event }) => event),
  },
}).createMachine({
  id: 'operationsCenter',
  initial: 'init',
  context: ({ input }) => ({
    notify: input.notify,
    // Whether the stream has connected before. The first connection needs no
    // reload: the notification center loads the list itself on mount.
    hasConnected: false,
  }),

  states: {
    init: {
      entry: 'spawnSubscriptionActor',
      always: 'idle',
    },
    idle: {
      on: {
        [OPERATION_CENTER_EVENTS.EVENT_RECEIVED_FROM_SERVER]: {
          actions: ['storeInRedux', 'notifyUI', 'emitBack'],
        },
        [OPERATION_CENTER_EVENTS.ERROR_OCCURRED_IN_SUBSCRIPTION]: {
          actions: [stopChild('subscriptionActor'), 'spawnSubscriptionActor'],
        },
        // Events raised while the stream was down are never replayed. On every
        // reconnection, reload the persisted list (and the unread summary) once
        // the server is listening again, so nothing raised in between is lost.
        [OPERATION_CENTER_EVENTS.SUBSCRIPTION_OPENED]: [
          {
            guard: ({ context }) => context.hasConnected,
            actions: ['invalidateRtk', 'reloadEvents'],
          },
          { actions: assign({ hasConnected: true }) },
        ],
      },
    },
  },
});
