import { fetchQuery, graphql } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

export default function MeshsyncStatusQuery(vars = { connectionID : vars.connectionID }) {
  const environment = createRelayEnvironment({});

  const MeshsyncStatusQueryNode = graphql`
    query MeshsyncStatusQuery($connectionID: String!)  {
        controller: getMeshsyncStatus(connectionID: $connectionID){
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, MeshsyncStatusQueryNode, vars);
}
