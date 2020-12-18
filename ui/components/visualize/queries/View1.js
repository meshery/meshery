import {fetchQuery, graphql} from 'relay-runtime';
import environment from '../../../lib/relayEnvironment';

export default function fetchView1(variables) {
  let query = graphql`
        query View1Query($nodeID: ID, $showClusters: Boolean!, $showNodes: Boolean!, $showPods: Boolean!) {
          cluster {
              id @include(if: $showClusters)
            clusterNodes(nodeid: $nodeID) {
                id  @include(if: $showNodes)
              pods {
                id @include(if: $showPods)
                name @include(if: $showPods)
              }
            }
          }
        }
      `;

  const data = fetchQuery(environment, query, variables)
  return data
}

// export const fetchLayout1