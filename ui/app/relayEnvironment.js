import {
  Environment,
  Network,
  Observable,
  RecordSource,
  Store,
} from "relay-runtime";
import { SubscriptionClient } from "subscriptions-transport-ws";
import { promisifiedDataFetch } from "../lib/data-fetch";

function fetchQuery(operation, variables) {
  return promisifiedDataFetch("/api/system/graphql/query", {
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    method: "POST",
    body: JSON.stringify({
      query: operation.text,
      variables,
    }),
  });
}

// let subscriptionClient;

// if (typeof window !== "undefined") {
//   subscriptionClient = new SubscriptionClient(
//     "ws://localhost:3000/api/system/graphql/query"
//   );
// }

function setupSubscription(config, variables) {
  const subscriptionClient = new SubscriptionClient(
    "ws://localhost:3000/api/system/graphql/query", // use a utitlity function to get the host dynamically
    { reconnect: true }
  );
  const query = config.text;

  const subscribeObservable = subscriptionClient.request({
    query,
    variables,
    operationName: config.name,
  });

  return Observable.from(subscribeObservable);
}

const environment = new Environment({
  network: Network.create(fetchQuery, setupSubscription),
  store: new Store(new RecordSource()),
});

export default environment;
