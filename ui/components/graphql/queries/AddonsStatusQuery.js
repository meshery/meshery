import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableAddons(variables) {
  const vars = {
    selector: variables.serviceMesh,
  };

  const query = graphql`
    query AddonsStatusQuery($selector: MeshType) {
      addonsState: getAvailableAddons(selector: $selector) {
        name
        owner
        endpoint
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
