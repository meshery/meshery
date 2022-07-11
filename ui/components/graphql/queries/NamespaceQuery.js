import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableNamespaces(vars) {
  const query = graphql`
        query NamespaceQuery($k8sClusterIDs: [String!]) {
          namespaces: getAvailableNamespaces(k8sClusterIDs: $k8sClusterIDs) {
            namespace
          }
        }
    `;

  return fetchQuery(environment, query, vars);
}
