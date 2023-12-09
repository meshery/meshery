import subscribeMesheryControllersStatus from '../graphql/subscriptions/MesheryControllersStatusSubscription';
import subscribeOperatorStatusEvents from '../graphql/subscriptions/OperatorStatusSubscription';
import {
  isMesheryControllerStateSubscriptionDataUpdated,
  isOperatorStateSubscriptionDataUpdated,
} from './comparatorFns';
import { mergeMesheryController, mergeOperatorStateSubscription } from './mergeFns';

export const MESHSYNC_EVENT_SUBSCRIPTION = 'MESHSYNC_EVENT_SUBSCRIPTION';
export const OPERATOR_EVENT_SUBSCRIPTION = 'OPERATOR_EVENT_SUBSCRIPTION';
export const MESHERY_CONTROLLER_SUBSCRIPTION = 'MESHERY_CONTROLLER_SUBSCRIPTION';

export const fnMapping = {
  // MESHSYNC_EVENT_SUBSCRIPTION : {
  //   eventName : "listenToMeshSyncEvents",
  //   comparatorFn : isMeshSyncSubscriptionDataUpdated,
  //   subscriptionFn : subscribeMeshSyncStatusEvents,
  //   mergeFn : mergeMeshSyncSubscription
  // },
  MESHERY_CONTROLLER_SUBSCRIPTION: {
    eventName: 'subscribeMesheryControllersStatus',
    subscriptionFn: subscribeMesheryControllersStatus,
    mergeFn: mergeMesheryController,
    comparatorFn: isMesheryControllerStateSubscriptionDataUpdated,
  },
  OPERATOR_EVENT_SUBSCRIPTION: {
    eventName: 'operator',
    comparatorFn: isOperatorStateSubscriptionDataUpdated,
    subscriptionFn: subscribeOperatorStatusEvents,
    mergeFn: mergeOperatorStateSubscription,
  },
};

export function isControllerObjectEqual(oldController, newController) {
  return (
    newController.contextId === oldController.contextId &&
    newController.controller === oldController.controller
  );
}
