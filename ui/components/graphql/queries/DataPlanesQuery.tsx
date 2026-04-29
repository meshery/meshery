import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchDataPlanes(variables) {
  const environment = createRelayEnvironment({});
  const vars = {
    filter: variables,
  };

  const DataPlanesQueryNode = graphql`
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

  return fetchQuery(environment, DataPlanesQueryNode, vars);
}
