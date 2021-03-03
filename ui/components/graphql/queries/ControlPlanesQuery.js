import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchControlPlanes(variables) {
  const vars = {
    filter: {
      "type": variables.serviceMesh
    }
  };

  const query = graphql`
    query ControlPlanesQuery($filter: ControlPlaneFilter) {
      controlPlanes: getControlPlanes(filter: $filter) {
        name
        members {
          version
          component
          namespace
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
