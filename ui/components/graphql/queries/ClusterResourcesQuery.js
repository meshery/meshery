import { graphql ,fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchClusterResources(ctxIDs, namespace) {
  const vars = {
    k8scontextIDs : ctxIDs,
    namespace : namespace
  }

  const query = graphql`
    query ClusterResourcesQuery($k8scontextIDs: [String!], $namespace: String!) {
      clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs, namespace: $namespace) {
        resources {
          kind
          count
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars)
}
