import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchControlPlanes(variables) {
  const vars = { filter : variables };

  const query = graphql`
    query ControlPlanesQuery($filter: ServiceMeshFilter) {
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
