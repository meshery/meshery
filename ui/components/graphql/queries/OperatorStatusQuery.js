import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchMesheryOperatorStatus(variables) {
  const environment = createRelayEnvironment({});
  const vars = { k8scontextID: variables.k8scontextID };

  const OperatorStatusQueryNode = graphql`
    query OperatorStatusQuery($k8scontextID: String!) {
      operator: getOperatorStatus(k8scontextID: $k8scontextID) {
        status
        version
        controllers {
          name
          version
          status
        }
        error {
          code
          description
        }
      }
    }
  `;

  return fetchQuery(environment, OperatorStatusQueryNode, vars);
}
