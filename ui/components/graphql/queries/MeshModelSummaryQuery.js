import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

export default function fetchMeshModelSummary(selector) {
  const environment = createRelayEnvironment({});

  const vars = { selector: selector };

  const MeshModelSummaryQueryNode = graphql`
    query MeshModelSummaryQuery($selector: MeshModelSummarySelector!) {
      meshmodelSummary: getMeshModelSummary(selector: $selector) {
        components {
          name
          count
        }
        relationships {
          name
          count
        }
      }
    }
  `;

  return fetchQuery(environment, MeshModelSummaryQueryNode, vars);
}
