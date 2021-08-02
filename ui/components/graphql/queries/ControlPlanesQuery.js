import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchControlPlanes(variables) {
  const vars = {
    filter: {
      "type": variables.type
    }
  };

  const query = graphql`
    query ControlPlanesQuery($filter: ControlPlaneFilter) {
      controlPlanesState: getControlPlanes(filter: $filter) {
        name
        members {
          name
          version
          component
          namespace
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
