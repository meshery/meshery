import { Environment, Network, RecordSource, Store } from "relay-runtime";
import { promisifiedDataFetch } from "../../lib/data-fetch";

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

const environment = new Environment({
  network: Network.create(fetchQuery),
  store: new Store(new RecordSource()),
});

export default environment;
