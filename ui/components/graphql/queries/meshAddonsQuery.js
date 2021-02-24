import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableAddons(variables) {
  const vars = {
    selector: variables.serviceMesh,
  };

  const query = graphql`
    query meshAddonsQuery($selector: MeshType) {
      addons: getAvailableAddons(selector: $selector) {
        type
        status
        config {
        serviceName
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
