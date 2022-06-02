import { graphql, fetchQuery } from "react-relay";
import environment from "../../../lib/relayEnvironment";

export default function fetchDataPlanes(variables) {
  const vars = {
    filter : variables,
  };

  const query = graphql`
    query DataPlanesQuery($filter: ServiceMeshFilter) {
      dataPlanesState: getDataPlanes(filter: $filter) {
        name
        proxies {
          controlPlaneMemberName
          containerName
          image
          status {
            containerStatusName
            image
            state
            lastState
            ready
            restartCount
            started
            imageID
            containerID
          }
          ports {
            name
            containerPort
            protocol
          }
          resources
        }
      }
    }
  `;

  return fetchQuery(environment, query, vars);
}
