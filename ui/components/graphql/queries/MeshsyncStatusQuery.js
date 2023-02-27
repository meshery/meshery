import { fetchQuery, graphql } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function MeshsyncStatusQuery(vars = { k8scontextID : vars.k8scontextID }) {
  const environment = createRelayEnvironment({});

  const MeshsyncStatusQueryNode = graphql`
    query MeshsyncStatusQuery($k8scontextID: String!)  {
        controller: getMeshsyncStatus(k8scontextID: $k8scontextID){
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, MeshsyncStatusQueryNode, vars);
}
