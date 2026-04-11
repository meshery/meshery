import { SEVERITY_TO_NOTIFICATION_TYPE_MAPPING } from '@/components/NotificationCenter/constants';
import subscribeEvents from '@/components/graphql/subscriptions/EventsSubscription';
import { emit, fromCallback, setup, spawnChild, stopChild } from 'xstate';
import { store } from '../store';
import { pushEvent } from '@/store/slices/events';
import { api as mesheryApi } from '../rtk-query';
import { PROVIDER_TAGS } from '@/rtk-query/notificationCenter';
export const OPERATION_CENTER_EVENTS = {
  EVENT_RECEIVED_FROM_SERVER: 'EVENT_RECEIVED_FROM_SERVER',
  ERROR_OCCURRED_IN_SUBSCRIPTION: 'ERROR_OCCURRED_IN_SUBSCRIPTION',
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
        if (!result.event) {
          console.error('Invalid event received', result);
          return;
        }

        const event = {
          ...result.event,
          user_id: result.event.userID,
          system_id: result.event.systemID,
          updated_at: result.event.updatedAt,
          created_at: result.event.createdAt,
          deleted_at: result.event.deletedAt,
          operation_id: result.event.operationID,
        };
        sendBack(events.eventReceivedFromServer(event));
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
      store.dispatch(pushEvent(event.data.event));
    },
    invalidateRtk: () => {
      store.dispatch(mesheryApi.util.invalidateTags([PROVIDER_TAGS.EVENT]));
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
