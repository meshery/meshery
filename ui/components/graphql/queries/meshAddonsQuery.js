import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableAddons(variables) {
  const vars = {
    selector: variables.serviceMesh
  }

  const query = graphql`
        query meshAddonsQuery($meshType: MeshType){
            getAvailableAddons(selector: $meshType) {
              type
              status
            }
        }
    `;

  return fetchQuery(environment, query, vars)
}
