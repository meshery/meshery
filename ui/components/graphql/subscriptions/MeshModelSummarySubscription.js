import { graphql, requestSubscription } from "react-relay";
import { createRelayEnvironment } from "../../../lib/relayEnvironment";

// not in use
export const meshmodelSummarySubscription = graphql`
  subscription MeshModelSummarySubscription($selector: MeshModelSummarySelector!) {
    meshmodelSummary: subscribeMeshModelSummary(selector: $selector) {
      components {
        name,
        count
      },
      relationships {
        name,
        count
      }
    }
  }
`;

export default function subscribeClusterResources(dataCB, variables) {
  const environment = createRelayEnvironment({});
  return requestSubscription(environment, {
    subscription : meshmodelSummarySubscription,
    variables : variables,
    onNext : dataCB,
    onError : (error) => console.log(`MeshModel Subscription error:`, error),
  });
}
