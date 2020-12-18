import {fetchQuery, graphql} from 'relay-runtime';
import environment from '../../../lib/relayEnvironment';

export default function fetchView2(variables) {
  let query = graphql`
        query View2Query($namespaceID: ID, $deploymentID: ID,$showClusters: Boolean!, $showNamespaces: Boolean!, $showDeployments: Boolean!, $showPods: Boolean!,) {
          cluster {
            id @include(if: $showClusters)
            namespaces(namespaceid: $namespaceID) {
              id @include(if: $showNamespaces)
              deployments(deploymentid: $deploymentID) {
                id @include(if: $showDeployments)
                pods {
                  id @include(if: $showPods)
                  name @include(if: $showPods)
                }
              }
            }
          }
        }
      `;

  const data = fetchQuery(environment, query, variables)
  return data
}