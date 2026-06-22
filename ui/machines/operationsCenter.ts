import { SEVERITY_TO_NOTIFICATION_TYPE_MAPPING } from '@/components/layout/NotificationCenter/constants';
import subscribeEvents from '@/graphql/subscriptions/EventsSubscription';
import { emit, fromCallback, setup, spawnChild, stopChild } from 'xstate';
import { store } from '../store';
import { pushEvent, pushEvents } from '@/store/slices/events';
import { api as mesheryApi } from '../rtk-query';
import { PROVIDER_TAGS } from '@/rtk-query/notificationCenter';

const SEVERITY_RANK = {
  error: 4,
  warning: 3,
  success: 2,
  informational: 1,
  info: 1,
};

export const OPERATION_CENTER_EVENTS = {
  EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
  ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
};

const normalizeServerEvents = (payload) => {
  if (!payload) {
    return [];
  }
  if (Array.isArray(payload)) {
    return payload.filter(Boolean);
  }

  const hasEventEnvelope =
    Object.prototype.hasOwnProperty.call(payload, 'events') ||
    Object.prototype.hasOwnProperty.call(payload, 'event');
  const eventsPayload = hasEventEnvelope ? (payload.events ?? payload.event) : payload;
  if (Array.isArray(eventsPayload)) {
    return eventsPayload.filter(Boolean);
  }
  return eventsPayload && Object.keys(eventsPayload).length > 0 ? [eventsPayload] : [];
};

const getRepresentativeEvent = (events) => {
  return events.reduce((selected, event) => {
    const selectedRank = SEVERITY_RANK[selected?.severity] || 0;
    const eventRank = SEVERITY_RANK[event?.severity] || 0;
    if (eventRank > selectedRank) {
      return event;
    }
    if (eventRank === selectedRank && event?.createdAt && selected?.createdAt) {
      return event.createdAt > selected.createdAt ? event : selected;
    }
    return selected;
  }, events[0]);
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
};

const subscriptionActor = fromCallback(({ sendBack }) => {
  const subscription = subscribeEvents(
    (result) => {
      try {
        const receivedEvents = normalizeServerEvents(result);
        if (receivedEvents.length === 0) {
          console.error('Invalid event received', result);
          return;
        }

        // GraphQL Event payload (canonical camelCase) is consumed as-is by
        // downstream UI; meshkit Event JSON tags also flipped to camelCase in
        // v1.0.7 so the REST /api/system/events list endpoint matches.
        sendBack(
          events.eventReceivedFromServer(
            receivedEvents.length === 1 ? receivedEvents[0] : receivedEvents,
          ),
        );
      } catch (error) {
        console.error('[operationsCenter] An error occurred in processing event', error);
      }
    },
    (error) => {
      console.error('[operationsCenter] An error occurred in subscription to events', error);
      sendBack(events.errorOccurredInSubscription(error));
    },
  );

  () => {
    subscription.dispose();
  };
});

export const operationsCenterActor = setup({
  actions: {
    spawnSubscriptionActor: spawnChild(subscriptionActor, {
      id: 'subscriptionActor',
    }),
    storeInRedux: ({ event }) => {
      const receivedEvents = normalizeServerEvents(event.data);
      if (receivedEvents.length === 1) {
        store.dispatch(pushEvent(receivedEvents[0]));
        return;
      }
      if (receivedEvents.length > 1) {
        store.dispatch(pushEvents(receivedEvents));
      }
    },
    invalidateRtk: () => {
      store.dispatch(mesheryApi.util.invalidateTags([PROVIDER_TAGS.EVENT]));
    },
    notifyUI: ({ context, event }) => {
      const receivedEvents = normalizeServerEvents(event.data);
      if (receivedEvents.length === 0) {
        return;
      }
      const validatedEvent = getRepresentativeEvent(receivedEvents);
      context.notify({
        message:
          receivedEvents.length === 1
            ? validatedEvent.description
            : `${receivedEvents.length} events received`,
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
      },
    },
  },
});
