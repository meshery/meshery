import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

/**
 *
 * @param {{k8scontextID: string}} variable
 * @returns
 */
export default function deployMeshSync(variable) {

  const query = graphql`
    query DeployMeshSyncQuery($k8scontextID: String!) {
    deployMeshsync (k8scontextID: $k8scontextID)
    }
  `;

  return fetchQuery(environment, query, variable);
}
