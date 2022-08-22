import { graphql ,fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchClusterResources(variables) {
  const vars = { k8scontextIDs : variables }

  const query = graphql`
    query ClusterResourcesQuery($k8scontextIDs: [String!]) {
      clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs) {
        resources {
          kind
          count
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars)
}
