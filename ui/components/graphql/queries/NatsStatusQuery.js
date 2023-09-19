import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function NatsStatusQuery(vars) {
  const environment = createRelayEnvironment({});

  const NatsStatusQueryNode = graphql`
    query NatsStatusQuery($k8scontextID: String!) {
      controller: getNatsStatus(k8scontextID: $k8scontextID) {
        name
        version
        status
      }
    }
  `;

  return fetchQuery(environment, NatsStatusQueryNode, vars);
}
