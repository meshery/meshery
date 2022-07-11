import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchAvailableAddons(variables) {
  const vars = { filter : variables };

  const query = graphql`
    query AddonsStatusQuery($filter: ServiceMeshFilter) {
      addonsState: getAvailableAddons(filter: $filter) {
        name
        owner
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
