import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";


export default function resetDatabase(variables) {
  const vars = { selector : variables.selector, k8scontextID : variables.k8scontextID };

  const query = graphql`
        query ResetDatabaseQuery($selector: ReSyncActions!, $k8scontextID: String!) {
           resetStatus: resyncCluster(selector: $selector, k8scontextID: $k8scontextID) 
        }
    `;

  return fetchQuery(environment, query, vars)
}
