import {
  Environment,
  Network,
  Observable,
  RecordSource,
  Store
} from 'relay-runtime';
import { createClient } from "graphql-ws";
import { promisifiedDataFetch } from "./data-fetch";

function fetchQuery(operation, variables) {
  return promisifiedDataFetch("/api/system/graphql/query", {
    headers : {
      "Content-Type" : "application/json",
    },
    credentials : "include",
    method : "POST",
    body : JSON.stringify({
      query : operation.text,
      variables,
    }),
  });
}

export let subscriptionClient;

if (typeof window !== 'undefined') {
  const isWss = window.location.protocol === "https:";
  const wsProtocol = isWss ? "wss://" : "ws://"
  subscriptionClient = createClient({
      url: wsProtocol + window.location.host + "/api/system/qraphql/query",
  })
}

function setupSubscription(
  operation,
  variables,
) {
  return Observable.create((sink) => {
      if (!operation.text) {
          return sink.error(new Error('Operation text cannot be empty'));
      }
      return subscriptionClient.subscribe(
          {
              operationName: operation.name,
              query: operation.text,
              variables,
          },
          sink,
      );
  });
}

const environment = new Environment({
  network: Network.create(fetchQuery, setupSubscription, setupSubscription),
  store: new Store(new RecordSource()),
});

export default environment;