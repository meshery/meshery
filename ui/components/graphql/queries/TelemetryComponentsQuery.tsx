import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface TelemetryQueryVariables {
  contexts?: string[];
}

export default function fetchTelemetryCompsQuery(variables: TelemetryQueryVariables) {
  const environment = createRelayEnvironment({});

  const TelemetryComponentsQueryNode = graphql`
    query TelemetryComponentsQuery($contexts: [String!]) {
      telemetryComps: fetchTelemetryComponents(contexts: $contexts) {
        name
        spec
        status
      }
    }
  `;
  return fetchQuery(environment, TelemetryComponentsQueryNode, variables);
}
