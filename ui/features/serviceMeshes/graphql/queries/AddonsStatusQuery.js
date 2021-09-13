import environment from "app/relayEnvironment";
import { graphql, fetchQuery } from "react-relay";

export default function fetchAvailableAddons(variables) {
  const vars = { selector : variables.serviceMesh, };

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
