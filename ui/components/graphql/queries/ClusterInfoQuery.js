import { graphql ,fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchClusterInfo(variables) {
  const vars = { k8scontextIDs : variables }

  const query = graphql`
    query ClusterInfoQuery($k8scontextIDs: [String!]) {
      clusterInfo: getClusterInfo(k8scontextIDs: $k8scontextIDs) {
        resources {
          kind
          number
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars)
}
