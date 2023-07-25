import { fetchQuery, graphql } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function fetchTelemetryCompsQuery(variables) {
  const environment = createRelayEnvironment({});

  const TelemetryComponentsQueryNode = graphql`
    query TelemetryComponentsQuery($contexts: [String!]) {
      telemetryComps: fetchTelemetryComponents(contexts: $contexts) {
          name
          spec
          status 
      }
    }`;
  return fetchQuery(environment, TelemetryComponentsQueryNode, variables);
}