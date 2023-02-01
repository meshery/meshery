import { fetchQuery, graphql } from "relay-runtime";
import environment from "../../../lib/relayEnvironment";

export default function fetchTelemetryCompsQuery(variables) {
  const query = graphql`
    query TelemetryComponentsQuery($contexts: [String!]) {
      telemetryComps: fetchTelemetryComponents(contexts: $contexts) {
          name
          spec
          status 
      }
    }`;
  return fetchQuery(environment, query, variables);
}