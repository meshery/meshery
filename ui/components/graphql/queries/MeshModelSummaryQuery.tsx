import { fetchQuery, graphql } from 'react-relay';
import { createRelayEnvironment } from '../../../lib/relayEnvironment';

interface MeshModelSummarySelector {
  [key: string]: unknown;
}

// not in use
export default function fetchMeshModelSummary(selector: MeshModelSummarySelector) {
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
