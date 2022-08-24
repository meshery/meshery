import { graphql ,fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchClusterResources(ctxIDs, namespaces) {
  const vars = {
    k8scontextIDs : ctxIDs,
    namespaces : namespaces
  }

  const query = graphql`
    query ClusterResourcesQuery($k8scontextIDs: [String!], $namespaces: [String!]) {
      clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs, namespaces: $namespaces) {
        resources {
          kind
          count
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars)
}
