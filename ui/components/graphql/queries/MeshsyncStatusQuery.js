import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function MeshsyncStatusQuery(variables) {
  const vars = { k8scontextID : variables.k8scontextID };

  const query = graphql`
    query MeshsyncStatusQuery($k8scontextID: String!)  {
        controller: getMeshsyncStatus(k8scontextID: $k8scontextID){
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, query, vars);
}
