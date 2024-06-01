import { enqueueActions, sendTo } from 'xstate';

export const sendToActor = (actor, event) =>
  sendTo(({ system }) => {
    console.log('sendToActor', actor, event);
    return system.get(actor);
  }, event);

export const sendToActors = (actorSystemIds, eventCreator) =>
  enqueueActions(({ context, event, enqueue, system, self }, params) => {
    actorSystemIds.forEach((actorSystemId) => {
      const actor = system.get(actorSystemId);
      if (!actor) {
        console.log('actor not found --sendToActors', actorSystemId);
      }
      enqueue.sendTo(actor, eventCreator({ context, event, system, self }, params));
    });
  });

export const forwardToActors = (actorSystemIds) =>
  enqueueActions(({ event, enqueue, system }) => {
    actorSystemIds.forEach((actorSystemId) => {
      const actor = system.get(actorSystemId);
      enqueue.sendTo(actor, event);
    });
  });

export const deadLetter = (event) => ({ type: 'DEAD_LETTER', event });

export const reply = (eventFn) => sendTo(({ context }) => context.returnAddress, eventFn);

export const XSTATE_DEBUG_EVENT = 'XSTATE_DEBUG_EVENT';
