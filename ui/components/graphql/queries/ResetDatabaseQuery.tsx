import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface ReSyncSelector {
  clearDB?: string;
  ReSync?: string;
  hardReset?: string;
}

interface ResetDatabaseVariables {
  selector: ReSyncSelector;
  k8scontextID: string;
}

export default function resetDatabase(variables: ResetDatabaseVariables) {
  const environment = createRelayEnvironment({});
  const vars = { selector: variables.selector, k8scontextID: variables.k8scontextID };

  const ResetDatabaseQueryNode = graphql`
    query ResetDatabaseQuery($selector: ReSyncActions!, $k8scontextID: String!) {
      resetStatus: resyncCluster(selector: $selector, k8scontextID: $k8scontextID)
    }
  `;

  return fetchQuery(environment, ResetDatabaseQueryNode, vars);
}
