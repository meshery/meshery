import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function NatsStatusQuery(vars) {
  const environment = createRelayEnvironment({});

  const NatsStatusQueryNode = graphql`
    query NatsStatusQuery($connectionID: String!) {
      controller: getNatsStatus(connectionID: $connectionID) {
        name
        version
        status
      }
    }
  `;

  return fetchQuery(environment, NatsStatusQueryNode, vars);
}
