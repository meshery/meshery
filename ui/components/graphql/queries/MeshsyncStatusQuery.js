import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function MeshsyncStatusQuery() {

  const query = graphql`
    query MeshsyncStatusQuery {
        controller: getMeshsyncStatus{
            name
            version
            status
        }
    }
  `;

  return fetchQuery(environment, query);
}
