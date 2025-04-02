import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchMesheryOperatorStatus(variables) {
  const environment = createRelayEnvironment({});
  const vars = { connectionID: variables.connectionID };

  const OperatorStatusQueryNode = graphql`
    query OperatorStatusQuery($connectionID: String!) {
      operator: getOperatorStatus(connectionID: $connectionID) {
        status
        controller
        connectionID
      }
    }
  `;

  return fetchQuery(environment, OperatorStatusQueryNode, vars);
}
