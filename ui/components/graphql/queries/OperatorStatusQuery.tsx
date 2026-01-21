import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface OperatorStatusVariables {
  connectionID: string;
}

export default function fetchMesheryOperatorStatus(variables: OperatorStatusVariables) {
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
