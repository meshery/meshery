import subscribeMeshSyncStatusEvents from "../graphql/subscriptions/MeshSyncStatusSubscription";
import { isMeshSyncSubscriptionDataUpdated } from "./comparatorFns";
import { mergeMeshSyncSubscription } from "./mergeFns";

export const MESHSYNC_EVENT_SUBSCRIPTION = "MESHSYNC_EVENT_SUBSCRIPTION";
export const OPERATOR_EVENT_SUBSCRIPTION = "OPERATOR_EVENT_SUBSCRIPTION";

export const fnMapping = {
  MESHSYNC_EVENT_SUBSCRIPTION : {
    eventName : "listenToMeshSyncEvents",
    comparatorFn : isMeshSyncSubscriptionDataUpdated,
    subscriptionFn : subscribeMeshSyncStatusEvents,
    mergeFn : mergeMeshSyncSubscription
  },
}