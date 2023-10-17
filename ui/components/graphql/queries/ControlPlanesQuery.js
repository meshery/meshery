import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchControlPlanes(variables) {
  const environment = createRelayEnvironment({});
  const vars = { filter: variables };

  const ControlPlanesQueryNode = graphql`
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

  return fetchQuery(environment, ControlPlanesQueryNode, vars);
}
