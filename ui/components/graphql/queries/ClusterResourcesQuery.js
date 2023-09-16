import { fetchQuery, graphql } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function fetchClusterResources(ctxIDs, namespace) {
  const environment = createRelayEnvironment({});
  const vars = { k8scontextIDs: ctxIDs, namespace };

  return fetchQuery(environment, graphql`
    query ClusterResourcesQuery($k8scontextIDs: [String!], $namespace: String!) {
      clusterResources: getClusterResources(k8scontextIDs: $k8scontextIDs, namespace: $namespace) {
        resources { kind count }
      }
    }
  `, vars);
}
